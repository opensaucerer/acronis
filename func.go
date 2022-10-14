package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// first trial at chunking file without using a bufio scanner
func SortLargeFile1(ram int) {
	file, err := os.Open("100GB.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	chunkSize := ram // in bytes
	totalRead := 0
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	sizedBufferReader := bufio.NewReaderSize(file, chunkSize)
	i := 0
	for {
		if fileSize-int64(totalRead) < int64(chunkSize) {
			chunkSize = int(fileSize) - totalRead
		}
		fmt.Println("chunkSize", chunkSize)

		chunk := make([]byte, chunkSize)

		rLen, err := io.ReadFull(sizedBufferReader, chunk)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("read chunk of size: ", rLen)

		fileName := fmt.Sprintf("sorted_%d.txt", i)
		f, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		wLen, err := w.Write(chunk)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("wrote chunk of size: ", wLen)
		w.Flush()

		totalRead += rLen

		i++

		if totalRead >= int(fileSize) {
			break
		}
	}
}
