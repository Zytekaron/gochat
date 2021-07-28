package types

var config *Config

func Init(cfg *Config) {
	config = cfg
}

type Opcode int

const (
	_ Opcode = iota
	OpHeartbeat
	OpHeartbeatAck
	OpInit
	OpHello
	OpInvalidOpcode
	OpInvalidSession
)

type Err int

const (
	ErrTimeout = 4000 + iota
	ErrInvalidOpcode
	ErrInvalidJson
)

//func containsString(slice []string, str string) bool {
//	for _, e := range slice {
//		if e == str {
//			return true
//		}
//	}
//	return false
//}
