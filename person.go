package main

import (
	"encoding/json"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

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
	phone := strings.ToLower(args[5])

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
	newPersonRecord := EMR{
		ObjectType: "emr",
		PatientID:  patientID,
		FirstName:  firstName,
		LastName:   lastName,
		DOB:        dob,
		Address:    address,
		Phone:      phone,
	}

	// convert struct to json bytes
	newPersonRecordAsBytes, err := json.Marshal(newPersonRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Create Index key to query for all people
	// allows us to query against all people
	indexName := "people"
	err = t.createIndex(stub, indexName, []string{"people", patientID})
	if err != nil {
		shim.Error(err.Error())
	}

	// submit person record as bytes to ledger
	err = stub.PutState(patientID, newPersonRecordAsBytes)
	if err != nil {
		return shim.Error("Error putting state in to ledger: " + err.Error())
	}

	return shim.Success(nil)
}

// getPerson
// input: patientID
// output: patientID, firstName, lastName, dob, address, phone
func (t *Chaincode) getPerson(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//	0
	//	personID
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguements. Expected 1")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non empty string")
	}

	id := args[0]

	patientRecordAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get record: " + err.Error())
	}

	initialEMR := EMR{}

	if err := json.Unmarshal(patientRecordAsBytes, &initialEMR); err != nil {
		return shim.Error(err.Error())
	}

	newPatientRecord := struct {
		PatientID string `json:"patientID"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		DOB       string `json:"dob"`
		Address   string `json:"address"`
		Phone     string `json:"phone"`
	}{
		PatientID: initialEMR.PatientID,
		FirstName: initialEMR.FirstName,
		LastName:  initialEMR.LastName,
		DOB:       initialEMR.DOB,
		Address:   initialEMR.Address,
		Phone:     initialEMR.Phone,
	}

	// Marshal patient record to bytes
	newPatientRecordAsBytes, err := json.Marshal(newPatientRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(newPatientRecordAsBytes)
}

// getPeople
// input: none
// output: all people in the database
func (t *Chaincode) getPeople(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// no input

	indexName := "people"

	personIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{"people"})
	if err != nil {
		return shim.Error("error getting people query result: " + err.Error())
	}
	defer personIterator.Close()

	type shortPersonRecord struct {
		PatientID string `json:"patientID,omitempty"`
		FirstName string `json:"firstName,omitempty"`
		LastName  string `json:"lastName,omitempty"`
	}

	responseStruct := struct {
		People []shortPersonRecord `json:"people,omitempty"`
	}{}

	for personIterator.HasNext() {
		response, err := personIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		_, components, err := stub.SplitCompositeKey(response.Key)

		patientID := components[1]

		// get person record
		tempPersonRecordAsBytes, err := stub.GetState(patientID)

		tempPersonRecord := EMR{}

		if err := json.Unmarshal(tempPersonRecordAsBytes, &tempPersonRecord); err != nil {
			return shim.Error(err.Error())
		}

		responseStruct.People = append(responseStruct.People,
			shortPersonRecord{
				PatientID: tempPersonRecord.PatientID,
				FirstName: tempPersonRecord.FirstName,
				LastName:  tempPersonRecord.LastName,
			})

	}

	responseJSONAsBytes, err := json.Marshal(responseStruct)
	if err != nil {
		return shim.Error("error marshalling people json bytes: " + err.Error())
	}

	return shim.Success(responseJSONAsBytes)
}
