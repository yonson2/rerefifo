package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/spf13/afero"
)

const regex string = "\\d{4}"

//GetNewName receives a name and returns a new name or an empty string
func getNewName(name string) string {
	r, _ := regexp.Compile(regex)
	allNameIndex := r.FindAllStringIndex(name, -1)
	if allNameIndex != nil {
		result := name
		for _, matchIndex := range allNameIndex {
			beggining := result[:matchIndex[0]]
			end := result[matchIndex[1]:]

			match := result[matchIndex[0]:matchIndex[1]]
			desiredOutcome := subtractOneYear(match)

			result = fmt.Sprintf("%v%v%v", beggining, desiredOutcome, end)
		}
		return result
	}
	return name
}

func subtractOneYear(match string) string {
	//right now, there is no need to check for errors as we are using a fixed regex pattern
	year, _ := strconv.Atoi(match)
	return strconv.Itoa(year - 1)
}

func travelFS(afs afero.Afero, root string, processElement func(afero.Fs, string)) {
	//TODO: cover error case
	elements, _ := afs.ReadDir(root)
	var dirs []os.FileInfo
	var files []os.FileInfo
	for _, e := range elements {
		if e.IsDir() {
			dirs = append(dirs, e)
		} else {
			files = append(files, e)
		}
	}
	//walk all sub-directories first
	for _, dir := range dirs {
		//TODO: cover error case
		travelFS(afs, filepath.Join(root, dir.Name()), processElement)
	}
	//then process all files in this directory.
	for _, file := range files {
		processElement(afs.Fs, filepath.Join(root, file.Name()))
	}

	if root != "" {
		processElement(afs.Fs, root)
	}

}

func renameElement(fs afero.Fs, name string) {
	newName := getNewName(name)
	if newName != name {
		//TODO: check for errors here
		fmt.Printf("renaming %s to %s\n", name, newName)
		fs.Rename(name, newName)
	}
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get current directory")
	}
	fmt.Println(wd)

	bp := afero.NewBasePathFs(afero.NewOsFs(), wd)
	fs := afero.Afero{Fs: bp}
	travelFS(fs, "", renameElement)
}
