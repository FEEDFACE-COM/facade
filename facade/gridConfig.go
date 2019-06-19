//
//
package facade




//
import (
    "flag"    
    "fmt"
//    gfx "../gfx"
    "strings"
)



var GridDefaults GridConfig = GridConfig {
    Width:        0,
    Height:       8,
    Downward: false,
    Speed:      0.4,
    Buffer:       2,
    Terminal: false,
    Vert:     "def",
    Frag:     "def",
    Fill:        "",
}




func (config *GridConfig) Desc() string { 
    ret := "grid["
    {
    	wok := config.GetCheckWidth(); 
	    hok := config.GetCheckHeight();
	    if wok { ret += fmt.Sprintf("%d",config.GetWidth()) }
    	if wok || hok { ret += "x" }
	    if hok { ret += fmt.Sprintf("%d",config.GetHeight()) }
	    if wok || hok { ret += " " }
	}
    
    {
        tok := config.GetCheckTerminal(); 
        bok := config.GetCheckBuffer();
        if bok { ret += fmt.Sprintf("+%d ",config.GetBuffer()) }
        if tok && config.GetTerminal() { ret += "TT " }
    }
    
    {
        tmp := "↑"
        dok := config.GetCheckDownward(); 
		pok := config.GetCheckSpeed();
        if dok && config.GetDownward() { tmp = "↓" } 
		if pok { ret += fmt.Sprintf("%.1f",config.GetSpeed()) }
        if dok { ret += tmp }
		if dok || pok { ret += " " }
	}
	{
		vok := config.GetCheckVert()
		fok := config.GetCheckFrag()
		if vok { ret += config.GetVert() }
		if vok || fok { ret += "," }
		if fok { ret += config.GetFrag() }
		if vok || fok { ret += " " }	
	}
    if config.GetCheckFill() { ret += config.GetFill() + " " } 
    ret = strings.TrimRight(ret, " ")
    ret += "]"
    return ret
}



func (config *GridConfig) AddFlags(flagset *flag.FlagSet) {
    
    flagset.Uint64Var( &config.Width, "w", GridDefaults.Width, "grid width" ) 
    flagset.Uint64Var( &config.Height,"h",GridDefaults.Height,"grid height")
    flagset.BoolVar(&config.Downward,"down",GridDefaults.Downward,"downward")
    flagset.Float64Var(&config.Speed,"speed",GridDefaults.Speed,"scroll speed")
    flagset.Uint64Var( &config.Buffer,"buffer",GridDefaults.Buffer,"buffer lines")
    flagset.BoolVar(&config.Terminal,"term",GridDefaults.Terminal,"ansi terminal")
    flagset.StringVar(&config.Vert,"vert",GridDefaults.Vert,"vertex shader")
    flagset.StringVar(&config.Frag,"frag",GridDefaults.Frag,"fragment shader")
    flagset.StringVar(&config.Fill,"fill",GridDefaults.Fill,"fill pattern")
    
}


func (config *GridConfig) VisitFlags(flagset *flag.FlagSet)  {
	flagset.Visit( func(flg *flag.Flag) {
        switch flg.Name {
            case "w":        { config.CheckWidth = true;    }
            case "h":        { config.CheckHeight = true;   }
            case "down":     { config.CheckDownward = true; }
            case "speed":    { config.CheckSpeed = true;    }
            case "buffer":   { config.CheckBuffer = true;   }
            case "term": { config.CheckTerminal = true; }
            case "vert":     { config.CheckVert = true;     }
            case "frag":     { config.CheckFrag = true;     }
            case "fill":     { config.CheckFill = true;     }
	   }
    })
}


