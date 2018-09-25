package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// iotData
// timestamp is id of transaction
type iotData struct {
	ObjectType string `json:"objType"`   // obj type should be iot
	PatientID  string `json:"patientID"` // id of the patient
	RecordID   string `json:"recordID"`
	HeartRate  int    `json:"heartRate"` // heart rate of the patient
	TimeStamp  int    `json:"timeStamp"` // timestamp of the record
}

// initIOTData
// input: id, heart rate and timestamp
// output: confirmation of record saved
// saved each entry of heart rate data as it comes in
func (t *Chaincode) newHeartRateRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//		0			1			2			3
	// "recordID", "patientID", heartRate, timestamp
	if len(args) < 4 {
		return shim.Error("Incorrect number of arguements. Expected 3")
	}

	// Input Sanitation
	fmt.Println("- start insertIOTDATA")
	if len(args[0]) <= 0 {
		return shim.Error("1st arguement myst be a non-empty string")
	}

	// convert recordID to lowercase
	recordID := strings.ToLower(args[0])

	// convert patiendID to lowercase
	patientID := strings.ToLower(args[1])

	// convert heartrate from string to integer
	heartRate, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("2nd arguement must be a numeric string")
	}

	// convert timestamp from string to integer
	timestamp, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("3rd arguement must be a numeric string")
	}

	// create struct from args variables
	// Record ID is a composotion of both patient id and the timestamp
	tempIotData := iotData{
		ObjectType: "iot",
		PatientID:  patientID,
		RecordID:   recordID,
		HeartRate:  heartRate,
		TimeStamp:  timestamp,
	}
	fmt.Printf("Converted args to iot struct: %v\n", tempIotData)

	// check if the patient record exists
	_, err = stub.GetState(tempIotData.PatientID)
	if err != nil {
		return shim.Error("Patient record does not exist")
	}

	// convert iot struct to JSON bytes
	iotJSONAsBytes, err := json.Marshal(tempIotData)
	if err != nil {
		return shim.Error("Error attempting to marshaling iot data")
	}
	fmt.Printf("iot data as json %s", string(iotJSONAsBytes))

	// submit data to ledger
	err = stub.PutState(tempIotData.RecordID, iotJSONAsBytes)
	if err != nil {
		return shim.Error("Error inserting iot data: " + err.Error())
	}

	// return success if reached to this point
	fmt.Println("- end of insertIOTData")
	return shim.Success(nil)
}

// getIOTHistory
// input: recordID
// output: array of iot records
func (t *Chaincode) getIOTHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// check for args of ID
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}

	// convert recordID to lowercase
	recordID := strings.ToLower(args[0])

	// iterate over history for each record state
	resultsIterator, err := stub.GetHistoryForKey(recordID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// create a slice of iot records
	iotSlice := []iotData{}

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON vehiclePart)
		if !response.IsDelete {
			iotRecord := iotData{}

			err := json.Unmarshal(response.Value, &iotRecord)
			if err != nil {
				return shim.Error("Error unmarshalling record in iot history")
			}

			iotSlice = append(iotSlice, iotRecord)
		}

	}

	iotSliceAsBytes, err := json.Marshal(iotSlice)
	if err != nil {
		return shim.Error("error marshalling iot data in iot history")
	}
	fmt.Printf("- getIOTHistory returning:\n%s\n", string(iotSliceAsBytes))

	return shim.Success(iotSliceAsBytes)
}
