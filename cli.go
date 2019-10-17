package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	tv "tv/lib"

	"github.com/iancoleman/orderedmap"
	"github.com/urfave/cli"
)

var tvVersion string
var versionInfos *orderedmap.OrderedMap

func getfirstKeyFromOrderedMap(o *orderedmap.OrderedMap) string {
	for _, k := range o.Keys() {
		return k
	}

	return ""
}

func readVersionFile(path string) (*orderedmap.OrderedMap, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	o := orderedmap.New()
	err = json.Unmarshal(buf, &o)

	return o, err
}

func writeVersionFile(o *orderedmap.OrderedMap, path string) error {
	buf, err := json.MarshalIndent(&o, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, buf, 0644)

	return err
}

func init() {
	versionInfos, _ = readVersionFile(tv.SemverFilePath)
}

func getTargetApp(c *cli.Context) string {
	targetApp := c.String("target")
	if targetApp == "" {
		targetApp = getfirstKeyFromOrderedMap(versionInfos)
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

		updateTags(c, v)
	} else {
		return fmt.Errorf("cannot find target app: %s", targetApp)
	}

	return nil
}

func updateTags(c *cli.Context, v *tv.Version) error {
	targetApp := getTargetApp(c)
	nextVer := v.GetVersion()
	appsToUpdate := []string{}

	if !c.Bool("all") {
		appsToUpdate = append(appsToUpdate, targetApp)
	} else {
		appsToUpdate = append(appsToUpdate, versionInfos.Keys()...)
	}

	for _, app := range appsToUpdate {
		versionInfos.Set(app, nextVer)
	}

	if !c.Bool("dry-run") {
		err := writeVersionFile(versionInfos, tv.SemverFilePath)
		if err != nil {
			return err
		}
	}

	tags := []string{}
	if c.Bool("pure") {
		tags = append(tags, nextVer)
	} else {
		for _, app := range appsToUpdate {
			tags = append(tags, v.GetTagStr(app))
		}
	}

	if !c.Bool("dry-run") {
		for _, tag := range tags {
			err := tv.TagVersion(tag)
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("git tag generated: %s\n", tags)

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "tv"
	app.Usage = "tag version for ur f** awesome project"
	if tvVersion != "" {
		app.Version = tvVersion
	}
	commandFlags := []cli.Flag{
		cli.StringFlag{Name: "target, t", Usage: "set target app"},
		cli.BoolFlag{Name: "pure, p", Usage: "create tag without app name"},
		cli.BoolFlag{Name: "all, a", Usage: "upgrade version of all apps"},
		cli.BoolFlag{Name: "dry-run", Usage: "do a fake action, won't create real tag"},
	}
	app.Commands = []cli.Command{
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
