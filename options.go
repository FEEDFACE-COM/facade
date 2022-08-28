package main

import (
	"FEEDFACE.COM/facade/facade"
	"FEEDFACE.COM/facade/gfx"
	"flag"
	"os"
)

// options provide a minimal set of flags and create a config from it

type Options struct {
	Mode      facade.Mode
	StyleName string
	FontName  string
	Mask      bool
	Zoom      float64

	Speed float64

	//lines+term
	LinesWidth    uint
	LinesHeight   uint
	LinesDownward bool

	// chars
	CharsCount uint

	//words
	WordsCount    uint
	WordsShuffle  bool
	WordsFillMark float64
	WordsLifeTime float64
}

var DefaultOptions = Options{
	Mode:      facade.DEFAULT_MODE,
	StyleName: facade.ShaderDefaults.Vert,
	FontName:  facade.FontDefaults.Name,
	Mask:      facade.MaskDefaults.Name == "mask",
	Zoom:      facade.CameraDefaults.Zoom,

	Speed: facade.LineDefaults.Speed,

	LinesWidth:    uint(facade.LineDefaults.Width),
	LinesHeight:   uint(facade.LineDefaults.Height),
	LinesDownward: facade.LineDefaults.Downward,

	CharsCount: uint(facade.CharDefaults.CharCount),

	WordsCount:    uint(facade.WordDefaults.Slots),
	WordsShuffle:  facade.WordDefaults.Shuffle,
	WordsFillMark: facade.WordDefaults.Watermark,
	WordsLifeTime: facade.WordDefaults.Lifetime,
}

func NewOptions(mode facade.Mode) *Options {
	var ret = DefaultOptions
	ret.Mode = mode
	return &ret
}

func (options *Options) AddFlags(flagset *flag.FlagSet, mode facade.Mode) {
	flagset.StringVar(&options.FontName, "font", DefaultOptions.FontName, "typeface ("+facade.AvailableFonts()+")")
	flagset.BoolVar(&options.Mask, "mask", DefaultOptions.Mask, "overlay mask?")
	flagset.Float64Var(&options.Zoom, "zoom", DefaultOptions.Zoom, "camera zoom")
	flagset.StringVar(&options.StyleName, "shape", DefaultOptions.StyleName, "shape ("+facade.AvailableShaders(facade.PrefixForMode(mode), ".vert")+")")

	switch mode {
	case facade.Mode_LINES:
		flagset.UintVar(&options.LinesWidth, "w", DefaultOptions.LinesWidth, "width: chars per line")
		flagset.UintVar(&options.LinesHeight, "h", DefaultOptions.LinesHeight, "height: line count")
		flagset.Float64Var(&options.Speed, "speed", DefaultOptions.Speed, "scroll speed")
		flagset.BoolVar(&options.LinesDownward, "down", DefaultOptions.LinesDownward, "scroll downward?")
	case facade.Mode_TERM:
		flagset.UintVar(&options.LinesWidth, "w", DefaultOptions.LinesWidth, "terminal width")
		flagset.UintVar(&options.LinesHeight, "h", DefaultOptions.LinesHeight, "terminal height")
	case facade.Mode_CHARS:
		flagset.UintVar(&options.CharsCount, "w", DefaultOptions.CharsCount, "width: chars in line")
		flagset.Float64Var(&options.Speed, "speed", DefaultOptions.Speed, "scroll speed")
	case facade.Mode_WORDS:
		flagset.UintVar(&options.WordsCount, "n", DefaultOptions.WordsCount, "number of word slots")
		flagset.Float64Var(&options.WordsLifeTime, "life", DefaultOptions.WordsLifeTime, "word lifetime")
		flagset.Float64Var(&options.WordsFillMark, "mark", DefaultOptions.WordsFillMark, "buffer fill mark")
		flagset.BoolVar(&options.WordsShuffle, "shuffle", DefaultOptions.WordsShuffle, "shuffle words?")

	}

}

