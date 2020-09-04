package main

import (
	"bufio"
	"flag"
	"strings"
	"time"

	"encoding/base64"
	"fmt"
	//"bufio"
	//"io"
	"os"
	"os/signal"
	"runtime"
	//	"runtime/debug"

	facade "./facade"
	log "./log"
)

const (
	DEBUG_PERIODIC = true
	DEBUG_EVENTS   = false
	DEBUG_DIAG     = false
	DEBUG_CLOCK    = true
	DEBUG_CONFIG   = false
	DEBUG_FONT     = false
	DEBUG_SHADER   = false
	DEBUG_MODE     = true
	DEBUG_MEMORY   = false
	DEBUG_BUFFER   = false
	DEBUG_RENDERER = false
)

const RENDER_FRAME_RATE = 60.0
const TEXT_BUFFER_SIZE = 1024

const DEFAULT_RECEIVE_HOST = "[::]"
const DEFAULT_CONNECT_HOST = "localhost"

const AUTHOR = `
   _   _   _   _   _   _      _   _   _   _   _   _   _   _     _   _        
  |_  |_| /   |_| | \ |_     |_  |_  |_  | \ |_  |_| /   |_    /   / \ |\/|  
  |   | | \_  | | |_/ |_  BY |   |_  |_  |_/ |   | | \_  |_  o \_  \_/ |  |  
`

var (
	BUILD_NAME     string = "facade"
	BUILD_VERSION  string = "0.0.0"
	BUILD_PLATFORM string = "os-arch"
	BUILD_DATE     string = "1970-01-01"
)

type Command string

const (
	SERVE  Command = "serve"
	PIPE   Command = "pipe"
	EXEC   Command = "exec"
	CONF   Command = "conf"
	HELP   Command = "help"
	README Command = "readme"
)

var (
	textPort       uint    = 0xfcd
	port           uint    = 0xfcc
	receiveHost    string  = ""
	connectHost    string  = ""
	connectTimeout float64 = 5.0
	readTimeout    float64 = 0.0
	noIPv4         bool    = false
	noIPv6         bool    = false
	noStdin        bool    = false
)

