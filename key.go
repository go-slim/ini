package ini

import (
	"bytes"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Key represents a key under a section.
type Key struct {
	s               *Section
	name            string
	value           string
	Comment         string
	isAutoIncrement bool
	isBooleanType   bool
}

// newKey simply return a key object with given values.
func newKey(s *Section, name, val string) *Key {
	return &Key{
		s:     s,
		name:  name,
		value: val,
	}
}

// Name returns name of key.
func (k *Key) Name() string {
	return k.name
}

// Value returns raw value of key for performance purpose.
func (k *Key) Value() string {
	return k.value
}

// String returns string representation of value.
func (k *Key) String() string {
	return transformValue(k)
}

// Validate accepts a validate function which can
// return modified result as key value.
func (k *Key) Validate(fn func(string) string) string {
	return fn(k.String())
}

// parseBool returns the boolean value represented by the string.
//
// It accepts 1, t, T, TRUE, true, True, YES, yes, Yes, y, ON, on, On,
// 0, f, F, FALSE, false, False, NO, no, No, n, OFF, off, Off.
// Any other value returns an error.
func parseBool(str string) (value bool, err error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "y", "ON", "on", "On":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "n", "OFF", "off", "Off":
		return false, nil
	}
	return false, fmt.Errorf("parsing \"%s\": invalid syntax", str)
}

// Bool returns bool type value.
func (k *Key) Bool() (bool, error) {
	return parseBool(k.String())
}

// Float64 returns float64 type value.
func (k *Key) Float64() (float64, error) {
	return strconv.ParseFloat(k.String(), 64)
}

// Int returns int type value.
func (k *Key) Int() (int, error) {
	v, err := strconv.ParseInt(k.String(), 0, 64)
	return int(v), err
}

// Int64 returns int64 type value.
func (k *Key) Int64() (int64, error) {
	return strconv.ParseInt(k.String(), 0, 64)
}

// Uint returns uint type valued.
func (k *Key) Uint() (uint, error) {
	u, e := strconv.ParseUint(k.String(), 0, 64)
	return uint(u), e
}

// Uint64 returns uint64 type value.
func (k *Key) Uint64() (uint64, error) {
	return strconv.ParseUint(k.String(), 0, 64)
}

// Duration returns time.Duration type value.
func (k *Key) Duration() (time.Duration, error) {
	return time.ParseDuration(k.String())
}

// TimeFormat parses with given format and returns time.Time type value.
func (k *Key) TimeFormat(format string) (time.Time, error) {
	return time.Parse(format, k.String())
}

// Time parses with RFC3339 format and returns time.Time type value.
func (k *Key) Time() (time.Time, error) {
	return k.TimeFormat(time.RFC3339)
}

// MustString returns default value if key value is empty.
func (k *Key) MustString(defaultVal string) string {
	val := k.String()
	if len(val) == 0 {
		k.value = defaultVal
		return defaultVal
	}
	return val
}

// MustBool always returns value without error,
// it returns false if error occurs.
func (k *Key) MustBool(defaultVal ...bool) bool {
	val, err := k.Bool()
	if len(defaultVal) > 0 && err != nil {
		k.value = strconv.FormatBool(defaultVal[0])
		return defaultVal[0]
	}
	return val
}

// MustFloat64 always returns value without error,
// it returns 0.0 if error occurs.
func (k *Key) MustFloat64(defaultVal ...float64) float64 {
	val, err := k.Float64()
	if len(defaultVal) > 0 && err != nil {
		k.value = strconv.FormatFloat(defaultVal[0], 'f', -1, 64)
		return defaultVal[0]
	}
	return val
}

// MustInt always returns value without error,
// it returns 0 if error occurs.
func (k *Key) MustInt(defaultVal ...int) int {
	val, err := k.Int()
	if len(defaultVal) > 0 && err != nil {
		k.value = strconv.FormatInt(int64(defaultVal[0]), 10)
		return defaultVal[0]
	}
	return val
}

// MustInt64 always returns value without error,
// it returns 0 if error occurs.
func (k *Key) MustInt64(defaultVal ...int64) int64 {
	val, err := k.Int64()
	if len(defaultVal) > 0 && err != nil {
		k.value = strconv.FormatInt(defaultVal[0], 10)
		return defaultVal[0]
	}
	return val
}

