package main

import (
	//	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/cbocovic/lavinia"
	"io"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var startaddr string

	//set up flags
	numPtr := flag.Int("num", 2, "the size of the DHT you wish to test")
	startPtr := flag.Int("start", 8888, "port number to start from")

	flag.Parse()
	num := *numPtr
	start := *startPtr
	fmt.Printf("Joining %d servers starting at %d!\n", num, start)

	list := make([]*lavinia.LaviniaServer, num)
	if start == 8888 {

		me := new(lavinia.LaviniaServer)
		startaddr = fmt.Sprintf("127.0.0.1:%d", start)
		me = lavinia.Create(fmt.Sprintf("/home/bocovich/FS/%d/", start), startaddr)
		list[0] = me
		lavinia.Publish("/home/bocovich/out", "127.0.0.1:8888")
	} else {
		me := new(lavinia.LaviniaServer)
		startaddr = fmt.Sprintf("127.0.0.1:%d", start)
		me = lavinia.Join(fmt.Sprintf("/home/bocovich/FS/%d/", start), startaddr, "127.0.0.1:8888")
		list[0] = me
	}

	for i := 1; i < num; i++ {
		//join node to network or start a new network
		time.Sleep(time.Second)

		addr := fmt.Sprintf("127.0.0.1:%d", start+i)
		node := lavinia.Join(fmt.Sprintf("/home/bocovich/FS/%d/", start+i), addr, startaddr)
		list[i] = node
		fmt.Printf("Joined server: %s.\n", addr)
	}
	//block until receive input
Loop:
	for {
		var cmd string
		var port int
		_, err := fmt.Scan(&cmd)
		switch {
		case cmd == "info":
			//print out successors and predecessors
			fmt.Printf("Node\t\t Successor\t\t Predecessor\n")
			for _, node := range list {
				fmt.Printf("%s\n", node.Info())
			}
		case cmd == "fingers":
			//print out finger table
			fmt.Printf("Enter port of desired node: ")
			fmt.Scan(&port)
			if port-start >= 0 && port-start < len(list) {
				node := list[port-start]
				fmt.Printf("\n%s", node.ShowFingers())
			}
		case cmd == "succ":
			//print out successor list
			fmt.Printf("Enter port of desired node: ")
			fmt.Scan(&port)
			if port-start >= 0 && port-start < len(list) {
				node := list[port-start]
				fmt.Printf("\n%s", node.ShowSucc())
			}
		case err == io.EOF:
			break Loop
		}

	}
	for _, node := range list {
		node.Finalize()
	}

}
