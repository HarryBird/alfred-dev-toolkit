package toolkit

import (
	"strings"
)

const (
	name string = "[Alfred-Dev-Toolkit]"
)

func WhoAMI() string {
	return name
}

func sign(msg string, infos ...string) string {
	if len(infos) > 0 {
		msg = msg + " | " + strings.Join(infos, " | ")
	}
	return name + " " + msg
}
