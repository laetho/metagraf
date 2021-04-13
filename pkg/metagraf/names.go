package metagraf

import (
	"github.com/blang/semver"
	"strconv"
	"strings"
)

// Returns the basename of a metagraf component based on conventions
// around semver version.
func (mg MetaGraf) Name(oname string, version string) string {
	var custver string
	var custname string

	if len(oname) > 0 {
		custname = oname
	} else {
		custname = strings.ToLower(mg.Metadata.Name)
	}

	if len(version) > 0 {
		sv, err := semver.Parse(version)
		if err != nil {
			custver = "-" + version
		} else {
			custver = "v" + strings.ToLower(strconv.FormatUint(sv.Major, 10))
		}
	} else {
		sv, err := semver.Parse(mg.Spec.Version)
		if err != nil {
			custver = "-" + mg.Spec.Version
		} else {
			custver = "v" + strings.ToLower(strconv.FormatUint(sv.Major, 10))
		}
	}
	return strings.ToLower(custname + custver)
}
