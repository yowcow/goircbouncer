package serverconn

import (
	"bufio"
	"bytes"
	"log"
	"testing"

	_ "github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goircbouncer/config"
	parser "github.com/yowcow/goircparser"
)

func Test_New(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", log.Lshortfile)
	cfg := new(config.Config)
	s := New(cfg, logger)

	assert.NotNil(t, s)
}

func Test_handleRow_responds_to_PING(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", log.Lshortfile)
	cfg := new(config.Config)
	s := New(cfg, logger)

	buf := new(bytes.Buffer)
	row := &parser.Row{
		Command: "PING",
		Suffix:  "hoge",
	}
	result := s.handleRow(buf, row)

	assert.Equal(t, true, result)
	assert.Equal(t, "PONG :hoge\r\n", buf.String())
}

func Test_handleRow_not_respond_to_NOTICE(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", log.Lshortfile)
	cfg := new(config.Config)
	s := New(cfg, logger)

	buf := new(bytes.Buffer)
	row := &parser.Row{
		Command: "NOTICE",
		Suffix:  "hoge",
	}
	result := s.handleRow(buf, row)

	assert.Equal(t, false, result)
	assert.Equal(t, "", buf.String())
}

func Test_initConn(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", log.Lshortfile)
	cfg := &config.Config{
		Server: &config.Server{
			Nick:     "mynick",
			User:     "myuser",
			Mode:     8,
			Channels: []string{"#hoge", "#fuga"},
		},
	}
	s := New(cfg, logger)

	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	s.initConn(w)

	assert.Equal(t, "NICK mynick\r\nUSER myuser 8 * :myuser\r\nJOIN #hoge,#fuga\r\n", buf.String())
}
