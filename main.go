package main

import (
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
	DEBUG_PERIODIC = false
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

	flag.Usage = ShowHelp

	flags := make(map[Command]*flag.FlagSet)

	var commands = []Command{CONF, PIPE, EXEC, INFO, TEST}
	if RENDERER_AVAILABLE {
		commands = append(commands, READ)
		commands = append(commands, RECV)
	}

	for _, cmd := range commands {
		flags[cmd] = flag.NewFlagSet(string(cmd), flag.ExitOnError)
	}

	for _, cmd := range []Command{PIPE, CONF, EXEC, INFO} {
		flags[cmd].UintVar(&port, "port", port, "connect to server at `port`")
		flags[cmd].StringVar(&host, "host", DEFAULT_CONNECT_HOST, "connect to server at `host`")
		flags[cmd].Float64Var(&connectTimeout, "timeout", connectTimeout, "timeout connect after `seconds`")
	}

	if flags[RECV] != nil {
		flags[RECV].UintVar(&port, "port", port, "listen on `port` for messages")
		flags[RECV].UintVar(&textPort, "textport", textPort, "listen on `port` for raw text")
		flags[RECV].StringVar(&host, "host", DEFAULT_LISTEN_HOST, "listen on `host`")
		flags[RECV].Float64Var(&readTimeout, "timeout", readTimeout, "timeout read after `seconds`")
	}

	if flags[TEST] != nil {
		flags[TEST].UintVar(&port, "port", port, "listen on `port` for messages")
		flags[TEST].UintVar(&textPort, "textport", textPort, "listen on `port` for raw text")
		flags[TEST].StringVar(&host, "host", DEFAULT_LISTEN_HOST, "listen on `host`")
		flags[TEST].Float64Var(&readTimeout, "timeout", readTimeout, "timeout read after `seconds`")
	}

	{
		flag.CommandLine.StringVar(&directory, "D", directory, "working directory")
	}

	flag.CommandLine.BoolVar(&verbose, "v", verbose, "show info messages?")
	flag.CommandLine.BoolVar(&debug, "d", debug, "show debug messages?")
	flag.CommandLine.BoolVar(&quiet, "q", quiet, "show warnings only?")

	flag.Parse()
	if flag.NArg() < 1 {
		ShowHelp()
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

	cmd := Command(flag.Args()[0])

	switch cmd {
	case READ, RECV:
		if !RENDERER_AVAILABLE {
			ShowHelp()
			os.Exit(-2)
		}
		fallthrough
	case PIPE, CONF, EXEC, TEST:
		flags[cmd].Usage = func() { ShowHelpCommand(cmd, *flags[cmd]) }
		flags[cmd].Parse(flag.Args()[1:])
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

		args = flags[cmd].Args()
		if len(args) > 0 && strings.ToUpper(args[0]) == facade.Mode_TERM.String() {
			args = args[1:]
		}

		modeFlags = flag.NewFlagSet("term", flag.ExitOnError)
		modeFlags.Usage = func() { ShowHelpMode(config.Mode, EXEC) }

		config.AddFlags(modeFlags)
		modeFlags.Parse(args)
		config.VisitFlags(modeFlags)

		args = modeFlags.Args()
		log.Debug("got flags %v", args)
		if len(args) <= 0 {
			ShowHelpMode(facade.Mode_TERM, EXEC)
			os.Exit(-2)
		}

		path = args[0]
		args = args[1:]

	case READ, RECV, PIPE, CONF, TEST:
		// parse mode, if given
		args = flags[cmd].Args()

		if cmd != INFO && len(args) > 0 {

			switch strings.ToUpper(args[0]) {

			case facade.Mode_TERM.String():
				config.SetMode = true
				config.Mode = facade.Mode_TERM
				config.Terminal = &facade.TermConfig{}
				config.Terminal.Grid = &facade.GridConfig{}

			case facade.Mode_LINE.String():
				config.SetMode = true
				config.Mode = facade.Mode_LINE
				config.Lines = &facade.LineConfig{}
				config.Lines.Grid = &facade.GridConfig{}

			case facade.Mode_DRAFT.String():
				config.SetMode = true
				config.Mode = facade.Mode_DRAFT

			default:
				ShowHelpCommand(cmd, *flags[cmd])
				os.Exit(-2)

			}

			args = args[1:]
			modeFlags = flag.NewFlagSet(strings.ToLower(config.Mode.String()), flag.ExitOnError)
			modeFlags.Usage = func() { ShowHelpMode(config.Mode, cmd) }
			modeFlags.Parse(args)
			config.VisitFlags(modeFlags)

		}

	case INFO:
		// no args, print local info
		if len(flag.Args()) <= 1 {
			ShowVersion()
			if log.InfoLogging() {
				ShowAssets()
			}
			fmt.Fprintf(os.Stderr, "\n\n")
			os.Exit(0)
		} else {

			// query remote host
			flags[INFO].Usage = func() { ShowHelpCommand(INFO, *flags[INFO]) }
			flags[INFO].Parse(flag.Args()[1:])

		}

	case HELP:
		ShowHelp()
		os.Exit(-2)

	default:
		ShowHelp()
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
			ShowHelpMode(facade.Mode_DRAFT, cmd)
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

func ShowHelpMode(mode facade.Mode, cmd Command) {
	ShowVersion()
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, "  %s %s %s [flags]\n", BUILD_NAME, cmd, strings.ToLower(mode.String()))

	fmt.Fprintf(os.Stderr, "\nFlags:\n")

	fmt.Fprintf(os.Stderr, "%s", facade.FontDefaults.Help())
	fmt.Fprintf(os.Stderr, "%s", facade.MaskDefaults.Help())
	fmt.Fprintf(os.Stderr, "%s", facade.CameraDefaults.Help())

	switch mode {
	case facade.Mode_LINE:
		fmt.Fprintf(os.Stderr, "%s", facade.LineDefaults.Help())
	case facade.Mode_TERM:
		fmt.Fprintf(os.Stderr, "%s", facade.TermDefaults.Help())
	}

	fmt.Fprintf(os.Stderr, "\n")
}

func ShowHelpCommand(cmd Command, flags flag.FlagSet) {
	var modes []string
	modes = append(modes, strings.ToLower(facade.Mode_TERM.String()))
	modes = append(modes, strings.ToLower(facade.Mode_LINE.String()))

	ShowVersion()
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	switch cmd {
	case INFO:
		fmt.Fprintf(os.Stderr, "  %s %s [flags]\n", BUILD_NAME, cmd)
	default:
		fmt.Fprintf(os.Stderr, "  %s %s [flags]  %s\n", BUILD_NAME, cmd, strings.Join(modes, " | "))
		ShowModes()
	}

	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	flags.VisitAll(func(f *flag.Flag) {
		name := f.Name
		if f.DefValue != "false" && f.DefValue != "true" {
			name = f.Name + "=" + f.DefValue
		}
		fmt.Fprintf(os.Stderr, "  -%-24s %-24s\n", name, f.Usage)
	})
	fmt.Fprintf(os.Stderr, "\n")
}

func ShowCommands() {
	fmt.Fprintf(os.Stderr, "\nCommands:\n")
	if RENDERER_AVAILABLE {
		fmt.Fprintf(os.Stderr, "%6s     %s\n", READ, "read text from stdin and render")
		fmt.Fprintf(os.Stderr, "%6s     %s\n", RECV, "receive text from client and render ")
	}
	fmt.Fprintf(os.Stderr, "%6s     %s\n", PIPE, "read text from stdin and send to server")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", CONF, "change server configuration")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", EXEC, "execute command and send stdio to server")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", INFO, "show available shaders and fonts of server ")
}

func ShowModes() {
	fmt.Fprintf(os.Stderr, "\nModes:\n")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", strings.ToLower(facade.Mode_TERM.String()), "text terminal")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", strings.ToLower(facade.Mode_LINE.String()), "line scroller")
}

