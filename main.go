package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode example simple Chaincode implementation
type Chaincode struct {
}

// EMR medical record containing PII, heart rate, blood pressure, etc
// summary:
// store patients record, heart rate data and insurance information
type EMR struct {
	ObjectType    string           `json:"objType"`
	PatientID     string           `json:"id"`                      // patient id must be in the format of "p###"
	FirstName     string           `json:"firstName"`               // will be lowercase
	LastName      string           `json:"lastName"`                // will be lowercas
	DOB           string           `json:"dob"`                     // format of MM/DD/YYYY
	Address       string           `json:"address"`                 // format is street address city, state, zip
	Phone         string           `json:"phone"`                   // format is ###-###-####
	HeartRate     heartRateMessage `json:"heartRate,omitempty"`     // current heart rate message
	RxList        []rx             `json:"rxList,omitempty"`        // list of prescriptions that the patient has currently
	Insurance     insurance        `json:"insurance,omitempty"`     // current insurance
	BloodPressure bloodPressure    `json:"bloodPressure,omitempty"` // current blood pressure
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

	response := t.initPerson(stub, []string{"p01", "john", "doe", "01/01/2000", "111 address city, state, zip", "111-111-1111"})
	fmt.Println(response.GetMessage())
	response = t.initPerson(stub, []string{"p02", "mary", "jane", "01/01/2000", "111 address city, state, zip", "111-111-1111"})
	fmt.Println(response.GetMessage())
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "insertRx" {
		// TESTED OK - check with approved attribute
		return t.insertRx(stub, args)
	} else if function == "getRxForPatient" {
		// TESTED OK - check with approved attribute
		return t.getRxForPatient(stub, args)
	} else if function == "getAllRx" {
		// TODO
		return t.getAllRx(stub, args)
	} else if function == "modifyRx" {
		// TESTED OK - check with approved attribute
		return t.modifyRx(stub, args)
	} else if function == "getRxHistoryOfPatient" {
		// bug found
		return t.getRxHistoryOfPatient(stub, args)
	} else if function == "newHeartRateMessage" {
		// TESTED OK
		return t.newHeartRateMessage(stub, args) // insert new heart rate message to blockchain
	} else if function == "getHeartRateHistory" {
		// TESTED OK
		return t.getHeartRateHistory(stub, args) // get history of heart rate data for a given patient
	} else if function == "getPerson" {
		// TESTED OK
		return t.getPerson(stub, args)
	} else if function == "getPeople" {
		// TESTED OK
		return t.getPeople(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// createIndex - create search index for ledger
// currently used to create a composite key
// used in getPeople by adding each person to a people index that we can query against
func (t *Chaincode) createIndex(stub shim.ChaincodeStubInterface, indexName string, attributes []string) error {
	fmt.Println("- start create index")
	var err error
	//  ==== Index the object to enable range queries, e.g. return all parts made by supplier b ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexKey, err := stub.CreateCompositeKey(indexName, attributes)
	if err != nil {
		return err
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of object.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(indexKey, value)

	fmt.Println("- end create index")
	return nil
}
