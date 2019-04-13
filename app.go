package svcledger

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/pkg/errors"

	"app"
	"app/httputil"
	"app/websocket"
	"app/websocket_gorilla"
	"common"
	"common/flags"
	"svcledger/job"
	"svcledger/helpers"
	"svcledger/store"
	"svcledger/blockchain"
)

const (
	TypeName 							= "ledger"
	servicePrivateKeyLength				= 2048
	autoCloseInterval uint64 			= 30
	writeFinancialBlockInterval uint64	= 60
	syncTransactionInterval	uint64		= 3
)

var (
	blockchainEndPointUrl = flags.String("ledger-blockchain-endpoint-url", "http://13.70.22.206:8080/", "Blockchain API endpoint")
	priceAmountIncome = flags.Uint64("ledger-price-amount-income", 500, "Min price of traffic")
	priceAmountSpend = flags.Uint64("ledger-price-amount-spend", 500, "Max price of traffic")
	amountInitial = flags.Uint64("ledger-amount-initial", 100000, "Initial tokens amount on wallet")
	amountSpend = flags.Uint64("ledger-amount-spend", 10000, "Tokens in spend address after it opened")
	priceQuantumPower = flags.Uint("ledger-quantum-power", 25, "Price power of quantum to transfer")
	servicePrivateKeyFileName = flags.String("ledger-private-key-file-name", "", "Service private key")
	servicePublicKeyFileName = flags.String("ledger-public-key-file-name", "", "Service public key")
	incomeChannelLifeTime = flags.Int64("ledger-income-channel-lifetime", 30, "Income channel lifetime, min")
	spendChannelLifeTime = flags.Int64("ledger-spend-channel-lifetime", 30, "Spend channel lifetime, min")
	transferChannelLifeTime = flags.Int64("ledger-transfer-channel-lifetime", 5, "Transfer channel lifetime, min")
)

type Config struct {
 	*app.Config
 	BlockchainEndPointUrl		string
	Amounts						*store.Amounts
	ServicePrivateKeyFileName	string
	ServicePublicKeyFileName	string
	LifeTimes					*store.LifeTimes
}

type Service struct {
	address				string
	keyPair				helpers.KeyPair
	queries				*store.Queries
	silentSch			*gocron.Scheduler
	silentSchStop		chan bool
	waiters				store.Waiters
	writer				blockchain.Writer

	BlockchainClient	*blockchain.Client
	Config				*Config
	Hub					websocket.Huber
	SocketHandler		websocket.SocketHandler
	Ledger				*store.Ledger
}

func NewService(config *app.Config) (*Service, error) {
	svcConfig := &Config{
		Config: config,
		BlockchainEndPointUrl: *blockchainEndPointUrl,
		Amounts: &store.Amounts{
			Initial: store.TransactionAmountType(*amountInitial),
			Spend:  store.TransactionAmountType(*amountSpend),
			PriceIncome: store.TransactionAmountType(*priceAmountIncome),
			PriceSpend: store.TransactionAmountType(*priceAmountSpend),
			PriceQuantumPower: store.QuantumPowerType(*priceQuantumPower),
		},
		ServicePrivateKeyFileName: *servicePrivateKeyFileName,
		ServicePublicKeyFileName: *servicePublicKeyFileName,
		LifeTimes: &store.LifeTimes{
			IncomeChannel: int64(time.Duration(*incomeChannelLifeTime) * time.Minute / time.Millisecond),
			SpendChannel: int64(time.Duration(*spendChannelLifeTime) * time.Minute / time.Millisecond),
			TransferChannel: int64(time.Duration(*transferChannelLifeTime) * time.Minute / time.Millisecond),
		},
	}

	client, err := blockchain.NewClient(svcConfig.BlockchainEndPointUrl)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		BlockchainClient: client,
		Config: svcConfig,
		waiters: store.NewWaiters(),
		writer: blockchain.NewWriter(client),
	}

	return svc, nil
}

func (svc *Service) startKeys() error {
	keyPair, err := helpers.NewECSDAKeyPair()

	if err != nil {
		return err
	}

	pk, err := keyPair.PublicKey()
	if err != nil {
		return err
	} else {
		common.Log.PrintFull(common.Printf("%v", pk), )
	}
	sk, err := keyPair.PrivateKey()
	if err != nil {
		return err
	} else {
		common.Log.PrintFull(common.Printf("%v", sk), )
	}

	svc.keyPair = keyPair
	svc.address = helpers.CreateCommonAddress(pk)

	return nil
}

func (svc *Service) startTransactionLoader() error {
	trans := make(blockchain.TranItems, 0, 0)
	for {
		t, err := svc.BlockchainClient.GetTrans()

		if err != nil {
			return err
		}

		if len(t.Data) > 0 {
			trans = append(trans, t.Data...)
		}

		if t.ReadAll {
			break
		}
	}

	var err error
	if err = svc.waiters.Start(); err != nil {
		return err
	}
	if err = svc.writer.Start(); err != nil {
		return err
	}

	svc.Ledger, err = store.NewLedger(
		svc.address,
		svc.keyPair,
		trans,
		svc.Config.Amounts,
		svc.Config.LifeTimes,
		svc.waiters,
		svc.writer,
	)
	if err != nil {
		return err
	}


	return nil
}

