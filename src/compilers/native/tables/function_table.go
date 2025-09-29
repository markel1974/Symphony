package tables

import (
	"fmt"
	"go/ast"

	"github.com/markel1974/c64emu/src/vm/objects"
)

const (
	UndefinedSymbol = "_"
)

type FunctionInput struct {
	Name string
	Type string
}

// FunctionDescription represents the structure to describe a Go function, its metadata, and related declarations.
type FunctionDescription struct {
	Name            string
	ReturnTypes     []string
	InputNames      []string
	InputTypes      []string
	StructName      string
	StructReceivers []string
	FuncDecl        *ast.FuncDecl
}

// NewFunctionDescription creates a new FunctionDescription using the provided *ast.FuncDecl.
func NewFunctionDescription(funcDecl *ast.FuncDecl) *FunctionDescription {
	return &FunctionDescription{
		FuncDecl: funcDecl,
	}
}

// FunctionTable is a type designed to manage a collection of function descriptions.
type FunctionTable struct {
	gk              objects.IGateKeeper
	scopes          *Scopes
	definitionTable *DefinitionTable
	container       []*FunctionDescription
}

// NewFunctionTable initializes and returns a new instance of FunctionTable.
func NewFunctionTable(gk objects.IGateKeeper, scopes *Scopes, definitionTable *DefinitionTable) *FunctionTable {
	return &FunctionTable{
		gk:              gk,
		scopes:          scopes,
		definitionTable: definitionTable,
	}
}

// Add appends a new function description, created from the provided function declaration, to the container.
func (f *FunctionTable) Add(node *ast.FuncDecl) {
	fd := NewFunctionDescription(node)
	f.container = append(f.container, fd)
}

// Len returns the number of elements currently stored in the container of the FunctionTable.
func (f *FunctionTable) Len() int {
	return len(f.container)
}

// Get retrieves a FunctionDescription at the specified index. Returns an error if the index is out of bounds.
func (f *FunctionTable) Get(index int) (*FunctionDescription, error) {
	if index < 0 || index >= f.Len() {
		return nil, fmt.Errorf("function handler: index %d out of bounds", index)
	}
	return f.container[index], nil
}

// GetByName retrieves a FunctionDescription from the container by its name. Returns nil if no match is found.
func (f *FunctionTable) GetByName(name string) *FunctionDescription {
	for _, v := range f.container {
		if name == v.Name {
			return v
		}
	}
	return nil
}

// CountParams correctly counts the number of parameters in a field list,
// handling both named (e.g., a, b int) and unnamed (e.g., int) parameters.
func (f *FunctionTable) CountParams(fieldList *ast.FieldList) int {
	if fieldList == nil || fieldList.List == nil {
		return 0
	}
	count := 0
	for _, field := range fieldList.List {
		if len(field.Names) == 0 {
			// Unnamed parameter, like func(int)
			count++
		} else {
			// Named parameters, like func(a, b int)
			count += len(field.Names)
		}
	}
	return count
}

// SymbolsFromFields creates and returns a slice of Symbols for the function parameters.
func (f *FunctionTable) SymbolsFromFields(fieldList *ast.FieldList) ([]*Symbol, error) {
	if fieldList == nil {
		return nil, nil
	}
	var symbols []*Symbol
	for _, p := range fieldList.List {
		var typeName string
		if ident := GetIdent(p.Type); ident != nil {
			typeName = ident.Name
		}
		for _, name := range p.Names {
			symbol, err := f.scopes.SymbolDefine(name.Name)
			if err != nil {
				return nil, err
			}
			if err = f.definitionTable.SymbolAssign(symbol, typeName); err != nil {
				return nil, err
			}
			symbols = append(symbols, symbol)
		}
	}
	return symbols, nil
}

// RangeKey returns the symbol for the range key, if any.
func (f *FunctionTable) RangeKey(node *ast.RangeStmt) (*Symbol, error) {
	if node.Key == nil {
		return nil, nil
	}
	switch k := node.Key.(type) {
	case *ast.Ident:
		if k.Name == UndefinedSymbol {
			return nil, nil
		}
		keySymbol, err := f.scopes.SymbolDefine(k.Name)
		if err != nil {
			return nil, err
		}
		return keySymbol, nil
	default:
		return nil, nil
	}
}

// RangeValue resolves and defines a symbol for the `Value` in a range statement, assigning it a type if specified.
func (f *FunctionTable) RangeValue(node *ast.RangeStmt, typeName string) (*Symbol, error) {
	if node.Value == nil {
		return nil, nil
	}
	switch k := node.Value.(type) {
	case *ast.Ident:
		if k.Name == UndefinedSymbol {
			return nil, nil
		}
		symbol, err := f.scopes.SymbolDefine(k.Name)
		if err != nil {
			return nil, err
		}
		if err = f.definitionTable.SymbolAssign(symbol, typeName); err != nil {
			return nil, err
		}
		return symbol, nil
	default:
		return nil, nil
	}
}
