package sdk

import (
	"runtime"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	register(NewRuntime)
}

type Runtime struct {
	*bytecode.Package
}

func NewRuntime(gk objects.IGateKeeper) bytecode.IPackage {
	const (
		defMemStats     = "MemStats"
		defReadMemStats = "ReadMemStats"
	)
	h := &Runtime{Package: bytecode.NewPackage("runtime")}
	h.Add(defMemStats, gk.NewStruct(objects.FrameStatic, defMemStats, nil))
	h.Add(defReadMemStats, gk.NewFuncImport(objects.FrameStatic, defReadMemStats, 1, h.readMemStats))
	return h
}

// encodeToString converts the given IObject into a hexadecimal string representation using the provided IGateKeeper.
// It requires exactly one argument. Returns an error if argument count is invalid or conversion fails.
func (f *Runtime) readMemStats(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	args1, ok := args[0].(*objects.ObjectPointer)
	if !ok {
		return 0, gk.UndefinedValue(), objects.ErrInvalidArgumentsNumber
	}
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	ret := make(map[string]interface{})
	ret["Alloc"] = m.Alloc
	ret["TotalAlloc"] = m.TotalAlloc
	ret["Sys"] = m.Sys
	ret["Lookups"] = m.Lookups
	ret["Mallocs"] = m.Mallocs
	ret["Frees"] = m.Frees
	ret["HeapAlloc"] = m.HeapAlloc
	ret["HeapSys"] = m.HeapSys
	ret["HeapIdle"] = m.HeapIdle
	ret["HeapInuse"] = m.HeapInuse
	ret["HeapReleased"] = m.HeapReleased
	ret["HeapObjects"] = m.HeapObjects
	ret["StackInuse"] = m.StackInuse
	ret["StackSys"] = m.StackSys
	ret["MSpanInuse"] = m.MSpanInuse
	ret["MSpanSys"] = m.MSpanSys
	ret["MCacheInuse"] = m.MCacheInuse
	ret["MCacheSys"] = m.MCacheSys
	ret["BuckHashSys"] = m.BuckHashSys
	ret["GCSys"] = m.GCSys
	ret["OtherSys"] = m.OtherSys
	ret["NextGC"] = m.NextGC
	ret["LastGC"] = m.LastGC
	ret["PauseTotalNs"] = m.PauseTotalNs
	ret["NumGC"] = m.NumGC
	ret["NumForcedGC"] = m.NumForcedGC
	ret["EnableGC"] = m.EnableGC
	ret["DebugGC"] = m.DebugGC

	res := gk.StructFromMap(frame, "MemStats", ret)
	value := args1.Value()
	if err := (*value).AssignValue(res); err != nil {
		return 0, gk.UndefinedValue(), err
	}
	return 0, gk.UndefinedValue(), nil
}
