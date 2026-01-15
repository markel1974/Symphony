package sdk

import (
	"time"

	"github.com/markel1974/symphony/src/vm/bytecode"
	"github.com/markel1974/symphony/src/vm/objects"
)

func init() {
	register(NewTime)
}

// Time represents a structure that manages a collection of modules implementing the IObject interface.
type Time struct {
	*bytecode.Package
}

// NewTime initializes and returns a new instance of Time with predefined constants and functions mapped to the module.
func NewTime(factory objects.IGateKeeper) bytecode.IPackage {
	const (
		defANSIC               = "ANSIC"
		defUnixDate            = "UnixDate"
		defRubyDate            = "RubyDate"
		defRFC822              = "RFC822"
		defRFC822Z             = "RFC822Z"
		defRFC850              = "RFC850"
		defRFC1123             = "RFC1123"
		defRFC1123Z            = "RFC1123Z"
		defRFC3339             = "RFC3339"
		defRFC3339Nano         = "RFC3339Nano"
		defKitchen             = "Kitchen"
		defStamp               = "Stamp"
		defStampMilli          = "StampMilli"
		defStampMicro          = "StampMicro"
		defStampNano           = "StampNano"
		defNanosecond          = "Nanosecond"
		defMicrosecond         = "Microsecond"
		defMillisecond         = "Millisecond"
		defSecond              = "Second"
		defMinute              = "Minute"
		defHour                = "Hour"
		defJanuary             = "January"
		defFebruary            = "February"
		defMarch               = "March"
		defApril               = "April"
		defMay                 = "May"
		defJune                = "June"
		defJuly                = "July"
		defAugust              = "August"
		defSeptember           = "September"
		defOctober             = "October"
		defNovember            = "November"
		defDecember            = "December"
		defSleep               = "Sleep"
		defParseDuration       = "ParseDuration"
		defSince               = "Since"
		defUntil               = "Until"
		defDurationHours       = "DurationHours"
		defDurationMinutes     = "DurationMinutes"
		defDurationNanoseconds = "DurationNanoseconds"
		defDurationSeconds     = "DurationSeconds"
		defDurationString      = "DurationString"
		defMonthString         = "MonthString"
		defDate                = "Date"
		defNow                 = "Now"
		defParse               = "Parse"
		defUnix                = "Unix"
		defadd                 = "add"
		defAddDate             = "AddDate"
		defSub                 = "Sub"
		defAfter               = "After"
		defBefore              = "Before"
		defTimeYear            = "TimeYear"
		defTimeMonth           = "TimeMonth"
		defTimeDay             = "TimeDay"
		defTimeWeekday         = "TimeWeekday"
		defTimeHour            = "TimeHour"
		defTimeMinute          = "TimeMinute"
		defTimeSecond          = "TimeSecond"
		defTimeNanosecond      = "TimeNanosecond"
		defTimeUnix            = "TimeUnix"
		defTimeUnixNano        = "TimeUnixNano"
		defTimeFormat          = "TimeFormat"
		defTimeLocation        = "TimeLocation"
		defTimeString          = "TimeString"
		defIsZero              = "IsZero"
		defToLocal             = "ToLocal"
		defToUTC               = "ToUTC"
	)
	t := &Time{Package: bytecode.NewPackage("time")}
	t.Add(defANSIC, factory.NewString(objects.FrameStatic, time.ANSIC))
	t.Add(defUnixDate, factory.NewString(objects.FrameStatic, time.UnixDate))
	t.Add(defRubyDate, factory.NewString(objects.FrameStatic, time.RubyDate))
	t.Add(defRFC822, factory.NewString(objects.FrameStatic, time.RFC822))
	t.Add(defRFC822Z, factory.NewString(objects.FrameStatic, time.RFC822Z))
	t.Add(defRFC850, factory.NewString(objects.FrameStatic, time.RFC850))
	t.Add(defRFC1123, factory.NewString(objects.FrameStatic, time.RFC1123))
	t.Add(defRFC1123Z, factory.NewString(objects.FrameStatic, time.RFC1123Z))
	t.Add(defRFC3339, factory.NewString(objects.FrameStatic, time.RFC3339))
	t.Add(defRFC3339Nano, factory.NewString(objects.FrameStatic, time.RFC3339Nano))
	t.Add(defKitchen, factory.NewString(objects.FrameStatic, time.Kitchen))
	t.Add(defStamp, factory.NewString(objects.FrameStatic, time.Stamp))
	t.Add(defStampMilli, factory.NewString(objects.FrameStatic, time.StampMilli))
	t.Add(defStampMicro, factory.NewString(objects.FrameStatic, time.StampMicro))
	t.Add(defStampNano, factory.NewString(objects.FrameStatic, time.StampNano))
	t.Add(defNanosecond, factory.NewInt(objects.FrameStatic, int64(time.Nanosecond)))
	t.Add(defMicrosecond, factory.NewInt(objects.FrameStatic, int64(time.Microsecond)))
	t.Add(defMillisecond, factory.NewInt(objects.FrameStatic, int64(time.Millisecond)))
	t.Add(defSecond, factory.NewInt(objects.FrameStatic, int64(time.Second)))
	t.Add(defMinute, factory.NewInt(objects.FrameStatic, int64(time.Minute)))
	t.Add(defHour, factory.NewInt(objects.FrameStatic, int64(time.Hour)))
	t.Add(defJanuary, factory.NewInt(objects.FrameStatic, int64(time.January)))
	t.Add(defFebruary, factory.NewInt(objects.FrameStatic, int64(time.February)))
	t.Add(defMarch, factory.NewInt(objects.FrameStatic, int64(time.March)))
	t.Add(defApril, factory.NewInt(objects.FrameStatic, int64(time.April)))
	t.Add(defMay, factory.NewInt(objects.FrameStatic, int64(time.May)))
	t.Add(defJune, factory.NewInt(objects.FrameStatic, int64(time.June)))
	t.Add(defJuly, factory.NewInt(objects.FrameStatic, int64(time.July)))
	t.Add(defAugust, factory.NewInt(objects.FrameStatic, int64(time.August)))
	t.Add(defSeptember, factory.NewInt(objects.FrameStatic, int64(time.September)))
	t.Add(defOctober, factory.NewInt(objects.FrameStatic, int64(time.October)))
	t.Add(defNovember, factory.NewInt(objects.FrameStatic, int64(time.November)))
	t.Add(defDecember, factory.NewInt(objects.FrameStatic, int64(time.December)))
	t.Add(defSleep, factory.NewFuncImport(objects.FrameStatic, defSleep, 1, t.sleep))
	t.Add(defParseDuration, factory.NewFuncImport(objects.FrameStatic, defParseDuration, 1, t.parseDuration))
	t.Add(defSince, factory.NewFuncImport(objects.FrameStatic, defSince, 1, t.since))
	t.Add(defUntil, factory.NewFuncImport(objects.FrameStatic, defUntil, 1, t.until))
	t.Add(defDurationHours, factory.NewFuncImport(objects.FrameStatic, defDurationHours, 1, t.durationHours))
	t.Add(defDurationMinutes, factory.NewFuncImport(objects.FrameStatic, defDurationMinutes, 1, t.durationMinutes))
	t.Add(defDurationNanoseconds, factory.NewFuncImport(objects.FrameStatic, defDurationNanoseconds, 1, t.durationNanoseconds))
	t.Add(defDurationSeconds, factory.NewFuncImport(objects.FrameStatic, defDurationSeconds, 1, t.durationSeconds))
	t.Add(defDurationString, factory.NewFuncImport(objects.FrameStatic, defDurationString, 1, t.durationString))
	t.Add(defMonthString, factory.NewFuncImport(objects.FrameStatic, defMonthString, 1, t.monthString))
	t.Add(defDate, factory.NewFuncImport(objects.FrameStatic, defDate, 7, t.date))
	t.Add(defNow, factory.NewFuncImport(objects.FrameStatic, defNow, 0, t.now))
	t.Add(defParse, factory.NewFuncImport(objects.FrameStatic, defParse, 2, t.parse))
	t.Add(defUnix, factory.NewFuncImport(objects.FrameStatic, defUnix, 2, t.unix))
	t.Add(defadd, factory.NewFuncImport(objects.FrameStatic, defadd, 2, t.add))
	t.Add(defAddDate, factory.NewFuncImport(objects.FrameStatic, defAddDate, 4, t.addDate))
	t.Add(defSub, factory.NewFuncImport(objects.FrameStatic, defSub, 2, t.sub))
	t.Add(defAfter, factory.NewFuncImport(objects.FrameStatic, defAfter, 2, t.after))
	t.Add(defBefore, factory.NewFuncImport(objects.FrameStatic, defBefore, 2, t.before))
	t.Add(defTimeYear, factory.NewFuncImport(objects.FrameStatic, defTimeYear, 1, t.timeYear))
	t.Add(defTimeMonth, factory.NewFuncImport(objects.FrameStatic, defTimeMonth, 1, t.timeMonth))
	t.Add(defTimeDay, factory.NewFuncImport(objects.FrameStatic, defTimeDay, 1, t.timeDay))
	t.Add(defTimeWeekday, factory.NewFuncImport(objects.FrameStatic, defTimeWeekday, 1, t.timeWeekday))
	t.Add(defTimeHour, factory.NewFuncImport(objects.FrameStatic, defTimeHour, 1, t.timeHour))
	t.Add(defTimeMinute, factory.NewFuncImport(objects.FrameStatic, defTimeMinute, 1, t.timeMinute))
	t.Add(defTimeSecond, factory.NewFuncImport(objects.FrameStatic, defTimeSecond, 1, t.timeSecond))
	t.Add(defTimeNanosecond, factory.NewFuncImport(objects.FrameStatic, defTimeNanosecond, 1, t.timeNanosecond))
	t.Add(defTimeUnix, factory.NewFuncImport(objects.FrameStatic, defTimeUnix, 1, t.timeUnix))
	t.Add(defTimeUnixNano, factory.NewFuncImport(objects.FrameStatic, defTimeUnixNano, 1, t.timeUnixNano))
	t.Add(defTimeFormat, factory.NewFuncImport(objects.FrameStatic, defTimeFormat, 2, t.timeFormat))
	t.Add(defTimeLocation, factory.NewFuncImport(objects.FrameStatic, defTimeLocation, 1, t.timeLocation))
	t.Add(defTimeString, factory.NewFuncImport(objects.FrameStatic, defTimeString, 1, t.timeString))
	t.Add(defIsZero, factory.NewFuncImport(objects.FrameStatic, defIsZero, 1, t.isZero))
	t.Add(defToLocal, factory.NewFuncImport(objects.FrameStatic, defToLocal, 1, t.toLocal))
	t.Add(defToUTC, factory.NewFuncImport(objects.FrameStatic, defToUTC, 1, t.toUTC))
	return t
}

