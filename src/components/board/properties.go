package board

import (
	"fmt"
	"reflect"
)

type RunFn func(cmd string, args []string) (map[string]interface{}, error)

// PropertyInfo represents metadata associated with a property, including its type, description, and access methods.
type PropertyInfo struct {
	id          string
	kind        string
	description string
	readOnly    bool
	getType     reflect.Type
	getValue    reflect.Value
	setType     reflect.Type
	setValue    reflect.Value
	set0Kind    reflect.Kind
}

// CreatePropertyInfo creates a new PropertyInfo instance with the specified type, description, and read-only configuration.
// It returns an error if the specified type is not supported.
func CreatePropertyInfo(id string, kind interface{}, desc string, ro bool, get interface{}, set interface{}) *PropertyInfo {
	p := &PropertyInfo{id: id, kind: fmt.Sprintf("%T", kind), description: desc, readOnly: ro}
	p.getType = reflect.TypeOf(get)
	p.getValue = reflect.ValueOf(get)
	p.setType = reflect.TypeOf(set)
	p.setValue = reflect.ValueOf(set)
	if get == nil || p.getType.Kind() != reflect.Func {
		panic(fmt.Errorf("get isn't a function"))
	}
	if p.getType.NumIn() != 0 || p.getType.NumOut() != 1 {
		panic("wrong get signature must be func get() <ret>")
	}
	if set == nil || p.setType.Kind() != reflect.Func {
		panic(fmt.Errorf("set function is nil"))
	}
	if p.setType.NumIn() != 1 || p.setType.NumOut() != 0 {
		panic("wrong get signature must be func set(v <arg>)")
	}
	if p.setType.NumOut() != 0 {
		panic("wrong set signature")
	}
	p.set0Kind = p.setType.In(0).Kind()
	return p
}

type Properties struct {
	properties map[string]*PropertyInfo
	run        RunFn
}

func NewProperties(run RunFn) *Properties {
	return &Properties{
		run:        run,
		properties: make(map[string]*PropertyInfo),
	}
}

func (p *Properties) Add(info *PropertyInfo) {
	if info == nil {
		panic(fmt.Errorf("property info is nil"))
	}
	p.properties[info.id] = info
}

func (p *Properties) Get(id string) *PropertyInfo {
	return p.properties[id]
}

func (p *Properties) Run(cmd string, values []string) (map[string]interface{}, error) {
	return p.run(cmd, values)
}

func (p *Properties) GetProperty(id string) (interface{}, error) {
	prop, ok := p.properties[id]
	if !ok || prop == nil {
		return nil, fmt.Errorf("property '%s' not found", id)
	}
	results := prop.getValue.Call([]reflect.Value{})
	if len(results) == 0 {
		return nil, fmt.Errorf("no results returned")
	}
	result := results[0].Interface()
	return result, nil
}

func (p *Properties) SetProperty(id string, arg interface{}) error {
	prop, ok := p.properties[id]
	if !ok || prop == nil {
		return fmt.Errorf("property '%s' not found", id)
	}
	argType := reflect.TypeOf(arg)
	if argType.Kind() != prop.set0Kind {
		return fmt.Errorf("wrong input signature")
	}
	argValue := reflect.ValueOf(arg)
	args := []reflect.Value{argValue}
	prop.setValue.Call(args)
	return nil
}

func (p *Properties) Dump() (map[string]interface{}, error) {
	//TODO IMPLEMENT
	return nil, fmt.Errorf("unimplemented")
}

func (p *Properties) Restore(d map[string]interface{}) error {
	//TODO IMPLEMENT
	return fmt.Errorf("unimplemented")
}

/*
func addOne(x int) int {
	return x + 1
}

func multiplyByTwo(x int) int {
	return x * 2
}

func invalidFunction(s string) string {
	return s + "!"
}

func callFunctionIfValid(target interface{}, arg interface{}, out interface{}) (interface{}, error) {
	funcType := reflect.TypeOf(target)
	funcValue := reflect.ValueOf(target)
	if funcType.Kind() != reflect.Func {
		return nil, fmt.Errorf("target isn't a function")
	}
	inputArg := reflect.TypeOf(arg).Kind()
	outputArg := reflect.TypeOf(out).Kind()
	// Verifica la firma della funzione dei dati di input
	if funcType.NumIn() != 1 || funcType.In(0).Kind() != inputArg {
		return nil, fmt.Errorf("wrong input signature")
	}
	// Verifica la firma della funzione dei dati di output
	if funcType.NumOut() != 1 || funcType.Out(0).Kind() != outputArg {
		return nil, fmt.Errorf("wrong output signature")
	}
	args := []reflect.Value{reflect.ValueOf(arg)}
	results := funcValue.Call(args)
	if len(results) == 0 {
		return nil, fmt.Errorf("no results returned")
	}
	result := results[0].Interface()
	return result, nil
}

func Test() {
	output := 0
	callFunctionIfValid(addOne, 5, output)
	callFunctionIfValid(multiplyByTwo, 10, output)
	callFunctionIfValid(invalidFunction, 12, output) // Fallisce
	callFunctionIfValid(123, 12, output)             //Fallisce

	//components := board.NewComponents()
	//s := mos6581.NewSID("test", "")
	//components.Register(s)
	//components.Dump(s.GetId())
	//os.Exit(1)
}


*/
