package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	tv "tv/lib"

	"github.com/urfave/cli"
)

var versionInfos map[string]string

func readVersionFile(path string) (map[string]string, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m map[string]string
	err = json.Unmarshal(buf, &m)
	return m, err
}

func writeVersionFile(m map[string]string, path string) error {
	buf, err := json.MarshalIndent(&m, "", "  ")
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
	if !tv.IsClean() {
		return fmt.Errorf("workspace not clean")
	}
	build := c.String("build")
	if version, ok := versionInfos[build]; ok {
		v, err := tv.Make(version)
		if err != nil {
			return err
		}
		reflect.ValueOf(v).MethodByName(action).Call([]reflect.Value{})
		versionInfos[build] = v.GetVersion()
		err = writeVersionFile(versionInfos, tv.SemverFilePath)
		if err != nil {
			return err
		}
		err = tv.TagVersion(v.GetTagStr(build))
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
