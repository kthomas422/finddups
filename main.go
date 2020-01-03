/*
 * Application:     Duplicate file finder
 * File:            main.go
 * Description:     Source file for "main" module
 * Language:        go
 * Dev Env:         Arch Linux 64-bit
 *
 * Author:          Kyle Thomas
 * Date Started:    December 2019
 *
 * Note: absolute paths probably won't work in windows
 */

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var pathDel string // "/" for linux "\\" for windows

type Dup struct {
	Size   int
	Fnames []string
}

type FileStat struct {
	Size int
	Name string
}

func writeDups(dups *[]*Dup, flag bool) {
	var (
		err      error
		i, total int
	)
	f, err := os.Create("dups.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for _, cat := range *dups {
		if len(cat.Fnames) > 1 { // only print out categories with more than 1 filename
			for j, dup := range cat.Fnames {
				fmt.Fprintf(f, "[%d.%d]", i+1, j+1)
				if flag && j > 0 {
					fmt.Fprintf(f, "X")
				} else {
					fmt.Fprintf(f, " ")
				}
				fmt.Fprintf(f, "\t%s\n", dup)
				if err != nil {
					log.Println(err)
				}
				total++
			}
			i++
		}
	}

	err = f.Close()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(total, "Duplicates found from", i, "categories.")
}

func deleteDups(dups *[]*Dup) {
	for _, cat := range *dups {
		if len(cat.Fnames) > 1 {
			for j := 1; j < len(cat.Fnames); j++ { // start at index 1 to keep one of the duplicates
				err := os.Remove(cat.Fnames[j])
				if err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func usage() {
	fmt.Println("$ finddups [flags] [directories...]")
	fmt.Println("-h flag will print this help menu")
	fmt.Println("-D flag will delete the duplicates")
	fmt.Println("\nanything else is assumed to be a directory, if no directory",
		"is passed then current directory is assumed.")
	fmt.Println()
}

func main() {
	var (
		dups []*Dup
		cwd  string
		del  = false
	)
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error: cannot get cwd. ", err)
	}
	if runtime.GOOS == "windows" {
		pathDel = "\\"
	} else {
		pathDel = "/"
	}

	if len(os.Args) == 1 { // no args, do current dir
		finddups(cwd+pathDel, &dups)
	} else {
		for i, arg := range os.Args {
			if i == 0 {
				continue
			}
			if i > 3 {
				log.Println("Error: too many arguments provided")
				usage()
				os.Exit(1)
			}
			if arg == "-h" {
				usage()
				os.Exit(0)
			} else if arg == "-D" {
				del = true
			} else {
				if arg[0] == pathDel[0] { // absolute path
					finddups(arg, &dups)
				} else {
					finddups(cwd+pathDel+arg, &dups)
				}
			}
		}
	}

	writeDups(&dups, del)
	if del {
		fmt.Println("Deleting duplicates...")
		deleteDups(&dups)
	}
	os.Exit(0)
}
