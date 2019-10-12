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

func doAction(c *cli.Context, action string) error {
	if c.Bool("dry-run") {
		fmt.Println("Start dry-run...")
	}
	build := c.String("build")
	if build == "" {
		build = getfirstKeyFromOrderedMap(versionInfos)
	}
	if version, ok := versionInfos.Get(build); ok {
		v, err := tv.Make(version.(string))
		if err != nil {
			return err
		}
		reflect.ValueOf(v).MethodByName(action).Call([]reflect.Value{})
		versionInfos.Set(build, v.GetVersion())
		tag := v.GetVersion()
		if !c.Bool("clear") {
			tag = v.GetTagStr(build)
		}
		fmt.Println("Generating git tag:", tag)
		if c.Bool("dry-run") {
			return nil
		}
		err = writeVersionFile(versionInfos, tv.SemverFilePath)
		if err != nil {
			return err
		}
		err = tv.TagVersion(tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "tv"
	app.Usage = "tag version for ur f** awesome project"
	flags := []cli.Flag{
		cli.StringFlag{Name: "build, b"},
		cli.BoolFlag{Name: "clear, c"},
		cli.BoolFlag{Name: "dry-run"},
	}
	app.Commands = []cli.Command{
		{
			Name:  "patch",
			Usage: "patch version, v0.0.1 -> v0.0.2",
			Flags: flags,
			Action: func(c *cli.Context) error {
				return doAction(c, "Patch")
			},
		},
		{
			Name:  "major",
			Usage: "major version, v0.0.1 -> v1.0.1",
			Flags: flags,
			Action: func(c *cli.Context) error {
				return doAction(c, "Major")
			},
		},
		{
			Name:  "minor",
			Usage: "minor version, v0.0.1 -> v0.1.1",
			Flags: flags,
			Action: func(c *cli.Context) error {
				return doAction(c, "Minor")
			},
		},
		{
			Name:  "prerelease",
			Usage: "prerelease version, v0.0.1-alpha.1 -> v0.0.1-alpha.2",
			Flags: flags,
			Action: func(c *cli.Context) error {
				return doAction(c, "Prerelease")
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
