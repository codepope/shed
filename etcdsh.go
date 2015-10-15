package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/client"
	"github.com/peterh/liner"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app      = kingpin.New("etcdsh", "An etcd shell")
	peerlist = app.Flag("peers", "etcd peers").Short('C').Default("http://127.0.0.1:4001,http://127.0.0.1:2379").OverrideDefaultFromEnvar("EX_PEERS").String()
	userpass = app.Flag("user", "etcd User").Short('u').OverrideDefaultFromEnvar("EX_USER").String()
	password = app.Flag("pass", "etcd Password").Short('p').OverrideDefaultFromEnvar("EX_PASS").String()
	debug    = app.Flag("debug", "debug messages").Short('d').Bool()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.MustParse(app.Parse(os.Args[1:]))

	commands := []ReplCommand{
		NewSetCommand(),
		NewGetCommand(),
		NewPwdCommand(),
		NewEnvCommand(),
		NewLsCommand(),
		NewCdCommand(),
	}

	shellState := getShellState()

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
			if *debug {
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

	peers := strings.Split(*peerlist, ",")

	if *debug {
		fmt.Println(peers)
	}

	myusername := ""
	mypassword := ""

	// username/pwd next
	if *userpass != "" {
		colon := strings.Index(*userpass, ":")
		if colon == -1 {
			myusername = *userpass
			if *password != "" {
				mypassword = *password
			} else {
				tmppass, err := getShellState().term.PasswordPrompt("Password:")
				if err != nil {
					fmt.Println(err)
					closeTerm()
					log.Fatal(err)
				}
				mypassword = tmppass
			}
		} else {
			myusername = (*userpass)[:colon]
			mypassword = (*userpass)[colon+1:]
		}
	}
	cfg := client.Config{
		Endpoints:               peers,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
		Username:                myusername,
		Password:                mypassword,
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
