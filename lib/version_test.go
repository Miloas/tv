package tv

import (
	"testing"

	"github.com/blang/semver"
)

func prstr(s string) semver.PRVersion {
	return semver.PRVersion{s, 0, false}
}

func prnum(i uint64) semver.PRVersion {
	return semver.PRVersion{"", i, true}
}

func TestMake(t *testing.T) {
	ver, err := Make("0.1.0-alpha.1+tv")
	if err != nil {
		t.Error("Make Version from string error")
	}
	if ver.v.String() != (semver.Version{
		0,
		1,
		0,
		[]semver.PRVersion{prstr("alpha"), prnum(1)},
		[]string{"tv"},
	}).String() {
		t.Error("Make get a wrong Version result")
	}
}

func TestMajor(t *testing.T) {
	ver, _ := Make("1.2.3")
	err := ver.Major()
	if err != nil {
		t.Error("Do Major error")
	}

	if ver.v.String() != (semver.Version{
		2,
		0,
		0,
		[]semver.PRVersion{},
		[]string{},
	}).String() {
		t.Error("Major get a wrong result")
	}

	if ver.GetTagStr("tv") != "2.0.0+tv" {
		t.Error("get wrong tag str result")
	}
}

func TestMinor(t *testing.T) {
	ver, _ := Make("1.2.3")
	err := ver.Minor()
	if err != nil {
		t.Error("Do Minor error")
	}

	if ver.v.String() != (semver.Version{
		1,
		3,
		0,
		[]semver.PRVersion{},
		[]string{},
	}).String() {
		t.Error("Minor get a wrong result")
	}

	if ver.GetTagStr("tv") != "1.3.0+tv" {
		t.Error("get wrong tag str result")
	}
}

func TestPatch(t *testing.T) {
	ver, _ := Make("1.2.3")
	err := ver.Patch()
	if err != nil {
		t.Error("Do Patch error")
	}

	if ver.v.String() != (semver.Version{
		1,
		2,
		4,
		[]semver.PRVersion{},
		[]string{},
	}).String() {
		t.Error("Patch get a wrong result")
	}

	if ver.GetTagStr("tv") != "1.2.4+tv" {
		t.Error("get wrong tag str result")
	}
}

func TestPrerelease(t *testing.T) {
	ver, _ := Make("1.2.3-alpha.1")
	err := ver.Prerelease()
	if err != nil {
		t.Error("Do Prerelease error")
	}

	if ver.v.String() != (semver.Version{
		1,
		2,
		3,
		[]semver.PRVersion{prstr("alpha"), prnum(2)},
		[]string{},
	}).String() {
		t.Error("Prerelease get a wrong result")
	}

	if ver.GetTagStr("tv") != "1.2.3-alpha.2+tv" {
		t.Error("get wrong tag str result")
	}
}
