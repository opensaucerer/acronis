package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
)

var (
	dest       = "demo.txt"
	out        = "demo_sorted.txt"
	dir        = "asc"
	size int64 = 0.5 * 1024 * 1024
	ram  int64 = 0.2 * 1000 * 1000
	i          = 0
)

func isSorted(file *os.File, dir string) bool {
	var prev int64
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		value, _ := strconv.ParseInt(scan.Text(), 10, 64)
		if dir == "asc" {
			if value < prev {
				return false
			}
		} else {
			if value > prev {
				return false
			}
		}
		prev = value
	}
	return true
}

// test function CreateTxt
func testCreateTxt(t *testing.T) {

	CreateTxt(dest, size)

	// check if file exists
	_, err := os.Stat(dest)
	if os.IsNotExist(err) {
		t.Error("File was not created")
	}

	// check if file size is correct
	file, err := os.Open(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	inf, _ := file.Stat()
	if inf.Size() < size {
		t.Error("File size is less than expected, expected >= ", size, " got: ", inf.Size())
	}

	// check if file contains only numbers
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		_, err := strconv.ParseInt(scan.Text(), 10, 64)
		if err != nil {
			t.Error("File contains non-numeric values")
		}
	}

}

// test function SortLargeFile
func testSortLargeFile(t *testing.T) {

	// sort the file
	SortLargeFile(dest, ram, &i, dir)

	// check if file chunks are sorted
	file, err := os.Open(fmt.Sprintf("%s_%d.txt", "chunk", i))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// check if file is sorted
	if !isSorted(file, dir) {
		t.Error("File is not sorted")
	}
}

// test function MergeKSortedFiles
func testMergeKSortedFiles(t *testing.T) {

	// merge the chunks
	MergeKSortedFiles(out, i, dir)

	// check if file is sorted
	file, err := os.Open(out)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	// check if file is sorted
	if !isSorted(file, dir) {
		t.Error("File is not sorted")
	}

}

// go test -v -run TestAcronis
func TestAcronis(t *testing.T) {
	t.Run("Should create a file of the specified size with random numbers", testCreateTxt)
	t.Run("Should ascendingly sorted chunks of the file", testSortLargeFile)
	t.Run("Should merge the chunks into one file", testMergeKSortedFiles)

	// clean up
	tearDown()
}

func tearDown() {
	// remove the main file
	os.Remove(dest)

	// remove the chunks
	for i := 0; i < int(size/ram); i++ {
		os.Remove(fmt.Sprintf("%s_%d.txt", "chunk", i))
	}

	// remove the sorted file
	os.Remove(out)
}
