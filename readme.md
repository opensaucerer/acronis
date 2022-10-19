# acronis

Simple CLI tool for sorting txt files with RAM restrictions.

## Testing
The tool contains a set of tests bundled as a suite using Go's built in testing package. After cloning the code, run the test suite using

```bash
go test -v
```

## Installation

Clone the repository and install the dependencies

```bash
git clone https://github.com/opensaucerer/acronis.git
cd acronis
go mod tidy
```

Build the binary and move to bin (Linux or MacOS)

```bash
go build -o acronis main.go
mv acronis /usr/local/bin
```

## Usage

Create a file of 2GB with random numbers

```bash
acronis create --size 2048
```

Sort a file with 1GB of RAM in ascending order

```bash
acronis sort --file /path/to/file.txt --ram 1024 --direction asc
```

### Problem:

There is a text file with a number in each line, the size of the file is 100GB. Need to write a GoLang command-line program to create a new file sorted in increasing order. It’s not allowed to use more than 1Gb of RAM.

### Restrictions:

1. sort the file in increasing order
2. use no more than 1GB of RAM
3. the file size is 100GB
4. the file contains only numbers
5. the file is not sorted

### Questions:

1. does file contain duplicate numbers?

### Solution:

1. External sort - read the file in chunks of 1GB, sort each chunk and write it to a new file. Repeat until the whole file is sorted.
2. K-way Merge - merge all of the sorted chunks into one file.
