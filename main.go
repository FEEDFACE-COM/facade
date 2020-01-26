package main

import (
	"bufio"
	"flag"
	"fmt"
	"strings"

	//"bufio"
	//"io"
	"os"
	"os/signal"
	"runtime"

	facade "./facade"
	log "./log"
)

const (
	DEBUG_PERIODIC = true
	DEBUG_DIAG     = false
	DEBUG_CLOCK    = true
	DEBUG_CONFIG   = false
	DEBUG_FONT     = false
	DEBUG_SHADER   = false
	DEBUG_MODE     = true
	DEBUG_MEMORY   = false
	DEBUG_BUFFER   = true
	DEBUG_RENDERER = false
)

const FRAME_RATE = 60.0

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
	READ Command = "read"
	RECV Command = "recv"
	PIPE Command = "pipe"
	EXEC Command = "exec"
	CONF Command = "conf"
	INFO Command = "info"
	HELP Command = "help"
	TEST Command = "test"
)

var (
	textPort       uint    = 0xfcd
	port           uint    = 0xfcc
	host           string  = ""
	connectTimeout float64 = 5.0
	readTimeout    float64 = 0.0
)

const DEFAULT_LISTEN_HOST = "0.0.0.0"
const DEFAULT_CONNECT_HOST = "localhost"

