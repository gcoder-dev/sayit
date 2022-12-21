package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var logFile, _ = os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
var pool = Pool{pool: map[*websocket.Conn]bool{}}

type OpCode uint16

const (
	closeConnection OpCode = iota
	register
	addUser
)

type Message struct {
	Type OpCode `json:"type"`
	Text string `json:"message"`
}

func (r *Message) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, r)
	if err != nil {
		log.Println("Unsupported data")
	}
	return err
}

func main() {
	log.SetOutput(logFile)
	defer logFile.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/ws", handleWebsocket)
	http.HandleFunc("/", landingView)
	log.Fatal(http.ListenAndServe("192.168.1.227:8888", nil))
}

//func middlewareBroadcast(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("Executing middlewareOne")
//		next.ServeHTTP(w, r)
//	})
//}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	upgrade := websocket.Upgrader{}
	c, _ := upgrade.Upgrade(w, r, nil)
	go reader(c)
}

func reader(c *websocket.Conn) {
	defer c.Close()
	defer pool.unregister(c)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return
		}
		res := Message{}
		err = res.Unmarshal(message)
		if err != nil {
			return
		}
		if res.Type == closeConnection {
			return
		} else if res.Type == register {
			pool.register(c)
			broadcastUsers(c.RemoteAddr().String())
		}
	}
}

func broadcastUsers(name string) {
	for i := range pool.pool {
		for j := range pool.pool {
			if i != j {
				err := i.WriteJSON(Message{Type: addUser, Text: name})
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
