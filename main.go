package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// TODO: Pick an appropriate buffer size. See Gorilla WebSocket documentation.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	hubChan := make(chan interface{})
	go hub(hubChan)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !websocket.IsWebSocketUpgrade(r) {
			http.ServeFile(w, r, "./main.html")
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		go pumpMessages(hubChan, conn)
	})
	http.HandleFunc("/main.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./main.js")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
	// TODO: Pick an appropriate port number.
}
