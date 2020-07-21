package tv

import (
	"testing"

	"github.com/blang/semver/v4"
)

func prStr(s string) semver.PRVersion {
	return semver.PRVersion{VersionStr: s, VersionNum: 0, IsNum: false}
}

func prNum(i uint64) semver.PRVersion {
	return semver.PRVersion{VersionStr: "", VersionNum: i, IsNum: true}
}

func TestMake(t *testing.T) {
	ver, err := Make("0.1.0-alpha.1+tv")
	if err != nil {
		t.Error("Make Version from string error")
	}
	if ver.v.String() != (semver.Version{
		Major: 0,
		Minor: 1,
		Patch: 0,
		Pre:   []semver.PRVersion{prStr("alpha"), prNum(1)},
		Build: []string{"tv"},
	}).String() {
		t.Error("Make get a wrong Version result")
	}
}

func TestSpecificVersion(t *testing.T) {
	ver, _ := Make("1.2.3")
	err := ver.SpecificVersion([]string{"2.5.0"})
	if err != nil {
		t.Error("Do SpecificVersion error")
	}

	if ver.v.String() != (semver.Version{
		Major: 2,
		Minor: 5,
		Patch: 0,
		Pre:   []semver.PRVersion{},
		Build: []string{},
	}).String() {
		t.Error("SpecificVersion get a wrong result")
	}

	if ver.GetTagStr("tv") != "2.5.0+tv" {
		t.Error("get wrong tag str result")
	}
}

func TestMajor(t *testing.T) {
	ver, _ := Make("1.2.3")
	err := ver.Major([]string{})
	if err != nil {
		t.Error("Do Major error")
	}

	if ver.v.String() != (semver.Version{
		Major: 2,
		Minor: 0,
		Patch: 0,
		Pre:   []semver.PRVersion{},
		Build: []string{},
	}).String() {
		t.Error("Major get a wrong result")
	}

	if ver.GetTagStr("tv") != "2.0.0+tv" {
		t.Error("get wrong tag str result")
	}
}

func TestMinor(t *testing.T) {
	ver, _ := Make("1.2.3")
	err := ver.Minor([]string{})
	if err != nil {
		t.Error("Do Minor error")
	}

	if ver.v.String() != (semver.Version{
		Major: 1,
		Minor: 3,
		Patch: 0,
		Pre:   []semver.PRVersion{},
		Build: []string{},
	}).String() {
		t.Error("Minor get a wrong result")
	}

	if ver.GetTagStr("tv") != "1.3.0+tv" {
		t.Error("get wrong tag str result")
	}
}

func TestPatch(t *testing.T) {
	ver, _ := Make("1.2.3")
	err := ver.Patch([]string{})
	if err != nil {
		t.Error("Do Patch error")
	}

	if ver.v.String() != (semver.Version{
		Major: 1,
		Minor: 2,
		Patch: 4,
		Pre:   []semver.PRVersion{},
		Build: []string{},
	}).String() {
		t.Error("Patch get a wrong result")
	}

	if ver.GetTagStr("tv") != "1.2.4+tv" {
		t.Error("get wrong tag str result")
	}
}

func TestPrerelease(t *testing.T) {
	ver, _ := Make("1.2.3-alpha.1")
	err := ver.Prerelease([]string{})
	if err != nil {
		t.Error("Do Prerelease error")
	}

	if ver.v.String() != (semver.Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
		Pre:   []semver.PRVersion{prStr("alpha"), prNum(2)},
		Build: []string{},
	}).String() {
		t.Error("Prerelease get a wrong result")
	}

	if ver.GetTagStr("tv") != "1.2.3-alpha.2+tv" {
		t.Error("get wrong tag str result")
	}
}

func TestRelease(t *testing.T) {
	ver, _ := Make("1.2.3-alpha.1")
	err := ver.Release([]string{})
	if err != nil {
		t.Error("Do Prerelease error")
	}

	if ver.v.String() != (semver.Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
		Pre:   []semver.PRVersion{},
		Build: []string{},
	}).String() {
		t.Error("Release get a wrong result")
	}

	if ver.GetTagStr("tv") != "1.2.3+tv" {
		t.Error("get wrong tag str result")
	}
}
