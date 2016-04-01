package namer

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	Dentry struct{ Prefix, Destination string }
	Dtab   []*Dentry
)

func (dentry *Dentry) String() string {
	return fmt.Sprintf("%s=>%s", dentry.Prefix, dentry.Destination)
}

var (
	commentRE   *regexp.Regexp = regexp.MustCompile("^\\s*#.*$")
	dentrySepRE *regexp.Regexp = regexp.MustCompile("\\s*;\\s*")
)

// ParseDtab reads a Dtab string into a list of Prefix and Destination pairs.
func ParseDtab(dtabStr string) (Dtab, error) {
	if dtabStr == "" {
		return Dtab([]*Dentry{}), nil
	}
	dentryStrs := dentrySepRE.Split(dtabStr, -1)
	dentries := []*Dentry{}
	for _, dentryStr := range dentryStrs {
		if dentryStr == "" {
			continue
		}
		parts := strings.SplitN(dentryStr, "=>", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid dentry: '%s'", dentryStr)
		}
		dentries = append(dentries, &Dentry{parts[0], parts[1]})
	}
	dtab := Dtab(dentries)
	return dtab, nil
}

func (dtab Dtab) String() string {
	out := ""
	for _, dentry := range dtab {
		out += dentry.String() + ";"
	}
	return out
}

func (dtab Dtab) Pretty() string {
	maxPfxLen := 0
	for _, d := range dtab {
		l := len(d.Prefix)
		if l > maxPfxLen {
			maxPfxLen = l
		}
	}

	str := ""
	for _, d := range dtab {
		arrow := "=>"
		w := maxPfxLen - len(d.Prefix) + 2
		if w != 0 {
			arrowfmt := fmt.Sprintf("%% %ds", w)
			arrow = fmt.Sprintf(arrowfmt, "=>")
		}
		str += fmt.Sprintf("%s  %s %s ;\n", d.Prefix, arrow, d.Destination)
	}
	return str
}
