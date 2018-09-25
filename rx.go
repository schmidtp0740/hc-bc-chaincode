package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// rx
type rx struct {
	ObjectType   string `json:"objType"`      //objType is used to distinguish the various types of objects in state database
	RXID         string `json:"rxid"`         // id of the prescription
	ID           string `json:"id"`           // id of the patient
	FirstName    string `json:"firstName"`    // first name of the patient
	LastName     string `json:"lastName"`     // last name of the patient
	Timestamp    int    `json:"timestamp"`    // timestamp of when prescription was prescribed and filled
	Doctor       string `json:"doctor"`       // name of the doctor
	Prescription string `json:"prescription"` // prescription name
	Refills      int    `json:"refills"`      // number of refills
	Status       string `json:"status"`       // current status of the prescription
}

// initPrescription: create a new prescription
func (t *Chaincode) insertRx(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0       1      2     		3		   4			5	       6			7			8
	// "rxid", "id", "firstName", "lastName", timestamp, "doctor", "prescription", refills, "status"
	if len(args) < 9 {
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
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return shim.Error("7th argument must be a non-empty string")
	}
	if len(args[8]) <= 0 {
		return shim.Error("9th argument must be a non-empty string")
	}

	// convert args to a rx struct
	tempRX, err := argsToRX(args)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Converted args to rx: %v\n", tempRX)

	// retrieve any state with received RXID
	emrAsBytes, err := stub.GetState(tempRX.RXID)
	if err != nil {
		return shim.Error("Failed to get record: " + err.Error())
	} else if emrAsBytes != nil {
		return shim.Error("This record already exists: " + tempRX.RXID)
	}

	// convert record to JSON bytes
	emrJSONAsBytes, err := json.Marshal(tempRX)
	if err != nil {
		return shim.Error("Error attempting to marshal rx: " + err.Error())
	}
	fmt.Printf("rx as json bytes: %s", string(emrAsBytes))

	// put record to state ledger
	err = stub.PutState(tempRX.RXID, emrJSONAsBytes)
	if err != nil {
		return shim.Error("Error putting prescription to ledger: " + err.Error())
	}
	fmt.Printf("Entered state")

	fmt.Println("- end insertObject (success)")
	return shim.Success(nil)
}

// modifyPrescription: modifies existing prescription
func (t *Chaincode) modifyRx(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0       		1      		2     	3		   4			5	       6			7			8
	// "rxid", "id", "firstName", "lastName", timestamp, "doctor", "prescription",refills,"status"
	if len(args) < 9 {
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
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return shim.Error("7th argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return shim.Error("8th argument must be a non-empty string")
	}
	if len(args[8]) <= 0 {
		return shim.Error("9th argument must be a non-empty string")
	}

	// convert args to rx record
	tempRX, err := argsToRX(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	// retrieve any state with retrieved rxid
	emrAsBytes, err := stub.GetState(tempRX.RXID)
	if err != nil {
		return shim.Error("Failed to get rx: " + tempRX.RXID)
	} else if emrAsBytes == nil {
		return shim.Error("prescription does not exist: " + tempRX.RXID)
	}

	// create an empty rx struct and unmarshal current state of rx
	emrToUpdate := rx{}
	err = json.Unmarshal(emrAsBytes, &emrToUpdate)
	if err != nil {
		return shim.Error(err.Error())
	}

	// update rx record with any new values
	emrToUpdate = tempRX

	// convert struct to json bytes
	emrJSONBytes, err := json.Marshal(emrToUpdate)
	if err != nil {
		return shim.Error(err.Error())
	}

	// send rx record to state ledger
	err = stub.PutState(tempRX.RXID, emrJSONBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end modifyObject (success)")
	return shim.Success(nil)
}

// argsToRX convert args to struct
// Input: args from shim
// Output: rx struct
func argsToRX(args []string) (rx, error) {
	var tempRX rx

	rxid := strings.ToLower(args[0])
	id := strings.ToLower(args[1])
	firstName := strings.ToLower(args[2])
	lastName := strings.ToLower(args[3])

	timestamp, err := strconv.Atoi(args[4])
	if err != nil {
		return tempRX, errors.New("5th argument must be a numeric string")
	}

	doctor := strings.ToLower(args[5])
	prescription := strings.ToLower(args[6])
	refills, err := strconv.Atoi(args[7])
	if err != nil {
		return tempRX, errors.New("8th argument must be a numeric string")
	}
	status := strings.ToLower(args[8])

	tempRX = rx{
		ObjectType:   "rx",
		RXID:         rxid,
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Timestamp:    timestamp,
		Doctor:       doctor,
		Prescription: prescription,
		Refills:      refills,
		Status:       status,
	}

	return tempRX, err
}

func (t *Chaincode) getRxHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// check for args of RXID
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}

	// convert rxid to lowercase
	rxid := strings.ToLower(args[0])

	// iterate over history for each record state
	resultsIterator, err := stub.GetHistoryForKey(rxid)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// create a buffer
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON vehiclePart)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForRecord returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (t *Chaincode) getAllRx(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

func (t *Chaincode) getRxForPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
