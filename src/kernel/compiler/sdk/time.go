package sdk

import (
	"time"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Time represents a structure that manages a collection of modules implementing the IObject interface.
type Time struct {
	factory *objects.GateKeeper
	*Package
}

// NewTime initializes and returns a new instance of Time with predefined constants and functions mapped to the module.
func NewTime(factory *objects.GateKeeper) *Time {
	t := &Time{factory: factory}
	constants := map[string]objects.IObject{
		"ANSIC":       factory.NewStringNoSize(objects.FrameStatic, time.ANSIC),
		"UnixDate":    factory.NewStringNoSize(objects.FrameStatic, time.UnixDate),
		"RubyDate":    factory.NewStringNoSize(objects.FrameStatic, time.RubyDate),
		"RFC822":      factory.NewStringNoSize(objects.FrameStatic, time.RFC822),
		"RFC822Z":     factory.NewStringNoSize(objects.FrameStatic, time.RFC822Z),
		"RFC850":      factory.NewStringNoSize(objects.FrameStatic, time.RFC850),
		"RFC1123":     factory.NewStringNoSize(objects.FrameStatic, time.RFC1123),
		"RFC1123Z":    factory.NewStringNoSize(objects.FrameStatic, time.RFC1123Z),
		"RFC3339":     factory.NewStringNoSize(objects.FrameStatic, time.RFC3339),
		"RFC3339Nano": factory.NewStringNoSize(objects.FrameStatic, time.RFC3339Nano),
		"Kitchen":     factory.NewStringNoSize(objects.FrameStatic, time.Kitchen),
		"Stamp":       factory.NewStringNoSize(objects.FrameStatic, time.Stamp),
		"StampMilli":  factory.NewStringNoSize(objects.FrameStatic, time.StampMilli),
		"StampMicro":  factory.NewStringNoSize(objects.FrameStatic, time.StampMicro),
		"StampNano":   factory.NewStringNoSize(objects.FrameStatic, time.StampNano),
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
	container := []*objects.FuncPackage{
		factory.NewFuncPackage(objects.FuncPackageDef, "Sleep", t.sleep),
		factory.NewFuncPackage(objects.FuncPackageDef, "ParseDuration", t.parseDuration),
		factory.NewFuncPackage(objects.FuncPackageDef, "Since", t.since),
		factory.NewFuncPackage(objects.FuncPackageDef, "Until", t.until),
		factory.NewFuncPackage(objects.FuncPackageDef, "DurationHours", t.durationHours),
		factory.NewFuncPackage(objects.FuncPackageDef, "DurationMinutes", t.durationMinutes),
		factory.NewFuncPackage(objects.FuncPackageDef, "DurationNanoseconds", t.durationNanoseconds),
		factory.NewFuncPackage(objects.FuncPackageDef, "DurationSeconds", t.durationSeconds),
		factory.NewFuncPackage(objects.FuncPackageDef, "DurationString", t.durationString),
		factory.NewFuncPackage(objects.FuncPackageDef, "MonthString", t.monthString),
		factory.NewFuncPackage(objects.FuncPackageDef, "Date", t.date),
		factory.NewFuncPackage(objects.FuncPackageDef, "Now", t.now),
		factory.NewFuncPackage(objects.FuncPackageDef, "Parse", t.parse),
		factory.NewFuncPackage(objects.FuncPackageDef, "Unix", t.unix),
		factory.NewFuncPackage(objects.FuncPackageDef, "Add", t.add),
		factory.NewFuncPackage(objects.FuncPackageDef, "AddDate", t.addDate),
		factory.NewFuncPackage(objects.FuncPackageDef, "Sub", t.sub),
		factory.NewFuncPackage(objects.FuncPackageDef, "After", t.after),
		factory.NewFuncPackage(objects.FuncPackageDef, "Before", t.before),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeYear", t.timeYear),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeMonth", t.timeMonth),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeDay", t.timeDay),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeWeekday", t.timeWeekday),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeHour", t.timeHour),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeMinute", t.timeMinute),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeSecond", t.timeSecond),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeNanosecond", t.timeNanosecond),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeUnix", t.timeUnix),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeUnixNano", t.timeUnixNano),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeFormat", t.timeFormat),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeLocation", t.timeLocation),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeString", t.timeString),
		factory.NewFuncPackage(objects.FuncPackageDef, "IsZero", t.isZero),
		factory.NewFuncPackage(objects.FuncPackageDef, "ToLocal", t.toLocal),
		factory.NewFuncPackage(objects.FuncPackageDef, "ToUTC", t.toUTC),
	}
	t.Package = NewPackage("time", container, constants)
	return t
}

