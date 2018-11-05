#!/bin/bash
clear

printf "newBloodPressure\n"
curl -H "Content-type:application/json" -X POST http://localhost:4001 -d '{"channel": "testchannel", "chaincode": "emrcc", "chaincodeVer": "v1", "method": "newBloodPressure", "args": ["p01", "10", "12", "1541440675318"]}'
printf "\n"