package dom

import (
	"fmt"
	"strings"
)

// GetEventID returns the id for a ElemEvent object
func GetEventID(m *ElemEvent) string {
	sel := strings.TrimSpace(m.EventSelector())
	return BuildEventID(sel, m.EventType())
}

// BuildEventID returns the string represent of the values using the select#event format
func BuildEventID(etype, eselect string) string {
	return fmt.Sprintf("%s#%s", eselect, etype)
}
