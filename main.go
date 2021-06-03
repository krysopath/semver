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

func main() {
	reader := bufio.NewReader(os.Stdin)
	version, _ := reader.ReadString('\n')
	trimmed := strings.TrimSpace(version)
	semver := SemanticVersion{trimmed}

	out := map[string]string{
		"canonical":  semver.Canonical(),
		"major":      semver.Major(),
		"majorminor": semver.MajorMinor(),
		"prerelease": semver.Prerelease(),
		"build":      semver.Build(),
	}
	data, _ := json.Marshal(out)
	fmt.Println(string(data))
}
