package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/mod/semver"
)

type SemanticVersion struct {
	value string
}

func (s *SemanticVersion) Canonical() string {
	return semver.Canonical(s.value)
}
func (s *SemanticVersion) Major() string {
	return semver.Major(s.value)
}
func (s *SemanticVersion) MajorMinor() string {
	return semver.MajorMinor(s.value)
}
func (s *SemanticVersion) Prerelease() string {
	return semver.Prerelease(s.value)
}
func (s *SemanticVersion) Build() string {
	return semver.Build(s.value)
}
func (s SemanticVersion) String() string {
	out := map[string]string{
		"canonical":  s.Canonical(),
		"major":      s.Major(),
		"majorminor": s.MajorMinor(),
		"prerelease": s.Prerelease(),
		"build":      s.Build(),
	}
	data, _ := json.Marshal(out)
	return fmt.Sprintf("%s", data)
}

func output(data string) {
	trimmed := strings.TrimSpace(data)
	semver := SemanticVersion{trimmed}

	fmt.Fprintln(os.Stdout, semver)
}

func input() string {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	var version string

	if fi.Mode()&os.ModeNamedPipe == 0 {
		version = os.Getenv("CI_COMMIT_TAG")
	} else {
		reader := bufio.NewReader(os.Stdin)
		version, _ = reader.ReadString('\n')
	}
	return version
}

func main() {
	output(input())

}
