package stdlib

import (
	"time"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// _timeModule is a map providing time-related constants and functions as IObject-compatible values and methods.
var _timeModule = map[string]objects.IObject{
	"ANSIC":               objects.NewStringNoSize(time.ANSIC),
	"UnixDate":            objects.NewStringNoSize(time.UnixDate),
	"RubyDate":            objects.NewStringNoSize(time.RubyDate),
	"RFC822":              objects.NewStringNoSize(time.RFC822),
	"RFC822Z":             objects.NewStringNoSize(time.RFC822Z),
	"RFC850":              objects.NewStringNoSize(time.RFC850),
	"RFC1123":             objects.NewStringNoSize(time.RFC1123),
	"RFC1123Z":            objects.NewStringNoSize(time.RFC1123Z),
	"RFC3339":             objects.NewStringNoSize(time.RFC3339),
	"RFC3339Nano":         objects.NewStringNoSize(time.RFC3339Nano),
	"Kitchen":             objects.NewStringNoSize(time.Kitchen),
	"Stamp":               objects.NewStringNoSize(time.Stamp),
	"StampMilli":          objects.NewStringNoSize(time.StampMilli),
	"StampMicro":          objects.NewStringNoSize(time.StampMicro),
	"StampNano":           objects.NewStringNoSize(time.StampNano),
	"Nanosecond":          objects.NewInt(int64(time.Nanosecond)),
	"Microsecond":         objects.NewInt(int64(time.Microsecond)),
	"Millisecond":         objects.NewInt(int64(time.Millisecond)),
	"Second":              objects.NewInt(int64(time.Second)),
	"Minute":              objects.NewInt(int64(time.Minute)),
	"Hour":                objects.NewInt(int64(time.Hour)),
	"January":             objects.NewInt(int64(time.January)),
	"February":            objects.NewInt(int64(time.February)),
	"March":               objects.NewInt(int64(time.March)),
	"April":               objects.NewInt(int64(time.April)),
	"May":                 objects.NewInt(int64(time.May)),
	"June":                objects.NewInt(int64(time.June)),
	"July":                objects.NewInt(int64(time.July)),
	"August":              objects.NewInt(int64(time.August)),
	"September":           objects.NewInt(int64(time.September)),
	"October":             objects.NewInt(int64(time.October)),
	"November":            objects.NewInt(int64(time.November)),
	"December":            objects.NewInt(int64(time.December)),
	"Sleep":               objects.NewFunctionModule(objects.FunctionModuleDef, "Sleep", timeSleep),                             // sleep(int)
	"ParseDuration":       objects.NewFunctionModule(objects.FunctionModuleDef, "ParseDuration", timeParseDuration),             // parse_duration(str) => int
	"Since":               objects.NewFunctionModule(objects.FunctionModuleDef, "Since", timeSince),                             // since(time) => int
	"Until":               objects.NewFunctionModule(objects.FunctionModuleDef, "Until", timeUntil),                             // until(time) => int
	"DurationHours":       objects.NewFunctionModule(objects.FunctionModuleDef, "DurationHours", timeDurationHours),             // duration_hours(int) => float
	"DurationMinutes":     objects.NewFunctionModule(objects.FunctionModuleDef, "DurationMinutes", timeDurationMinutes),         // duration_minutes(int) => float
	"DurationNanoseconds": objects.NewFunctionModule(objects.FunctionModuleDef, "DurationNanoseconds", timeDurationNanoseconds), // duration_nanoseconds(int) => int
	"DurationSeconds":     objects.NewFunctionModule(objects.FunctionModuleDef, "DurationSeconds", timeDurationSeconds),         // duration_seconds(int) => float
	"DurationString":      objects.NewFunctionModule(objects.FunctionModuleDef, "DurationString", timeDurationString),           // duration_string(int) => string
	"MonthString":         objects.NewFunctionModule(objects.FunctionModuleDef, "MonthString", timeMonthString),                 // month_string(int) => string
	"Date":                objects.NewFunctionModule(objects.FunctionModuleDef, "Date", timeDate),                               // date(year, month, day, hour, min, sec, nsec) => time
	"Now":                 objects.NewFunctionModule(objects.FunctionModuleDef, "Now", timeNow),                                 // now() => time
	"Parse":               objects.NewFunctionModule(objects.FunctionModuleDef, "Parse", timeParse),                             // parse(format, str) => time
	"Unix":                objects.NewFunctionModule(objects.FunctionModuleDef, "Unix", timeUnix),                               // unix(sec, nsec) => time
	"Add":                 objects.NewFunctionModule(objects.FunctionModuleDef, "Add", timeAdd),                                 // add(time, int) => time
	"AddDate":             objects.NewFunctionModule(objects.FunctionModuleDef, "AddDate", timeAddDate),                         // add_date(time, years, months, days) => time
	"Sub":                 objects.NewFunctionModule(objects.FunctionModuleDef, "Sub", timeSub),                                 // sub(t time, u time) => int
	"After":               objects.NewFunctionModule(objects.FunctionModuleDef, "After", timeAfter),                             // after(t time, u time) => bool
	"Before":              objects.NewFunctionModule(objects.FunctionModuleDef, "Before", timeBefore),                           // before(t time, u time) => bool
	"TimeYear":            objects.NewFunctionModule(objects.FunctionModuleDef, "TimeYear", timeTimeYear),                       // time_year(time) => int
	"TimeMonth":           objects.NewFunctionModule(objects.FunctionModuleDef, "TimeMonth", timeTimeMonth),                     // time_month(time) => int
	"TimeDay":             objects.NewFunctionModule(objects.FunctionModuleDef, "TimeDay", timeTimeDay),                         // time_day(time) => int
	"TimeWeekday":         objects.NewFunctionModule(objects.FunctionModuleDef, "TimeWeekday", timeTimeWeekday),                 // time_weekday(time) => int
	"TimeHour":            objects.NewFunctionModule(objects.FunctionModuleDef, "TimeHour", timeTimeHour),                       // time_hour(time) => int
	"TimeMinute":          objects.NewFunctionModule(objects.FunctionModuleDef, "TimeMinute", timeTimeMinute),                   // time_minute(time) => int
	"TimeSecond":          objects.NewFunctionModule(objects.FunctionModuleDef, "TimeSecond", timeTimeSecond),                   // time_second(time) => int
	"TimeNanosecond":      objects.NewFunctionModule(objects.FunctionModuleDef, "TimeNanosecond", timeTimeNanosecond),           // time_nanosecond(time) => int
	"TimeUnix":            objects.NewFunctionModule(objects.FunctionModuleDef, "TimeUnix", timeTimeUnix),                       // time_unix(time) => int
	"TimeUnixNano":        objects.NewFunctionModule(objects.FunctionModuleDef, "TimeUnixNano", timeTimeUnixNano),               // time_unix_nano(time) => int
	"TimeFormat":          objects.NewFunctionModule(objects.FunctionModuleDef, "TimeFormat", timeTimeFormat),                   // time_format(time, format) => string
	"TimeLocation":        objects.NewFunctionModule(objects.FunctionModuleDef, "TimeLocation", timeTimeLocation),               // time_location(time) => string
	"TimeString":          objects.NewFunctionModule(objects.FunctionModuleDef, "TimeString", timeTimeString),                   // time_string(time) => string
	"IsZero":              objects.NewFunctionModule(objects.FunctionModuleDef, "is_zero", timeIsZero),                          // is_zero(time) => bool
	"ToLocal":             objects.NewFunctionModule(objects.FunctionModuleDef, "to_local", timeToLocal),                        // to_local(time) => time
	"ToUTC":               objects.NewFunctionModule(objects.FunctionModuleDef, "to_utc", timeToUTC),                            // to_utc(time) => time
}

// timeSleep pauses program execution for the specified number of milliseconds, passed as a single argument.
// Returns an error if the argument count is not 1 or if the conversion to int64 fails.
// Returns `objects.UndefinedValue` upon successful completion.
func timeSleep(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Duration(i1))
	return objects.UndefinedValue, nil
}

