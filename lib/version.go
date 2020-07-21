package tv

import (
	"errors"
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/urfave/cli"
)

// SemverFileName : store version info
const SemverFileName = "semver.json"

// SemverFilePath : path of SemverFile
const SemverFilePath = "./" + SemverFileName

// Version : struct to store semver.Version
type Version struct {
	v *semver.Version
}

// Make : create Version from version string
func Make(vStr string) (*Version, error) {
	v, err := semver.Make(vStr)
	if err != nil {
		return nil, err
	}

	return &Version{
		v: &v,
	}, nil
}

// Compare : using to compare two version string
func Compare(vStr1 string, vStr2 string) (int, error) {
	v1, err := semver.Make(vStr1)
	if err != nil {
		return 0, err
	}

	v2, err := semver.Make(vStr2)
	if err != nil {
		return 0, err
	}

	return v1.Compare(v2), nil
}

// SpecificVersion : tag specific version
func (v *Version) SpecificVersion(args cli.Args) error {
	if len(args) != 1 {
		return errors.New("unacceptable arguments for specific version")
	}

	ver, err := semver.Make(args[0])
	if err != nil {
		return err
	}

	v.v = &ver

	return err
}

// Major : increment Major version
func (v *Version) Major(args cli.Args) error {
	err := checkArgsEmpty(args)
	if err != nil {
		return err
	}

	err = v.v.IncrementMajor()
	v.v.Pre = nil

	return err
}

// Minor : increment Minor version
func (v *Version) Minor(args cli.Args) error {
	err := checkArgsEmpty(args)
	if err != nil {
		return err
	}

	err = v.v.IncrementMinor()
	v.v.Pre = nil

	return err
}

// Patch : using to increment Patch version
func (v *Version) Patch(args cli.Args) error {
	err := checkArgsEmpty(args)
	if err != nil {
		return err
	}

	err = v.v.IncrementPatch()
	v.v.Pre = nil

	return err
}

// Release : Release version
func (v *Version) Release(args cli.Args) error {
	err := checkArgsEmpty(args)
	if err != nil {
		return err
	}

	preVersions := v.v.Pre
	if len(preVersions) != 2 {
		return fmt.Errorf("%q is not a prerelease version", v.v.String())
	}
	v.v.Pre = nil
	return nil
}

// Prerelease : using to increment Prerelease version
func (v *Version) Prerelease(args cli.Args) error {
	err := checkArgsEmpty(args)
	if err != nil {
		return err
	}

	preVersions := v.v.Pre
	err = fmt.Errorf("prerelease version can not be incremented for %q", v.v.String())

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

// Prepatch: using to increment Prepatch version
func (v *Version) Prepatch(args cli.Args) error {
	err := checkArgsEmpty(args)
	if err != nil {
		return err
	}

	preVersions := v.v.Pre
	if len(preVersions) != 0 {
		return fmt.Errorf("Prepatch version can not be applied for a prerelease version %q", v.v.String())
	}
	v.v.Patch += 1
	v.v.Pre = []semver.PRVersion{{
		VersionStr: "alpha",
		IsNum:      false,
	}, {
		VersionNum: 0,
		IsNum:      true,
	}}
	return nil
}

// Preminor: using to increment Preminor version
func (v *Version) Preminor(args cli.Args) error {
	err := checkArgsEmpty(args)
	if err != nil {
		return err
	}

	preVersions := v.v.Pre
	if len(preVersions) != 0 {
		return fmt.Errorf("Preminor version can not be applied for a prerelease version %q", v.v.String())
	}
	v.v.Minor += 1
	v.v.Patch = 0
	v.v.Pre = []semver.PRVersion{{
		VersionStr: "alpha",
		IsNum:      false,
	}, {
		VersionNum: 0,
		IsNum:      true,
	}}
	return nil
}

// GetTagStr : get tag str with build info
func (v *Version) GetTagStr(build string) string {
	v.v.Build = []string{build}
	return v.v.String()
}

// String : Version -> string
func (v *Version) String() string {
	return v.v.String()
}

func checkArgsEmpty(args cli.Args) error {
	if len(args) > 0 {
		return errors.New("too many arguments")
	}

	return nil
}