// sleep pauses the execution for a specified duration provided as an argument in nanoseconds.
// Returns an error if the number of arguments is incorrect or the argument is not an integer.
// On success, returns the undefined value.
func (t *Time) sleep(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	time.Sleep(time.Duration(i1))
	return 1, gk.UndefinedValue(), nil
}

// parseDuration parses a duration string and converts it into an integer representation of nanoseconds as IObject.
// Accepts exactly one argument of type string. Returns an error object if parsing or type conversion fails.
func (t *Time) parseDuration(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
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
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(time.Since(t1))), nil
}

// until calculates the duration from the current time to a specified time object and returns it as an integer object.
// Returns an error if the argument is missing, invalid, or not a time-compatible object.
func (t *Time) until(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(time.Until(t1))), nil
}

// durationHours calculates the duration in hours from an integer argument representing a duration in nanoseconds.
// It returns a Float object containing the result or an error if the input is invalid.
func (t *Time) durationHours(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewFloat(frame, time.Duration(i1).Hours()), nil
}

// durationMinutes calculates the duration in minutes based on the given integer argument and returns it as a float.
// Returns an error if the number of arguments is incorrect or if the argument type is invalid.
func (t *Time) durationMinutes(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewFloat(frame, time.Duration(i1).Minutes()), nil
}

