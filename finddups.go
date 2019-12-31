/*
 * Application:     Duplicate file finder
 * File:            finddups.go
 * Description:     Source file for "finddups" module
 * Language:        go
 * Dev Env:         Arch Linux 64-bit
 *
 * Author:          Kyle Thomas
 * Date Started:    December 2019
 */

package main

import (
	"io"
	"log"
	"os"
)

/*
 * Will put duplicates in 2d list of strings containing filenames... that should be all you need
 */

// compareFiles takes in two file names, opens the files and returs true if they're the same, false otherwise
func compareFiles(fn1, fn2 string) bool {
	const chunkSize = 64000

	// open files
	f1, err := os.Open(fn1)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(fn2)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	// check if filesize is the same
	stat1, err1 := f1.Stat()
	if err1 != nil {
		log.Fatal(err1)
	}
	stat2, err2 := f2.Stat()
	if err2 != nil {
		log.Fatal(err2)
	}
	if stat1.Size() != stat2.Size() {
		return false
	}

	// compare that every byte is the same
	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)
		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)
		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}
		if len(b1) != len(b2) {
			return false
		}
		for i := 1; i < len(b1); i++ {
			if b1[i] != b2[i] {
				return false
			}
		}
	}
}
