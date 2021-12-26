package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	// "os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/jlaffaye/ftp"
)

var (
	device            = "en0"
	snapshotLen int32 = 1024
	err         error
	timeout     = 30 * time.Second
	handle      *pcap.Handle
	buffer      gopacket.SerializeBuffer
	options     gopacket.SerializeOptions
)

func attackFTP(targetIP string, LoginStall bool) {
	c, err := ftp.Dial(targetIP, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	// If logging in is relevant, edit LoginStall to true.
	if LoginStall {

		err = c.Login("login", "password")
		if err != nil {
			log.Fatal(err)
		}

		if err := c.Quit(); err != nil {
			log.Fatal(err)
		}
	}
}

func attackHTTP(httpReq string) {
	// Get request URL
	reqURL, _ := url.Parse(httpReq)

	// Create request body
	reqBody := ioutil.NopCloser(strings.NewReader(`
					{
						"test":"test me"
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

	// Check for response error
	if err != nil {
		log.Fatal("Error:", err)
	}

	// Read response body
	data, _ := ioutil.ReadAll(res.Body)

	// Close response body
	err = res.Body.Close()
	if err != nil {
		return
	}

	// Print response status and body
	fmt.Printf("status: %d\n", res.StatusCode)
	fmt.Printf("body: %s\n", data)

}

func selectPort(ipTarget string, portTarget string) {

	switch portTarget {
	case "ftp":
		attackFTP(ipTarget, false)
	case "http":
		attackHTTP(ipTarget)
	}

}

func spoofIP() {

	myip := net.ParseIP("150.0.0.1:8000")
	addrspoof := &net.IPAddr{myip, ""}

	TransportExp := &http.Transport{
		DialContext: (&net.Dialer{
			LocalAddr: addrspoof,
		}).DialContext,
	}
	client := &http.Client{
		Transport: TransportExp,
	}
	//	targetIP2 := net.ParseIP("128.0.0.1")
	//	client, err := ioutil.ReadAll(client)

	req, err := http.NewRequest("GET", "0.0.0.0", nil)
	if err != nil {
		return
	}
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	resp, err := client.Do(req)
	// TransportExp = buffer.Bytes()

	fmt.Println(resp.Body)
}

func main() {

	var targetIP = os.Args[0]
	var targetPort = os.Args[1]
	selectPort(targetIP, targetPort)

}
