package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// Car describes basic details of what makes up a car
type Device struct {
	ID            string `json:"id"`
	Uuid          string `json:"uuid"`
	Servicenumber string `json:"servicenumber"`
	Version       string `json:"version"`
	Date          string `json:"date"` // updateing date
}

type QueryResult struct {
	Key    string `json:"Key"`
	Record *Device
}

// InitLedger adds a base set of cars to the ledger (Do nothing)
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreateDevice(ctx contractapi.TransactionContextInterface, deviceNumber string, id string, uuid string, servicenumber string, version string, date string) error {
	device := Device{
		ID:            id,
		Uuid:          uuid,
		Servicenumber: servicenumber,
		Version:       version,
		Date:          date,
	}

	deviceAsBytes, _ := json.Marshal(device)

	return ctx.GetStub().PutState(deviceNumber, deviceAsBytes)
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryDevice(ctx contractapi.TransactionContextInterface, deviceNumber string) (*Device, error) {
	deviceAsBytes, err := ctx.GetStub().GetState(deviceNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if deviceAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", deviceNumber)
	}

	device := new(Device)
	_ = json.Unmarshal(deviceAsBytes, device)

	return device, nil
}

func (s *SmartContract) QueryAllDevices(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		device := new(Device)
		_ = json.Unmarshal(queryResponse.Value, device)

		queryResult := QueryResult{Key: queryResponse.Key, Record: device}
		results = append(results, queryResult)
	}

	return results, nil
}

// Delete the device module when it should be deleted
func (s *SmartContract) DeleteDevice(ctx contractapi.TransactionContextInterface, deviceNumber string) error {
	exists, err := s.QueryDevice(ctx, deviceNumber)

	if err != nil {
		return err
	}

	if exists == nil {
		return fmt.Errorf("The Device %s does not exist", deviceNumber)
	}

	return ctx.GetStub().DelState(deviceNumber)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
