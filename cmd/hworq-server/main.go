package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yuuki/hworq/pkg/db"
	"github.com/yuuki/hworq/pkg/web"
)

const (
	DefaultPort = "9000"
)

// CLI is the command line object.
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		port    string
		version bool
	)

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.Usage = func() {
		fmt.Fprint(cli.errStream, helpText)
	}
	flags.StringVar(&port, "port", DefaultPort, "")
	flags.StringVar(&port, "P", DefaultPort, "")
	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&version, "v", false, "")

	if err := flags.Parse(args[1:]); err != nil {
		return 1
	}

	if version {
		fmt.Fprintf(cli.errStream, "%s version %s, build %s \n", Name, Version, GitCommit)
		return 0
	}

	log.Println("Connecting postgres ...")
	db, err := db.New()
	if err != nil {
		log.Printf("postgres initialize error: %v\n", err)
		return 2
	}

	handler := web.New(&web.Option{
		Port: port,
		DB:   db,
	})
	go handler.Run()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGTERM, syscall.SIGINT)
	s := <-sigch
	if err := handler.Shutdown(s); err != nil {
		log.Println(err)
		return 3
	}

	return 0
}

var helpText = `
Usage: hworq-server [options]

  simple HTTP-based Job Queue

Options:
  --port, -P           Listen port
`
