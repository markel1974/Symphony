package stdlib

import (
	"time"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// timesModule is a map providing time-related constants, formats, durations, months, and various time utility functions.
var _timesModule = map[string]objects.IObject{
	"format_ansic":         objects.NewStringNoSize(time.ANSIC),
	"format_unix_date":     objects.NewStringNoSize(time.UnixDate),
	"format_ruby_date":     objects.NewStringNoSize(time.RubyDate),
	"format_rfc822":        objects.NewStringNoSize(time.RFC822),
	"format_rfc822z":       objects.NewStringNoSize(time.RFC822Z),
	"format_rfc850":        objects.NewStringNoSize(time.RFC850),
	"format_rfc1123":       objects.NewStringNoSize(time.RFC1123),
	"format_rfc1123z":      objects.NewStringNoSize(time.RFC1123Z),
	"format_rfc3339":       objects.NewStringNoSize(time.RFC3339),
	"format_rfc3339_nano":  objects.NewStringNoSize(time.RFC3339Nano),
	"format_kitchen":       objects.NewStringNoSize(time.Kitchen),
	"format_stamp":         objects.NewStringNoSize(time.Stamp),
	"format_stamp_milli":   objects.NewStringNoSize(time.StampMilli),
	"format_stamp_micro":   objects.NewStringNoSize(time.StampMicro),
	"format_stamp_nano":    objects.NewStringNoSize(time.StampNano),
	"nanosecond":           objects.NewInt(int64(time.Nanosecond)),
	"microsecond":          objects.NewInt(int64(time.Microsecond)),
	"millisecond":          objects.NewInt(int64(time.Millisecond)),
	"second":               objects.NewInt(int64(time.Second)),
	"minute":               objects.NewInt(int64(time.Minute)),
	"hour":                 objects.NewInt(int64(time.Hour)),
	"january":              objects.NewInt(int64(time.January)),
	"february":             objects.NewInt(int64(time.February)),
	"march":                objects.NewInt(int64(time.March)),
	"april":                objects.NewInt(int64(time.April)),
	"may":                  objects.NewInt(int64(time.May)),
	"june":                 objects.NewInt(int64(time.June)),
	"july":                 objects.NewInt(int64(time.July)),
	"august":               objects.NewInt(int64(time.August)),
	"september":            objects.NewInt(int64(time.September)),
	"october":              objects.NewInt(int64(time.October)),
	"november":             objects.NewInt(int64(time.November)),
	"december":             objects.NewInt(int64(time.December)),
	"sleep":                objects.NewFunctionUser("sleep", timesSleep),                              // sleep(int)
	"parse_duration":       objects.NewFunctionUser("parse_duration", timesParseDuration),             // parse_duration(str) => int
	"since":                objects.NewFunctionUser("since", timesSince),                              // since(time) => int
	"until":                objects.NewFunctionUser("until", timesUntil),                              // until(time) => int
	"duration_hours":       objects.NewFunctionUser("duration_hours", timesDurationHours),             // duration_hours(int) => float
	"duration_minutes":     objects.NewFunctionUser("duration_minutes", timesDurationMinutes),         // duration_minutes(int) => float
	"duration_nanoseconds": objects.NewFunctionUser("duration_nanoseconds", timesDurationNanoseconds), // duration_nanoseconds(int) => int
	"duration_seconds":     objects.NewFunctionUser("duration_seconds", timesDurationSeconds),         // duration_seconds(int) => float
	"duration_string":      objects.NewFunctionUser("duration_string", timesDurationString),           // duration_string(int) => string
	"month_string":         objects.NewFunctionUser("month_string", timesMonthString),                 // month_string(int) => string
	"date":                 objects.NewFunctionUser("date", timesDate),                                // date(year, month, day, hour, min, sec, nsec) => time
	"now":                  objects.NewFunctionUser("now", timesNow),                                  // now() => time
	"parse":                objects.NewFunctionUser("parse", timesParse),                              // parse(format, str) => time
	"unix":                 objects.NewFunctionUser("unix", timesUnix),                                // unix(sec, nsec) => time
	"add":                  objects.NewFunctionUser("add", timesAdd),                                  // add(time, int) => time
	"add_date":             objects.NewFunctionUser("add_date", timesAddDate),                         // add_date(time, years, months, days) => time
	"sub":                  objects.NewFunctionUser("sub", timesSub),                                  // sub(t time, u time) => int
	"after":                objects.NewFunctionUser("after", timesAfter),                              // after(t time, u time) => bool
	"before":               objects.NewFunctionUser("before", timesBefore),                            // before(t time, u time) => bool
	"time_year":            objects.NewFunctionUser("time_year", timesTimeYear),                       // time_year(time) => int
	"time_month":           objects.NewFunctionUser("time_month", timesTimeMonth),                     // time_month(time) => int
	"time_day":             objects.NewFunctionUser("time_day", timesTimeDay),                         // time_day(time) => int
	"time_weekday":         objects.NewFunctionUser("time_weekday", timesTimeWeekday),                 // time_weekday(time) => int
	"time_hour":            objects.NewFunctionUser("time_hour", timesTimeHour),                       // time_hour(time) => int
	"time_minute":          objects.NewFunctionUser("time_minute", timesTimeMinute),                   // time_minute(time) => int
	"time_second":          objects.NewFunctionUser("time_second", timesTimeSecond),                   // time_second(time) => int
	"time_nanosecond":      objects.NewFunctionUser("time_nanosecond", timesTimeNanosecond),           // time_nanosecond(time) => int
	"time_unix":            objects.NewFunctionUser("time_unix", timesTimeUnix),                       // time_unix(time) => int
	"time_unix_nano":       objects.NewFunctionUser("time_unix_nano", timesTimeUnixNano),              // time_unix_nano(time) => int
	"time_format":          objects.NewFunctionUser("time_format", timesTimeFormat),                   // time_format(time, format) => string
	"time_location":        objects.NewFunctionUser("time_location", timesTimeLocation),               // time_location(time) => string
	"time_string":          objects.NewFunctionUser("time_string", timesTimeString),                   // time_string(time) => string
	"is_zero":              objects.NewFunctionUser("is_zero", timesIsZero),                           // is_zero(time) => bool
	"to_local":             objects.NewFunctionUser("to_local", timesToLocal),                         // to_local(time) => time
	"to_utc":               objects.NewFunctionUser("to_utc", timesToUTC),                             // to_utc(time) => time
}

// timesSleep pauses execution for the specified duration, given in nanoseconds, as an integer argument.
// Returns an error if the argument count is not one or if the argument is not int-compatible.
// The function returns an undefined value after the sleep duration.
func timesSleep(args ...objects.IObject) (objects.IObject, error) {
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

// timesParseDuration parses a duration string and returns its value as an integer representing nanoseconds or an error.
func timesParseDuration(args ...objects.IObject) (objects.IObject, error) {
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

// timesSince calculates the duration in nanoseconds between the provided time argument and the current time.
// Returns an integer object representing the difference or an error if the argument is invalid.
func timesSince(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(time.Since(t1))), nil
}

// timesUntil returns the number of nanoseconds until the provided time argument.
// Expects exactly one time or time-compatible argument; otherwise, it returns an error.
// If the argument is not a valid time-compatible object, an invalid argument type error is returned.
func timesUntil(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(time.Until(t1))), nil
}

// timesDurationHours converts an integer input representing nanoseconds to a float representing hours and returns it.
// Returns an error if input is not a single integer value or cannot be converted to int64.
func timesDurationHours(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewFloat(time.Duration(i1).Hours()), nil
}

// timesDurationMinutes converts the given integer argument to a time.Duration and returns its value in minutes as a Float.
// Returns an error if the argument count is not 1 or if the argument is not int-compatible.
func timesDurationMinutes(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewFloat(time.Duration(i1).Minutes()), nil
}

// timesDurationNanoseconds converts an integer input argument to a time duration in nanoseconds and returns it as an object.
// Returns an error if the wrong number of arguments is provided or if the input argument is not integer-compatible.
func timesDurationNanoseconds(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(time.Duration(i1).Nanoseconds()), nil
}

// timesDurationSeconds converts an integer argument to a time.Duration and returns its value in seconds as a float.
// Returns an error if the argument count is incorrect or if the provided argument is not an integer-compatible type.
func timesDurationSeconds(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewFloat(time.Duration(i1).Seconds()), nil
}

// timesDurationString converts an integer duration in nanoseconds to its string representation using Go's time.Duration.
// It takes exactly one argument of int-compatible type; otherwise, it returns an error.
// The result is returned as a string object representing the duration in a readable format.
func timesDurationString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewString(time.Duration(i1).String())
}

