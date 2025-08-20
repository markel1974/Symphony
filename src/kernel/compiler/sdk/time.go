package sdk

import (
	"time"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Time represents a structure that manages a collection of modules implementing the IObject interface.
type Time struct {
	factory *objects.Factory
	*Package
}

// NewTime initializes and returns a new instance of Time with predefined constants and functions mapped to the module.
func NewTime(factory *objects.Factory) *Time {
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
		factory.NewFuncPackage(objects.FuncPackageDef, "Sleep", t.Sleep),
		factory.NewFuncPackage(objects.FuncPackageDef, "ParseDuration", t.ParseDuration),
		factory.NewFuncPackage(objects.FuncPackageDef, "Since", t.Since),
		factory.NewFuncPackage(objects.FuncPackageDef, "Until", t.Until),
		factory.NewFuncPackage(objects.FuncPackageDef, "DurationHours", t.DurationHours),
		factory.NewFuncPackage(objects.FuncPackageDef, "DurationMinutes", t.DurationMinutes),
		factory.NewFuncPackage(objects.FuncPackageDef, "DurationNanoseconds", t.DurationNanoseconds),
		factory.NewFuncPackage(objects.FuncPackageDef, "DurationSeconds", t.DurationSeconds),
		factory.NewFuncPackage(objects.FuncPackageDef, "DurationString", t.DurationString),
		factory.NewFuncPackage(objects.FuncPackageDef, "MonthString", t.MonthString),
		factory.NewFuncPackage(objects.FuncPackageDef, "Date", t.Date),
		factory.NewFuncPackage(objects.FuncPackageDef, "Now", t.Now),
		factory.NewFuncPackage(objects.FuncPackageDef, "Parse", t.Parse),
		factory.NewFuncPackage(objects.FuncPackageDef, "Unix", t.Unix),
		factory.NewFuncPackage(objects.FuncPackageDef, "Add", t.Add),
		factory.NewFuncPackage(objects.FuncPackageDef, "AddDate", t.AddDate),
		factory.NewFuncPackage(objects.FuncPackageDef, "Sub", t.Sub),
		factory.NewFuncPackage(objects.FuncPackageDef, "After", t.After),
		factory.NewFuncPackage(objects.FuncPackageDef, "Before", t.Before),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeYear", t.TimeYear),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeMonth", t.TimeMonth),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeDay", t.TimeDay),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeWeekday", t.TimeWeekday),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeHour", t.TimeHour),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeMinute", t.TimeMinute),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeSecond", t.TimeSecond),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeNanosecond", t.TimeNanosecond),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeUnix", t.TimeUnix),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeUnixNano", t.TimeUnixNano),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeFormat", t.TimeFormat),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeLocation", t.TimeLocation),
		factory.NewFuncPackage(objects.FuncPackageDef, "TimeString", t.TimeString),
		factory.NewFuncPackage(objects.FuncPackageDef, "is_zero", t.IsZero),
		factory.NewFuncPackage(objects.FuncPackageDef, "to_local", t.ToLocal),
		factory.NewFuncPackage(objects.FuncPackageDef, "to_utc", t.ToUTC),
	}
	t.Package = NewPackage("time", container, constants)
	return t
}

// Sleep pauses the execution for a specified duration provided as an argument in nanoseconds.
// Returns an error if the number of arguments is incorrect or the argument is not an integer.
// On success, returns the undefined value.
func (t *Time) Sleep(args ...objects.IObject) (objects.IObject, error) {
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

// ParseDuration parses a duration string and converts it into an integer representation of nanoseconds as IObject.
// Accepts exactly one argument of type string. Returns an error object if parsing or type conversion fails.
func (t *Time) ParseDuration(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := t.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	dur, err := time.ParseDuration(s1)
	if err != nil {
		return t.factory.NewObjectError(objects.FrameReturnValue, err), nil
	}
	return t.factory.NewInt(objects.FrameReturnValue, int64(dur)), nil
}

// Since calculates the time duration between the current Time instance and the given time argument as an integer.
// Expects exactly one argument of a compatible time type, returns an error if the argument is invalid or missing.
func (t *Time) Since(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, int64(time.Since(t1))), nil
}

