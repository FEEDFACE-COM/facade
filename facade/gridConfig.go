package facade

import (
	"FEEDFACE.COM/facade/gfx"
	"FEEDFACE.COM/facade/log"
	"FEEDFACE.COM/facade/math32"
	"flag"
	"fmt"
	"strings"
)

var GridDefaults GridConfig = GridConfig{
	Width:  32,
	Height: 8,
	Fill:   "",
}

func (config *GridConfig) Desc() string {
	ret := "grid["
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

	if config.GetSetFill() {
		ret += "f" + config.GetFill() + " "
	}
	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *GridConfig) AddFlags(flagset *flag.FlagSet) {
	patterns := "title,grid,alpha,clear"
	flagset.Uint64Var(&config.Width, "w", GridDefaults.Width, "grid width")
	flagset.Uint64Var(&config.Height, "h", GridDefaults.Height, "grid height")
	flagset.StringVar(&config.Fill, "fill", GridDefaults.Fill, "fill pattern ("+patterns+")")

}

func (config *GridConfig) VisitFlags(flagset *flag.FlagSet) bool {
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "w":
			{
				config.SetWidth = true
			}
		case "h":
			{
				config.SetHeight = true
			}
		case "fill":
			{
				config.SetFill = true
			}
		}
	})
	return config.SetWidth ||
		config.SetHeight ||
		config.SetFill
}

func (config *GridConfig) autoSize(cameraRatio float32, fontRatio float32) {

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

func (config *GridConfig) Help() string {
	ret := ""
	tmp := flag.NewFlagSet("grid", flag.ExitOnError)
	config.AddFlags(tmp)
	for _, s := range []string{"w", "h", "fill"} {
		if flg := tmp.Lookup(s); flg != nil {
			ret += gfx.FlagHelp(flg)
		}
	}
	return ret
}
