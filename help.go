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
	for _, c := range []Command{PIPE, CONF, EXEC /*README,*/} {
		cmds = append(cmds, string(c))
	}

	fmt.Fprintf(os.Stderr, InfoAuthor())
	fmt.Fprintf(os.Stderr, InfoVersion())
	fmt.Fprintf(os.Stderr, "\nusage:\n")
	fmt.Fprintf(os.Stderr, "  %s [flags]  %s\n", BUILD_NAME, strings.Join(cmds, "|"))
	ShowCommands()
	fmt.Fprintf(os.Stderr, "\nflags:\n")
	str := ""
	for _,s := range[]string{"q","d", "D"} {
		if flg := flags.Lookup(s); flg != nil {
			str += gfx.FlagHelp(flg)
		}
	}
	fmt.Fprintf(os.Stderr, "%s\n", str)
}

func ShowHelpMode(cmd Command, mode facade.Mode, basicOptions bool) {
	fmt.Fprintf(os.Stderr, InfoAuthor())
	fmt.Fprintf(os.Stderr, InfoVersion())
	fmt.Fprintf(os.Stderr, "\nusage:\n")
	if cmd == EXEC {
		fmt.Fprintf(os.Stderr, "  %s exec term [flags] /path/to/executable [args]\n", BUILD_NAME)
	} else {
		fmt.Fprintf(os.Stderr, "  %s %s %s [flags]\n", BUILD_NAME, cmd, strings.ToLower(mode.String()))
	}

	fmt.Fprintf(os.Stderr, "\nflags:\n")
	fmt.Fprintf(os.Stderr, facade.Defaults.Help(mode,basicOptions) )

	fmt.Fprintf(os.Stderr, "\n")
}

func ShowHelpCommand(cmd Command, flags flag.FlagSet) {
	var modes = []string{
		strings.ToLower(facade.Mode_LINES.String()),
		strings.ToLower(facade.Mode_CHARS.String()),
		strings.ToLower(facade.Mode_WORDS.String()),
		strings.ToLower(facade.Mode_TERM.String()),
	}
	fmt.Fprintf(os.Stderr, InfoAuthor())
	fmt.Fprintf(os.Stderr, InfoVersion())
	fmt.Fprintf(os.Stderr, "\nusage:\n")
	switch cmd {
	case EXEC:
		fmt.Fprintf(os.Stderr, "  %s exec [flags] term /path/to/executable [args]\n", BUILD_NAME)
	default:
		fmt.Fprintf(os.Stderr, "  %s %s [flags] %s\n", BUILD_NAME, cmd, strings.Join(modes, "|"))
		ShowModes()
	}
	fmt.Fprintf(os.Stderr, "\nflags:\n")
	str := ""
	flags.VisitAll(func(f *flag.Flag) { str += gfx.FlagHelp(f) })
	fmt.Fprintf(os.Stderr, "%s\n", str)
}

func ShowCommands() {
	fmt.Fprintf(os.Stderr, "\ncommands:\n")
	if RENDERER_AVAILABLE {
		fmt.Fprintf(os.Stderr, "  %6s     %s\n", SERVE, "receive text from network and draw to screen")
	}
	fmt.Fprintf(os.Stderr, "%  6s     %s\n", PIPE, "read text from stdin and send to renderer")
	fmt.Fprintf(os.Stderr, "%  6s     %s\n", CONF, "send config to renderer")
	fmt.Fprintf(os.Stderr, "%  6s     %s\n", EXEC, "execute command and send stdout/stderr to renderer")
	//	fmt.Fprintf(os.Stderr, "%6s     %s\n", README, "show documentation")
}

func ShowModes() {
	fmt.Fprintf(os.Stderr, "\nmodes:\n")
	fmt.Fprintf(os.Stderr, "  %5s     %s\n", strings.ToLower(facade.Mode_LINES.String()), "scrolling lines of text")
	fmt.Fprintf(os.Stderr, "  %5s     %s\n", strings.ToLower(facade.Mode_CHARS.String()), "scrolling letters")
	fmt.Fprintf(os.Stderr, "  %5s     %s\n", strings.ToLower(facade.Mode_WORDS.String()), "fading words")
	fmt.Fprintf(os.Stderr, "  %5s     %s\n", strings.ToLower(facade.Mode_TERM.String()), "ansi terminal")
}

func InfoAuthor() string {
	return fmt.Sprintf("%s\n", strings.TrimLeft(AUTHOR, "\n"))
}

func InfoVersion() string {
	return fmt.Sprintf("%s version %s for %s built %s\n", BUILD_NAME, BUILD_VERSION, BUILD_PLATFORM, BUILD_DATE)

}
