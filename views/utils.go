package views

import (
	"fmt"

	"github.com/influx6/flux"
)

// MakeBlueprintName generates a new name forthe blueprint
func MakeBlueprintName(b Blueprint) string {
	return fmt.Sprintf("%s:%s", b.Type(), flux.RandString(5))
}