// Until calculates the duration from the current time to a specified time object and returns it as an integer object.
// Returns an error if the argument is missing, invalid, or not a time-compatible object.
func (t *Time) Until(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, int64(time.Until(t1))), nil
}

// DurationHours calculates the duration in hours from an integer argument representing a duration in nanoseconds.
// It returns a Float object containing the result or an error if the input is invalid.
func (t *Time) DurationHours(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewFloat(objects.FrameReturnValue, time.Duration(i1).Hours()), nil
}

// DurationMinutes calculates the duration in minutes based on the given integer argument and returns it as a float.
// Returns an error if the number of arguments is incorrect or if the argument type is invalid.
func (t *Time) DurationMinutes(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewFloat(objects.FrameReturnValue, time.Duration(i1).Minutes()), nil
}

// DurationNanoseconds returns the nanosecond representation of a given duration argument as an IObject, or an error for invalid input.
func (t *Time) DurationNanoseconds(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, time.Duration(i1).Nanoseconds()), nil
}

// DurationSeconds converts the given integer argument (in nanoseconds) to a float representation of seconds.
func (t *Time) DurationSeconds(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewFloat(objects.FrameReturnValue, time.Duration(i1).Seconds()), nil
}

// DurationString converts a duration given as an integer to its string representation and returns it as an IObject.
// Returns an error if not exactly one argument is provided or if the argument is not a valid integer.
func (t *Time) DurationString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewString(objects.FrameReturnValue, time.Duration(i1).String())
}

// MonthString takes a single integer argument, converts it to a month, and returns its string representation.
// Returns an error if the argument count is incorrect or if the conversion fails.
func (t *Time) MonthString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := t.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewStringNoSize(objects.FrameReturnValue, time.Month(i1).String()), nil
}

// Date creates a new Time object using the specified year, month, day, hour, minute, second, and nanosecond values.
// It requires exactly 7 integer arguments and returns an error if the argument count or types are incorrect.
func (t *Time) Date(args ...objects.IObject) (objects.IObject, error) {
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
	return t.factory.NewTime(objects.FrameReturnValue, time.Date(int(i1), time.Month(i2), int(i3), int(i4), int(i5), int(i6), int(i7), time.Now().Location())), nil
}

// Now retrieves the current time as a Time object. Returns an error if any arguments are provided.
func (t *Time) Now(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 0 {
		return nil, objects.ErrWrongNumArguments
	}
	return t.factory.NewTime(objects.FrameReturnValue, time.Now()), nil
}

// Parse parses a time string using the given format and returns a new Time object or an error if parsing fails.
func (t *Time) Parse(args ...objects.IObject) (objects.IObject, error) {
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
		return t.factory.NewObjectError(objects.FrameReturnValue, err), nil
	}
	return t.factory.NewTime(objects.FrameReturnValue, parsed), nil
}

// Unix creates a new Time object based on the given Unix timestamp and nanoseconds, or returns an error for invalid arguments.
func (t *Time) Unix(args ...objects.IObject) (objects.IObject, error) {
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
	return t.factory.NewTime(objects.FrameReturnValue, time.Unix(i1, i2)), nil
}

// Add adds a duration (int64) to a Time object and returns a new Time object or an error if the inputs are invalid.
func (t *Time) Add(args ...objects.IObject) (objects.IObject, error) {
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
	return t.factory.NewTime(objects.FrameReturnValue, t1.Add(time.Duration(i2))), nil
}

// Sub calculates the duration between two time arguments and returns it as an Int object or an error if invalid arguments.
func (t *Time) Sub(args ...objects.IObject) (objects.IObject, error) {
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
	return t.factory.NewInt(objects.FrameReturnValue, int64(t1.Sub(t2))), nil
}

