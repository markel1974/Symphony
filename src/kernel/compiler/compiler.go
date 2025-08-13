package compiler

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
	"github.com/markel1974/c64emu/src/kernel/vm/stdlib"
)

const (
	mainFnName = "main"
)

// maxScope defines the maximum allowable depth for scopes in the compiler to prevent excessive nesting or stack overflow.
const (
	maxScope = 1024
)

// Compiler manages the compilation process, including constant storage, scopes, and symbol resolution during program compilation.
type Compiler struct {
	constants     []objects.IObject
	symbolTable   *SymbolTable
	scopes        []*CompilationScope
	scopeIndex    int
	mainFn        *objects.FunctionCompiled
	constantCache map[string]int
}

// New creates and initializes a new Compiler instance with a fresh symbol table and main compilation scope.
func New() *Compiler {
	mainScope := NewCompilationScope()
	symbolTable := NewSymbolTable()
	for i, fn := range stdlib.GetAllBuiltinFunctions() {
		symbolTable.DefineBuiltin(fn.Name(), i)
	}
	return &Compiler{
		constants:     []objects.IObject{},
		symbolTable:   symbolTable,
		scopes:        []*CompilationScope{mainScope},
		scopeIndex:    0,
		mainFn:        nil,
		constantCache: make(map[string]int),
	}
}

