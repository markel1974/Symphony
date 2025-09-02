package objects

const (
	maxDepth     = 256
	maxBytesLen  = 268435455
	MaxStringLen = 268435455
	maxMapLen    = 100000
	maxArrayLen  = 100000
	maxStructLen = 100000
)

// GateKeeper is a type responsible for creating and managing IObject instances, including primitive and complex types.
// It provides pre-instantiated objects for `true`, `false`, and `undefined` values for efficient reuse.
// The GateKeeper may also include object pooling for specific types to optimize memory usage and performance.
type GateKeeper struct {
	IGateAllocator
	IGateConverter
	IGateAdapter
	IGateCall
}

func (f *GateKeeper) FromMap(frame int, v map[string]interface{}) map[string]IObject {
	//TODO implement me
	panic("implement me")
}

const (
	FrameStatic = -1
)

// NewGateKeeper initializes a new GateKeeper instance and sets up default bool and undefined values.
func NewGateKeeper(maxAllocations int64) IGateKeeper {
	f := &GateKeeper{}
	f.IGateAllocator = NewGateAllocator(f, maxAllocations)
	f.IGateConverter = NewGateConverter(f)
	f.IGateAdapter = NewGateAdapter(f)
	f.IGateCall = NewGateCall(f)
	return f
}

func (f *GateKeeper) Reset() {
	f.IGateAllocator.Reset()
}
