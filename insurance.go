package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type insurance struct {
	Name           string `json:"insuranceName,omitempty"`
	ExpirationDate int    `json:"expDate,omitempty"`
	PolicyID       string `json:"policyID,omitempty"`
}

// TODO ASAP
func (t *Chaincode) insertInsurance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//	0			1		2					3
	// "patientID", "name", expirationDate, "policyID"

	if len(args) < 4 {
		return shim.Error("Incorrect number of arguements, Expecting 4")
	}

	fmt.Println("---start insertInsurance----")
	if len(args[0]) <= 0 {
		return shim.Error("1st arguement must be non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd arguement must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("3rd arguement must be a non-empty string")
	}

	patientID := args[0]
	insuranceName := args[1]

	expirationDate, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("unable to convert 3rd aguement to integer")
	}

	policyID := args[3]

	patientRecord := EMR{}

	patientRecordAsBytes, err := stub.GetState(patientID)
	if err != nil {
		return shim.Error(err.Error())
	}

	if err := json.Unmarshal(patientRecordAsBytes, &patientRecord); err != nil {
		return shim.Error("Patient record does not exist")
	}

	newInsurance := insurance{
		Name:           insuranceName,
		ExpirationDate: expirationDate,
		PolicyID:       policyID,
	}

	if patientRecord.Insurance.PolicyID == newInsurance.PolicyID && patientRecord.Insurance.ExpirationDate == newInsurance.ExpirationDate {
		return shim.Error("Insurance policy already exists: " + newInsurance.PolicyID)
	}

	patientRecord.Insurance = newInsurance

	patientRecordAsBytes, err = json.Marshal(patientRecord)
	if err != nil {
		return shim.Error("Error attempting to marshal insurance")
	}

	// put record to state ledger
	if err = stub.PutState(patientRecord.PatientID, patientRecordAsBytes); err != nil {
		return shim.Error("unable to put insurance to state")
	}

	return shim.Success(nil)
}

func (t *Chaincode) getInsurance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0
	// "patientID"

	if len(args) < 1 {
		return shim.Error("Expecting 1 arguement: patientID")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st arguement must be a non empty string")
	}

	patientID := args[0]

	patientRecord := EMR{}

	// get current state of the given patient record
	patientRecordAsBytes, err := stub.GetState(patientID)
	if err != nil {
		return shim.Error("Unable to get record: " + err.Error())
	}

	// convert patient record as bytes to struct
	if err := json.Unmarshal(patientRecordAsBytes, &patientRecord); err != nil {
		return shim.Error(err.Error())
	}

	// create custom struct for response of insurance for a patient
	response := struct {
		PatientID string    `json:"patientID"`
		Insurance insurance `json:"insurance,omitempty"`
	}{
		PatientID: patientRecord.PatientID,
		Insurance: patientRecord.Insurance,
	}

	// convert reponse to bytes
	responseAsBytes, err := json.Marshal(response)
	if err != nil {
		return shim.Error(err.Error())
	}

	// return results
	return shim.Success(responseAsBytes)
}

// TODO sprint 2
func (t *Chaincode) getInsuranceHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

// TODO sprint 2
func (t *Chaincode) newClaim(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)
}

// TODO sprint 2
func (t *Chaincode) getClaim(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

// TODO sprint 2
func (t *Chaincode) getClaimHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