// timesMonthString converts an integer to its corresponding month name as a string. Returns an error for invalid arguments.
func timesMonthString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewStringNoSize(time.Month(i1).String()), nil
}

// timesDate creates a new time object using year, month, day, hour, minute, second, and nanosecond values as arguments.
// timesDate returns an error if the number of arguments is not seven or if any argument is not an int-compatible type.
func timesDate(args ...objects.IObject) (objects.IObject, error) {
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

// timesNow returns the current time as an object of type Time or an error if called with arguments.
func timesNow(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 0 {
		return nil, objects.ErrWrongNumArguments
	}
	return objects.NewTime(time.Now()), nil
}

// timesParse parses a time string based on a provided format string and returns a `Time` object or an error.
// Accepts exactly two arguments: a format string and a time string, both compatible with string types.
// Returns an error if the number of arguments is incorrect or their types are invalid.
// Wraps parsing errors as `Error` objects for consistent error handling.
func timesParse(args ...objects.IObject) (objects.IObject, error) {
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

// timesUnix creates a new Time object by interpreting two integers as Unix seconds and nanoseconds since the epoch.
// The first argument must be an integer for seconds, and the second argument must be an integer for nanoseconds.
// Returns an error if arguments are not exactly two or cannot be converted to integers.
func timesUnix(args ...objects.IObject) (objects.IObject, error) {
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

// timesAdd adds a duration (in nanoseconds) to a time object and returns the resulting time or an error if invalid arguments.
func timesAdd(args ...objects.IObject) (objects.IObject, error) {
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

// timesSub subtracts the second time argument from the first and returns the duration as an integer in nanoseconds.
// Returns an error if the number of arguments is not 2 or if the arguments are not time-compatible objects.
func timesSub(args ...objects.IObject) (objects.IObject, error) {
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

// timesAddDate adds years, months, and days to a time object and returns the resulting time as a new object.
// The function expects 4 arguments: a time object, an integer for years, an integer for months, and an integer for days.
// An error is returned if arguments are of the wrong types or if the number of arguments is not 4.
func timesAddDate(args ...objects.IObject) (objects.IObject, error) {
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

// timesAfter checks if the first time argument occurs after the second time argument and returns a boolean result.
// Returns an error if the number of arguments is not two or if the provided types are not time-compatible.
func timesAfter(args ...objects.IObject) (objects.IObject, error) {
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

// timesBefore checks if the first time argument is before the second time argument, returning true or false.
// Returns an error if the number of arguments is not 2 or if argument types are invalid.
func timesBefore(args ...objects.IObject) (objects.IObject, error) {
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

// timesTimeYear extracts the year from a given time object and returns it as an integer. It expects exactly one argument.
func timesTimeYear(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Year())), nil
}

// timesTimeMonth extracts the month from a given time object as an integer and returns it. Errors on invalid input or argument count.
func timesTimeMonth(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Month())), nil
}

// timesTimeDay retrieves the day of the month from a time object and returns it as an integer.
// Returns an error if the argument count is incorrect or if the input is not a valid time-compatible object.
func timesTimeDay(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Day())), nil
}

// timesTimeWeekday returns the weekday of the given time object as an integer.
// It requires exactly one argument of a type convertible to a time object.
// Returns an error if the argument count is incorrect or if the type is invalid.
func timesTimeWeekday(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Weekday())), nil
}

