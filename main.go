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

func main() {
	reader := bufio.NewReader(os.Stdin)
	version, _ := reader.ReadString('\n')
	trimmed := strings.TrimSpace(version)
	semver := SemanticVersion{trimmed}
	fmt.Fprintln(os.Stdout, semver)
}
