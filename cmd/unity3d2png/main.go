package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mix3/unity3d2png"
	"github.com/pkg/errors"
)

func main() {
	app := cli.NewApp()
	app.Name = "unity3d2png"
	app.UsageText = "unity3d2png <file>"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "disunity",
			Value: "disunity.jar",
		},
		cli.StringFlag{
			Name:  "imagemagick",
			Value: "convert",
		},
		cli.StringFlag{
			Name:  "java",
			Value: "java",
		},
	}
	app.Action = func(c *cli.Context) error {
		if len(c.Args()) <= 0 {
			return fmt.Errorf("Parameters error")
		}

		s := unity3d2png.Service{
			Java:     c.String("java"),
			Disunity: c.String("disunity"),
			Convert:  c.String("imagemagick"),
		}

		_, err := s.Extract(c.Args()[0])
		if err != nil {
			return errors.Wrap(err, "failed extract unity3d")
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(-1)
	}
}
