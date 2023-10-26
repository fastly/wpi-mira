package main

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"fmt"
	"io"
	"log"
	"os"

	mrt "github.com/osrg/gobgp/pkg/packet/mrt"
)

func main() {
	file, err := os.Open("bgpfiles/rib.20160629.1600.bz2")
	if err != nil {
		log.Fatal(err)
	}

	bz2Scanner := bzip2.NewReader(file)

	// Put the decompressed data into a byte array
	var data []byte
	for {
		buffer := make([]byte, 1024) // You can adjust the buffer size as needed
		n, err := bz2Scanner.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Fatal("Error reading BZ2 data:", err)
			}
			break
		}
		data = append(data, buffer[:n]...)
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	scanner.Split(mrt.SplitMrt)

	count := 0

	for {
		scanSuccess := scanner.Scan()
		if !scanSuccess {
			break
		}

		body := scanner.Bytes()

		hdr, err := mrt.NewMRTHeader(0, mrt.BGP4MP, mrt.RIB_IPV4_UNICAST, 0)
		if err != nil {
			fmt.Println(err)
			break
		}
		hdr.DecodeFromBytes(body)

		msg, err := mrt.ParseMRTBody(hdr, body[mrt.MRT_COMMON_HEADER_LEN:])
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(msg)
		fmt.Println()

		count += 1
	}

	fmt.Println("Message Count: ", count)
}
