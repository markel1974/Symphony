package mos6581

// Lead
// data is an array for managing lead data used in audio rendering.
// pos specifies the current position within the lead array.
// smooth determines the smoothing level applied to the lead buffer.
// hiWater defines the high-water threshold for lead buffer management.
// loWater sets the low-water threshold for lead buffer management.
type Lead struct {
	data    []int
	pos     int
	smooth  int
	hiWater int
	loWater int
}

func NewLead(fragFreq int) *Lead {
	fragInterval := 1000 / fragFreq // in milliseconds
	maxLeadAvg := fragFreq
	return &Lead{
		data:    make([]int, maxLeadAvg),
		smooth:  LatencyAvg / fragInterval,
		hiWater: LatencyMax / fragInterval,
		loWater: LatencyMin / fragInterval,
	}
}

func (l *Lead) Reset() {
	for x := range l.data {
		l.data[x] = 0
	}
	l.pos = 0
}

func (l *Lead) Average(leadInFrags int) (int, bool) {
	l.data[l.pos] = leadInFrags
	l.pos++
	if l.pos == l.smooth {
		l.pos = 0
	}
	// Compute the average lead in frags.
	avgLead := 0
	for i := 0; i < l.smooth; i++ {
		avgLead += l.data[i]
	}
	avgLead /= l.smooth
	//fmt.Printf("lead = %d, avg = %d\n", leadInFrags, avgLead)
	if avgLead > l.hiWater {
		//too far ahead of the audio skip a frag.
		for i := 0; i < l.smooth; i++ {
			l.data[i]--
		}
		//fmt.Printf("Skipping a frag...\n")
		return 0, false
	}
	return avgLead, true
}

func (l *Lead) GetLoWater() int {
	return l.loWater
}

func (l *Lead) Update() {
	for i := 0; i < l.smooth; i++ {
		l.data[i]++
	}
}
