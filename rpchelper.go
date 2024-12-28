package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"go.bug.st/serial"
)

// Possibly not needed
const (
	RPC_CHANNEL_NONE   = 0
	RPC_CHANNEL_SERIAL = 1
	RPC_CHANNEL_TCP    = 2
)

const (
	RPC_BROADCAST_PORT = ":54789"
)

var CurrentRpcChannelType int = RPC_CHANNEL_NONE
var CurrentRpcChannel RpcChannel = nil

type RpcChannel interface {
	CallRpcChannel(rpcData map[string]interface{}) (map[string]interface{}, error)
	Open() error
	Close()
}

func CallRpc(rpcData map[string]interface{}) (map[string]interface{}, error) {
	if CurrentRpcChannel != nil {
		return CurrentRpcChannel.CallRpcChannel(rpcData)
	}

	return nil, fmt.Errorf("no rpc channel open")
}

func CloseRpcChannel() {
	if CurrentRpcChannel != nil {
		CurrentRpcChannel.Close()
	}

	CurrentRpcChannelType = RPC_CHANNEL_NONE
	CurrentRpcChannel = nil
}

type TcpConnectionInfo struct {
	IpAddress string
	Port      int
}

type TcpChannel struct {
	connectionInfo TcpConnectionInfo
	connection     *net.TCPConn
}

func (c *TcpChannel) CallRpcChannel(rpcData map[string]interface{}) (map[string]interface{}, error) {
	if c.connection == nil {
		err := c.Open()
		if err != nil {
			return nil, err
		}
	}

	data, err := msgpack.Marshal(rpcData)
	if err != nil {
		return nil, err
	}

	c.connection.SetWriteDeadline(time.Now().Add(time.Second * 5))

	_, err = c.connection.Write(data)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024)
	timeoutMs := 3000
	startTime := time.Now()
	bytesReceived := 0

	// Wait for bytes and attempt to decode the message pack
	for time.Since(startTime) < time.Millisecond*time.Duration(timeoutMs) {
		// Set deadline to loop time remaining
		timeRemaining := int64(timeoutMs) - time.Since(startTime).Milliseconds()
		c.connection.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(timeRemaining)))

		numBytes, err := c.connection.Read(buffer[bytesReceived:])
		if err != nil {
			return nil, err
		}

		bytesReceived += numBytes

		// if numBytes == 0 {
		// 	return nil, fmt.Errorf("tcp response was empty")
		// }

		// fmt.Printf("Received %d bytes\n", numBytes)

		// Print bytes
		// for i := 0; i < bytesReceived; i++ {
		// 	fmt.Printf("%02x ", buffer[i])
		// }
		// fmt.Println()

		response := make(map[string]interface{})
		err = msgpack.Unmarshal(buffer[:bytesReceived], &response)
		if err != nil {
			continue
		}

		return response, nil
	}

	return nil, fmt.Errorf("timed out waiting for tcp response")

	// c.connection.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(timeoutMs)))
	// numBytes, err := c.connection.Read(buffer)
	// if err != nil {
	// 	return nil, err
	// }

	// if numBytes == 0 {
	// 	return nil, fmt.Errorf("no response from tcp channel")
	// }

	// // Print bytes
	// for i := 0; i < numBytes; i++ {
	// 	fmt.Printf("%02x ", buffer[i])
	// }
	// fmt.Println()

	// response := make(map[string]interface{})
	// err = msgpack.Unmarshal(buffer[:numBytes], &response)
	// if err != nil {
	// 	return nil, err
	// }

	// return response, nil
}

func (c *TcpChannel) Open() error {
	connection, err := net.DialTCP("tcp4", &net.TCPAddr{
		IP: net.ParseIP("0.0.0.0"),
		// Port: c.connectionInfo.Port,
		Port: 0,
	}, &net.TCPAddr{
		IP:   net.ParseIP(c.connectionInfo.IpAddress),
		Port: c.connectionInfo.Port,
	})
	if err != nil {
		return err
	}

	c.connection = connection
	return nil
}

func (c *TcpChannel) Close() {
	if c.connection != nil {
		connFile, _ := c.connection.File()
		if connFile != nil {
			connFile.Close()
		}
		c.connection.Close()
	}
}

var BroadcastedTcpChannels []TcpConnectionInfo

var SerialChannels []string
var activeSerialPort string = ""

