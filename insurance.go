package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type insurance struct {
	Name           string `json:"insuranceName,omitempty"`
	ExpirationDate string `json:"expirationDate,omitempty"`
	PolicyID       string `json:"policyID,omitempty"`
}

// TODO ASAP
func (t *Chaincode) insertInsurance(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)
}

// TODO ASAP
func (t *Chaincode) getInsurance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
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
