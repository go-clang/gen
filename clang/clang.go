package clang

import (
	"fmt"
	"regexp"
	"strconv"
)

// Version represents a Clang version.
type Version struct {
	Major    int
	Minor    int
	Subminor int
}

// ParseVersion parses Version from s.
func ParseVersion(s []byte) *Version {
	m := regexp.MustCompile(`^(\d+)\.(\d+)(?:\.(\d+))?`).FindSubmatch(s)
	if m == nil {
		return nil
	}

	var err error
	var v Version

	if v.Major, err = strconv.Atoi(string(m[1])); err != nil {
		return nil
	}
	if v.Minor, err = strconv.Atoi(string(m[2])); err != nil {
		return nil
	}
	if len(m[3]) != 0 {
		if v.Subminor, err = strconv.Atoi(string(m[3])); err != nil {
			return nil
		}
	} else {
		v.Subminor = 0
	}

	return &v
}

// String returns a string representation of the Version.
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Subminor)
}

// StringMinor returns a string representation of the minor Version.
func (v Version) StringMinor() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}
