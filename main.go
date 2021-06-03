package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
)

var (
	isSorting = flag.Bool("sort", false, "sort strings from $* inputs with semver algo")
)

type SerialzedSemVer struct {
	Canonical  string `json:"canonical"`
	Major      string `json:"major"`
	MajorMinor string `json:"majorminor"`
	PreRelease string `json:"prerelease"`
	Build      string `json:"build"`
	Source     string `json:"source"`
}

type SemanticVersion struct {
	value string
}

func (s *SemanticVersion) Canonical() string  { return semver.Canonical(s.value) }
func (s *SemanticVersion) Major() string      { return semver.Major(s.value) }
func (s *SemanticVersion) MajorMinor() string { return semver.MajorMinor(s.value) }
func (s *SemanticVersion) Prerelease() string { return semver.Prerelease(s.value) }
func (s *SemanticVersion) Build() string      { return semver.Build(s.value) }

func (s SemanticVersion) String() string {
	asJson, _ := s.MarshalJSON()
	data, _ := json.Marshal(asJson)
	return fmt.Sprintf("%s", data)
}

func (s SemanticVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		SerialzedSemVer{
			Canonical:  s.Canonical(),
			Major:      s.Major(),
			MajorMinor: s.MajorMinor(),
			PreRelease: s.Prerelease(),
			Build:      s.Build(),
			Source:     s.value,
		})
}

type ByVersion []SemanticVersion

func (a ByVersion) Len() int           { return len(a) }
func (a ByVersion) Less(i, j int) bool { return semver.Compare(a[i].value, a[j].value) < 0 }
func (a ByVersion) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func outputSingle(data string) {
	ver := SemanticVersion{data}
	out, _ := json.Marshal(ver)
	fmt.Fprintln(
		os.Stdout,
		string(out),
	)
}

func outputSorted(data string) {
	vers := strings.Split(data, " ")
	var sorted []SemanticVersion
	for _, v := range vers {
		sorted = append(sorted, SemanticVersion{v})
	}
	sort.Sort(ByVersion(sorted))

	out, _ := json.Marshal(sorted)
	fmt.Fprintln(os.Stdout, string(out))
}

func input() string {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	var data string

	if fi.Mode()&os.ModeNamedPipe == 0 {
		data = os.Getenv("CI_COMMIT_TAG")
	} else {
		reader := bufio.NewReader(os.Stdin)
		data, _ = reader.ReadString('\n')
	}
	trimmed := strings.TrimSpace(data)
	return trimmed
}

func main() {
	flag.Parse()

	if *isSorting {
		outputSorted(input())
	} else {
		outputSingle(input())
	}
}
