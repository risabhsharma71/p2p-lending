/*
Copyright IBM Corp 2016 All Rights Reserved.

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

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var userIndexStr = "_userindex"
//var campaignIndexStr= "_campaignindex"
//var transactionIndexStr= "_transactionindex"

type User struct {
	Name  string `json:"name"` //the fieldtags of user are needed to store in the ledger
    PhoneNo int    `json:"phoneno"`
	Email string `json:"email"`
	User_Type string `json:"usertype"`
    Pin int `json:"pin"`
    Pan_No string `json:"panno"`
    Upi string `json:"upi"`
}

/*type AllUsers struct{
	Userlist []User `json:"userlist"`
}*/

type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	//_, args := stub.GetFunctionAndParameters()
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval))) //making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty) //marshal an emtpy array of strings to clear the index
	err = stub.PutState(userIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	

	return nil, nil
}

// Invoke is ur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)

	} else if function == "User_register" {
		return t.User_register(stub, args)

	} 

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "readuser" { //read a variable
		return t.readuser(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// read - query function to read key/value pair

func (t *SimpleChaincode) readuser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error
    //var campaign_title,jsonResp string
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}


func (t *SimpleChaincode) User_register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	// "lol", "1", "323323", "r@r.com"
	if len(args) != 7 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	
	name := args[0]
	
	phone, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Failed to get phone as cannot convert it to int")
	}
	email := args[2]
	usertype:=args[3]
	pin, err := strconv.Atoi(args[4])
	if err != nil {
		return nil, errors.New("Failed to get phone as cannot convert it to int")
	}
	panno:=args[5]
	upi:=args[6]

	

	//build the marble json string manually
	str := `{"name": "` + name + `", "phone": "` + strconv.Itoa(phone) + `", "email": ` + email + `, "usertype": "` + usertype + `,"pin":"`+ strconv.Itoa(pin) + `,"panno":"`+ panno+ `,"upi":"`+ upi +`"}`
	var alluser []string
	//json.Unmarshal(AllUsersAsBytes, &user)										//un stringify it aka JSON.parse()
	
alluser = append(alluser, str); 
fmt.Println("alluser",alluser);
	
 	stringByte := "\x00" + strings.Join(alluser, "\x20\x00")
	err = stub.PutState("getusers", []byte(stringByte)) //store User with name as key
	
	fmt.Println([]byte(stringByte))
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	

	fmt.Println("- end user_register")
	return nil, nil
}