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
	gk             objects.IGateKeeper
	scopes         *Scopes
	structTable    *StructTable
	interfaceTable *InterfaceTable
	container      []*FunctionDescription
}

// NewFunctionTable initializes and returns a new instance of FunctionTable.
func NewFunctionTable(gk objects.IGateKeeper, scopes *Scopes, structTable *StructTable, interfaceTable *InterfaceTable) *FunctionTable {
	return &FunctionTable{
		gk:             gk,
		scopes:         scopes,
		structTable:    structTable,
		interfaceTable: interfaceTable,
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

// SymbolsFromParameters creates and returns a slice of Symbols for the function parameters.
func (f *FunctionTable) SymbolsFromParameters(fieldList *ast.FieldList) ([]*Symbol, error) {
	var symbols []*Symbol
	if fieldList != nil {
		for _, p := range fieldList.List {
			typeName := f.TypeName(p.Type)
			isStruct := f.structTable.Has(typeName)
			isInterface := f.interfaceTable.Has(typeName)
			for _, name := range p.Names {
				symbol, err := f.scopes.SymbolDefine(name.Name)
				if err != nil {
					return nil, err
				}
				symbol.SetScope(LocalScope)
				symbol.SetReturnTypes([]string{typeName})
				if isStruct {
					f.structTable.BindSymbol(symbol, typeName)
					symbol.SetObject(f.gk.NewString(objects.FrameStatic, typeName+":"+symbol.Name()))
				} else if isInterface {
					symbol.SetInterface(typeName)
					symbol.SetObject(f.gk.NewString(objects.FrameStatic, "interface:"+symbol.Name()))
				}
				symbols = append(symbols, symbol)
			}
		}
	}
	return symbols, nil
}

// RangeKey returns the symbol for the range key, if any.
func (f *FunctionTable) RangeKey(node *ast.RangeStmt) (*Symbol, error) {
	if node.Key != nil {
		if ident, ok := node.Key.(*ast.Ident); ok && ident.Name != UndefinedSymbol {
			keySymbol, err := f.scopes.SymbolDefine(ident.Name)
			if err != nil {
				return nil, err
			}
			return keySymbol, nil
		}
	}
	return nil, nil
}

// RangeValue resolves and defines a symbol for the `Value` in a range statement, assigning it a type if specified.
func (f *FunctionTable) RangeValue(node *ast.RangeStmt, typeName string) (*Symbol, error) {
	if node.Value != nil {
		if ident, ok := node.Value.(*ast.Ident); ok && ident.Name != UndefinedSymbol {
			valueSymbol, err := f.scopes.SymbolDefine(ident.Name)
			if err != nil {
				return nil, err
			}
			if len(typeName) > 0 {
				valueSymbol.SetReturnTypes([]string{typeName})
				valueSymbol.SetObject(f.gk.NewString(objects.FrameStatic, typeName+":"+valueSymbol.Name()))
				f.structTable.BindSymbol(valueSymbol, typeName)
			}
			return valueSymbol, nil
		}
	}
	return nil, nil
}

func (f *FunctionTable) TypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return f.TypeName(t.X)
	}
	return ""
}
