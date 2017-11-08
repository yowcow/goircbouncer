package command

import (
	"fmt"
	"io"
	"strings"
)

// Pong writes a PONG command to writer
func Pong(w io.Writer, server string) error {
	_, err := fmt.Fprintf(w, "PONG :%s\r\n", server)
	return err
}

// Nick writes a NICK command to writer
func Nick(w io.Writer, nick string) error {
	_, err := fmt.Fprintf(w, "NICK %s\r\n", nick)
	return err
}

// User writes a USER command to writer
func User(w io.Writer, username string, mode int, realname string) error {
	_, err := fmt.Fprintf(w, "USER %s %d * :%s\r\n", username, mode, realname)
	return err
}

// Join writes a JOIN command to writer
func Join(w io.Writer, channels []string) error {
	_, err := fmt.Fprintf(w, "JOIN %s\r\n", strings.Join(channels, ","))
	return err
}

// Names writes a NAMES command to writer
func Names(w io.Writer, channels []string) error {
	_, err := fmt.Fprintf(w, "NAMES :%s\r\n", strings.Join(channels, ","))
	return err
}
