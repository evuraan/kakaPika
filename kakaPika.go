package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type wemoStruct struct {
	name    string
	port    int
	connPtr *net.Conn
}

const (
	wemoTurnOffThingy = `POST /upnp/control/basicevent1 HTTP/1.0
Content-type: text/xml; charset="utf-8"
SOAPACTION: "urn:Belkin:service:basicevent:1#SetBinaryState"
Content-Length: 299

<?xml version="1.0" encoding="utf-8"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body><u:SetBinaryState xmlns:u="urn:Belkin:service:basicevent:1"><BinaryState>0</BinaryState></u:SetBinaryState></s:Body></s:Envelope>`

	wemoTurnONThingy = `POST /upnp/control/basicevent1 HTTP/1.0
Content-type: text/xml; charset="utf-8"
SOAPACTION: "urn:Belkin:service:basicevent:1#SetBinaryState"
Content-Length: 299

<?xml version="1.0" encoding="utf-8"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body><u:SetBinaryState xmlns:u="urn:Belkin:service:basicevent:1"><BinaryState>1</BinaryState></u:SetBinaryState></s:Body></s:Envelope>`

	wemoStatusThingy = `POST /upnp/control/basicevent1 HTTP/1.0
Content-type: text/xml; charset="utf-8"
SOAPACTION: "urn:Belkin:service:basicevent:1#GetBinaryState"
Content-Length: 299

<?xml version="1.0" encoding="utf-8"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body><u:GetBinaryState xmlns:u="urn:Belkin:service:basicevent:1"><BinaryState>1</BinaryState></u:GetBinaryState></s:Body></s:Envelope>`
)

func findWemoPort(device string) *wemoStruct {
	// pick between 49152 and 49153
	wemoDevice := wemoStruct{}
	x := false
	wg := sync.WaitGroup{}
	checkThese := []int{49152, 49153}
	for i := range checkThese {
		wg.Add(1)
		port := checkThese[i]
		connectTo := fmt.Sprintf("%s:%d", device, port)

		go func(wd *wemoStruct) {
			defer wg.Done()
			conn, err := net.DialTimeout("tcp", connectTo, 5*time.Second)
			if err != nil {
				return
			}
			wd.name = device
			wd.port = port
			wd.connPtr = &conn
			x = true

		}(&wemoDevice)
	}
	wg.Wait()

	if x {
		return &wemoDevice
	}
	return nil

}

func operateWemo(device, action string, onChan chan string) {

	devPtr := findWemoPort(device)
	if devPtr == nil {
		return
	}
	var thingy string
	switch action {
	case "on":
		thingy = wemoTurnONThingy
	case "off":
		thingy = wemoTurnOffThingy
	case "stat":
		thingy = wemoStatusThingy
	default:
		return
	}

	conn := *devPtr.connPtr
	defer conn.Close()

	n, err := conn.Write([]byte(thingy))
	if err != nil {
		return
	}

	if n < 1 {
		return
	}

	readInto := make([]byte, 8192)

	// After POST, read until EOF - if we can.
	for true {
		n, err = conn.Read(readInto)
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		if n < 1 {
			return
		}
	}

	if bytes.Contains(readInto, []byte("<BinaryState>0</BinaryState>")) {
		onChan <- "off"
		return
	} else if bytes.Contains(readInto, []byte("<BinaryState>1</BinaryState>")) {
		onChan <- "on"
		return
	}

	// because readInto for wemo devices are മൊണഞ്ഞത്,
	// some recursion, eh?
	operateWemo(device, "stat", onChan)
	// onChan <- "Err"
}
