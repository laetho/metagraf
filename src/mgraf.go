package main

import (
	"fmt"
	"os"
	"metagraf/src/metagraf"
	"path/filepath"
	"strings"
)



func main() {
	var files []string
	basepath := "/home/a01595/go/src/metagraf/collections/poc"

	// Wal the directory provided in basepath
	err := filepath.Walk(basepath, func( path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Loop through all files in the directory given in basepath,
	// ignore basepath itself, and ignore file names not containing "json"
	// Parse each json file
	for _, file := range files {
		if file == basepath { continue }
		if !strings.Contains( file, "json") { continue }
		fmt.Println("Parsing MetaGraf ", file)
		metagraf.Parse( file )
	}
}
