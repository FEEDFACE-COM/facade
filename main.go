package main

import (
	"FEEDFACE.COM/facade/facade"
	"FEEDFACE.COM/facade/log"
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	DEBUG_RENDERER = false
	DEBUG_BUFFER   = true
	DEBUG_CHANGES  = false
	DEBUG_DIAG     = false
)

const RENDER_FRAME_RATE = 60.0
const TEXT_BUFFER_SIZE = 1024

const DEFAULT_RECEIVE_HOST = "[::]"
const DEFAULT_CONNECT_HOST = "localhost"

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
	stdin          bool    = false
	noTitle        bool    = true
)

func main() {
	runtime.LockOSThread()
	quiet, debug := false, false
	directory := facade.DEFAULT_DIRECTORY
	var err error

	confs := make(chan facade.Config)
	texts := make(chan facade.TextSeq)
	ticks := make(chan bool, 2)

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
		commandFlags[cmd].BoolVar(&noIPv4, "noinet", noIPv4, "disable IPv4 networking")
		commandFlags[cmd].BoolVar(&noIPv6, "noinet6", noIPv6, "disable IPv6 networking")
	}

	if RENDERER_AVAILABLE {
		globalFlags.StringVar(&directory, "D", directory, "asset directory")

		commandFlags[SERVE].UintVar(&port, "port", port, "listen on `port` for config")
		commandFlags[SERVE].UintVar(&textPort, "textport", textPort, "listen on `textport` for text")
		commandFlags[SERVE].StringVar(&receiveHost, "host", DEFAULT_RECEIVE_HOST, "listen on `host` for config and text")
		commandFlags[SERVE].Float64Var(&readTimeout, "timeout", readTimeout, "timeout read after `seconds`")
		commandFlags[SERVE].BoolVar(&noIPv4, "noinet", noIPv4, "disable IPv4 networking")
		commandFlags[SERVE].BoolVar(&noIPv6, "noinet6", noIPv6, "disable IPv6 networking")
		commandFlags[SERVE].BoolVar(&stdin, "stdin", stdin, "also read text from stdin")
		commandFlags[SERVE].BoolVar(&noTitle, "notitle", noTitle, "no title on startup")
	}

	{
		globalFlags.BoolVar(&debug, "d", debug, "show debug info")
		globalFlags.BoolVar(&quiet, "q", quiet, "show errors only")
	}

	globalFlags.Parse(os.Args[1:])
	if globalFlags.NArg() < 1 {
		ShowHelp(*globalFlags)
		os.Exit(-2)
	}
	if debug {
		log.SetVerbosity(log.DEBUG)
	} else if quiet {
		log.SetVerbosity(log.ERROR)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		for {
			sig := <-signals
			log.Notice("signal %s", sig)
			ticks <- false
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
	config.Shader = &facade.ShaderConfig{}

	var args []string
	var modeFlags *flag.FlagSet

	switch cmd {

	case EXEC:

		config.SetMode = true
		config.Mode = facade.Mode_TERM
		config.Term = &facade.TermConfig{}

		args = commandFlags[cmd].Args()
		if len(args) > 0 && strings.ToUpper(args[0]) == facade.Mode_TERM.String() {
			args = args[1:]
		}

		modeFlags = flag.NewFlagSet("exec", flag.ExitOnError)
		modeFlags.SetOutput(bufio.NewWriter(nil))
		modeFlags.Usage = func() { ShowHelpMode(EXEC, config.Mode, *modeFlags) }

		config.AddFlags(modeFlags, facade.Mode_TERM)
		modeFlags.Parse(args)
		config.VisitFlags(modeFlags)

		args = modeFlags.Args()
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
				config.Term = &facade.TermConfig{}

			case facade.Mode_LINES.String():
				config.SetMode = true
				config.Mode = facade.Mode_LINES
				config.Lines = &facade.LineConfig{}

			case facade.Mode_WORDS.String():
				config.SetMode = true
				config.Mode = facade.Mode_WORDS
				config.Words = &facade.WordConfig{}

			case facade.Mode_CHARS.String():
				config.SetMode = true
				config.Mode = facade.Mode_CHARS
				config.Chars = &facade.CharConfig{}

			default:
				ShowHelpCommand(cmd, *commandFlags[cmd])
				os.Exit(-2)

			}
			modeFlags = flag.NewFlagSet(strings.ToLower(config.Mode.String()), flag.ExitOnError)
			modeFlags.Usage = func() { ShowHelpMode(cmd, config.Mode, *modeFlags) }
			modeFlags.SetOutput(bufio.NewWriter(nil))
			config.AddFlags(modeFlags, config.Mode)
			modeFlags.Parse(args[1:])
			config.VisitFlags(modeFlags)
		} else {
			config = nil
		}

	case README:
		readme, err := base64.StdEncoding.DecodeString(facade.Asset["README"])
		if err != nil {
			log.PANIC("fail to decode readme: %s", err)
		}
		os.Stdout.Write(readme)
		os.Exit(0)

	case HELP:
		fallthrough

	default:
		ShowHelp(*globalFlags)
		os.Exit(-2)

	}

	switch cmd {

	case SERVE:
		// REM: use DEBUGINFO_HEIGHT constant here..
		fmt.Fprintf(os.Stderr,"\033[2J\033[H") // clear screen, jump to origin
		fmt.Fprintf(os.Stderr,"\033[%d;1H", 16+1) // cursor down


		log.Notice(strings.TrimLeft(AUTHOR,"\n"))
		runtime.LockOSThread()

		server = NewServer(receiveHost, port, textPort, readTimeout, noIPv4, noIPv6)
		renderer = NewRenderer(directory, ticks)
		if stdin {
			scanner = NewScanner()
			go scanner.ScanText(texts)
		}
		go server.Listen(confs, texts)
		go server.ListenText(texts)
		runtime.LockOSThread()
		if err = renderer.Init(); err != nil {
			log.PANIC("fail to initialize renderer: %s", err)
		}

		renderer.Configure(config)
		go renderer.ProcessTextSeqs(texts)

		//if !noTitle {
		//	titleConfig := &facade.Config{}
		//	if renderer.mode == facade.Mode_TERM {
		//		if config.Terminal == nil || !config.Terminal.Grid.GetSetFill() {
		//			gridConfig := &facade.GridConfig{SetFill: true, Fill: "title"}
		//			titleConfig.Terminal = &facade.TermConfig{Grid: gridConfig}
		//		}
		//	} else if renderer.mode == facade.Mode_LINES {
		//		if config.Lines == nil || !config.Lines.Grid.GetSetFill() {
		//			gridConfig := &facade.GridConfig{SetFill: true, Fill: "title"}
		//			titleConfig.Lines = &facade.LineConfig{Grid: gridConfig}
		//		}
		//	}
		//	renderer.Configure(titleConfig)
		//}

		err = renderer.Render(confs,debug)
		if err != nil {
			log.Error("fail to render: %s", err)
		}
		renderer.Finish()

	case PIPE:
		client = NewClient(connectHost, port, connectTimeout, noIPv4, noIPv6)
		if err = client.Dial(); err != nil {
			log.Error("fail to dial: %s", err)
		}
		defer client.Close()
		if config != nil {
			log.Info("configure %s", config.Desc())
			if client.SendConf(config); err != nil {
				log.PANIC("fail to send conf: %s", err)
			}
		}
		if err = client.OpenTextStream(); err != nil {
			log.PANIC("fail to open stream: %s", err)
		}
		defer client.CloseTextStream()
		if err = client.ScanAndSendText(); err != nil {
			log.Error("fail to scan and send: %s", err)
		}
		time.Sleep(time.Duration(int64(time.Second / 10.))) //wait until all text flushed

	case CONF:
		client = NewClient(connectHost, port, connectTimeout, noIPv4, noIPv6)
		if err = client.Dial(); err != nil {
			log.Error("fail to dial: %s", err)
		}
		defer client.Close()
		if config != nil {
			log.Info("configure %s", config.Desc())
			if err = client.SendConf(config); err != nil {
				log.PANIC("fail to conf: %s", err)
			}
		}

	case EXEC:

		var cols, rows = facade.TermDefaults.Width, facade.TermDefaults.Height
		if config.GetTerm() == nil {
			config.Term = &facade.TermConfig{}
		}

		if config.GetTerm().GetSetWidth() {
			cols = config.GetTerm().GetWidth()
		}
		if config.GetTerm().GetSetHeight() {
			rows = config.GetTerm().GetHeight()
		}

		config.Term.Width = cols
		config.Term.SetWidth = true

		config.Term.Height = rows
		config.Term.SetHeight = true

		client = NewClient(connectHost, port, connectTimeout, noIPv4, noIPv6)
		executor = NewExecutor(client, uint(cols), uint(rows), path, args)

		if err = client.Dial(); err != nil {
			log.Error("fail to dial: %s", err)
		}
		defer client.Close()
		if config != nil {
			log.Info("configure %s", config.Desc())
			if client.SendConf(config); err != nil {
				log.PANIC("fail to conf: %s", err)
			}
		}
		if err = client.OpenTextStream(); err != nil {
			log.PANIC("fail to open stream: %s", err)
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

const AUTHOR = `
   _   _   _   _   _   _      _   _   _   _   _   _   _   _     _   _        
  |_  |_| /   |_| | \ |_     |_  |_  |_  | \ |_  |_| /   |_    /   / \ |\/|  
  |   | | \_  | | |_/ |_  BY |   |_  |_  |_/ |   | | \_  |_  o \_  \_/ |  |  
`
