package util

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/davecgh/go-spew/spew"
)

func BaseP(s string) string {
	return filepath.Base(s)
}

func SanitizeBackticks(s string) string {
	reg, err := regexp.Compile("`")
	if err != nil {
		panic(err.Error())
	}
	s = reg.ReplaceAllString(s, "'")
	return s
}

func FromMapStringKeys(m map[string]interface{}, key string) interface{} {
	return m[key]
}

func Slice(val []interface{}, index int) interface{} {
	return val[index]
}

func Inpect(val interface{}) string {
	return spew.Sdump(val)
}

func UnsnakeCase(name string) string {
	return strings.Replace(name, "_", "", -1)
}

func CamelCase(name string) string {
	in := strings.Split(name, "_")
	if len(in) == 0 {
		return strings.Title(name)
	}
	out := make([]string, 0, len(in))
	for _, word := range in {
		out = append(out, strings.Title(word))
	}
	return strings.TrimSpace(strings.Join(out, ""))
}

func LowerFirst(name string) string {
	for i, v := range name {
		return string(unicode.ToLower(v)) + name[i+1:]
	}
	return ""
}

func LengthOf(params ...interface{}) int {
	return len(params) - 1 // wtf
}

func FirstOf(opts ...string) string {
	for _, opt := range opts {
		if opt != "" {
			return opt
		}
	}
	return ""
}