// timeParseDuration parses a duration string and returns it as an Int object; errors if input is invalid or arguments mismatch.
func timeParseDuration(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	dur, err := time.ParseDuration(s1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return objects.NewInt(int64(dur)), nil
}

// timeSince calculates the time duration since the provided time argument and returns it as an integer object.
// Accepts a single argument that must be convertible to a time.Time object, otherwise returns an error.
// Returns an error if the number of arguments is not exactly one.
func timeSince(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(time.Since(t1))), nil
}

// timeUntil calculates the duration in nanoseconds from the current time until the provided future time argument.
// Returns a new integer object representing the duration or an error if the argument is invalid.
func timeUntil(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(time.Until(t1))), nil
}

// timeDurationHours converts an integer argument in nanoseconds to a floating-point representation in hours.
// Returns an error if the argument count is not 1 or the argument is not a valid integer.
func timeDurationHours(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewFloat(time.Duration(i1).Hours()), nil
}

// timeDurationMinutes computes the duration in minutes from the provided argument, which must be convertible to an int64.
// Returns a floating-point IObject representing the minutes or an error on invalid input or incorrect argument count.
func timeDurationMinutes(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewFloat(time.Duration(i1).Minutes()), nil
}

// timeDurationNanoseconds converts a single integer argument to a time.Duration and returns its value in nanoseconds.
// Returns an error if the number of arguments is not 1 or if the argument is not a valid integer.
func timeDurationNanoseconds(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(time.Duration(i1).Nanoseconds()), nil
}

