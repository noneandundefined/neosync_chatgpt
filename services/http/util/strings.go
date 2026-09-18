package util

import "neomatica/neosync/config"

func EscapeString(s string) string {
	escaped, _ := config.JSON.Marshal(s)
	return string(escaped[1 : len(escaped)-1])
}
