package main

import (
	"fmt"
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"html/template"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r * http.Request) bool { return true },
}

func homePage(w http.ResponseWriter, r * http.Request) {
	fmt.Fprintf(w, "Home page")
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
		}
	}
}

func wsPage(w http.ResponseWriter, r * http.Request) {
	fmt.Println(r.Host)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client successfully connected...")

	reader(ws)
}

func chatPage(w http.ResponseWriter, r * http.Request) {
	var t, err = template.ParseFiles("chat.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		t.Execute(w, nil)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsPage)
	http.HandleFunc("/chat", chatPage)
}

func main() {
	fmt.Println("Init App")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8081", nil))
}