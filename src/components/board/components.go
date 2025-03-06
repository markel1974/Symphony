package board

import "fmt"

// Components manages a collection of IComponent instances using a map for storing and retrieving by their unique IDs.
type Components struct {
	container map[string]IComponent
}

// NewComponents initializes and returns a pointer to a new Components instance with an empty container map.
func NewComponents() *Components {
	return &Components{
		container: make(map[string]IComponent),
	}
}

// Register adds the given IComponent to the container using its unique ID as the key.
func (s *Components) Register(component IComponent) {
	id := component.GetId()
	s.container[id] = component
}

// Dump retrieves a component by its ID, initializes a dumper with its properties, and exports the component's state.
func (m *Components) Dump(id string) error {
	c, ok := m.container[id]
	if !ok {
		return fmt.Errorf("component '%s' not found", id)
	}
	dumper := NewDumper(c.GetProperties())
	return c.Dump(dumper)
}

// Restore attempts to restore the state of a component by its ID using a Dumper instance. Returns an error if not found.
func (m *Components) Restore(id string) error {
	c, ok := m.container[id]
	if !ok {
		return fmt.Errorf("component '%s' not found", id)
	}
	dumper := NewDumper(c.GetProperties())
	return c.Restore(dumper)
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
