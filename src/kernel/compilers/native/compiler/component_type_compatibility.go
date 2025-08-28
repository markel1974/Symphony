package compiler

import (
	"go/ast"
	"go/token"
	"reflect"
)

// TypeCompatibility gestisce la verifica della compatibilità tra tipi,
// come la verifica dell'implementazione di interfacce da parte degli struct.
type TypeCompatibility struct {
	fileSet        *token.FileSet
	structTable    *StructTable
	interfaceTable *InterfaceTable
	functionTable  *FunctionTable
	// Mappa per memorizzare le implementazioni: structName -> []interfaceName
	implementations map[string][]string
}

// NewTypeCompatibility crea una nuova istanza di TypeCompatibility.
func NewTypeCompatibility(structTable *StructTable, interfaceTable *InterfaceTable, functionTable *FunctionTable) *TypeCompatibility {
	return &TypeCompatibility{
		structTable:     structTable,
		interfaceTable:  interfaceTable,
		functionTable:   functionTable,
		implementations: make(map[string][]string),
	}
}

// Setup non fa nulla in questo componente, ma è richiesto dall'interfaccia IComponent.
func (tc *TypeCompatibility) Setup(fileSet *token.FileSet, _ func(node ast.Node) error) error {
	tc.fileSet = fileSet
	return nil
}

// Prepare è dove eseguiamo l'analisi. Deve essere chiamato dopo che
// tutte le definizioni di struct, interfacce e metodi sono state caricate.
func (tc *TypeCompatibility) Prepare() error {
	//fmt.Println("Running interface implementation check...")
	// Itera su ogni struct definito
	for structName, _ := range tc.structTable.container {
		// Itera su ogni interfaccia definita
		for interfaceName, interfaceDesc := range tc.interfaceTable.container {
			implements, err := tc.checkStructImplementsInterface(structName, interfaceDesc)
			if err != nil {
				// In una versione reale, potremmo voler raccogliere tutti gli errori
				return err
			}
			if implements {
				// Memorizza la relazione
				tc.implementations[structName] = append(tc.implementations[structName], interfaceName)
				//fmt.Printf("=> SUCCESS: Struct '%s' implements interface '%s'\n", structName, interfaceName)
			}
		}
	}
	// Aggiungiamo le implementazioni trovate alla StructTable per un accesso centralizzato
	tc.structTable.SetImplementations(tc.implementations)
	return nil
}

// Compile non fa nulla in questo componente.
func (tc *TypeCompatibility) Compile() error {
	return nil
}

// checkStructImplementsInterface verifica se un dato struct implementa un'interfaccia.
func (tc *TypeCompatibility) checkStructImplementsInterface(structName string, interfaceDesc *InterfaceDescription) (bool, error) {
	for _, requiredMethod := range interfaceDesc.Methods {
		// Il nome "mangled" del metodo dello struct è "StructName.MethodName"
		mangledMethodName := GetMangledName(structName, requiredMethod.Name)

		// Cerca la descrizione della funzione (metodo) nella functionTable
		var structMethod *FunctionDescription
		for i := 0; i < tc.functionTable.Len(); i++ {
			fd, _ := tc.functionTable.Get(i)
			if fd.Name == mangledMethodName {
				structMethod = fd
				break
			}
		}

		if structMethod == nil {
			// Metodo richiesto non trovato sullo struct
			return false, nil
		}

		// Confronta le firme dei metodi
		// NOTA: il ricevitore conta come primo parametro per il metodo dello struct
		numStructParams := len(structMethod.InputParams)
		if len(requiredMethod.InputParams) != numStructParams {
			return false, nil
		}
		if !reflect.DeepEqual(requiredMethod.InputParams, structMethod.InputParams) {
			return false, nil
		}
		if !reflect.DeepEqual(requiredMethod.ReturnTypes, structMethod.ReturnTypes) {
			return false, nil
		}
	}
	// Se siamo arrivati qui, tutti i metodi richiesti sono stati trovati e le firme corrispondono
	return true, nil
}
