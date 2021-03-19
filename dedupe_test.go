package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const (
	testPath string = "testdata/foo"
)

func TestExistsNotEmpty(t *testing.T) {
	t.Run("returns an error if path not found", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		d := &Dedupe{fs: fs, path: testPath}
		expectedErr := fmt.Errorf("unable to find path: \"%s\"", testPath)

		err := d.existsNotEmpty()

		assert.Equal(t, expectedErr, err)
	})

	t.Run("returns an error if path contains no file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		fs.MkdirAll(testPath, 0755)
		d := &Dedupe{fs: fs, path: testPath}
		expectedErr := fmt.Errorf("no files found in path: \"%s\"", testPath)

		err := d.existsNotEmpty()

		assert.Equal(t, expectedErr, err)
	})

	t.Run("returns nil if path exists and is not empty", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		fs.MkdirAll(testPath, 0755)
		afero.WriteFile(fs, fmt.Sprintf("%s/file.png", testPath), []byte("foo"), 0644)
		d := &Dedupe{fs: fs, path: testPath}

		err := d.existsNotEmpty()

		assert.Nil(t, err)
	})
}

func TestFileDetails(t *testing.T) {
	t.Run("returns data of files found", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		fs.MkdirAll(testPath, 0755)
		afero.WriteFile(fs, fmt.Sprintf("%s/file1.png", testPath), []byte("foo"), 0644)
		afero.WriteFile(fs, fmt.Sprintf("%s/file2.png", testPath), []byte("foo"), 0644)
		d := &Dedupe{fs: fs, path: testPath}

		err := d.fileDetails()

		assert.Nil(t, err)
		for i := 0; i < 2; i++ {
			assert.Equal(t, fmt.Sprintf("file%d.png", i+1), d.files[i].Name())
		}
	})
}

func TestGroupMaybeDupes(t *testing.T) {
	testCases := []struct {
		name          string
		testFiles     []os.FileInfo
		expectedNames [][]string
	}{
		{
			name: "does not group files that do not have similar names",
			testFiles: []os.FileInfo{
				testFileInfo{name: "264.CR2"},
				testFileInfo{name: "265.CR2"},
			},
			expectedNames: [][]string{
				{"264.CR2"}, {"265.CR2"},
			},
		},
		{
			name: "groups files of names like '264.CR2'",
			testFiles: []os.FileInfo{
				testFileInfo{name: "264.CR2"},
				testFileInfo{name: "264-2.CR2"},
			},
			expectedNames: [][]string{
				{"264.CR2", "264-2.CR2"},
			},
		},
		{
			name: "groups files of names like 'DSC_2093.NEF'",
			testFiles: []os.FileInfo{
				testFileInfo{name: "DSC_2093.NEF"},
				testFileInfo{name: "DSC_2093-001.NEF"},
				testFileInfo{name: "DSC_2093-002.NEF"},
			},
			expectedNames: [][]string{
				{"DSC_2093.NEF", "DSC_2093-001.NEF", "DSC_2093-002.NEF"},
			},
		},
		{
			name: "groups different files of different name styles",
			testFiles: []os.FileInfo{
				testFileInfo{name: "P3190152.ORF"},
				testFileInfo{name: "P3190152 (2).ORF"},
				testFileInfo{name: "IMG_1766.CR2"},
				testFileInfo{name: "IMG_1766_2.CR2"},
			},
			expectedNames: [][]string{
				{"P3190152.ORF", "P3190152 (2).ORF"},
				{"IMG_1766.CR2", "IMG_1766_2.CR2"},
			},
		},
		{
			name: "groups files of similar names regardless if in reverse order",
			testFiles: []os.FileInfo{
				testFileInfo{name: "264-2.CR2"},
				testFileInfo{name: "264.CR2"},
			},
			expectedNames: [][]string{
				{"264-2.CR2", "264.CR2"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := &Dedupe{files: tc.testFiles}

			actual := d.groupMaybeDupes()

			// only check contents if initial length check passes to cut down
			// on panics over index out of range errors
			if assert.Len(t, actual, len(tc.expectedNames)) {

				for gi, group := range tc.expectedNames {
					if assert.Len(t, group, len(tc.expectedNames[gi])) {

						for di, expected := range group {
							assert.Equal(t, expected, actual[gi][di].Name())
						}
					}
				}
			}
		})
	}
}

type testFileInfo struct {
	name string
	data []byte
}

func (fi testFileInfo) Name() string       { return fi.name }
func (fi testFileInfo) Size() int64        { return int64(len(fi.data)) }
func (fi testFileInfo) Mode() os.FileMode  { return 0444 }        // Read for all
func (fi testFileInfo) ModTime() time.Time { return time.Time{} } // Return whatever you want
func (fi testFileInfo) IsDir() bool        { return false }
func (fi testFileInfo) Sys() interface{}   { return nil }
