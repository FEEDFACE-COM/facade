package facade

import (
	"flag"
	"fmt"
)

type DraftConfig struct {
}

func NewDraftConfig() *DraftConfig {
	return &DraftConfig{}
}

func (config *DraftConfig) AddFlags(flags *flag.FlagSet) {
}

func (config *DraftConfig) Desc() string { return fmt.Sprintf("draft[]") }
