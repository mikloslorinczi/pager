package util

import (
	"strings"
)

var HunToEng = map[string]string{
	"á": "a",
	"Á": "a",
	"é": "e",
	"É": "é",
	"í": "i",
	"Í": "i",
	"ó": "o",
	"Ó": "o",
	"ö": "o",
	"Ö": "o",
	"ő": "o",
	"Ő": "o",
	"ú": "u",
	"Ú": "u",
	"ü": "u",
	"Ü": "u",
	"ű": "u",
	"Ű": "u",
}

func Englisher(in string) string {
	out := ""
	for _, ch := range strings.ToLower(in) {
		if val, ok := HunToEng[string(ch)]; ok {
			out += val
		} else {
			out += string(ch)
		}
	}
	return out
}
