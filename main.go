package main

import (
  "os"
  "github.com/codegangsta/cli"
)

const (
	appName = "hranoprovod"
	appUsage = "Lifestyle tracker"
	appVersion = "2.0.0"
	appAuthor = "aquilax"
	appEmail = "aquilax@gmail.com"
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = appUsage
	app.Version = appVersion
	app.Author = appAuthor
	app.Email = appEmail
	app.Commands = []cli.Command{
		{
			Name: "register",
			ShortName: "reg",
			Usage: "Shows the register",
			Action: func(c *cli.Context) {
				handleExit(NewHranoprovod().register())
			},
		},
		{
			Name: "search",
			Usage: "Search for food",
			Action: func(c *cli.Context) {
				handleExit(NewHranoprovod().search(c.Args().First()))
			},
		},
	}
	app.Run(os.Args)
}

func handleExit(err error) {

}