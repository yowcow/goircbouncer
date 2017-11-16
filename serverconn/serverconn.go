package serverconn

import (
	"bufio"
	"io"
	"log"
	"net"
	_ "sync"
	"time"

	"github.com/yowcow/goircbouncer/config"
	command "github.com/yowcow/goirccommand"
	parser "github.com/yowcow/goircparser"
)

type EventHandler func(w io.Writer, row *parser.Row) bool

type Events map[string]EventHandler

type ServerConn struct {
	cfg    *config.Config
	logger *log.Logger
	events Events
}

func New(cfg *config.Config, logger *log.Logger) *ServerConn {
	s := &ServerConn{cfg, logger, make(Events)}
	s.RegisterEvent("PING", func(w io.Writer, row *parser.Row) bool {
		command.Pong(w, "", "", row.Suffix)
		return true
	})
	return s
}

func (s *ServerConn) RegisterEvent(command string, function EventHandler) {
	s.events[command] = function
}

func (s ServerConn) Start() {
	for {
		if err := s.ConnectOnce(); err != nil {
			s.logger.Println("failed connecting to server: ", err)
		}
		s.logger.Println("retrying in 10 seconds")
		time.Sleep(10 * time.Second)
	}
}

func (s ServerConn) ConnectOnce() error {
	conn, err := net.Dial("tcp", s.cfg.Server.Host+s.cfg.Server.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	s.initConn(w)
	return s.watchConn(r, w)
}

func (s ServerConn) initConn(w *bufio.Writer) {
	command.Nick(w, s.cfg.Server.Nick)
	command.User(w, s.cfg.Server.User, s.cfg.Server.Mode, s.cfg.Server.User)
	if len(s.cfg.Server.Channels) > 0 {
		command.Join(w, s.cfg.Server.Channels, []string{})
	}
	w.Flush()
}

func (s ServerConn) watchConn(r *bufio.Reader, w *bufio.Writer) error {
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return err
		}
		row := parser.Parse(string(line))
		s.logger.Println(row.RawLine)
		if wrote := s.handleRow(w, row); wrote {
			w.Flush()
		}
	}
	return nil
}

func (s ServerConn) handleRow(w io.Writer, row *parser.Row) bool {
	if f, ok := s.events[row.Command]; ok {
		return f(w, row)
	}
	return false
}
