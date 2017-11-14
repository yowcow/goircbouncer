package serverconn

import (
	"bufio"
	"io"
	"log"
	"net"
	"time"

	"github.com/yowcow/goircbouncer/config"
	command "github.com/yowcow/goirccommand"
	parser "github.com/yowcow/goircparser"
)

type EventHandler func(w io.Writer, row *parser.Row)

type Events map[string]EventHandler

type ServerConn struct {
	cfg    *config.Config
	logger *log.Logger
	quit   chan<- bool
	events Events
}

func New(cfg *config.Config, logger *log.Logger, quit chan<- bool) *ServerConn {
	events := make(Events)
	events["PING"] = func(w io.Writer, row *parser.Row) {
		command.Pong(w, "", "", row.Suffix)
	}
	return &ServerConn{cfg, logger, quit, events}
}

func (s *ServerConn) RegisterEvent(command string, function EventHandler) {
	s.events[command] = function
}

func (s ServerConn) Start() {
	for {
		conn, err := net.Dial("tcp", s.cfg.Server.Host+s.cfg.Server.Addr)
		if err != nil {
			s.logger.Println("failed connecting to server: ", err)
		} else {
			if err = s.handleConn(conn); err != nil {
				s.logger.Println("failed while handling a connection to server: ", err)
			}
		}
		s.logger.Println("retrying in 10 secs")
		time.Sleep(time.Second * 10)
	}
}

func (s ServerConn) handleConn(conn net.Conn) error {
	defer conn.Close()

	command.Nick(conn, s.cfg.Server.Nick)
	command.User(conn, s.cfg.Server.User, 0, s.cfg.Server.User)
	if len(s.cfg.Server.Channels) > 0 {
		command.Join(conn, s.cfg.Server.Channels, []string{})
	}

	comChan := make(chan *parser.Row)
	errChan := make(chan error)
	go watchConn(conn, comChan, errChan)

	w := bufio.NewWriter(conn)
	for {
		select {
		case row := <-comChan:
			if f, ok := s.events[row.Command]; ok {
				f(w, row)
				w.Flush()
				s.logger.Println("handled event: ", row.Command)
			} else {
				s.logger.Println(row.RawLine)
			}
		case err := <-errChan:
			close(errChan)
			close(comChan)
			return err
		}
	}

	return nil
}

func watchConn(conn net.Conn, comChan chan<- *parser.Row, errChan chan<- error) {
	r := bufio.NewReader(conn)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			errChan <- err
			return
		}
		row := parser.Parse(string(line))
		comChan <- row
	}
}
