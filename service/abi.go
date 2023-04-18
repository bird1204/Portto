package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetContractAbiFromFile(client *ethclient.Client, contractAddress common.Address) (*abi.ABI, error) {
	filename := "abiData/" + contractAddress.Hex() + ".json"

	// Check whether the system has ABI data
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		saveContractAbiFromEtherscan(client, contractAddress)
	}

	// Read ABI data from file
	abiJson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Create an ABI object from the ABI data
	contractAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		panic("Error creating ABI object: " + err.Error())
	}

	if err != nil {
		panic(fmt.Errorf("error creating ABI object: %v", err))
	}

	return &contractAbi, nil
}

func saveContractAbiFromEtherscan(client *ethclient.Client, contractAddress common.Address) {
	// Get contract address and ABI
	response, err := http.Get("https://api.etherscan.io/api?module=contract&action=getabi&address=" + contractAddress.Hex())
	if err != nil {
		panic(err.Error())
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}

	var abiData map[string]interface{}
	err = json.Unmarshal(body, &abiData)
	if err != nil {
		panic("Error decoding ABI data: " + err.Error())
	}

	if abiData["result"] == "Contract source code not verified" {
		fmt.Print(abiData["result"])
		return
	}

	filepath := "abiData/" + contractAddress.Hex() + ".json"
	err = ioutil.WriteFile(filepath, []byte(abiData["result"].(string)), 0644)
	if err != nil {
		panic("Error decoding ABI data: " + err.Error())
	}

	fmt.Println("ABI data saved to", filepath)
}
