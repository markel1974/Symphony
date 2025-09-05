package sdk

import (
	"time"

	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	RegisterPackage(NewTime)
}

// Time represents a structure that manages a collection of modules implementing the IObject interface.
type Time struct {
	container map[string]objects.IObject
}

// NewTime initializes and returns a new instance of Time with predefined constants and functions mapped to the module.
func NewTime(factory objects.IGateKeeper) IPackage {
	t := &Time{}
	constants := map[string]objects.IObject{
		"ANSIC":       factory.NewString(objects.FrameStatic, time.ANSIC),
		"UnixDate":    factory.NewString(objects.FrameStatic, time.UnixDate),
		"RubyDate":    factory.NewString(objects.FrameStatic, time.RubyDate),
		"RFC822":      factory.NewString(objects.FrameStatic, time.RFC822),
		"RFC822Z":     factory.NewString(objects.FrameStatic, time.RFC822Z),
		"RFC850":      factory.NewString(objects.FrameStatic, time.RFC850),
		"RFC1123":     factory.NewString(objects.FrameStatic, time.RFC1123),
		"RFC1123Z":    factory.NewString(objects.FrameStatic, time.RFC1123Z),
		"RFC3339":     factory.NewString(objects.FrameStatic, time.RFC3339),
		"RFC3339Nano": factory.NewString(objects.FrameStatic, time.RFC3339Nano),
		"Kitchen":     factory.NewString(objects.FrameStatic, time.Kitchen),
		"Stamp":       factory.NewString(objects.FrameStatic, time.Stamp),
		"StampMilli":  factory.NewString(objects.FrameStatic, time.StampMilli),
		"StampMicro":  factory.NewString(objects.FrameStatic, time.StampMicro),
		"StampNano":   factory.NewString(objects.FrameStatic, time.StampNano),
		"Nanosecond":  factory.NewInt(objects.FrameStatic, int64(time.Nanosecond)),
		"Microsecond": factory.NewInt(objects.FrameStatic, int64(time.Microsecond)),
		"Millisecond": factory.NewInt(objects.FrameStatic, int64(time.Millisecond)),
		"Second":      factory.NewInt(objects.FrameStatic, int64(time.Second)),
		"Minute":      factory.NewInt(objects.FrameStatic, int64(time.Minute)),
		"Hour":        factory.NewInt(objects.FrameStatic, int64(time.Hour)),
		"January":     factory.NewInt(objects.FrameStatic, int64(time.January)),
		"February":    factory.NewInt(objects.FrameStatic, int64(time.February)),
		"March":       factory.NewInt(objects.FrameStatic, int64(time.March)),
		"April":       factory.NewInt(objects.FrameStatic, int64(time.April)),
		"May":         factory.NewInt(objects.FrameStatic, int64(time.May)),
		"June":        factory.NewInt(objects.FrameStatic, int64(time.June)),
		"July":        factory.NewInt(objects.FrameStatic, int64(time.July)),
		"August":      factory.NewInt(objects.FrameStatic, int64(time.August)),
		"September":   factory.NewInt(objects.FrameStatic, int64(time.September)),
		"October":     factory.NewInt(objects.FrameStatic, int64(time.October)),
		"November":    factory.NewInt(objects.FrameStatic, int64(time.November)),
		"December":    factory.NewInt(objects.FrameStatic, int64(time.December)),
	}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "Sleep", t.sleep),
		factory.NewFuncImport(objects.FrameStatic, "ParseDuration", t.parseDuration),
		factory.NewFuncImport(objects.FrameStatic, "Since", t.since),
		factory.NewFuncImport(objects.FrameStatic, "Until", t.until),
		factory.NewFuncImport(objects.FrameStatic, "DurationHours", t.durationHours),
		factory.NewFuncImport(objects.FrameStatic, "DurationMinutes", t.durationMinutes),
		factory.NewFuncImport(objects.FrameStatic, "DurationNanoseconds", t.durationNanoseconds),
		factory.NewFuncImport(objects.FrameStatic, "DurationSeconds", t.durationSeconds),
		factory.NewFuncImport(objects.FrameStatic, "DurationString", t.durationString),
		factory.NewFuncImport(objects.FrameStatic, "MonthString", t.monthString),
		factory.NewFuncImport(objects.FrameStatic, "Date", t.date),
		factory.NewFuncImport(objects.FrameStatic, "Now", t.now),
		factory.NewFuncImport(objects.FrameStatic, "Parse", t.parse),
		factory.NewFuncImport(objects.FrameStatic, "Unix", t.unix),
		factory.NewFuncImport(objects.FrameStatic, "add", t.add),
		factory.NewFuncImport(objects.FrameStatic, "AddDate", t.addDate),
		factory.NewFuncImport(objects.FrameStatic, "Sub", t.sub),
		factory.NewFuncImport(objects.FrameStatic, "After", t.after),
		factory.NewFuncImport(objects.FrameStatic, "Before", t.before),
		factory.NewFuncImport(objects.FrameStatic, "TimeYear", t.timeYear),
		factory.NewFuncImport(objects.FrameStatic, "TimeMonth", t.timeMonth),
		factory.NewFuncImport(objects.FrameStatic, "TimeDay", t.timeDay),
		factory.NewFuncImport(objects.FrameStatic, "TimeWeekday", t.timeWeekday),
		factory.NewFuncImport(objects.FrameStatic, "TimeHour", t.timeHour),
		factory.NewFuncImport(objects.FrameStatic, "TimeMinute", t.timeMinute),
		factory.NewFuncImport(objects.FrameStatic, "TimeSecond", t.timeSecond),
		factory.NewFuncImport(objects.FrameStatic, "TimeNanosecond", t.timeNanosecond),
		factory.NewFuncImport(objects.FrameStatic, "TimeUnix", t.timeUnix),
		factory.NewFuncImport(objects.FrameStatic, "TimeUnixNano", t.timeUnixNano),
		factory.NewFuncImport(objects.FrameStatic, "TimeFormat", t.timeFormat),
		factory.NewFuncImport(objects.FrameStatic, "TimeLocation", t.timeLocation),
		factory.NewFuncImport(objects.FrameStatic, "TimeString", t.timeString),
		factory.NewFuncImport(objects.FrameStatic, "IsZero", t.isZero),
		factory.NewFuncImport(objects.FrameStatic, "ToLocal", t.toLocal),
		factory.NewFuncImport(objects.FrameStatic, "ToUTC", t.toUTC),
	}
	t.container = BuildContainer(container, constants)
	return t
}