// MustUint always returns value without error,
// it returns 0 if error occurs.
func (k *Key) MustUint(defaultVal ...uint) uint {
	val, err := k.Uint()
	if len(defaultVal) > 0 && err != nil {
		k.value = strconv.FormatUint(uint64(defaultVal[0]), 10)
		return defaultVal[0]
	}
	return val
}

// MustUint64 always returns value without error,
// it returns 0 if error occurs.
func (k *Key) MustUint64(defaultVal ...uint64) uint64 {
	val, err := k.Uint64()
	if len(defaultVal) > 0 && err != nil {
		k.value = strconv.FormatUint(defaultVal[0], 10)
		return defaultVal[0]
	}
	return val
}

// MustDuration always returns value without error,
// it returns zero value if error occurs.
func (k *Key) MustDuration(defaultVal ...time.Duration) time.Duration {
	val, err := k.Duration()
	if len(defaultVal) > 0 && err != nil {
		k.value = defaultVal[0].String()
		return defaultVal[0]
	}
	return val
}

// MustTimeFormat always parses with given format and returns value without error,
// it returns zero value if error occurs.
func (k *Key) MustTimeFormat(format string, defaultVal ...time.Time) time.Time {
	val, err := k.TimeFormat(format)
	if len(defaultVal) > 0 && err != nil {
		k.value = defaultVal[0].Format(format)
		return defaultVal[0]
	}
	return val
}

// MustTime always parses with RFC3339 format and returns value without error,
// it returns zero value if error occurs.
func (k *Key) MustTime(defaultVal ...time.Time) time.Time {
	return k.MustTimeFormat(time.RFC3339, defaultVal...)
}

// In always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *Key) In(defaultVal string, candidates []string) string {
	val := k.String()
	for _, cand := range candidates {
		if val == cand {
			return val
		}
	}
	return defaultVal
}

// InFloat64 always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *Key) InFloat64(defaultVal float64, candidates []float64) float64 {
	val := k.MustFloat64()
	for _, cand := range candidates {
		if val == cand {
			return val
		}
	}
	return defaultVal
}

// InInt always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *Key) InInt(defaultVal int, candidates []int) int {
	val := k.MustInt()
	for _, cand := range candidates {
		if val == cand {
			return val
		}
	}
	return defaultVal
}

// InInt64 always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *Key) InInt64(defaultVal int64, candidates []int64) int64 {
	val := k.MustInt64()
	for _, cand := range candidates {
		if val == cand {
			return val
		}
	}
	return defaultVal
}

// InUint always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *Key) InUint(defaultVal uint, candidates []uint) uint {
	val := k.MustUint()
	for _, cand := range candidates {
		if val == cand {
			return val
		}
	}
	return defaultVal
}

// InUint64 always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *Key) InUint64(defaultVal uint64, candidates []uint64) uint64 {
	val := k.MustUint64()
	if slices.Contains(candidates, val) {
		return val
	}
	return defaultVal
}

// InTimeFormat always parses with given format and returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *Key) InTimeFormat(format string, defaultVal time.Time, candidates []time.Time) time.Time {
	val := k.MustTimeFormat(format)
	if slices.Contains(candidates, val) {
		return val
	}
	return defaultVal
}

// InTime always parses with RFC3339 format and returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *Key) InTime(defaultVal time.Time, candidates []time.Time) time.Time {
	return k.InTimeFormat(time.RFC3339, defaultVal, candidates)
}

// RangeFloat64 checks if value is in given range inclusively,
// and returns default value if it's not.
func (k *Key) RangeFloat64(defaultVal, min, max float64) float64 {
	val := k.MustFloat64()
	if val < min || val > max {
		return defaultVal
	}
	return val
}

// RangeInt checks if value is in given range inclusively,
// and returns default value if it's not.
func (k *Key) RangeInt(defaultVal, min, max int) int {
	val := k.MustInt()
	if val < min || val > max {
		return defaultVal
	}
	return val
}

