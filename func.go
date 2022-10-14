package main

// first trial at chunking file without using a bufio scanner
// func SortLargeFile1(ram int) {
// 	file, err := os.Open("100GB.txt")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	chunkSize := ram // in bytes
// 	totalRead := 0
// 	fileInfo, _ := file.Stat()
// 	fileSize := fileInfo.Size()

// 	sizedBufferReader := bufio.NewReaderSize(file, chunkSize)
// 	i := 0
// 	for {
// 		if fileSize-int64(totalRead) < int64(chunkSize) {
// 			chunkSize = int(fileSize) - totalRead
// 		}
// 		fmt.Println("chunkSize", chunkSize)

// 		chunk := make([]byte, chunkSize)

// 		rLen, err := io.ReadFull(sizedBufferReader, chunk)

// 		if err == io.EOF {
// 			break
// 		}

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		fmt.Println("read chunk of size: ", rLen)

// 		fileName := fmt.Sprintf("sorted_%d.txt", i)
// 		f, err := os.Create(fileName)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer f.Close()

// 		w := bufio.NewWriter(f)
// 		wLen, err := w.Write(chunk)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println("wrote chunk of size: ", wLen)
// 		w.Flush()

// 		totalRead += rLen

// 		i++

// 		if totalRead >= int(fileSize) {
// 			break
// 		}
// 	}
// }

// 	if !scan.Scan() {

// 		// sort the chunk
// 		sort.SliceStable(chunk, func(i, j int) bool {
// 			if sortDirection == "asc" {
// 				return chunk[i] < chunk[j]
// 			} else {
// 				return chunk[i] > chunk[j]
// 			}
// 		})

// 		// write the first chunk to a new file
// 		f, err := os.Create("chunk_" + strconv.Itoa(*i) + ".txt")
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		w := bufio.NewWriter(f)
// 		for _, v := range chunk {
// 			// conver int64 to string
// 			w.WriteString(strconv.FormatInt(v, 10) + "\n")
// 		}
// 		w.Flush()

// 		// reset the chunk
// 		chunk = []int64{}
// 		totalRead = 0

// 		f.Close()
// 	} else {
// 		// append to chunk
// 		value, _ := strconv.ParseInt(scan.Text(), 10, 64)
// 		chunk = append(chunk, value)
// 		bytesLen := len(scan.Bytes())
// 		totalRead += int64(bytesLen)
// 	}
