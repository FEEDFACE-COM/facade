package gfx

import (
	"flag"
	"fmt"
)

func FlagHelp(f *flag.Flag) string {
	name := f.Name
	if f.DefValue == "false" {
		name = f.Name + ""
	} else if f.DefValue == "true" {
		name = f.Name + "=f"
	} else {
		name = f.Name + "=" + f.DefValue
	}
	return fmt.Sprintf("  -%-24s %-24s\n", name, f.Usage)
}
