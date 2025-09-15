package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
)

// getDotFilePath returns the dot file for the repos list
// creates it and the enclosing folder if it does not exist
func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dotFile := usr.HomeDir + "/.gogitlocalstats"
	return dotFile
}

// openFile opens the file located at `filePath`. Creates it if not existing.
func openFile(filePath string) *os.File {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Errorf("openFile %q: %w", filePath, err))
	}
	return f
}

// given a file path string, get the content of each line and parses it to a slice of strings
func parseFileLinesToSlice(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Errorf("open %q: %w", filePath, err))
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("scan %q: %w", filePath, err))
	}
	return lines
}

// slice contains return true if slice contains value
func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// joinSlices adds the element of the new slice into the exisiting slice
func joinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !sliceContains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

func dumpStringsSliceToFile(repos []string, filePath string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(filePath, []byte(content), 0755)

}

// parse the existing repos storede in the file to a slice
// add the new items to the slice, without dupes
// store a the slice file and overwriting the exisiting content

func addNewSliceElementsToFile(filePath string, newRepos []string) {
	extistingRepos := parseFileLinesToSlice(filePath)
	repos := joinSlices(newRepos, extistingRepos)
	dumpStringsSliceToFile(repos, filePath)
}

// starts the recursive search of git repos
// living in the folder subtree
func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

func scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repositories := recursiveScanFolder(folder)
	filePath := getDotFilePath()
	addNewSliceElementsToFile(filePath, repositories)
	fmt.Print("\n\n Successfully added \n\n")
}

// scanGitFolders return a list of subfolders
// returns the base folder of the repo, the .git folder parent.
// recursively searches the subfolders
func scanGitFolders(folders []string, folder string) []string {
	folder = strings.TrimSuffix(folder, "/")
	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
	var path string

	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}
			folders = scanGitFolders(folders, path)
		}
	}
	return folders
}

type User struct {
	//userID
	Uid string
	// primary group ID
	Gid string
	// login name
	Username string
	// real name or display name
	Name string
	//hoe directory
	HomeDir string
}
