package versions

import (
	"encoding/json"
	"testing"
)

var (
	ver0 = "v0.1.2-prerelease.0+build.999"
	ver1 = "v0.1.3-prerelease.0+build.999"
	ver2 = "v0.2.2-prerelease.0+build.999"
	sem0 = SemanticVersion{ver0}
	sem1 = SemanticVersion{ver1}
	sem2 = SemanticVersion{ver2}
)

func TestSemanticStruct(t *testing.T) {
	want := ver0
	success := sem0.Value == want
	if !success {
		t.Fatalf(`SemanticVersion{"v0.1.2-prerelease.0+build.999"}.Value == %q, 
compare for %t, failed: %+v`, want, true, sem0.Value)
	}
}
func TestSemanticMethodMajor(t *testing.T) {
	want := "v0"
	success := string(sem0.Major()) == want
	if !success {
		t.Fatalf(`SemanticVersion{"v0.1.2-prerelease.0+build.999"}.Major() == %q, 
compare for %t, failed: %+v`, want, true, sem0.Major())
	}
}
func TestSemanticMethodMajorMinor(t *testing.T) {
	want := "v0.1"
	success := string(sem0.MajorMinor()) == "v0.1"
	if !success {
		t.Fatalf(`SemanticVersion{"v0.1.2-prerelease.0+build.999"}.MajorMinor() == %q,
compare for %t, failed: %+v`, want, true, sem0.MajorMinor())
	}
}
func TestSemanticMethodCanonical(t *testing.T) {
	want := "v0.1.2-prerelease.0"
	success := string(sem0.Canonical()) == want
	if !success {
		t.Fatalf(`SemanticVersion{"v0.1.2-prerelease.0+build.999"}.Canonical() == %q,
compare for %t, failed: %+v`, want, true, sem0.Canonical())
	}
}
func TestSemanticMethodPrerelease(t *testing.T) {
	want := "-prerelease.0"
	success := string(sem0.Prerelease()) == want
	if !success {
		t.Fatalf(`SemanticVersion{"v0.1.2-prerelease.0+build.999"}.Prerelease() == %q,
compare for %t, failed: %+v`, want, true, sem0.Prerelease())
	}
}
func TestSemanticMethodBuild(t *testing.T) {
	want := "+build.999"
	success := string(sem0.Build()) == want
	if !success {
		t.Fatalf(`SemanticVersion{"v0.1.2-prerelease.0+build.999"}.Build() == %q,
compare for %t, failed: %+v`, want, true, sem0.Build())
	}
}

func TestSemanticJson(t *testing.T) {
	want := `{"canonical":"v0.1.2-prerelease.0","major":"v0","majorminor":"v0.1","prerelease":"-prerelease.0","build":"+build.999","source":"v0.1.2-prerelease.0+build.999"}`
	res, err := json.Marshal(sem0)

	success := string(res) == want
	if !success || err != nil {
		t.Fatalf(`json.Marshal(SemanticVersion{"v0.1.2-prerelease.0+build.999"}) == %q,
compare for %t, failed: %+v`, want, true, sem0.Build())
	}
}
func TestSemanticMethodReleaseMajor(t *testing.T) {
	want := "v1.0.0"
	has := string(sem0.Release("major").Value)
	success := want == has
	if !success {
		t.Fatalf(`SemanticVersion{"v0.1.2-prerelease.0+build.999"} == %q,
compare for %t, failed: %+v`, want, true, has)
	}
}
func TestSemanticMethodReleaseMinor(t *testing.T) {
	want := "v0.2.0"
	has := string(sem0.Release("minor").Value)
	success := want == has
	if !success {
		t.Fatalf(`SemanticVersion{"v0.1.2-prerelease.0+build.999"} == %q,
compare for %t, failed: %+v`, want, true, has)
	}
}
func TestSemanticMethodReleasePatch(t *testing.T) {
	want := "v0.1.3"
	has := string(sem0.Release("patch").Value)
	success := want == has
	if !success {
		t.Fatalf(`SemanticVersion{"v0.1.2-prerelease.0+build.999"} == %q,
compare for %t, failed: %+v`, want, true, has)
	}
}
