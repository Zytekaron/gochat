package server

import (
	"gochat/src/types"
	"log"
)

func handleOp(client *types.Client, packet *types.Packet) {
	switch packet.Opcode {
	case types.OpHeartbeat:
		client.Heartbeat()
		client.send(&types.Packet{
			Opcode: types.OpHeartbeatAck,
		})
	case types.OpHeartbeatAck:
		// todo impl client heartbeat ack
	case types.OpInit:
		client.Init(packet)
	default:
		err := client.Disconnect(types.ErrInvalidOpcode)
		if err != nil {
			log.Println(err)
		}
	}
}