// Name returns the name of the Math module as a string.
func (t *Time) Name() string {
	return "time"
}

// Get retrieves an object associated with the given name from the container. It returns the object and a boolean indicating success.
func (t *Time) Get(name string) (objects.IObject, bool) {
	v, ok := t.container[name]
	return v, ok
}

// sleep pauses the execution for a specified duration provided as an argument in nanoseconds.
// Returns an error if the number of arguments is incorrect or the argument is not an integer.
// On success, returns the undefined value.
func (t *Time) sleep(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	i1, err := gk.ToInt64Arg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	time.Sleep(time.Duration(i1))
	return 1, gk.UndefinedValue(), nil
}

// parseDuration parses a duration string and converts it into an integer representation of nanoseconds as IObject.
// Accepts exactly one argument of type string. Returns an error object if parsing or type conversion fails.
func (t *Time) parseDuration(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	dur, err := time.ParseDuration(s1)
	if err != nil {
		return 1, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewInt(frame, int64(dur)), nil
}

// since calculates the time duration between the current Time instance and the given time argument as an integer.
// Expects exactly one argument of a compatible time type, returns an error if the argument is invalid or missing.
func (t *Time) since(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(time.Since(t1))), nil
}

// until calculates the duration from the current time to a specified time object and returns it as an integer object.
// Returns an error if the argument is missing, invalid, or not a time-compatible object.
func (t *Time) until(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(time.Until(t1))), nil
}

