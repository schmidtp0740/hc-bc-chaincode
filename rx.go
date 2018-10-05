package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// rx
type rx struct {
	RXID         string `json:"rxid"`             // id of the prescription
	Timestamp    int    `json:"timestamp"`        // timestamp of when prescription was prescribed and filled
	Doctor       string `json:"doctor,omitempty"` // name of the doctor
	Pharmacist   string `json:"pharmacist,omitempty"`
	Prescription string `json:"prescription,omitempty"` // prescription name
	Refills      int    `json:"refills,emitempty"`      // number of refills
	ExpirateDate int    `json:"expDate,omitempty"`
	Status       string `json:"status,emitempty"` // current status of the prescription
}

// initPrescription: create a new prescription
func (t *Chaincode) insertRx(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0       		1      2     	3		   4	       		5		6
	// "patientID", "rxid", timestamp, "doctor", "prescription", refills, "status"
	if len(args) < 7 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init inserRx")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return shim.Error("9th argument must be a non-empty string")
	}

	patientID := args[0]
	rxid := args[1]

	timestamp, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("3rd arguement must be non empty integer string")
	}

	doctor := args[3]
	prescription := args[4]

	refills, err := strconv.Atoi(args[5])
	if err != nil {
		return shim.Error("5th arguement must be a non empty integer string")
	}

	status := args[6]

	// get patient Record
	patientRecord := EMR{}

	// retrieve patient record as bytes
	patientRecordAsBytes, err := stub.GetState(patientID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// convert patient record as bytes to struct
	if err := json.Unmarshal(patientRecordAsBytes, &patientRecord); err != nil {
		return shim.Error(err.Error())
	}

	// return error if the patient record does not exist
	if patientRecordAsBytes == nil {
		return shim.Error("Patient Record does not exist: " + err.Error())
	}

	newRx := rx{
		RXID:         rxid,
		Timestamp:    timestamp,
		Doctor:       doctor,
		Prescription: prescription,
		Refills:      refills,
		Status:       status,
	}

	// see if rxid already exists in patient record
	for _, tempRX := range patientRecord.RxList {
		if tempRX.RXID == newRx.RXID {
			return shim.Error("RXID already exists: " + tempRX.RXID)
		}
	}

	// add new prescription to patient record
	patientRecord.RxList = append(patientRecord.RxList, newRx)

	// convert record to JSON bytes
	patientRecordAsBytes, err = json.Marshal(patientRecord)
	if err != nil {
		return shim.Error("Error attempting to marshal rx: " + err.Error())
	}
	fmt.Printf("rx as json bytes: %s", string(patientRecordAsBytes))

	// put record to state ledger
	err = stub.PutState(patientRecord.PatientID, patientRecordAsBytes)
	if err != nil {
		return shim.Error("Error putting prescription to ledger: " + err.Error())
	}
	fmt.Printf("Entered state")

	fmt.Println("- end insertObject (success)")
	return shim.Success(nil)
}

