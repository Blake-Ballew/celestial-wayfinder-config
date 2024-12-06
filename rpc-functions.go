package main

import (
	"encoding/json"
	"os"
)

func ExecRpc(rpcData map[string]interface{}) (map[string]interface{}, error) {
	return CallRpc(rpcData)
}

func ExecRpcFromJsonFile(fileName string) (map[string]interface{}, error) {
	var rpcData map[string]interface{} = make(map[string]interface{})
	data, err := os.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &rpcData)
	if err != nil {
		return nil, err
	}

	return ExecRpc(rpcData)
}

func AddSavedMessage(message string) (map[string]interface{}, error) {
	var rpcData map[string]interface{} = make(map[string]interface{})
	rpcData["F"] = "AddSavedMessage"
	rpcData["message"] = message

	result, err := CallRpc(rpcData)

	return result, err
}

func AddSavedMessages(messages []string) (map[string]interface{}, error) {
	var rpcData map[string]interface{} = make(map[string]interface{})
	rpcData["F"] = "AddSavedMessages"
	rpcData["messages"] = messages

	result, err := CallRpc(rpcData)

	return result, err
}

func DeleteSavedMessage(idx int) (map[string]interface{}, error) {
	var rpcData map[string]interface{} = make(map[string]interface{})
	rpcData["F"] = "DeleteSavedMessage"
	rpcData["idx"] = idx

	result, err := CallRpc(rpcData)

	return result, err
}

func DeleteSavedMessages() (map[string]interface{}, error) {
	var rpcData map[string]interface{} = make(map[string]interface{})
	rpcData["F"] = "DeleteSavedMessages"

	result, err := CallRpc(rpcData)

	return result, err
}

func GetSavedMessage(idx int) (map[string]interface{}, error) {
	var rpcData map[string]interface{} = make(map[string]interface{})
	rpcData["F"] = "GetSavedMessage"
	rpcData["idx"] = idx

	result, err := CallRpc(rpcData)

	return result, err
}

func GetSavedMessages() (map[string]interface{}, error) {
	var rpcData map[string]interface{} = make(map[string]interface{})
	rpcData["F"] = "GetSavedMessages"

	result, err := CallRpc(rpcData)

	return result, err
}

func UpdateSavedMessage(idx int, message string) (map[string]interface{}, error) {
	var rpcData map[string]interface{} = make(map[string]interface{})
	rpcData["F"] = "UpdateSavedMessage"
	rpcData["idx"] = idx
	rpcData["message"] = message

	result, err := CallRpc(rpcData)

	return result, err
}
