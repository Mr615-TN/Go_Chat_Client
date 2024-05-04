package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Open or create SQLite database file
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	// Create a table if it does not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
                            id INTEGER PRIMARY KEY,
                            content TEXT
                        );`)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// Use the msg variable
		log.Println(string(msg))
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}
