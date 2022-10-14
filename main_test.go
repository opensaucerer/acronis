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
	size int64 = 0.5 * 1024 * 1024
	i          = 0
)

// test function CreateTxt - go test -v -run TestCreateTxt
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

// test function SortLargeFile - go test -v -run TestSortLargeFile
func testSortLargeFile(t *testing.T) {

	// sort the file
	SortLargeFile(dest, size, &i, "asc")

	// check if file chunks are sorted
	file, err := os.Open(fmt.Sprintf("%s_%d.txt", "chunk", i))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// check if file is sorted
	var prev int64
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		value, _ := strconv.ParseInt(scan.Text(), 10, 64)
		if value < prev {
			t.Error("File chunks are not sorted")
		}
		prev = value
	}

}

func TestAcronis(t *testing.T) {
	t.Run("Should create a file of the specified size with random numbers", testCreateTxt)
	t.Run("Should ascendingly sorted chunks of the file", testSortLargeFile)

	tearDown()
}

func tearDown() {
	// remove the main file
	os.Remove(dest)

	// remove the chunks
	for i := 0; i < 10; i++ {
		os.Remove(fmt.Sprintf("%s_%d.txt", "chunk", i))
	}
}
