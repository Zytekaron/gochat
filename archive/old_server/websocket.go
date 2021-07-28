package old_server

import (
	"fmt"
	"gochat/src/types"
)

func process(message types.PacketMessage) {
	fmt.Println("Incoming Packet:", string(message.Bytes), "(old_server.go:46)")
}