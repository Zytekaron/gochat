package old_server

import (
	"bytes"
	"encoding/json"
	gonanoid "github.com/matoous/go-nanoid"
	"gochat/src/types"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 60 * time.Second

	pongWait = 10 * time.Second
	delay    = (pongWait * 9) / 10

	maxMessageSize = 512
)

var (
	newline   = []byte{'\n'}
	space     = []byte{' '}
	roomRegex = regexp.MustCompile("\\w{1,32}")
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(server *types.Server, w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	if !roomRegex.MatchString(room) {
		_ = types.WriteErrorJson(w, 400, -1, "That room is taken!")
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // fixme danger CORS Policy
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	id := gonanoid.MustID(21)
	client := &types.Client{
		User:     &types.User{ID: id, Username: room},
		Server:   server,
		Conn:     conn,
		SendChan: make(chan []byte, 256),
	}
	client.Server.Register <- client

	init, _ := json.Marshal(&types.Packet{
		Event:  "INIT",
		Opcode: 1,
		Data:   struct{ id string }{id},
	})
	client.SendChan <- init

	go readPump(client)
	go writePump(client)
}

func readPump(client *types.Client) {
	defer func() {
		client.Server.Unregister <- client
		_ = client.Conn.Close()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	_ = client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		_ = client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Error: %v\n", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		client.Server.Messages <- types.PacketMessage{Client: client, Bytes: message}
	}
}

func writePump(client *types.Client) {
	ticker := time.NewTicker(delay)
	defer func() {
		ticker.Stop()
		_ = client.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.SendChan:
			_ = client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			n := len(client.SendChan)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-client.SendChan)
			}

			if e := w.Close(); e != nil {
				return
			}
		case <-ticker.C:
			_ = client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