// RangeInt64 checks if value is in given range inclusively,
// and returns default value if it's not.
func (k *Key) RangeInt64(defaultVal, min, max int64) int64 {
	val := k.MustInt64()
	if val < min || val > max {
		return defaultVal
	}
	return val
}

// RangeTimeFormat checks if value with given format is in given range inclusively,
// and returns default value if it's not.
func (k *Key) RangeTimeFormat(format string, defaultVal, min, max time.Time) time.Time {
	val := k.MustTimeFormat(format)
	if val.Unix() < min.Unix() || val.Unix() > max.Unix() {
		return defaultVal
	}
	return val
}

// RangeTime checks if value with RFC3339 format is in given range inclusively,
// and returns default value if it's not.
func (k *Key) RangeTime(defaultVal, min, max time.Time) time.Time {
	return k.RangeTimeFormat(time.RFC3339, defaultVal, min, max)
}

// Strings returns list of string divided by given delimiter.
func (k *Key) Strings(delim string) []string {
	str := k.String()
	if len(str) == 0 {
		return []string{}
	}

	runes := []rune(str)
	vals := make([]string, 0, 2)
	var buf bytes.Buffer
	escape := false
	idx := 0
	for {
		if escape {
			escape = false
			if runes[idx] != '\\' && !strings.HasPrefix(string(runes[idx:]), delim) {
				buf.WriteRune('\\')
			}
			buf.WriteRune(runes[idx])
		} else {
			if runes[idx] == '\\' {
				escape = true
			} else if strings.HasPrefix(string(runes[idx:]), delim) {
				idx += len(delim) - 1
				vals = append(vals, strings.TrimSpace(buf.String()))
				buf.Reset()
			} else {
				buf.WriteRune(runes[idx])
			}
		}
		idx++
		if idx == len(runes) {
			break
		}
	}

	if buf.Len() > 0 {
		vals = append(vals, strings.TrimSpace(buf.String()))
	}

	return vals
}

// Float64s returns list of float64 divided by given delimiter. Any invalid input will be treated as zero value.
func (k *Key) Float64s(delim string) []float64 {
	vals, _ := k.parseFloat64s(k.Strings(delim), true, false)
	return vals
}

// Ints returns list of int divided by given delimiter. Any invalid input will be treated as zero value.
func (k *Key) Ints(delim string) []int {
	vals, _ := k.parseInts(k.Strings(delim), true, false)
	return vals
}

// Int64s returns list of int64 divided by given delimiter. Any invalid input will be treated as zero value.
func (k *Key) Int64s(delim string) []int64 {
	vals, _ := k.parseInt64s(k.Strings(delim), true, false)
	return vals
}

// Uints returns list of uint divided by given delimiter. Any invalid input will be treated as zero value.
func (k *Key) Uints(delim string) []uint {
	vals, _ := k.parseUints(k.Strings(delim), true, false)
	return vals
}

// Uint64s returns list of uint64 divided by given delimiter. Any invalid input will be treated as zero value.
func (k *Key) Uint64s(delim string) []uint64 {
	vals, _ := k.parseUint64s(k.Strings(delim), true, false)
	return vals
}

// Bools returns list of bool divided by given delimiter. Any invalid input will be treated as zero value.
func (k *Key) Bools(delim string) []bool {
	vals, _ := k.parseBools(k.Strings(delim), true, false)
	return vals
}

// TimesFormat parses with given format and returns list of time.Time divided by given delimiter.
// Any invalid input will be treated as zero value (0001-01-01 00:00:00 +0000 UTC).
func (k *Key) TimesFormat(format, delim string) []time.Time {
	vals, _ := k.parseTimesFormat(format, k.Strings(delim), true, false)
	return vals
}

// Times parses with RFC3339 format and returns list of time.Time divided by given delimiter.
// Any invalid input will be treated as zero value (0001-01-01 00:00:00 +0000 UTC).
func (k *Key) Times(delim string) []time.Time {
	return k.TimesFormat(time.RFC3339, delim)
}

// ValidFloat64s returns list of float64 divided by given delimiter. If some value is not float, then
// it will not be included to result list.
func (k *Key) ValidFloat64s(delim string) []float64 {
	vals, _ := k.parseFloat64s(k.Strings(delim), false, false)
	return vals
}

