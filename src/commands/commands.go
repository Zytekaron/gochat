package commands

import "gochat/src/types"

var server *types.Server

func Init(s *types.Server) {
	server = s
}
