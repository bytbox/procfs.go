package procfs

import (
	"os"
	"regexp"
)

func exists(pathname string) bool {
	_, err := os.Stat(pathname)
	return err != os.ENOENT
}

func isNumeric(s string) bool {
	a, _ := regexp.Match("[0-9]+", []byte(s))
	return a
}

