package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

// create a txt file of size 100GB with each line containing a random
func CreateTxt() {
	f, err := os.Create("100GB.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for i := 0; i < 10000000; i++ {
		w.WriteString(strconv.Itoa(rand.Intn(1000000000)) + "\r")
	}
	w.Flush()
}

// Problem: There is a text file with a number in each line, the size of the file is 100GB. Need to write a GoLang command-line program to create a new file sorted in increasing order. It’s not allowed to use more than 1Gb of RAM.

// Restrictions:
// 1. sort the file in increasing order
// 2. use no more than 1GB of RAM
// 3. the file size is 100GB
// 4. the file contains only numbers
// 5. the file is not sorted

// Questions:
// 1. does file contain duplicate numbers?

// Possible Approach:
// 1. External sort - read the file in chunks of 1GB, sort each chunk and write it to a new file. Repeat until the whole file is sorted. Then merge the sorted chunks into one file.

func SortLargeFile(ram int, i *int) {
	file, err := os.Open("100GB.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	chunkSize := int64(float64(ram) * 0.8) // in bytes
	var totalRead int64
	var totalChunkRead int64

	scan := bufio.NewScanner(file)

	// read the first chunk of the file
	var chunk []int64

	for scan.Scan() {
		// append to chunk
		value, _ := strconv.ParseInt(scan.Text(), 10, 64)
		chunk = append(chunk, value)
		bytesLen := len(scan.Bytes())
		totalRead += int64(bytesLen)
		totalChunkRead += int64(bytesLen)

		if totalRead >= chunkSize {

			// sort the chunk
			sort.Slice(chunk, func(i, j int) bool {
				return chunk[i] < chunk[j]
			})

			// write the first chunk to a new file
			f, err := os.Create("chunk_" + strconv.Itoa(*i) + ".txt")
			if err != nil {
				log.Fatal(err)
			}

			w := bufio.NewWriter(f)
			for _, v := range chunk {
				w.WriteString(strconv.FormatInt(v, 10) + "\r")
			}
			w.Flush()

			// reset the chunk
			chunk = []int64{}
			totalRead = 0
			*i++

			f.Chmod(0600)
			f.Close()
		}

		if !scan.Scan() {
			// sort the chunk
			sort.Slice(chunk, func(i, j int) bool {
				return chunk[i] < chunk[j]
			})

			// write the first chunk to a new file
			f, err := os.Create("chunk_" + strconv.Itoa(*i) + ".txt")
			if err != nil {
				log.Fatal(err)
			}

			w := bufio.NewWriter(f)
			for _, v := range chunk {
				// conver int64 to string
				w.WriteString(strconv.FormatInt(v, 10) + "\r")
			}
			w.Flush()

			// reset the chunk
			chunk = []int64{}
			totalRead = 0

			f.Chmod(0600)
			f.Close()
		} else {
			// append to chunk
			value, _ := strconv.ParseInt(scan.Text(), 10, 64)
			chunk = append(chunk, value)
			bytesLen := len(scan.Bytes())
			totalRead += int64(bytesLen)
			totalChunkRead += int64(bytesLen)
		}

	}
}

func MergeKSortedFiles(i int) {
	outfile, err := os.OpenFile("sorted.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	// iterate over sorted chunk files
	for {

		type Heap struct {
			Val int64
			Idx int
		}

		heap := make([]Heap, 0)
		for j := 0; j <= i; j++ {
			// open each file
			file, err := os.OpenFile("chunk_"+strconv.Itoa(j)+".txt", os.O_RDONLY, 0600)
			if err != nil {
				log.Fatal(err)
			}

			// read the first line
			scan := bufio.NewScanner(file)

			for scan.Scan() {
				value, _ := strconv.ParseInt(scan.Text(), 10, 64)
				fmt.Println("value", value)
				heap = append(heap, Heap{
					Val: value,
					Idx: j,
				})
				break
			}

			file.Close()
		}

		if len(heap) == 0 {
			break
		}

		// sort the heap
		sort.Slice(heap, func(i, j int) bool {
			return heap[i].Val < heap[j].Val
		})

		outfile.WriteString(strconv.FormatInt(heap[0].Val, 10) + "\r")

		// open the file with the smallest value and remove the first line
		file, err := os.OpenFile("chunk_"+strconv.Itoa(heap[0].Idx)+".txt", os.O_RDONLY, 0600)
		if err != nil {
			log.Fatal(err)
		}

		// read the first line
		scan := bufio.NewScanner(file)
		var lines []string
		scan.Scan()
		for scan.Scan() {
			lines = append(lines, scan.Text())
		}
		file.Close()

		if len(lines) > 1 {

			// remove the first line
			lines = lines[1:]

			// write the lines back to the file
			f, err := os.Create("chunk_" + strconv.Itoa(heap[0].Idx) + ".txt")
			if err != nil {
				log.Fatal(err)
			}

			w := bufio.NewWriter(f)
			for _, v := range lines {
				w.WriteString(v + "\r")
			}
			w.Flush()

			f.Chmod(0600)
			f.Close()
		}

		// reset the heap
		heap = []Heap{}
	}

	// delete the chunk files
	// for j := 0; j <= i; j++ {
	// 	os.Remove("chunk_" + strconv.Itoa(j) + ".txt")
	// }

	outfile.Close()
}

func main() {
	// args ---
	// 1. ram size in GB
	// 2. sort direction
	// CreateTxt()
	i := 0
	SortLargeFile(int(0.5*1024*1024), &i)
	MergeKSortedFiles(i)
}
