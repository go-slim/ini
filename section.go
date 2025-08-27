package ini

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type Section struct {
	m        *Manager
	name     string
	keys     map[string]*Key
	keyList  []string
	keysHash map[string]string
	Comment  string
}

func newSection(m *Manager, name string) *Section {
	return &Section{
		m:        m,
		name:     name,
		keys:     make(map[string]*Key),
		keyList:  make([]string, 0),
		keysHash: make(map[string]string),
	}
}

// Name returns name of Section.
func (s *Section) Name() string {
	return s.name
}

// Parent returns the parent section.
func (s *Section) Parent() (*Section, bool) {
	if i := strings.LastIndex(s.name, s.m.options.ChildSectionDelimiter); i > -1 {
		return s.m.Section(s.name[:i]), true
	}
	return nil, false
}

// NewKey creates a new key to given section.
func (s *Section) NewKey(name, value string) *Key {
	if s.m.options.Insensitive || s.m.options.InsensitiveKeys {
		name = strings.ToLower(name)
	}

	s.m.mutex.Lock()
	defer s.m.mutex.Unlock()

	if slices.Contains(s.keyList, name) {
		return s.keys[name]
	}

	s.keyList = append(s.keyList, name)
	s.keys[name] = newKey(s, name, value)
	s.keysHash[name] = value

	return s.keys[name]
}

func (s *Section) NewBooleanKey(name string) *Key {
	key := s.NewKey(name, "true")
	key.isBooleanType = true
	return key
}

// GetKey returns key in section by given name.
func (s *Section) GetKey(name string) (*Key, error) {
	s.m.mutex.RLock()
	if s.m.options.Insensitive || s.m.options.InsensitiveKeys {
		name = strings.ToLower(name)
	}
	key := s.keys[name]
	s.m.mutex.RUnlock()

	if key == nil {
		// Check if it is a child-section.
		sname := s.name
		for {
			if i := strings.LastIndex(sname, s.m.options.ChildSectionDelimiter); i > -1 {
				sname = sname[:i]
				sec, err := s.m.GetSection(sname)
				if err != nil {
					continue
				}
				return sec.GetKey(name)
			}
			break
		}
		return nil, fmt.Errorf("error when getting key of section %q: key %q not exists", s.name, name)
	}
	return key, nil
}

// HasKey returns true if section contains a key with given name.
func (s *Section) HasKey(name string) bool {
	key, _ := s.GetKey(name)
	return key != nil
}

// HasValue returns true if section contains given raw value.
func (s *Section) HasValue(value string) bool {
	s.m.mutex.RLock()
	defer s.m.mutex.RUnlock()
	for _, k := range s.keys {
		if value == k.value {
			return true
		}
	}
	return false
}

// Key assumes named Key exists in section and returns a zero-value when not.
func (s *Section) Key(name string) *Key {
	key, err := s.GetKey(name)
	if err != nil {
		// It's OK here because the only possible error is empty key name,
		// but if it's empty, this piece of code won't be executed.
		key = newKey(s, name, "")
	}
	return key
}

// Keys returns list of keys of section.
func (s *Section) Keys() []*Key {
	keys := make([]*Key, len(s.keyList))
	for i := range s.keyList {
		keys[i] = s.Key(s.keyList[i])
	}
	return keys
}

// String returns string representation of value.
func (s *Section) String(name string) string {
	return s.Key(name).String()
}

// Validate accepts a validate function which can
// return modifed result as key value.
func (s *Section) Validate(name string, fn func(string) string) string {
	return s.Key(name).Validate(fn)
}

// Bool returns bool type value.
func (s *Section) Bool(name string) (bool, error) {
	return s.Key(name).Bool()
}

// Float64 returns float64 type value.
func (s *Section) Float64(name string) (float64, error) {
	return s.Key(name).Float64()
}

// Int returns int type value.
func (s *Section) Int(name string) (int, error) {
	return s.Key(name).Int()
}

// Int64 returns int64 type value.
func (s *Section) Int64(name string) (int64, error) {
	return s.Key(name).Int64()
}

// Uint returns uint type valued.
func (s *Section) Uint(name string) (uint, error) {
	return s.Key(name).Uint()
}

// Uint64 returns uint64 type value.
func (s *Section) Uint64(name string) (uint64, error) {
	return s.Key(name).Uint64()
}

// Duration returns time.Duration type value.
func (s *Section) Duration(name string) (time.Duration, error) {
	return s.Key(name).Duration()
}

// TimeFormat parses with given format and returns time.Time type value.
func (s *Section) TimeFormat(name string, format string) (time.Time, error) {
	return s.Key(name).TimeFormat(format)
}

// Time parses with RFC3339 format and returns time.Time type value.
func (s *Section) Time(name string) (time.Time, error) {
	return s.Key(name).Time()
}

// MustString returns default value if key value is empty.
func (s *Section) MustString(name string, defaultVal ...string) string {
	if len(defaultVal) > 0 {
		return s.Key(name).MustString(defaultVal[0])
	}
	return s.Key(name).String()
}

