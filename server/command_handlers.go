package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/nicolasacquaviva/nicolasacquaviva.github.io/models"
)

// get the last part of a path, used to get the filename or the directory
func getPathLastPart(path string) string {
	splittedPath := strings.Split(path, "/")

	return splittedPath[len(splittedPath)-1]
}

// display
func NewDisplayImage(db models.Datastore) func(string) string {
	return func(imagePath string) string {
		content, err := db.GetContentByPath(imagePath)

		if err != nil {
			log.Println("Cannot get image:", err)
			return fmt.Sprintf("display: %s: No such image", getPathLastPart(imagePath))
		}

		return fmt.Sprintf("display:status:1:%s", content.Content)
	}
}

// cd
func NewChangeDirectory(db models.Datastore) func(string, string) string {
	return func(currDir string, dirToGo string) string {
		// go to home dir (~) if no given dir name
		if dirToGo == "" {
			return "cd:status:1:~"
		}

		errMessage := "cd: not a directory: " + dirToGo

		if dirToGo == ".." {
			parent, err := db.GetContentsParentByChild(getPathLastPart(currDir))

			if err != nil {
				log.Println("Cannot list directory:", err)
				return errMessage
			}

			return "cd:status:1:" + parent.Path
		}

		currentDirContent, err := db.GetContentByParentDir(currDir)

		if err != nil {
			log.Println("Cannot list directory:", err)
			return errMessage
		}

		// checks if the dir to go is part of the current directory
		for _, content := range currentDirContent {
			// if the given dir is equal to the dir name with or without
			// the ending forward slash
			// and the last char of one of the content is a forward slash (means it is a dir)
			if (dirToGo == content || dirToGo == content[:len(content)-1]) && content[len(content)-1:] == "/" {
				dirToGoContent, err := db.GetContentByName(dirToGo)

				if err != nil {
					log.Println("Cannot list directory:", err)
					return errMessage
				}

				return "cd:status:1:" + dirToGoContent.Path
			}
		}

		return errMessage
	}
}

// help
func NewHelp(db models.Datastore) func() string {
	return func() string {
		return `available commands:
		- cat: Print file content
		- cd: Change directory
		- clear: Clear the console
		- display: Display image file
		- help: Show help about how to use this site
		- ls: List directory contents`
	}
}

// ls
func NewListDirectory(db models.Datastore) func(string, string) string {
	return func(dir string, params string) string {
		content, err := db.GetContentByParentDir(getPathLastPart(dir))

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
