package formats

type fmtFlags struct {
	widPresent       bool
	precisionPresent bool
	minus            bool
	plus             bool
	sharp            bool
	space            bool
	zero             bool

	// For the formats %+v %#v, we set the plusV/sharpV flags
	// and clear the plus/sharp flags since %+v and %#v are in effect
	// different, flagless formats set at the top level.
	plusV  bool
	sharpV bool

	// error-related flags.
	inDetail    bool
	needNewline bool
	needColon   bool
}
