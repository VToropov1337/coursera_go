package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	return dirChecker(out, path, printFiles, "")
}

func dirChecker(out io.Writer, path string, printFiles bool, prefix string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(list); i++ {
		if !printFiles && !list[i].IsDir() {
			list = append(list[:i], list[i+1:]...)
			i--
		}
	}

	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	for idx, file := range list {

		filePrefix := ""
		if idx == len(list)-1 {
			filePrefix = "└───"
		} else {
			filePrefix = "├───"
		}
		if file.Size() > 0 && !file.IsDir() {
			fmt.Fprintln(out, prefix+filePrefix+file.Name()+" ("+strconv.Itoa(int(file.Size()))+"b)")
		} else if file.Size() == 0 && !file.IsDir() {
			fmt.Fprintln(out, prefix+filePrefix+file.Name()+" (empty)")
		} else {
			fmt.Fprintln(out, prefix+filePrefix+file.Name())
		}

		if file.IsDir() {
			var err error
			if file.Name() == list[len(list)-1].Name() {
				dirChecker(out, path+"/"+file.Name(), printFiles, prefix+"\t")
			} else {
				dirChecker(out, path+"/"+file.Name(), printFiles, prefix+"│\t")
			}
			if err != nil {
				panic(err.Error())
			}
		}

	}

	return nil
}