// durationNanoseconds returns the nanosecond representation of a given duration argument as an IObject, or an error for invalid input.
func (t *Time) durationNanoseconds(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, time.Duration(i1).Nanoseconds()), nil
}

// durationSeconds converts the given integer argument (in nanoseconds) to a float representation of seconds.
func (t *Time) durationSeconds(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewFloat(frame, time.Duration(i1).Seconds()), nil
}

// durationString converts a duration given as an integer to its string representation and returns it as an IObject.
// Returns an error if not exactly one argument is provided or if the argument is not a valid integer.
func (t *Time) durationString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, time.Duration(i1).String()), nil
}

// monthString takes a single integer argument, converts it to a month, and returns its string representation.
// Returns an error if the argument count is incorrect or if the conversion fails.
func (t *Time) monthString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, time.Month(i1).String()), nil
}

// Date creates a new Time object using the specified year, month, day, hour, minute, second, and nanosecond values.
// It requires exactly 7 integer arguments and returns an error if the argument count or types are incorrect.
func (t *Time) date(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	i3, err := gk.ToInt64Arg(2, args)
	if err != nil {
		return 0, nil, err
	}
	i4, err := gk.ToInt64Arg(3, args)
	if err != nil {
		return 0, nil, err
	}
	i5, err := gk.ToInt64Arg(4, args)
	if err != nil {
		return 0, nil, err
	}
	i6, err := gk.ToInt64Arg(5, args)
	if err != nil {
		return 0, nil, err
	}
	i7, err := gk.ToInt64Arg(6, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewTime(frame, time.Date(int(i1), time.Month(i2), int(i3), int(i4), int(i5), int(i6), int(i7), time.Now().Location())), nil
}

// Now retrieves the current time as a Time object. Returns an error if any arguments are provided.
func (t *Time) now(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	return 1, gk.NewTime(frame, time.Now()), nil
}

// Parse parses a time string using the given format and returns a new Time object or an error if parsing fails.
func (t *Time) parse(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	s2, err := gk.ToStringArg(1, args)
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
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewTime(frame, time.Unix(i1, i2)), nil
}

// add adds a duration (int64) to a Time object and returns a new Time object or an error if the inputs are invalid.
func (t *Time) add(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewTime(frame, t1.Add(time.Duration(i2))), nil
}

// Sub calculates the duration between two time arguments and returns it as an Int object or an error if invalid arguments.
func (t *Time) sub(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	t2, err := gk.ToTimeArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Sub(t2))), nil
}

