package mos6581

// DivisorTableData represents a structure that holds a divisor and its corresponding output value for calculations.
type DivisorTableData struct {
	divisor int32
	toOut   int32
}

// DivisorTable is a structure that holds a collection of divisor and output data entries.
type DivisorTable struct {
	divisorTableData []*DivisorTableData
}

// NewDivisorTable creates and initializes a DivisorTable for managing divisors based on the given rasters and fragFreq values.
func NewDivisorTable(rasters int, fragFreq int) *DivisorTable {
	d := int32(rasters * fragFreq)
	dt := &DivisorTable{
		divisorTableData: make([]*DivisorTableData, SampleFreq+1),
	}
	for x := range dt.divisorTableData {
		dtd := &DivisorTableData{}
		dtd.divisor = int32(x)
		dtd.toOut = 0
		for dtd.divisor >= 0 {
			dtd.divisor -= d
			dtd.toOut++
		}
		dt.divisorTableData[x] = dtd
	}
	return dt
}

// GetDivisor retrieves the divisor value associated with the given index from the DivisorTable.
func (dt *DivisorTable) GetDivisor(x int) int32 {
	return dt.divisorTableData[x].divisor
}

// GetOut returns the `toOut` value from the DivisorTableData at index `x`.
func (dt *DivisorTable) GetOut(x int) int32 {
	return dt.divisorTableData[x].toOut
}
