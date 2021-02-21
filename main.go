package main

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	version = "kakaPika ver 1.0b  Copyright (c) Evuraan <evuraan@gmail.com>"
	// encoded {"system":{"set_relay_state":{"state":1}}}
	payloadOn = "AAAAKtDygfiL/5r31e+UtsWg1Iv5nPCR6LfEsNGlwOLYo4HyhueT9tTu36Lfog=="

	//encoded {"system":{"set_relay_state":{"state":0}}}
	payloadOff = "AAAAKtDygfiL/5r31e+UtsWg1Iv5nPCR6LfEsNGlwOLYo4HyhueT9tTu3qPeow=="

	// encoded { "system":{ "get_sysinfo":null } }
	payloadQuery = "AAAAI9Dw0qHYq9+61/XPtJS20bTAn+yV5o/hh+jK8J7rh+vLtpbr"

	// the encoded request { "emeter":{ "get_realtime":null } }
	payloadEmeter = "AAAAJNDw0rfav8uu3P7Ev5+92r/LlOaD4o76k/6buYPtmPSYuMXlmA=="
	tplugON       = `relay_state":1`
	tplugOFF      = `relay_state":0`
	tplugOpOK     = `{"err_code":0}`
	timeOut       = 5
)

var (
	device        = ""
	cmd           = ""
	tplugCommands = map[string]string{"on": payloadOn, "off": payloadOff, "stat": payloadQuery}
)

func main() {
	parseArgs()
	fmt.Printf("device: %s, cmd: %s\n", device, cmd)

	switch cmd {
	case "stat", "on", "off":
		fmt.Printf("%s: ", cmd)
		smartPlugOp(device, cmd)

	default:
		fmt.Fprintf(os.Stderr, "Err: Unknown operative %v\n", cmd)
		usage()
		os.Exit(1)
	}

	// fmt.Printf("Reply: %q\n", tram)
	// fmt.Printf("netConn type: %T\n", netConn)
}

func decrypt(encrypted []byte) string {
	out := ""
	key := byte(171)
	for i := range encrypted {
		c := encrypted[i]
		x := c ^ key
		key = c
		out = fmt.Sprintf("%s%c", out, x)
	}
	return out
}

func usage() {
	fmt.Println("Usage: ")
	fmt.Println("  -h  --help           print this usage and exit")
	fmt.Println("  -d  --device         smartPlug hostname or address")
	fmt.Println("  -c  --cmd            cmd to run: [on,off,stat]")
	fmt.Println("  -v  --version        print version information and exit")
}

func parseArgs() {
	argc := len(os.Args)
	if argc < 2 {
		usage()
		os.Exit(1)
	}

	for i := range os.Args {
		arg := os.Args[i]
		arg = strings.ToLower(arg)
		if strings.Contains(arg, "help") || arg == "h" || arg == "--h" || arg == "-h" || arg == "?" {
			usage()
			os.Exit(0)
		}
		if strings.Contains(arg, "version") || arg == "v" || arg == "--v" || arg == "-v" {
			fmt.Println("Version:", version)
			os.Exit(0)
		}
		// look for -d or --device
		if strings.Contains(arg, "dev") || arg == "--d" || arg == "-d" || arg == "d" {
			nextArg := i + 1
			if argc > nextArg {
				device = os.Args[nextArg]
			} else {
				invalidUsage("device")
			}
		}

		// look for -c or --cmd
		if strings.Contains(arg, "cmd") || arg == "--c" || arg == "-c" || arg == "c" {
			nextArg := i + 1
			if argc > nextArg {
				cmd = os.Args[nextArg]
			} else {
				invalidUsage("cmd")
			}
		}
	}

	// by now, cmd and device must be filled.
	if len(device) < 1 || len(cmd) < 1 {
		fmt.Println("Error parsing cmd args")
		usage()
		os.Exit(1)
	}

}

func invalidUsage(cmdArg string) {
	fmt.Printf("Invalid Usage. Error parsing %s\n", cmdArg)
	usage()
	os.Exit(1)
}

func checkerr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func smartPlugOp(device, action string) {

	var verb string
	switch action {
	case "on":
		verb = action
	case "off":
		verb = action
	default:
		verb = "stat"
	}

	textChannel := make(chan string, 1)
	go operateTplug(device, verb, textChannel)
	go operateWemo(device, verb, textChannel)
	select {
	case opStat := <-textChannel:
		fmt.Println(opStat)
	case <-time.After(6 * time.Second):
		fmt.Println("Timeout!")
	}
}

func operateTplug(device string, action string, onChan chan string) bool {

	cmd, ok := tplugCommands[action]
	if !ok {
		return false
	}
	connectTo := fmt.Sprintf("%s:9999", device)
	netConn, err := net.DialTimeout("tcp", connectTo, timeOut*time.Second)

	if err != nil {
		return false
	}
	defer netConn.Close()

	query, err := base64.StdEncoding.DecodeString(cmd)
	if err != nil {
		return false
	}
	sent, err := netConn.Write(query)
	if err != nil {
		return false
	}

	if sent < 1 {
		return false
	}

	readInto := make([]byte, 8192)
	bytesRead, err := netConn.Read(readInto)
	if err != nil {
		return false
	}

	if bytesRead < 1 {
		fmt.Println("Err 331.5")
		return false
	}
	tram := readInto[:bytesRead]
	tram = tram[4:]
	textStatus := decrypt(tram)
	switch action {
	case "stat":
		if strings.Contains(textStatus, tplugON) {
			onChan <- "on"
			return true
		} else if strings.Contains(textStatus, tplugOFF) {
			onChan <- "off"
			return true
		}
	case "on", "off":
		if strings.Contains(textStatus, tplugOpOK) {
			// onChan <- "OK"
			onChan <- action
			return true
		}
	default:
		fmt.Printf("Unknown action: %v\n", action)
		return false

	}

	return false
}
