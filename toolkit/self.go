package toolkit

import (
	"strings"
)

const (
	SelfName string = "[Alfred-Dev-Toolkit]"
)

func WhoAMI() string {
	return SelfName
}

func sign(msg string, infos ...string) string {
	if len(infos) > 0 {
		msg = msg + " | " + strings.Join(infos, " | ")
	}
	return SelfName + " " + msg
}
