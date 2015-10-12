package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/peterh/liner"
)

var debug = false

func main() {
	app := cli.NewApp()
	app.Name = "etcdsh"
	app.Usage = "Shell into etcd"
	app.Version = "0.0.1"
	app.Action = func(c *cli.Context) {
		replCommand(c)
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "peers, C", Value: "", Usage: "A comma delimited list of machine addresses"},
		cli.StringFlag{Name: "username, u", Value: "", Usage: "Provide username[:password] for etcd, prompts with no password"},
		cli.BoolFlag{Name: "debug, d", Usage: "Turn on shell debugging"},
	}

	app.Run(os.Args)
}

func replCommand(cliContext *cli.Context) {

	debug = cliContext.Bool("debug")

	commands := []ReplCommand{
		NewSetCommand(),
		NewGetCommand(),
		NewPwdCommand(),
		NewEnvCommand(),
		NewLsCommand(),
		NewCdCommand(),
	}

	shellState := getShellState()
	shellState.cliContext = cliContext

	initTerm()
	shellState.kapi = getKapi()

	replExit := false

	prompt := "> "

	for !replExit {
		prompt = shellState.pwd + "> "
		line, err := shellState.term.Prompt(prompt)
		if err != nil {
			break
		}
		if line == "exit" {
			break
		}
		shellState.term.AppendHistory(line)
		// This will need to be replaced with a smarter parser

		replYell, err := makeReplCommand(commands, line)
		if err == nil {
			if debug {
				fmt.Println(replYell)
			}
			replYell.Command.Action(replYell)
		} else {
			fmt.Println(err)
		}
	}

	closeTerm()
}

func initTerm() {
	getShellState().term = liner.NewLiner()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		closeTerm()
		fmt.Println("Exiting etcdsh")
		os.Exit(1)
	}()
}

func getKapi() client.KeysAPI {

	peerlist := getShellState().cliContext.String("peers")
	if peerlist == "" {
		peerlist = "http://127.0.0.1:4001,http://127.0.0.1:2379"
	}
	peers := strings.Split(peerlist, ",")

	if debug {
		fmt.Println(peers)
	}

	username := ""
	password := ""

	// username/pwd next
	userpass := getShellState().cliContext.String("username")

	if userpass != "" {
		colon := strings.Index(userpass, ":")
		if colon == -1 {
			username = userpass
			// Slight hack below
			tmppass, err := getShellState().term.PasswordPrompt("Password:")
			if err != nil {
				fmt.Println(err)
				closeTerm()
				log.Fatal(err)
			}
			password = tmppass
		} else {
			username = userpass[:colon]
			password = userpass[colon+1:]
		}
	}
	cfg := client.Config{
		Endpoints:               peers,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
		Username:                username,
		Password:                password,
	}

	etcdclient, err := client.New(cfg)
	if err != nil {
		closeTerm()
		log.Fatal(err)
	}

	kapi := client.NewKeysAPI(etcdclient)

	_, err = kapi.Get(context.TODO(), "/", nil)

	if err != nil {
		fmt.Println(err)
	}

	return kapi
}

// ShellState for all your variable needs
type ShellState struct {
	cliContext *cli.Context
	term       *liner.State
	pwd        string
	format     string
	kapi       client.KeysAPI
	etcdClient client.Client
}

var currShellState ShellState

func getShellState() *ShellState {
	if currShellState.pwd == "" {
		currShellState.pwd = "/"
		currShellState.format = "short"
	}
	return &currShellState
}

func closeTerm() {
	getShellState().term.Close()
}
