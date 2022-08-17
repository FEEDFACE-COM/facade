package main

import (
	"FEEDFACE.COM/facade/facade"
	"FEEDFACE.COM/facade/gfx"
	"flag"
	"fmt"
	"os"
	"strings"
)

func ShowHelp(flags flag.FlagSet) {
	cmds := []string{}
	if RENDERER_AVAILABLE {
		cmds = append(cmds, string(SERVE))
	}
	for _, c := range []Command{PIPE, CONF, EXEC, README} {
		cmds = append(cmds, string(c))
	}

	fmt.Fprintf(os.Stderr, InfoAuthor())
	fmt.Fprintf(os.Stderr, InfoVersion())
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, "  %s [flags]  %s\n", BUILD_NAME, strings.Join(cmds, "|"))
	ShowCommands()
	fmt.Fprintf(os.Stderr, "\nCommand Flags:\n")
	str := ""
	flags.VisitAll(func(f *flag.Flag) { str += gfx.FlagHelp(f) })
	fmt.Fprintf(os.Stderr, "%s\n", str)
}

func ShowHelpMode(cmd Command, mode facade.Mode, flags flag.FlagSet) {
	fmt.Fprintf(os.Stderr, InfoAuthor())
	fmt.Fprintf(os.Stderr, InfoVersion())
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	if cmd == EXEC {
		fmt.Fprintf(os.Stderr, "  %s exec [flags] /path/to/executable [args]\n", BUILD_NAME)
	} else {
		fmt.Fprintf(os.Stderr, "  %s %s %s [flags]\n", BUILD_NAME, cmd, strings.ToLower(mode.String()))
	}

	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fmt.Fprintf(os.Stderr, "%s", facade.FontDefaults.Help())
	fmt.Fprintf(os.Stderr, "%s", facade.CameraDefaults.Help())
	fmt.Fprintf(os.Stderr, "%s", facade.MaskDefaults.Help())
	fmt.Fprintf(os.Stderr, "%s", facade.Defaults.Help(mode))
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "%s", facade.ShaderDefaults.Help(mode))
	fmt.Fprintf(os.Stderr, "\n")

	switch mode {
	case facade.Mode_LINES:
		fmt.Fprintf(os.Stderr, "%s", facade.LineDefaults.Help())
	case facade.Mode_TERM:
		fmt.Fprintf(os.Stderr, "%s", facade.TermDefaults.Help())
	case facade.Mode_WORDS:
		fmt.Fprintf(os.Stderr, "%s", facade.WordDefaults.Help())
	case facade.Mode_CHARS:
		fmt.Fprintf(os.Stderr, "%s", facade.CharDefaults.Help())
	}

	fmt.Fprintf(os.Stderr, "\n")
}

func ShowHelpCommand(cmd Command, flags flag.FlagSet) {
	var modes []string
	modes = append(modes, strings.ToLower(facade.Mode_TERM.String()))
	modes = append(modes, strings.ToLower(facade.Mode_LINES.String()))
	modes = append(modes, strings.ToLower(facade.Mode_WORDS.String()))
	modes = append(modes, strings.ToLower(facade.Mode_CHARS.String()))
	fmt.Fprintf(os.Stderr, InfoAuthor())
	fmt.Fprintf(os.Stderr, InfoVersion())
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	switch cmd {
	case EXEC:
		fmt.Fprintf(os.Stderr, "  %s exec [flags] /path/to/executable [args]\n", BUILD_NAME)
	default:
		fmt.Fprintf(os.Stderr, "  %s %s [flags] %s\n", BUILD_NAME, cmd, strings.Join(modes, " | "))
		ShowModes()
	}
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	str := ""
	flags.VisitAll(func(f *flag.Flag) { str += gfx.FlagHelp(f) })
	fmt.Fprintf(os.Stderr, "%s\n", str)
}

func ShowCommands() {
	fmt.Fprintf(os.Stderr, "\nCommands:\n")
	if RENDERER_AVAILABLE {
		fmt.Fprintf(os.Stderr, "%6s     %s\n", SERVE, "receive text from network and render")
	}
	fmt.Fprintf(os.Stderr, "%6s     %s\n", PIPE, "read text from stdin and send to network")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", CONF, "send config to server")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", EXEC, "execute command and send stdout/stderr to server")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", README, "show documentation")
}

func ShowModes() {
	fmt.Fprintf(os.Stderr, "\nModes:\n")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", strings.ToLower(facade.Mode_LINES.String()), "scrolling lines of text")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", strings.ToLower(facade.Mode_CHARS.String()), "scrolling letters")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", strings.ToLower(facade.Mode_WORDS.String()), "fading words")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", strings.ToLower(facade.Mode_TERM.String()), "ansi terminal")
}

func InfoAuthor() string {
	return fmt.Sprintf("%s\n", strings.TrimLeft(AUTHOR,"\n") )
}

func InfoVersion() string {
	return fmt.Sprintf("%s version %s for %s built %s\n", BUILD_NAME, BUILD_VERSION, BUILD_PLATFORM, BUILD_DATE)

}
