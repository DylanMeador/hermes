package main

import (
	"fmt"
	"github.com/jonas747/dca"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path"
	"strings"
)

var filePath string

func main() {
	var cmd = &cobra.Command{
		Use:   "encodedca",
		Short: "encodes an audio file to dca format",
		RunE:   run,
	}

	flags := cmd.Flags()
	flags.StringVarP(&filePath, "file", "f", "", "path/to/file.mp3")
	cmd.MarkFlagRequired("file")

	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func run(c *cobra.Command, args []string) error {
	// Encoding a file and saving it to disk
	encodeSession, err := dca.EncodeFile(filePath, dca.StdEncodeOptions)
	if err != nil {
		return err
	}

	// Make sure everything is cleaned up, that for example the encoding process if any issues happened isnt lingering around
	defer encodeSession.Cleanup()

	fileName := strings.TrimSuffix(path.Base(filePath), path.Ext(filePath))
	outputPath := path.Dir(filePath) + "\\" + fileName + ".dca"

	fmt.Println(outputPath)

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(output, encodeSession)
	return err
}
