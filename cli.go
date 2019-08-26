package main

import (
  "fmt"
  "log"
  "os"

  "github.com/urfave/cli"
)


func main() {
  app := cli.NewApp()
	app.Name = "tv"
	app.Usage = "tag version for ur f** awesome project"
  app.Commands = []cli.Command{
    {
      Name:  "patch",
      Usage: "patch version, v0.0.1 -> v0.0.2",
      Flags: []cli.Flag{
        cli.StringFlag{Name: "build, b"},
      },
      Action: func(c *cli.Context) error {
        fmt.Println("build:", c.String("build"))
				return nil
      },
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
