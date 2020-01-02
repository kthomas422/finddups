/*
 * Application:     Duplicate file finder
 * File:            traverseDir.go
 * Description:     Source file for "Directory Traversal" module
 * Language:        go
 * Dev Env:         Arch Linux 64-bit
 *
 * Author:          Kyle Thomas
 * Date Started:    December 2019
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func finddups(dir string, dups *[]*Dup) {
	var (
		listOfFiles []FileStat
		tmpDup      *Dup
		i, nDups    int
	)
	d, err := os.Open(dir)
	if err != nil {
		log.Fatal("Error: Cannot open directory ", "[", dir, "] ", err)
	}
	stat, err := d.Stat()
	if err != nil {
		log.Fatal("Error: Cannot get directory stats: ", "[", dir, "] ", err)
	}
	d.Close()
	if !stat.IsDir() {
		log.Fatal("Error: ", "[", dir, "]", " is not a directory.")
	}

	// read dir and add all filenames found into list  (this isn't the most efficient way but it works)
	listOfFiles = readFiles(dir)
	if len(listOfFiles) == 0 {
		log.Fatal("Error: no files found.")
	}

	/* Add first file in list of files to 'dups' list */
	tmpDup = new(Dup)
	tmpDup.Size = listOfFiles[0].Size
	tmpDup.Fnames = append(tmpDup.Fnames, listOfFiles[0].Name)
	*dups = append(*dups, tmpDup)

	/* iterater over the rest of the list, put in new category or dup category if appropiate */
	fmt.Println("Checking", len(listOfFiles), "files for duplicates.")
	for i = 1; i < len(listOfFiles); i++ {
		//fmt.Println("File", i+1, "of", len(listOfFiles), "("+listOfFiles[i].Name+")")
		if !checkCategories(listOfFiles[i], dups) {
			tmpDup = new(Dup)
			tmpDup.Size = listOfFiles[i].Size
			tmpDup.Fnames = append(tmpDup.Fnames, listOfFiles[i].Name)
			*dups = append(*dups, tmpDup)
		} else {
			nDups++
		}
	}
}

func readFiles(rootDir string) (returnFiles []FileStat) {
	var fileStat *FileStat
	nc0 := 0
	if rootDir[len(rootDir)-1] != '/' {
		rootDir += "/"
	}
	filesList, err := ioutil.ReadDir(rootDir)
	if err != nil {
		log.Println(err)
	}
	for _, file := range filesList {
		if file.IsDir() {
			returnFiles = append(returnFiles, readFiles(rootDir+file.Name())...)
		} else {
			nc0++
			fileStat = new(FileStat)
			fileStat.Name = rootDir + file.Name()
			fileStat.Size = int(file.Size())
			returnFiles = append(returnFiles, *fileStat)
		}
	}
	return returnFiles
}
