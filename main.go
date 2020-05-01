package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	//"golang.org/x/net/websocket"

	"github.com/WolvenSpirit/job_scheduler/rxobservable"
	"golang.org/x/net/websocket"
)

var broadcastObs rxobservable.Observable
var jsonCodec = websocket.Codec{Marshal: jsonMarshal, Unmarshal: jsonUnmarshal}

type pulse struct {
	Time string `json:"time"`
}
type message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}
type hub interface {
	keepConn(string, *websocket.Conn)
	keepAlive(*websocket.Conn)
	getConn(string) *websocket.Conn
}
type hubx struct {
	ClientPool sync.Map
}
type msg struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var clientHub *hubx

func jsonMarshal(v interface{}) ([]byte, byte, error) {
	msg, err := json.Marshal(v)
	return msg, websocket.TextFrame, err
}
func jsonUnmarshal(msg []byte, payloadType byte, v interface{}) error {
	return json.Unmarshal(msg, v)
}

func (h *hubx) keepConn(client string, ws *websocket.Conn) {
	h.ClientPool.Store(client, ws)
}
func (h *hubx) getConn(client string) *websocket.Conn {
	if wsint, ok := h.ClientPool.Load(client); ok {
		return wsint.(*websocket.Conn)
	}
	return nil
}
func init() {
	clientHub = &hubx{}
	broadcastObs = rxobservable.Observable{}
}
func template(wr http.ResponseWriter, r *http.Request) {}
func wsoc(ws *websocket.Conn) {
	username := make([]byte, 1024)
	if _, err := ws.Read(username); err != nil {
		log.Println(err.Error())
	}
	clientHub.keepConn(string(username), ws)
	broadcastChan := make(chan interface{}, 1)
	broadcastObs.Subscribe(&broadcastChan)
	go func() {
		for {
			msg := <-broadcastChan
			if err := jsonCodec.Send(ws, msg); err != nil {
				log.Println(err.Error())
			}
		}
	}()

	for {
		msg := message{}
		if err := jsonCodec.Receive(ws, &msg); err != nil {
			log.Println(err.Error())
			break
		}
		broadcastObs.Next(msg)
	}

}
func listen() *http.Server {
	http.HandleFunc("/", template)
	http.Handle("/websocket", websocket.Handler(wsoc))
	s := &http.Server{Addr: "localhost:8080"}
	go func() {
		s.ListenAndServe()
	}()
	return s
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	s := listen()
	<-interrupt
	s.Shutdown(context.Background())
}