// MustBool always returns value without error,
// it returns false if error occurs.
func (s *Section) MustBool(name string, defaultVal ...bool) bool {
	return s.Key(name).MustBool(defaultVal...)
}

// MustFloat64 always returns value without error,
// it returns 0.0 if error occurs.
func (s *Section) MustFloat64(name string, defaultVal ...float64) float64 {
	return s.Key(name).MustFloat64(defaultVal...)
}

// MustInt always returns value without error,
// it returns 0 if error occurs.
func (s *Section) MustInt(name string, defaultVal ...int) int {
	return s.Key(name).MustInt(defaultVal...)
}

// MustInt64 always returns value without error,
// it returns 0 if error occurs.
func (s *Section) MustInt64(name string, defaultVal ...int64) int64 {
	return s.Key(name).MustInt64(defaultVal...)
}

// MustUint always returns value without error,
// it returns 0 if error occurs.
func (s *Section) MustUint(name string, defaultVal ...uint) uint {
	return s.Key(name).MustUint(defaultVal...)
}

// MustUint64 always returns value without error,
// it returns 0 if error occurs.
func (s *Section) MustUint64(name string, defaultVal ...uint64) uint64 {
	return s.Key(name).MustUint64(defaultVal...)
}

// MustDuration always returns value without error,
// it returns zero value if error occurs.
func (s *Section) MustDuration(name string, defaultVal ...time.Duration) time.Duration {
	return s.Key(name).MustDuration(defaultVal...)
}

// MustTimeFormat always parses with given format and returns value without error,
// it returns zero value if error occurs.
func (s *Section) MustTimeFormat(name string, format string, defaultVal ...time.Time) time.Time {
	return s.Key(name).MustTimeFormat(format, defaultVal...)
}

// MustTime always parses with RFC3339 format and returns value without error,
// it returns zero value if error occurs.
func (s *Section) MustTime(name string, defaultVal ...time.Time) time.Time {
	return s.Key(name).MustTime(defaultVal...)
}

// In always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (s *Section) In(name string, defaultVal string, candidates []string) string {
	return s.Key(name).In(defaultVal, candidates)
}

// InFloat64 always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (s *Section) InFloat64(name string, defaultVal float64, candidates []float64) float64 {
	return s.Key(name).InFloat64(defaultVal, candidates)
}

// InInt always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (s *Section) InInt(name string, defaultVal int, candidates []int) int {
	return s.Key(name).InInt(defaultVal, candidates)
}

// InInt64 always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (s *Section) InInt64(name string, defaultVal int64, candidates []int64) int64 {
	return s.Key(name).InInt64(defaultVal, candidates)
}

// InUint always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (s *Section) InUint(name string, defaultVal uint, candidates []uint) uint {
	return s.Key(name).InUint(defaultVal, candidates)
}

// InUint64 always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (s *Section) InUint64(name string, defaultVal uint64, candidates []uint64) uint64 {
	return s.Key(name).InUint64(defaultVal, candidates)
}

// InTimeFormat always parses with given format and returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (s *Section) InTimeFormat(name string, format string, defaultVal time.Time, candidates []time.Time) time.Time {
	return s.Key(name).InTimeFormat(format, defaultVal, candidates)
}

// InTime always parses with RFC3339 format and returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (s *Section) InTime(name string, defaultVal time.Time, candidates []time.Time) time.Time {
	return s.Key(name).InTime(defaultVal, candidates)
}

// RangeFloat64 checks if value is in given range inclusively,
// and returns default value if it's not.
func (s *Section) RangeFloat64(name string, defaultVal, min, max float64) float64 {
	return s.Key(name).RangeFloat64(defaultVal, min, max)
}

// RangeInt checks if value is in given range inclusively,
// and returns default value if it's not.
func (s *Section) RangeInt(name string, defaultVal, min, max int) int {
	return s.Key(name).RangeInt(defaultVal, min, max)
}

// RangeInt64 checks if value is in given range inclusively,
// and returns default value if it's not.
func (s *Section) RangeInt64(name string, defaultVal, min, max int64) int64 {
	return s.Key(name).RangeInt64(defaultVal, min, max)
}

// RangeTimeFormat checks if value with given format is in given range inclusively,
// and returns default value if it's not.
func (s *Section) RangeTimeFormat(name string, format string, defaultVal, min, max time.Time) time.Time {
	return s.Key(name).RangeTimeFormat(format, defaultVal, min, max)
}

// RangeTime checks if value with RFC3339 format is in given range inclusively,
// and returns default value if it's not.
func (s *Section) RangeTime(name string, defaultVal, min, max time.Time) time.Time {
	return s.Key(name).RangeTime(defaultVal, min, max)
}

// Strings returns list of string divided by given delimiter.
func (s *Section) Strings(name string, delim string) []string {
	return s.Key(name).Strings(delim)
}

// Float64s returns list of float64 divided by given delimiter. Any invalid input will be treated as zero value.
func (s *Section) Float64s(name string, delim string) []float64 {
	return s.Key(name).Float64s(delim)
}

// Ints returns list of int divided by given delimiter. Any invalid input will be treated as zero value.
func (s *Section) Ints(name string, delim string) []int {
	return s.Key(name).Ints(delim)
}

