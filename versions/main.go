package versions

import (
	"encoding/json"
	"flag"
	"fmt"

	"golang.org/x/mod/semver"
	"gopkg.in/alessio/shellescape.v1"
)

var (
	isSorting  = flag.Bool("sort", false, "sort strings from $* inputs with semver algo")
	emitFormat = flag.String("format", "json", "output format on stdout")
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

func (s *SemanticVersion) Canonical() Canonical   { return Canonical(semver.Canonical(s.Value)) }
func (s *SemanticVersion) Major() Major           { return Major(semver.Major(s.Value)) }
func (s *SemanticVersion) MajorMinor() MajorMinor { return MajorMinor(semver.MajorMinor(s.Value)) }
func (s *SemanticVersion) Prerelease() Prerelease { return Prerelease(semver.Prerelease(s.Value)) }
func (s *SemanticVersion) Build() Build           { return Build(semver.Build(s.Value)) }

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

	//	var properties []func() string = []func() string{
	//		s.Major, s.MajorMinor, s.Canonical, s.Prerelease, s.Build,
	//	}
	//
	//	for _, fn := range properties {
	//		out = fmt.Sprintf("%s\nexport MAJOR='%s'", out, shellescape.Quote(s.Major()))
	//	}

	out = fmt.Sprintf("%s\nexport MAJOR='%s'", out, shellescape.Quote(string(s.Major())))
	out = fmt.Sprintf("%s\nexport MAJORMINOR='%s'", out, shellescape.Quote(string(s.MajorMinor())))
	out = fmt.Sprintf("%s\nexport CANONICAL='%s'", out, shellescape.Quote(string(s.Canonical())))
	out = fmt.Sprintf("%s\nexport PRERELEASE='%s'", out, shellescape.Quote(string(s.Prerelease())))
	out = fmt.Sprintf("%s\nexport BUILD='%s'", out, shellescape.Quote(string(s.Build())))

	return []byte(out), nil
}

type ByVersion []SemanticVersion

func (a ByVersion) Len() int           { return len(a) }
func (a ByVersion) Less(i, j int) bool { return semver.Compare(a[i].Value, a[j].Value) < 0 }
func (a ByVersion) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
