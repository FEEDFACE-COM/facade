//
//
package facade

//
import (
	gfx "../gfx"
	log "../log"
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
		ret += config.GetFill() + " "
	}
	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *GridConfig) AddFlags(flagset *flag.FlagSet) {
	patterns := " (" + availableFillPatterns() + ")"
	flagset.Uint64Var(&config.Width, "w", GridDefaults.Width, "grid width")
	flagset.Uint64Var(&config.Height, "h", GridDefaults.Height, "grid height")
	flagset.StringVar(&config.Fill, "fill", GridDefaults.Fill, "grid fill pattern"+patterns)

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

func (config *GridConfig) autoWidth(cameraRatio float32, fontRatio float32) {

	if !config.GetSetWidth() {
		if !config.GetSetHeight() {
			return
		}

		height := config.GetHeight()
		w := cameraRatio / fontRatio * float32(height)
		if height == 1 {
			w = 5.
		} //special case

		config.SetWidth = true
		config.Width = uint64(w)

	}
	log.Debug("%s calculated width", config.Desc())
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

func availableFillPatterns() string {
	ret := ""
	for _, name := range FillPatterns {
		ret += name
		ret += ", "
	}
	ret = strings.TrimSuffix(ret, ", ")
	return ret
}