// durationHours calculates the duration in hours from an integer argument representing a duration in nanoseconds.
// It returns a Float object containing the result or an error if the input is invalid.
func (t *Time) durationHours(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	i1, err := gk.ToInt64Arg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewFloat(frame, time.Duration(i1).Hours()), nil
}

// durationMinutes calculates the duration in minutes based on the given integer argument and returns it as a float.
// Returns an error if the number of arguments is incorrect or if the argument type is invalid.
func (t *Time) durationMinutes(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	i1, err := gk.ToInt64Arg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewFloat(frame, time.Duration(i1).Minutes()), nil
}

// durationNanoseconds returns the nanosecond representation of a given duration argument as an IObject, or an error for invalid input.
func (t *Time) durationNanoseconds(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	i1, err := gk.ToInt64Arg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, time.Duration(i1).Nanoseconds()), nil
}

// durationSeconds converts the given integer argument (in nanoseconds) to a float representation of seconds.
func (t *Time) durationSeconds(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	i1, err := gk.ToInt64Arg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewFloat(frame, time.Duration(i1).Seconds()), nil
}

// durationString converts a duration given as an integer to its string representation and returns it as an IObject.
// Returns an error if not exactly one argument is provided or if the argument is not a valid integer.
func (t *Time) durationString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	i1, err := gk.ToInt64Arg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, time.Duration(i1).String()), nil
}

// monthString takes a single integer argument, converts it to a month, and returns its string representation.
// Returns an error if the argument count is incorrect or if the conversion fails.
func (t *Time) monthString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	i1, err := gk.ToInt64Arg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, time.Month(i1).String()), nil
}

// Date creates a new Time object using the specified year, month, day, hour, minute, second, and nanosecond values.
// It requires exactly 7 integer arguments and returns an error if the argument count or types are incorrect.
func (t *Time) date(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 7 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	i1, err := gk.ToInt64Arg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	i3, err := gk.ToInt64Arg(2, args[2])
	if err != nil {
		return 0, nil, err
	}
	i4, err := gk.ToInt64Arg(3, args[3])
	if err != nil {
		return 0, nil, err
	}
	i5, err := gk.ToInt64Arg(4, args[4])
	if err != nil {
		return 0, nil, err
	}
	i6, err := gk.ToInt64Arg(5, args[5])
	if err != nil {
		return 0, nil, err
	}
	i7, err := gk.ToInt64Arg(6, args[6])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewTime(frame, time.Date(int(i1), time.Month(i2), int(i3), int(i4), int(i5), int(i6), int(i7), time.Now().Location())), nil
}

// Now retrieves the current time as a Time object. Returns an error if any arguments are provided.
func (t *Time) now(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 0 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	return 1, gk.NewTime(frame, time.Now()), nil
}

// Parse parses a time string using the given format and returns a new Time object or an error if parsing fails.
func (t *Time) parse(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 2 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	s2, err := gk.ToStringArg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	parsed, err := time.Parse(s1, s2)
	if err != nil {
		return 1, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewTime(frame, parsed), nil
}

// Unix creates a new Time object based on the given Unix timestamp and nanoseconds, or returns an error for invalid arguments.
func (t *Time) unix(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 2 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	i1, err := gk.ToInt64Arg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewTime(frame, time.Unix(i1, i2)), nil
}

// add adds a duration (int64) to a Time object and returns a new Time object or an error if the inputs are invalid.
func (t *Time) add(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 2 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewTime(frame, t1.Add(time.Duration(i2))), nil
}

// Sub calculates the duration between two time arguments and returns it as an Int object or an error if invalid arguments.
func (t *Time) sub(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 2 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	t2, err := gk.ToTimeArg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Sub(t2))), nil
}

