package server

import (
	"log"
	"strings"

	"github.com/nicolasacquaviva/nicolasacquaviva.github.io/models"
)

// cd
func ChangeDirectory(db models.Datastore) func(string, string) string {
	return func(currDir string, dirToGo string) string {
		if dirToGo == "" {
			return "cd:status:1:~"
		}

		errMessage := "cd: not a directory: " + dirToGo
		currentDirContent, err := db.GetContentByParentDir(currDir)

		if err != nil {
			log.Println("Cannot list directory:", err)
			return errMessage
		}

		for _, content := range currentDirContent {
			// if the given dir is equal to the dir name with or without
			// the ending forward slash
			// and the last char of one of the content is a forward slash (means it is a dir)
			if (dirToGo == content || dirToGo == content[:len(content)-1]) && content[len(content)-1:] == "/" {
				return "cd:status:1:" + content
			}
		}

		return errMessage
	}
}

// help
func NewHelp(db models.Datastore) func() string {
	return func() string {
		return `available commands:
		- clear: Clear the console
		- ls: List directory contents
		- cat: Print file content
		- help: Show help about how to use this site`
	}
}

// ls
func NewListDirectory(db models.Datastore) func(string, string) string {
	return func(dir string, params string) string {
		content, err := db.GetContentByParentDir(dir)

		if err != nil {
			log.Println("Cannot list directory:", err)
			return ""
		}

		return strings.Join(content[:], " ")
	}
}

// cat
func NewPrintFileContent(db models.Datastore) func(string) string {
	return func(name string) string {
		if name == "" {
			return "usage: cat [file_name]"
		}

		content := db.GetFileContent(name)

		if content == "" {
			return "cat: " + name + ": No such file or directory"
		}

		return content
	}
}
