/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// SmartContract implements a simple chaincode to manage an asset
type SmartContract struct {
}

// DataIPFS Struct
type DataIPFS struct {
	Hash string `json:"hash"`
	Key  string `json:"key"`
}

// MedicalInfo Struct
type MedicalInfo struct {
	PatientID string `json:"patient_id"`
	Gender    string `json:"gender"`
	Diagnosis string `json:"diagnosis"`
	Infotype  string `json:"infotype"`
}

// TxStruct Struct
type TxStruct struct {
	ObjectType    string   `json:"docType"`
	CreateReqTime string   `json:"create_request_time"`
	HospitalID    string   `json:"hospital_id"`
	TimeToSend    string   `json:"time_to_send"`
	MedicalInfo   string   `json:"medical_info"`
	MetaData      DataIPFS `json:"meta_data"`
	SecurityLv    int      `json:"securitylv, int"`
	Signature     string   `json:"signature"`
}

// HospitalInfo struct
type HospitalInfo struct {
	ObjectType string `json:"docType"`
	HospitalID string `json:"hospital_id"`
	HospitalIP string `json:"hospital_ip"`
}

// SendData struct
type SendData struct {
	Hash   string `json:"hash"`
	DestIP string `json:"dest_ip"`
	DestID string `json:"dest_id"`
}

// Init is called during chaincode instantiation to initialize any data.
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	Data := &HospitalInfo{"hospitalInfo", "000001", "http://10.0.2.15:9001"}
	DataJSON, _ := json.Marshal(Data)
	_ = stub.PutState("000001", DataJSON)
	Data = &HospitalInfo{"hospitalInfo", "000002", "http://10.0.2.15:9002"}
	DataJSON, _ = json.Marshal(Data)
	_ = stub.PutState("000002", DataJSON)
	return shim.Success(nil)
}

// Invoke 1
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	if function == "Init" {
		return s.Init(stub)
	} else if function == "GetData" {
		return s.getData(stub, args)
	} else if function == "AddData" {
		return s.addData(stub, args)
	} else if function == "QueryMedicalInfo" {
		return s.queryMedicalInfo(stub, args)
	} else if function == "ShareData" {
		return s.dataShare(stub, args)
	}

	fmt.Println("invoke did not find func " + function)
	return shim.Error("Received unkown function invocation")
}

func (s *SmartContract) getData(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		shim.Error("Incorrent number of arguments. Expection 1 (Name)")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	Name := args[0]

	Data, err := stub.GetState(Name)
	if err != nil {
		return shim.Error("Failed to get Data: " + err.Error())
	} else if Data == nil {
		fmt.Println("This Data does not exist: " + Name)
		return shim.Error("This Data does not exist: " + Name)
	}

	return shim.Success(Data)
}

func (s *SmartContract) addData(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 6 {
		shim.Error("Incorrent number of arguments. Expection 8 (Key, HospitalID, TimeToSend, MedicalInfo, MetaDataHash, MetaDataFileName, SecurityLevel, Signature)")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	} else if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	} else if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	} else if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	} else if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	} else if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	} else if len(args[6]) <= 0 {
		return shim.Error("7th argument must be a non-empty string")
	}
	Key := args[0]
	HospitalID := args[1]

	TimeToSend, _ := time.Parse(time.RFC3339, args[2])
	loc, _ := time.LoadLocation("Asia/Seoul")
	TimeToSend = TimeToSend.In(loc)

	MedicalInfo := args[3]
	MetaDataHash := args[4]
	SecurityLv, _ := strconv.Atoi(args[5])
	Signature := args[6]

	MetaData := &DataIPFS{MetaDataHash, Key}

	chkExist, err := stub.GetState(MetaDataHash)
	if err != nil {
		return shim.Error("Failed to get product: " + err.Error())
	} else if chkExist != nil {
		return shim.Error("This product already exists: " + MetaDataHash)
	}

	CreateRequestTime := time.Now().In(loc).Format(time.RFC3339)

	objectType := "MedicalData"
	Data := &TxStruct{objectType, CreateRequestTime, HospitalID, TimeToSend.Format(time.RFC3339), MedicalInfo, *MetaData, SecurityLv, Signature}
	DataJSON, err := json.Marshal(Data)

	if err != nil {
		return shim.Error("JSON Marshal Error: " + err.Error())
	}

	err = stub.PutState(MetaDataHash, DataJSON)

	if err != nil {
		return shim.Error("PutState Error: " + err.Error())
	}

	return shim.Success(DataJSON)
}

