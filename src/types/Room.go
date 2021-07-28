package types

import "log"

type Room struct {
	ID       string       `json:"id"`
	OwnerID  string       `json:"owner_id"`
	Settings RoomSettings `json:"settings"`
	Members  map[string]*RoomMember
}

type RoomSettings struct {
	Password string `json:"-"`

}

type RoomMember struct {
	// todo consider refactor
	Client *Client `json:"-"`
	User   *User   `json:"user"`
	Role   string  `json:"role"` // member, mod, admin, owner
}

func (r *Room) Broadcast(event string, data interface{}) {
	for _, c := range r.Members {
		go func(client *Client) {
			err := client.Send(event, data)
			if err != nil {
				log.Println(err)
			}
		}(c.Client)
	}
}

func (r *Room) GetInRole(role string) []*RoomMember {
	var inRole []*RoomMember
	for _, m := range r.Members {
		if m.Role == role {
			inRole = append(inRole, m)
		}
	}
	return inRole
}

func (r *Room) UserJoin(client *Client) {
	r.Members[client.User.ID] = &RoomMember{
		Client: client,
		User:   client.User,
		Role:   "member",
	}
	r.Broadcast("USER_JOIN", client.User)
}

func (r *Room) UserLeave(client *Client) {
	delete(r.Members, client.User.ID)
	r.Broadcast("USER_LEAVE", client.User)
}

func (r *Room) GetMember(query string) *RoomMember {
	if m, ok := r.Members[query]; ok {
		return m
	}
	for _, m := range r.Members {
		if m.User.Username == query {
			return m
		}
	}
	return nil
}
