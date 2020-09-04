package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	facade "./facade"
	gfx "./gfx"
)

func ShowHelp(flags flag.FlagSet) {
	cmds := []string{}
	if RENDERER_AVAILABLE {
		cmds = append(cmds, string(SERVE))
	}
	for _, c := range []Command{PIPE, CONF, EXEC, README} {
		cmds = append(cmds, string(c))
	}

	ShowVersion()
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, "  %s [flags]  %s\n", BUILD_NAME, strings.Join(cmds, " | "))
	ShowCommands()
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	str := ""
	flags.VisitAll(func(f *flag.Flag) { str += gfx.FlagHelp(f) })
	fmt.Fprintf(os.Stderr, "%s\n", str)
}

func ShowHelpMode(cmd Command, mode facade.Mode, flags flag.FlagSet) {
	ShowVersion()
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	if cmd == EXEC {
		fmt.Fprintf(os.Stderr, "  %s %s term [flags] /path/to/executable [args]\n", BUILD_NAME, cmd)
	} else {
		fmt.Fprintf(os.Stderr, "  %s %s %s [flags]\n", BUILD_NAME, cmd, strings.ToLower(mode.String()))
	}

	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fmt.Fprintf(os.Stderr, "%s", facade.FontDefaults.Help())
	fmt.Fprintf(os.Stderr, "%s", facade.CameraDefaults.Help())
	fmt.Fprintf(os.Stderr, "%s", facade.MaskDefaults.Help())
	fmt.Fprintf(os.Stderr, "%s", facade.ShaderDefaults.Help())
	fmt.Fprintf(os.Stderr, "%s", facade.Defaults.Help())
	fmt.Fprintf(os.Stderr, "\nMode Flags:\n")

	switch mode {
	case facade.Mode_LINES:
		fmt.Fprintf(os.Stderr, "%s", facade.LineDefaults.Help())
	case facade.Mode_TERM:
		fmt.Fprintf(os.Stderr, "%s", facade.TermDefaults.Help())
	case facade.Mode_WORDS:
		fmt.Fprintf(os.Stderr, "%s", facade.WordDefaults.Help())
	case facade.Mode_TAGS:
		fmt.Fprintf(os.Stderr, "%s", facade.TagDefaults.Help())
	}

	fmt.Fprintf(os.Stderr, "\n")
}

func ShowHelpCommand(cmd Command, flags flag.FlagSet) {
	var modes []string
	modes = append(modes, strings.ToLower(facade.Mode_TERM.String()))
	modes = append(modes, strings.ToLower(facade.Mode_LINES.String()))
	modes = append(modes, strings.ToLower(facade.Mode_WORDS.String()))
	modes = append(modes, strings.ToLower(facade.Mode_TAGS.String()))

	ShowVersion()
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	switch cmd {
	case EXEC:
		fmt.Fprintf(os.Stderr, "  %s %s term [flags] /path/to/executable [args]\n", BUILD_NAME, cmd)
	default:
		fmt.Fprintf(os.Stderr, "  %s %s [flags]  %s\n", BUILD_NAME, cmd, strings.Join(modes, " | "))
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
		fmt.Fprintf(os.Stderr, "%6s     %s\n", SERVE, "receive text from client and render ")
	}
	fmt.Fprintf(os.Stderr, "%6s     %s\n", PIPE, "read text from stdin and send to server")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", CONF, "send configuration to server")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", EXEC, "execute command and send stdio to server")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", README, "show documentation")
}

func ShowModes() {
	fmt.Fprintf(os.Stderr, "\nModes:\n")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", strings.ToLower(facade.Mode_LINES.String()), "line scroller")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", strings.ToLower(facade.Mode_TERM.String()), "text terminal")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", strings.ToLower(facade.Mode_WORDS.String()), "word scroller")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", strings.ToLower(facade.Mode_TAGS.String()), "tag cloud")
}

func ShowAssets(directory string) {

	fontService := gfx.NewFontService(directory+"/font", facade.FontAsset)
	programService := gfx.NewProgramService(directory+"/shader", facade.ShaderAsset)

	fmt.Fprintf(os.Stderr, InfoAssets(programService.GetAvailableNames(), fontService.GetAvailableNames()))
}

func InfoAssets(shaders, fonts []string) string {
	ret := ""

	modes := map[facade.Mode]string{
		facade.Mode_LINES: "grid/",
		facade.Mode_TERM:  "grid/",
		facade.Mode_WORDS: "set/",
		facade.Mode_TAGS:  "set/",
	}

	ret += fmt.Sprintf("\nfacade conf <MODE> -font= ")
	for _, font := range fonts {
		ret += font
		ret += " "
	}

	ret += fmt.Sprintf("\nfacade conf <MODE> -mask= ")
	for _, shader := range shaders {
		if strings.HasPrefix(shader, "mask/") && strings.HasSuffix(shader, "frag") {
			ret += strings.TrimSuffix(strings.TrimPrefix(shader, "mask/"), ".frag")
			ret += " "
		}
	}

	for _, mode := range []facade.Mode{facade.Mode_LINES, facade.Mode_TERM, facade.Mode_WORDS, facade.Mode_TAGS} {

		prefix := modes[mode]
		for _, suffix := range []string{".vert", ".frag"} {
			tmp := fmt.Sprintf("facade conf %5s -%s", strings.ToLower(mode.String()), strings.TrimPrefix(suffix, "."))
			ret += fmt.Sprintf("\n%12s= ", tmp)
			for _, shader := range shaders {
				if strings.HasPrefix(shader, prefix) && strings.HasSuffix(shader, suffix) {
					ret += strings.TrimSuffix(strings.TrimPrefix(shader, prefix), suffix)
					ret += " "
				}

			}
		}
	}

	ret += "\n"
	return ret
}

func ShowVersion() {
	fmt.Fprintf(os.Stderr, "%s", AUTHOR)
	fmt.Fprintf(os.Stderr, "%s", InfoVersion())
}

func InfoVersion() string {
	ret := ""
	ret += fmt.Sprintf("\n%s version %s for %s built %s\n", BUILD_NAME, BUILD_VERSION, BUILD_PLATFORM, BUILD_DATE)
	return ret
}
