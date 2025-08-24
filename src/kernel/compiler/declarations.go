package compiler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

type Declarations struct {
	gk         *objects.GateKeeper
	references *Constants
	constants  *Constants
	scopes     *Scopes
	fileSet    *token.FileSet
}

func NewDeclarations(gk *objects.GateKeeper, references *Constants, constants *Constants, scopes *Scopes) *Declarations {
	return &Declarations{
		gk: gk, references: references, constants: constants, scopes: scopes,
	}
}

func (c *Declarations) Initialize(fileSet *token.FileSet) {
	c.fileSet = fileSet
}

// compile traverses the provided AST node and compiles it into bytecode, handling various node types in a switch block.
func (c *Declarations) compile(in ast.Node) error {
	var err error = nil
	switch node := in.(type) {
	case *ast.GenDecl:
		err = c.GenDecl(node) // for `var` and `const` which are handled by AssignStmt
	case *ast.DeclStmt:
		err = c.DeclStmt(node)
	case *ast.TypeSpec:
		err = c.TypeSpec(node)
	case *ast.ValueSpec: // handles 'var x = 10'
		err = c.ValueSpec(node)
	case *ast.CompositeLit:
		err = c.CompositeLit(node)
	case *ast.BasicLit:
		err = c.BasicLit(node)
	case *ast.Ident:
		err = c.Ident(node)
	case *ast.AssignStmt:
		err = c.AssignStmt(node)
	default:
		err = fmt.Errorf("unsupported expression type: %T", node)
	}
	return err
}

// DeclStmt processes a declaration statement node by compiling its declaration content. Returns an error if compilation fails.
func (c *Declarations) DeclStmt(node *ast.DeclStmt) error {
	if err := c.compile(node.Decl); err != nil {
		return err
	}
	return nil
}

// GenDecl processes a general declaration node by compiling each specification within the node. It returns an error if any occur.
func (c *Declarations) GenDecl(node *ast.GenDecl) error {
	for _, spec := range node.Specs {
		if err := c.compile(spec); err != nil {
			return err
		}
	}
	return nil
}

// TypeSpec processes a type specification node, validating and defining struct types in the current scope.
func (c *Declarations) TypeSpec(node *ast.TypeSpec) error {
	structType, isStruct := node.Type.(*ast.StructType)
	if !isStruct {
		return nil
	}
	structName := node.Name.Name
	if _, ok := c.scopes.SymbolResolve(structName); ok {
		return fmt.Errorf("type '%s' already defined", structName)
	}
	var fields []*FieldDef
	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			var typeNameBuf bytes.Buffer
			if err := printer.Fprint(&typeNameBuf, c.fileSet, field.Type); err != nil {
				return fmt.Errorf("failed to resolve type for field in struct '%s'", structName)
			}
			fieldType := typeNameBuf.String()
			for _, name := range field.Names {
				// here we could add a check for duplicate fields.
				fields = append(fields, NewFieldDef(name.Name, fieldType, nil))
			}
		}
	}
	symbol, err := c.scopes.SymbolDefine(structName, UnknownScope, true)
	if err != nil {
		return err
	}

	var structData []string
	for _, field := range fields {
		structData = append(structData, "["+field.Name()+" "+field.Type()+"]")
	}
	symbol.SetObject(c.gk.NewString(objects.FrameStatic, structName+":"+strings.Join(structData, ",")))
	symbol.Fields = fields
	return nil
}

// ValueSpec processes a ValueSpec node to handle variable declarations and assignments within a given scope.
func (c *Declarations) ValueSpec(node *ast.ValueSpec) error {
	// handles 'var x = 10'
	for i, name := range node.Names {
		if i > len(node.Values)-1 {
			return fmt.Errorf("too few values for %s", name.Name)
		}
		if err := c.compile(node.Values[i]); err != nil {
			return err
		}

		isStruct := false

		// 3. Inferenza del tipo, ora coerente con la nuova logica
		var assignedTypeNames []string
		if compLit, ok := node.Values[i].(*ast.CompositeLit); ok {
			if ident, ok := compLit.Type.(*ast.Ident); ok {
				if typeSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok && typeSymbol.IsStruct() {
					assignedTypeNames = []string{typeSymbol.Name()}
					isStruct = true
				}
			}
		} else if callExpr, ok := node.Values[i].(*ast.CallExpr); ok {
			var funcName string
			if ident, isIdent := callExpr.Fun.(*ast.Ident); isIdent {
				funcName = ident.Name
			}
			if funcName != "" {
				if funcSymbol, ok := c.scopes.SymbolResolve(funcName); ok {
					returnTypes := funcSymbol.Types()
					if len(returnTypes) != 1 {
						return fmt.Errorf("assignment mismatch: 'var' declaration expects 1 value, but function %s returns %d", funcName, len(returnTypes))
					}
					assignedTypeNames = []string{returnTypes[0]}
				}
			}
		}

		symbol2, err := c.scopes.SymbolDefine(name.Name, UnknownScope, isStruct)
		if err != nil {
			return err
		}

		if len(assignedTypeNames) > 0 {
			symbol2.SetTypes(assignedTypeNames)
			symbol2.SetObject(c.gk.NewString(objects.FrameStatic, symbol2.Name()+":"+strings.Join(assignedTypeNames, " ")))
		} else {
			symbol2.SetObject(c.gk.NewString(objects.FrameStatic, symbol2.Name()))
		}

		// 4. Emette bytecode per assegnare il valore dalla cima dello stack alla variabile.
		if err = c.scopes.EmitSymbolDefine(symbol2); err != nil {
			return err
		}
		// 5. Pulisce lo stack dal valore ora che è stato assegnato.
		if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
	}
	return nil
}

