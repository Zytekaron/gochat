package types

type CommandContext struct {
	Client   *Client
	Username string
	Room     string
	Args     []string
}

func (c *CommandContext) Reply(text string) {
	panic("unimplemented")
}
