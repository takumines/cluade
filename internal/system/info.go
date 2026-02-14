package system

import (
	"os"
	"os/user"
	"strings"
)

func GetUsername() string {
	u, err := user.Current()
	if err != nil {
		return "user"
	}
	return u.Username
}

func GetCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "~"
	}

	u, err := user.Current()
	if err != nil {
		return dir
	}

	if strings.HasPrefix(dir, u.HomeDir) {
		return "~" + dir[len(u.HomeDir):]
	}
	return dir
}
