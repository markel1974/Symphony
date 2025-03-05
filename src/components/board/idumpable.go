package board

import (
	"log"
	"strings"
)

// IDumpable represents an interface for dumping and restoring an object's state using a map of string to interface{}.
type IDumpable interface {
	Dump(d map[string]interface{}) error
	Restore(d map[string]interface{}) error
}

// CreateIdFromPath joins the elements of the given path slice into a single string, separated by dots.
func CreateIdFromPath(path []string) string {
	return strings.Join(path, ".")
}

// CreatePathFromKey splits a given key string by '.' and returns a slice of strings representing the path segments.
func CreatePathFromKey(key string) []string {
	return strings.Split(key, ".")
}

// dumpSet sets a value in the map `d` with a key generated from the `path` array using CreateIdFromPath.
func dumpSet(d map[string]interface{}, path []string, val interface{}) {
	key := CreateIdFromPath(path)
	d[key] = val
}

// DumpSetUint8 sets a uint8 value into the map `d` at the specified `path` by converting it to a float64 internally.
func DumpSetUint8(d map[string]interface{}, path []string, val uint8) {
	DumpSetNumber(d, path, float64(val))
}

// DumpSetUint16 sets a uint16 value in a nested map structure at the specified path by converting it to a float64.
func DumpSetUint16(d map[string]interface{}, path []string, val uint16) {
	DumpSetNumber(d, path, float64(val))
}

// DumpSetUint sets a uint value in a nested map using a specified path by converting it to a float64 and delegating to DumpSetNumber.
func DumpSetUint(d map[string]interface{}, path []string, val uint) {
	DumpSetNumber(d, path, float64(val))
}

// DumpSetInt sets an integer value in a nested map structure at a specific path, converting it to float64 internally.
func DumpSetInt(d map[string]interface{}, path []string, val int) {
	DumpSetNumber(d, path, float64(val))
}

// DumpSetNumber sets a float64 value in the map at the specified path by converting the path into a unique key.
func DumpSetNumber(d map[string]interface{}, path []string, val float64) {
	dumpSet(d, path, val)
}

// DumpSetBool sets a boolean value in the given map using a unique key derived from the provided path.
func DumpSetBool(d map[string]interface{}, path []string, val bool) {
	dumpSet(d, path, val)
}

// DumpSetString sets a string value in the provided map at a location defined by the path.
func DumpSetString(d map[string]interface{}, path []string, val string) {
	dumpSet(d, path, val)
}

// DumpSetByteArray sets a byte array value in a hierarchical map-like structure at the specified path.
func DumpSetByteArray(d map[string]interface{}, path []string, val []byte) {
	c := make([]uint8, len(val))
	copy(c, val)
	dumpSet(d, path, c)
}

// dumpGet retrieves a value from a map using a key derived from the provided path and checks if the value exists and is non-nil.
func dumpGet(d map[string]interface{}, path []string) (interface{}, bool) {
	if d == nil {
		return 0, false
	}
	val, found := d[CreateIdFromPath(path)]
	if !found || val == nil {
		return nil, false
	}
	return val, true
}

// DumpGetFloat64 retrieves a float64 value from a nested map using the specified path and stores it in the provided pointer.
// Returns true if successful; otherwise false if the key does not exist or the value is not a float64.
// Logs an error message if the value at the key is of an unexpected type.
func DumpGetFloat64(d map[string]interface{}, path []string, val *float64) bool {
	data, ok := dumpGet(d, path)
	if !ok {
		return false
	}
	t, ok := data.(float64)
	if !ok {
		log.Printf("DumpGetFloat64: invalid type for key '%v': expected float64, got %T", path, t)
		return false
	}
	*val = t
	return true
}

// DumpGetInt retrieves an integer value from a nested map using the specified path, storing it in the provided pointer.
// Returns true if successful or false if the key is not found or the value is not an integer.
// Logs an error message if the type assertion fails.
func DumpGetInt(d map[string]interface{}, path []string, val *int) bool {
	data, ok := dumpGet(d, path)
	if !ok {
		return false
	}
	t, ok := data.(int)
	if !ok {
		log.Printf("DumpGetInt: invalid type for key '%v': expected int, got %T", path, t)
		return false
	}
	*val = t
	return true
}

// DumpGetUint retrieves a uint value from a nested map based on the provided path and assigns it to the provided pointer.
// Returns true if the value is successfully retrieved and matches the uint type; otherwise, returns false.
func DumpGetUint(d map[string]interface{}, path []string, val *uint) bool {
	data, ok := dumpGet(d, path)
	if !ok {
		return false
	}
	t, ok := data.(uint)
	if !ok {
		log.Printf("DumpGetUint8: invalid type for key '%v': expected uint, got %T", path, t)
		return false
	}
	*val = t
	return true
}

