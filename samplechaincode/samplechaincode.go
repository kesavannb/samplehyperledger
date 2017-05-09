package main

import (
"errors"
	"fmt"
	//"strconv"
	
	"encoding/json"
		
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

var customerIndexStr = "_customerdetails"

type Customer struct{

Index string `json:"Index"`
ID string `json:"ID"`
Name string `json:"Name"`
Details string `json:"Details"`

}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
fmt.Println("init is running")
	var ID, Name, Details string    
		var err error
	// Initialize the chaincode

	ID = args[0]
	//ID, err = strconv.Atoi(args[0])
	//if err != nil {
		//return nil, errors.New("Expecting integer value for asset holding")
	//}
	
	Name = args[1]
	
		Details =args[2]
	
	
	fmt.Printf("IDvalue = %d, NameValue = %d, Detailsvalue = %d\n", ID, Name,Details)

	// Write the state to the ledger
	err = stub.PutState(ID, []byte(ID))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(Name, []byte(Name))
	if err != nil {
		return nil, err
	}

	err =stub.PutState(Details, []byte(Details))
	if err !=nil{
	return nil, err
	}

	var empty []string
	
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	
	err = stub.PutState(customerIndexStr, jsonAsBytes)
	
	if err != nil {
		return nil, err
	}
	
fmt.Println("deploying is result",customerIndexStr)

    return nil, nil
}
 
//Invoke Method
 
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {


  var err error
  
	if function == "save_data" {
	var ID, Name, Details string    

	ID = args[0]
	Name = args[1]
    Details =args[2]
   
   fmt.Printf("ID %d Name %d Details",ID,Name,Details)
   
	valueID := string(ID)
	valueName := string(Name)
	valueDetails := string(Details)
	Index := args[3]
			
	fmt.Printf("IDvalue = %d, NameValue = %d, Detailsvalue = %d\n", valueID,valueName,valueDetails)
		
	str := `{"Index": "` +Index+ `", "ID": "` + valueID+ `","Name": "` +valueName+ `","Details": "` +valueDetails+ `"}`			
	
	fmt.Println("str inside invoke",str)
	

	err = stub.PutState(Index, []byte(str))									//store marble with id as key
	fmt.Println("err",err)
	if err != nil {
		return nil, err
	}
	
	//get the math index

	customerBytes, err := stub.GetState(customerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get math index")
	}
	fmt.Println("customerBytes",customerBytes)
		
	var customerIndex []string
	json.Unmarshal(customerBytes, &customerIndex)	
	
	//store and append the index to Index
	fmt.Println("Index invoke",Index)
	customerIndex = append(customerIndex, Index)									//add math name to index list
	fmt.Println("! Customer index: ", customerIndex)
	
	jsonAsBytes, _ := json.Marshal(customerIndex)
	err = stub.PutState(customerIndexStr, jsonAsBytes)						//store name of marble

	return jsonAsBytes, nil
	
	
	}
	if function == "update" {											//writes a value to the chaincode state
		return t.update(stub, args)
	}
	if function == "delete" {											//writes a value to the chaincode state
		return t.delete(stub, args)
	}
	
	return nil, nil
	
}

//Update


func (t *SimpleChaincode)update(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
var Index, Name, Details string 
var err error

fmt.Println("Update is running--")

Index =args[0]
Name= args[1]
Details = args[2]
fmt.Println("Updating value --",Index,Name,Details)
updatestr :=`{"Index":"`+Index+`","Name":"`+Name+`","Details":"`+Details+`"}`
fmt.Println("Updated string",updatestr)
err = stub.PutState(Index,[]byte(updatestr))
fmt.Println("err--",err)
if (err != nil){
return nil, err
}
return nil, nil

}


// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Delete is running--")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	Index := args[0]
	fmt.Println("Delete is ID--",Index)
	// Delete the key from the state in ledger
	err := stub.DelState(Index)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}


//Query
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running")
	
		
		
		var Index, jsonResp string
	//var err error
	
	if function == "query" {
		//return nil, errors.New("Invalid query function name. Expecting \"query\"")
	

	Index = args[0]
	valAsbytes, err := stub.GetState(Index)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + Index + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
	
	
}

if function == "queryall" {

//========================================
	//for loop for incrementing the index
	//===========================
	
	//get the math index
	
	customerAsBytes, err := stub.GetState(customerIndexStr)
	if err != nil {
		//return fail, errors.New("Failed to get math index")
	}
	
	
	var cusIndex []string
	var data []byte
		var jsonRespAll string
	
	json.Unmarshal(customerAsBytes, &cusIndex)
	
	
	for i:= range cusIndex{													//iter through all the math		
		
		customerAsBytes, err := stub.GetState(cusIndex[i])						//grab this math
		if err != nil {
			//return fail, errors.New("Failed to get ")
		}
		fmt.Printf("cusIndex:",cusIndex[i])
			
		res := Customer{}
		json.Unmarshal(customerAsBytes, &res)										//un stringify it aka JSON.parse()
		fmt.Printf("res data:",res)
		
		jsonResp := "{\"Index\":\"" + res.Index + "\",\"ID\":\"" + res.ID + "\",\"Name\":\"" + res.Name + "\",\"Details\":\"" + res.Details + "\"}"
	
	    fmt.Printf("Query Response:%s\n", jsonResp)
		
		jsonRespAll = jsonRespAll+jsonResp
	    fmt.Printf("Query Response:%s\n", jsonResp)
		C:= string(jsonRespAll)
		data = []byte(C)
		
	
	}

return data, nil

}
return nil, nil

}
 
func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Println("Could not start SimpleChaincode")
    } else {
        fmt.Println("SimpleChaincode successfully started")
    }
 
}
