package views

// Views defines the view interface with its render method rules which all views must confirm to
type Views interface {
	Render() string
	//String must calls .Render() underneath
	String() string
}
