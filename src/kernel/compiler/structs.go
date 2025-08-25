package compiler

import (
	"fmt"
	"go/ast"
)

type StructProperty struct {
	name string
	base string
	kind string
	node ast.Node
}

func NewStructProperty(name string, base string, kind string, node ast.Node) *StructProperty {
	return &StructProperty{
		name: name,
		base: base,
		kind: kind,
		node: node,
	}
}

// Structs gestisce una collezione di tutte le definizioni di struct trovate
// durante la compilazione, permettendo un accesso rapido tramite nome.
type Structs struct {
	container map[string][]*StructProperty
}

// NewStructs crea un nuovo gestore di struct.
func NewStructs() *Structs {
	return &Structs{
		container: make(map[string][]*StructProperty),
	}
}

// Add aggiunge una nuova definizione di struct al contenitore.
func (s *Structs) Add(name string, def []*StructProperty) {
	s.container[name] = def
	fmt.Printf("Struct %s defined with fields:\n", name)
	for _, p := range def {
		fmt.Printf("\t%s %s %s\n", p.name, p.base, p.kind)
	}
	fmt.Printf("\n")
}

func (s *Structs) Get(name string) []*StructProperty {
	return s.container[name]
}

/*
// Get recupera la definizione di uno struct tramite il suo nome.
func (s *Structs) Get(name string) (*StructDef, bool) {
	def, ok := s.container[name]
	return def, ok
}

*/
