package trees

// This contains printers for the tree dom definition structures

// AttrPrinter defines a printer interface for writing out a Attribute objects into a string form
type AttrPrinter interface {
	Print([]*Attribute) string
}

// StylePrinter defines a printer interface for writing out a style objects into a string form
type StylePrinter interface {
	Print([]*Style) string
}

// MarkupPrinter defines a printer interface for writing out a markup object into a string form
type MarkupPrinter interface {
	Read([]Markup)
}
