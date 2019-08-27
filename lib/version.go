package tv

import (
	"fmt"

	"github.com/blang/semver"
)

const SemverFileName = "semver.json"
const SemverFilePath = "./" + SemverFileName

type Version struct {
	v *semver.Version
}

func Make(vStr string) (*Version, error) {
	v, err := semver.Make(vStr)
	if err != nil {
		return nil, err
	}
	return &Version{
		v: &v,
	}, nil
}

func (v *Version) Major() error {
	err := v.v.IncrementMajor()
	v.v.Pre = nil
	return err
}

func (v *Version) Minor() error {
	err := v.v.IncrementMinor()
	v.v.Pre = nil
	return err
}

func (v *Version) Patch() error {
	err := v.v.IncrementPatch()
	v.v.Pre = nil
	return err
}

func (v *Version) Prerelease() error {
	preVersions := v.v.Pre
	err := fmt.Errorf("Prerelease version can not be incremented for %q", v.v.String())

	if len(preVersions) != 2 {
		return err
	}

	if preVersions[0].VersionStr != "alpha" && preVersions[0].VersionStr != "beta" || !preVersions[1].IsNum {
		return err
	}

	preVersionNum := v.v.Pre[1].VersionNum
	v.v.Pre[1].VersionNum = preVersionNum + 1

	return nil
}

func (v *Version) GetTagStr(build string) string {
	v.v.Build = []string{build}
	return v.v.String()
}

func (v *Version) GetVersion() string {
	return v.v.String()
}
