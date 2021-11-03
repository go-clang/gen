package gen

// Lookup represents a in-memory lookup store.
type Lookup struct {
	lookupEnum        map[string]*Enum
	lookupNonTypedefs map[string]string
	lookupStruct      map[string]*Struct
}

// NewLookup returns the intitialised Lookup.
func NewLookup() Lookup {
	return Lookup{
		lookupEnum:        map[string]*Enum{},
		lookupNonTypedefs: map[string]string{},
		lookupStruct: map[string]*Struct{
			"cxstring": {
				Name:  "cxstring",
				CName: "CXString",
			},
		},
	}
}

// RegisterEnum registers e *Enum to Lookup.
func (l *Lookup) RegisterEnum(e *Enum) {
	if _, ok := l.lookupEnum[e.Name]; !ok {
		l.lookupEnum[e.Name] = e
		l.lookupNonTypedefs["enum "+e.CName] = e.Name
		l.lookupEnum[e.CName] = e
	}
}

// HasEnum reports whether the n has Enum.
func (l *Lookup) HasEnum(n string) (*Enum, bool) {
	e, ok := l.lookupEnum[n]

	return e, ok
}

// RegisterStruct registers s *Struct to Lookup.
func (l *Lookup) RegisterStruct(s *Struct) {
	if _, ok := l.lookupStruct[s.Name]; !ok {
		l.lookupStruct[s.Name] = s
		l.lookupNonTypedefs["struct "+s.CName] = s.Name
		l.lookupStruct[s.CName] = s
	}
}

// HasStruct reports whether the n has Struct.
func (l *Lookup) HasStruct(n string) (*Struct, bool) {
	s, ok := l.lookupStruct[n]

	return s, ok
}

// RemoveStruct removes s *Struct from Lookup.
func (l *Lookup) RemoveStruct(s *Struct) {
	delete(l.lookupStruct, s.Name)
	delete(l.lookupNonTypedefs, "struct "+s.CName)
	delete(l.lookupStruct, s.CName)
}

// LookupNonTypedef lookups non typedef from Lookup.
func (l *Lookup) LookupNonTypedef(s string) (string, bool) {
	n, ok := l.lookupNonTypedefs[s]

	return n, ok
}

// IsEnumOrStruct reports whether the name is Enum or Struct.
func (l *Lookup) IsEnumOrStruct(name string) bool {
	if _, ok := l.HasEnum(name); ok {
		return true
	} else if _, ok := l.HasStruct(name); ok {
		return true
	}

	return false
}