func (s *SmartContract) queryMedicalInfo(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//   0
	// "bob"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	MedicalInfo := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"MedicalData\",\"medical_info\":\"%s\"}}", MedicalInfo)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

func (s *SmartContract) dataShare(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	Hash := args[0]
	//Medical_Info := args[1]
	//Patient_ID := args[2]
	SecureLv, _ := strconv.Atoi(args[1])
	HospitalID := args[2]

	var Data TxStruct

	Tx, err := stub.GetState(Hash)

	err = json.Unmarshal(Tx, &Data)

	if err != nil {
		return shim.Error(err.Error())
	}

	if Data.SecurityLv > SecureLv {
		return shim.Error("Security Level is not Accept")
	}

	var SDest HospitalInfo

	Source, err := stub.GetState(HospitalID)

	if err != nil {
		return shim.Error("Failed to get Data: " + err.Error())
	} else if Source == nil {
		return shim.Error("This Data does not exist: ")
	}

	err = json.Unmarshal(Source, &SDest)

	if err != nil {
		return shim.Error("Debug 1" + err.Error())
	}

	var Dest HospitalInfo
	Hospitalinfo, err := stub.GetState(Data.HospitalID)

	err = json.Unmarshal(Hospitalinfo, &Dest)

	if err != nil {
		return shim.Error("Debug 2" + err.Error())
	}

	sendData := &SendData{Data.MetaData.Hash, SDest.HospitalIP, SDest.HospitalID}
	DataJSON, err := json.Marshal(sendData)

	if err != nil {
		return shim.Error("Debug 3" + err.Error())
	}

	buffer := bytes.NewBuffer(DataJSON)
	_, err = http.Post(Dest.HospitalIP, "application/json", buffer)

	return shim.Success(nil)
}

// func (s *SmartContract) updateProduct(stub shim.ChaincodeStubInterface, args []string) peer.Response {

// 	if len(args) != 2 {
// 		shim.Error("Incorrent number of arguments. Expection 2 (ProductName, Weight)")
// 	}

// 	if len(args[0]) <= 0 {
// 		return shim.Error("1st argument must be a non-empty string")
// 	} else if len(args[1]) <= 0 {
// 		return shim.Error("2st argument must be a non-empty string")
// 	}

// 	productName := args[0]
// 	weightpath := args[1]

// 	file, err := os.Open(weightpath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	stats, statsErr := file.Stat()
// 	if statsErr != nil {
// 		return nil, statsErr
// 	}

// 	var size int64 = stats.Size()
// 	f := make([]byte, size)

// 	bufr := bufio.NewReader(file)
// 	_, err = bufr.Read(f)

// 	chkExist, err := stub.GetState(productName)
// 	if err != nil {
// 		return shim.Error("Failed to get product: " + err.Error())
// 	} else if chkExist == nil {
// 		fmt.Println("This product does not exist: " + productName)
// 		return shim.Error("This product does not exist: " + productName)
// 	}

// 	productToUpdate := Product{}
// 	err = json.Unmarshal(chkExist, &productToUpdate)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	productToUpdate.Weight = f

// 	productJSON, _ := json.Marshal(productToUpdate)

// 	err = stub.PutState(productName, productJSON)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	return shim.Success(nil)
// }

// func (s *SmartContract) deleteProduct(stub shim.ChaincodeStubInterface, args []string) peer.Response {

// 	if len(args) != 1 {
// 		shim.Error("Incorrent number of arguments. Expection 1 (ProductName)")
// 	}

// 	if len(args[0]) <= 0 {
// 		return shim.Error("1st argument must be a non-empty string")
// 	}

// 	productName := args[0]

// 	product, err := stub.GetState(productName)
// 	if err != nil {
// 		return shim.Error("Failed to get product: " + err.Error())
// 	} else if product == nil {
// 		fmt.Println("This product does not exist: " + productName)
// 		return shim.Error("This product does not exist: " + productName)
// 	}

// 	productToDelete := Product{}

// 	err = json.Unmarshal([]byte(product), &productToDelete)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	err = stub.DelState(productName)
// 	if err != nil {
// 		return shim.Error("Falied to delete state: " + err.Error())
// 	}

// 	return shim.Success(nil)
// }
