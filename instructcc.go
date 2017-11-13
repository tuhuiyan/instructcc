package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"

	"instructcc/contract"
)



func main() {
	err := shim.Start(new(contract.SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
