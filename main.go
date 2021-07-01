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
	isSorting  = flag.Bool("sort", false, "sort strings from $* inputs with semver algo")
	emitFormat = flag.String("format", "json", "output format on stdout")
)

type SerialzedSemVer struct {
	Canonical  string `json:"canonical"`
	Major      string `json:"major"`
	MajorMinor string `json:"majorminor"`
	Prerelease string `json:"prerelease"`
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
	return fmt.Sprintf("%s", s.value)
}

func (s SemanticVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		SerialzedSemVer{
			Canonical:  s.Canonical(),
			Major:      s.Major(),
			MajorMinor: s.MajorMinor(),
			Prerelease: s.Prerelease(),
			Build:      s.Build(),
			Source:     s.value,
		})
}

func (s SemanticVersion) MarshalEVAL() ([]byte, error) {
	var out string

	out = fmt.Sprintf("%s\nexport %s=%s", out, "MAJOR", s.Major())
	out = fmt.Sprintf("%s\nexport %s=%s", out, "MAJORMINOR", s.MajorMinor())
	out = fmt.Sprintf("%s\nexport %s=%s", out, "CANONICAL", s.Canonical())

	return []byte(out), nil
}

type ByVersion []SemanticVersion

func (a ByVersion) Len() int           { return len(a) }
func (a ByVersion) Less(i, j int) bool { return semver.Compare(a[i].value, a[j].value) < 0 }
func (a ByVersion) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func outputSingle(data string) {
	ver := SemanticVersion{data}

	var out []byte

	switch *emitFormat {
	case "json":
		out, _ = json.Marshal(ver)
	case "eval":
		out, _ = ver.MarshalEVAL()
	}

	fmt.Fprintln(
		os.Stdout,
		strings.TrimSpace(shellescape.Quote(string(out))),
	)
}

func outputSorted(data []string) {
	var sorted []SemanticVersion
	for _, v := range data {
		sorted = append(sorted, SemanticVersion{v})
	}
	sort.Sort(ByVersion(sorted))
	out, _ := json.Marshal(sorted)
	fmt.Fprintln(os.Stdout, string(out))
}

func outputOrdered(data []string) {
	var sorted []SemanticVersion
	for _, v := range data {
		sorted = append(sorted, SemanticVersion{v})
	}
	out, _ := json.Marshal(sorted)
	fmt.Fprintln(os.Stdout, string(out))
}

func input() []string {
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
	trimmed := strings.Fields(data)

	return trimmed
}

func main() {
	flag.Parse()
	data := input()
	if len(data) > 2 {
		if *isSorting {
			outputSorted(data)
		} else {
			outputOrdered(data)
		}
	} else if len(data) == 1 {
		outputSingle(data[0])
	} else {
		flag.Usage()
	}
}
