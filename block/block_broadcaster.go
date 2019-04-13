package block

import (
	"common"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	nodr "svcnodr/ws"
)

const ErrorConnectToBlockProducer = "error connect to block"

type BlockBroadcaster struct {
	addr          string
	addrBroadcast string
	conn          *websocket.Conn
	hub           *Hub
	send          chan nodr.BlockInfo
}

func NewWsBlockBroadcaster(addr string, addresBroadcast string) *BlockBroadcaster {
	return &BlockBroadcaster{
		addr:          addr,
		addrBroadcast: addresBroadcast,
		hub:           newHub(),
		send:          make(chan nodr.BlockInfo, 100),
	}
}

func (ws BlockBroadcaster) Start() error {
	uri, err := url.Parse(ws.addr)
	uri.Path = "blocks"

	for ws.conn == nil {
		ws.conn, _, err = websocket.DefaultDialer.Dial(uri.String(), nil)
		if err != nil {
			common.Log.Error(ErrorConnectToBlockProducer, common.Printf("Error: %v", err))
		}
	}

	go func() {
		defer close(ws.send)
		for {
			block := &nodr.BlockInfo{}
			err := ws.conn.ReadJSON(block)
			if err != nil {
				log.Println("read:", err)
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}

			bytes, err := json.Marshal(block)
			fmt.Printf("\nBlockInfo Receive:%v\n", string(bytes))

			if err != nil {
				continue
			}
			ws.hub.broadcast <- bytes
		}
	}()

	go func() {
		go ws.hub.run()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			serveWs(ws.hub, w, r)
		})
		err := http.ListenAndServe(ws.addrBroadcast, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	return nil
}

func (ws BlockBroadcaster) Stop() error {
	err := ws.conn.Close()

	if err != nil {
		return err
	}
	return nil
}
