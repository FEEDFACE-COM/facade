package main

import (
	"FEEDFACE.COM/facade/facade"
	"FEEDFACE.COM/facade/log"
	"bufio"
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
	DEFAULT_DIRECTORY = ""
	DEFAULT_RECEIVE_HOST = "[::]"
	DEFAULT_CONNECT_HOST = "localhost"
	DEFAULT_BASIC_OPTIONS = true
	DEFAULT_DEBUG_VIEW   = false

)

const PAUSABLE = true
const DEBUG_BUFFERS = true

var (
	BUILD_NAME     string = "facade"
	BUILD_VERSION  string = "0.0.0"
	BUILD_PLATFORM string = "os-arch"
	BUILD_DATE     string = "1970-01-01"
)

type Command string

const (
	SERVE  Command = "render"
	PIPE   Command = "pipe"
	EXEC   Command = "exec"
	CONF   Command = "conf"
	HELP   Command = "help"
	//README Command = "readme"
)

type Tick int

const (
	TICK Tick = iota
	QUIT Tick = iota
	STOP Tick = iota
)

var (
	textPort       uint    = 0xfcd
	port           uint    = 0xfcc
	receiveHost    string  = ""
	connectHost    string  = ""
	connectTimeout float64 = 5.0
	readTimeout    float64 = 0.0
	ipv4           bool    = true
	ipv6           bool    = true
	stdin          bool    = false
	title          bool    = true
)

