package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"encoding/json"

	"encoding/hex"
	"github.com/twmb/franz-go/pkg/kgo"
	obmp "github.com/sdstrowes/go-openbmp"
	gobmp "github.com/sbezverk/gobmp/pkg/bmp"
)

func makePretty(obmp *obmp.OpenBMPHeader) string {
	tmpJson, err := json.MarshalIndent(*obmp, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(tmpJson)
}

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func main() {
	var broker string  
	flag.StringVar(&broker, "b", "stream.routeviews.org", "Kafka broker")
	flag.Parse()

	seeds := []string{broker}
	// One client can both produce and consume!
	// Consuming can either be direct (no consumer group), or through a group. Below, we use a group.
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
	//	kgo.ConsumerGroup("this-is-a-test-2"),
		kgo.ConsumeTopics("routeviews.route-views2.31019.bmp_raw"),
	)
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	ctx := context.Background()

	// 2.) Consuming messages from a topic
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context closed for kafka consumer with topics %v", "foo")
			return
		default:
		}

		fmt.Println("Entering loop")
		fetches := cl.PollFetches(ctx)
		fmt.Println("Polled")
		if fetches.IsClientClosed() {
			return
		}
		fetches.EachError(func(t string, p int32, err error) {
			die("fetch err topic %s partition %d: %v", t, p, err)
		})

		// We can iterate through a record iterator...
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			str := hex.EncodeToString(record.Value)

			data := record.Value

			fmt.Println(str)

			message, error := obmp.ParseHeader(data)
			if error != nil {
				fmt.Println(error.Error())
				log.Fatalf("Error parsing: ", data)
			}

			fmt.Println("Message:")
			fmt.Printf("%s\n", makePretty(message))

			// Now do the BMP part!
			bmp := message.BMPMessage
			outbmp, err := gobmp.UnmarshalCommonHeader(bmp[0:6])
			if err != nil {
				log.Fatalf("Failed decoding BMPheader")
			}
			fmt.Printf("%+v\n", outbmp)

			perpeerdata := bmp[6:6+42+1]
			fmt.Println("--> ", len(perpeerdata), perpeerdata)
			perpeerheader, err := gobmp.UnmarshalPerPeerHeader(perpeerdata)
			if err != nil {
				log.Fatalf("Failed decoding BMPPerPeerHeader")
			}
			fmt.Printf("%+v\n", perpeerheader)


			switch outbmp.MessageType {
			case gobmp.PeerDownMsg:
				fmt.Println(bmp[48:])
				tmp, err := gobmp.UnmarshalPeerDownMessage(bmp[48:])
				if err != nil {
					log.Fatalf("Failed decoding PeerDownMessage")
				}
				fmt.Printf("%+v\n", tmp)
			}

		}
	}
}