// timeDurationSeconds converts an integer argument representing nanoseconds into a float object representing seconds.
// It requires exactly one argument; otherwise, it returns an error.
// If the argument is not convertible to an int64, an error is returned.
func timeDurationSeconds(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewFloat(time.Duration(i1).Seconds()), nil
}

// timeDurationString converts an integer argument to a time.Duration and returns its string representation.
// Returns an error if the argument count is not 1 or if the conversion fails.
func timeDurationString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewString(time.Duration(i1).String())
}

// timeMonthString converts a single integer argument to its corresponding month name as a string, returning an error if invalid.
// Expects exactly one argument of type objects.IObject that can be converted to int64, or returns ErrWrongNumArguments if not.
func timeMonthString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewStringNoSize(time.Month(i1).String()), nil
}

// timeDate creates a new date and time object from seven integer arguments representing year, month, day, hour, minute, second, and nanosecond.
// Returns an IObject representing the constructed time or an error if the arguments are invalid or insufficient.
func timeDate(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 7 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	i3, err := objects.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	i4, err := objects.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	i5, err := objects.ToInt64Arg(4, args[4])
	if err != nil {
		return nil, err
	}
	i6, err := objects.ToInt64Arg(5, args[5])
	if err != nil {
		return nil, err
	}
	i7, err := objects.ToInt64Arg(6, args[6])
	if err != nil {
		return nil, err
	}
	return objects.NewTime(time.Date(int(i1), time.Month(i2), int(i3), int(i4), int(i5), int(i6), int(i7), time.Now().Location())), nil
}

// timeNow returns the current time as a new Time object.
// It expects no arguments and returns an error if any are provided.
func timeNow(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 0 {
		return nil, objects.ErrWrongNumArguments
	}
	return objects.NewTime(time.Now()), nil
}

// timeParse parses a time string according to a provided layout and returns a wrapped time object or an error.
func timeParse(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	parsed, err := time.Parse(s1, s2)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return objects.NewTime(parsed), nil
}

// timeUnix creates a new Time object using Unix seconds and nanoseconds provided as int64 arguments.
// Returns an error if the number of arguments is not exactly 2 or if conversion to int64 fails.
func timeUnix(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	return objects.NewTime(time.Unix(i1, i2)), nil
}

// timeAdd adds a duration (as int64) to a time object and returns the resulting time. Returns an error for invalid arguments.
func timeAdd(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	return objects.NewTime(t1.Add(time.Duration(i2))), nil
}

// timeSub subtracts two time values passed as arguments and returns the result as an Int object or an error.
// Expects exactly two arguments of IObject type.
// Returns ErrWrongNumArguments if the number of arguments is incorrect.
// Converts arguments to time.Time and calculates the duration difference.
func timeSub(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	t2, err := objects.ToTimeArg(1, args[1])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Sub(t2))), nil
}

// timeAddDate adjusts a given time by adding years, months, and days using the provided integer arguments.
// Accepts exactly 4 arguments: a time object and three integers for years, months, and days, respectively.
// Returns the adjusted time object or an error if arguments are incorrect or incompatible.
func timeAddDate(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 4 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	i3, err := objects.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	i4, err := objects.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	v := t1.AddDate(int(i2), int(i3), int(i4))
	return objects.NewTime(v), nil
}

