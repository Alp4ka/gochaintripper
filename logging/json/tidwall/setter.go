package tidwall

import (
	"github.com/Alp4ka/gochaintripper/logging"
	"github.com/tidwall/sjson"
)

type Setter struct{}

func NewSetter() *Setter {
	return &Setter{}
}

// SetValue implements logging.JSONSetter interface.
func (s *Setter) SetValue(json string, fieldPattern, replaceString string) (string, error) {
	return sjson.Set(json, fieldPattern, replaceString)
}

var _ logging.JSONSetter = (*Setter)(nil)
