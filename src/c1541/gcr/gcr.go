package gcr

type GCR struct {
	data      []uint8
	errorInfo []uint8
}

func NewGCR() *GCR {
	d := &GCR{
		data:      make([]uint8, DiskSize),
		errorInfo: make([]uint8, NumSectors),
	}
	for x := range d.data {
		d.data[x] = 0x55
	}
	for x := range d.errorInfo {
		d.errorInfo[x] = 1
	}
	return d
}

func (gcr *GCR) GetData() []uint8 {
	return gcr.data
}

func (gcr *GCR) GetErrorInfo() []uint8 {
	return gcr.errorInfo
}
