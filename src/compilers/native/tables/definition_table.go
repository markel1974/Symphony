package tables

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// DefinitionTable is a composite structure for managing scopes, structs, interfaces, and access control through an IGateKeeper.
type DefinitionTable struct {
	gk             objects.IGateKeeper
	scopes         *Scopes
	structTable    *StructTable
	interfaceTable *InterfaceTable
}

// NewDefinitionTable initializes and returns a pointer to a DefinitionTable with the provided IGateKeeper and Scopes instances.
func NewDefinitionTable(gk objects.IGateKeeper, scopes *Scopes) *DefinitionTable {
	structTable := NewStructTable(gk, scopes)
	interfaceTable := NewInterfaceTable(gk, scopes)
	return &DefinitionTable{
		gk:             gk,
		scopes:         scopes,
		structTable:    structTable,
		interfaceTable: interfaceTable,
	}
}

// Assign assigns a type to the given symbol based on whether the type is an interface or not.
func (f *DefinitionTable) Assign(symbol *Symbol, typeName string) error {
	isInterface := f.interfaceTable.Has(typeName)
	if isInterface {
		f.InterfaceAssign(symbol, typeName)
	} else {
		f.TypeAssign(symbol, typeName)
	}
	return nil
}

// InferAssign assigns a type to the given symbol by inferring the type from an expression if not explicitly provided.
func (f *DefinitionTable) InferAssign(symbol *Symbol, inferredTypeName string, rhsIn ast.Expr) {
	if len(inferredTypeName) == 0 {
		if inferredTypeName, _ = f.TypeInference(rhsIn); len(inferredTypeName) == 0 {
			return
		}
	}
	f.TypeAssign(symbol, inferredTypeName)
}

// InterfaceAssign assigns an interface type to the given symbol, updating its return types, interface, and object representation.
func (f *DefinitionTable) InterfaceAssign(symbol *Symbol, interfaceName string) {
	symbol.SetReturnTypes([]string{interfaceName})
	symbol.SetInterface(interfaceName)
	symbol.SetObject(f.gk.NewString(objects.FrameStatic, interfaceName+":"+symbol.Name()))
}

// TypeAssign assigns a type to a symbol by setting its return types and object, and binds it to a struct if applicable.
func (f *DefinitionTable) TypeAssign(symbol *Symbol, typeName string) {
	symbol.SetReturnTypes([]string{typeName})
	symbol.SetObject(f.gk.NewString(objects.FrameStatic, typeName+":"+symbol.Name()))
	//assign struct if present in struct table
	f.structTable.BindSymbol(symbol, typeName)
}

// StructBindSymbol binds a given struct symbol to its corresponding type name in the struct table.
func (f *DefinitionTable) StructBindSymbol(symbol *Symbol, typeName string) {
	symbol.SetObject(f.gk.NewString(objects.FrameStatic, typeName+":"+symbol.Name()))
	f.structTable.BindSymbol(symbol, typeName)
}

// StructHas checks if a struct with the given name exists in the struct table.
func (f *DefinitionTable) StructHas(name string) bool {
	return f.structTable.Has(name)
}

// StructIsBuiltin checks if a struct with the given name is considered a built-in type in the struct table.
func (f *DefinitionTable) StructIsBuiltin(name string) bool {
	return f.structTable.IsBuiltin(name)
}

// StructAdd adds a new struct with the given name to the struct table.
func (f *DefinitionTable) StructAdd(name string) {
	f.structTable.AddStruct(name)
}

// StructAddField adds a new field to a struct definition in the struct table using the provided field details and base type.
func (f *DefinitionTable) StructAddField(name string, fieldName string, baseStruct string, node ast.Node, kind ast.Expr) {
	// Surface Check
	isPointer := false
	//isMap := false
	//isArray := false
	switch kind.(type) {
	case *ast.StarExpr:
		isPointer = true
	case *ast.MapType:
		//isMap = true
		//fmt.Println("Is MapType")
	case *ast.ArrayType:
		//isArray = true
		//fmt.Println("Is ArrayType")
	}

	sd := f.structTable.AddStruct(name)
	var field IStructField = nil
	if s, ok := f.structTable.Get(baseStruct); ok {
		field = s.FieldClone()
	} else if i, ok := f.interfaceTable.Get(baseStruct); ok {
		field = i.FieldClone()
	} else {
		field = NewStructField(baseStruct)
	}
	field.SetFieldName(fieldName)
	field.SetFieldNode(node)
	field.SetIsPointer(isPointer)
	sd.AddField(field)
}

