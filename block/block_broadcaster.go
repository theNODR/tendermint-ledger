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
	"time"
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

func (ws *BlockBroadcaster) Start() error {
	common.Log.Event("EventStartBlockBroadcaster", common.Printf("Broadcaster start"))
	uri, err := url.Parse(ws.addr)
	uri.Scheme = "ws"
	for ws.conn == nil {
		uriStr := uri.String()
		ws.conn, _, err = websocket.DefaultDialer.Dial(uriStr, nil)
		if err != nil {
			fmt.Printf("Connect to blockchain error: %s", err)
			common.Log.Error(ErrorConnectToBlockProducer, common.Printf("Error: %v", err))
			time.Sleep(30 * time.Second)
		}
	}

	go func() {
		go ws.hub.run()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			common.Log.Event("EventHandleRoot", common.Printf("Request:%v", r.RemoteAddr))
			serveWs(ws.hub, w, r)
		})

		err := http.ListenAndServe(ws.addrBroadcast, nil)
		if err != nil {
			common.Log.Error("ErrorListenAndServe: ", common.Printf("Error listening: %v", err))
		}
	}()

	go func() {
		defer close(ws.send)
		for {
			block := &nodr.BlockInfo{}
			err := ws.conn.ReadJSON(block)
			if err != nil {
				common.Log.Error("ErrorReadFromBlockWs", common.Printf("Error listening: %v", err))
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
					break
				}
			}

			bytes, err := json.Marshal(block)
			if err != nil {
				common.Log.Error("ErrorMarshalBlock", common.Printf("Error marshal block info: %v", block))
				continue
			}
			ws.hub.broadcast <- bytes
		}
	}()

	return nil
}

func (ws *BlockBroadcaster) Stop() error {
	err := ws.conn.Close()

	if err != nil {
		return err
	}
	return nil
}
