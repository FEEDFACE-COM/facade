//
//
package facade




//
import (
    "flag"    
    "fmt"
    "strings"
)



var GridDefaults GridConfig = GridConfig{
    Width:       32,
    Height:       8,
    Downward: false,
    Speed:      1.0,
    Adaptive:  true,
    Drop:      true,
    Smooth:    true,
    Buffer:       8,
    Terminal: false,
    Vert:     "def",
    Frag:     "def",
    Fill:        "",
}




func (config *GridConfig) Desc() string { 
    ret := "grid["
    {
    	wok := config.GetSetWidth(); 
	    hok := config.GetSetHeight();
	    if wok { ret += fmt.Sprintf("%d",config.GetWidth()) }
    	if wok || hok { ret += "x" }
	    if hok { ret += fmt.Sprintf("%d",config.GetHeight()) }
	    if wok || hok { ret += " " }
	}
    
    {
        tok := config.GetSetTerminal(); 
        bok := config.GetSetBuffer();
        if bok { ret += fmt.Sprintf("+%d ",config.GetBuffer()) }
        if tok && config.GetTerminal() { ret += "TT " }
    }
    
    {
        down,adapt,drop,smooth := "","","",""
        dok := config.GetSetDownward(); 
		sok := config.GetSetSpeed();
		aok := config.GetSetAdaptive();
		pok := config.GetSetDrop()
		mok := config.GetSetSmooth()
		
        if dok {
            if config.GetDownward() { down = "↓" } 
            if ! config.GetDownward() { down = "↑" }
        }
        if aok { 
            if  config.GetAdaptive() {adapt = "a" }
            if ! config.GetAdaptive() { adapt = "á" }
        }
        if pok {
            if config.GetDrop() { drop = "p" }
            if ! config.GetDrop() { drop = "ṕ" }
        }

        if mok { 
            if config.GetSmooth() { smooth = "s" }
            if ! config.GetSmooth() { smooth = "ś" }
        }
        
        if dok { ret += down }
		if sok { ret += fmt.Sprintf("%.1f",config.GetSpeed()) }
        if aok { ret += adapt }
        if pok { ret += drop }
        if mok { ret += smooth }
		if dok || sok || aok || pok || mok { ret += " " }
	}

	{
		vok := config.GetSetVert()
		fok := config.GetSetFrag()
		if vok { ret += config.GetVert() }
		if vok || fok { ret += "," }
		if fok { ret += config.GetFrag() }
		if vok || fok { ret += " " }	
	}

    if config.GetSetFill() { ret += config.GetFill() + " " } 
    ret = strings.TrimRight(ret, " ")
    ret += "]"
    return ret
}



func (config *GridConfig) AddFlags(flagset *flag.FlagSet) {
    
    flagset.Uint64Var( &config.Width, "w", GridDefaults.Width, "grid width" ) 
    flagset.Uint64Var( &config.Height,"h",GridDefaults.Height,"grid height")
    flagset.BoolVar(&config.Downward,"down",GridDefaults.Downward,"scroll downward?")
    flagset.BoolVar(&config.Drop,"drop",GridDefaults.Drop,"drop lines?")
    flagset.BoolVar(&config.Smooth,"smooth",GridDefaults.Smooth,"smooth speed?")
    flagset.Float64Var(&config.Speed,"speed",GridDefaults.Speed,"scroll speed")
    flagset.BoolVar(&config.Adaptive,"adapt",GridDefaults.Adaptive,"adapt speed?")
    flagset.Uint64Var( &config.Buffer,"buffer",GridDefaults.Buffer,"buffer lines")
    flagset.BoolVar(&config.Terminal,"term",GridDefaults.Terminal,"ansi terminal?")
    flagset.StringVar(&config.Vert,"vert",GridDefaults.Vert,"vertex shader")
    flagset.StringVar(&config.Frag,"frag",GridDefaults.Frag,"fragment shader")
    flagset.StringVar(&config.Fill,"fill",GridDefaults.Fill,"fill pattern")
    
}


func (config *GridConfig) VisitFlags(flagset *flag.FlagSet) bool {
	flagset.Visit( func(flg *flag.Flag) {
        switch flg.Name {
            case "w":        { config.SetWidth = true;    }
            case "h":        { config.SetHeight = true;   }
            case "down":     { config.SetDownward = true; }
            case "drop":     { config.SetDrop = true;    }
            case "smooth":   { config.SetSmooth = true;   }
            case "speed":    { config.SetSpeed = true;    }
            case "adapt":    { config.SetAdaptive = true; }
            case "buffer":   { config.SetBuffer = true;   }
            case "term":     { config.SetTerminal = true; }
            case "vert":     { config.SetVert = true;     }
            case "frag":     { config.SetFrag = true;     }
            case "fill":     { config.SetFill = true;     }
	   }
    })
    return config.SetWidth    ||
           config.SetHeight   ||
           config.SetDownward ||
           config.SetSpeed    ||
           config.SetAdaptive ||
           config.SetDrop     ||
           config.SetSmooth   ||
           config.SetBuffer   ||
           config.SetTerminal ||
           config.SetVert     ||
           config.SetFrag     ||
           config.SetFill
}


func (config *GridConfig) autoWidth(cameraRatio float32, fontRatio float32) {
    
    if ! config.GetSetWidth() {
    
        if ! config.GetSetHeight() {
            return
        }
        
        height := config.GetHeight()
        w := cameraRatio / fontRatio * float32(height)
        if height == 1 { w = 5. } //special case

        config.SetWidth = true
        config.Width = uint64(w)
        
    }
}

