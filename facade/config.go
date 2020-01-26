//
package facade

import (
	"flag"
	"fmt"
	"strings"
)

var DEFAULT_DIRECTORY = "~/src/gfx/facade"

var DEFAULT_MODE Mode = Mode_LINE

var Defaults = Config{
	Mode:  DEFAULT_MODE,
	Debug: false,
}

func (config *Config) Desc() string {
	ret := "config["
	if config.GetSetMode() {
		ret += "mode=" + strings.ToLower(config.GetMode().String()) + " "
	}

	if font := config.GetFont(); font != nil {
		ret += font.Desc() + " "
	}
	if cam := config.GetCamera(); cam != nil {
		ret += cam.Desc() + " "
	}
	if mask := config.GetMask(); mask != nil {
		ret += mask.Desc() + " "
	}

	if terminal := config.GetTerminal(); terminal != nil {
		ret += terminal.Desc() + " "
	}
	if lines := config.GetLines(); lines != nil {
		ret += lines.Desc() + " "
	}
	if tags := config.GetTags(); tags != nil {
		ret += tags.Desc()
	}

	if config.GetSetDebug() {
		if config.GetDebug() {
			ret += "DEBUG "
		} else {
			ret += "nobug "
		}
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *Config) AddFlags(flagset *flag.FlagSet) {
	flagset.BoolVar(&config.Debug, "D", Defaults.Debug, "draw debug?")
	if tags := config.GetTags(); tags != nil {
		tags.AddFlags(flagset)
	}
	if terminal := config.GetTerminal(); terminal != nil {
		terminal.AddFlags(flagset)
	}
	if lines := config.GetLines(); lines != nil {
		lines.AddFlags(flagset)
	}
	if font := config.GetFont(); font != nil {
		font.AddFlags(flagset)
	}
	if cam := config.GetCamera(); cam != nil {
		cam.AddFlags(flagset)
	}
	if mask := config.GetMask(); mask != nil {
		mask.AddFlags(flagset)
	}
}

func (config *Config) VisitFlags(flagset *flag.FlagSet) {

	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "D":
			{
				config.SetDebug = true
			}
		}
	})

	if terminal := config.GetTerminal(); terminal != nil {
		if !terminal.VisitFlags(flagset) {
			config.Terminal = nil
		} // no flags used
	}

	if lines := config.GetLines(); lines != nil {
		if !lines.VisitFlags(flagset) {
			config.Lines = nil
		} // no flags used
	}

	if tags := config.GetTags(); tags != nil {
		if !tags.VisitFlags(flagset) {
			config.Tags = nil
		}

	}

	if font := config.GetFont(); font != nil {
		if !font.VisitFlags(flagset) {
			config.Font = nil
		} // no flags used
	}

	if cam := config.GetCamera(); cam != nil {
		if !cam.VisitFlags(flagset) {
			config.Camera = nil
		} // no flags used
	}

	if mask := config.GetMask(); mask != nil {
		if !mask.VisitFlags(flagset) {
			config.Mask = nil
		} // no flags used
	}

}

func (config *Config) Help() string {
	ret := ""
	fun := func(f *flag.Flag) {
		name := f.Name
		if f.DefValue != "false" && f.DefValue != "true" {
			name = f.Name + "=" + f.DefValue
		}
		ret += fmt.Sprintf("  -%-24s %-24s\n", name, f.Usage)
	}

	tmp := flag.NewFlagSet("facade", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll(fun)
	return ret
}
