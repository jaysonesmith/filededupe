package main

import (
	"errors"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("returns an error if path is empty", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		d := &Dedupe{fs: fs, path: ""}
		expectedErr := errors.New("path is required")

		err := Run(d)

		assert.Equal(t, expectedErr, err)
	})
}
