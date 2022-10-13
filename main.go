package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

// create a txt file of size 100GB with each line containing a random
func CreateTxt(filepath string, size int64) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	for i := int64(0); i < size; i++ {
		w.WriteString(strconv.FormatInt(int64(rand.Intn(int(size))), 10) + "\r")
	}
	w.Flush()
	f.Close()
}

func SortLargeFile(filepath string, ram int64, i *int, sortDirection string) {
	file, err := os.Open(filepath)
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
			sort.SliceStable(chunk, func(i, j int) bool {
				if sortDirection == "asc" {
					return chunk[i] < chunk[j]
				} else {
					return chunk[i] > chunk[j]
				}
			})

			// write the first chunk to a new file
			f, err := os.OpenFile("chunk_"+strconv.Itoa(*i)+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

			f.Close()
		}

		if !scan.Scan() {
			// sort the chunk
			sort.SliceStable(chunk, func(i, j int) bool {
				if sortDirection == "asc" {
					return chunk[i] < chunk[j]
				} else {
					return chunk[i] > chunk[j]
				}
			})

			// write the first chunk to a new file
			f, err := os.OpenFile("chunk_"+strconv.Itoa(*i)+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func MergeKSortedFiles(outpath string, i int, sortDirection string) {
	outfile, err := os.OpenFile(outpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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
			file, err := os.Open("chunk_" + strconv.Itoa(j) + ".txt")
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
		sort.SliceStable(heap, func(i, j int) bool {
			if sortDirection == "asc" {
				return heap[i].Val < heap[j].Val
			} else {
				return heap[i].Val > heap[j].Val
			}
		})

		outfile.WriteString(strconv.FormatInt(heap[0].Val, 10) + "\r")

		// open the file with the smallest value and remove the first line
		file, err := os.Open("chunk_" + strconv.Itoa(heap[0].Idx) + ".txt")
		if err != nil {
			log.Fatal(err)
		}

		// read the first line
		scan := bufio.NewScanner(file)
		var lines []string
		for scan.Scan() {
			lines = append(lines, scan.Text())
		}
		file.Close()

		if len(lines) > 1 {

			// remove the first line
			lines = lines[1:]

			// write the lines back to the file
			f, err := os.OpenFile("chunk_"+strconv.Itoa(heap[0].Idx)+".txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}

			w := bufio.NewWriter(f)
			for _, v := range lines {
				w.WriteString(v + "\r")
			}
			w.Flush()

			f.Close()
		}

		// reset the heap
		heap = []Heap{}
	}

	// delete the chunk files
	for j := 0; j <= i; j++ {
		os.Remove("chunk_" + strconv.Itoa(j) + ".txt")
	}

	outfile.Close()
}

var app = cli.NewApp()

func info() {
	app.Name = "Simple CLI tool for sorting txt files"
	app.Usage = `Sort a txt file with RAM restrictions
	Usage: acronis --file <file> --direction <sort direction> --ram <max memory in MB> -outpath <output file>

	Example:
	acronis sort --file /path/to/input.txt --ram 100 --direction asc
	acronis sort --file /path/to/input.txt --ram 100 --direction asc --outpath /path/to/output.txt`
	app.Author = "Perfection (https://github.com/opensaucerer)"
	app.Version = "0.0.1"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:    "sort",
			Aliases: []string{"s"},
			Usage: `Sort a txt file with RAM restrictions
						acronis sort --file /path/to/input.txt --ram 100 --direction asc
						acronis sort --file /path/to/input.txt --ram 100 --direction asc --outpath /path/to/output.txt`,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "file",
					Usage:    "Path to the input file",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "direction",
					Usage:    "Sort direction (asc or desc)",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "outpath",
					Usage:    "Output path for the sorted file (defaults to sorted.txt in curr dir)",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "ram",
					Usage:    "RAM limit in MB",
					Required: true,
				},
			},
			Action: func(c *cli.Context) {
				// get the RAM limit
				ram := c.String("ram")
				var r int64
				if strings.Contains(ram, ".") {
					// parse the float
					ramFloat, err := strconv.ParseFloat(ram, 64)
					if err != nil {
						log.Fatal(err)
					}
					r = int64(ramFloat * 1024 * 1024)
				} else {
					// parse the int
					ramInt, err := strconv.ParseInt(ram, 10, 64)
					if err != nil {
						log.Fatal(err)
					}
					r = ramInt * 1024 * 1024
				}

				dir := c.String("direction")

				outpath := c.String("outpath")
				if outpath == "" {
					outpath = "sorted.txt"
				}

				fmt.Printf("Sorting %s with %s MB (%d bytes) of RAM in %s direction. Output file: %s\n", c.String("file"), ram, r, dir, outpath)

				filepath := c.String("file")
				// validate that it's a txt file
				s := strings.Split(filepath, ".")
				if s[len(s)-1] != "txt" {
					log.Fatal("File is not a txt file")
				}

				// sort and merge
				i := 0
				SortLargeFile(filepath, r, &i, dir)
				MergeKSortedFiles(outpath, i, dir)

				fmt.Println("Sorted file saved to", outpath)
			},
		},

		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create a txt file with random numbers",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "size",
					Usage:    "Size of the file in MB",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "outpath",
					Usage:    "Output path for the sorted file (defaults to lacronis.txt in curr dir)",
					Required: false,
				},
			},
			Action: func(c *cli.Context) {
				// get the file size
				size := c.String("size")
				var s int64
				if strings.Contains(size, ".") {
					// parse the float
					sizeFloat, err := strconv.ParseFloat(size, 64)
					if err != nil {
						log.Fatal(err)
					}
					s = int64(sizeFloat * 1024 * 1024)
				} else {
					// parse the int
					sizeInt, err := strconv.ParseInt(size, 10, 64)
					if err != nil {
						log.Fatal(err)
					}
					s = sizeInt * 1024 * 1024
				}

				outpath := c.String("outpath")
				if outpath == "" {
					outpath = "lacronis.txt"
				}

				fmt.Printf("Creating a random txt file of size %s MB (%d bytes) at %s\n", size, s, outpath)

				// create the file
				CreateTxt(outpath, s)

				fmt.Println("Random txt file saved to", outpath)
			},
		},
	}
}

func main() {
	// args ---
	// 1. ram size in GB
	// 2. sort direction
	// CreateTxt()
	// i := 0
	// dir := "asc"
	// SortLargeFile(int(0.5*1024*1024), &i, dir)
	// MergeKSortedFiles(i, dir)
	// for i < 2 {
	// 	f, _ := os.Open(fmt.Sprintf("chunk_%d.txt", i))
	// 	info, _ := f.Stat()
	// 	// get file mode
	// 	fmt.Println(info.Mode())
	// 	// get file permission
	// 	fmt.Println(info.Mode().Perm())
	// 	// get mode type
	// 	fmt.Println(info.Mode().Type())
	// 	// get mode in string
	// 	fmt.Println(info.Mode().IsRegular())

	// 	b := bufio.NewScanner(f)
	// 	for b.Scan() {
	// 		fmt.Printf("from chunk %d: %s\n", i, b.Text())
	// 		break
	// 	}

	// 	i++
	// }
	info()
	commands()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