// StructKeys returns a list of all the keys present in the struct table.
func (f *DefinitionTable) StructKeys() []string {
	return f.structTable.Keys()
}

// StructSetImplementations sets the struct implementations mapping in the struct table.
func (f *DefinitionTable) StructSetImplementations(impls map[string][]string) {
	f.structTable.SetImplementations(impls)
}

// StructFieldsFromLiteral retrieves the fields of a struct from its literal representation using provided expressions.
func (f *DefinitionTable) StructFieldsFromLiteral(structName string, eltS []ast.Expr) ([]IStructField, error) {
	return f.structTable.FieldsFromLiteral(structName, eltS)
}

// StructTypeNameFromSymbolField retrieves the type name of a struct field by its symbol and field name.
// Returns the type name and a boolean indicating success.
func (f *DefinitionTable) StructTypeNameFromSymbolField(name string, fieldName string) (string, bool) {
	return f.structTable.TypeNameFromSymbolField(name, fieldName)
}

// StructImplements checks if a given struct implements a specified interface and returns a boolean result.
func (f *DefinitionTable) StructImplements(structName string, interfaceName string) bool {
	return f.structTable.Implements(structName, interfaceName)
}

// InterfaceAdd adds a new interface with the specified name and AST node to the interface table. Returns an error if it fails.
func (f *DefinitionTable) InterfaceAdd(name string, node *ast.InterfaceType) error {
	return f.interfaceTable.Add(name, node)
}

// InterfaceGet retrieves the interface definition by name from the interface table and indicates its existence.
func (f *DefinitionTable) InterfaceGet(name string) (*InterfaceDescription, bool) {
	return f.interfaceTable.Get(name)
}

// InterfaceContainer retrieves the internal container that maps interface names to their descriptions.
func (f *DefinitionTable) InterfaceContainer() map[string]*InterfaceDescription {
	return f.interfaceTable.Container()
}

// TypeInference determines the type of the given expression and returns the type name and a boolean indicating success.
func (f *DefinitionTable) TypeInference(expr ast.Expr) (string, bool) {
	var ret string
	switch rhs := expr.(type) {
	case *ast.Ident:
		if v, ok := f.scopes.SymbolResolve(rhs.Name); ok {
			ret = v.StructName()
		}
	case *ast.CompositeLit: // es. MyStruct{...}
		if ident := GetIdent(rhs.Type); ident != nil {
			ret = ident.Name
		}
	case *ast.CallExpr: // es. NewStruct()
		if ident := GetIdent(rhs.Fun); ident != nil {
			if funcSymbol, ok := f.scopes.SymbolResolve(ident.Name); ok && len(funcSymbol.ReturnTypes()) > 0 {
				// We assume the first return type
				returnType := funcSymbol.ReturnTypes()[0]
				// Verify if the returned type is a struct
				if typeSymbol, ok := f.scopes.SymbolResolve(returnType); ok && typeSymbol.IsStruct() {
					ret = returnType
				}
			}
		}
	case *ast.UnaryExpr: // es. &MyStruct{}
		if rhs.Op == token.AND {
			if compLit, ok := rhs.X.(*ast.CompositeLit); ok {
				if ident := GetIdent(compLit.Type); ident != nil {
					if returnSymbol, ok := f.scopes.SymbolResolve(ident.Name); ok && returnSymbol.IsStruct() {
						ret = returnSymbol.Name()
					}
				}
			}
		}
	case *ast.TypeAssertExpr:
		if ident := GetIdent(rhs.Type); ident != nil {
			if returnSymbol, ok := f.scopes.SymbolResolve(ident.Name); ok && returnSymbol.IsStruct() {
				ret = returnSymbol.Name()
			}
		}
	}
	if len(ret) == 0 {
		return "", false
	}
	return ret, true
	//return "", nil, false
}

// Finalize performs the final resolution and layout of all struct and interface definitions. It resolves references and computes sizes.
func (f *DefinitionTable) Finalize() error {
	// PHASE 1: Linker (Reference Resolution)
	keys := f.structTable.Keys() // Assume it returns []string
	for _, structName := range keys {
		structDef, _ := f.structTable.Get(structName)

		// Note: structDef must expose Fields() as []IStructField
		for _, field := range structDef.Fields() {
			// If the field is not yet linked to a real definition...
			if field.IsPlaceholder() {
				baseName := field.FieldBase()
				// Look for the definition in the tables
				if baseDef, ok := f.structTable.Get(baseName); ok {
					field.BindDefinition(baseDef)
				} else if interfaceDef, ok := f.interfaceTable.Get(baseName); ok {
					// If you support interfaces as fields
					field.BindDefinition(interfaceDef)
				} else {
					//check if embedded type or undefined
					continue
					//return fmt.Errorf("type '%s' undefined in struct '%s'", baseName, structName)
				}
			}
		}
	}

	// PHASE 2: Layout Engine (Size Calculation)
	visiting := make(map[string]bool)
	for _, structName := range keys {
		if err := f.computeLayout(structName, visiting); err != nil {
			return err
		}
	}
	return nil
}

