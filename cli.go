package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	tv "tv/lib"

	"github.com/iancoleman/orderedmap"
	"github.com/urfave/cli"
)

var version string
var versionInfos *orderedmap.OrderedMap

func getFirstKeyFromOrderedMap(o *orderedmap.OrderedMap) string {
	for _, k := range o.Keys() {
		return k
	}

	return ""
}

func init() {
	versionInfos, _ = tv.ReadVersionFile(tv.SemverFilePath)
}

func getTargetApp(c *cli.Context) string {
	targetApp := c.String("target")
	if targetApp == "" {
		targetApp = getFirstKeyFromOrderedMap(versionInfos)
	}
	// if `--all` is set, choose the highest version
	if c.Bool("all") {
		currentVersion := "0.0.0"
		for _, k := range versionInfos.Keys() {
			version, _ := versionInfos.Get(k)
			result, _ := tv.Compare(version.(string), currentVersion)
			if result == 1 {
				currentVersion = version.(string)
				targetApp = k
			}
		}
	}

	return targetApp
}

func doAction(c *cli.Context, action string) error {
	if c.Bool("dry-run") {
		fmt.Println("start dry-run...")
	}

	targetApp := getTargetApp(c)

	if version, ok := versionInfos.Get(targetApp); ok {
		v, err := tv.Make(version.(string))
		if err != nil {
			return err
		}

		in := []reflect.Value{reflect.ValueOf(c.Args())}
		result := reflect.ValueOf(v).MethodByName(action).Call(in)

		if result[0].Interface() != nil {
			return result[0].Interface().(error)
		}

		err = updateTags(c, v)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("cannot find target app: %s", targetApp)
	}

	return nil
}

func updateTags(c *cli.Context, v *tv.Version) error {
	targetApp := getTargetApp(c)
	nextVer := v.String()
	var appsToUpdate []string

	if !c.Bool("all") {
		appsToUpdate = append(appsToUpdate, targetApp)
	} else {
		appsToUpdate = append(appsToUpdate, versionInfos.Keys()...)
	}

	for _, app := range appsToUpdate {
		versionInfos.Set(app, nextVer)
	}

	if !c.Bool("dry-run") {
		err := tv.WriteVersionFile(versionInfos, tv.SemverFilePath)
		if err != nil {
			return err
		}
	}

	var tags []string
	if c.Bool("pure") {
		tags = append(tags, nextVer)
	} else {
		for _, app := range appsToUpdate {
			tags = append(tags, v.GetTagStr(app))
		}
	}

	if !c.Bool("dry-run") {
		if len(tags) == 1 {
			_ = tv.TagVersion(tags[0])
		} else {
			_ = tv.TagVersions(nextVer, tags)
		}
	}

	fmt.Printf("git tag generated: %s\n", tags)

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "tv"
	app.Usage = "tag version for ur f** awesome project"
	if version != "" {
		app.Version = version
	}
	commandFlags := []cli.Flag{
		cli.StringFlag{Name: "target, t", Usage: "set target app"},
		cli.BoolFlag{Name: "pure, p", Usage: "create tag without app name"},
		cli.BoolFlag{Name: "all, a", Usage: "upgrade version of all apps"},
		cli.BoolFlag{Name: "dry-run", Usage: "do a fake action, won't create real tag"},
	}
	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "init tv for ur project",
			Action: func(c *cli.Context) error {
				args := c.Args()
				if len(args) == 0 {
					return errors.New("init need at least one module name")
				}
				if tv.IsFileExist(tv.SemverFilePath) {
					return errors.New("semver.json is exist")
				}
				o := orderedmap.New()
				for _, moduleName := range args {
					o.Set(moduleName, "0.0.0")
				}
				tv.WriteVersionFile(o, tv.SemverFilePath)
				return nil
			},
		},
		{
			Name:  "patch",
			Usage: "patch version, v0.0.1 -> v0.0.2",
			Flags: commandFlags,
			Action: func(c *cli.Context) error {
				return doAction(c, "Patch")
			},
		},
		{
			Name:  "major",
			Usage: "major version, v0.0.1 -> v1.0.1",
			Flags: commandFlags,
			Action: func(c *cli.Context) error {
				return doAction(c, "Major")
			},
		},
		{
			Name:  "minor",
			Usage: "minor version, v0.0.1 -> v0.1.1",
			Flags: commandFlags,
			Action: func(c *cli.Context) error {
				return doAction(c, "Minor")
			},
		},
		{
			Name:  "prerelease",
			Usage: "prerelease version, v0.0.1-alpha.1 -> v0.0.1-alpha.2",
			Flags: commandFlags,
			Action: func(c *cli.Context) error {
				return doAction(c, "Prerelease")
			},
		},
		{
			Name:  "release",
			Usage: "release version, v0.0.1-alpha.1 -> v0.0.1",
			Flags: commandFlags,
			Action: func(c *cli.Context) error {
				return doAction(c, "Release")
			},
		},
		{
			Name:  "version",
			Usage: "set specific version",
			Flags: commandFlags,
			Action: func(c *cli.Context) error {
				return doAction(c, "SpecificVersion")
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