// AddDate adjusts the date by adding years, months, and days to the provided time object, returning the result as an IObject.
// It expects four arguments: a time object, and three integers representing years, months, and days respectively.
// Returns an error if the wrong number of arguments is provided or a type conversion fails.
func (t *Time) addDate(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 4 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	i3, err := gk.ToInt64Arg(2, args[2])
	if err != nil {
		return 0, nil, err
	}
	i4, err := gk.ToInt64Arg(3, args[3])
	if err != nil {
		return 0, nil, err
	}
	v := t1.AddDate(int(i2), int(i3), int(i4))
	return 1, gk.NewTime(frame, v), nil
}

// After compares two time values and returns TrueValue if the first is after the second, otherwise returns FalseValue.
func (t *Time) after(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 2 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	t2, err := gk.ToTimeArg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	if t1.After(t2) {
		return 1, gk.TrueValue(), nil
	}
	return 1, gk.FalseValue(), nil
}

// Before determines if the first time argument occurs before the second time argument and returns a boolean result.
func (t *Time) before(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 2 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	t2, err := gk.ToTimeArg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	if t1.Before(t2) {
		return 1, gk.TrueValue(), nil
	}
	return 1, gk.FalseValue(), nil
}

// TimeYear returns the year component of a given time object as an integer. Accepts a single argument of type IObject.
func (t *Time) timeYear(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Year())), nil
}

// TimeMonth extracts the month from a time object and returns it as an integer. It requires exactly one argument.
// Returns an error if the argument count is incorrect or the type conversion fails.
func (t *Time) timeMonth(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Month())), nil
}

// TimeDay extracts and returns the day of the month as an integer from a given time object. It requires exactly one argument.
func (t *Time) timeDay(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Day())), nil
}

// TimeWeekday returns the weekday of a time object as an integer. Returns an error if the argument count is invalid.
func (t *Time) timeWeekday(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Weekday())), nil
}

// TimeHour extracts the hour from the given time object and returns it as an Int. Returns an error if arguments are invalid.
func (t *Time) timeHour(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Hour())), nil
}

// TimeMinute extracts the minute component from a time object and returns it as an integer.
func (t *Time) timeMinute(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Minute())), nil
}

// TimeSecond extracts the second component from a time object passed as an argument and returns it as an Int object.
// Returns an error if the argument count is not 1 or if the conversion to a time object fails.
func (t *Time) timeSecond(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Second())), nil
}

// timeNanosecond returns the nanosecond component of the given time object as an integer.
// Expects a single argument of a time-compatible object. Returns an error for invalid arguments or conversion failures.
func (t *Time) timeNanosecond(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Nanosecond())), nil
}

// timeUnix converts a provided time object into its Unix timestamp and returns it as an integer object.
func (t *Time) timeUnix(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, t1.Unix()), nil
}

// timeUnixNano returns the Unix time in nanoseconds as an IObject for the given time argument. An error is returned for invalid input.
func (t *Time) timeUnixNano(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, t1.UnixNano()), nil
}

// timeFormat formats a time object using the provided format and returns the formatted string as an IObject.
func (t *Time) timeFormat(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 2 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	s2, err := gk.ToStringArg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	s := t1.Format(s2)
	return 1, gk.NewString(frame, s), nil
}

// isZero checks if the provided time argument is zero and returns TrueValue if it is, otherwise returns FalseValue.
func (t *Time) isZero(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	if t1.IsZero() {
		return 1, gk.TrueValue(), nil
	}
	return 1, gk.FalseValue(), nil
}

// toLocal converts the given IObject argument to a local time zone Time object or returns an error if conversion fails.
func (t *Time) toLocal(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewTime(frame, t1.Local()), nil
}

// toUTC converts the provided IObject time argument to UTC and returns a new IObject representing the UTC time.
func (t *Time) toUTC(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewTime(frame, t1.UTC()), nil
}

// timeLocation returns the location (timezone) from the given time object as a string. Takes exactly one argument.
func (t *Time) timeLocation(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, t1.Location().String()), nil
}

// timeString converts a time instance to its string representation. It requires exactly one argument of type IObject.
func (t *Time) timeString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 1 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	t1, err := gk.ToTimeArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, t1.String()), nil
}
