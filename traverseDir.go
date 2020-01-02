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
	"sync"
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

	// get list of files from directory and sub directories
	filesChan := make(chan FileStat, 2048)
	var wg sync.WaitGroup
	wg.Add(1)
	go readFiles(dir, filesChan, &wg)
	go func() {
		wg.Wait()
		close(filesChan)
	}()
	for file := range filesChan {
		listOfFiles = append(listOfFiles, file)
	}
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

func readFiles(rootDir string, filesChan chan<- FileStat, wg *sync.WaitGroup) {
	defer wg.Done()
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
			wg.Add(1)
			go readFiles(rootDir+file.Name(), filesChan, wg)
		} else {
			nc0++
			fileStat = new(FileStat)
			fileStat.Name = rootDir + file.Name()
			fileStat.Size = int(file.Size())
			filesChan <- *fileStat
		}
	}
	return
}
