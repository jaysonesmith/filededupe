package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.StringP("path", "p", "", "path to directory to dedupe (required)")
	pflag.BoolP("dry", "d", false, "whether or not to actually make changes on run")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

func main() {
	d := &Dedupe{
		fs:   afero.NewOsFs(),
		path: viper.GetString("path"),
		dry:  viper.GetBool("dry"),
	}

	if err := Run(d); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Run(d *Dedupe) error {
	if d.path == "" {
		return errors.New("path is required")
	}

	err := d.existsNotEmpty()
	if err != nil {
		return err
	}

	err = d.fileDetails()
	if err != nil {
		return err
	}

	_ = d.groupMaybeDupes()

	return nil
}
