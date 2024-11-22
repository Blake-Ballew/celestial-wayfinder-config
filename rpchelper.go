package main

import (
	"fmt"
	"net"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"go.bug.st/serial"
)

const (
	RPC_CHANNEL_NONE   = 0
	RPC_CHANNEL_SERIAL = 1
	RPC_CHANNEL_TCP    = 2
)

const (
	RPC_BROADCAST_PORT = 12345
)

var CurrentRpcChannel int = RPC_CHANNEL_NONE

type TcpChannel struct {
	IpAddress string
	Port      int
}

var BroadcastedTcpChannels []TcpChannel
var activeTcpConnection *net.TCPConn = nil

var SerialChannels []string
var activeSerialPort string = ""

func FindTcpChannels(timeoutMs int) {
	pc, err := net.ListenPacket("udp", ":"+fmt.Sprint(RPC_BROADCAST_PORT))
	if err != nil {
		fmt.Println("Error listening for broadcast:", err)
		return
	}
	defer pc.Close()

	start := time.Now()

	for {
		if time.Since(start).Milliseconds() > int64(timeoutMs) {
			fmt.Println()
			break
		}

		fmt.Print(".")

		pc.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(100)))
		buffer := make([]byte, 1024)
		numBytes, _, err := pc.ReadFrom(buffer)

		if numBytes == 0 {
			continue
		}

		if err != nil {
			fmt.Println()
			fmt.Println("Error reading from broadcast:", err)
			return
		}

		tcpChannel := TcpChannel{
			IpAddress: "",
			Port:      0,
		}

		err = msgpack.Unmarshal(buffer, &tcpChannel)
		if err != nil {
			continue
		}

		if tcpChannel.Port != 0 && tcpChannel.IpAddress != "" {
			BroadcastedTcpChannels = append(BroadcastedTcpChannels, tcpChannel)
		}
	}
}

func ConnectToTcpChannel(channel TcpChannel) {
	connection, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   net.ParseIP(channel.IpAddress),
		Port: channel.Port,
	})

	if err != nil {
		fmt.Println("Error connecting to TCP channel:", err)
		CurrentRpcChannel = RPC_CHANNEL_NONE
		return
	}

	if activeTcpConnection != nil {
		activeTcpConnection.Close()
	}

	CurrentRpcChannel = RPC_CHANNEL_TCP
	activeTcpConnection = connection
}

func RefreshSerialChannels() {
	channels, err := serial.GetPortsList()
	if err != nil {
		fmt.Println(err)
	}
	SerialChannels = channels
}

func SelectSerialPort(port string) {
	activeSerialPort = port
	CurrentRpcChannel = RPC_CHANNEL_SERIAL
}