// ValidInts returns list of int divided by given delimiter. If some value is not integer, then it will
// not be included to result list.
func (k *Key) ValidInts(delim string) []int {
	vals, _ := k.parseInts(k.Strings(delim), false, false)
	return vals
}

// ValidInt64s returns list of int64 divided by given delimiter. If some value is not 64-bit integer,
// then it will not be included to result list.
func (k *Key) ValidInt64s(delim string) []int64 {
	vals, _ := k.parseInt64s(k.Strings(delim), false, false)
	return vals
}

// ValidUints returns list of uint divided by given delimiter. If some value is not unsigned integer,
// then it will not be included to result list.
func (k *Key) ValidUints(delim string) []uint {
	vals, _ := k.parseUints(k.Strings(delim), false, false)
	return vals
}

// ValidUint64s returns list of uint64 divided by given delimiter. If some value is not 64-bit unsigned
// integer, then it will not be included to result list.
func (k *Key) ValidUint64s(delim string) []uint64 {
	vals, _ := k.parseUint64s(k.Strings(delim), false, false)
	return vals
}

// ValidBools returns list of bool divided by given delimiter. If some value is not 64-bit unsigned
// integer, then it will not be included to result list.
func (k *Key) ValidBools(delim string) []bool {
	vals, _ := k.parseBools(k.Strings(delim), false, false)
	return vals
}

// ValidTimesFormat parses with given format and returns list of time.Time divided by given delimiter.
func (k *Key) ValidTimesFormat(format, delim string) []time.Time {
	vals, _ := k.parseTimesFormat(format, k.Strings(delim), false, false)
	return vals
}

// ValidTimes parses with RFC3339 format and returns list of time.Time divided by given delimiter.
func (k *Key) ValidTimes(delim string) []time.Time {
	return k.ValidTimesFormat(time.RFC3339, delim)
}

// StrictFloat64s returns list of float64 divided by given delimiter or error on first invalid input.
func (k *Key) StrictFloat64s(delim string) ([]float64, error) {
	return k.parseFloat64s(k.Strings(delim), false, true)
}

// StrictInts returns list of int divided by given delimiter or error on first invalid input.
func (k *Key) StrictInts(delim string) ([]int, error) {
	return k.parseInts(k.Strings(delim), false, true)
}

// StrictInt64s returns list of int64 divided by given delimiter or error on first invalid input.
func (k *Key) StrictInt64s(delim string) ([]int64, error) {
	return k.parseInt64s(k.Strings(delim), false, true)
}

// StrictUints returns list of uint divided by given delimiter or error on first invalid input.
func (k *Key) StrictUints(delim string) ([]uint, error) {
	return k.parseUints(k.Strings(delim), false, true)
}

// StrictUint64s returns list of uint64 divided by given delimiter or error on first invalid input.
func (k *Key) StrictUint64s(delim string) ([]uint64, error) {
	return k.parseUint64s(k.Strings(delim), false, true)
}

// StrictBools returns list of bool divided by given delimiter or error on first invalid input.
func (k *Key) StrictBools(delim string) ([]bool, error) {
	return k.parseBools(k.Strings(delim), false, true)
}

// StrictTimesFormat parses with given format and returns list of time.Time divided by given delimiter
// or error on first invalid input.
func (k *Key) StrictTimesFormat(format, delim string) ([]time.Time, error) {
	return k.parseTimesFormat(format, k.Strings(delim), false, true)
}

// StrictTimes parses with RFC3339 format and returns list of time.Time divided by given delimiter
// or error on first invalid input.
func (k *Key) StrictTimes(delim string) ([]time.Time, error) {
	return k.StrictTimesFormat(time.RFC3339, delim)
}

// parseBools transforms strings to bools.
func (k *Key) parseBools(strs []string, addInvalid, returnOnInvalid bool) ([]bool, error) {
	vals := make([]bool, 0, len(strs))
	parser := func(str string) (interface{}, error) {
		val, err := parseBool(str)
		return val, err
	}
	rawVals, err := k.doParse(strs, addInvalid, returnOnInvalid, parser)
	if err == nil {
		for _, val := range rawVals {
			vals = append(vals, val.(bool))
		}
	}
	return vals, err
}

