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

var TermDefaults = TermConfig{
	Width:  32,
	Height: 8,
}

func (config *TermConfig) Desc() string {
	ret := "term["
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
	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *TermConfig) FillPatterns() []string {
	return []string{"title", "index", "alpha", "clear"}
}

func (config *TermConfig) AddFlags(flagset *flag.FlagSet, basicFlags bool) {
	flagset.Uint64Var(&config.Width, "w", TermDefaults.Width, "terminal width")
	flagset.Uint64Var(&config.Height, "h", TermDefaults.Height, "terminal height")
}

func (config *TermConfig) VisitFlags(flagset *flag.FlagSet, basicFlags bool) bool {
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
		}
	})
	return ret

}

func (config *TermConfig) Help(basicOptions bool) string {
	ret := ""
	tmp := flag.NewFlagSet("term", flag.ExitOnError)
	config.AddFlags(tmp,basicOptions)
	for _, s := range []string{"w", "h"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}
	return ret
}

func (config *TermConfig) autoSize(cameraRatio float32, fontRatio float32) {

	if config.GetSetHeight() && !config.GetSetWidth() {

		height := config.GetHeight()
		w := math32.Round((cameraRatio / fontRatio) * float32(height))
		if height == 1 { //special case
			w = 6.
		}

		log.Info("%s calculated %.0f width", config.Desc(), w)
		config.SetWidth = true
		config.Width = uint64(w)

	} else if config.GetSetWidth() && !config.GetSetHeight() {

		width := config.GetWidth()
		h := math32.Round(float32(width) / (cameraRatio / fontRatio))
		if width <= 6 { //special case
			h = 1.
		}

		log.Info("%s calculated %.0f height", config.Desc(), h)
		config.SetHeight = true
		config.Height = uint64(h)
	}
}
