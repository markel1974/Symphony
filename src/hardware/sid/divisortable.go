package mos6581

type DivisorTableData struct {
	divisor int32
	toOut   int32
}

type DivisorTable struct {
	divisorTableData []*DivisorTableData
}

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

func (dt *DivisorTable) GetDivisor(x int) int32 {
	return dt.divisorTableData[x].divisor
}

func (dt *DivisorTable) GetOut(x int) int32 {
	return dt.divisorTableData[x].toOut
}
