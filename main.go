package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/yowcow/goircbouncer/config"
	"github.com/yowcow/goircbouncer/serverconn"
	command "github.com/yowcow/goirccommand"
	parser "github.com/yowcow/goircparser"
)

func main() {
	var configfile string
	flag.StringVar(&configfile, "config", "", "path to config file")
	flag.Parse()

	cfg, err := config.Load(configfile)
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	svr := serverconn.New(cfg, logger)
	svr.RegisterEvent("PRIVMSG", func(w io.Writer, row *parser.Row) bool {
		command.Notice(w, row.Params[0], fmt.Sprintf("You said '%s'", row.Suffix))
		return true
	})
	svr.Start()
}
