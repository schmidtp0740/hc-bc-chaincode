package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type hack struct {
	ObjectType string `json:"ObjType"`
	IsHacked   string `json:"isHacked"`
}

// isHacked
// input: nothing
// output: a boolean of whether or not the system is hacked
// Summary: get the status of whether or not the system is hacked or not
func (t *Chaincode) isHacked(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	hackRecordAsBytes, err := stub.GetState("hack")
	if err != nil {
		return shim.Error(err.Error())
	}
	if hackRecordAsBytes == nil {
		hackRecord := hack{
			ObjectType: "hack",
			IsHacked:   "false",
		}

		hackRecordAsBytes, err = json.Marshal(hackRecord)
		if err != nil {
			return shim.Error(err.Error())
		}

		if err := stub.PutState(hackRecord.ObjectType, hackRecordAsBytes); err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(hackRecordAsBytes)
}

// hack
// input: nothing
// output: success that the value has been swapped between true and false for isHacked
// Summary: set the value of isHacked to the opposite to transition betweent the system to being
// hacked or not
func (t *Chaincode) hack(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	hackRecord := hack{}
	hackRecordAsBytes, err := stub.GetState("hack")
	if err != nil {
		return shim.Error("Error retrieving state" + err.Error())
	}
	// there was no hack record to begin with so make a default hack record of IsHacked = false
	if hackRecordAsBytes == nil {
		hackRecord.ObjectType = "hack"
		hackRecord.IsHacked = "false"
	} else {
		// record existed so grabbing current state
		if err := json.Unmarshal(hackRecordAsBytes, &hackRecord); err != nil {
			return shim.Error(err.Error())
		}
	}

	if hackRecord.IsHacked == "true" {
		hackRecord.IsHacked = "false"
	} else {
		hackRecord.IsHacked = "true"
	}

	hackRecordAsBytes, err = json.Marshal(hackRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	if err := stub.PutState(hackRecord.ObjectType, hackRecordAsBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(hackRecordAsBytes)
}