func main() {
	quiet, verbose, debug := false, false, false
	directory := facade.DEFAULT_DIRECTORY

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		sig := <-signals
		log.Notice("%s", sig)
		os.Exit(0)
	}()

	log.SetVerbosity(log.NOTICE)

	globalFlags := flag.NewFlagSet("", flag.ExitOnError)
	globalFlags.SetOutput(bufio.NewWriter(nil))
	globalFlags.Usage = func() { ShowHelp(*globalFlags) }

	commandFlags := make(map[Command]*flag.FlagSet)

	var commands = []Command{CONF, PIPE, EXEC, INFO, TEST}
	if RENDERER_AVAILABLE {
		commands = append(commands, READ)
		commands = append(commands, RECV)
	}

	for _, cmd := range commands {
		commandFlags[cmd] = flag.NewFlagSet(string(cmd), flag.ExitOnError)
		commandFlags[cmd].Usage = func() { ShowHelp(*globalFlags) }
		commandFlags[cmd].SetOutput(bufio.NewWriter(nil))
	}

	for _, cmd := range []Command{PIPE, CONF, EXEC, INFO} {
		commandFlags[cmd].UintVar(&port, "port", port, "connect to server at `port`")
		commandFlags[cmd].StringVar(&host, "host", DEFAULT_CONNECT_HOST, "connect to server at `host`")
		commandFlags[cmd].Float64Var(&connectTimeout, "timeout", connectTimeout, "timeout connect after `seconds`")
	}

	if commandFlags[RECV] != nil {
		commandFlags[RECV].UintVar(&port, "port", port, "listen on `port` for messages")
		commandFlags[RECV].UintVar(&textPort, "textport", textPort, "listen on `port` for raw text")
		commandFlags[RECV].StringVar(&host, "host", DEFAULT_LISTEN_HOST, "listen on `host`")
		commandFlags[RECV].Float64Var(&readTimeout, "timeout", readTimeout, "timeout read after `seconds`")
	}

	if commandFlags[TEST] != nil {
		commandFlags[TEST].UintVar(&port, "port", port, "listen on `port` for messages")
		commandFlags[TEST].UintVar(&textPort, "textport", textPort, "listen on `port` for raw text")
		commandFlags[TEST].StringVar(&host, "host", DEFAULT_LISTEN_HOST, "listen on `host`")
		commandFlags[TEST].Float64Var(&readTimeout, "timeout", readTimeout, "timeout read after `seconds`")
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

	var client *Client
	var server *Server
	var scanner *Scanner
	var renderer *Renderer
	var tester *Tester
	var executor *Executor
	var path string

	cmd := Command(globalFlags.Args()[0])

	switch cmd {
	case READ, RECV:
		if !RENDERER_AVAILABLE {
			ShowHelp(*globalFlags)
			os.Exit(-2)
		}
		fallthrough
	case PIPE, CONF, EXEC, TEST:
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

	case READ, RECV, PIPE, CONF, TEST:
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

			case facade.Mode_LINE.String():
				config.SetMode = true
				config.Mode = facade.Mode_LINE
				config.Lines = &facade.LineConfig{}
				config.Lines.Shader = &facade.ShaderConfig{}
				config.Lines.Grid = &facade.GridConfig{}

			case facade.Mode_TAGS.String():
				config.SetMode = true
				config.Mode = facade.Mode_TAGS
				config.Tags = &facade.TagConfig{}
				config.Tags.Shader = &facade.ShaderConfig{}

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

	case INFO:
		// no args, print local info
		if len(globalFlags.Args()) <= 1 {
			ShowVersion()
			ShowAssets(directory)
			fmt.Fprintf(os.Stderr, "\n\n")
			os.Exit(0)
		} else {

			// query remote host
			commandFlags[INFO].Usage = func() { ShowHelpCommand(INFO, *commandFlags[INFO]) }
			commandFlags[INFO].SetOutput(bufio.NewWriter(nil))
			commandFlags[INFO].Parse(globalFlags.Args()[1:])

		}

	case HELP:
		ShowHelp(*globalFlags)
		os.Exit(-2)

	default:
		ShowHelp(*globalFlags)
		os.Exit(-2)

	}

	var err error
	confs := make(chan facade.Config)
	texts := make(chan facade.TextSeq)
	quers := make(chan (chan string))
	switch cmd {

	case READ:
		log.Info(AUTHOR)
		scanner = NewScanner()
		renderer = NewRenderer(directory)
		go scanner.ScanText(texts)
		runtime.LockOSThread()
		if err = renderer.Init(); err != nil {
			log.PANIC("fail to initialize renderer: %s", err)
		}
		renderer.Configure(config)
		go renderer.ProcessTextSeqs(texts)
		err = renderer.Render(nil)

	case RECV:
		log.Info(AUTHOR)
		server = NewServer(host, port, textPort, readTimeout)
		renderer = NewRenderer(directory)
		go server.Listen(confs, texts, quers)
		go server.ListenText(texts)
		runtime.LockOSThread()
		if err = renderer.Init(); err != nil {
			log.PANIC("fail to initialize renderer: %s", err)
		}
		renderer.Configure(config)
		go renderer.ProcessTextSeqs(texts)
		go renderer.ProcessQueries(quers)
		err = renderer.Render(confs)

	case PIPE:
		client = NewClient(host, port, connectTimeout)
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
		client = NewClient(host, port, connectTimeout)
		if err = client.Dial(); err != nil {
			log.Error("fail to dial: %s", err)
		}
		defer client.Close()
		if err = client.SendConf(config); err != nil {
			log.Error("fail to conf: %s", err)
		}

	case INFO:
		client = NewClient(host, port, connectTimeout)
		if err = client.Dial(); err != nil {
			log.Error("fail to dial: %s", err)
		}
		defer client.Close()
		info, err := client.QueryInfo()
		if err != nil {
			log.Error("fail to query: %s", err)
		} else {
			log.Notice("%s", info)
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

		client = NewClient(host, port, connectTimeout)
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

	case TEST:
		log.Info(AUTHOR)
		tester = NewTester(directory)
		scanner = NewScanner()
		go scanner.ScanText(texts)

		server = NewServer(host, port, textPort, readTimeout)
		go server.Listen(confs, texts, quers)
		go server.ListenText(texts)

		runtime.LockOSThread()
		tester.Init()
		tester.Configure(config)

		//start processing only after init!
		go tester.ProcessTextSeqs(texts)
		go tester.ProcessQueries(quers)
		err = tester.Test(confs)

	default:
		log.PANIC("unexpected command %s", cmd)
	}

	if err != nil {
		log.Error("fail to %s: %s", cmd, err)
		os.Exit(-1)
	}

	os.Exit(0)
}
