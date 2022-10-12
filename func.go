package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

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

func SortLargeFile2(ram int) {
	file, err := os.Open("100GB.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	chunkSize := int64(float64(ram) * 0.8) // in bytes
	var totalRead int64
	var totalChunkRead int64
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	i := 0

	scan := bufio.NewScanner(file)

	// read the first chunk of the file
	var chunk []string

	for scan.Scan() {
		chunk = append(chunk, scan.Text())
		bytesLen := len(chunk)
		totalRead += int64(bytesLen)
		totalChunkRead += int64(bytesLen)

		fmt.Println("chunk size", chunkSize)
		fmt.Println("totalRead", totalRead)
		fmt.Println("fileSize", fileSize)
		fmt.Println("totalChunkRead", totalChunkRead)
		if totalRead >= chunkSize {

			// write the first chunk to a new file
			f, err := os.Create("chunk_" + strconv.Itoa(i) + ".txt")
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			w := bufio.NewWriter(f)
			for _, v := range chunk {
				w.WriteString(v + "\r")
			}
			w.Flush()

			// reset the chunk
			chunk = []string{}
			totalRead = 0
			i++

			fmt.Println("chunk size:", chunkSize, "chunk size left:", fileSize-totalChunkRead)
			// if fileSize-totalChunkRead < chunkSize {
			// 	chunkSize = fileSize - totalChunkRead
			// 	fmt.Println("new chunk size:", chunkSize)
			// }

		}

		if !scan.Scan() {
			// write the first chunk to a new file
			f, err := os.Create("chunk_" + strconv.Itoa(i) + ".txt")
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			w := bufio.NewWriter(f)
			for _, v := range chunk {
				w.WriteString(v + "\r")
			}
			w.Flush()

			// reset the chunk
			chunk = []string{}
			totalRead = 0
			i++
		}

	}
}
