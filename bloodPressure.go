package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type bloodPressure struct {
	Low       int `json:"low,omitempty"`
	High      int `json:"high,omitempty"`
	Timestamp int `json:"timestamp,omitempty"`
}

// newBloodPressure
// input: low and high blood pressure
// output: confirmation that data is on ledger
func (t *Chaincode) newBloodPressure(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0			1	2		3
	// "patientID", low, high, timestamp
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguements. Expected 4")
	}

	// input sanitation
	fmt.Println("- start newBloodPressure")
	if len(args[0]) <= 0 {
		return shim.Error("1st arguement must be a non-empty string")
	}

	// convert patientID to lowercase
	patientID := strings.ToLower(args[0])

	// convert high blood pressure from string to int
	high, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("2nd arguement must be a integer string")
	}

	// convert low blood pressure from string to int
	low, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("3rd arguement must be a integer string")
	}

	timestamp, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("4th arguement must be an integer string")
	}

	// create blood pressure struct from args variables
	initialBP := bloodPressure{
		High:      high,
		Low:       low,
		Timestamp: timestamp,
	}

	// create initialEMR with patientID
	initialEMR := EMR{
		PatientID: patientID,
	}

	// get EMR record with patient ID
	initialEMRAsBytes, err := stub.GetState(initialEMR.PatientID)
	if err != nil {
		return shim.Error("unable to get state" + err.Error())
	}

	// unmarshal EMR bytes to initialEMR struct
	if err = json.Unmarshal(initialEMRAsBytes, &initialEMR); err != nil {
		return shim.Error("unable to unmarshal struct: " + err.Error())
	}

	// modify blood pressure attribute with new blood pressure
	initialEMR.BloodPressure = initialBP

	// convert EMR to bytes
	modifiedEMRAsBytes, err := json.Marshal(initialEMR)
	if err != nil {
		return shim.Error(err.Error())
	}

	// submit new EMR record
	if err = stub.PutState(initialEMR.PatientID, modifiedEMRAsBytes); err != nil {
		return shim.Error("error putting state" + err.Error())
	}
	// return success

	return shim.Success(nil)
}

// getBloodPressureHistory
// input: patientID
// output: history of blood pressure for patient
func (t *Chaincode) getBloodPressureHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//	1
	// "patientID"

	// check that there is the correct number of arguements
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguements, expecting 1")
	}

	// check that the arguement is non empty
	if len(args[0]) <= 0 {
		return shim.Error("1st arguement must be a non empty string")
	}

	// convert args to patientID
	patientID := strings.ToLower(args[0])

	// get patient history
	resultsIterator, err := stub.GetHistoryForKey(patientID)
	if err != nil {
		return shim.Error("Patient record does not exist: " + err.Error())
	}
	defer resultsIterator.Close()

	// create patient record with just patientID and blood pressure history
	patientBloodPressureHistory := struct {
		PatientID            string          `json:"patientID,omitempty"`
		BloodPressureHistory []bloodPressure `json:"bloodPressureHistory,omitempty"`
	}{
		PatientID: patientID,
	}

	// iterate through patient history
	for resultsIterator.HasNext() {
		// create empty EMR record to marshall
		patientRecord := EMR{}

		// get iterators result
		result, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// unmarshal result's value to patientRecord interface
		if err := json.Unmarshal(result.Value, &patientRecord); err != nil {
			return shim.Error("unable to unmarshal value" + err.Error())
		}

		// check if blood history is not the same as last
		if len(patientBloodPressureHistory.BloodPressureHistory) == 0 && patientRecord.BloodPressure.Timestamp != 0 {
			patientBloodPressureHistory.BloodPressureHistory = append(patientBloodPressureHistory.BloodPressureHistory, patientRecord.BloodPressure)
		} else if len(patientBloodPressureHistory.BloodPressureHistory) > 0 &&
			patientBloodPressureHistory.BloodPressureHistory[len(patientBloodPressureHistory.BloodPressureHistory)-1].Timestamp != patientRecord.BloodPressure.Timestamp {
			patientBloodPressureHistory.BloodPressureHistory = append(patientBloodPressureHistory.BloodPressureHistory, patientRecord.BloodPressure)
		}
	}

	bloodPressureHistoryAsBytes, err := json.Marshal(patientBloodPressureHistory)
	if err != nil {
		return shim.Error("error marshalling blood pressure history" + err.Error())
	}

	return shim.Success(bloodPressureHistoryAsBytes)
}
