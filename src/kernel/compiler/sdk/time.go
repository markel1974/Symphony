package sdk

import (
	"time"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Time represents a structure that manages a collection of modules implementing the IObject interface.
type Time struct {
	*Module
}

// NewTime initializes and returns a new instance of Time with predefined constants and functions mapped to the module.
func NewTime() *Time {
	t := &Time{
		Module: NewModule(),
	}
	t.attrs = map[string]objects.IObject{
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
		"Sleep":               objects.NewFunctionModule(objects.FunctionModuleDef, "Sleep", t.Sleep),                             // sleep(int)
		"ParseDuration":       objects.NewFunctionModule(objects.FunctionModuleDef, "ParseDuration", t.ParseDuration),             // parse_duration(str) => int
		"Since":               objects.NewFunctionModule(objects.FunctionModuleDef, "Since", t.Since),                             // since(time) => int
		"Until":               objects.NewFunctionModule(objects.FunctionModuleDef, "Until", t.Until),                             // until(time) => int
		"DurationHours":       objects.NewFunctionModule(objects.FunctionModuleDef, "DurationHours", t.DurationHours),             // duration_hours(int) => float
		"DurationMinutes":     objects.NewFunctionModule(objects.FunctionModuleDef, "DurationMinutes", t.DurationMinutes),         // duration_minutes(int) => float
		"DurationNanoseconds": objects.NewFunctionModule(objects.FunctionModuleDef, "DurationNanoseconds", t.DurationNanoseconds), // duration_nanoseconds(int) => int
		"DurationSeconds":     objects.NewFunctionModule(objects.FunctionModuleDef, "DurationSeconds", t.DurationSeconds),         // duration_seconds(int) => float
		"DurationString":      objects.NewFunctionModule(objects.FunctionModuleDef, "DurationString", t.DurationString),           // duration_string(int) => string
		"MonthString":         objects.NewFunctionModule(objects.FunctionModuleDef, "MonthString", t.MonthString),                 // month_string(int) => string
		"Date":                objects.NewFunctionModule(objects.FunctionModuleDef, "Date", t.Date),                               // date(year, month, day, hour, min, sec, nsec) => time
		"Now":                 objects.NewFunctionModule(objects.FunctionModuleDef, "Now", t.Now),                                 // now() => time
		"Parse":               objects.NewFunctionModule(objects.FunctionModuleDef, "Parse", t.Parse),                             // parse(format, str) => time
		"Unix":                objects.NewFunctionModule(objects.FunctionModuleDef, "Unix", t.Unix),                               // unix(sec, nsec) => time
		"Add":                 objects.NewFunctionModule(objects.FunctionModuleDef, "Add", t.Add),                                 // add(time, int) => time
		"AddDate":             objects.NewFunctionModule(objects.FunctionModuleDef, "AddDate", t.AddDate),                         // add_date(time, years, months, days) => time
		"Sub":                 objects.NewFunctionModule(objects.FunctionModuleDef, "Sub", t.Sub),                                 // sub(t time, u time) => int
		"After":               objects.NewFunctionModule(objects.FunctionModuleDef, "After", t.After),                             // after(t time, u time) => bool
		"Before":              objects.NewFunctionModule(objects.FunctionModuleDef, "Before", t.Before),                           // before(t time, u time) => bool
		"TimeYear":            objects.NewFunctionModule(objects.FunctionModuleDef, "TimeYear", t.TimeYear),                       // time_year(time) => int
		"TimeMonth":           objects.NewFunctionModule(objects.FunctionModuleDef, "TimeMonth", t.TimeMonth),                     // time_month(time) => int
		"TimeDay":             objects.NewFunctionModule(objects.FunctionModuleDef, "TimeDay", t.TimeDay),                         // time_day(time) => int
		"TimeWeekday":         objects.NewFunctionModule(objects.FunctionModuleDef, "TimeWeekday", t.TimeWeekday),                 // time_weekday(time) => int
		"TimeHour":            objects.NewFunctionModule(objects.FunctionModuleDef, "TimeHour", t.TimeHour),                       // time_hour(time) => int
		"TimeMinute":          objects.NewFunctionModule(objects.FunctionModuleDef, "TimeMinute", t.TimeMinute),                   // time_minute(time) => int
		"TimeSecond":          objects.NewFunctionModule(objects.FunctionModuleDef, "TimeSecond", t.TimeSecond),                   // time_second(time) => int
		"TimeNanosecond":      objects.NewFunctionModule(objects.FunctionModuleDef, "TimeNanosecond", t.TimeNanosecond),           // time_nanosecond(time) => int
		"TimeUnix":            objects.NewFunctionModule(objects.FunctionModuleDef, "TimeUnix", t.TimeUnix),                       // time_unix(time) => int
		"TimeUnixNano":        objects.NewFunctionModule(objects.FunctionModuleDef, "TimeUnixNano", t.TimeUnixNano),               // time_unix_nano(time) => int
		"TimeFormat":          objects.NewFunctionModule(objects.FunctionModuleDef, "TimeFormat", t.TimeFormat),                   // time_format(time, format) => string
		"TimeLocation":        objects.NewFunctionModule(objects.FunctionModuleDef, "TimeLocation", t.TimeLocation),               // time_location(time) => string
		"TimeString":          objects.NewFunctionModule(objects.FunctionModuleDef, "TimeString", t.TimeString),                   // time_string(time) => string
		"IsZero":              objects.NewFunctionModule(objects.FunctionModuleDef, "is_zero", t.IsZero),                          // is_zero(time) => bool
		"ToLocal":             objects.NewFunctionModule(objects.FunctionModuleDef, "to_local", t.ToLocal),                        // to_local(time) => time
		"ToUTC":               objects.NewFunctionModule(objects.FunctionModuleDef, "to_utc", t.ToUTC),                            // to_utc(time) => time
	}

	return t
}

// Name returns the name of Time module.
func (t *Time) Name() string {
	return "time"
}

// Sleep pauses the execution for a specified duration provided as an argument in nanoseconds.
// Returns an error if the number of arguments is incorrect or the argument is not an integer.
// On success, returns the undefined value.
func (t *Time) Sleep(args ...objects.IObject) (objects.IObject, error) {
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

// ParseDuration parses a duration string and converts it into an integer representation of nanoseconds as IObject.
// Accepts exactly one argument of type string. Returns an error object if parsing or type conversion fails.
func (t *Time) ParseDuration(args ...objects.IObject) (objects.IObject, error) {
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

// Since calculates the time duration between the current Time instance and the given time argument as an integer.
// Expects exactly one argument of a compatible time type, returns an error if the argument is invalid or missing.
func (t *Time) Since(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(time.Since(t1))), nil
}

// Until calculates the duration from the current time to a specified time object and returns it as an integer object.
// Returns an error if the argument is missing, invalid, or not a time-compatible object.
func (t *Time) Until(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(time.Until(t1))), nil
}

// DurationHours calculates the duration in hours from an integer argument representing a duration in nanoseconds.
// It returns a Float object containing the result or an error if the input is invalid.
func (t *Time) DurationHours(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewFloat(time.Duration(i1).Hours()), nil
}

// DurationMinutes calculates the duration in minutes based on the given integer argument and returns it as a float.
// Returns an error if the number of arguments is incorrect or if the argument type is invalid.
func (t *Time) DurationMinutes(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewFloat(time.Duration(i1).Minutes()), nil
}

// DurationNanoseconds returns the nanosecond representation of a given duration argument as an IObject, or an error for invalid input.
func (t *Time) DurationNanoseconds(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(time.Duration(i1).Nanoseconds()), nil
}

// DurationSeconds converts the given integer argument (in nanoseconds) to a float representation of seconds.
func (t *Time) DurationSeconds(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewFloat(time.Duration(i1).Seconds()), nil
}

// DurationString converts a duration given as an integer to its string representation and returns it as an IObject.
// Returns an error if not exactly one argument is provided or if the argument is not a valid integer.
func (t *Time) DurationString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewString(time.Duration(i1).String())
}

