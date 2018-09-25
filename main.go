package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode example simple Chaincode implementation
type Chaincode struct {
}

// Main
func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting File Trace chaincode: %s", err)
	}
}

// Init initializes chaincode
func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "getAllRx" {
		// TODO
		return t.getAllRx(stub, args)
	} else if function == "getRxForPatient" {
		// TODO
		return t.getRxForPatient(stub, args)
	} else if function == "insertRx" { //insert an object into the ledger
		return t.insertRx(stub, args)
	} else if function == "modifyRx" { //modify an attribute of an object
		return t.modifyRx(stub, args)
	} else if function == "getRxHistory" {
		return t.getRxHistory(stub, args) // get history of prescription
	} else if function == "newHeartRateRecord" {
		return t.newHeartRateRecord(stub, args) // insert IOT Data
	} else if function == "getIOTHistory" {
		return t.getIOTHistory(stub, args) // get history of IOT data
	} else if function == "getPerson" {
		// TODO
		return t.getPerson(stub, args)
	} else if function == "getPeople" {
		// TODO
		return t.getPeople(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}
