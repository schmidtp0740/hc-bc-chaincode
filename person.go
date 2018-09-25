package main

import (
	"encoding/json"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type person struct {
	ObjectType string `json:"objType"`
	ID         string `json:"id"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	DOB        string `json:"dob,omitempty"`
	Address    string `json:"address,omitempty"`
	Phone      string `json:"phone,omitempty"`
}

// initPerson
// input: patientID, firstname, last name, date of birth, address, and phone
// output: success or failure
// enter in a new patient record
func (t *Chaincode) initPerson(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//		0			1			2			3		4		5
	//	"patientID", "firstName", "lastName", "dob", "address", "phone"

	// Input Sanitation
	for key, value := range args {
		if len(value) <= 0 {
			return shim.Error(string(key) + "argument must be a non empty string")
		}
	}

	// convert all arguments to lower case
	patientID := strings.ToLower(args[0])
	firstName := strings.ToLower(args[1])
	lastName := strings.ToLower(args[2])
	dob := strings.ToLower(args[3])
	address := strings.ToLower(args[4])
	phone := strings.ToLower(args[4])

	// see if person id already exists
	personRecordAsBytes, err := stub.GetState(patientID)
	if err != nil {
		// error if failed to get record
		return shim.Error("Failed to get person record: " + err.Error())
	} else if personRecordAsBytes != nil {
		// error if record exists
		return shim.Error("This person's record already exists: " + patientID)
	}

	// create a person struct with all arguments
	tempPersonRecord := person{
		ObjectType: "person",
		ID:         patientID,
		FirstName:  firstName,
		LastName:   lastName,
		DOB:        dob,
		Address:    address,
		Phone:      phone,
	}

	// convert struct to json bytes
	personRecordAsBytes, err = json.Marshal(tempPersonRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	// submit person record as bytes to ledger
	err = stub.PutState(patientID, personRecordAsBytes)
	if err != nil {
		return shim.Error("Error putting state in to ledger: " + err.Error())
	}

	return shim.Success(nil)
}

// getPerson
// input
// output
func (t *Chaincode) getPerson(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//	0
	//	id
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguements. Expected 1")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non empty string")
	}

	id := args[0]

	personAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get record: " + err.Error())
	}

	return shim.Success(personAsBytes)
}

func (t *Chaincode) getPeople(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)
}