func main() {
	quiet, verbose, debug := false, false, false
	directory := facade.DEFAULT_DIRECTORY

	var err error
	confs := make(chan facade.Config)
	texts := make(chan facade.TextSeq)
	quers := make(chan (chan string))
	pause := make(chan bool, 10)

	log.SetVerbosity(log.NOTICE)

	globalFlags := flag.NewFlagSet("", flag.ExitOnError)
	globalFlags.SetOutput(bufio.NewWriter(nil))
	globalFlags.Usage = func() { ShowHelp(*globalFlags) }

	commandFlags := make(map[Command]*flag.FlagSet)

	var commands = []Command{CONF, PIPE, EXEC}
	if RENDERER_AVAILABLE {
		commands = append(commands, SERVE)
	}

	for _, cmd := range commands {
		commandFlags[cmd] = flag.NewFlagSet(string(cmd), flag.ExitOnError)
		commandFlags[cmd].Usage = func() { ShowHelp(*globalFlags) }
		commandFlags[cmd].SetOutput(bufio.NewWriter(nil))
	}

	for _, cmd := range []Command{PIPE, CONF, EXEC} {
		commandFlags[cmd].UintVar(&port, "port", port, "connect to server at `port`")
		commandFlags[cmd].StringVar(&connectHost, "host", DEFAULT_CONNECT_HOST, "connect to server at `host`")
		commandFlags[cmd].Float64Var(&connectTimeout, "timeout", connectTimeout, "timeout connect after `seconds`")
		commandFlags[cmd].BoolVar(&noIPv4, "no4", noIPv4, "disable IPv4 networking")
		commandFlags[cmd].BoolVar(&noIPv6, "no6", noIPv6, "disable IPv6 networking")
	}

	if commandFlags[SERVE] != nil {
		commandFlags[SERVE].UintVar(&port, "port", port, "listen on `port` for config")
		commandFlags[SERVE].UintVar(&textPort, "textport", textPort, "listen on `port` for text")
		commandFlags[SERVE].StringVar(&receiveHost, "host", DEFAULT_RECEIVE_HOST, "listen on `host` for config and text")
		commandFlags[SERVE].Float64Var(&readTimeout, "timeout", readTimeout, "timeout read after `seconds`")
		commandFlags[SERVE].BoolVar(&noIPv4, "noinet", noIPv4, "disable IPv4 networking")
		commandFlags[SERVE].BoolVar(&noIPv6, "noinet6", noIPv6, "disable IPv6 networking")
		commandFlags[SERVE].BoolVar(&noStdin, "nostdin", noStdin, "do not read from stdin")
	}

	{
		globalFlags.StringVar(&directory, "D", directory, "working directory")
		globalFlags.BoolVar(&verbose, "v", verbose, "show info messages?")
		globalFlags.BoolVar(&debug, "d", debug, "show debug messages?")
		globalFlags.BoolVar(&quiet, "q", quiet, "show warnings only?")
	}

	globalFlags.Parse(os.Args[1:])
	if globalFlags.NArg() < 1 {
		ShowHelp(*globalFlags)
		os.Exit(-2)
	}
	if debug {
		log.SetVerbosity(log.DEBUG)
	} else if verbose {
		log.SetVerbosity(log.INFO)
	} else if quiet {
		log.SetVerbosity(log.WARNING)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		last := time.Now()
		for {
			sig := <-signals
			if !debug || (time.Now().Sub(last) < 250*time.Millisecond) {
				log.Notice("SIGNAL %s", sig)
				os.Exit(0)
			} else {
				//        		log.Notice("signal %s", sig)
				select {
				case pause <- true:
				default:
				}
			}
			last = time.Now()
		}
	}()

	var client *Client
	var server *Server
	var scanner *Scanner
	var renderer *Renderer
	var executor *Executor
	var path string

	cmd := Command(globalFlags.Args()[0])

	switch cmd {
	case SERVE:
		if !RENDERER_AVAILABLE {
			ShowHelp(*globalFlags)
			os.Exit(-2)
		}
		fallthrough
	case PIPE, CONF, EXEC:
		commandFlags[cmd].Usage = func() { ShowHelpCommand(cmd, *commandFlags[cmd]) }
		commandFlags[cmd].Parse(globalFlags.Args()[1:])
	}

	var config *facade.Config = &facade.Config{}
	config.Font = &facade.FontConfig{}
	config.Camera = &facade.CameraConfig{}
	config.Mask = &facade.MaskConfig{}

	var args []string
	var modeFlags *flag.FlagSet

	switch cmd {

	case EXEC:

		config.SetMode = true
		config.Mode = facade.Mode_TERM
		config.Terminal = &facade.TermConfig{}
		config.Terminal.Shader = &facade.ShaderConfig{}
		config.Terminal.Grid = &facade.GridConfig{}

		args = commandFlags[cmd].Args()
		if len(args) > 0 && strings.ToUpper(args[0]) == facade.Mode_TERM.String() {
			args = args[1:]
		}

		modeFlags = flag.NewFlagSet("exec", flag.ExitOnError)
		modeFlags.SetOutput(bufio.NewWriter(nil))
		modeFlags.Usage = func() { ShowHelpMode(EXEC, config.Mode, *modeFlags) }

		config.AddFlags(modeFlags)
		modeFlags.Parse(args)
		config.VisitFlags(modeFlags)

		args = modeFlags.Args()
		log.Debug("execute %v", args)
		if len(args) <= 0 {
			ShowHelpMode(EXEC, facade.Mode_TERM, *modeFlags)
			os.Exit(-2)
		}

		path = args[0]
		args = args[1:]

	case SERVE, PIPE, CONF:
		args = commandFlags[cmd].Args()

		// parse mode, if given
		if len(args) > 0 {

			switch strings.ToUpper(args[0]) {

			case facade.Mode_TERM.String():
				config.SetMode = true
				config.Mode = facade.Mode_TERM
				config.Terminal = &facade.TermConfig{}
				config.Terminal.Shader = &facade.ShaderConfig{}
				config.Terminal.Grid = &facade.GridConfig{}

			case facade.Mode_LINES.String():
				config.SetMode = true
				config.Mode = facade.Mode_LINES
				config.Lines = &facade.LineConfig{}
				config.Lines.Shader = &facade.ShaderConfig{}
				config.Lines.Grid = &facade.GridConfig{}

			case facade.Mode_WORDS.String():
				config.SetMode = true
				config.Mode = facade.Mode_WORDS
				config.Words = &facade.WordConfig{}
				config.Words.Shader = &facade.ShaderConfig{}
				config.Words.Set = &facade.SetConfig{}

			case facade.Mode_TAGS.String():
				config.SetMode = true
				config.Mode = facade.Mode_TAGS
				config.Tags = &facade.TagConfig{}
				config.Tags.Shader = &facade.ShaderConfig{}
				config.Tags.Set = &facade.SetConfig{}

			case facade.Mode_DRAFT.String():
				config.SetMode = true
				config.Mode = facade.Mode_DRAFT

			default:
				ShowHelpCommand(cmd, *commandFlags[cmd])
				os.Exit(-2)

			}

			modeFlags = flag.NewFlagSet(strings.ToLower(config.Mode.String()), flag.ExitOnError)
			modeFlags.Usage = func() { ShowHelpMode(cmd, config.Mode, *modeFlags) }
			modeFlags.SetOutput(bufio.NewWriter(nil))
			config.AddFlags(modeFlags)
			modeFlags.Parse(args[1:])
			config.VisitFlags(modeFlags)
		}

	case README:
		readme, err := base64.StdEncoding.DecodeString(facade.Asset["README"])
		if err != nil {
			log.PANIC("fail to decode readme: %s", err)
		}
		fmt.Fprintf(os.Stdout, string(readme))
		os.Exit(0)

	case HELP:
		fallthrough

	default:
		ShowHelp(*globalFlags)
		os.Exit(-2)

	}

	switch cmd {

	case SERVE:
		log.Info(AUTHOR)

		server = NewServer(receiveHost, port, textPort, readTimeout, noIPv4, noIPv6)
		renderer = NewRenderer(directory, ticks)
		if !noStdin {
			scanner = NewScanner()
			go scanner.ScanText(texts)
		}
		go server.Listen(confs, texts, quers)
		go server.ListenText(texts)
		runtime.LockOSThread()
		if err = renderer.Init(); err != nil {
			log.PANIC("fail to initialize renderer: %s", err)
		}
		renderer.Configure(config)
		go renderer.ProcessTextSeqs(texts)
		go renderer.ProcessQueries(quers)
		err = renderer.Render(confs, pause)

	case PIPE:
		client = NewClient(connectHost, port, connectTimeout, noIPv4, noIPv6)
		if err = client.Dial(); err != nil {
			log.Error("fail to dial: %s", err)
		}
		defer client.Close()
		if config != nil {
			log.Debug("configure %s", config.Desc())
			if client.SendConf(config); err != nil {
				log.Error("fail to send conf: %s", err)
			}
		}
		if err = client.OpenTextStream(); err != nil {
			log.Error("fail to open stream: %s", err)
		}
		defer client.CloseTextStream()
		if err = client.ScanAndSendText(); err != nil {
			log.Error("fail to scan and send: %s", err)
		}

	case CONF:
		if config == nil {
			ShowHelpMode(cmd, facade.Mode_DRAFT, *modeFlags)
			os.Exit(-1)
		}
		log.Debug("configure %s", config.Desc())
		client = NewClient(connectHost, port, connectTimeout, noIPv4, noIPv6)
		if err = client.Dial(); err != nil {
			log.Error("fail to dial: %s", err)
		}
		defer client.Close()
		if err = client.SendConf(config); err != nil {
			log.Error("fail to conf: %s", err)
		}

	case EXEC:

		var cols, rows = facade.GridDefaults.Width, facade.GridDefaults.Height
		if config.GetTerminal() == nil {
			config.Terminal = &facade.TermConfig{}
			config.Terminal.Grid = &facade.GridConfig{}
		}

		if config.GetTerminal().GetGrid().GetSetWidth() {
			cols = config.GetTerminal().GetGrid().GetWidth()
		}
		if config.GetTerminal().GetGrid().GetSetHeight() {
			rows = config.GetTerminal().GetGrid().GetHeight()
		}

		config.Terminal.Grid.Width = cols
		config.Terminal.Grid.SetWidth = true

		config.Terminal.Grid.Height = rows
		config.Terminal.Grid.SetHeight = true

		client = NewClient(connectHost, port, connectTimeout, noIPv4, noIPv6)
		executor = NewExecutor(client, uint(cols), uint(rows), path, args)

		if err = client.Dial(); err != nil {
			log.Error("fail to dial: %s", err)
		}
		defer client.Close()
		if config != nil {
			log.Debug("configure %s", config.Desc())
			if client.SendConf(config); err != nil {
				log.Error("fail to send conf: %s", err)
			}
		}
		if err = client.OpenTextStream(); err != nil {
			log.Error("fail to open stream: %s", err)
		}
		defer client.CloseTextStream()
		err = executor.Execute()

	default:
		log.PANIC("unexpected command %s", cmd)
	}

	if err != nil {
		log.Error("fail to %s: %s", cmd, err)
		os.Exit(-1)
	}

	os.Exit(0)
}
