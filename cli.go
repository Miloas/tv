package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	tv "tv/lib"

	"github.com/AlecAivazis/survey/v2"
	"github.com/iancoleman/orderedmap"
	"github.com/urfave/cli"
)

var version string
var versionInfos *orderedmap.OrderedMap

func init() {
	versionInfos, _ = tv.ReadVersionFile(tv.SemverFilePath)
}

// parse the cli arguments and return a list of apps that should be updated.
func getTargetApps(c *cli.Context) []string {
	if !c.Bool("all") {
		targetApp := c.String("target")
		if targetApp == "" {
			// user does not input target app. try to be clever:
			// 1. if there is only one app in semver.json, then we simply choose it for the user.
			// 2. if multiple apps, then launch a select panel for the user to choose the app he/she wants to tag.
			if len(versionInfos.Keys()) == 1 {
				targetApp = versionInfos.Keys()[0]
			} else {
				question := []*survey.Question{
					{
						Name: "app",
						Prompt: &survey.Select{
							Message: "choose an app you want to tag:",
							Options: versionInfos.Keys(),
						},
					},
				}

				answer := struct {
					App string
				}{}

				_ = survey.Ask(question, &answer)
				targetApp = answer.App
			}
		}
		return []string{targetApp}
	} else {
		return versionInfos.Keys()
	}
}

func doAction(c *cli.Context, action string) error {
	if c.Bool("dry-run") {
		fmt.Println("start dry-run...")
	}

	appsToUpdate := getTargetApps(c)

	highestVersion := "0.0.0"

	// find the highest version of these apps.
	for _, app := range appsToUpdate {
		if version, ok := versionInfos.Get(app); ok {
			result, _ := tv.Compare(version.(string), highestVersion)
			if result == 1 {
				highestVersion = version.(string)
			}
		} else {
			return fmt.Errorf("cannot find target app: %s", app)
		}
	}

	v, err := tv.Make(highestVersion)
	if err != nil {
		return err
	}

	in := []reflect.Value{reflect.ValueOf(c.Args())}
	result := reflect.ValueOf(v).MethodByName(action).Call(in)

	if result[0].Interface() != nil {
		return result[0].Interface().(error)
	}

	err = updateTags(c, v, appsToUpdate)
	if err != nil {
		return err
	}

	return nil
}

func updateTags(c *cli.Context, v *tv.Version, appsToUpdate []string) error {
	nextVer := v.String()

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
				var o *orderedmap.OrderedMap
				if tv.IsFileExist(tv.SemverFilePath) {
					o = versionInfos
				} else {
					o = orderedmap.New()
				}
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
			Name:  "prepatch",
			Usage: "prepatch version, v0.0.1 -> v0.0.2-alpha.0",
			Flags: commandFlags,
			Action: func(c *cli.Context) error {
				return doAction(c, "Prepatch")
			},
		},
		{
			Name:  "preminor",
			Usage: "preminor version, v0.0.1 -> v0.1.0-alpha.0",
			Flags: commandFlags,
			Action: func(c *cli.Context) error {
				return doAction(c, "Preminor")
			},
		},
		{
			Name:  "premajor",
			Usage: "preminor version, v0.0.1 -> v1.0.0-alpha.0",
			Flags: commandFlags,
			Action: func(c *cli.Context) error {
				return doAction(c, "Premajor")
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
