package ini

import (
	"os"
	"regexp"
	"strings"
)

// Maximum allowed depth when recursively substituing variable names.
const depthValues = 99

var (
	// Variable regexp pattern: %(variable)s
	varPattern = regexp.MustCompile(`%\(([^)]+)\)s`)

	// Environment variable regexp pattern:
	//   ${variable}
	//   ${variable||default}
	//   ${variable??default}
	envPattern = regexp.MustCompile(`\$\{([^}]+)\}`)
)

type ValueTransformer func(m *Manager, s *Section, k *Key) string

// transformValue takes a key and transforms to its final string.
func transformValue(k *Key) string {
	val := transformCustom(k)
	val = transformReference(k, val)
	val = transformEnvironment(val)
	return val
}

func transformCustom(k *Key) string {
	if k.s.m.options.Transformer != nil {
		return k.s.m.options.Transformer(k.s.m, k.s, k)
	}
	return k.value
}

func transformReference(k *Key, val string) string {
	// Fail-fast if no indicate char found for recursive value
	if !strings.Contains(val, "%") {
		return val
	}

	for range depthValues {
		vr := varPattern.FindString(val)
		if len(vr) == 0 {
			break
		}

		// Take off leading '%(' and trailing ')s'.
		noption := vr[2 : len(vr)-2]

		// Search in the same section.
		// If not found or found the key itself, then search again in default section.
		nk, err := k.s.GetKey(noption)
		if err != nil || k == nk {
			nk, _ = k.s.m.Section("").GetKey(noption)
			if nk == nil {
				// Stop when no results found in the default section,
				// and returns the value as-is.
				break
			}
		}

		// Substitute by new value and take off leading '%(' and trailing ')s'.
		val = strings.Replace(val, vr, nk.value, -1)
	}

	return val
}

func transformEnvironment(val string) string {
	// Fail-fast if no indicate char found for recursive value
	if !strings.Contains(val, "$") {
		return val
	}

	for range depthValues {
		vr := envPattern.FindString(val)
		if len(vr) == 0 {
			break
		}

		// Take off leading '${' and trailing '}'.
		noption := vr[2 : len(vr)-1]

		// Split the option into key and default value.
		// If no default value found, then use empty string.
		var force bool
		parts := strings.SplitN(noption, "??", 2)
		if len(parts) == 1 {
			force = true
			parts = strings.SplitN(noption, "||", 2)
		}

		// Get the key and default value.
		key := strings.TrimSpace(parts[0])
		def := ""
		if len(parts) == 2 {
			def = trimQuote(strings.TrimSpace(parts[1]))
		}

		// Get the value from environment.
		// If no value found, then use default value.
		value, ok := os.LookupEnv(key)
		if !ok || (value == "" && force) {
			value = def
		}

		// Substitute by new value and take off leading '${' and trailing '}'.
		val = strings.Replace(val, vr, value, -1)
	}

	return val
}

func trimQuote(s string) string {
	if hasSurroundedQuote(s, '\'') || hasSurroundedQuote(s, '"') {
		return s[1 : len(s)-1]
	}
	return s
}
