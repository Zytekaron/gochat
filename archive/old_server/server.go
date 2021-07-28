package old_server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gochat/src/commands"
	"gochat/src/types"
	"log"
	"net/http"
)

var (
	server   *types.Server
	cmds map[string]types.Command
)

func Start(addr string) {
	loadCommands()

	server = newServer()
	prepareUsers(server)
	go runWs()

	fmt.Println("Listening on", addr)
	err := http.ListenAndServe(addr, newRouter())
	if err != nil {
		panic(err)
	}
}

type RouteHandler func(*types.Server, http.ResponseWriter, *http.Request)

func wrap(handler RouteHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler(server, w, req)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	http.ServeFile(w, r, "index.html")
}

func runWs() {
	for {
		select {
		case client := <-server.Register:
			fmt.Println("Client Connected:", client.User.Id)
			server.Clients[client.User.Id] = client
		case client := <-server.Unregister:
			if _, ok := server.Clients[client.User.Id]; ok {
				fmt.Println("Client Disconnected:", client.User.Id)
				delete(server.Clients, client.User.Id)
				close(client.Send)
			}
		case message := <-server.Messages:
			process(message)
		}
	}
}

func prepareUsers(server *types.Server) {
	ch := make(chan []byte)
	server.Register <- &types.Client{
		User:     &types.User{ID: "0", Username: "admin"},
		Server:   server,
		Conn:     &websocket.Conn{},
		SendChan: ch,
	}
	server.Register <- &types.Client{
		User:     &types.User{ID: "1", Username: "Zytekaron"},
		Server:   server,
		Conn:     &websocket.Conn{},
		SendChan: ch,
	}
}

func loadCommands() {
	cmds = make(map[string]types.Command)
	cmds["ping"] = types.Command{Name: "ping", Run: commands.PingCommand}
	cmds["leave"] = types.Command{Name: "leave", Run: commands.LeaveCommand}
}

func newServer() *types.Server {
	return &types.Server{
		Messages:   make(chan types.PacketMessage),
		Register:   make(chan *types.Client),
		Unregister: make(chan *types.Client),
		Clients:    make(map[string]*types.Client),
	}
}
