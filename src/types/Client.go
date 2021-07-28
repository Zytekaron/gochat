package types

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type Client struct {
	// Client information
	User   *User
	Server *Server
	Conn   *websocket.Conn
	Rooms  map[string]*Room

	// WebSocket variables
	Recv      chan *Packet
	seq       int
	heartbeat *time.Timer
	sendMutex sync.Mutex
}

func NewClient(user *User, server *Server, conn *websocket.Conn) *Client {
	return &Client{
		User:      user,
		Server:    server,
		Conn:      conn,
		Rooms:     make(map[string]*Room),
		Recv:      make(chan *Packet),
	}
}

func (c *Client) SendOp(op Opcode) error {
	return c.SendData(op, "", nil)
}

func (c *Client) Send(event string, data interface{}) error {
	return c.SendData(0, event, data)
}

func (c *Client) SendData(op Opcode, event string, data interface{}) error {
	if op != 0 {
		return c.send(&Packet{Opcode: op})
	}

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	p := &Packet{
		Event: event,
		Data:  d,
	}
	p.Seq = c.seq
	c.seq++

	return c.send(p)
}

func (c *Client) send(p *Packet) error {
	c.sendMutex.Lock()
	defer c.sendMutex.Unlock()
	return c.Conn.WriteJSON(p)
}

func (c *Client) InRoom(id string) bool {
	_, ok := c.Rooms[id]
	return ok
}

func (c *Client) SendMsg(msg *Message) {

}

// Heartbeat is called when the server receives an
// OpcodeHeartbeat
func (c *Client) Heartbeat() {
	if c.heartbeat != nil {
		c.heartbeat.Stop()
	}

	c.heartbeat = time.AfterFunc(config.Heartbeat, func() {
		log.Println("Client", c.User.ID, "disconnecting: timeout")
		err := c.Disconnect(ErrTimeout)
		if err != nil {
			log.Println(err)
		}
	})
}

func (c *Client) Init(packet *Packet) { // todo fix
	// todo clients and such
	var data *struct {
		Username string `json:"username"`
	}
	err := json.Unmarshal(packet.Data, &data)
	if err != nil {
		err = c.Disconnect(ErrInvalidJson)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// todo register to c.Server and do username checks


	// todo send a bunch of data
	c.Send("READY", nil)
}

func (c *Client) Disconnect(code Err) error {
	c.heartbeat.Stop()
	time.Sleep(config.CloseWait)
	return c.Conn.Close()
}