func (svc *Service) startWebSocket() error {
	var err error

	svc.queries = store.NewQueries()

	handler, err := handlerFunc(
		svc.keyPair,
		svc.Ledger,
		svc.queries,
	)
	if err != nil {
		return err
	}

	//svc.SocketHandler = websocket_gobwas.NewSocketHandler(ln, handler, nil)
	router := httputil.NewEmptyRouter(svc.Config.IsDevelopment())

	svc.SocketHandler, err = websocket_gorilla.NewSocketHandler(
		svc.Config.ListenNetwork,
		svc.Config.Listen,
		handler,
		router,
		nil,
		nil,
		true,
	)
	if err != nil {
		return err
	}
	svc.SocketHandler.Execute()

	return nil
}

func (svc *Service) createSilentScheduler() (*gocron.Scheduler, chan bool) {
	scheduler := gocron.NewScheduler()

	scheduler.
		Every(autoCloseInterval).
		Seconds().
		Do(job.AutoCloseJob(svc.Ledger))

	scheduler.
		Every(writeFinancialBlockInterval).
		Seconds().
		Do(job.WriteFinancialBlock(svc.Ledger))

	scheduler.
		Every(syncTransactionInterval).
		Seconds().
		Do(job.Syncer(svc.Ledger, svc.BlockchainClient))

	schedulerStop := scheduler.Start()

	return scheduler, schedulerStop
}

func (svc *Service) Start() error {
	var err error

	err = svc.startKeys()
	if err != nil {
		return err
	}

	err = svc.startTransactionLoader()
	if err != nil {
		return err
	}

	if svc.silentSch == nil {
		svc.silentSch, svc.silentSchStop = svc.createSilentScheduler()
	} else {
		return errors.New("silent scheduler already started")
	}

	common.Log.PrintFull(
		common.Printf("try to init tokens on wallet"),
	)
	status, err := svc.Ledger.Init(svc.Config.Amounts.Initial)

	if status == store.InvalidTransactionStatus {
		return err
	} else {
		common.Log.PrintFull(
			common.Printf("init tokens on wallet result: %v", status),
		)
	}

	err = svc.startWebSocket()
	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) Stop() error {
	var err error

	if svc.SocketHandler != nil {
		svc.SocketHandler.Shutdown()
	}

	if svc.writer != nil {
		svc.writer.Stop()
	}

	if svc.waiters != nil {
		svc.waiters.Stop()
	}

	if svc.queries != nil {
		svc.queries.Stop()
	}

	if svc.silentSchStop != nil {
		svc.silentSchStop <- true
	}

	return err
}

func (svc *Service) getRsaPrivateKey() (*rsa.PrivateKey, error) {
	var privateKey *rsa.PrivateKey
	var err error

	if svc.Config.ServicePrivateKeyFileName == "" ||
		(svc.Config.EnvironmentName == "development" && svc.Config.Listen != ":8085") {
		privateKey, err = rsa.GenerateKey(rand.Reader, servicePrivateKeyLength)

		pemdata := pem.EncodeToMemory(
			&pem.Block{
				Type: "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
			},
		)

		common.Log.PrintFull(common.Printf("%v", string(pemdata[:]), ), )

		if err != nil {
			return nil, err
		}
	} else {
		key, err := ioutil.ReadFile(svc.Config.ServicePrivateKeyFileName)

		if err != nil {
			return nil, err
		}

		block, _ := pem.Decode(key)
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)

		if err != nil {
			return nil, err
		}
	}

	return privateKey, nil
}

func (svc *Service) getRSAPublicKey(privateKey *rsa.PrivateKey) (crypto.PublicKey, error) {
	var publicKey crypto.PublicKey
	if svc.Config.ServicePublicKeyFileName == "" ||
		(svc.Config.EnvironmentName == "development" && svc.Config.Listen != ":8085") {
		publicKey = privateKey.Public()

		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)

		if err != nil {
			return nil, err
		}

		pemdata := pem.EncodeToMemory(
			&pem.Block{
				Type: "PUBLIC KEY",
				Bytes: publicKeyBytes,
			},
		)

		common.Log.PrintFull(common.Printf("%v", string(pemdata[:]), ), )
	} else {
		key, err := ioutil.ReadFile(svc.Config.ServicePublicKeyFileName)

		if err != nil {
			return nil, err
		}

		block, _ := pem.Decode(key)
		publicKey, err = x509.ParsePKIXPublicKey(block.Bytes)

		if err != nil {
			return nil, err
		}
	}

	return publicKey, nil
}
