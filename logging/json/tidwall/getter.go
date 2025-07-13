package tidwall

import (
	"github.com/Alp4ka/gochaintripper/logging"
	"github.com/tidwall/gjson"
)

type Getter struct{}

func NewGetter() *Getter {
	return &Getter{}
}

// IsJSON implements logging.JSONGetter interface.
func (g *Getter) IsJSON(json string) bool {
	return gjson.Valid(json)
}

// Exists implements logging.JSONGetter interface.
func (g *Getter) Exists(json string, fieldPattern string) bool {
	return gjson.Get(json, fieldPattern).Exists()
}

var _ logging.JSONGetter = (*Getter)(nil)
