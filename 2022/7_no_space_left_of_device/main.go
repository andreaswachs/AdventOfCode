package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexForCD                 = regexp.MustCompile(`^\$ cd (?P<dir>.+)$`)
	regexForLS                 = regexp.MustCompile(`^\$ ls`)
	regexForDirListing         = regexp.MustCompile(`^dir (?P<dir>.*)$`)
	regexForFileListing        = regexp.MustCompile(`^(?P<size>\d+) (?P<name>\S+)$`)
	part1AnswerBuffer   uint64 = 0
	directorySizes      []uint64
)

const (
	diskSpace = 70000000
	needSpace = 30000000
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

// Size returns the size of the directory and all subdirectories summed together
func (d *Directory) Size() uint64 {
	var directorySize uint64 = 0
	for _, file := range d.Files {
		directorySize += file.Size
	}

	for _, subdir := range d.Subdirs {
		directorySize += subdir.Size()
	}

	// This is for part 1
	if directorySize <= 100000 {
		part1AnswerBuffer += directorySize
	}

	// This is for part 2
	directorySizes = append(directorySizes, directorySize)

	return directorySize
}

// NewDir creates a new directory with the given name and parent directory.
// It also initializes the subdirectories and files maps.
func NewDir(name string, parentDir *Directory) *Directory {
	subdirs := make(map[string]*Directory)
	files := make(map[string]*File)

	return &Directory{Name: name, Files: files, Subdirs: subdirs, ParentDir: parentDir}
}

// NewFile creates a new file with the given name and size.
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

	// We manually create the root directory as we assume that the input will always start with `cd /`
	rootDirectory := NewDir("/", nil)
	currentDir := rootDirectory

	// We split the input into lines for easier processing
	lines := strings.Split(string(readFile), "\n")
	// We ignore the first two lines as they are the instructions to cd into the root directory and `ls`
	lines = lines[2:]

	// Go throguh all the the lines in the input and build the directory structure
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
	spaceUsed := rootDirectory.Size()
	fmt.Printf("Result part 1: %d\n", part1AnswerBuffer)

	// Solve for part 2
	// Calculate the unused amount of space
	unusedSpace := diskSpace - spaceUsed
	// Calculate the amount of space we need to free
	needToFree := int64(needSpace - unusedSpace)

	// Find out the least amount of space we need to free in order to
	// free the necessary amount of space
	var spaceToFree int64
	for _, size := range directorySizes {
		directorySize := int64(size)
		if needToFree-directorySize < 0 {
			if spaceToFree == 0 || spaceToFree > directorySize {
				spaceToFree = directorySize
			}
		}
	}
	fmt.Printf("Result part 2: %d\n", spaceToFree)
}
