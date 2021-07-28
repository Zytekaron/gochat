package server

import (
	"github.com/gorilla/websocket"
	"gochat/src/types"
	"log"
)

func reader(client *types.Client) {
	for {
		var packet *types.Packet
		err := client.Conn.ReadJSON(&packet)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println(err)
			}
			break
		}

		if packet.Opcode > 0 {
			handleOp(client, packet)
		}

		client.Recv <- packet
	}
}
