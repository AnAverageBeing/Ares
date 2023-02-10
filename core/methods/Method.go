package methods

import (
	"Ares/core"
	"fmt"
	"strings"
)

type Method interface {
	Start()
	Stop()
}

func GetMethod(name string, conf *core.AttackConfig) (Method, error) {
	switch strings.ToLower(name) {
	case "join":
		return Join{Config: conf}, nil
	case "ping":
		return Ping{Config: conf}, nil
	default:
		return nil, fmt.Errorf("no method with name %s found", name)
	}
}