func ShowHelp() {
	cmds := []string{}
	if RENDERER_AVAILABLE {
		for _, c := range []Command{READ, RECV} {
			cmds = append(cmds, string(c))
		}
	}
	for _, c := range []Command{PIPE, CONF, EXEC, INFO} {
		cmds = append(cmds, string(c))
	}

	ShowVersion()
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, "  %s [flags]  %s\n", BUILD_NAME, strings.Join(cmds, " | "))
	ShowCommands()
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	flag.VisitAll(func(f *flag.Flag) {
		name := f.Name
		if f.DefValue != "false" && f.DefValue != "true" {
			name = f.Name + "=" + f.DefValue
		}
		fmt.Fprintf(os.Stderr, "  -%-24s %-24s\n", name, f.Usage)
	})
	fmt.Fprintf(os.Stderr, "\n")

}

func ShowAssets() {

	//FIXME
	fmt.Fprintf(os.Stderr, InfoAssets(nil, nil))
}

func InfoAssets(shaders, fonts []string) string {
	ret := ""

	for _, prefix := range []string{"grid/"} {
		for _, suffix := range []string{".vert", ".frag"} {
			tmp := strings.TrimSuffix(prefix, "/") + " -" + strings.TrimPrefix(suffix, ".")
			ret += fmt.Sprintf("\n%12s= ", tmp)
			for _, shader := range shaders {
				if strings.HasPrefix(shader, prefix) && strings.HasSuffix(shader, suffix) {
					ret += strings.TrimSuffix(strings.TrimPrefix(shader, prefix), suffix)
					ret += " "
				}

			}
		}
	}

	ret += fmt.Sprintf("\n%12s= ", "-mask")
	for _, shader := range shaders {
		if strings.HasPrefix(shader, "mask/") && strings.HasSuffix(shader, "frag") {
			ret += strings.TrimSuffix(strings.TrimPrefix(shader, "mask/"), ".frag")
			ret += " "
		}
	}

	ret += fmt.Sprintf("\n%12s= ", "-font")
	for _, font := range fonts {
		ret += font
		ret += " "
	}
	ret += "\n"
	return ret
}

func ShowVersion() { fmt.Fprintf(os.Stderr, InfoVersion()) }

func InfoVersion() string {
	ret := ""
	ret += AUTHOR
	ret += fmt.Sprintf("\n%s version %s for %s built %s\n", BUILD_NAME, BUILD_VERSION, BUILD_PLATFORM, BUILD_DATE)
	return ret
}