// DumpGetUint8 retrieves a uint8 value from the map at the specified path and assigns it to the provided pointer.
// Returns true if the retrieval and type assertion succeed, otherwise logs an error and returns false.
func DumpGetUint8(d map[string]interface{}, path []string, val *uint8) bool {
	data, ok := dumpGet(d, path)
	if !ok {
		return false
	}
	t, ok := data.(uint8)
	if !ok {
		log.Printf("DumpGetUint8: invalid type for key '%v': expected uint8, got %T", path, t)
		return false
	}
	*val = t
	return true
}

// DumpGetUint16 retrieves a uint16 value from a map based on a given path and stores it in the provided pointer.
// Returns true if the value is successfully retrieved and is of type uint16, otherwise returns false.
func DumpGetUint16(d map[string]interface{}, path []string, val *uint16) bool {
	data, ok := dumpGet(d, path)
	if !ok {
		return false
	}
	t, ok := data.(uint16)
	if !ok {
		log.Printf("DumpGetUint16: invalid type for key '%v': expected uint16, got %T", path, t)
		return false
	}
	*val = t
	return true
}

// DumpGetBool retrieves a boolean value from a nested map based on the provided path and stores it in the provided pointer.
// If the path does not exist or the value is not a boolean, it logs an error and returns false.
// Returns true if the value is successfully retrieved and assigned.
func DumpGetBool(d map[string]interface{}, path []string, val *bool) bool {
	t, ok := dumpGet(d, path)
	if !ok {
		return false
	}
	r, ok := t.(bool)
	if !ok {
		log.Printf("DumpGetBool: invalid type for key '%v': expected bool, got %T", path, t)
		return false
	}
	*val = r
	return true
}

// DumpGetString retrieves a string value from a nested map using a path of keys and assigns it to the provided pointer.
// Returns false if the path is invalid or the type is not a string.
func DumpGetString(d map[string]interface{}, path []string, val *string) bool {
	t, ok := dumpGet(d, path)
	if !ok {
		return false
	}
	r, ok := t.(string)
	if !ok {
		log.Printf("DumpGetString: invalid type for key '%v': expected string, got %T", path, t)
		return false
	}
	*val = r
	return true
}

// DumpGetByteArray retrieves a byte array from a map using a specified path and validates its length if count is provided.
func DumpGetByteArray(d map[string]interface{}, path []string, res *[]byte, count int) bool {
	val, ok := dumpGet(d, path)
	if !ok {
		return false
	}
	r, ok := val.([]uint8)
	if !ok {
		log.Printf("DumpGetByteArray: invalid type for key '%v': expected []uint8, got %T", path, val)
		return false
	}
	if count > -1 {
		if len(r) != count {
			return false
		}
	}
	if len(r) == 0 {
		return false
	}
	*res = make([]byte, len(r))
	copy(*res, r)
	return true
}

/*
func (b *Board) RestoreComponent(componentID string, state map[string]interface{}) error {
    parts, err := splitPath(componentID) // Es: ["cia1", "timerA", "cr"]
    if err != nil {
        return err
    }
    if len(parts) == 0 {
        return fmt.Errorf("invalid component ID: %s", componentID)
    }

    componentName := parts[0] // "cia1"
    component, ok := b.components[componentName]
    if !ok {
        return fmt.Errorf("unknown component: %s", componentName)
    }

    if len(parts) == 1 {
        // Ripristina l'intero componente.
        if dumpable, ok := component.(IDumpable); ok {
            return dumpable.Restore(state)
        }
        return fmt.Errorf("component %s does not support restoring", componentID)
    }

    // Passa il resto del percorso al componente.  *NON* fare uno switch qui.
    return component.RestoreProperty(parts[1:], state) // Ipotetico metodo
}

func (c *CIA) RestoreProperty(path []string, state map[string]interface{}) error {
    if len(path) == 0 {
        return fmt.Errorf("empty path") // Errore: percorso vuoto.
    }

    propertyName := path[0] // "timerA"

    switch propertyName {
    case "timerA":
        return c.timerA.RestoreProperty(path[1:], state) // Chiama RestoreProperty di TimerA
    case "timerB":
        return c.timerB.RestoreProperty(path[1:], state) // Chiama RestoreProperty di TimerB
    // ... altri casi (per le porte, ecc.) ...
     case "pra": // Esempio: accesso diretto a un registro (se è *veramente* necessario)
         if len(path) > 1 {
             return fmt.Errorf("invalid path for CIA register: %v", path)
         }
         var value uint8 // Usa il tipo *corretto*
         if !board.DumpGetUint8(state, "pra", &value) { // Supponendo che la mappa contenga "cia1.pra"
             return fmt.Errorf("invalid or missing value for CIA register 'pra'")
         }
         c.WriteRegister(0, value) // Scrivi nel registro (usa il metodo del CIA!)
         return nil

    default:
        return fmt.Errorf("unknown property: %s", propertyName)
    }
}
*/