// modifyPrescription: modifies existing prescription
func (t *Chaincode) modifyRx(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0       	1      	2     		3		   4			5	       		6		7
	// "patientid", "rxid", timestamp, "doctor", "pharmacist", "prescription", refills,"status"
	if len(args) < 8 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init modifyObject")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}

	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("7th argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return shim.Error("8th argument must be a non-empty string")
	}

	patientID := args[0]
	rxid := args[1]
	timestamp, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	doctor := args[3]
	pharmacist := args[4]
	prescription := args[5]

	refills, err := strconv.Atoi(args[6])
	if err != nil {
		return shim.Error(err.Error())
	}

	status := args[7]

	// retrieve patient record
	patientRecordAsBytes, err := stub.GetState(patientID)
	if err != nil {
		return shim.Error("Failed to get record: " + patientID)
	} else if patientRecordAsBytes == nil {
		return shim.Error("patient record does not exist: " + patientID)
	}

	// create a patient record interface to load bytes into
	patientRecord := EMR{}
	err = json.Unmarshal(patientRecordAsBytes, &patientRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	// check if prescription record exists
	IfExists := false
	for key, tempRx := range patientRecord.RxList {
		// update rx record with new details
		if tempRx.RXID == rxid {
			patientRecord.RxList[key].Doctor = doctor
			patientRecord.RxList[key].Pharmacist = pharmacist
			patientRecord.RxList[key].Prescription = prescription
			patientRecord.RxList[key].Refills = refills
			patientRecord.RxList[key].Status = status
			patientRecord.RxList[key].Timestamp = timestamp
			IfExists = true
		}
	}

	if IfExists == false {
		return shim.Error("RXID does not exist: " + rxid)
	}

	// convert struct to json bytes
	patientRecordAsBytes, err = json.Marshal(patientRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	// send rx record to state ledger
	err = stub.PutState(patientID, patientRecordAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end modifyObject (success)")
	return shim.Success(nil)
}

func (t *Chaincode) getRxHistoryOfPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 	0
	// "patientid"

	// check for args of RXID
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}

	// convert patientID to lowercase
	patientID := strings.ToLower(args[0])

	// retrieve iterator of the history for a patient record
	resultsIterator, err := stub.GetHistoryForKey(patientID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// create struct that returns a history of past rx transactions
	// RxHistory is a list that contains a point in time for every time the patients
	// records prescription list changed
	rxHistoryResponse := struct {
		PatientID string `json:"patientID"`
		RxHistory [][]rx `json:"rxHistory"`
	}{
		PatientID: patientID,
	}

	// check for results in the iterator
	for resultsIterator.HasNext() {
		// retrieve the results
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// create a temporary patient record to hold the iterators values
		tempPatientRecord := EMR{}

		// unmarshal bytes to temporary patient record
		if err := json.Unmarshal(response.Value, &tempPatientRecord); err != nil {
			return shim.Error(err.Error())
		}

		// check if the length of th last rxList is the same as the temporary
		// patients record list
		// this will check if a new prescription was added or deleted
		// in any case add the point in time state to the rx history response
		// else iterate through each prescription to see if it changed

		if len(rxHistoryResponse.RxHistory) >= 0 && len(tempPatientRecord.RxList) == 0 {
			continue
		} else if len(rxHistoryResponse.RxHistory) == 0 && len(tempPatientRecord.RxList) >= 1 {
			rxHistoryResponse.RxHistory = append(rxHistoryResponse.RxHistory, tempPatientRecord.RxList)
		} else if len(rxHistoryResponse.RxHistory[len(rxHistoryResponse.RxHistory)-1]) != len(tempPatientRecord.RxList) {
			rxHistoryResponse.RxHistory = append(rxHistoryResponse.RxHistory, tempPatientRecord.RxList)
		} else {
			for key, tempPatientRx := range tempPatientRecord.RxList {
				if rxHistoryResponse.RxHistory[len(rxHistoryResponse.RxHistory)-1][key] != tempPatientRx {
					rxHistoryResponse.RxHistory = append(rxHistoryResponse.RxHistory, tempPatientRecord.RxList)
				}
			}
		}
		// else {
		// 	for key, rx := range tempPatientRecord.RxList {
		// 		if rx != rxHistoryResponse.RxHistory[len(rxHistoryResponse.RxHistory)-1][key] {
		// 			rxHistoryResponse.RxHistory = append(rxHistoryResponse.RxHistory, tempPatientRecord.RxList)
		// 		}
		// 	}
		// }
	}

	rxHistoryResponseAsBytes, err := json.Marshal(rxHistoryResponse)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(rxHistoryResponseAsBytes)
}

func (t *Chaincode) getAllRx(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

func (t *Chaincode) getRxForPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//	0
	// "patientID"

	if len(args) < 1 {
		return shim.Error("Expectin 1 arguement: patientID")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st arguement must be a non empty string")
	}

	patientID := args[0]

	// create empty patient record interface
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

	// create custom struct for response of list of prescriptions for a given patient
	response := struct {
		PatientID string `json:"patientID"`
		RxList    []rx   `json:"rxList,omitempty"`
	}{
		PatientID: patientRecord.PatientID,
		RxList:    patientRecord.RxList,
	}

	// convert reponse to bytes
	responseAsBytes, err := json.Marshal(response)
	if err != nil {
		return shim.Error(err.Error())
	}

	// return results
	return shim.Success(responseAsBytes)
}