// sleep pauses the execution for a specified duration provided as an argument in nanoseconds.
// Returns an error if the number of arguments is incorrect or the argument is not an integer.
// On success, returns the undefined value.
func (t *Time) sleep(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Duration(i1))
	return t.factory.UndefinedValue(), nil
}

// parseDuration parses a duration string and converts it into an integer representation of nanoseconds as IObject.
// Accepts exactly one argument of type string. Returns an error object if parsing or type conversion fails.
func (t *Time) parseDuration(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := t.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	dur, err := time.ParseDuration(s1)
	if err != nil {
		return t.factory.NewError(frame, err.Error()), nil
	}
	return t.factory.NewInt(frame, int64(dur)), nil
}

// since calculates the time duration between the current Time instance and the given time argument as an integer.
// Expects exactly one argument of a compatible time type, returns an error if the argument is invalid or missing.
func (t *Time) since(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, int64(time.Since(t1))), nil
}

// until calculates the duration from the current time to a specified time object and returns it as an integer object.
// Returns an error if the argument is missing, invalid, or not a time-compatible object.
func (t *Time) until(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, int64(time.Until(t1))), nil
}

// durationHours calculates the duration in hours from an integer argument representing a duration in nanoseconds.
// It returns a Float object containing the result or an error if the input is invalid.
func (t *Time) durationHours(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewFloat(frame, time.Duration(i1).Hours()), nil
}

// durationMinutes calculates the duration in minutes based on the given integer argument and returns it as a float.
// Returns an error if the number of arguments is incorrect or if the argument type is invalid.
func (t *Time) durationMinutes(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewFloat(frame, time.Duration(i1).Minutes()), nil
}

// durationNanoseconds returns the nanosecond representation of a given duration argument as an IObject, or an error for invalid input.
func (t *Time) durationNanoseconds(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, time.Duration(i1).Nanoseconds()), nil
}

// durationSeconds converts the given integer argument (in nanoseconds) to a float representation of seconds.
func (t *Time) durationSeconds(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewFloat(frame, time.Duration(i1).Seconds()), nil
}

// durationString converts a duration given as an integer to its string representation and returns it as an IObject.
// Returns an error if not exactly one argument is provided or if the argument is not a valid integer.
func (t *Time) durationString(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewString(frame, time.Duration(i1).String())
}

// monthString takes a single integer argument, converts it to a month, and returns its string representation.
// Returns an error if the argument count is incorrect or if the conversion fails.
func (t *Time) monthString(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewStringNoSize(frame, time.Month(i1).String()), nil
}

// Date creates a new Time object using the specified year, month, day, hour, minute, second, and nanosecond values.
// It requires exactly 7 integer arguments and returns an error if the argument count or types are incorrect.
func (t *Time) date(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 7 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := t.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	i3, err := t.factory.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	i4, err := t.factory.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	i5, err := t.factory.ToInt64Arg(4, args[4])
	if err != nil {
		return nil, err
	}
	i6, err := t.factory.ToInt64Arg(5, args[5])
	if err != nil {
		return nil, err
	}
	i7, err := t.factory.ToInt64Arg(6, args[6])
	if err != nil {
		return nil, err
	}
	return t.factory.NewTime(frame, time.Date(int(i1), time.Month(i2), int(i3), int(i4), int(i5), int(i6), int(i7), time.Now().Location())), nil
}

// Now retrieves the current time as a Time object. Returns an error if any arguments are provided.
func (t *Time) now(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 0 {
		return nil, objects.ErrWrongNumArguments
	}
	return t.factory.NewTime(frame, time.Now()), nil
}

// Parse parses a time string using the given format and returns a new Time object or an error if parsing fails.
func (t *Time) parse(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := t.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := t.factory.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	parsed, err := time.Parse(s1, s2)
	if err != nil {
		return t.factory.NewError(frame, err.Error()), nil
	}
	return t.factory.NewTime(frame, parsed), nil
}

// Unix creates a new Time object based on the given Unix timestamp and nanoseconds, or returns an error for invalid arguments.
func (t *Time) unix(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := t.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	return t.factory.NewTime(frame, time.Unix(i1, i2)), nil
}

// Add adds a duration (int64) to a Time object and returns a new Time object or an error if the inputs are invalid.
func (t *Time) add(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := t.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	return t.factory.NewTime(frame, t1.Add(time.Duration(i2))), nil
}