// MonthString takes a single integer argument, converts it to a month, and returns its string representation.
// Returns an error if the argument count is incorrect or if the conversion fails.
func (t *Time) MonthString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewStringNoSize(time.Month(i1).String()), nil
}

// Date creates a new Time object using the specified year, month, day, hour, minute, second, and nanosecond values.
// It requires exactly 7 integer arguments and returns an error if the argument count or types are incorrect.
func (t *Time) Date(args ...objects.IObject) (objects.IObject, error) {
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

// Now retrieves the current time as a Time object. Returns an error if any arguments are provided.
func (t *Time) Now(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 0 {
		return nil, objects.ErrWrongNumArguments
	}
	return objects.NewTime(time.Now()), nil
}

// Parse parses a time string using the given format and returns a new Time object or an error if parsing fails.
func (t *Time) Parse(args ...objects.IObject) (objects.IObject, error) {
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

// Unix creates a new Time object based on the given Unix timestamp and nanoseconds, or returns an error for invalid arguments.
func (t *Time) Unix(args ...objects.IObject) (objects.IObject, error) {
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

// Add adds a duration (int64) to a Time object and returns a new Time object or an error if the inputs are invalid.
func (t *Time) Add(args ...objects.IObject) (objects.IObject, error) {
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

// Sub calculates the duration between two time arguments and returns it as an Int object or an error if invalid arguments.
func (t *Time) Sub(args ...objects.IObject) (objects.IObject, error) {
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

// AddDate adjusts the date by adding years, months, and days to the provided time object, returning the result as an IObject.
// It expects four arguments: a time object, and three integers representing years, months, and days respectively.
// Returns an error if the wrong number of arguments is provided or a type conversion fails.
func (t *Time) AddDate(args ...objects.IObject) (objects.IObject, error) {
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

// After compares two time values and returns TrueValue if the first is after the second, otherwise returns FalseValue.
func (t *Time) After(args ...objects.IObject) (objects.IObject, error) {
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

// Before determines if the first time argument occurs before the second time argument and returns a boolean result.
func (t *Time) Before(args ...objects.IObject) (objects.IObject, error) {
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

// TimeYear returns the year component of a given time object as an integer. Accepts a single argument of type IObject.
func (t *Time) TimeYear(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Year())), nil
}

// TimeMonth extracts the month from a time object and returns it as an integer. It requires exactly one argument.
// Returns an error if the argument count is incorrect or the type conversion fails.
func (t *Time) TimeMonth(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Month())), nil
}

// TimeDay extracts and returns the day of the month as an integer from a given time object. It requires exactly one argument.
func (t *Time) TimeDay(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Day())), nil
}

// TimeWeekday returns the weekday of a time object as an integer. Returns an error if the argument count is invalid.
func (t *Time) TimeWeekday(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Weekday())), nil
}

