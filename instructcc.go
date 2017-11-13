package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"


	"strconv"
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/protos/peer"
)

type Instruct struct {
	Sender    string
	TradeNo   string
	MetaReq   string
	MetaResp  string
	Timestamp string
	Status    string
	Remark    string
	Id        string
}

type InstructIds_Type []string

const ISTKEY string = "IST_"
const IST_SIZE_KEY string = "IST_SIZE"
//const TODO_IST_SIZE_KEY string = "TODO_IST_SIZE"
//const DOING_IST_SIZE_KEY string = "DOING_IST_SIZE"
const IST_IDS_KEY string = "IST_IDS"

//var todoStatus string = "0"
//var doingStatus string = "1"

//mapping(bytes32 => Instruct) public instructMap;
//uint public instructCount;
//bytes32[] public instructIdArray;
//
//bytes32[] public todoIdArray;
//uint public todoInstructCount;
//string public todoStatus = "0";  //status for TODO item
//
//bytes32[] public doingIdArray;
//uint public doingInstructCount;
//string public doingStatus = "1";  //status for DOING item

type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return t.init(stub)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "AddInstruct" { //
		return t.AddInstruct(stub, args)
	} else if function == "RemoveInstruct" { //
		return t.RemoveInstruct(stub, args)
	} else if function == "UpdateInstruct" { //transfer all marbles of a certain color
		return t.UpdateInstruct(stub, args)
	} else if function == "GetInstruct" { //
		return t.GetInstruct(stub, args)
	} else if function == "GetInstructIds" { //
		return t.GetInstructIds(stub,args)
	} else if function == "GetInstructSize" { //
		return t.GetInstructSize(stub,args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleChaincode) init(stub shim.ChaincodeStubInterface) peer.Response {
	err := stub.PutState(IST_SIZE_KEY, []byte(strconv.Itoa(0)))
	if err != nil {
		return shim.Error(err.Error())
	}
	//err = stub.PutState(TODO_IST_SIZE_KEY, []byte(strconv.Itoa(0)))
	//if err != nil {
	//	return shim.Error(err.Error())
	//}
	//err = stub.PutState(DOING_IST_SIZE_KEY, []byte(strconv.Itoa(0)))
	//if err != nil {
	//	return shim.Error(err.Error())
	//}
	instructIds := InstructIds_Type{}
	instructIdsBytes, _ := json.Marshal(instructIds)
	err = stub.PutState(IST_IDS_KEY, instructIdsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) AddInstruct(stub shim.ChaincodeStubInterface, args []string) peer.Response  {
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}
	var instructId, senderId , tradeNo, metaReq, metaResp, timestamp, status, remark string
	instructId = args[0]
	senderId = args[1]
	tradeNo = args[2]
	metaReq = args[3]
	metaResp = args[4]
	timestamp = args[5]
	status = args[6]
	remark = args[7]

	exist, err := stub.GetState(ISTKEY + instructId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if exist != nil {
		return shim.Error("Add Instruct Error. Instruct alread exist!")
	}

	var instruct = Instruct{
		Sender:	senderId,
		TradeNo: tradeNo,
		MetaReq: metaReq,
		MetaResp: metaResp,
		Timestamp: timestamp,
		Status: status,
		Remark: remark,
		Id: instructId,
	}
	instructBytes, _ := json.Marshal(instruct)
	fmt.Println("put instruct: " + string(instructBytes))
	err = stub.PutState(ISTKEY + instructId, instructBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = t.addInstructId(stub, instructId)
	if err != nil {
		return shim.Error(err.Error())
	}

	//if (status == todoStatus) {
	//
	//} else if (status == doingStatus) {
	//
	//}

	return shim.Success([]byte(instructId))
}

func (t *SimpleChaincode) RemoveInstruct(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var instructId string
	instructId = args[0]
	err := stub.DelState(ISTKEY + instructId)
	if err != nil {
		return shim.Error("Failed to delete instruct:" + err.Error())
	}
	err = t.decInstructSize(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) UpdateInstruct(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}
	var instructId, senderId , tradeNo, metaReq, metaResp, timestamp, status, remark string
	instructId = args[0]
	senderId = args[1]
	tradeNo = args[2]
	metaReq = args[3]
	metaResp = args[4]
	timestamp = args[5]
	status = args[6]
	remark = args[7]

	instructBytes, err := stub.GetState(ISTKEY + instructId)
	if err != nil || instructBytes == nil{
		return shim.Error("Failed to get instruct for " + instructId)
	}
	instruct := Instruct{}
	err = json.Unmarshal(instructBytes, &instruct)
	if err != nil {
		return shim.Error(err.Error())
	}

	//if instruct.status != status {
	//	if (instruct.status == todoStatus) {
	//
	//	} else if (instruct.status == doingStatus) {
	//
	//	}
	//	if (status == todoStatus) {
	//
	//	} else if (status == doingStatus) {
	//
	//	}
	//
	//}

	instruct = Instruct{
		Sender:	senderId,
		TradeNo: tradeNo,
		MetaReq: metaReq,
		MetaResp: metaResp,
		Timestamp: timestamp,
		Status: status,
		Remark: remark,
		Id: instructId,
	}
	instructBytes, _ = json.Marshal(instruct)
	err = stub.PutState(ISTKEY + instructId, instructBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(instructId))
}

func (t *SimpleChaincode) GetInstruct(stub shim.ChaincodeStubInterface, args []string) peer.Response  {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var instructId string
	instructId = args[0]
	instructBytes, err := stub.GetState(ISTKEY + instructId)
	if err != nil {
		return shim.Error("Failed to get instruct for " + instructId)
	}
	fmt.Println("GetInstruct: " + string(instructBytes))
	return shim.Success(instructBytes)
}

func (t *SimpleChaincode) addInstructId(stub shim.ChaincodeStubInterface, instructId string) error {
	instructIdsBytes, err := stub.GetState(IST_IDS_KEY)
	if err != nil {
		return errors.New("Failed to get instruct ids")
	}
	instructIds := InstructIds_Type{}
	err = json.Unmarshal(instructIdsBytes, &instructIds)
	if err != nil {
		return errors.New("Failed to Unmarshal instruct ids")
	}
	instructIds = append(instructIds, instructId)
	err = t.incInstructSize(stub)
	if err != nil {
		return err
	}
	instructIdsBytes, _ = json.Marshal(instructIds)
	err = stub.PutState(IST_IDS_KEY, instructIdsBytes)
	if err != nil {
		return err
	}
	return nil
}

func (t *SimpleChaincode) GetInstructIds(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	instructIdsBytes, err := stub.GetState(IST_IDS_KEY)
	if err != nil {
		return shim.Error("Failed to get instruct ids")
	}
	return shim.Success(instructIdsBytes)
}

func (t *SimpleChaincode) GetInstructSize(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	istSizeByte, err := stub.GetState(IST_SIZE_KEY)
	if err != nil {
		return shim.Error(err.Error())
	}
	if istSizeByte == nil {
		return shim.Error("IST_SIZE_KEY not exist!")
	}
	return shim.Success(istSizeByte)
}

func (t *SimpleChaincode) incInstructSize(stub shim.ChaincodeStubInterface) error {
	istSizeByte, err := stub.GetState(IST_SIZE_KEY)
	if err != nil {
		return err
	}
	istSize, err := strconv.Atoi(string(istSizeByte))
	if err != nil {
		return err
	}
	istSize++
	if istSize < 1 {
		return errors.New("InstructSize Error")
	}
	err = stub.PutState(IST_SIZE_KEY, []byte(strconv.Itoa(istSize)))
	if err != nil {
		return err
	}
	return nil
}

func (t *SimpleChaincode) decInstructSize(stub shim.ChaincodeStubInterface) error {
	istSizeByte, err := stub.GetState(IST_SIZE_KEY)
	if err != nil {
		return err
	}
	istSize, err := strconv.Atoi(string(istSizeByte))
	if err != nil {
		return err
	}
	istSize--
	if istSize < 0 {
		return errors.New("InstructSize Error")
	}
	err = stub.PutState(IST_SIZE_KEY, []byte(strconv.Itoa(istSize)))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