func main() {
	runtime.LockOSThread()
	quiet, debug := false, false
	basicOptions, debugView := DEFAULT_BASIC_OPTIONS, DEFAULT_DEBUG_VIEW
	directory := DEFAULT_DIRECTORY
	var err error

	confs := make(chan facade.Config)
	texts := make(chan facade.TextSeq)
	ticks := make(chan Tick, 2)

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
		commandFlags[cmd].BoolVar(&ipv4, "inet", ipv4, "use IPv4 networking")
		commandFlags[cmd].BoolVar(&ipv6, "inet6", ipv6, "use IPv6 networking")
	}

	if RENDERER_AVAILABLE {
		commandFlags[SERVE].StringVar(&directory, "dir", DEFAULT_DIRECTORY, "asset directory")
		commandFlags[SERVE].UintVar(&port, "port-fcd", port, "listen on `port-fcd`")
		commandFlags[SERVE].UintVar(&textPort, "port-txt", textPort, "listen on `port-txt` for raw text")
		commandFlags[SERVE].StringVar(&receiveHost, "host", DEFAULT_RECEIVE_HOST, "listen on `host` for config and text")
		commandFlags[SERVE].Float64Var(&readTimeout, "timeout", readTimeout, "timeout read after `seconds`")
		commandFlags[SERVE].BoolVar(&ipv4, "inet", ipv4, "use IPv4 networking")
		commandFlags[SERVE].BoolVar(&ipv6, "inet6", ipv6, "use IPv6 networking")
		commandFlags[SERVE].BoolVar(&stdin, "stdin", stdin, "read text from stdin")
		commandFlags[SERVE].BoolVar(&title, "title", title, "show title on startup")
	}

	{
		globalFlags.BoolVar(&debug, "d", debug, "debug messages and flags")
		globalFlags.BoolVar(&quiet, "q", quiet, "error messages only")
		globalFlags.BoolVar(&debugView, "D", debugView, "debug view")
	}

	globalFlags.Parse(os.Args[1:])
	if globalFlags.NArg() < 1 {
		ShowHelp(*globalFlags)
		os.Exit(-2)
	}
	if debug {
		log.SetVerbosity(log.DEBUG)
		basicOptions = false
	} else if quiet {
		log.SetVerbosity(log.ERROR)
	}

	var client *Client
	var server *Server
	var scanner *Scanner
	var renderer *Renderer
	var executor *Executor

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

	var config = &facade.Config{}

	var args []string
	var modeFlags *flag.FlagSet
	var path string

	switch cmd {

	case EXEC:

		args = commandFlags[cmd].Args()
		if len(args) <= 0 || strings.ToUpper(args[0]) != facade.Mode_TERM.String() {
			ShowHelpMode(EXEC, facade.Mode_TERM, basicOptions)
			os.Exit(-2)
		}
		args = args[1:]

		modeFlags = flag.NewFlagSet("exec", flag.ExitOnError)
		modeFlags.SetOutput(bufio.NewWriter(nil))

		config = facade.NewConfig(facade.Mode_TERM)
		modeFlags.Usage = func() { ShowHelpMode(EXEC, facade.Mode_TERM, basicOptions) }
		config.AddFlags(modeFlags, facade.Mode_TERM, basicOptions)
		modeFlags.Parse(args)
		config.VisitFlags(modeFlags,facade.Mode_TERM, basicOptions)

		args = modeFlags.Args()
		if len(args) <= 0 {
			ShowHelpMode(EXEC, facade.Mode_TERM, basicOptions)
			os.Exit(-2)
		}

		path = args[0]
		args = args[1:]

	case SERVE, PIPE, CONF:
		args = commandFlags[cmd].Args()
		mode := facade.DEFAULT_MODE

		// parse mode, if given
		if len(args) > 0 {

			switch strings.ToUpper(args[0]) {

			case facade.Mode_TERM.String():
				mode = facade.Mode_TERM
			case facade.Mode_LINES.String():
				mode = facade.Mode_LINES
			case facade.Mode_WORDS.String():
				mode = facade.Mode_WORDS
			case facade.Mode_CHARS.String():
				mode = facade.Mode_CHARS
			default:
				ShowHelpCommand(cmd, *commandFlags[cmd])
				os.Exit(-2)

			}

			args = args[1:]

			modeFlags = flag.NewFlagSet(strings.ToLower(mode.String()), flag.ExitOnError)
			modeFlags.SetOutput(bufio.NewWriter(nil))

			config = facade.NewConfig(mode)
			modeFlags.Usage = func() { ShowHelpMode(cmd, config.Mode, basicOptions) }
			config.AddFlags(modeFlags, config.Mode, basicOptions)
			modeFlags.Parse(args)
			config.VisitFlags(modeFlags, config.Mode, basicOptions)

			if len( modeFlags.Args() ) > 0 {
				ShowHelpMode(cmd,config.Mode, basicOptions)
				os.Exit(-2)
			}


		}

	//case README:
	//	readme, err := base64.StdEncoding.DecodeString(facade.Asset["README"])
	//	if err != nil {
	//		log.PANIC("fail to decode readme: %s", err)
	//	}
	//	os.Stdout.Write(readme)
	//	os.Exit(0)

	case HELP:
		fallthrough

	default:
		ShowHelp(*globalFlags)
		os.Exit(-2)

	}

	switch cmd {

	case SERVE:
		if DEBUG_BUFFERS && debugView {
			// REM: use DEBUGINFO_HEIGHT constant here..
			fmt.Fprintf(os.Stderr, "\033[2J\033[H")    // clear screen, jump to origin
			fmt.Fprintf(os.Stderr, "\033[%d;1H", 16+1) // cursor down
		} else {
			log.Notice(strings.TrimLeft(AUTHOR, "\n"))
		}

		// install signal handler
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGQUIT)
		go handleSignals(signals, ticks, debugView)

		runtime.LockOSThread()

		// add title if not given
		if title && !config.SetFill {
			config.SetFill = true
			config.Fill = "title"
		}

		server = NewServer(receiveHost, port, textPort, readTimeout, ipv4, ipv6)
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

		err = renderer.Render(confs, debugView)
		if err != nil {
			log.Error("fail to render: %s", err)
		}
		renderer.Finish()

	case PIPE:
		client = NewClient(connectHost, port, connectTimeout, ipv4, ipv6)
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
		client = NewClient(connectHost, port, connectTimeout, ipv4, ipv6)
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

		client = NewClient(connectHost, port, connectTimeout, ipv4, ipv6)
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

func handleSignals(signals chan os.Signal, ticks chan Tick, pausable bool) {
	for {
		sig := <-signals
		log.Notice("signal %s", sig)
		if PAUSABLE && pausable && sig == syscall.SIGQUIT {
			ticks <- STOP
			continue
		}
		ticks <- QUIT
	}
}


const AUTHOR = `
   _   _   _   _   _   _                  
  |_  |_| /   |_| | \ |_                  
  |   | | \_  | | |_/ |_  by FEEDFACE.COM 
`
