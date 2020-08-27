//
//
package facade

//
import (
	"flag"
	"fmt"
	"strings"
	gfx "../gfx"
)

var LineDefaults LineConfig = LineConfig{
	Downward: false,
	Speed:    1.0,
	Fixed:    false,
	Drop:     true,
	Smooth:   true,
	Buffer:   8,
}

func (config *LineConfig) Desc() string {
	ret := "line["

	if shader := config.GetShader(); shader != nil {
		ret += shader.Desc() + " "
	}

	if grid := config.GetGrid(); grid != nil {
		ret += " " + grid.Desc() + " "
	}
	{
		bok := config.GetSetBuffer()
		if bok {
			ret += fmt.Sprintf("+%d ", config.GetBuffer())
		}
	}

	{
		down, fixed, drop, smooth := "", "", "", ""
		dok := config.GetSetDownward()
		sok := config.GetSetSpeed()
		aok := config.GetSetFixed()
		pok := config.GetSetDrop()
		mok := config.GetSetSmooth()

		if dok {
			if config.GetDownward() {
				down = "↓"
			}
			if !config.GetDownward() {
				down = "↑"
			}
		}
		if aok {
			if config.GetFixed() {
				fixed = "F"
			}
			if !config.GetFixed() {
				fixed = ""
			}
		}
		if pok {
			if config.GetDrop() {
				drop = ""
			}
			if !config.GetDrop() {
				drop = "ṕ"
			}
		}

		if mok {
			if config.GetSmooth() {
				smooth = "S"
			}
			if !config.GetSmooth() {
				smooth = ""
			}
		}

		if dok {
			ret += down
		}
		if sok {
			ret += fmt.Sprintf("%.1f", config.GetSpeed())
		}
		if aok {
			ret += fixed
		}
		if pok {
			ret += drop
		}
		if mok {
			ret += smooth
		}
		if dok || sok || aok || pok || mok {
			ret += " "
		}
	}

	ret = strings.TrimRight(ret, " ")
	ret += "]"
	return ret
}

func (config *LineConfig) AddFlags(flagset *flag.FlagSet) {

	if config.GetShader() != nil {
		config.GetShader().AddFlags(flagset)
	}
	if config.GetGrid() != nil {
		config.GetGrid().AddFlags(flagset)
	}

	flagset.BoolVar(&config.Downward, "down", LineDefaults.Downward, "line scroll downward?")
	flagset.BoolVar(&config.Drop, "drop", LineDefaults.Drop, "line drop lines?")
	flagset.BoolVar(&config.Smooth, "smooth", LineDefaults.Smooth, "line scroll smooth?")
	flagset.Float64Var(&config.Speed, "speed", LineDefaults.Speed, "line scroll speed")
	flagset.BoolVar(&config.Fixed, "fixed", LineDefaults.Fixed, "line fixed scroll speed?")
	flagset.Uint64Var(&config.Buffer, "buffer", LineDefaults.Buffer, "line buffer length")

}

func (config *LineConfig) VisitFlags(flagset *flag.FlagSet) bool {
	ret := false
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
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

	if shader := config.GetShader(); shader != nil {
		if shader.VisitFlags(flagset) {
			ret = true
		}
	}

	if grid := config.GetGrid(); grid != nil {
		if grid.VisitFlags(flagset) {
			ret = true
		}
	}
	return ret

}

func (config *LineConfig) Help() string {
	ret := GridDefaults.Help()
	tmp := flag.NewFlagSet("line", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll( func (f *flag.Flag) { ret += gfx.FlagHelp(f) } )
	return ret
}    
    
    
