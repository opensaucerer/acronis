package main_test

import (
	"main"
	"testing"
)

// func TestSortLargeFile(t *testing.T) {
// 	// test function SortLargeFile

// }

func TestCreateTxt(t *testing.T) {
	// test function CreateTxt

	dest := "demo.txt"
	size := 10 * 1024 * 1024

	main.CreateTxt(dest, size)

}
