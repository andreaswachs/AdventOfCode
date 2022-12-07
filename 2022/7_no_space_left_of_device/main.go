package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexForCD          = regexp.MustCompile(`^\$ cd (?P<dir>.+)$`)
	regexForLS          = regexp.MustCompile(`^\$ ls`)
	regexForDirListing  = regexp.MustCompile(`^dir (?P<dir>.*)$`)
	regexForFileListing = regexp.MustCompile(`^(?P<size>\d+) (?P<name>\S+)$`)
)

type lineType uint8

const (
	changeDirectory lineType = iota
	listDirectory
	directoryEntry
	fileEntry
)

type File struct {
	Name string
	Size uint64
}

type Directory struct {
	Name      string
	Files     map[string]*File
	Subdirs   map[string]*Directory
	ParentDir *Directory
}

func (d *Directory) AddDir(directory *Directory) {
	d.Subdirs[directory.Name] = directory
}

func (d *Directory) AddFile(file *File) {
	d.Files[file.Name] = file
}

func (d *Directory) Size() uint64 {
	var subDirsSize uint64 = 0

	for _, subDir := range d.Subdirs {
		subDirsSize += subDir.Size()
	}

	var thisDirSize uint64 = 0
	for _, file := range d.Files {
		thisDirSize += file.Size
	}

	if thisDirSize > 100000 {
		return 0 + subDirsSize
	}

	return thisDirSize + subDirsSize
}

func NewDir(name string, parentDir *Directory) *Directory {
	subdirs := make(map[string]*Directory)
	files := make(map[string]*File)

	return &Directory{Name: name, Files: files, Subdirs: subdirs, ParentDir: parentDir}
}

func NewFile(name string, size uint64) *File {
	return &File{Name: name, Size: size}
}

func isCmdCD(input string) (bool, string) {
	return getRegexMatch(regexForCD, input, "dir")
}

func isCmdLS(input string) bool {
	return regexForLS.Match([]byte(input))
}

func isListingDir(input string, parentDir *Directory) (bool, *Directory) {
	isDirListing, dirName := getRegexMatch(regexForDirListing, input, "dir")
	if !isDirListing {
		return false, nil
	}

	return true, NewDir(dirName, parentDir)
}

func isListingFile(input string) (bool, *File) {
	isFileListing, fileName := getRegexMatch(regexForFileListing, input, "name")
	if !isFileListing {
		return false, nil
	}

	// Since we didnt' exit early, we know this is a file listing
	_, size := getRegexMatch(regexForFileListing, input, "size")
	sizeAsInt, err := strconv.Atoi(size)
	if err != nil {
		panic(err)
	}

	return true, NewFile(fileName, uint64(sizeAsInt))
}

func getRegexMatch(rexp *regexp.Regexp, input string, group string) (bool, string) {
	if !rexp.Match([]byte(input)) {
		return false, ""
	}

	matches := rexp.FindStringSubmatch(input)
	index := rexp.SubexpIndex(group)
	return true, matches[index]

}

// Could probably make effort to use generics but oh well
func getLineTypeAndPayload(input string, currentDir *Directory) (lineType, interface{}) {
	if isCD, dir := isCmdCD(input); isCD {
		return changeDirectory, dir
	}

	if isCmdLS(input) {
		return listDirectory, nil
	}

	if isDirListing, dir := isListingDir(input, currentDir); isDirListing {
		return directoryEntry, dir
	}

	// We could default to return a fileEntry, but I want it to panic
	// if we don't explicit verify that it is a file, cause then I did something wrong
	if isFileListing, file := isListingFile(input); isFileListing {
		return fileEntry, file
	}

	panic(input)
}

func main() {
	readFile, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	rootDirectory := NewDir("/", nil)
	currentDir := rootDirectory

	lines := strings.Split(string(readFile), "\n")
	lines = lines[2:] // We ignore the first two lines of the input which is just `cd /` and `ls`.

	for _, line := range lines {
		lt, payload := getLineTypeAndPayload(line, currentDir)
		switch lt {
		case changeDirectory:
			if payload == ".." {
				currentDir = currentDir.ParentDir
				continue
			}
			currentDir = currentDir.Subdirs[payload.(string)]
			continue
		case listDirectory:
			// We don't really do anything interesting in this case
			continue
		case directoryEntry:
			currentDir.AddDir(payload.(*Directory))
		case fileEntry:
			currentDir.AddFile(payload.(*File))
		}
	}

	// Solve for part 1

	fmt.Printf("Result part 1: %d\n", rootDirectory.Size())

}
