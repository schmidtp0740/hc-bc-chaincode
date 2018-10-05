#!/bin/bash
printf "Testing getPerson\n"

# Testing getPerson
curl -H "Content-type:application/json" -X POST http://localhost:4001/bcsgw/rest/v1/transaction/query -d '{"channel":"mychannel","chaincode":"emrcc","method":"getPerson","args":["p01"],"chaincodeVer":"v5"}'
printf "\n\n"

# # Testing getPerson
# curl -H "Content-type:application/json" -X POST http://localhost:4001/bcsgw/rest/v1/transaction/query -d '{"channel":"mychannel","chaincode":"emrcc","method":"getPerson","args":["p02"],"chaincodeVer":"v5"}'
# printf "\n\n"

# # Testing getPeople
# curl -H "Content-type:application/json" -X POST http://localhost:4001/bcsgw/rest/v1/transaction/query -d '{"channel":"mychannel","chaincode":"emrcc","method":"getPeople","args":[],"chaincodeVer": "v5" }'
# printf "\n\n"

# Testing insertRx
# printf "Testing insertRx\n"
# curl -H "Content-type:application/json" -X POST http://localhost:4001/bcsgw/rest/v1/transaction/invocation -d '{"channel":"mychannel","chaincode":"emrcc","method":"insertRx","args":["p01", "rx01", "1538511402080", "dr sloan", "atenolol", "0", "prescribed" ],"chaincodeVer": "v5" }'
# printf "\n\n"

# Confirming insertRx
# curl -H "Content-type:application/json" -X POST http://localhost:4001/bcsgw/rest/v1/transaction/query -d '{"channel":"mychannel","chaincode":"emrcc","method":"getRxForPatient","args":["p01"],"chaincodeVer":"v5"}'
# printf "\n\n"


# Testing modifyRx
printf "Testing modifyRx\n"
# curl -H "Content-type:application/json" -X POST http://localhost:4001/bcsgw/rest/v1/transaction/invocation -d '{"channel":"mychannel","chaincode":"emrcc","method":"modifyRx","args":["p01", "rx01", "1538511402080", "dr sloan", "blake","atenolol", "0", "filled" ],"chaincodeVer": "v5" }'
# printf "\n\n"

# Confirming modifyRx
curl -H "Content-type:application/json" -X POST http://localhost:4001/bcsgw/rest/v1/transaction/query -d '{"channel":"mychannel","chaincode":"emrcc","method":"getRxForPatient","args":["p01"],"chaincodeVer":"v5"}'
printf "\n\n"

# Testing getRxHistoryOfPatient
printf "Testing getRxHistoryOfPatient\n"
curl -H "Content-type:application/json" -X POST http://localhost:4001/bcsgw/rest/v1/transaction/query -d '{"channel":"mychannel","chaincode":"emrcc","method":"getRxHistoryOfPatient","args":["p01"],"chaincodeVer":"v5"}'
printf "\n\n"


