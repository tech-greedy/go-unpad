package main

import (
	"github.com/filecoin-project/go-state-types/abi"
	"log"
	"os"
	"github.com/urfave/cli/v2"
	"github.com/filecoin-project/lotus/extern/sector-storage/fr32"
	"github.com/filecoin-project/lotus/extern/sector-storage/partialfile"
)

func convertPiece(input string, output string, offset int64, length int64) error {
	stat, err := os.Stat(input)
	if err != nil {
		return err
	}
	size := stat.Size()
	ssize := abi.SectorSize(34359738368)
	if size >= 68719476736 {
		ssize = abi.SectorSize(68719476736)
	}
	maxPieceSize := abi.PaddedPieceSize(ssize)
	pf, err := partialfile.OpenPartialFile(maxPieceSize, path.Unsealed)
}

func main() {
	app := &cli.App{
		Name:  "unpad",
		Usage: "unpad the unsealed sector into car file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "Read unsealed sector file from `INPUT`",
				Required: true,
			},
			&cli.Int64Flag{
				Name:  "offset",
				Usage: "Start position of the deal",
				Value: 0,
			},
			&cli.Int64Flag{
				Name:  "length",
				Usage: "length of the deal",
				Value: int64(32) * 1024 * 1024 * 1024,
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Write car file to `OUTPUT`",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			input := c.String("input")
			output := c.String("output")
			offset := c.Int64("offset")
			length := c.Int64("length")
			err := convertPiece(input, output, offset, length)
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