// Compile traverses the provided AST node and compiles it into bytecode for execution by the virtual machine.
func (c *Compiler) Compile(in ast.Node) error {
	switch node := in.(type) {
	case *ast.File:
		for _, s := range node.Decls {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
	case *ast.DeclStmt:
		if err := c.Compile(node.Decl); err != nil {
			return err
		}
	case *ast.GenDecl: // For `var` and `const` which are handled by AssignStmt
		for _, spec := range node.Specs {
			if err := c.Compile(spec); err != nil {
				return err
			}
		}
	case *ast.ValueSpec: // Handles 'var x = 10'
		for i, name := range node.Names {
			if err := c.Compile(node.Values[i]); err != nil {
				return err
			}
			symbol := c.symbolTable.Define(name.Name)
			if err := c.emitSymbolDefine(symbol); err != nil {
				return err
			}
		}

	// --- Statements ---
	case *ast.BlockStmt:
		for _, s := range node.List {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
	case *ast.ExprStmt:
		if err := c.Compile(node.X); err != nil {
			return err
		}
		// Remove value from stack if unused
		if _, err := c.emit(bytecode.OpPop); err != nil {
			return err
		}
	case *ast.AssignStmt:
		for i, expr := range node.Rhs {
			if err := c.Compile(expr); err != nil {
				return err
			}
			ident := node.Lhs[i].(*ast.Ident)
			if node.Tok == token.DEFINE { // Handles 'x := 10'
				symbol := c.symbolTable.Define(ident.Name)
				// AND USE THE NEW FUNCTION HERE TOO
				if err := c.emitSymbolDefine(symbol); err != nil {
					return err
				}
			} else { // Handles 'x = 20'
				symbol, ok := c.symbolTable.Resolve(ident.Name)
				if !ok {
					return fmt.Errorf("undefined variable: %s", ident.Name)
				}
				// Assignment continues to use the old function
				if err := c.emitSymbolSet(symbol); err != nil {
					return err
				}
			}
		}
	case *ast.IfStmt:
		// Compile condition
		if err := c.Compile(node.Cond); err != nil {
			return err
		}
		// Emit conditional jump with temporary address
		jumpNotTruthyPos, err := c.emit(bytecode.OpJumpFalsy, 9999)
		if err != nil {
			return err
		}
		// Compile 'then' block
		if err = c.Compile(node.Body); err != nil {
			return err
		}
		// If there's an 'else' block, emit jump to skip it
		jumpToEndPos := 0
		if node.Else != nil {
			jumpToEndPos, err = c.emit(bytecode.OpJump, 9999)
			if err != nil {
				return err
			}
		}
		scope, err := c.scopeCurrent()
		if err != nil {
			return err
		}
		// Update conditional jump address
		if err = c.changeOperand(jumpNotTruthyPos, scope.InstructionsLen()); err != nil {
			return err
		}
		// Compile 'else' block if it exists
		if node.Else != nil {
			if err = c.Compile(node.Else); err != nil {
				return err
			}
			scope, err = c.scopeCurrent()
			if err != nil {
				return err
			}
			// Update jump address to skip else
			if err = c.changeOperand(jumpToEndPos, scope.InstructionsLen()); err != nil {
				return err
			}
		}
	case *ast.ForStmt:
		scope, err := c.scopeCurrent()
		if err != nil {
			return err
		}
		// Starting position for loop condition
		loopStartPos := scope.InstructionsLen()
		// Compile condition
		if err = c.Compile(node.Cond); err != nil {
			return err
		}
		// Emit conditional jump to exit loop
		jumpNotTruthyPos, err := c.emit(bytecode.OpJumpFalsy, 9999)
		if err != nil {
			return err
		}
		// Compile loop body
		if err = c.Compile(node.Body); err != nil {
			return err
		}
		// Emit unconditional jump to return to condition start
		if _, err = c.emit(bytecode.OpJump, loopStartPos); err != nil {
			return err
		}
		// Update (back-patching) OpJumpFalsy address to point to loop end
		scope, err = c.scopeCurrent()
		if err != nil {
			return err
		}
		afterLoopPos := scope.InstructionsLen()
		if err = c.changeOperand(jumpNotTruthyPos, afterLoopPos); err != nil {
			return err
		}
		// Remove condition value from stack after loop terminates
		if _, err = c.emit(bytecode.OpPop); err != nil {
			return err
		}
	case *ast.BinaryExpr:
		if err := c.Compile(node.X); err != nil {
			return err
		}
		if err := c.Compile(node.Y); err != nil {
			return err
		}
		if err := c.emitBinaryOp(node.Op); err != nil {
			return err
		}
	case *ast.UnaryExpr:
		if err := c.Compile(node.X); err != nil {
			return err
		}
		if err := c.emitUnaryOp(node.Op); err != nil {
			return err
		}
	case *ast.BasicLit:
		if err := c.compileLiteral(node); err != nil {
			return err
		}
	case *ast.Ident:
		symbol, ok := c.symbolTable.Resolve(node.Name)
		if !ok {
			return fmt.Errorf("undefined variable: %s", node.Name)
		}
		if err := c.emitSymbolGet(symbol); err != nil {
			return err
		}
	case *ast.CompositeLit:
		// Check type to determine if array or map
		switch node.Type.(type) {
		case *ast.ArrayType:
			// Array literal (e.g. []int{1, 2, 3})
			for _, elt := range node.Elts {
				if err := c.Compile(elt); err != nil {
					return err
				}
			}
			if _, err := c.emit(bytecode.OpArray, len(node.Elts)); err != nil {
				return err
			}
		case *ast.MapType:
			// Map literal (e.g. map[string]int{"a": 1})
			for _, elt := range node.Elts {
				kve := elt.(*ast.KeyValueExpr)
				if err := c.Compile(kve.Key); err != nil {
					return err
				}
				if err := c.Compile(kve.Value); err != nil {
					return err
				}
			}
			if _, err := c.emit(bytecode.OpMap, len(node.Elts)); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unsupported composite literal type")
		}
	case *ast.FuncDecl:
		//fmt.Println("Compiling FuncDecl:", node.Name.Name)
		// Global function declaration
		if err := c.scopeEnter(); err != nil {
			return err
		}
		// Receiver (methods) not supported in this version Parameters
		for _, p := range node.Type.Params.List {
			for _, name := range p.Names {
				c.symbolTable.Define(name.Name)
			}
		}
		// Body
		if err := c.Compile(node.Body); err != nil {
			return err
		}
		// Implicit return if missing
		if _, err := c.emit(bytecode.OpReturn, 0); err != nil {
			return err
		}
		numLocals := c.symbolTable.NumDefinitions()
		instructions, err := c.scopeLeave()
		if err != nil {
			return err
		}

		numParams := 0
		varArgs := false
		if paramL := node.Type.Params.List; paramL != nil {
			if numParams = len(paramL); numParams > 0 {
				if _, ok := paramL[numParams-1].Type.(*ast.Ellipsis); ok {
					varArgs = true
				}
			}
		}
		// Create compiled function object
		//TODO sourceMap
		compiledFn := objects.NewFunctionCompiled(node.Name.String(), instructions, numLocals, numParams, varArgs, nil, c.symbolTable.ConvertFreeSymbols())
		fnIndex := c.addConstant("", compiledFn)
		if _, err = c.emit(bytecode.OpClosure, fnIndex, c.symbolTable.FreeSymbolsLen()); err != nil {
			return err
		}
		// Define function in current scope
		symbol := c.symbolTable.Define(node.Name.Name)
		if err = c.emitSymbolSet(symbol); err != nil {
			return err
		}
		if node.Name.Name == mainFnName {
			c.mainFn = compiledFn
		}
	case *ast.CallExpr:
		// Compile function to call (e.g. identifier)
		if err := c.Compile(node.Fun); err != nil {
			return err
		}
		// Compile arguments
		for _, arg := range node.Args {
			if err := c.Compile(arg); err != nil {
				return err
			}
		}
		// 0 for non-spread call
		if _, err := c.emit(bytecode.OpCall, len(node.Args), 0); err != nil {
			return err
		}
	case *ast.ReturnStmt:
		if len(node.Results) == 0 {
			// Return 'undefined'
			if _, err := c.emit(bytecode.OpReturn, 0); err != nil {
				return err
			}
		} else {
			if err := c.Compile(node.Results[0]); err != nil {
				return err
			}
			// Return a value
			if _, err := c.emit(bytecode.OpReturn, 1); err != nil {
				return err
			}
		}
	case *ast.SelectorExpr:
		// 1. Compila l'espressione a sinistra (es. 'fmt').
		// Questo emetterà un 'OpGetGlobal' che a runtime metterà
		// l'oggetto modulo MapImmutable sullo stack.
		if err := c.Compile(node.X); err != nil {
			return err
		}
		moduleIdent, ok := node.X.(*ast.Ident)
		if !ok {
			// Per ora, gestiamo solo il caso semplice come 'fmt.print'
			// e non casi complessi come 'a[0].print()'.
			return fmt.Errorf("unsupported selector expression: %T", node.X)
		}
		mName := moduleIdent.Name
		sName := node.Sel.Name
		cacheKey := "selector:" + mName + "." + sName
		nameIndex, found := c.constantCache[cacheKey]
		if !found {
			attrArray := objects.NewArray([]objects.IObject{objects.NewStringNoSize(mName), objects.NewStringNoSize(sName)})
			nameIndex = c.addConstant(cacheKey, attrArray)
		}
		if _, err := c.emit(bytecode.OpGetAttr, nameIndex); err != nil {
			return err
		}
		/*
			//nameIndex := c.addConstant(objects.NewStringNoSize(mName + "." + sName))
			symbol := objects.NewArray([]objects.IObject{objects.NewStringNoSize(mName), objects.NewStringNoSize(sName)})
			nameIndex := c.addConstant(symbol)
			if _, err := c.emit(bytecode.OpGetAttr, nameIndex); err != nil {
				return err
			}

		*/
	case *ast.ImportSpec:
		moduleName := node.Path.Value
		c.symbolTable.Define(strings.Trim(moduleName, "\"'"))
	default:
		return fmt.Errorf("unsupported expression type: %T", node)
	}
	return nil
}

// Bytecode generates and returns the compiled bytecode representation from the compiler's current state.
// It encapsulates the main function and constants into the bytecode structure to produce an executable output.
// Returns an error if there's an issue retrieving the current compilation scope.
func (c *Compiler) Bytecode() (*bytecode.Bytecode, error) {
	bc := bytecode.NewBytecode()
	if c.mainFn == nil {
		return nil, errors.New("main function not found")
	}
	bc.SetMainFunction(c.mainFn)
	bc.SetConstants(c.constants)
	return bc, nil
}

// compileLiteral compiles a basic literal node into bytecode and adds the constant to the constant pool.
// It handles integer, float, and string literal types, emitting the corresponding constant operation.
// Returns an error if the literal type is not supported or an issue occurs during bytecode emission.
func (c *Compiler) compileLiteral(node *ast.BasicLit) error {
	switch node.Kind {
	case token.INT:
		val, _ := strconv.ParseInt(node.Value, 0, 64)
		if _, err := c.emit(bytecode.OpConstant, c.addConstant("", objects.NewInt(val))); err != nil {
			return err
		}
	case token.FLOAT:
		val, _ := strconv.ParseFloat(node.Value, 64)
		if _, err := c.emit(bytecode.OpConstant, c.addConstant("", objects.NewFloat(val))); err != nil {
			return err
		}
	case token.STRING:
		val, _ := strconv.Unquote(node.Value)
		s, err := objects.NewString(val)
		if err != nil {
			return err
		}
		if _, err = c.emit(bytecode.OpConstant, c.addConstant("", s)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unhandled literal: %s", node.Kind)
	}
	return nil
}

// addConstant adds a new object to the constants pool and returns its index in the pool.
func (c *Compiler) addConstant(id string, obj objects.IObject) int {
	c.constants = append(c.constants, obj)
	nameIndex := len(c.constants) - 1
	if len(id) > 0 {
		c.constantCache[id] = nameIndex
	}
	return nameIndex
}

// addInstructions adds a sequence of byte instructions to the current compilation scope.
// It returns the position of the newly appended instructions or an error if appending fails.
func (c *Compiler) addInstructions(ins []byte) (int, error) {
	scope, err := c.scopeCurrent()
	if err != nil {
		return 0, err
	}
	posNewInstruction := scope.InstructionsLen()
	if err = scope.InstructionsAppend(ins); err != nil {
		return 0, err
	}
	return posNewInstruction, nil
}

// setLastInstruction updates the last instruction emitted in the current compilation scope with the given opcode and position.
func (c *Compiler) setLastInstruction(op bytecode.Opcode, pos int) error {
	scope, err := c.scopeCurrent()
	if err != nil {
		return err
	}
	previous := scope.LastInstruction()
	last := NewEmittedInstruction(op, pos)
	scope.SetPreviousInstruction(previous)
	scope.SetLastInstruction(last)
	return nil
}

// changeOperand modifies the operand of an instruction at the specified position within the current scope's instructions.
func (c *Compiler) changeOperand(opPos int, operand int) error {
	op, err := c.instructionGet(opPos)
	if err != nil {
		return err
	}
	newInstruction := bytecode.MakeInstruction(op, operand)
	if err = c.instructionReplace(opPos, newInstruction); err != nil {
		return err
	}
	return nil
}

// instructionReplace modifies the current instruction set starting from the given position with the provided new instruction.
// It returns an error if the modification cannot be applied.
func (c *Compiler) instructionReplace(pos int, newInstruction []byte) error {
	scope, err := c.scopeCurrent()
	if err != nil {
		return err
	}
	if err = scope.InstructionsReplace(pos, newInstruction); err != nil {
		return err
	}
	return nil
}

// instructionSet sets an instruction at a specific position in the current scope's instructions.
func (c *Compiler) instructionSet(pos int, instruction byte) error {
	scope, err := c.scopeCurrent()
	if err != nil {
		return err
	}
	if err = scope.InstructionsSet(pos, instruction); err != nil {
		return err
	}
	return nil
}

// instructionGet retrieves a byte of instruction data at the specified position from the current scope.
func (c *Compiler) instructionGet(pos int) (byte, error) {
	scope, err := c.scopeCurrent()
	if err != nil {
		return 0, err
	}
	data, err := scope.InstructionsGet(pos)
	if err != nil {
		return 0, err
	}
	return data, nil
}

// scopeCurrent retrieves the current CompilationScope of the compiler. Returns an error if the scope index is invalid.
func (c *Compiler) scopeCurrent() (*CompilationScope, error) {
	if c.scopeIndex < 0 || c.scopeIndex >= len(c.scopes) {
		return nil, fmt.Errorf("invalid scope index: %d", c.scopeIndex)
	}
	return c.scopes[c.scopeIndex], nil
}

// scopeEnter adds a new compilation scope to the stack and updates the symbol table with an enclosed scope.
// Returns an error if the maximum scope depth is exceeded.
func (c *Compiler) scopeEnter() error {
	if c.scopeIndex > maxScope {
		return fmt.Errorf("maximum scope depth exceeded: %d", maxScope)
	}
	scope := NewCompilationScope() //CompilationScope{instructions: []byte{}, lastInstruction: EmittedInstruction{}, previousInstruction: EmittedInstruction{}}
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
	return nil
}

// scopeLeave exits the current compilation scope, returning its bytecode instructions and any error encountered.
func (c *Compiler) scopeLeave() ([]byte, error) {
	scopesL := len(c.scopes)
	if scopesL <= 0 {
		return nil, errors.New("no scopes to leave")
	}
	scope, err := c.scopeCurrent()
	if err != nil {
		return nil, err
	}
	c.scopes = c.scopes[:scopesL-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer()
	return scope.Instructions(), nil
}

// emitSymbolSet generates bytecode instructions to set the value of a symbol in its appropriate scope (global, local, or free).
func (c *Compiler) emitSymbolSet(s *Symbol) error {
	//fmt.Println("Emitting Symbol:", s)
	switch s.Scope {
	case GlobalScope:
		if _, err := c.emit(bytecode.OpSetGlobal, s.Index); err != nil {
			return err
		}
	case LocalScope:
		if _, err := c.emit(bytecode.OpSetLocal, s.Index); err != nil {
			return err
		}
	case FreeScope:
		if _, err := c.emit(bytecode.OpSetFree, s.Index); err != nil {
			return err
		}
	}
	return nil
}

// emitSymbolDefine emits the opcode for *defining* a variable.
func (c *Compiler) emitSymbolDefine(s *Symbol) error {
	switch s.Scope {
	case GlobalScope:
		// For global scope, define is same as assign
		if _, err := c.emit(bytecode.OpSetGlobal, s.Index); err != nil {
			return err
		}
	case LocalScope:
		// Use new opcode for local variables
		if _, err := c.emit(bytecode.OpDefineLocal, s.Index); err != nil {
			return err
		}
	}
	return nil
}

// emitSymbolGet generates bytecode instructions to retrieve a symbol's value based on its scope and index.
func (c *Compiler) emitSymbolGet(s *Symbol) error {
	switch s.Scope {
	case GlobalScope:
		if _, err := c.emit(bytecode.OpGetGlobal, s.Index); err != nil {
			return err
		}
	case LocalScope:
		if _, err := c.emit(bytecode.OpGetLocal, s.Index); err != nil {
			return err
		}
	case BuiltinScope:
		if _, err := c.emit(bytecode.OpGetBuiltin, s.Index); err != nil {
			return err
		}
	case FreeScope:
		if _, err := c.emit(bytecode.OpGetFree, s.Index); err != nil {
			return err
		}
	}
	return nil
}

// emitBinaryOp compiles a binary operation by emitting the corresponding bytecode based on the provided token operator.
func (c *Compiler) emitBinaryOp(op token.Token) error {
	switch op {
	case token.ADD:
		if _, err := c.emit(bytecode.OpBinaryOp, int(objects.OperatorAdd)); err != nil {
			return err
		}
	case token.SUB:
		if _, err := c.emit(bytecode.OpBinaryOp, int(objects.OperatorSub)); err != nil {
			return err
		}
	case token.MUL:
		if _, err := c.emit(bytecode.OpBinaryOp, int(objects.OperatorMul)); err != nil {
			return err
		}
	case token.QUO:
		if _, err := c.emit(bytecode.OpBinaryOp, int(objects.OperatorQuo)); err != nil {
			return err
		}
	case token.EQL:
		if _, err := c.emit(bytecode.OpEqual); err != nil {
			return err
		}
	case token.NEQ:
		if _, err := c.emit(bytecode.OpNotEqual); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unhandled binary op: %s", op)
	}
	return nil
}

// emitUnaryOp compiles unary operations by emitting the corresponding bytecode for the given operator.
func (c *Compiler) emitUnaryOp(op token.Token) error {
	switch op {
	case token.SUB:
		if _, err := c.emit(bytecode.OpMinus); err != nil {
			return err
		}
	case token.NOT: // !
		if _, err := c.emit(bytecode.OpLNot); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unhandled unary op: %s", op)
	}
	return nil
}

// emit generates and adds a bytecode instruction to the instruction sequence, storing the opcode and operands.
func (c *Compiler) emit(op bytecode.Opcode, operands ...int) (int, error) {
	ins := bytecode.MakeInstruction(op, operands...)
	pos, err := c.addInstructions(ins)
	if err != nil {
		return 0, err
	}
	if err = c.setLastInstruction(op, pos); err != nil {
		return 0, err
	}
	return pos, nil
}
