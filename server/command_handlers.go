package server

import (
	"log"
	"strings"

	"github.com/nicolasacquaviva/nicolasacquaviva.github.io/models"
)

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