// parseFloat64s transforms strings to float64s.
func (k *Key) parseFloat64s(strs []string, addInvalid, returnOnInvalid bool) ([]float64, error) {
	vals := make([]float64, 0, len(strs))
	parser := func(str string) (any, error) {
		val, err := strconv.ParseFloat(str, 64)
		return val, err
	}
	rawVals, err := k.doParse(strs, addInvalid, returnOnInvalid, parser)
	if err == nil {
		for _, val := range rawVals {
			vals = append(vals, val.(float64))
		}
	}
	return vals, err
}

// parseInts transforms strings to ints.
func (k *Key) parseInts(strs []string, addInvalid, returnOnInvalid bool) ([]int, error) {
	vals := make([]int, 0, len(strs))
	parser := func(str string) (any, error) {
		val, err := strconv.ParseInt(str, 0, 64)
		return val, err
	}
	rawVals, err := k.doParse(strs, addInvalid, returnOnInvalid, parser)
	if err == nil {
		for _, val := range rawVals {
			vals = append(vals, int(val.(int64)))
		}
	}
	return vals, err
}

// parseInt64s transforms strings to int64s.
func (k *Key) parseInt64s(strs []string, addInvalid, returnOnInvalid bool) ([]int64, error) {
	vals := make([]int64, 0, len(strs))
	parser := func(str string) (any, error) {
		val, err := strconv.ParseInt(str, 0, 64)
		return val, err
	}

	rawVals, err := k.doParse(strs, addInvalid, returnOnInvalid, parser)
	if err == nil {
		for _, val := range rawVals {
			vals = append(vals, val.(int64))
		}
	}
	return vals, err
}

// parseUints transforms strings to uints.
func (k *Key) parseUints(strs []string, addInvalid, returnOnInvalid bool) ([]uint, error) {
	vals := make([]uint, 0, len(strs))
	parser := func(str string) (any, error) {
		val, err := strconv.ParseUint(str, 0, 64)
		return val, err
	}

	rawVals, err := k.doParse(strs, addInvalid, returnOnInvalid, parser)
	if err == nil {
		for _, val := range rawVals {
			vals = append(vals, uint(val.(uint64)))
		}
	}
	return vals, err
}

// parseUint64s transforms strings to uint64s.
func (k *Key) parseUint64s(strs []string, addInvalid, returnOnInvalid bool) ([]uint64, error) {
	vals := make([]uint64, 0, len(strs))
	parser := func(str string) (any, error) {
		val, err := strconv.ParseUint(str, 0, 64)
		return val, err
	}
	rawVals, err := k.doParse(strs, addInvalid, returnOnInvalid, parser)
	if err == nil {
		for _, val := range rawVals {
			vals = append(vals, val.(uint64))
		}
	}
	return vals, err
}

type Parser func(str string) (any, error)

// parseTimesFormat transforms strings to times in given format.
func (k *Key) parseTimesFormat(format string, strs []string, addInvalid, returnOnInvalid bool) ([]time.Time, error) {
	vals := make([]time.Time, 0, len(strs))
	parser := func(str string) (any, error) {
		val, err := time.Parse(format, str)
		return val, err
	}
	rawVals, err := k.doParse(strs, addInvalid, returnOnInvalid, parser)
	if err == nil {
		for _, val := range rawVals {
			vals = append(vals, val.(time.Time))
		}
	}
	return vals, err
}

// doParse transforms strings to different types
func (k *Key) doParse(strs []string, addInvalid, returnOnInvalid bool, parser Parser) ([]any, error) {
	vals := make([]any, 0, len(strs))
	for _, str := range strs {
		val, err := parser(str)
		if err != nil && returnOnInvalid {
			return nil, err
		}
		if err == nil || addInvalid {
			vals = append(vals, val)
		}
	}
	return vals, nil
}

// SetValue changes key value.
func (k *Key) SetValue(v string) {
	k.s.m.mutex.Lock()
	defer k.s.m.mutex.Unlock()

	k.value = v
	k.s.keysHash[k.name] = v
}
