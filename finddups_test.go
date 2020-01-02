/*
 * Application:     Duplicate file finder
 * File:            test.go
 * Description:     Test file for "finddups" program
 * Language:        go
 * Dev Env:         Arch Linux 64-bit
 *
 * Author:          Kyle Thomas
 * Date Started:    December 2019
 */

package main

import "testing"

type compareFilesPairs struct {
	f1, f2 string
	same   bool
}

var (
	compareFilesTest = [...]compareFilesPairs{
		{"Tests/t1.txt", "Tests/t2.txt", true},
		{"Tests/t1.txt", "Tests/t3.txt", false},
		{"Tests/t1.txt", "Tests/t4.txt", false},
		{"Tests/t4.txt", "Tests/t1.txt", false},
	}
)

func TestCompareFiles(t *testing.T) {
	for _, test := range compareFilesTest {
		if compareFiles(test.f1, test.f2) != test.same {
			t.Error("Test error:", test.f1, test.f2, test.same)
		}
	}
}

func TestWriteDups(t *testing.T) {
	var dups = []Dup{
		{0, []string{"a1"}},
		{0, []string{"b1", "b2"}},
		{0, []string{"c1", "c2", "c3", "c4"}},
		{0, []string{"d1"}},
		{0, []string{"e1", "e2", "e3"}},
	}
	writeDups(&dups)
}
