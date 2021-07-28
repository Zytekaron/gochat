package main

import (
	"flag"
	"github.com/spf13/pflag"
	"gochat/src/server"
	"gochat/src/types"
	"strconv"
	"time"
)

// to-do list
// - add mutex locks if necessary, especially on Server. see Room#Members and Client#Rooms

var port int

func init() {
	pflag.IntVarP(&port, "port", "p", 1337, "http port")
	pflag.Parse()
}

func main() {
	flag.Parse()

	// temporary, will pull from file later
	config := &types.Config{
		Heartbeat: 60 * time.Second,
		CloseWait: time.Second,
	}
	types.Init(config)
	server.Init(config)

	server.Start(":" + strconv.Itoa(port))
}
