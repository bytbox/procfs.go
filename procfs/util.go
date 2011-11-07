package procfs

import (
	"bytes"
	"os"
	"regexp"
)

func splitNull(b []byte) []string {
	null := []byte{0}
	rb := bytes.Split(b, null)
	r := []string{}
	for _, x := range rb[:len(rb)-1] {
		r = append(r, string(x))
	}
	return r
}

func exists(pathname string) bool {
	_, err := os.Stat(pathname)
	return err != os.ENOENT
}

func isNumeric(s string) bool {
	a, _ := regexp.Match("^[0-9]+$", []byte(s))
	return a
}

