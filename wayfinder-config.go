package main

import (
	"fmt"
	"os"

	"github.com/nexidian/gocliselect"
	"github.com/vmihailenco/msgpack/v5"
	"go.bug.st/serial"
)

func main() {
	fmt.Println("Hello, World!")
	fmt.Println("msgpack version:", msgpack.Version())
	ports, err := serial.GetPortsList()
	if err != nil {
		fmt.Println(err)
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
	} else {
		for _, port := range ports {
			fmt.Printf("Found port: %v\n", port)
		}
	}

	if len(os.Args) < 2 || os.Args[1] == "help" {
		interactiveMode()
		return
	}

}

func printUsage() {
	fmt.Println("Usage: wayfinder-config [command]")
}

func interactiveMode() {
	menu := gocliselect.NewMenu("Celestial Wayfinder Utils")

	menu.AddItem("RPC", "rpc")
	menu.AddItem("Quit", "quit")

	choice := menu.Display()

	switch choice {
	case "rpc":
		// rpcMode()
	case "quit":
		os.Exit(0)
	}
}
