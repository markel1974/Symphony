package board

import (
	"fmt"
	"strings"
)

// Components manages a collection of IComponent instances, facilitating their registration and interaction.
// It provides methods to manipulate properties, execute commands, and serialize/restore component states.
type Components struct {
	container map[string]IComponent
}

// NewComponents creates and returns a new instance of Components with an initialized container map.
func NewComponents() *Components {
	return &Components{
		container: make(map[string]IComponent),
	}
}

// Register adds an IComponent to the Components container using its unique identifier as the key.
func (m *Components) Register(component IComponent) {
	id := component.GetId()
	m.container[id] = component
}

// RunCommand executes a command for a specific component by ID, using the provided arguments, and returns the results or an error.
func (m *Components) RunCommand(id string, cmd string, args string) (map[string]interface{}, error) {
	c, ok := m.container[id]
	if !ok {
		return nil, fmt.Errorf("component '%s' not found", id)
	}
	values := strings.Split(args, " ")
	d, err := c.GetProperties().Run(cmd, values)
	return d, err
}

// SetProperty sets the specified property of a component to the provided value. Returns an error if the component or property is not found.
func (m *Components) SetProperty(id string, prop string, val interface{}) error {
	c, ok := m.container[id]
	if !ok {
		return fmt.Errorf("component '%s' not found", id)
	}
	err := c.GetProperties().SetProperty(prop, val)
	return err
}

// GetProperty retrieves the value of a specified property from a given component by its id.
// Returns the property value and an error if the component or property is not found.
func (m *Components) GetProperty(id string, prop string) (interface{}, error) {
	c, ok := m.container[id]
	if !ok {
		return fmt.Errorf("component '%s' not found", id), nil
	}
	v, err := c.GetProperties().GetProperty(prop)
	return v, err
}

// Dump retrieves and returns a snapshot of all properties of the specified component by its ID. Returns an error if the component is not found.
func (m *Components) Dump(id string) (map[string]interface{}, error) {
	c, ok := m.container[id]
	if !ok {
		return nil, fmt.Errorf("component '%s' not found", id)
	}
	d, err := c.GetProperties().Dump()
	return d, err
}

// Restore updates the properties of a component identified by the given id using the provided data map. Returns an error if the component is not found or the restore operation fails.
func (m *Components) Restore(id string, d map[string]interface{}) error {
	c, ok := m.container[id]
	if !ok {
		return fmt.Errorf("component '%s' not found", id)
	}
	err := c.GetProperties().Restore(d)
	return err
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
