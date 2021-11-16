package gen

func (g *Generation) API() *API {
	return g.api
}

func (g *Generation) Enums() []*Enum {
	return g.enums
}

func (g *Generation) Functions() []*Function {
	return g.functions
}

func (g *Generation) Structs() []*Struct {
	return g.structs
}