func (options *Options) VisitFlags(cmd Command, flagset *flag.FlagSet) *facade.Config {
	ret := &facade.Config{SetMode: true, Mode: options.Mode}

	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "font":
			if _, ok := facade.FontAsset[options.FontName]; ok {
				ret.Font = &facade.FontConfig{SetName: true, Name: options.FontName}
			} else {

			}
		case "mask":
			if options.Mask {
				ret.Mask = &facade.MaskConfig{SetName: true, Name: "mask"}
			} else {
				ret.Mask = &facade.MaskConfig{SetName: true, Name: facade.ShaderDefaults.Frag}
			}
		case "zoom":
			ret.Camera = &facade.CameraConfig{SetZoom: true, Zoom: options.Zoom}
		case "shape":
			vert := facade.PrefixForMode(options.Mode) + options.StyleName + "." + string(gfx.VertType)
			frag := facade.PrefixForMode(options.Mode) + options.StyleName + "." + string(gfx.FragType)
			ret.Shader = &facade.ShaderConfig{SetVert: true, SetFrag: true}
			if _, ok := facade.ShaderAsset[vert]; ok {
				ret.Shader.Vert = options.StyleName
				if _, ok := facade.ShaderAsset[frag]; ok {
					ret.Shader.Frag = options.StyleName
				} else {
					ret.Shader.Frag = facade.ShaderDefaults.Frag
				}
			} else {
				ShowHelpMode(cmd, options.Mode, *flagset, nil, options)
				os.Exit(-1)

			}

		}

	})

	switch options.Mode {
	case facade.Mode_LINES:
		ret.Lines = &facade.LineConfig{}
		flagset.Visit(func(flg *flag.Flag) {
			switch flg.Name {
			case "w":
				ret.Lines.Width = uint64(options.LinesWidth)
				ret.Lines.SetWidth = true
			case "h":
				ret.Lines.Height = uint64(options.LinesHeight)
				ret.Lines.SetHeight = true
			case "speed":
				ret.Lines.Speed = options.Speed
				ret.Lines.SetSpeed = true
			case "down":
				ret.Lines.Downward = options.LinesDownward
				ret.Lines.SetDownward = true
			}
		})
	case facade.Mode_CHARS:
		ret.Chars = &facade.CharConfig{}
		flagset.Visit(func(flg *flag.Flag) {
			switch flg.Name {
			case "w":
				ret.Chars.CharCount = uint64(options.CharsCount)
				ret.Chars.SetCharCount = true
			case "speed":
				ret.Chars.Speed = options.Speed
				ret.Chars.SetSpeed = true
			}
		})
	case facade.Mode_WORDS:
		ret.Words = &facade.WordConfig{}
		flagset.Visit(func(flg *flag.Flag) {
			switch flg.Name {
			case "n":
				ret.Words.Slots = uint64(options.WordsCount)
				ret.Words.SetSlots = true
			case "life":
				ret.Words.Lifetime = options.WordsLifeTime
				ret.Words.SetLifetime = true
			case "mark":
				ret.Words.Watermark = options.WordsFillMark
				ret.Words.SetWatermark = true
			case "shuffle":
				ret.Words.Shuffle = options.WordsShuffle
				ret.Words.SetShuffle = true
			}
		})
	case facade.Mode_TERM:
		ret.Term = &facade.TermConfig{}
		ret.Lines = &facade.LineConfig{}
		flagset.Visit(func(flg *flag.Flag) {
			switch flg.Name {
			case "w":
				ret.Term.Width = uint64(options.LinesWidth)
				ret.Term.SetWidth = true
			case "h":
				ret.Term.Height = uint64(options.LinesHeight)
				ret.Term.SetHeight = true
			}

		})
	}

	return ret
}

func (options *Options) Help(mode facade.Mode) string {
	ret := ""
	tmp := flag.NewFlagSet("facade", flag.ExitOnError)
	options.AddFlags(tmp, mode)
	ret += gfx.FlagHelp(tmp.Lookup("font"))
	ret += gfx.FlagHelp(tmp.Lookup("zoom"))
	ret += gfx.FlagHelp(tmp.Lookup("mask"))
	ret += "\n"
	ret += gfx.FlagHelp(tmp.Lookup("shape"))
	ret += "\n"
	switch mode {
	case facade.Mode_TERM:
		ret += gfx.FlagHelp(tmp.Lookup("w"))
		ret += gfx.FlagHelp(tmp.Lookup("h"))
	case facade.Mode_LINES:
		ret += gfx.FlagHelp(tmp.Lookup("w"))
		ret += gfx.FlagHelp(tmp.Lookup("h"))
		ret += gfx.FlagHelp(tmp.Lookup("speed"))
		ret += gfx.FlagHelp(tmp.Lookup("down"))
	case facade.Mode_CHARS:
		ret += gfx.FlagHelp(tmp.Lookup("w"))
		ret += gfx.FlagHelp(tmp.Lookup("speed"))
	case facade.Mode_WORDS:
		ret += gfx.FlagHelp(tmp.Lookup("n"))
		ret += gfx.FlagHelp(tmp.Lookup("life"))
		ret += gfx.FlagHelp(tmp.Lookup("mark"))
		ret += gfx.FlagHelp(tmp.Lookup("shuffle"))
	}
	ret += "\n"
	return ret
}
