package main

import (
	"errors"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/extern/sector-storage/fr32"
	"github.com/filecoin-project/lotus/extern/sector-storage/partialfile"
	"github.com/filecoin-project/lotus/extern/sector-storage/storiface"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"io"
	"log"
	"os"
)

func unpadded(i storiface.PaddedByteIndex) storiface.UnpaddedByteIndex {
	return storiface.UnpaddedByteIndex(abi.PaddedPieceSize(i).Unpadded())
}

func convertPiece(input string, output string, offset storiface.PaddedByteIndex, length abi.UnpaddedPieceSize) error {
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
	log.Print("open partial file")
	pf, err := partialfile.OpenPartialFile(maxPieceSize, input)
	if err != nil {
		return err
	}
	log.Print("check allocation")
	ok, err := pf.HasAllocated(unpadded(offset), length)
	if err != nil {
		_ = pf.Close()
		return err
	}

	if !ok {
		_ = pf.Close()
		return errors.New("not allocated")
	}

	log.Print("setup reader")
	f, err := pf.Reader(offset, length.Padded())
	if err != nil {
		_ = pf.Close()
		return xerrors.Errorf("getting partial file reader: %w", err)
	}

	log.Print("setup unpad reader")
	upr, err := fr32.NewUnpadReader(f, length.Padded())
	if err != nil {
		return xerrors.Errorf("creating unpadded reader: %w", err)
	}

	log.Print("open output file")
	writer, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
	defer writer.Close()

	log.Print("copy stream")
	_, err = io.CopyN(writer, upr, int64(length))
	return err
}

func main() {
	app := &cli.App{
		Name:  "unpad",
		Usage: "unpad the unsealed sector into car file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "Read unsealed sector file from `INPUT`",
				Required: true,
			},
			&cli.Uint64Flag{
				Name:  "offset",
				Usage: "Start position of the deal",
				Value: 0,
			},
			&cli.Uint64Flag{
				Name:     "length",
				Usage:    "length of the deal",
				Value:    uint64(32) * 1024 * 1024 * 1024,
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Write car file to `OUTPUT`",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			input := c.String("input")
			output := c.String("output")
			offset := c.Uint64("offset")
			length := c.Uint64("length")
			err := convertPiece(input, output, storiface.PaddedByteIndex(offset), abi.UnpaddedPieceSize(length))
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