func FindTcpChannels(timeoutMs int) {
	fmt.Println("Listening for broadcast on:", RPC_BROADCAST_PORT)

	conn, err := net.ListenPacket("udp4", RPC_BROADCAST_PORT)
	if err != nil {
		fmt.Println("Error listening for broadcast:", err)
		return
	}
	defer conn.Close()

	start := time.Now()

	for {
		timeRemaining := int64(timeoutMs) - time.Since(start).Milliseconds()
		if timeRemaining <= 0 {
			fmt.Println()
			break
		}

		fmt.Print(".")

		// Set deadline to loop time remaining
		conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(100)))
		buffer := make([]byte, 1024)
		numBytes, _, err := conn.ReadFrom(buffer)

		if numBytes == 0 {
			continue
		} else {
			fmt.Println("Packet received")
		}

		if err != nil {
			fmt.Println()
			fmt.Println("Error reading from broadcast:", err)
			continue
			// return
		}

		tcpChannel := TcpConnectionInfo{
			IpAddress: "",
			Port:      0,
		}

		err = msgpack.Unmarshal(buffer, &tcpChannel)
		if err != nil {
			fmt.Println("Error reading into TcpChannel:", err)
			continue
		}

		if tcpChannel.Port != 0 && tcpChannel.IpAddress != "" {
			duplicateFound := false
			for _, channel := range BroadcastedTcpChannels {
				if channel.IpAddress == tcpChannel.IpAddress && channel.Port == tcpChannel.Port {
					fmt.Println("Duplicate TCP channel:", tcpChannel.IpAddress+":"+fmt.Sprint(tcpChannel.Port))
					duplicateFound = true
					break
				}
			}

			if !duplicateFound {
				BroadcastedTcpChannels = append(BroadcastedTcpChannels, tcpChannel)
			}

		} else {
			fmt.Println("Invalid TCP channel:", tcpChannel.IpAddress+":"+fmt.Sprint(tcpChannel.Port))
		}
	}
}

func ConnectToTcpChannel(channel TcpConnectionInfo) {

	tcpchannel := TcpChannel{
		connectionInfo: channel,
	}

	err := tcpchannel.Open()

	if err != nil {
		fmt.Println("Error connecting to TCP channel:", err)
		CurrentRpcChannelType = RPC_CHANNEL_NONE
		return
	}

	fmt.Println()
	fmt.Println("Connected to TCP channel:", channel.IpAddress+":"+fmt.Sprint(channel.Port))

	if CurrentRpcChannel != nil {
		CurrentRpcChannel.Close()
	}

	CurrentRpcChannelType = RPC_CHANNEL_TCP
	CurrentRpcChannel = &tcpchannel
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
	CurrentRpcChannelType = RPC_CHANNEL_SERIAL
}

// RPC Functions

func GenerateRpcFunctionsMenu() *MenuPage {
	functionMenu := NewMenuPage("RPC Functions").
		AssignMenuSelection("exec-json-file", "Send Hello World", func(key string) (int, error) {

			// prompt user for file name
			// reader := bufio.NewReader(os.Stdin)
			// fmt.Print("Enter file name: ")
			// fileName, err := reader.ReadString('\n')
			// if err != nil {
			// 	return 0, err
			// }

			// Test code
			testRpcCall := make(map[string]interface{})
			testRpcCall["F"] = "AddSavedMessage"
			testRpcCall["message"] = "Hello World"

			result, err := ExecRpc(testRpcCall)

			//result, err := ExecRpcFromJsonFile(fileName)
			if err != nil {
				return 0, err
			}

			jsonData, err := json.Marshal(result)
			if err != nil {
				return 0, err
			}

			fmt.Println()
			fmt.Println("Result:")
			fmt.Println(string(jsonData))
			fmt.Println()

			return 0, nil
		}).AssignMenuSelection("delete-message", "Delete Messages", func(key string) (int, error) {
		messagePacket, err := GetSavedMessages()
		if err != nil {
			fmt.Println("Error getting saved messages:", err)
			return 0, err
		}

		maxIdx := 0

		if messages, ok := messagePacket["messages"]; ok {
			if messagesSlice, ok := messages.([]interface{}); ok {
				maxIdx = len(messagesSlice) - 1
				for idx, message := range messagesSlice {
					if messageStr, ok := message.(string); ok {
						fmt.Println(idx, ": ", messageStr)
					}
				}

				fmt.Print("Enter message index to delete: ")
				reader := bufio.NewReader(os.Stdin)
				key, err := reader.ReadString('\n')
				if err != nil {
					return 0, err
				}

				key = strings.TrimSpace(key)
				if key == "" {
					return 0, nil
				}

				idx, err := strconv.Atoi(key)
				if err != nil {
					return 0, err
				}

				if idx < 0 || idx > maxIdx {
					fmt.Println("Invalid index")
					return 0, nil
				}

				_, err = DeleteSavedMessage(idx)
				if err != nil {
					return 0, err
				}
			}
		}

		return 0, nil
	}).AssignMenuSelection("add-message", "Add Message", func(key string) (int, error) {
		fmt.Print("Enter message: ")
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}

		result, err := AddSavedMessage(message)
		if err != nil {
			return 0, err
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			return 0, err
		}

		fmt.Println()
		fmt.Println("Result:")
		fmt.Println(string(jsonData))
		fmt.Println()

		return 0, nil

	}).AssignMenuSelection("back", "Back", func(key string) (int, error) {
		return WINDOW_BACK, nil
	})

	return functionMenu
}