// Int64s returns list of int64 divided by given delimiter. Any invalid input will be treated as zero value.
func (s *Section) Int64s(name string, delim string) []int64 {
	return s.Key(name).Int64s(delim)
}

// Uints returns list of uint divided by given delimiter. Any invalid input will be treated as zero value.
func (s *Section) Uints(name string, delim string) []uint {
	return s.Key(name).Uints(delim)
}

// Uint64s returns list of uint64 divided by given delimiter. Any invalid input will be treated as zero value.
func (s *Section) Uint64s(name string, delim string) []uint64 {
	return s.Key(name).Uint64s(delim)
}

// Bools returns list of bool divided by given delimiter. Any invalid input will be treated as zero value.
func (s *Section) Bools(name string, delim string) []bool {
	return s.Key(name).Bools(delim)
}

// TimesFormat parses with given format and returns list of time.Time divided by given delimiter.
// Any invalid input will be treated as zero value (0001-01-01 00:00:00 +0000 UTC).
func (s *Section) TimesFormat(name string, format, delim string) []time.Time {
	return s.Key(name).TimesFormat(format, delim)
}

// Times parses with RFC3339 format and returns list of time.Time divided by given delimiter.
// Any invalid input will be treated as zero value (0001-01-01 00:00:00 +0000 UTC).
func (s *Section) Times(name string, delim string) []time.Time {
	return s.Key(name).Times(delim)
}

// ValidFloat64s returns list of float64 divided by given delimiter. If some value is not float, then
// it will not be included to result list.
func (s *Section) ValidFloat64s(name string, delim string) []float64 {
	return s.Key(name).ValidFloat64s(delim)
}

// ValidInts returns list of int divided by given delimiter. If some value is not integer, then it will
// not be included to result list.
func (s *Section) ValidInts(name string, delim string) []int {
	return s.Key(name).ValidInts(delim)
}

// ValidInt64s returns list of int64 divided by given delimiter. If some value is not 64-bit integer,
// then it will not be included to result list.
func (s *Section) ValidInt64s(name string, delim string) []int64 {
	return s.Key(name).ValidInt64s(delim)
}

// ValidUints returns list of uint divided by given delimiter. If some value is not unsigned integer,
// then it will not be included to result list.
func (s *Section) ValidUints(name string, delim string) []uint {
	return s.Key(name).ValidUints(delim)
}

// ValidUint64s returns list of uint64 divided by given delimiter. If some value is not 64-bit unsigned
// integer, then it will not be included to result list.
func (s *Section) ValidUint64s(name string, delim string) []uint64 {
	return s.Key(name).ValidUint64s(delim)
}

// ValidBools returns list of bool divided by given delimiter. If some value is not 64-bit unsigned
// integer, then it will not be included to result list.
func (s *Section) ValidBools(name string, delim string) []bool {
	return s.Key(name).ValidBools(delim)
}

// ValidTimesFormat parses with given format and returns list of time.Time divided by given delimiter.
func (s *Section) ValidTimesFormat(name string, format, delim string) []time.Time {
	return s.Key(name).ValidTimesFormat(format, delim)
}

// ValidTimes parses with RFC3339 format and returns list of time.Time divided by given delimiter.
func (s *Section) ValidTimes(name string, delim string) []time.Time {
	return s.Key(name).ValidTimes(delim)
}

// StrictFloat64s returns list of float64 divided by given delimiter or error on first invalid input.
func (s *Section) StrictFloat64s(name string, delim string) ([]float64, error) {
	return s.Key(name).StrictFloat64s(delim)
}

// StrictInts returns list of int divided by given delimiter or error on first invalid input.
func (s *Section) StrictInts(name string, delim string) ([]int, error) {
	return s.Key(name).StrictInts(delim)
}

// StrictInt64s returns list of int64 divided by given delimiter or error on first invalid input.
func (s *Section) StrictInt64s(name string, delim string) ([]int64, error) {
	return s.Key(name).StrictInt64s(delim)
}

// StrictUints returns list of uint divided by given delimiter or error on first invalid input.
func (s *Section) StrictUints(name string, delim string) ([]uint, error) {
	return s.Key(name).StrictUints(delim)
}

// StrictUint64s returns list of uint64 divided by given delimiter or error on first invalid input.
func (s *Section) StrictUint64s(name string, delim string) ([]uint64, error) {
	return s.Key(name).StrictUint64s(delim)
}

// StrictBools returns list of bool divided by given delimiter or error on first invalid input.
func (s *Section) StrictBools(name string, delim string) ([]bool, error) {
	return s.Key(name).StrictBools(delim)
}

// StrictTimesFormat parses with given format and returns list of time.Time divided by given delimiter
// or error on first invalid input.
func (s *Section) StrictTimesFormat(name string, format, delim string) ([]time.Time, error) {
	return s.Key(name).StrictTimesFormat(format, delim)
}

// StrictTimes parses with RFC3339 format and returns list of time.Time divided by given delimiter
// or error on first invalid input.
func (s *Section) StrictTimes(name string, delim string) ([]time.Time, error) {
	return s.Key(name).StrictTimes(delim)
}
