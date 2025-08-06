package stats

// bToMb converts bytes to megabytes by dividing the input value in bytes by 1024 twice.
func bToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

/*
// getCPUSample reads CPU statistics from /proc/stat and returns the idle ticks and total ticks as uint64 values.
func getCPUSample() (uint64, uint64) {
	var idle uint64 = 0  //(rand.Intn(max - min) + min)
	var total uint64 = 0 //(rand.Intn(max - min) + min)

	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return idle, total
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err == nil {
					total += val // tally up all the numbers to get total ticks
					if i == 4 {  // idle is the 5th field in the cpu line
						idle = val
					}
				}
			}
			return idle, total
		}
	}
	return total, total
}
*/
