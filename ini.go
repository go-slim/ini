package ini

import "sync"

// Options contains all customized options used for load data source(s).
type Options struct {
	// Loose indicates whether the parser should ignore nonexistent files or return error.
	Loose bool
	// Insensitive indicates whether the parser forces all section and key names to lowercase.
	Insensitive bool
	// InsensitiveSections indicates whether the parser forces all section to lowercase.
	InsensitiveSections bool
	// InsensitiveKeys indicates whether the parser forces all key names to lowercase.
	InsensitiveKeys bool
	// IgnoreContinuation indicates whether to ignore continuation lines while parsing.
	IgnoreContinuation bool
	// IgnoreInlineComment indicates whether to ignore comments at the end of value and treat it as part of value.
	IgnoreInlineComment bool
	// AllowBooleanKeys indicates whether to allow boolean type keys or treat as value is missing.
	// This type of keys are mostly used in my.cnf.
	AllowBooleanKeys bool
	// AllowPythonMultilineValues indicates whether to allow Python-like multi-line values.
	// Docs: https://docs.python.org/3/library/configparser.html#supported-ini-file-structure
	// Relevant quote:  Values can also span multiple lines, as long as they are indented deeper
	// than the first line of the value.
	AllowPythonMultilineValues bool
	// SpaceBeforeInlineComment indicates whether to allow comment symbols (\# and \;) inside value.
	// Docs: https://docs.python.org/2/library/configparser.html
	// Quote: Comments may appear on their own in an otherwise empty line, or may be entered in lines holding values or section names.
	// In the latter case, they need to be preceded by a whitespace character to be recognized as a comment.
	SpaceBeforeInlineComment bool
	// UnescapeValueDoubleQuotes indicates whether to unescape double quotes inside value to regular format
	// when value is surrounded by double quotes, e.g. key="a \"value\"" => key=a "value"
	UnescapeValueDoubleQuotes bool
	// UnescapeValueCommentSymbols indicates to unescape comment symbols (\# and \;) inside value to regular format
	// when value is NOT surrounded by any quotes.
	// Note: UNSTABLE, behavior might change to only unescape inside double quotes but may noy necessary at all.
	UnescapeValueCommentSymbols bool
	// KeyValueDelimiters is the sequence of delimiters that are used to separate key and value. By default, it is "=:".
	KeyValueDelimiters string
	// ChildSectionDelimiter is the delimiter that is used to separate child sections. By default, it is ".".
	ChildSectionDelimiter string
	// PreserveSurroundedQuote indicates whether to preserve surrounded quote (single and double quotes).
	PreserveSurroundedQuote bool
	// DebugFunc is called to collect debug information (currently only useful to debug parsing Python-style multiline values).
	DebugFunc func(message string)
	// ReaderBufferSize is the buffer size of the reader in bytes.
	ReaderBufferSize int
	// AllowNonUniqueSections indicates whether to allow sections with the same name multiple times.
	AllowNonUniqueSections bool
	// AllowDuplicateShadowValues indicates whether values for shadowed keys should be deduplicated.
	AllowDuplicateShadowValues bool
	// Mutex Should make things safe, but sometimes doesn't matter.
	Mutex Mutex
	// ValueMapper represents a mapping function for values
	ValueMapper func(m *Manager, s *Section, k *Key) string
	Transformer ValueTransformer
}

type Mutex interface {
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

func New(opts Options) *Manager {
	if len(opts.KeyValueDelimiters) == 0 {
		opts.KeyValueDelimiters = "=:"
	}
	if len(opts.ChildSectionDelimiter) == 0 {
		opts.ChildSectionDelimiter = "."
	}
	if opts.Mutex == nil {
		opts.Mutex = &sync.RWMutex{}
	}
	return &Manager{
		options:  opts,
		sections: make(map[string]*Section),
		mutex:    opts.Mutex,
	}
}
