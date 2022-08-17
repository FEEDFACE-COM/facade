//
package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
	"strings"
)

var DEFAULT_DIRECTORY = "."
var DEFAULT_MODE Mode = Mode_LINES

var Defaults = Config{
	Mode:  DEFAULT_MODE,
	Debug: false,
}

func (config *Config) Desc() string {
	ret := "cfg["
	if config.GetSetMode() {
		ret += strings.ToUpper(config.GetMode().String()) + " "
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
	if shader := config.GetShader(); shader != nil {
		ret += shader.Desc() + " "
	}
	if term := config.GetTerm(); term != nil {
		ret += term.Desc() + " "
	}
	if lines := config.GetLines(); lines != nil {
		ret += lines.Desc() + " "
	}
	if words := config.GetWords(); words != nil {
		ret += words.Desc()
	}
	if chars := config.GetChars(); chars != nil {
		ret += chars.Desc()
	}
	if config.GetSetFill() {
		ret += "f:" + config.GetFill() + " "
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

func (config *Config) AddFlags(flagset *flag.FlagSet, mode Mode) {
	flagset.BoolVar(&config.Debug, "D", Defaults.Debug, "draw debug?")
	var shader *ShaderConfig = config.GetShader()
	var patterns = []string{}
	switch mode {
	case Mode_LINES:
		patterns = LineDefaults.FillPatterns()
	case Mode_WORDS:
		patterns = WordDefaults.FillPatterns()
	case Mode_CHARS:
		patterns = CharDefaults.FillPatterns()
	case Mode_TERM:
		patterns = TermDefaults.FillPatterns()
	}

	if term := config.GetTerm(); term != nil {
		term.AddFlags(flagset)
		shader.AddFlags(flagset, Mode_TERM)
	}
	if lines := config.GetLines(); lines != nil {
		lines.AddFlags(flagset)
		shader.AddFlags(flagset, Mode_LINES)
	}
	if words := config.GetWords(); words != nil {
		words.AddFlags(flagset)
		shader.AddFlags(flagset, Mode_WORDS)
	}
	if chars := config.GetChars(); chars != nil {
		chars.AddFlags(flagset)
		shader.AddFlags(flagset, Mode_CHARS)
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

	flagset.StringVar(&config.Fill, "fill", "", "fill pattern ("+strings.Join(patterns, ",")+")")

}

func (config *Config) VisitFlags(flagset *flag.FlagSet) {

	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "D":
			config.SetDebug = true
		case "fill":
			config.SetFill = true
		}
	})

	if term := config.GetTerm(); term != nil {
		if !term.VisitFlags(flagset) {
			config.Term = nil
		} // no flags used
	}
	if lines := config.GetLines(); lines != nil {
		if !lines.VisitFlags(flagset) {
			config.Lines = nil
		} // no flags used
	}
	if words := config.GetWords(); words != nil {
		if !words.VisitFlags(flagset) {
			config.Words = nil
		}
	}
	if chars := config.GetChars(); chars != nil {
		if !chars.VisitFlags(flagset) {
			config.Chars = nil
		} // no flags used
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

	if shader := config.GetShader(); shader != nil {
		if !shader.VisitFlags(flagset) {
			config.Shader = nil
		} // no flags used
	}

}

func (config *Config) Help(mode Mode) string {
	ret := ""
	tmp := flag.NewFlagSet("facade", flag.ExitOnError)
	config.AddFlags(tmp, mode)
	for _, s := range []string{"fill", "D"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}
	return ret
}