// timeAfter checks if the first time argument is after the second time argument and returns a boolean result or an error.
func timeAfter(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	t2, err := objects.ToTimeArg(1, args[1])
	if err != nil {
		return nil, err
	}
	if t1.After(t2) {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// timeBefore checks if the first time argument occurs before the second time argument.
// Returns TrueValue if true, otherwise returns FalseValue.
// Accepts exactly two arguments; returns ErrWrongNumArguments if the argument count is incorrect.
func timeBefore(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	t2, err := objects.ToTimeArg(1, args[1])
	if err != nil {
		return nil, err
	}
	if t1.Before(t2) {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// timeTimeYear returns the year extracted from a time.Time object provided as the single argument.
// Returns an error if the argument count is not 1 or if conversion to time.Time fails.
func timeTimeYear(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Year())), nil
}

// timeTimeMonth returns the month of a given time object as an integer. It expects exactly one time argument.
// Returns an error if the number of arguments is not one or if the argument cannot be converted to a time object.
func timeTimeMonth(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Month())), nil
}

// timeTimeDay returns the day of the month from a time object passed as an argument.
// Expects exactly one argument of a compatible time type.
// Returns an error if the argument count or type is incorrect.
func timeTimeDay(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Day())), nil
}

// timeTimeWeekday returns the weekday of a given time object as an integer (0 for Sunday, 1 for Monday, ..., 6 for Saturday).
// Accepts one argument of type objects.IObject representing the time.
// Returns an objects.IObject representing the weekday or an error if the argument is invalid.
// Errors if the number of arguments is not exactly one.
func timeTimeWeekday(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Weekday())), nil
}

// timeTimeHour extracts the hour (0-23) from a time IObject argument and returns it as a new integer IObject.
func timeTimeHour(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Hour())), nil
}

// timeTimeMinute extracts the minute component from a given time object and returns it as an integer.
// Accepts exactly one argument of a time-compatible object. Returns an error for invalid arguments or types.
func timeTimeMinute(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Minute())), nil
}

// timeTimeSecond returns the second component (0-59) of the provided time object as an integer.
// The function expects exactly one argument of a time-compatible type and returns an error otherwise.
func timeTimeSecond(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Second())), nil
}

// timeTimeNanosecond returns the nanosecond component of the given time.IObject as an integer.
// Expects exactly one argument of a time-compatible IObject.
// Returns an error if the argument count is incorrect or conversion to time fails.
func timeTimeNanosecond(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Nanosecond())), nil
}

// timeTimeUnix converts a given IObject to a Unix timestamp (seconds since epoch) and returns it as an Int object.
// Returns an error if the input argument count is not exactly one or if the argument cannot be converted to a time.
func timeTimeUnix(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(t1.Unix()), nil
}

// timeTimeUnixNano returns the Unix timestamp in nanoseconds for a given time object. Accepts exactly one argument.
func timeTimeUnixNano(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(t1.UnixNano()), nil
}

// timeTimeFormat formats a time.Time object using a specified format string and returns the formatted string as IObject.
// Accepts two arguments: a time object and a format string.
// Returns an error if the argument count is incorrect or if type conversion fails.
func timeTimeFormat(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	s := t1.Format(s2)
	return objects.NewString(s)
}

// timeIsZero checks if the given time object is a zero value and returns a corresponding boolean object or an error.
func timeIsZero(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if t1.IsZero() {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// timeToLocal converts a UTC time object to its local time equivalent and returns the transformed object.
// Returns an error if the provided argument is not a valid time object or the argument count is incorrect.
func timeToLocal(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewTime(t1.Local()), nil
}

// timeToUTC converts the given time object to its UTC equivalent and returns it.
// Returns an error if the number of arguments is incorrect or the input cannot be converted to a time object.
func timeToUTC(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewTime(t1.UTC()), nil
}

// timeTimeLocation returns the string representation of the location of the given time object.
// Accepts one argument of type IObject.
// Returns a string object containing the time location or an error if the input is invalid.
// Returns ErrWrongNumArguments if the number of arguments is not exactly one.
func timeTimeLocation(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewString(t1.Location().String())
}

// timeTimeString converts a single time object argument to its string representation and returns it as a new string object.
func timeTimeString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewString(t1.String())
}
