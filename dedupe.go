package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"
)

type deduper interface {
	run() error
	existsNotEmpty() error
}

type Dedupe struct {
	fs    afero.Fs
	path  string
	dry   bool
	files []os.FileInfo
}

func (d *Dedupe) existsNotEmpty() error {
	exists, err := afero.DirExists(d.fs, d.path)
	if err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("unable to find path: %q", d.path)
	}

	empty, err := afero.IsEmpty(d.fs, d.path)
	if err != nil {
		return err
	} else if empty {
		return fmt.Errorf("no files found in path: %q", d.path)
	}

	return nil
}

func (d *Dedupe) fileDetails() error {
	details, err := afero.ReadDir(d.fs, d.path)
	d.files = details

	return err
}

func (d *Dedupe) groupMaybeDupes() [][]os.FileInfo {
	var out [][]os.FileInfo

	var comparePos int
	var appendPos int
	for i := 0; i < len(d.files); i++ {
		if i == 0 {
			out = append(out, []os.FileInfo{d.files[i]})
			continue
		}

		if namesSimilar(d.files[comparePos].Name(), d.files[i].Name()) {
			out[appendPos] = append(out[appendPos], d.files[i])
			continue
		}

		comparePos = i
		appendPos = len(out)
		out = append(out, []os.FileInfo{d.files[i]})
	}

	return out
}

func namesSimilar(name1, name2 string) bool {
	split1 := strings.Split(name1, ".")
	split2 := strings.Split(name2, ".")

	return (strings.Contains(split2[0], split1[0]) || strings.Contains(split1[0], split2[0])) && split2[1] == split1[1]
}