// Sub calculates the duration between two time arguments and returns it as an Int object or an error if invalid arguments.
func (t *Time) sub(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	t2, err := t.factory.ToTimeArg(1, args[1])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, int64(t1.Sub(t2))), nil
}

// AddDate adjusts the date by adding years, months, and days to the provided time object, returning the result as an IObject.
// It expects four arguments: a time object, and three integers representing years, months, and days respectively.
// Returns an error if the wrong number of arguments is provided or a type conversion fails.
func (t *Time) addDate(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 4 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := t.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	i3, err := t.factory.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	i4, err := t.factory.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	v := t1.AddDate(int(i2), int(i3), int(i4))
	return t.factory.NewTime(frame, v), nil
}

// After compares two time values and returns TrueValue if the first is after the second, otherwise returns FalseValue.
func (t *Time) after(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	t2, err := t.factory.ToTimeArg(1, args[1])
	if err != nil {
		return nil, err
	}
	if t1.After(t2) {
		return t.factory.TrueValue(), nil
	}
	return t.factory.FalseValue(), nil
}

// Before determines if the first time argument occurs before the second time argument and returns a boolean result.
func (t *Time) before(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	t2, err := t.factory.ToTimeArg(1, args[1])
	if err != nil {
		return nil, err
	}
	if t1.Before(t2) {
		return t.factory.TrueValue(), nil
	}
	return t.factory.FalseValue(), nil
}

// TimeYear returns the year component of a given time object as an integer. Accepts a single argument of type IObject.
func (t *Time) timeYear(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, int64(t1.Year())), nil
}

// TimeMonth extracts the month from a time object and returns it as an integer. It requires exactly one argument.
// Returns an error if the argument count is incorrect or the type conversion fails.
func (t *Time) timeMonth(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, int64(t1.Month())), nil
}

// TimeDay extracts and returns the day of the month as an integer from a given time object. It requires exactly one argument.
func (t *Time) timeDay(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, int64(t1.Day())), nil
}

// TimeWeekday returns the weekday of a time object as an integer. Returns an error if the argument count is invalid.
func (t *Time) timeWeekday(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, int64(t1.Weekday())), nil
}

// TimeHour extracts the hour from the given time object and returns it as an Int. Returns an error if arguments are invalid.
func (t *Time) timeHour(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, int64(t1.Hour())), nil
}

// TimeMinute extracts the minute component from a time object and returns it as an integer.
func (t *Time) timeMinute(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, int64(t1.Minute())), nil
}

// TimeSecond extracts the second component from a time object passed as an argument and returns it as an Int object.
// Returns an error if the argument count is not 1 or if the conversion to a time object fails.
func (t *Time) timeSecond(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, int64(t1.Second())), nil
}

// timeNanosecond returns the nanosecond component of the given time object as an integer.
// Expects a single argument of a time-compatible object. Returns an error for invalid arguments or conversion failures.
func (t *Time) timeNanosecond(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, int64(t1.Nanosecond())), nil
}

// timeUnix converts a provided time object into its Unix timestamp and returns it as an integer object.
func (t *Time) timeUnix(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, t1.Unix()), nil
}

// timeUnixNano returns the Unix time in nanoseconds as an IObject for the given time argument. An error is returned for invalid input.
func (t *Time) timeUnixNano(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(frame, t1.UnixNano()), nil
}

// timeFormat formats a time object using the provided format and returns the formatted string as an IObject.
func (t *Time) timeFormat(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := t.factory.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	s := t1.Format(s2)
	return t.factory.NewString(frame, s)
}

// isZero checks if the provided time argument is zero and returns TrueValue if it is, otherwise returns FalseValue.
func (t *Time) isZero(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if t1.IsZero() {
		return t.factory.TrueValue(), nil
	}
	return t.factory.FalseValue(), nil
}

// toLocal converts the given IObject argument to a local time zone Time object or returns an error if conversion fails.
func (t *Time) toLocal(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewTime(frame, t1.Local()), nil
}

// toUTC converts the provided IObject time argument to UTC and returns a new IObject representing the UTC time.
func (t *Time) toUTC(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewTime(frame, t1.UTC()), nil
}

// timeLocation returns the location (timezone) from the given time object as a string. Takes exactly one argument.
func (t *Time) timeLocation(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewString(frame, t1.Location().String())
}

// timeString converts a time instance to its string representation. It requires exactly one argument of type IObject.
func (t *Time) timeString(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewString(frame, t1.String())
}
