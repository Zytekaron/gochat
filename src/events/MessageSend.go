package events

import (
	"encoding/json"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid"
	"gochat/src/types"
	"net/http"
	"strings"
)

func MessageSend(server *types.Server, cmds map[string]types.Command, w http.ResponseWriter, req *http.Request) {
	message, ok := getMessage(w, req)
	if !ok {
		return
	}

	message.Id = gonanoid.MustID(21)

	if len(message.Content) < 1 {
		_ = types.WriteErrorJson(w, 400, -1, "Malformed post body: Content is too short")
		return
	}

	if len(message.Content) > 2048 {
		_ = types.WriteErrorJson(w, 400, -1, "Malformed post body: Content is too long")
		return
	}

	if strings.HasPrefix(message.Content, "/") {
		fmt.Println("processing command")
		processCommand(cmds, message)
	}

	fmt.Printf("Bytes: [%s] %s: %s\n", message.Id, message.Author, message.Content)

	res := types.Packet{
		Event:  "SEND_MESSAGE",
		Opcode: 0,
		Data:   message,
	}

	server.BroadcastJson(&res)
}

func getMessage(w http.ResponseWriter, req *http.Request) (*types.Message, bool) {
	var message *types.Message
	err := json.NewDecoder(req.Body).Decode(&message)
	if err != nil {
		_ = types.WriteErrorJson(w, 400, -1, "Malformed post body: body does not conform to Bytes type")
		return message, false
	}

	message.Content = strings.Trim(message.Content, " ")
	return message, true
}

func processCommand(cmds map[string]types.Command, message *types.Message) {
	// Slice off the '/' prefix and any extra whitespace
	str := strings.TrimPrefix(message.Content, "/")
	str = strings.TrimSpace(str)

	// Split into arguments
	args := strings.Split(str, " ")
	command, args := args[0], args[1:]

	// Grab the command and run it (if it exists)
	cmd, ok := cmds[command]
	if ok {
		cmd.Run(args)
	}
}