package dom

import "fmt"

// GetEventID returns the id for a ElemEvent object
func GetEventID(m *ElemEvent) string {
	return BuildEventID(m.EventSelector(), m.EventType())
}

// BuildEventID returns the string represent of the values using the select#event format
func BuildEventID(etype, eselect string) string {
	return fmt.Sprintf("%s#%s", eselect, etype)
}
