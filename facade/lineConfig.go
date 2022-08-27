//
//
package facade

//
import (
	"FEEDFACE.COM/facade/gfx"
	"FEEDFACE.COM/facade/log"
	"FEEDFACE.COM/facade/math32"
	"flag"
	"fmt"
	"strings"
)

var LineDefaults LineConfig = LineConfig{
	Width:    32,
	Height:   8,
	Downward: false,
	Speed:    1.0,
	Fixed:    false,
	Drop:     true,
	Smooth:   true,
	Buffer:   8,
}

func (config *LineConfig) Desc() string {
	ret := "line["

	{
		wok := config.GetSetWidth()
		hok := config.GetSetHeight()
		if wok {
			ret += fmt.Sprintf("%d", config.GetWidth())
		}
		if wok || hok {
			ret += "x"
		}
		if hok {
			ret += fmt.Sprintf("%d", config.GetHeight())
		}
		if wok || hok {
			ret += " "
		}
	}

	{
		bok := config.GetSetBuffer()
		if bok {
			ret += fmt.Sprintf("+%d ", config.GetBuffer())
		}
	}

	{
		if config.GetSetDownward() {
			if config.GetDownward() {
				ret += "↓ "
			} else {
				ret += "↑ "
			}
		}
		if config.GetSetFixed() {
			if !config.GetFixed() {
				ret += "!"
			}
			ret += "ƒ "
		}

		if config.GetSetDrop() {
			if !config.GetDrop() {
				ret += "!"
			}
			ret += "ṕ "
		}

		if config.GetSetSmooth() {
			if !config.GetSmooth() {
				ret += "!"
			}
			ret += "ß "
		}

		if config.GetSetSpeed() {
			ret += fmt.Sprintf("s%.1f ", config.GetSpeed())
		}
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *LineConfig) FillPatterns() []string {
	return []string{"title", "index", "alpha", "clear"}
}

func (config *LineConfig) AddFlags(flagset *flag.FlagSet) {
	flagset.Uint64Var(&config.Width, "w", LineDefaults.Width, "width: chars per line")
	flagset.Uint64Var(&config.Height, "h", LineDefaults.Height, "height: line count")
	flagset.BoolVar(&config.Downward, "down", LineDefaults.Downward, "scroll downward?")
	flagset.BoolVar(&config.Drop, "drop", LineDefaults.Drop, "drop lines?")
	flagset.BoolVar(&config.Smooth, "smooth", LineDefaults.Smooth, "scroll smooth?")
	flagset.Float64Var(&config.Speed, "speed", LineDefaults.Speed, "scroll speed")
	flagset.BoolVar(&config.Fixed, "fixed", LineDefaults.Fixed, "fixed scroll speed?")
	flagset.Uint64Var(&config.Buffer, "buffer", LineDefaults.Buffer, "buffer length")
}

func (config *LineConfig) VisitFlags(flagset *flag.FlagSet) bool {
	ret := false
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "w":
			{
				config.SetWidth = true
				ret = true
			}
		case "h":
			{
				config.SetHeight = true
				ret = true
			}
		case "down":
			{
				config.SetDownward = true
				ret = true
			}
		case "drop":
			{
				config.SetDrop = true
				ret = true
			}
		case "smooth":
			{
				config.SetSmooth = true
				ret = true
			}
		case "speed":
			{
				config.SetSpeed = true
				ret = true
			}
		case "fixed":
			{
				config.SetFixed = true
				ret = true
			}
		case "buffer":
			{
				config.SetBuffer = true
				ret = true
			}
		}
	})

	return ret

}

func (config *LineConfig) Help() string {
	ret := ""
	tmp := flag.NewFlagSet("lines", flag.ExitOnError)
	config.AddFlags(tmp)
	for _, s := range []string{"w", "h", "down", "drop", "smooth", "speed", "fixed", "buffer"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}
	return ret
}

func (config *LineConfig) autoSize(cameraRatio float32, fontRatio float32) {

	if config.GetSetHeight() && !config.GetSetWidth() {

		height := config.GetHeight()
		w := math32.Round((cameraRatio / fontRatio) * float32(height))
		if height == 1 { //special case
			w = 6.
		}

		config.SetWidth = true
		config.Width = uint64(w)
		log.Info("%s calculated %.0f width", config.Desc(), w)

	} else if config.GetSetWidth() && !config.GetSetHeight() {

		width := config.GetWidth()
		h := math32.Round(float32(width) / (cameraRatio / fontRatio))
		if width <= 6 { //special case
			h = 1.
		}

		config.SetHeight = true
		config.Height = uint64(h)
		log.Info("%s calculated %.0f height", config.Desc(), h)
	}
}