// AddDate adjusts the date by adding years, months, and days to the provided time object, returning the result as an IObject.
// It expects four arguments: a time object, and three integers representing years, months, and days respectively.
// Returns an error if the wrong number of arguments is provided or a type conversion fails.
func (t *Time) AddDate(args ...objects.IObject) (objects.IObject, error) {
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
	return t.factory.NewTime(objects.FrameReturnValue, v), nil
}

// After compares two time values and returns TrueValue if the first is after the second, otherwise returns FalseValue.
func (t *Time) After(args ...objects.IObject) (objects.IObject, error) {
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
func (t *Time) Before(args ...objects.IObject) (objects.IObject, error) {
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
func (t *Time) TimeYear(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, int64(t1.Year())), nil
}

// TimeMonth extracts the month from a time object and returns it as an integer. It requires exactly one argument.
// Returns an error if the argument count is incorrect or the type conversion fails.
func (t *Time) TimeMonth(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, int64(t1.Month())), nil
}

// TimeDay extracts and returns the day of the month as an integer from a given time object. It requires exactly one argument.
func (t *Time) TimeDay(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, int64(t1.Day())), nil
}

// TimeWeekday returns the weekday of a time object as an integer. Returns an error if the argument count is invalid.
func (t *Time) TimeWeekday(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, int64(t1.Weekday())), nil
}

// TimeHour extracts the hour from the given time object and returns it as an Int. Returns an error if arguments are invalid.
func (t *Time) TimeHour(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, int64(t1.Hour())), nil
}

// TimeMinute extracts the minute component from a time object and returns it as an integer.
func (t *Time) TimeMinute(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, int64(t1.Minute())), nil
}

// TimeSecond extracts the second component from a time object passed as an argument and returns it as an Int object.
// Returns an error if the argument count is not 1 or if the conversion to a time object fails.
func (t *Time) TimeSecond(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, int64(t1.Second())), nil
}

// TimeNanosecond returns the nanosecond component of the given time object as an integer.
// Expects a single argument of a time-compatible object. Returns an error for invalid arguments or conversion failures.
func (t *Time) TimeNanosecond(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, int64(t1.Nanosecond())), nil
}

// TimeUnix converts a provided time object into its Unix timestamp and returns it as an integer object.
func (t *Time) TimeUnix(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, t1.Unix()), nil
}

// TimeUnixNano returns the Unix time in nanoseconds as an IObject for the given time argument. An error is returned for invalid input.
func (t *Time) TimeUnixNano(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewInt(objects.FrameReturnValue, t1.UnixNano()), nil
}

// TimeFormat formats a time object using the provided format and returns the formatted string as an IObject.
func (t *Time) TimeFormat(args ...objects.IObject) (objects.IObject, error) {
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
	return t.factory.NewString(objects.FrameReturnValue, s)
}

// IsZero checks if the provided time argument is zero and returns TrueValue if it is, otherwise returns FalseValue.
func (t *Time) IsZero(args ...objects.IObject) (objects.IObject, error) {
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

// ToLocal converts the given IObject argument to a local time zone Time object or returns an error if conversion fails.
func (t *Time) ToLocal(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewTime(objects.FrameReturnValue, t1.Local()), nil
}

// ToUTC converts the provided IObject time argument to UTC and returns a new IObject representing the UTC time.
func (t *Time) ToUTC(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewTime(objects.FrameReturnValue, t1.UTC()), nil
}

// TimeLocation returns the location (timezone) from the given time object as a string. Takes exactly one argument.
func (t *Time) TimeLocation(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewString(objects.FrameReturnValue, t1.Location().String())
}

// TimeString converts a time instance to its string representation. It requires exactly one argument of type IObject.
func (t *Time) TimeString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := t.factory.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return t.factory.NewString(objects.FrameReturnValue, t1.String())
}