// timesTimeHour returns the hour component of the provided time object as an integer.
// Expects exactly one argument of type time-compatible object.
// Returns an error if the argument is not time-compatible or if the number of arguments is not one.
func timesTimeHour(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Hour())), nil
}

// timesTimeMinute returns the minute component of a time object as an integer.
// Expects exactly one argument of a time-compatible object.
func timesTimeMinute(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Minute())), nil
}

// timesTimeSecond extracts the second component of a time object, returning it as an integer. It expects one argument only.
func timesTimeSecond(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Second())), nil
}

// timesTimeNanosecond returns the nanosecond component of a time object as an integer.
// Returns an error if the argument is not a valid time object or if the number of arguments is not exactly one.
func timesTimeNanosecond(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Nanosecond())), nil
}

// timesTimeUnix converts a single time object to its Unix epoch timestamp representation and returns it as an integer.
// It returns an error if the number of arguments is not one or if the argument is not time-compatible.
func timesTimeUnix(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(t1.Unix()), nil
}

// timesTimeUnixNano converts a time object into its Unix time in nanoseconds and returns it as an integer object.
// An error is returned if the number or type of arguments is invalid.
func timesTimeUnixNano(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(t1.UnixNano()), nil
}

// timesTimeFormat formats a time object using a specified string-based layout and returns the formatted string.
// Returns an error if the arguments count is incorrect or the argument types are invalid.
// Returns an error if the resultant string exceeds the allowed maximum string length.
func timesTimeFormat(args ...objects.IObject) (objects.IObject, error) {
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

// timesIsZero checks if the provided time object is a zero time. Returns true if zero, false otherwise.
// Accepts one argument of type time-compatible object. Returns an error if the argument is invalid or missing.
func timesIsZero(args ...objects.IObject) (objects.IObject, error) {
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

// timesToLocal converts a time object to its local time equivalent or returns an error if the argument is invalid.
func timesToLocal(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewTime(t1.Local()), nil
}

// timesToUTC converts a single time-compatible object to its equivalent UTC time and returns it as a Time object.
// Returns an error if the argument count is incorrect or if the input is not a valid time-compatible type.
func timesToUTC(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewTime(t1.UTC()), nil
}

// timesTimeLocation extracts the location (timezone) of the given time object and returns it as a string.
// Accepts exactly one time-compatible argument and returns an error for invalid types or argument count.
func timesTimeLocation(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewString(t1.Location().String())
}

// timesTimeString converts a time object to its string representation and returns it as a string object.
// Returns an error if the argument is not a time-compatible object or if the number of arguments is not exactly one.
func timesTimeString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewString(t1.String())
}
