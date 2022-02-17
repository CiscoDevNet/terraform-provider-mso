package mso

import (
	"fmt"
)

const version = 1

func toStringList(configured interface{}) []string {
	vs := make([]string, 0, 1)
	val, ok := configured.(string)
	if ok && val != "" {
		vs = append(vs, val)
	}
	return vs
}

func makeTestVariable(s string) string {
	return fmt.Sprintf("acctest_%s", s)
}
