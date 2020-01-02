/*
 * Application:     Duplicate file finder
 * File:            main.go
 * Description:     Source file for "main" module
 * Language:        go
 * Dev Env:         Arch Linux 64-bit
 *
 * Author:          Kyle Thomas
 * Date Started:    December 2019
 */

package main

import (
	"fmt"
	"log"
	"os"
)

type Dup struct {
	Size   int
	Fnames []string
}

type FileStat struct {
	Size int
	Name string
}

func writeDups(dups *[]*Dup) {
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
				fmt.Fprintf(f, "[%d.%d]\t%s\n", i+1, j+1, dup)
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

func usage() {
	fmt.Println("$ finddups [-h or directory]")
	fmt.Println("-h flag will print this help menu")
	fmt.Println("\nanything else is assumed to be a directory, if no directory",
		"is passed then current directory is assumed.")
	fmt.Println()
}

func main() {
	var (
		dups []*Dup
		cwd  string
	)
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error: cannot get cwd. ", err)
	}
	if len(os.Args) > 2 {
		log.Fatal("Error: too many arguments passed in.")
		usage()
		os.Exit(1)
	}
	if len(os.Args) == 2 {
		if os.Args[1] == "-h" {
			usage()
			os.Exit(0)
		}
		if os.Args[1][len(os.Args[1])-1] != '/' {
			os.Args[1] += "/"
		}
		if os.Args[1][0] == '/' { // absolute path
			finddups(os.Args[1], &dups)
		} else { // relative path
			finddups(cwd+"/"+os.Args[1], &dups)
		}
	} else {
		finddups(cwd+"/", &dups)
	}
	writeDups(&dups)
	os.Exit(0)
}