// CompositeLit processes the given composite literal node and compiles it into bytecode representation.
// Handles struct, array, and map literals by resolving types, validating fields, and emitting appropriate instructions.
// Returns an error if the composite literal type is unsupported or if any validation or compilation step fails.
func (c *Declarations) CompositeLit(node *ast.CompositeLit) error {
	switch t := node.Type.(type) {
	case *ast.Ident:
		// struct literal (es. MyStruct{...})
		symbol, ok := c.scopes.SymbolResolve(t.Name)
		if !ok {
			return fmt.Errorf("undefined type: %s", t.Name)
		}
		if !symbol.IsStruct() {
			return fmt.Errorf("unknown composite literal type: %s", t.Name)
		}
		if len(node.Elts) > len(symbol.Fields) {
			return fmt.Errorf("too many values in positional struct literal for type '%s'", symbol.Name())
		}
		for idx := range symbol.Fields {
			symbol.Fields[idx].SetNode(nil)
		}
		isKeyed := false
		if len(node.Elts) > 0 {
			if _, ok := node.Elts[0].(*ast.KeyValueExpr); ok {
				isKeyed = true
			}
		}
		if isKeyed {
			// key literal (es. Home{Name: "Alfa", Address: "Shanghai"})
			providedFields := make(map[string]ast.Expr)
			for _, elt := range node.Elts {
				kvExpr, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					return fmt.Errorf("cannot mix keyed and unkeyed values in struct literal")
				}
				keyIdent, ok := kvExpr.Key.(*ast.Ident)
				if !ok {
					return fmt.Errorf("invalid field name in struct literal")
				}
				providedFields[keyIdent.Name] = kvExpr.Value
			}
			for idx := range symbol.Fields {
				if valueExpr, ok := providedFields[symbol.Fields[idx].Name()]; ok {
					symbol.Fields[idx].SetNode(valueExpr)
				}
			}
		} else {
			// positional literal (es. Home{"Alfa", 20, "Shanghai"}) ---
			for i, elt := range node.Elts {
				symbol.Fields[i].SetNode(elt)
			}
		}
		for idx := range symbol.Fields {
			fieldName := symbol.Fields[idx].Name()
			fieldNode := symbol.Fields[idx].Node()
			keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
			if _, err := c.scopes.Emit(bytecode.OpConstant, keyConst); err != nil {
				return err
			}
			if fieldNode != nil {
				if err := c.compile(fieldNode); err != nil {
					return err
				}
			} else {
				if _, err := c.scopes.Emit(bytecode.OpNull); err != nil {
					return err
				}
			}
		}
		structLen := len(symbol.Fields) * 2
		if _, err := c.scopes.Emit(bytecode.OpStruct, structLen); err != nil {
			return err
		}
		return nil
	case *ast.ArrayType:
		for _, elt := range node.Elts {
			if err := c.compile(elt); err != nil {
				return err
			}
		}
		if _, err := c.scopes.Emit(bytecode.OpArray, len(node.Elts)); err != nil {
			return err
		}
		return nil
	case *ast.MapType:
		for _, elt := range node.Elts {
			kve := elt.(*ast.KeyValueExpr)
			if err := c.compile(kve.Key); err != nil {
				return err
			}
			if err := c.compile(kve.Value); err != nil {
				return err
			}
		}
		if _, err := c.scopes.Emit(bytecode.OpMap, len(node.Elts)*2); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported composite literal type: %T", node.Type)
	}
}

// BasicLit processes an AST BasicLit node and emits the corresponding literal into the current scope.
func (c *Declarations) BasicLit(node *ast.BasicLit) error {
	var obj objects.IObject
	switch node.Kind {
	case token.INT:
		val, _ := strconv.ParseInt(node.Value, 0, 64)
		obj = c.gk.NewInt(objects.FrameStatic, val)
	case token.FLOAT:
		val, _ := strconv.ParseFloat(node.Value, 64)
		obj = c.gk.NewFloat(objects.FrameStatic, val)
	case token.STRING:
		val, _ := strconv.Unquote(node.Value)
		obj = c.gk.NewString(objects.FrameStatic, val)
	default:
		return fmt.Errorf("unhandled literal: %s", node.Kind)
	}
	id := c.constants.Add("", obj)
	if _, err := c.scopes.Emit(bytecode.OpConstant, id); err != nil {
		return err
	}
	return nil
}

