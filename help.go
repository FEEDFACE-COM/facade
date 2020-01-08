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
	flags.VisitAll(func(f *flag.Flag) {
		name := f.Name
		if f.DefValue != "false" && f.DefValue != "true" {
			name = f.Name + "=" + f.DefValue
		}
		fmt.Fprintf(os.Stderr, "  -%-24s %-24s\n", name, f.Usage)
	})
	fmt.Fprintf(os.Stderr, "\n")
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
	fmt.Fprintf(os.Stderr, "%s", facade.MaskDefaults.Help())
	fmt.Fprintf(os.Stderr, "%s", facade.CameraDefaults.Help())
	fmt.Fprintf(os.Stderr, "%s", facade.Defaults.Help())

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
    case EXEC:
    	fmt.Fprintf(os.Stderr, "  %s %s term [flags] /path/to/executable [args]\n", BUILD_NAME, cmd)    
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
	fmt.Fprintf(os.Stderr, "%6s     %s\n", CONF, "send configuration to server")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", EXEC, "execute command and send stdio to server")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", INFO, "show available shaders and fonts of server ")
}

func ShowModes() {
	fmt.Fprintf(os.Stderr, "\nModes:\n")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", strings.ToLower(facade.Mode_LINE.String()), "line scroller")
	fmt.Fprintf(os.Stderr, "%6s     %s\n", strings.ToLower(facade.Mode_TERM.String()), "text terminal")
}

func ShowAssets(directory string) {

	fontService := gfx.NewFontService(directory+"/font", facade.FontAsset)
	programService := gfx.NewProgramService(directory+"/shader", facade.ShaderAsset)

	fmt.Fprintf(os.Stderr, InfoAssets(programService.GetAvailableNames(), fontService.GetAvailableNames()))
}

func InfoAssets(shaders, fonts []string) string {
	ret := ""

	ret += fmt.Sprintf("\n%12s= ", "-font")
	for _, font := range fonts {
		ret += font
		ret += " "
	}

	ret += fmt.Sprintf("\n%12s= ", "-mask")
	for _, shader := range shaders {
		if strings.HasPrefix(shader, "mask/") && strings.HasSuffix(shader, "frag") {
			ret += strings.TrimSuffix(strings.TrimPrefix(shader, "mask/"), ".frag")
			ret += " "
		}
	}

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
