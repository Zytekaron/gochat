package types

import (
	"errors"
	"log"
)

type Server struct {
	Rooms   map[string]*Room
	Clients map[string]*Client
}

func (s *Server) UserJoinRoom(userID, roomID, password string) error {
	client, ok := s.Clients[userID]
	if !ok {
		return errors.New("user does not exist")
	}

	room, ok := s.Rooms[roomID]
	if !ok {
		return errors.New("room does not exist")
	}

	if room.Password != "" {
		if password == "" {
			return errors.New("password is required")
		}
		return errors.New("password is invalid")
	}

	// todo consider delegating to user
	client.Rooms[roomID] = room
	room.UserJoin(client)

	return nil
}

func (s *Server) UserLeaveRoom(userID, roomID string) error {
	c, ok := s.Clients[userID]
	if !ok {
		return errors.New("user does not exist")
	}

	r, ok := s.Rooms[roomID]
	if !ok {
		return errors.New("room does not exist")
	}
	delete(c.Rooms, roomID)
	delete(r.Members, userID)

	return nil
}

func (s *Server) Broadcast(event string, data interface{}) {
	for _, c := range s.Clients {
		go func(client *Client) {
			err := client.Send(event, data)
			if err != nil {
				log.Println(err)
			}
		}(c)
	}
}

func (s *Server) GetClient(query string) *Client {
	if c, ok := s.Clients[query]; ok {
		return c
	}
	for _, c := range s.Clients {
		if c.User.Username == query {
			return c // todo refactor with a set? also see: Room#GetMember
		}
	}
	return nil
}

func (s *Server) IsUsernameTaken(username string) bool {
	return s.GetClient(username) != nil
}
