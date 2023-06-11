//
package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"flag"
	"strings"
)

var DEFAULT_MODE Mode = Mode_LINES

var Defaults = Config{
	Mode:  DEFAULT_MODE,
	Debug: false,
}

func NewConfig(mode Mode) *Config {
	ret := &Config{
		SetMode: true,
		Mode:    mode,
		Font:    &FontConfig{},
		Camera:  &CameraConfig{},
		Mask:    &MaskConfig{},
		Shader:  &ShaderConfig{},
	}
	switch mode {
	case Mode_LINES:
		ret.Lines = &LineConfig{}
	case Mode_TERM:
		ret.Term = &TermConfig{}
	case Mode_WORDS:
		ret.Words = &WordConfig{}
	case Mode_CHARS:
		ret.Chars = &CharConfig{}
	}
	return ret
}

func (config *Config) Desc() string {
	ret := "cfg["
	if config.GetSetMode() {
		ret += strings.ToLower(config.GetMode().String()) + " "
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
			ret += "debug "
		} else {
			ret += "nobug "
		}
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *Config) AddFlags(flagset *flag.FlagSet, mode Mode, basicOptions bool) {
	if !basicOptions {
		flagset.BoolVar(&config.Debug, "D", Defaults.Debug, "draw debug?")
	}
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
		term.AddFlags(flagset, basicOptions)
		shader.AddFlags(flagset, Mode_TERM, basicOptions)
	}
	if lines := config.GetLines(); lines != nil {
		lines.AddFlags(flagset, basicOptions)
		shader.AddFlags(flagset, Mode_LINES, basicOptions)
	}
	if words := config.GetWords(); words != nil {
		words.AddFlags(flagset, basicOptions)
		shader.AddFlags(flagset, Mode_WORDS, basicOptions)
	}
	if chars := config.GetChars(); chars != nil {
		chars.AddFlags(flagset, basicOptions)
		shader.AddFlags(flagset, Mode_CHARS, basicOptions)
	}
	if font := config.GetFont(); font != nil {
		font.AddFlags(flagset, basicOptions)
	}
	if cam := config.GetCamera(); cam != nil {
		cam.AddFlags(flagset, basicOptions)
	}
	if mask := config.GetMask(); mask != nil {
		mask.AddFlags(flagset, basicOptions)
	}

	if !basicOptions {
		flagset.StringVar(&config.Fill, "fill", "", "fill pattern ("+strings.Join(patterns, ", ")+")")
	}

}

func (config *Config) VisitFlags(flagset *flag.FlagSet, mode Mode, basicOptions bool) {

	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "D":
			config.SetDebug = true
		case "fill":
			config.SetFill = true
		}
	})

	if term := config.GetTerm(); term != nil {
		if !term.VisitFlags(flagset, basicOptions) {
			config.Term = nil
		} // no flags used
	}
	if lines := config.GetLines(); lines != nil {
		if !lines.VisitFlags(flagset, basicOptions) {
			config.Lines = nil
		} // no flags used
	}
	if words := config.GetWords(); words != nil {
		if !words.VisitFlags(flagset, basicOptions) {
			config.Words = nil
		}
	}
	if chars := config.GetChars(); chars != nil {
		if !chars.VisitFlags(flagset, basicOptions) {
			config.Chars = nil
		} // no flags used
	}

	if font := config.GetFont(); font != nil {
		if !font.VisitFlags(flagset, basicOptions) {
			config.Font = nil
		} // no flags used
	}

	if cam := config.GetCamera(); cam != nil {
		if !cam.VisitFlags(flagset, basicOptions) {
			config.Camera = nil
		} // no flags used
	}

	if mask := config.GetMask(); mask != nil {
		if !mask.VisitFlags(flagset, basicOptions) {
			config.Mask = nil
		} // no flags used
	}

	if shader := config.GetShader(); shader != nil {
		if !shader.VisitFlags(flagset, mode, basicOptions) {
			config.Shader = nil
		} // no flags used
	}

}

func (config *Config) Help(mode Mode, basicOptions bool) string {
	ret := ""

	switch mode {
	case Mode_LINES:
		ret += LineDefaults.Help(basicOptions)
	case Mode_WORDS:
		ret += WordDefaults.Help(basicOptions)
	case Mode_CHARS:
		ret += CharDefaults.Help(basicOptions)
	case Mode_TERM:
		ret += TermDefaults.Help(basicOptions)
	}

	ret += "\n"
	ret += ShaderDefaults.Help(mode, basicOptions)
	ret += "\n"

	ret += FontDefaults.Help(basicOptions)
	ret += CameraDefaults.Help(basicOptions)
	ret += MaskDefaults.Help(basicOptions)

	tmp := flag.NewFlagSet("facade", flag.ExitOnError)
	config.AddFlags(tmp, mode, basicOptions)
	for _, s := range []string{"fill", "D"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}

	ret += "\n"

	return ret
}
