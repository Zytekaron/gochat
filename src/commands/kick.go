package commands

import (
	"gochat/src/types"
	"strings"
)

func KickCommand(ctx *types.CommandContext) {
	if len(ctx.Args) == 0 {
		return
	}

	user := server.GetClient(ctx.Args[0])
	if user == nil {
		ctx.Reply("That user does not exist")
		return
	}

	if user.GetRoom(ctx.Room) == nil {
		ctx.Reply("That user is not in this room")
		return
	}

	// check roles/perms

	reason := strings.Join(ctx.Args[1:], " ")
	if reason == "" {
		reason = "No reason provided"
	}

}
