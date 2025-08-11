package format

// fmtFlags represents formatting flags used in formatting operations.
// It encapsulates settings like width, precision, and various flags for alignment, sign, and formatting style.
type fmtFlags struct {
	widPresent       bool
	precisionPresent bool
	minus            bool
	plus             bool
	sharp            bool
	space            bool
	zero             bool
	plusV            bool
	sharpV           bool
	inDetail         bool
	needNewline      bool
	needColon        bool
}
