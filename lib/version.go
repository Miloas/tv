package tv

import (
	"errors"
	"fmt"

	"github.com/Miloas/semver"
	"github.com/urfave/cli"
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

func (v *Version) Major(args cli.Args) error {
	err := checkArgsEmpty(args)
	if err != nil {
		return err
	}

	err = v.v.IncrementMajor()
	v.v.Pre = nil

	return err
}

func (v *Version) Minor(args cli.Args) error {
	err := checkArgsEmpty(args)
	if err != nil {
		return err
	}

	err = v.v.IncrementMinor()
	v.v.Pre = nil

	return err
}

func (v *Version) Patch(args cli.Args) error {
	err := checkArgsEmpty(args)
	if err != nil {
		return err
	}

	err = v.v.IncrementPatch()
	v.v.Pre = nil

	return err
}

func (v *Version) Prerelease(args cli.Args) error {
	err := checkArgsEmpty(args)
	if err != nil {
		return err
	}

	preVersions := v.v.Pre
	err = fmt.Errorf("Prerelease version can not be incremented for %q", v.v.String())

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

func checkArgsEmpty(args cli.Args) error {
	if len(args) > 0 {
		return fmt.Errorf("too many arguments")
	}

	return nil
}
