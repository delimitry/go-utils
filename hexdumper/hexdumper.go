package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

const bytesPerLine = 16
const hexAlphabet = "0123456789abcdef"

func bytesToHex(data []byte) []byte {
	dataLen := len(data)
	out := make([]byte, dataLen*3-1)
	for i := 0; i < dataLen; i++ {
		out[i*3] = hexAlphabet[data[i]>>4]
		out[i*3+1] = hexAlphabet[data[i]&0x0f]
		if i != dataLen-1 {
			out[i*3+2] = ' '
		}
	}
	return out
}

func bytesToPrintable(data []byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		if b < 32 || b > 126 {
			out[i] = '.'
		} else {
			out[i] = b
		}
	}
	return out
}

func dumpFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fi.Size()

	offset := 0
	for {
		data := make([]byte, bytesPerLine)
		n, err := f.Read(data)
		if err != nil {
			return err
		}
		if fileSize < 0xffffffff {
			fmt.Printf("%08x | %-47s | %s\n", offset, bytesToHex(data[:n]), bytesToPrintable(data))
		} else {
			fmt.Printf("%016x | %-47s | %s\n", offset, bytesToHex(data[:n]), bytesToPrintable(data))
		}
		if n < bytesPerLine {
			break
		}
		offset += n
	}
	return err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <filename>\n", path.Base(os.Args[0]))
		os.Exit(1)
	}
	filename := os.Args[1]
	err := dumpFile(filename)
	if err != nil {
		log.Fatal(err)
	}
}
