package util

import (
	"strings"
	"unicode"
)

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
