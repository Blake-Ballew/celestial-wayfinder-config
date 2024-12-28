package main

import (
	"fmt"
	"os"

	"go.bug.st/serial"
)

var menuStack *MenuStack

func main() {
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

	if len(os.Args) < 2 {
		interactiveMode()
		return
	}

	if os.Args[1] == "help" {
		printUsage()
		return
	}

}

func printUsage() {
	fmt.Println("Usage: wayfinder-config [command]")
}

func interactiveMode() {

	// Return if platform is Windows
	if os.Getenv("GOOS") == "windows" {
		fmt.Println("Interactive mode is not supported on Windows")
		os.Exit(0)
	}

	menuStack = NewMenuStack()

	homePage := NewMenuPage("Celestial Wayfinder Configuration")

	homePage.OnDisplay = func(thisPage *MenuPage) {
		thisPage.ClearMenuSelections().
			AssignMenuSelection("rpc-channel", "Select RPC Channel", func(key string) (int, error) {
				return WINDOW_SELECT, nil
			}).AssignAdjacentMenu("rpc-channel", NewRpcChannelMenu)

		if CurrentRpcChannel != nil {
			thisPage.AssignMenuSelection("rpc-functions", "RPC Functions", func(key string) (int, error) {
				return WINDOW_SELECT, nil
			}).AssignAdjacentMenu("rpc-functions", GenerateRpcFunctionsMenu)
		}

		thisPage.AssignMenuSelection("quit", "Quit", func(key string) (int, error) {
			if CurrentRpcChannel != nil {
				CurrentRpcChannel.Close()
			}
			os.Exit(0)
			return 0, nil
		})

	}

	menuStack.Push(homePage)

	for {
		menu := menuStack.MenuPageStack.Front()
		if menu == nil {
			os.Exit(0)
		}

		if menu.Value.(*MenuPage).OnDisplay != nil {
			menu.Value.(*MenuPage).OnDisplay(menu.Value.(*MenuPage))
		}
		choice := menu.Value.(*MenuPage).MenuObject.Display()
		selectionObj := menu.Value.(*MenuPage).SelectionMap[choice]

		result, err := selectionObj.ExecuteMenuSelection(choice)
		if err != nil {
			fmt.Println(err)
		}

		if result > 0 {
			switch result {
			case WINDOW_BACK:
				menuStack.Pop()
			case WINDOW_SELECT:
				if menu.Value.(*MenuPage).AdjacentMenu[choice] != nil {
					menuStack.Push(menu.Value.(*MenuPage).AdjacentMenu[choice]())
				}
			}
		}
	}

}
