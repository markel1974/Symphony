package tables

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

const (
	UndefinedSymbol = "_"
)

// FunctionDescription represents the structure to describe a Go function, its metadata, and related declarations.
type FunctionDescription struct {
	Name            string
	ReturnTypes     []string
	InputParams     []string
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

// DefineFunctionVariables handles the definition and type inference of variables from function call assignments.
// It ensures correct assignment based on token type (:= or =) and aligns the number of variables to return types.
// Emits appropriate bytecode for defining or setting variables and supports scoped symbol management.
func (f *FunctionTable) DefineFunctionVariables(tok token.Token, callExpr *ast.CallExpr, lhs []ast.Expr) error {
	var funcReturnTypes []string
	if ident, isIdent := callExpr.Fun.(*ast.Ident); isIdent {
		if ident.Name != "" {
			if funcSymbol, ok := f.scopes.SymbolResolve(ident.Name); ok {
				funcReturnTypes = funcSymbol.ReturnTypes()
			}
		}
	}
	if len(lhs) != len(funcReturnTypes) {
		return fmt.Errorf("assignment mismatch: %d variables but %d return values", len(lhs), len(funcReturnTypes))
	}
	if len(lhs) == 0 {
		return nil
	}
	// Check type of each variable e.g.: a, b := Test()
	for i := len(lhs) - 1; i >= 0; i-- {
		ident, ok := lhs[i].(*ast.Ident)
		if !ok {
			return fmt.Errorf("unsupported multiple assignment to type %T", lhs)
		}
		var symbol *Symbol
		if tok == token.DEFINE {
			var err error
			if symbol, err = f.scopes.SymbolDefine(ident.Name); err != nil {
				return err
			}
		} else {
			var found bool
			if symbol, found = f.scopes.SymbolResolve(ident.Name); !found {
				return fmt.Errorf("undefined variable: %s", ident.Name)
			}
		}
		// Complete type inference for each variable
		inferredTypeName := funcReturnTypes[i]
		if err := f.structTable.AssignSymbol(symbol, inferredTypeName, []string{inferredTypeName}); err != nil {
			return err
		}
		// Emit correct opcode based on ':=' or '='
		if tok == token.DEFINE {
			if err := f.scopes.EmitSymbolDefine(symbol); err != nil {
				return err
			}
		} else {
			if err := f.scopes.EmitSymbolSet(symbol); err != nil {
				return err
			}
		}
		if _, err := f.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
	}
	return nil
}

// SymbolsFromParameters creates and returns a slice of Symbols for the function parameters.
func (f *FunctionTable) SymbolsFromParameters(fieldList *ast.FieldList) ([]*Symbol, error) {
	var symbols []*Symbol
	if fieldList != nil {
		for _, p := range fieldList.List {
			var typeName string
			switch t := p.Type.(type) {
			case *ast.Ident:
				typeName = t.Name
			case *ast.StarExpr:
				if ident, ok := t.X.(*ast.Ident); ok {
					typeName = ident.Name
				}
			}
			// Determina se il tipo è uno struct o un'interfaccia conosciuta
			isStruct := f.structTable.Has(typeName)
			isInterface := f.interfaceTable.Has(typeName)
			for _, name := range p.Names {
				symbol, err := f.scopes.SymbolDefine(name.Name)
				if err != nil {
					return nil, err
				}
				symbol.SetScope(LocalScope)
				if isStruct {
					// Se è uno struct, assegna le informazioni dello struct
					if err = f.structTable.AssignSymbol(symbol, typeName, []string{typeName}); err != nil {
						return nil, err
					}
				} else if isInterface {
					// Se è un'interfaccia, contrassegnalo come tale!
					symbol.SetInterface(typeName)
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
func (f *FunctionTable) RangeValue(node *ast.RangeStmt, returnTypeName string) (*Symbol, error) {
	if node.Value != nil {
		if ident, ok := node.Value.(*ast.Ident); ok && ident.Name != UndefinedSymbol {
			valueSymbol, err := f.scopes.SymbolDefine(ident.Name)
			if err != nil {
				return nil, err
			}
			if len(returnTypeName) > 0 {
				if err = f.structTable.AssignSymbol(valueSymbol, returnTypeName, []string{returnTypeName}); err != nil {
					return nil, err
				}
			}
			return valueSymbol, nil
		}
	}
	return nil, nil
}
