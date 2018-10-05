package main

type insurance struct {
	Name           string `json:"insuranceName,omitempty"`
	ExpirationDate string `json:"expirationDate,omitempty"`
	PolicyID       string `json:"policyID"`
}
