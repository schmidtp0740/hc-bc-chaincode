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
type heartRateMessage struct {
	HeartRate int `json:"heartRate,omitempty"` // heart rate of the patient
	Timestamp int `json:"timestamp,omitempty"` // timestamp of the record
}

// newHeartRateMessage
// input: id, heart rate and timestamp
// output: confirmation of record saved
// summary: insert each entry of a new heart rate message
func (t *Chaincode) newHeartRateMessage(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Initiating newHeartRateMessage")
	//		0			1			2
	// "patientID", heartRate, timestamp
	if len(args) < 3 {
		return shim.Error("Incorrect number of arguements. Expected 3")
	}

	// Input Sanitation
	fmt.Println("- start newHeartRateMessage")
	if len(args[0]) <= 0 {
		return shim.Error("1st arguement myst be a non-empty string")
	}

	// convert patiendID to lowercase
	patientID := strings.ToLower(args[0])

	// convert heartrate from string to integer
	heartRate, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("2nd arguement must be a numeric string")
	}

	// convert timestamp from string to integer
	timestamp, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("3rd arguement must be a numeric string")
	}

	// create struct from args variables
	// Record ID is a composotion of both patient id and the timestamp
	newHeartRateMessage := heartRateMessage{
		HeartRate: heartRate,
		Timestamp: timestamp,
	}
	fmt.Printf("Converted args to heartRateMessage struct: %v\n", newHeartRateMessage)

	initialEMR := EMR{
		PatientID: patientID,
	}

	// check if the patient record exists
	initialEMRAsBytes, err := stub.GetState(initialEMR.PatientID)
	if err != nil {
		return shim.Error("Patient record does not exist")
	}

	//convert inital EMR bytes to struct
	err = json.Unmarshal(initialEMRAsBytes, &initialEMR)
	if err != nil {
		return shim.Error("unable to marshal emr bytes")
	}

	// add new heart rate message to initial EMR struct
	initialEMR.HeartRate = newHeartRateMessage

	// convert modified EMR to JSON bytes
	modifiedEMRAsBytes, err := json.Marshal(initialEMR)
	if err != nil {
		return shim.Error("Error attempting to marshaling iot data")
	}
	fmt.Printf("iot data as json %s", string(modifiedEMRAsBytes))

	// submit modified EMR data to ledger
	err = stub.PutState(initialEMR.PatientID, modifiedEMRAsBytes)
	if err != nil {
		return shim.Error("Error inserting iot data: " + err.Error())
	}

	// return success if reached to this point
	fmt.Println("- end of newHeartRateMessage")
	return shim.Success(nil)
}

// getHeartRateHistory
// input: recordID
// output: array of iot records
// summary: get history of heart rate data for patient
func (t *Chaincode) getHeartRateHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("-------- init getHeartRateHistory-----------")
	// check for args of patientID
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}

	// convert patientID to lowercase
	patientID := strings.ToLower(args[0])

	// get patient record
	// output will be list of patient records ( history )
	resultsIterator, err := stub.GetHistoryForKey(patientID)
	if err != nil {
		return shim.Error("Patient record does not exist: " + err.Error())
	}
	defer resultsIterator.Close()

	// create patient record heart rate history struct
	patientHeartRateHistory := struct {
		PatientID        string             `json:"patientID"`
		HeartRateHistory []heartRateMessage `json:"heartRateHistory,omitempty"`
	}{
		PatientID: patientID,
	}

	for resultsIterator.HasNext() {
		// create patient record interface
		patientRecord := EMR{}

		// get response
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// unmarshall response's value to patient interface
		if err := json.Unmarshal(response.Value, &patientRecord); err != nil {
			return shim.Error(err.Error())
		}

		fmt.Printf("tempPatientRecord: %v\n", patientRecord)

		// add new heart rate message if timestamp is not equal to last entry in heart rate history
		if len(patientHeartRateHistory.HeartRateHistory) == 0 && patientRecord.HeartRate.Timestamp != 0 {
			patientHeartRateHistory.HeartRateHistory = append(patientHeartRateHistory.HeartRateHistory, patientRecord.HeartRate)
		} else if len(patientHeartRateHistory.HeartRateHistory) > 0 &&
			patientHeartRateHistory.HeartRateHistory[len(patientHeartRateHistory.HeartRateHistory)-1].Timestamp != patientRecord.HeartRate.Timestamp {
			patientHeartRateHistory.HeartRateHistory = append(patientHeartRateHistory.HeartRateHistory, patientRecord.HeartRate)
		}

	}

	// convert the heart rate history struct to json bytes to be returned
	heartRateHistoryAsBytes, err := json.Marshal(patientHeartRateHistory)
	if err != nil {
		return shim.Error("error marshalling iot data in iot history")
	}
	fmt.Printf("- getHeartRateHistory returning:\n%s\n", string(heartRateHistoryAsBytes))

	return shim.Success(heartRateHistoryAsBytes)
}
