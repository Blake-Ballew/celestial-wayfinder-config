package main

import (
	"fmt"
	"strconv"
)

func NewRpcChannelMenu() *MenuPage {
	RpcChannelMenu := NewMenuPage("Select RPC Channel")
	RpcChannelMenu.AssignMenuSelection("tcp", "TCP", func(key string) (int, error) {
		return WINDOW_SELECT, nil
	}).AssignMenuSelection("serial", "Serial", func(key string) (int, error) {
		return WINDOW_SELECT, nil
	}).AssignMenuSelection("back", "Back", func(key string) (int, error) {
		return WINDOW_BACK, nil
	}).AssignAdjacentMenu("tcp", NewTcpConnectionMenu).
		AssignAdjacentMenu("serial", NewSerialConnectionMenu)

	return RpcChannelMenu
}

func NewTcpConnectionMenu() *MenuPage {
	menupg := NewMenuPage("Select TCP Connection")

	menupg.OnDisplay = func(thisPage *MenuPage) {
		// Clear all selections
		menupg.ClearMenuSelections()

		menupg.AssignMenuSelection("refresh", "Refresh", func(key string) (int, error) {
			FindTcpChannels(1000)
			return 0, nil
		})

		menupg.AssignMenuSelection("back", "Back", func(key string) (int, error) {
			return WINDOW_BACK, nil
		})

		for idx, channel := range BroadcastedTcpChannels {
			menupg.AssignMenuSelection(fmt.Sprint(idx), channel.IpAddress+":"+fmt.Sprint(channel.Port), func(key string) (int, error) {
				tcpIdx, err := strconv.Atoi(key)
				if err != nil {
					return 0, err
				}

				ConnectToTcpChannel(BroadcastedTcpChannels[tcpIdx])
				return WINDOW_BACK, nil
			})
		}
	}

	return menupg
}

func NewSerialConnectionMenu() *MenuPage {
	menupg := NewMenuPage("Select Serial Port")

	menupg.OnDisplay = func(thisPage *MenuPage) {
		// Clear all selections
		menupg.ClearMenuSelections()

		menupg.AssignMenuSelection("refresh", "Refresh", func(key string) (int, error) {
			RefreshSerialChannels()
			return 0, nil
		})

		menupg.AssignMenuSelection("back", "Back", func(key string) (int, error) {
			return WINDOW_BACK, nil
		})

		for idx, port := range SerialChannels {
			menupg.AssignMenuSelection(fmt.Sprint(idx), port, func(key string) (int, error) {
				portIdx, err := strconv.Atoi(key)
				if err != nil {
					return 0, err
				}

				SelectSerialPort(SerialChannels[portIdx])
				return WINDOW_BACK, nil
			})
		}
	}

	return menupg
}