// AddDate adjusts the date by adding years, months, and days to the provided time object, returning the result as an IObject.
// It expects four arguments: a time object, and three integers representing years, months, and days respectively.
// Returns an error if the wrong number of arguments is provided or a type conversion fails.
func (t *Time) addDate(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	i3, err := gk.ToInt64Arg(2, args)
	if err != nil {
		return 0, nil, err
	}
	i4, err := gk.ToInt64Arg(3, args)
	if err != nil {
		return 0, nil, err
	}
	v := t1.AddDate(int(i2), int(i3), int(i4))
	return 1, gk.NewTime(frame, v), nil
}

// After compares two time values and returns TrueValue if the first is after the second, otherwise returns FalseValue.
func (t *Time) after(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	t2, err := gk.ToTimeArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	if t1.After(t2) {
		return 1, gk.TrueValue(), nil
	}
	return 1, gk.FalseValue(), nil
}

// Before determines if the first time argument occurs before the second time argument and returns a boolean result.
func (t *Time) before(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	t2, err := gk.ToTimeArg(1, args)
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
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Year())), nil
}

// TimeMonth extracts the month from a time object and returns it as an integer. It requires exactly one argument.
// Returns an error if the argument count is incorrect or the type conversion fails.
func (t *Time) timeMonth(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Month())), nil
}

// TimeDay extracts and returns the day of the month as an integer from a given time object. It requires exactly one argument.
func (t *Time) timeDay(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Day())), nil
}

// TimeWeekday returns the weekday of a time object as an integer. Returns an error if the argument count is invalid.
func (t *Time) timeWeekday(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Weekday())), nil
}

// TimeHour extracts the hour from the given time object and returns it as an Int. Returns an error if arguments are invalid.
func (t *Time) timeHour(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Hour())), nil
}

// TimeMinute extracts the minute component from a time object and returns it as an integer.
func (t *Time) timeMinute(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Minute())), nil
}

// TimeSecond extracts the second component from a time object passed as an argument and returns it as an Int object.
// Returns an error if the argument count is not 1 or if the conversion to a time object fails.
func (t *Time) timeSecond(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Second())), nil
}

// timeNanosecond returns the nanosecond component of the given time object as an integer.
// Expects a single argument of a time-compatible object. Returns an error for invalid arguments or conversion failures.
func (t *Time) timeNanosecond(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, int64(t1.Nanosecond())), nil
}

// timeUnix converts a provided time object into its Unix timestamp and returns it as an integer object.
func (t *Time) timeUnix(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, t1.Unix()), nil
}

// timeUnixNano returns the Unix time in nanoseconds as an IObject for the given time argument. An error is returned for invalid input.
func (t *Time) timeUnixNano(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewInt(frame, t1.UnixNano()), nil
}

// timeFormat formats a time object using the provided format and returns the formatted string as an IObject.
func (t *Time) timeFormat(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	s2, err := gk.ToStringArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	s := t1.Format(s2)
	return 1, gk.NewString(frame, s), nil
}

// isZero checks if the provided time argument is zero and returns TrueValue if it is, otherwise returns FalseValue.
func (t *Time) isZero(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
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
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewTime(frame, t1.Local()), nil
}

// toUTC converts the provided IObject time argument to UTC and returns a new IObject representing the UTC time.
func (t *Time) toUTC(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewTime(frame, t1.UTC()), nil
}

// timeLocation returns the location (timezone) from the given time object as a string. Takes exactly one argument.
func (t *Time) timeLocation(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, t1.Location().String()), nil
}

// timeString converts a time instance to its string representation. It requires exactly one argument of type IObject.
func (t *Time) timeString(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	t1, err := gk.ToTimeArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, t1.String()), nil
}
