package catdog_env

import (
	"os"
	"strings"

	"github.com/pubgo/xerror"
)

var prefix string

func Prefix(p string) {
	prefix = p
}

func Set(key, value string) error {
	if prefix != "" {
		key = prefix + "_" + key
	}

	return xerror.Wrap(os.Setenv(strings.ToUpper(key), value))
}

func Get(val *string, names ...string) {
	for _, name := range names {
		nm := name
		if prefix != "" {
			nm = prefix + "_" + nm
		}

		env, ok := os.LookupEnv(strings.ToUpper(nm))
		env = strings.TrimSpace(env)
		if ok && env != "" {
			*val = env
		}
	}
}

// Expand
// replaces ${var} or $var in the string according to the values
// of the current environment variables. References to undefined
// variables are replaced by the empty string.
func Expand(data string) string {
	return os.Expand(data, func(s string) string {
		if prefix != "" {
			return prefix + "_" + s
		}
		return s
	})
}

func Clear() {
	os.Clearenv()
}

func Lookup(key string) (string, bool) {
	if prefix != "" {
		key = prefix + "_" + key
	}

	return os.LookupEnv(key)
}

func Unsetenv(key string) error {
	if prefix != "" {
		key = prefix + "_" + key
	}

	return os.Unsetenv(key)
}
