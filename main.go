package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

// create a txt file of size with each line containing a random
func CreateTxt(filepath string, size int64) {
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	log.Println("Filling file with random numbers...")
	for i := int64(0); i < size; i++ {
		inf, _ := f.Stat()
		if inf.Size() >= size {
			break
		}
		w.WriteString(strconv.FormatInt(int64(rand.Intn(int(size))), 10) + "\n")
	}
	w.Flush()
	f.Close()
}

func SortLargeFile(filepath string, ram int64, i *int, sortDirection string) {
	log.Println("Opening file for sorting...")
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	chunkSize := int64(float64(ram) * 0.8) // in bytes
	var totalRead int64

	scan := bufio.NewScanner(file)

	// read the first chunk of the file
	var chunk []int64

	log.Println("Creating chunks and performing initial sort...")
	for scan.Scan() {
		// append to chunk
		value, _ := strconv.ParseInt(scan.Text(), 10, 64)
		chunk = append(chunk, value)
		bytesLen := len(scan.Bytes())
		totalRead += int64(bytesLen)

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
			f, err := os.Create("chunk_" + strconv.Itoa(*i) + ".txt")
			if err != nil {
				log.Fatal(err)
			}

			w := bufio.NewWriter(f)
			for _, v := range chunk {
				w.WriteString(strconv.FormatInt(v, 10) + "\n")
			}
			w.Flush()

			// reset the chunk
			chunk = []int64{}
			totalRead = 0
			*i++

			f.Close()
		}

	}

	// perform flush
	if totalRead > 0 && len(chunk) > 0 {
		log.Println("Performing flush into last chunk...")

		// sort the chunk
		sort.SliceStable(chunk, func(i, j int) bool {
			if sortDirection == "asc" {
				return chunk[i] < chunk[j]
			} else {
				return chunk[i] > chunk[j]
			}
		})

		// write the first chunk to a new file
		f, err := os.Create("chunk_" + strconv.Itoa(*i) + ".txt")
		if err != nil {
			log.Fatal(err)
		}

		w := bufio.NewWriter(f)
		for _, v := range chunk {
			w.WriteString(strconv.FormatInt(v, 10) + "\n")
		}
		w.Flush()

		// reset the chunk
		chunk = []int64{}
		totalRead = 0

		f.Close()
	}
}

func MergeKSortedFiles(outpath string, i int, sortDirection string) {
	outfile, err := os.Create(outpath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Merging %d sorted chunks into out file...\n", i+1)
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
				os.Remove(outpath)
				log.Fatal(err)
			}

			// read the first line
			scan := bufio.NewScanner(file)

			for scan.Scan() {
				// save the first line to the heap
				value, err := strconv.ParseInt(scan.Text(), 10, 64)
				if err != nil {
					break
				}
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

		outfile.WriteString(strconv.FormatInt(heap[0].Val, 10) + "\n")

		// open the file with the smallest value and remove the first line
		file, err := os.Open("chunk_" + strconv.Itoa(heap[0].Idx) + ".txt")
		if err != nil {
			os.Remove(outpath)
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
		} else {
			// don't remove, just set to empty
			lines = []string{}
		}

		// write the lines back to the file
		f, err := os.Create("chunk_" + strconv.Itoa(heap[0].Idx) + ".txt")
		if err != nil {
			os.Remove(outpath)
			log.Fatal(err)
		}

		w := bufio.NewWriter(f)
		for _, v := range lines {
			w.WriteString(v + "\n")
		}
		w.Flush()
		f.Close()

		// reset the heap
		heap = []Heap{}
	}

	log.Println("Almost done...just performing cleaups...")
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

				log.Printf("Sorting %s with %s MB (%d bytes) of RAM in %s direction. Output file: %s...\n", c.String("file"), ram, r, dir, outpath)

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

				log.Println("Sorted file saved to", outpath)
			},
		},

		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage: `Create a txt file with random numbers
						acronis create --size 10
						acronis create --size 10 --outpath /path/to/output.txt`,
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

				log.Printf("Creating a random txt file of size %s MB (%d bytes) at %s...\n", size, s, outpath)

				// create the file
				CreateTxt(outpath, s)

				log.Println("Random txt file saved to", outpath)
			},
		},
	}
}

func main() {
	info()
	commands()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
