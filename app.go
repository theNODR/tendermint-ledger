package svcledger

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"svcledger/block"
	"svcledger/blockchain/tendermint"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/pkg/errors"

	"app"
	"app/httputil"
	"app/websocket"
	"app/websocket_gorilla"
	"common"
	"common/flags"
	"svcledger/helpers"
	"svcledger/job"
	"svcledger/store"
)

const (
	TypeName                           = "ledger"
	servicePrivateKeyLength            = 2048
	autoCloseInterval           uint64 = 30
	writeFinancialBlockInterval uint64 = 60
	syncTransactionInterval     uint64 = 3
)

var (
	tendermintEndPointUrl        = flags.String("ledger-tendermint-endpoint-url", "tcp://localhost:26657", "Tendermint CLI endpoint")
	priceAmountIncome            = flags.Uint64("ledger-price-amount-income", 500, "Min price of traffic")
	priceAmountSpend             = flags.Uint64("ledger-price-amount-spend", 500, "Max price of traffic")
	amountInitial                = flags.Uint64("ledger-amount-initial", 1000000, "Initial tokens amount on wallet")
	amountSpend                  = flags.Uint64("ledger-amount-spend", 10000, "Tokens in spend address after it opened")
	priceQuantumPower            = flags.Uint("ledger-quantum-power", 25, "Price power of quantum to transfer")
	servicePrivateKeyFileName    = flags.String("ledger-private-key-file-name", "", "Service private key")
	servicePublicKeyFileName     = flags.String("ledger-public-key-file-name", "", "Service public key")
	incomeChannelLifeTime        = flags.Int64("ledger-income-channel-lifetime", 30, "Income channel lifetime, min")
	spendChannelLifeTime         = flags.Int64("ledger-spend-channel-lifetime", 30, "Spend channel lifetime, min")
	transferChannelLifeTime      = flags.Int64("ledger-transfer-channel-lifetime", 5, "Transfer channel lifetime, min")
	chainId                      = flags.String("chain-id", "testchain", "Chain ID for thendermint node")
	blockInfoServeEndPointUrl    = flags.String("block-listen-endpoint", ":2000", "Endpoint for listen blocks info")
	tendermintBlocksInfoEndPoint = flags.String("tendermint-blocks-endpoint", "ws://localhost:8081", "Endpoint for connect to tenderming for listen Block infos")
)

type Config struct {
	*app.Config
	TendermintEndPointUrl     string
	Amounts                   *store.Amounts
	ServicePrivateKeyFileName string
	ServicePublicKeyFileName  string
	LifeTimes                 *store.LifeTimes
	ChainId                   string
}

type Service struct {
	address       string
	keyPair       helpers.KeyPair
	silentSch     *gocron.Scheduler
	silentSchStop chan bool

	BlockBroadcast   *block.BlockBroadcaster
	BlockchainClient *tendermint.TendermintClient
	Config           *Config
	Hub              websocket.Huber
	SocketHandler    websocket.SocketHandler
	Ledger           *store.Ledger
}

func NewService(config *app.Config) (*Service, error) {
	svcConfig := &Config{
		Config:                config,
		TendermintEndPointUrl: *tendermintEndPointUrl,
		Amounts: &store.Amounts{
			Initial:           store.TransactionAmountType(*amountInitial),
			Spend:             store.TransactionAmountType(*amountSpend),
			PriceIncome:       store.TransactionAmountType(*priceAmountIncome),
			PriceSpend:        store.TransactionAmountType(*priceAmountSpend),
			PriceQuantumPower: store.QuantumPowerType(*priceQuantumPower),
		},
		ServicePrivateKeyFileName: *servicePrivateKeyFileName,
		ServicePublicKeyFileName:  *servicePublicKeyFileName,
		LifeTimes: &store.LifeTimes{
			IncomeChannel:   int64(time.Duration(*incomeChannelLifeTime) * time.Minute / time.Millisecond),
			SpendChannel:    int64(time.Duration(*spendChannelLifeTime) * time.Minute / time.Millisecond),
			TransferChannel: int64(time.Duration(*transferChannelLifeTime) * time.Minute / time.Millisecond),
		},
		ChainId: *chainId,
	}

	client := tendermint.NewTendermintClient(svcConfig.ChainId, svcConfig.TendermintEndPointUrl)

	svc := &Service{
		BlockchainClient: client,
		Config:           svcConfig,
		BlockBroadcast:   block.NewWsBlockBroadcaster(*tendermintBlocksInfoEndPoint, *blockInfoServeEndPointUrl),
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
		common.Log.PrintFull(common.Printf("%v", pk))
	}
	sk, err := keyPair.PrivateKey()
	if err != nil {
		return err
	} else {
		common.Log.PrintFull(common.Printf("%v", sk))
	}

	svc.keyPair = keyPair
	svc.address = helpers.CreateCommonAddress(pk)

	return nil
}

func (svc *Service) startLedger() error {
	var err error

	ledger, err := store.NewLedger(
		svc.address,
		svc.keyPair,
		svc.Config.Amounts,
		svc.Config.LifeTimes,
		svc.BlockchainClient,
	)

	if err != nil {
		return err
	}

	svc.Ledger = ledger
	return nil
}

func (svc *Service) startBlockBroadcaster() error {
	return svc.BlockBroadcast.Start()
}

func (svc *Service) startWebSocket() error {
	var err error

	handler, err := handlerFunc(
		svc.keyPair,
		svc.Ledger,
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

	schedulerStop := scheduler.Start()

	return scheduler, schedulerStop
}

func (svc *Service) Start() error {
	var err error

	err = svc.startKeys()
	if err != nil {
		return err
	}

	err = svc.startLedger()
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

	_, err = svc.Ledger.Init(svc.Config.Amounts.Initial)
	if err != nil {
		common.Log.PrintFull(
			common.Printf("init tokens on wallet error: %v", err),
		)
	}

	err = svc.startWebSocket()
	if err != nil {
		return err
	}

	err = svc.startBlockBroadcaster()
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

	if svc.silentSchStop != nil {
		svc.silentSchStop <- true
	}

	if svc.BlockBroadcast != nil {
		svc.BlockBroadcast.Stop()
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
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
			},
		)

		common.Log.PrintFull(common.Printf("%v", string(pemdata[:])))

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
				Type:  "PUBLIC KEY",
				Bytes: publicKeyBytes,
			},
		)

		common.Log.PrintFull(common.Printf("%v", string(pemdata[:])))
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
