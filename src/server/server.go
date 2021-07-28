package server

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	gonanoid "github.com/matoous/go-nanoid"
	"gochat/src/commands"
	"gochat/src/types"
	"log"
	"net/http"
)

var server = &types.Server{
	Clients: make(map[string]*types.Client),
}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var config *types.Config

func Init(cfg *types.Config) {
	config = cfg

	commands.Init(server) // todo reorder
}

func Start(addr string) {
	router := chi.NewRouter()

	router.Get("/ws", serveWs)

	fmt.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(*http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	id := gonanoid.MustID(32)
	log.Println("new client with id", id)

	user := types.NewUser(id, "")
	client := types.NewClient(user, server, ws)

	server.Clients[id] = client

	client.Heartbeat()
	go reader(client)
}
