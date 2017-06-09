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
	"errors"
	"fmt"
 "encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
	var ideaNameIndexStr = "_ideaNameIndex"
	 var coll_NameIndexStr = "coll_NameIndex"
// SimpleChaincode example simple Chaincode implementation
type IdeaChaincode struct {
	 
}
type idea struct {
	Owner string  `json:"owner"`
	IdeaName string  `json:"IdeaName"`
	Collaboraters  string  `json:"coll_name"`
	Description string    `json:"description"`
	 
}


// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {			
	err := shim.Start(new(IdeaChaincode))
	if err != nil {
		fmt.Printf("Error starting Form management chaincode: %s", err)
	}
}

// Init resets all the things
func (t *IdeaChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var EmptyIdeaName []string
	IdeajsonAsBytes, _ := json.Marshal(EmptyIdeaName)								//marshal an emtpy array of strings to clear the index
	err := stub.PutState(ideaNameIndexStr, IdeajsonAsBytes)
	if err != nil {
		return nil, err
	}
/*	var EmptyIdeaName []string
	IdeajsonAsBytes, _ := json.Marshal(Idea)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(ideaNameIndexStr, IdeajsonAsBytes)
	if err != nil {
		return nil, err
	}*/
	return nil, nil
}


func (t *IdeaChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "createNew_idea" {											//create a new Form
		return t.createNew_idea(stub, args)
	} else if function == "update_idea" {											//create a new Form
		return t.update_idea(stub, args)
	}else if function == "set_collaboraters" {											//create a new Form
		return t.set_collaboraters(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)	
	jsonResp := "Error : Received unknown function invocation: "+ function 				//error
	return nil, errors.New(jsonResp)
}
// Query - legacy function
/*func (t *IdeaChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return shim.Error("query running")
}*/
func (t *IdeaChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	return nil, errors.New("running")

	// Handle different functions
	if function == "view_idea" {													//Read a Form by FormID
		return t.view_idea(stub, args)
	} 

	fmt.Println("query did not find func: " + function)				//error
	jsonResp := "Error : Received unknown function query: "+ function 
	return nil, errors.New(jsonResp)
}
func (t *IdeaChaincode) createNew_idea(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var err error
	fmt.Println("starting creating idea")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. key of the variable and value to set")
	}

	
	owner := args[0]
	IdeaName := args[1]                                   //rename for funsies
	collaboraters := args[2]
	description := args[3]
		chaincodeData := 	`{`+
		`"owner": "` + owner + `" , `+
		`"IdeaName": "` + IdeaName + `" , `+ 
		`"coll_name": "` + collaboraters + `" , `+ 
		`"description": "` + description + `"  `+
		`}`
	err = stub.PutState(IdeaName, []byte(chaincodeData))         //write the variable into the ledger
	if err != nil {
		return nil, errors.New(err.Error())
	}
	//get the name index
	ideaNameIndexAsBytes, err := stub.GetState(ideaNameIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Tier-2 Form index")
	}
	var ideaNameIndex []string
	json.Unmarshal(ideaNameIndexAsBytes, &ideaNameIndex)
	//append
	ideaNameIndex = append(ideaNameIndex, IdeaName)
	fmt.Println("ideaName index : ", ideaNameIndex)	
	jsonAsBytes, _ := json.Marshal(ideaNameIndex)
	fmt.Print("jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(ideaNameIndexStr, jsonAsBytes)						//store name of Form
	if err != nil {
		return nil, err
	}	
	fmt.Println("- New idea Created")
	return nil, nil
}



func (t *IdeaChaincode) set_collaboraters(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("starting init_owner")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	user :=  args[0]
	IdeaName :=  args[1]
	collaboraters :=  args[2]

	//check if user already exists
	ideaAsBytes, err := stub.GetState(IdeaName)
	if err == nil {
		fmt.Println("This get_collaborater already exists - " )
		return nil, errors.New("This get_collaborater already exists - ")
	}
	var ideaIndex idea
	json.Unmarshal(ideaAsBytes, &ideaIndex)
	                //convert to array of bytes
if user == ideaIndex.Owner{
	setCollaborator := 	`{`+
		`"owner": "` + ideaIndex.Owner + `" , `+
		`"IdeaName": "` + ideaIndex.IdeaName + `" , `+ 
		`"coll_name": "` + collaboraters + `" , `+ 
		`"description": "` + ideaIndex.Description + `"  `+
		`}`

		//store user
	err = stub.PutState(IdeaName, []byte(setCollaborator))           //store owner by its Id
	if err != nil {
		fmt.Println("Could not store user")
		return nil, errors.New(err.Error())
	}
}
	

	fmt.Println("- Collaborater Created")
	return nil, nil
}

func (t *IdeaChaincode) update_idea(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("starting updating idea")
var err error
	fmt.Println("starting init_owner")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	user :=  args[0]
	owner := args[1]
	IdeaName :=  args[2]
	collaboraters :=  args[3]
	description := args[4]
	//check if user already exists
	ideaAsBytes, err := stub.GetState(IdeaName)
	if err == nil {
		fmt.Println("This get_collaborater already exists - " )
		return nil, errors.New("This get_collaborater already exists - ")
	}
	var ideaUpdateIndex idea
	json.Unmarshal(ideaAsBytes, &ideaUpdateIndex)
	                //convert to array of bytes
if user == ideaUpdateIndex.Owner{
	update_idea := 	`{`+
		`"owner": "` + owner + `" , `+
		`"IdeaName": "` + IdeaName + `" , `+ 
		`"coll_name": "` + collaboraters + `" , `+ 
		`"description": "` + description + `"  `+
		`}`

		//store user
	err = stub.PutState(IdeaName, []byte(update_idea))           //store owner by its Id
	if err != nil {
		fmt.Println("Could not store user")
		return nil, errors.New(err.Error())
	}
} else if user == ideaUpdateIndex.Collaboraters{
		update_idea_2 := `{`+
			`"owner": "` + ideaUpdateIndex.Owner + `" , `+
			`"IdeaName": "` + ideaUpdateIndex.IdeaName + `" , `+ 
			`"coll_name": "` + ideaUpdateIndex.Collaboraters + `" , `+ 
			`"description": "` + description + `"  `+
			`}`
	err = stub.PutState(IdeaName, []byte(update_idea_2))           //store owner by its Id
	if err != nil {
		fmt.Println("Could not store user")
		return nil, errors.New(err.Error())
	}
}
	

	fmt.Println("- Collaborater Created")
	return nil, nil
}



func (t *IdeaChaincode) view_idea(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
var IdeaName, jsonResp string
	var err error
	fmt.Println("Fetching Form by inea name")
	if len(args) != 3{
		return nil, errors.New("Incorrect number of arguments. Expecting Form ID to query")
	}
	
	IdeaName = args[0]
	valAsbytes, err := stub.GetState(IdeaName)									//get the FAA_formNumber from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + IdeaName + "\"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Print("valAsbytes : ")
	fmt.Println(valAsbytes)
	fmt.Println("Fetched")
	return valAsbytes, nil													//send it onward
}