//
//
package facade

//
import (
	"flag"
	"fmt"
	"strings"
)

var LineDefaults LineConfig = LineConfig{
	Downward: false,
	Speed:    1.0,
	Adaptive: true,
	Drop:     true,
	Smooth:   true,
	Buffer:   8,
	Grid:     nil,
}

func (config *LineConfig) Desc() string {
	ret := "line["

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
		down, adapt, drop, smooth := "", "", "", ""
		dok := config.GetSetDownward()
		sok := config.GetSetSpeed()
		aok := config.GetSetAdaptive()
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
			if config.GetAdaptive() {
				adapt = "a"
			}
			if !config.GetAdaptive() {
				adapt = "á"
			}
		}
		if pok {
			if config.GetDrop() {
				drop = "p"
			}
			if !config.GetDrop() {
				drop = "ṕ"
			}
		}

		if mok {
			if config.GetSmooth() {
				smooth = "s"
			}
			if !config.GetSmooth() {
				smooth = "ś"
			}
		}

		if dok {
			ret += down
		}
		if sok {
			ret += fmt.Sprintf("%.1f", config.GetSpeed())
		}
		if aok {
			ret += adapt
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

	if config.GetGrid() != nil {
		config.GetGrid().AddFlags(flagset)
	}

	flagset.BoolVar(&config.Downward, "down", LineDefaults.Downward, "line scroll downward?")
	flagset.BoolVar(&config.Drop, "drop", LineDefaults.Drop, "line drop lines?")
	flagset.BoolVar(&config.Smooth, "smooth", LineDefaults.Smooth, "line continuous scroll?")
	flagset.Float64Var(&config.Speed, "speed", LineDefaults.Speed, "line scroll speed")
	flagset.BoolVar(&config.Adaptive, "adapt", LineDefaults.Adaptive, "line adaptive scroll?")
	flagset.Uint64Var(&config.Buffer, "buffer", LineDefaults.Buffer, "line buffer length")

}

func (config *LineConfig) VisitFlags(flagset *flag.FlagSet) bool {
	flagset.Visit(func(flg *flag.Flag) {
		switch flg.Name {
		case "down":
			{
				config.SetDownward = true
			}
		case "drop":
			{
				config.SetDrop = true
			}
		case "smooth":
			{
				config.SetSmooth = true
			}
		case "speed":
			{
				config.SetSpeed = true
			}
		case "adapt":
			{
				config.SetAdaptive = true
			}
		case "buffer":
			{
				config.SetBuffer = true
			}
		}
	})
	setGrid := false
	if grid := config.GetGrid(); grid != nil {
		setGrid = grid.VisitFlags(flagset)
	}
	return setGrid ||
		config.SetDownward ||
		config.SetSpeed ||
		config.SetAdaptive ||
		config.SetDrop ||
		config.SetSmooth ||
		config.SetBuffer

}

func (config *LineConfig) Help() string {
	ret := ""
	fun := func(f *flag.Flag) {
		name := f.Name
		if f.DefValue != "false" && f.DefValue != "true" {
			name = f.Name + "=" + f.DefValue
		} else if f.DefValue == "true" {
    		name = f.Name + "=f"
        }
		ret += fmt.Sprintf("  -%-24s %-24s\n", name, f.Usage)
	}

	ret += GridDefaults.Help()
	tmp := flag.NewFlagSet("line", flag.ExitOnError)
	config.AddFlags(tmp)
	tmp.VisitAll(fun)
	return ret
}
