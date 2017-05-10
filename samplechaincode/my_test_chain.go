/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var empId = "ID"

type Employee struct {

ID int `json:"id"`
Name string `json:"name"`
Company string `json:"company"`

}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error
	
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return nil, err
	}
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)
	err = stub.PutState(empId, jsonAsBytes)
	
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}
	if function == "createEmployee" {
		// creates an entity from its state
		return t.createEmployee(stub, args)
	}
if function == "updateEmployee" {
		// updates an entity from its state
		return t.updateEmployee(stub, args)
	}	
	return nil, nil
}

func (t *SimpleChaincode) createEmployee(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {	
	var id int
	var name,company string
	var err error	
	
	 id,err = strconv.Atoi(args[0])
	 if err != nil {
		return nil, err
	}
	 name = args[1]
	 company = args[2]
	 fmt.Println("Employees data",id,name,company)
	 
	
	//Adding ID to array
	//start
	
	IDAsBytes, err := stub.GetState(empId)
	if err != nil {
		return nil, errors.New("Failed to get Employee IDs")
	}
	
	var empIds []string
	json.Unmarshal(IDAsBytes, &empIds)	
	fmt.Println("empIds  ",empIds)
	
	//store and append the index to empIds
	
	empIds = append(empIds, args[0])									
	fmt.Println("! math index: ", empIds)
	
	empIdAsBytes, _ := json.Marshal(empIds)
	err = stub.PutState(empId, empIdAsBytes)
	
	EmployeeStr := `{"ID": "` + args[0] + `","Name": "` + name + `", "Company": "` + company + `"}`
	
	err = stub.PutState(args[0], []byte(EmployeeStr))
	if err != nil {
		return nil, err
	}
	 
	return empIdAsBytes, nil
}
func (t *SimpleChaincode) updateEmployee(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {	
	var id int
	var name,company string
	var err error	
	
	 id,err = strconv.Atoi(args[0])
	 if err != nil {
		return nil, err
	}
	 name = args[1]
	 company = args[2]
	 fmt.Println("Employees data",id,name,company)

	updateStr := `{"ID": "` + args[0] + `","Name": "` + name + `", "Company": "` + company + `"}`
	
	err = stub.PutState(args[0], []byte(updateStr))									
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "queryAll" {
		var jsonRespAll string
		fmt.Println("Inside List Employees")
	var data []byte
	var IDs []string
	
	IDAsBytes, err := stub.GetState(empId)
	if err != nil {
		//return fail, errors.New("Failed to get math index")
	}	
	json.Unmarshal(IDAsBytes, &IDs)
	fmt.Println("IDAsBytes :",IDs)
	
	
	for i := range(IDs){
       // stu := Employees[i]	
	   empIdAsBytes, err := stub.GetState(IDs[i])						//grab this math
		if err != nil {
			//return fail, errors.New("Failed to get ")
		}
		fmt.Printf("IDs :",IDs[i])
			
		res := Employee{}
		json.Unmarshal(empIdAsBytes, &res)										//un stringify it aka JSON.parse()
		fmt.Printf("res data:",res)
		
		jsonResp := "{\"ID\":\"" + IDs[i] + "\",\"Name\":\"" + res.Name + "\",\"Company\":\"" + res.Company + "\"}"	
		jsonRespAll = jsonRespAll+jsonResp
	    fmt.Printf("Query Response:%s\n", jsonRespAll)
		
		data = []byte(jsonRespAll)		
       
    }
		return data, nil
	}
	if function == "query" {
		//return nil, errors.New("Invalid query function name. Expecting \"query\"")
	var Id,jsonResp string

	Id = args[0]
	IDAsbytes, err := stub.GetState(Id)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + Id + "\"}"
		return nil, errors.New(jsonResp)
	}

	return IDAsbytes, nil										
}
	
	return nil, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