// Ident processes an identifier node, resolving its symbol in the current scope and emitting a symbol get operation.
func (c *Declarations) Ident(node *ast.Ident) error {
	switch node.Name {
	case "true":
		if _, err := c.scopes.Emit(bytecode.OpTrue); err != nil {
			return err
		}
		return nil
	case "false":
		if _, err := c.scopes.Emit(bytecode.OpFalse); err != nil {
			return err
		}
		return nil
	}
	symbol, ok := c.scopes.SymbolResolve(node.Name)
	if !ok {
		return fmt.Errorf("undefined variable: %s", node.Name)
	}
	if err := c.scopes.EmitSymbolGet(symbol); err != nil {
		return err
	}
	return nil
}

// AssignStmt processes an assignment statement by compiling the right-hand side and resolving variable symbols.
// It also updates the type information for symbols or emits appropriate bytecode for assignments.
func (c *Declarations) AssignStmt(node *ast.AssignStmt) error {
	//TODO UNIFY function-only and multi assignment

	// multi assignment function-only (es. x, y := f())
	if callExpr, ok := node.Rhs[0].(*ast.CallExpr); ok && len(node.Lhs) > 1 {
		if err := c.compile(callExpr); err != nil {
			return err
		}
		var funcReturnTypes []string
		if ident, isIdent := callExpr.Fun.(*ast.Ident); isIdent {
			if ident.Name != "" {
				if funcSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok {
					funcReturnTypes = funcSymbol.Types()
				}
			}
		}
		if len(node.Lhs) != len(funcReturnTypes) {
			return fmt.Errorf("assignment mismatch: %d variables but %d return values", len(node.Lhs), len(funcReturnTypes))
		}
		for i := len(node.Lhs) - 1; i >= 0; i-- {
			lhs := node.Lhs[i]
			ident, ok := lhs.(*ast.Ident)
			if !ok {
				return fmt.Errorf("unsupported multiple assignment to type %T", lhs)
			}
			var symbol *Symbol
			if node.Tok == token.DEFINE {
				var err error
				symbol, err = c.scopes.SymbolDefine(ident.Name, UnknownScope, false)
				if err != nil {
					return err
				}
			} else {
				var found bool
				symbol, found = c.scopes.SymbolResolve(ident.Name)
				if !found {
					return fmt.Errorf("undefined variable: %s", ident.Name)
				}
			}
			symbol.SetTypes([]string{funcReturnTypes[i]})
			if err := c.scopes.EmitSymbolSet(symbol); err != nil {
				return err
			}
			if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
				return err
			}
		}
		return nil
	}
	//single assignment
	if err := c.compile(node.Rhs[0]); err != nil {
		return err
	}
	//inference type check
	var assignedTypeName []string
	switch rhs := node.Rhs[0].(type) {
	case *ast.CompositeLit: // check for variable assignment
		if ident, ok := rhs.Type.(*ast.Ident); ok {
			if typeSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok && typeSymbol.IsStruct() {
				assignedTypeName = []string{typeSymbol.Name()}
			}
		}
	case *ast.CallExpr: // check for function call assignment
		if ident, ok := rhs.Fun.(*ast.Ident); ok {
			if funcSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok && len(funcSymbol.Types()) > 0 {
				assignedTypeName = []string{funcSymbol.Types()[0]}
			}
		}
	}

	switch lhs := node.Lhs[0].(type) {
	case *ast.Ident:
		name := lhs.Name
		var symbol *Symbol
		if node.Tok == token.DEFINE {
			var err error
			symbol, err = c.scopes.SymbolDefine(name, UnknownScope, false)
			if err != nil {
				return err
			}
		} else {
			var ok bool
			symbol, ok = c.scopes.SymbolResolve(name)
			if !ok {
				return fmt.Errorf("undefined variable: %s", name)
			}
		}
		// Updates the symbol type in both cases (:= and =)
		if len(assignedTypeName) > 0 {
			symbol.SetTypes(assignedTypeName)
		}
		if err := c.scopes.EmitSymbolSet(symbol); err != nil {
			return err
		}
		// OpPop is only needed for simple variable assignments because OpSetLocal/Global don't clean up the stack.
		if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
		return nil
	case *ast.SelectorExpr:
		if node.Tok == token.DEFINE {
			return fmt.Errorf("cannot define a field with :=")
		}
		receiverIdent, ok := lhs.X.(*ast.Ident)
		if !ok {
			return fmt.Errorf("unsupported receiver for field assignment")
		}
		symbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
		if !ok {
			return fmt.Errorf("undefined variable: %s", receiverIdent.Name)
		}
		fieldName := lhs.Sel.Name
		keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
		if _, err := c.scopes.Emit(bytecode.OpConstant, keyConst); err != nil {
			return err
		}
		if symbol.Scope() == GlobalScope {
			if _, err := c.scopes.Emit(bytecode.OpSetSelGlobal, symbol.Index(), 1); err != nil {
				return err
			}
		} else {
			if _, err := c.scopes.Emit(bytecode.OpSetSelLocal, symbol.Index(), 1); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported left-hand side in assignment: %T", node.Lhs[0])
	}
}
