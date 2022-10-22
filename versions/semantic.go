package versions

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

type (
	Canonical  string
	Major      string
	MajorMinor string
	Prerelease string
	Build      string
)

type SerialzedSemVer struct {
	Canonical  Canonical  `json:"canonical"`
	Major      Major      `json:"major"`
	MajorMinor MajorMinor `json:"majorminor"`
	Prerelease Prerelease `json:"prerelease"`
	Build      Build      `json:"build"`
	Source     string     `json:"source"`
}

type SemanticVersion struct {
	Value string
}

func (s *SemanticVersion) IsValid() bool {
	if strings.HasPrefix(s.Value, "v") {
		return semver.IsValid(s.Value)
	}
	return semver.IsValid(fmt.Sprintf("v%s", s.Value))
}

func (s *SemanticVersion) Canonical() Canonical {
	if strings.HasPrefix(s.Value, "v") {
		return Canonical(semver.Canonical(s.Value))
	}
	return Canonical(
		strings.TrimPrefix(
			semver.Canonical(fmt.Sprintf("v%s", s.Value)),
			"v",
		),
	)
}

func (s *SemanticVersion) Major() Major {
	if strings.HasPrefix(s.Value, "v") {
		return Major(semver.Major(s.Value))
	}
	return Major(
		strings.TrimPrefix(
			semver.Major(fmt.Sprintf("v%s", s.Value)),
			"v",
		),
	)
}

func (s *SemanticVersion) MajorMinor() MajorMinor {
	if strings.HasPrefix(s.Value, "v") {
		return MajorMinor(semver.MajorMinor(s.Value))
	}
	return MajorMinor(
		strings.TrimPrefix(
			semver.MajorMinor(fmt.Sprintf("v%s", s.Value)),
			"v",
		),
	)
}

func (s *SemanticVersion) Prerelease() Prerelease {
	if strings.HasPrefix(s.Value, "v") {
		return Prerelease(semver.Prerelease(s.Value))
	}
	return Prerelease(
		strings.TrimPrefix(
			semver.Prerelease(fmt.Sprintf("v%s", s.Value)),
			"v",
		),
	)
}

func (s *SemanticVersion) Build() Build {
	if strings.HasPrefix(s.Value, "v") {
		return Build(semver.Build(s.Value))
	}
	return Build(
		strings.TrimPrefix(
			semver.Build(fmt.Sprintf("v%s", s.Value)),
			"v",
		),
	)
}

func (s SemanticVersion) String() string {
	return fmt.Sprintf("%s", s.Value)
}

func (s SemanticVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		SerialzedSemVer{
			Canonical:  s.Canonical(),
			Major:      s.Major(),
			MajorMinor: s.MajorMinor(),
			Prerelease: s.Prerelease(),
			Build:      s.Build(),
			Source:     s.Value,
		})
}

func (s SemanticVersion) MarshalEVAL() ([]byte, error) {
	var out string

	out = fmt.Sprintf("%s\nexport MAJOR='%s'", out, string(s.Major()))
	out = fmt.Sprintf("%s\nexport MAJORMINOR='%s'", out, string(s.MajorMinor()))
	out = fmt.Sprintf("%s\nexport CANONICAL='%s'", out, string(s.Canonical()))
	out = fmt.Sprintf("%s\nexport PRERELEASE='%s'", out, string(s.Prerelease()))
	out = fmt.Sprintf("%s\nexport BUILD='%s'", out, string(s.Build()))

	return []byte(out), nil
}

func (s *SemanticVersion) Release(rType string) SemanticVersion {
	switch rType {
	case "major":
		major, err := strconv.Atoi(strings.TrimPrefix(string(s.Major()), "v"))
		if err != nil {
			panic(err)
		}
		major += 1
		return SemanticVersion{fmt.Sprintf("v%d.0.0", major)}
	case "minor":
		major, err := strconv.Atoi(strings.TrimPrefix(string(s.Major()), "v"))
		if err != nil {
			panic(err)
		}
		minor, err := strconv.Atoi(strings.TrimPrefix(string(s.MajorMinor()), fmt.Sprintf("v%d.", major)))
		if err != nil {
			panic(err)
		}
		minor += 1
		return SemanticVersion{fmt.Sprintf("v%d.%d.0", major, minor)}
	case "patch":
		major, err := strconv.Atoi(strings.TrimPrefix(string(s.Major()), "v"))
		if err != nil {
			panic(err)
		}
		minor, err := strconv.Atoi(strings.TrimPrefix(string(s.MajorMinor()), fmt.Sprintf("v%d.", major)))
		if err != nil {
			panic(err)
		}
		noPrefix := strings.TrimPrefix(string(s.Canonical()), fmt.Sprintf("v%d.%d.", major, minor))
		patch := strings.TrimSuffix(noPrefix, string(s.Prerelease()))
		newPatch, err := strconv.Atoi(patch)
		if err != nil {
			panic(err)
		}
		newPatch += 1
		return SemanticVersion{fmt.Sprintf("v%d.%d.%d", major, minor, newPatch)}
	default:
		panic(fmt.Sprintf("err: %s unkown release type", rType))
	}
}

type ByVersion []SemanticVersion

func (a ByVersion) Len() int           { return len(a) }
func (a ByVersion) Less(i, j int) bool { return semver.Compare(a[i].Value, a[j].Value) < 0 }
func (a ByVersion) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
