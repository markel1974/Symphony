package compiler

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/markel1974/c64emu/src/compilers/native/tables"
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// Imports is a structure that manages a set of imported items, built-in functions, and related compilation resources.
type Imports struct {
	gk         objects.IGateKeeper
	loader     bytecode.ILoader
	references *Constants
	scopes     *tables.Scopes
	imports    map[string]bool
	builtin    map[string]objects.CallId
	container  []ast.Decl
	fileSet    *token.FileSet
	compile    func(node ast.Node) error
}

// NewImports creates and initializes a new Imports instance with provided GateKeeper, Constants, and Scopes references.
func NewImports(gk objects.IGateKeeper, loader bytecode.ILoader, references *Constants, scopes *tables.Scopes) *Imports {
	i := &Imports{
		gk:         gk,
		loader:     loader,
		references: references,
		scopes:     scopes,
		imports:    make(map[string]bool),
		builtin:    make(map[string]objects.CallId),
		compile:    nil,
	}
	return i
}

// Setup initializes the Imports instance by configuring compilation and loading built-in functions from the provided loader.
func (i *Imports) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	i.fileSet = fileSet
	i.compile = compile

	i.builtin = map[string]objects.CallId{
		"len":         objects.CallIdLen,
		"copy":        objects.CallIdCopy,
		"append":      objects.CallIdAppend,
		"delete":      objects.CallIdDelete,
		"splice":      objects.CallIdSplice,
		"panic":       objects.CallIdPanic,
		"recover":     objects.CallIdRecover,
		"int":         objects.CallIdInt,
		"bool":        objects.CallIdBool,
		"float":       objects.CallIdFloat,
		"char":        objects.CallIdChar,
		"string":      objects.CallIdString,
		"time":        objects.CallIdTime,
		"typeName":    objects.CallIdTypeName,
		"isInt":       objects.CallIdIsInt,
		"isFloat":     objects.CallIdIsFloat,
		"isString":    objects.CallIdIsString,
		"isBool":      objects.CallIdIsBool,
		"isChar":      objects.CallIdIsChar,
		"isBytes":     objects.CallIdIsBytes,
		"isArray":     objects.CallIdIsArray,
		"isMap":       objects.CallIdIsMap,
		"isIterable":  objects.CallIdIsIterable,
		"isTime":      objects.CallIdIsTime,
		"isError":     objects.CallIdIsError,
		"isUndefined": objects.CallIdIsUndefined,
		"isFunction":  objects.CallIdIsFunction,
		"isCallable":  objects.CallIdIsCallable,
		"printf":      objects.CallIdPrintf,
		"sprintf":     objects.CallIdSprintf,
		"make":        objects.CallIdMake,
	}

	/*


		for idx := 0; idx < i.loader.BuiltinLen(); idx++ {
			bi := i.loader.Builtin(idx)
			if bi == nil {
				return fmt.Errorf("builtin %d not found", idx)
			}
			builtinId := i.references.Add(bi.Name(), bi)
			i.builtin[bi.Name()] = builtinId
		}

	*/
	return nil
}

// Prepare initializes the Imports instance, ensuring it is ready for further use in the compilation process.
func (i *Imports) Prepare() error {
	return nil
}

// Compile processes and compiles all stored declarations. Returns an error if any declaration fails to compile.
func (i *Imports) Compile() error {
	for _, decl := range i.container {
		if err := i.compile(decl); err != nil {
			return err
		}
	}
	return nil
}

// Declare appends the provided declaration to the container slice of the Imports structure.
func (i *Imports) Declare(decls ast.Decl) {
	i.container = append(i.container, decls)
}

// HasPackage checks if the specified package name exists in the imports map and returns true if it exists, otherwise false.
func (i *Imports) HasPackage(name string) bool {
	return i.imports[name]
}

// HasBuiltin checks if the given name exists in the builtin map and returns true if found, otherwise false.
func (i *Imports) HasBuiltin(name string) bool {
	_, ok := i.builtin[name]
	return ok
}

// Emit attempts to attach a function reference or emit a builtin reference, returning true if successful.
func (i *Imports) Emit(name string, selName string) bool {
	if len(selName) > 0 {
		_, ok := i.imports[name]
		if !ok {
			return false
		}
		_, nameIndex, err := i.packageFunctionAttach(name, selName)
		if err != nil {
			return false
		}
		if _, err = i.scopes.Emit(bytecode.OpFuncImport, nameIndex); err != nil {
			return false
		}
		return true
	}
	id, ok := i.builtin[name]
	if !ok {
		return false
	}
	if _, err := i.scopes.Emit(bytecode.OpFuncInternal, int(id)); err != nil {
		return false
	}
	return true
}

// Attach attempts to resolve a name and optional selector from imports or built-in references, returning details if found.
func (i *Imports) Attach(name string, selName string) (string, int, bool) {
	if len(selName) > 0 {
		_, ok := i.imports[name]
		if !ok {
			return "", -1, false
		}
		mangledName, nameIndex, err := i.packageFunctionAttach(name, selName)
		if err != nil {
			return "", 0, false
		}
		return mangledName, nameIndex, true
	}
	id, ok := i.builtin[name]
	if !ok {
		return "", 0, false
	}
	return name, int(id), true
}

// PackageFunctionAttach registers and attaches a function from a given package, returning its mangled name, index, and any error.
func (i *Imports) packageFunctionAttach(pkgName string, fnName string) (string, int, error) {
	mangledName := tables.GetMangledName(pkgName, fnName)
	nameIndex, found := i.references.Get(mangledName)
	if !found {
		attrArray := i.gk.NewArray(objects.FrameStatic, []objects.IObject{
			i.gk.NewString(objects.FrameStatic, pkgName),
			i.gk.NewString(objects.FrameStatic, fnName)},
		)
		nameIndex = i.references.Add(mangledName, attrArray)
	}
	return mangledName, nameIndex, nil
}

// doGenDecl processes a generic declaration node, compiling each specification it contains. Returns an error if compilation fails.
func (i *Imports) doGenDecl(node *ast.GenDecl) error {
	for _, spec := range node.Specs {
		if err := i.compile(spec); err != nil {
			return err
		}
	}
	return nil
}

// ImportSpec processes an ast.ImportSpec node, extracts the module name, and adds it to the imports map as a key.
func (i *Imports) ImportSpec(node *ast.ImportSpec) error {
	moduleName := node.Path.Value
	moduleName = strings.Trim(moduleName, "\"'")
	i.imports[moduleName] = true
	return nil
}
