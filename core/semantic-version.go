package core

import (
	"fmt"
	"regexp"
	"strconv"
)

type SemanticVersion struct {
	Prefix string
	Major  uint32
	Minor  uint32
	Patch  uint32
	Suffix string
}

func (v SemanticVersion) String() string {
	return fmt.Sprintf("%s%d.%d.%d%s", v.Prefix, v.Major, v.Minor, v.Patch, v.Suffix)
}

var semverExp = regexp.MustCompile(`(\D+)?([\d]+)\.([\d]+)\.([\d]+)(\D+)?`)

func ParseSemanticVersion(s string) (SemanticVersion, error) {
	m := semverExp.FindStringSubmatch(s)
	if len(m) != 6 {
		return SemanticVersion{}, fmt.Errorf("%s does not follow the format '<prefix><major>.<minor>.<patch><suffix> where major, minor, and patch must be digits'", s)
	}

	major, err := strconv.Atoi(m[2])
	if err != nil {
		return SemanticVersion{}, err
	}
	minor, err := strconv.Atoi(m[3])
	if err != nil {
		return SemanticVersion{}, err
	}
	patch, err := strconv.Atoi(m[4])
	if err != nil {
		return SemanticVersion{}, err
	}

	return SemanticVersion{
		Prefix: m[1],
		Major:  uint32(major),
		Minor:  uint32(minor),
		Patch:  uint32(patch),
		Suffix: m[5],
	}, nil
}

type ByVersion []SemanticVersion

func (v ByVersion) Len() int      { return len(v) }
func (v ByVersion) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v ByVersion) Less(i, j int) bool {
	a := v[i]
	b := v[j]
	//check in the order 'major', 'minor', 'patch', 'suffix'. Don't care about the prefix
	if a.Major < b.Major {
		return true
	} else if a.Major == b.Major {
		if a.Minor < b.Minor {
			return true
		} else if a.Minor == b.Minor {
			if a.Patch < b.Patch {
				return true
			} else if a.Patch == b.Patch {
				if a.Suffix < b.Suffix {
					return true
				}
			}
		}
	}
	return false
}