// computeLayout calculates the memory layout of a struct by determining field offsets and total size recursively.
// It checks for recursive value types and handles both primitive and embedded value fields.
func (f *DefinitionTable) computeLayout(structName string, visiting map[string]bool) error {
	structDef, _ := f.structTable.Get(structName)

	// If the struct already has a calculated size, exit
	// (Assume StructDefinition has an IsFinalized method or check size > 0)
	if structDef.IsFinalized() {
		return nil
	}

	if visiting[structName] {
		return fmt.Errorf("recursive value type detected: %s (use a pointer instead)", structName)
	}
	visiting[structName] = true

	currentOffset := 0
	for _, field := range structDef.Fields() {
		field.SetOffset(currentOffset)

		if field.IsPointer() {
			currentOffset += 8 // Pointer (64-bit)
		} else {
			// It's an embedded value. We need to know how large the base type is.

			// Retrieve the linked definition (guaranteed to exist thanks to PHASE 1)
			def := field.Definition()

			// Cast to *StructDefinition (you need to handle the Interface case separately if needed)
			if baseDef, ok := def.(interface {
				TotalSize() int
				IsFinalized() bool
				Name() string
			}); ok {
				// Recursion: calculate the child's layout if not ready
				if !baseDef.IsFinalized() {
					if err := f.computeLayout(baseDef.Name(), visiting); err != nil {
						return err
					}
				}
				currentOffset += baseDef.TotalSize()
			} else {
				// Fallback for primitive or unhandled types?
				// If 'baseDef' is nil or not struct, how large is it?
				// Here you should handle base types (int, string) if they are treated as structFields
				// or assume a default size.
			}
		}
	}

	structDef.SetTotalSize(currentOffset)
	structDef.SetFinalized(true)
	delete(visiting, structName)
	return nil
}

/*
// JsonSchemaEncoder trasforma la definizione gerarchica dei campi in JSON.
type JsonSchemaEncoder struct{}

// Encode parte dalla lista radice di campi (lo schema della Struct principale)
func (e *JsonSchemaEncoder) Encode(fields []IStructField) ([]byte, error) {
	// 1. Costruiamo la mappa radice
	rootMap := make(map[string]interface{})

	// 2. Popoliamo ricorsivamente
	// Nota: Un campo struct è definito dai suoi campi, quindi la "mappa" JSON
	// rappresenta i valori di default o la struttura stessa.
	for _, field := range fields {
		//rootMap[field.Name()] = e.buildValue(field.Value())
		////rootMap[field.FieldName()] = e.buildValue(field.FieldName())
	}

	// 3. Serializziamo
	return json.Marshal(rootMap)
}

// buildValue è la funzione ricorsiva che "non conosce IObjects"
// Gestisce tipi primitivi, mappe, slice e soprattutto strutture annidate ([]IStructField)
func (e *JsonSchemaEncoder) buildValue(val interface{}) interface{} {
	if val == nil {
		return nil
	}

	switch v := val.(type) {

	// --- Caso Chiave: Struct Annidata ---
	// Se il valore è una lista di campi, significa che è una struct annidata.
	// Dobbiamo trasformarla in una mappa JSON.
	case []IStructField:
		subMap := make(map[string]interface{})
		for _, f := range v {
			subMap[f.Name()] = e.buildValue(f.Value()) // <--- Ricorsione
		}
		return subMap

	// --- Container Generici ---
	case map[string]interface{}:
		// Mappa generica: ricorsione sui valori
		out := make(map[string]interface{}, len(v))
		for k, subVal := range v {
			out[k] = e.buildValue(subVal)
		}
		return out

	case []interface{}:
		// Array generico: ricorsione sugli elementi
		out := make([]interface{}, len(v))
		for i, subVal := range v {
			out[i] = e.buildValue(subVal)
		}
		return out

	// --- Primitivi ---
	// Passano direttamente al JSON
	case int, int64, float64, string, bool:
		return v

	default:
		// Fallback per tipi custom o imprevisti (es. Stringer)
		if s, ok := v.(fmt.Stringer); ok {
			return s.String()
		}
		return v
	}
}
*/
