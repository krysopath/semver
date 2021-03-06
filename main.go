package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

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
			Canonical:  Canonical(s.Canonical()),
			Major:      Major(s.Major()),
			MajorMinor: MajorMinor(s.MajorMinor()),
			Prerelease: Prerelease(s.Prerelease()),
			Build:      Build(s.Build()),
			Source:     s.value,
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

	out = fmt.Sprintf("%s\nexport MAJOR='%s'", out, shellescape.Quote(s.Major()))
	out = fmt.Sprintf("%s\nexport MAJORMINOR='%s'", out, shellescape.Quote(s.MajorMinor()))
	out = fmt.Sprintf("%s\nexport CANONICAL='%s'", out, shellescape.Quote(s.Canonical()))
	out = fmt.Sprintf("%s\nexport PRERELEASE='%s'", out, shellescape.Quote(s.Prerelease()))
	out = fmt.Sprintf("%s\nexport BUILD='%s'", out, shellescape.Quote(s.Build()))

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
		strings.TrimSpace(string(out)),
	)
}

func toSemVersion(in []string) []SemanticVersion {
	var out []SemanticVersion
	for _, v := range in {
		out = append(out, SemanticVersion{v})
	}
	return out
}

func cleanVersions(versions []SemanticVersion) []SemanticVersion {
	var cleaned []SemanticVersion
	for _, v := range versions {
		if len(v.Canonical()) > 0 {
			cleaned = append(cleaned, v)
		} else {
			fmt.Fprintf(os.Stderr, "WARN: truncated input %s\n", v.value)
		}
	}
	return cleaned
}
func outputSorted(data []string) {
	sorted := cleanVersions(toSemVersion(data))
	sort.Sort(ByVersion(sorted))
	out, _ := json.Marshal(sorted)
	fmt.Fprintln(os.Stdout, string(out))
}

func outputOrdered(data []string) {
	sorted := cleanVersions(toSemVersion(data))
	out, _ := json.Marshal(sorted)
	fmt.Fprintln(os.Stdout, string(out))
}

func input() []string {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	var trimmed []string

	if fi.Mode()&os.ModeNamedPipe == 0 {
		data := os.Getenv("CI_COMMIT_TAG")
		trimmed = strings.Fields(data)
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			for _, v := range strings.Fields(scanner.Text()) {
				trimmed = append(trimmed, v)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}
	return trimmed
}

func main() {
	flag.Parse()
	data := input()
	if len(data) > 1 {
		if *isSorting {
			outputSorted(data)
		} else {
			outputOrdered(data)
		}
	} else if len(data) == 1 {
		outputSingle(data[0])
	} else {
		log.Fatalf("WARN: no input %s", flag.Args())
		flag.Usage()
	}
}