// TimeHour extracts the hour from the given time object and returns it as an Int. Returns an error if arguments are invalid.
func (t *Time) TimeHour(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Hour())), nil
}

// TimeMinute extracts the minute component from a time object and returns it as an integer.
func (t *Time) TimeMinute(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Minute())), nil
}

// TimeSecond extracts the second component from a time object passed as an argument and returns it as an Int object.
// Returns an error if the argument count is not 1 or if the conversion to a time object fails.
func (t *Time) TimeSecond(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Second())), nil
}

// TimeNanosecond returns the nanosecond component of the given time object as an integer.
// Expects a single argument of a time-compatible object. Returns an error for invalid arguments or conversion failures.
func (t *Time) TimeNanosecond(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(int64(t1.Nanosecond())), nil
}

// TimeUnix converts a provided time object into its Unix timestamp and returns it as an integer object.
func (t *Time) TimeUnix(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(t1.Unix()), nil
}

// TimeUnixNano returns the Unix time in nanoseconds as an IObject for the given time argument. An error is returned for invalid input.
func (t *Time) TimeUnixNano(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewInt(t1.UnixNano()), nil
}

// TimeFormat formats a time object using the provided format and returns the formatted string as an IObject.
func (t *Time) TimeFormat(args ...objects.IObject) (objects.IObject, error) {
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

// IsZero checks if the provided time argument is zero and returns TrueValue if it is, otherwise returns FalseValue.
func (t *Time) IsZero(args ...objects.IObject) (objects.IObject, error) {
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

// ToLocal converts the given IObject argument to a local time zone Time object or returns an error if conversion fails.
func (t *Time) ToLocal(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewTime(t1.Local()), nil
}

// ToUTC converts the provided IObject time argument to UTC and returns a new IObject representing the UTC time.
func (t *Time) ToUTC(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewTime(t1.UTC()), nil
}

// TimeLocation returns the location (timezone) from the given time object as a string. Takes exactly one argument.
func (t *Time) TimeLocation(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewString(t1.Location().String())
}

// TimeString converts a time instance to its string representation. It requires exactly one argument of type IObject.
func (t *Time) TimeString(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	t1, err := objects.ToTimeArg(0, args[0])
	if err != nil {
		return nil, err
	}
	return objects.NewString(t1.String())
}
