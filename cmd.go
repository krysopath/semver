package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type SemanticVersion struct {
	value string
}

func (s *SemanticVersion) Canonical() string {
	return Canonical(s.value)
}
func (s *SemanticVersion) Major() string {
	return Major(s.value)
}
func (s *SemanticVersion) MajorMinor() string {
	return MajorMinor(s.value)
}
func (s *SemanticVersion) Prerelease() string {
	return Prerelease(s.value)
}
func (s *SemanticVersion) Build() string {
	return Build(s.value)
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
