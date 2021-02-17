package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func TestRename(t *testing.T) {
	t.Run("name should match specified regex", func(t *testing.T) {
		name := "4_11_3_1 4x 2013 edit(1)"
		got := getNewName(name)

		if got == "" {
			t.Errorf("didn't get a result but one was expected")
		}

	})
	t.Run("Should return a new name, subtracting all years by one", func(t *testing.T) {
		cases := []struct {
			oldName string
			want    string
		}{
			{"4_11_3_1 4x 2014 edit(1)", "4_11_3_1 4x 2013 edit(1)"},
			{"11_3_1 4x 2019 edit(1)", "11_3_1 4x 2018 edit(1)"},
			{"11_3_1 4x 2019 2019edit(1)", "11_3_1 4x 2018 2018edit(1)"},
			{"11_3_1 4x (1)", "11_3_1 4x (1)"},
			{"", ""},
		}

		for _, test := range cases {
			t.Run(fmt.Sprintf("should rename %q to %q", test.oldName, test.want), func(t *testing.T) {
				got := getNewName(test.oldName)
				if got != test.want {
					t.Errorf("got %q but expected %q", got, test.want)
				}
			})
		}
	})
}

func TestDirectoryWalk(t *testing.T) {
	t.Run("Should process files from bottom to top", func(t *testing.T) {
		appFS := afero.NewMemMapFs()
		afs := afero.Afero{Fs: appFS}
		afs.MkdirAll("src/a/b/c", 0755)
		afs.WriteFile("src/a/b/c/1.png", []byte(""), 0755)
		afs.WriteFile("src/a/b/2.png", []byte(""), 0755)

		got := []string{}

		processElement := func(fs afero.Fs, e string) {
			got = append(got, e)
		}

		travelFS(afs, "", processElement)
		want := []string{
			"src/a/b/c/1.png",
			"src/a/b/c",
			"src/a/b/2.png",
			"src/a/b",
			"src/a",
			"src",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, expected %v", got, want)
		}

	})
}
