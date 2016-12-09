/*/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	//"time"
	//"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ManagePO example simple Chaincode implementation
type ManagePO struct {
}

var POIndexStr = "_POindex"				//name for the key/value that will store a list of all known marbles
//var openTradesStr = "_opentrades"				//name for the key/value that will store all open trades

type PO struct{
	TransID string `json:"transId"`					//the fieldtags are needed to keep case from bouncing around
	SellerName string `json:"sellerName"`
	BuyerName string `json:"buyerName"`					//the fieldtags are needed to keep case from bouncing around
	ExpectedDeliveryDate string `json:"expectedDeliveryDate"`
	PO_status string `json:"po_status"`
	PO_date string `json:"po_date"`
	ItemID string `json:"item_id"`
	Item_name string `json:"item_name"`
	Item_quantity int `json:"item_quantity"`
}

/*type allPO struct{
	all_po []PO `json:"all_PO"`
}*/

/*type Item struct{
	ID string `json:"id"`					//the fieldtags are needed to keep case from bouncing around
	Name string `json:"name"`
	Quantity int `json:"quantity"`
}*/

/*type AnOpenTrade struct{
	User string `json:"user"`					//user who created the open trade order
	Timestamp int64 `json:"timestamp"`			//utc timestamp of creation
	Want Description  `json:"want"`				//description of desired marble
	Willing []Description `json:"willing"`		//array of marbles willing to trade away
}*/

type allPO struct{
	PO_list []PO `json:"po_list"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(ManagePO))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManagePO) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var msg string
	//var PO_list allPO

	var err error
	if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	// Initialize the chaincode
	
	msg = args[0]
	fmt.Println("ManagePO chaincode is deployed with the message : "+ msg);
	/*Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}*/

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(msg))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(POIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	/*var trades AllTrades
	jsonAsBytes, _ = json.Marshal(trades)								//clear the open trade struct
	err = stub.PutState(openTradesStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}*/
	
	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *ManagePO) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *ManagePO) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "create_po" {											//writes a value to the chaincode state
		return t.create_po(stub, args)
	}else if function == "delete_po" {									//create a new marble
		return t.delete_po(stub, args)
	}else if function == "update_po" {									//create a new trade order
		return t.update_po(stub, args)
	}/* else if function == "perform_trade" {									//forfill an open trade order
		res, err := t.perform_trade(stub, args)
		cleanTrades(stub)													//lets clean just in case
		return res, err
	} else if function == "remove_trade" {									//cancel an open trade order
		return t.remove_trade(stub, args)
	}*/
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *ManagePO) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getPO_byID" {													//read a variable
		return t.getPO_byID(stub, args)
	} else if function == "getPO_byBuyer" {													//read a variable
		return t.getPO_byBuyer(stub, args)
	} else if function == "getPO_bySeller" {													//read a variable
		return t.getPO_bySeller(stub, args)
	} else if function == "get_AllPO" {													//read a variable
		return t.get_AllPO(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// getPO_byID - display PO details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManagePO) getPO_byID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var transId, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ID of the var to query")
	}

	transId = args[0]
	valAsbytes, err := stub.GetState(transId)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + transId + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
//  getPO_byBuyer - get PO details by buyer name from chaincode state
// ============================================================================================================================
func (t *ManagePO) getPO_byBuyer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var errResp, buyerName string
	var poIndex, valIndex []string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	buyerName = args[0]
	poAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get PO index")
	}
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all PO")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		json.Unmarshal(valueAsBytes, &valIndex)
		for j,value := range valIndex{
			fmt.Println(strconv.Itoa(j) + " - looking at " + value + " for all PO")
			if value == buyerName {
				valAsBytes, err := stub.GetState(value)
				if err != nil {
					errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
					return nil, errors.New(errResp)
				}
				return valAsBytes,nil
			}
		}
	}
	return nil, nil													//send it onward
}
// ============================================================================================================================
//  getPO_bySeller - display PO details for a specific Seller from chaincode state
// ============================================================================================================================
func (t *ManagePO) getPO_bySeller(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var errResp, sellerName string
	var poIndex, valIndex []string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	sellerName = args[0]
	poAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get PO index")
	}
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	//jsonResp = "{"
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all PO")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		json.Unmarshal(valueAsBytes, &valIndex)
		for j,value := range valIndex{
			fmt.Println(strconv.Itoa(j) + " - looking at " + value + " for all PO")
			if value == sellerName {
				valAsBytes, err := stub.GetState(value)
				if err != nil {
					errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
					return nil, errors.New(errResp)
				}
				return valAsBytes,nil
			}
		}
	}
	return nil, nil													//send it onward
}
// ============================================================================================================================
//  get_AllPO- display details of all PO from chaincode state
// ============================================================================================================================
func (t *ManagePO) get_AllPO(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var poIndex []string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	poAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get PO index")
	}
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	jsonResp = "{"
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all PO")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i != 0 {
			jsonResp = jsonResp + ","
		}
		jsonResp = jsonResp + "}"
	}
	//jsonAsBytes, _ := json.Marshal(valueAsBytes)
	return []byte(jsonResp), nil
											//send it onward
}
// ============================================================================================================================
// Delete - remove a key/value pair from state
// ============================================================================================================================
func (t *ManagePO) delete_po(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	transId := args[0]
	err := stub.DelState(transId)													//remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the marble index
	poAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get marble index")
	}
	fmt.Println("poAsBytes in delete po")
	fmt.Println(poAsBytes);
	var poIndex []string
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	fmt.Println("poIndex in delete po")
	fmt.Println(poIndex);
	//remove marble from index
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + transId)
		if val == transId{															//find the correct marble
			fmt.Println("found marble")
			poIndex = append(poIndex[:i], poIndex[i+1:]...)			//remove it
			for x:= range poIndex{											//debug prints...
				fmt.Println(string(x) + " - " + poIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(poIndex)									//save new index
	err = stub.PutState(POIndexStr, jsonAsBytes)
	return nil, nil
}

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *ManagePO) update_po(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error
	fmt.Println("running update_po()")

	if len(args) != 9 {
		return nil, errors.New("Incorrect number of arguments. Expecting 9.")
	}
	transId := args[0]
	poAsBytes, err := stub.GetState(transId)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + transId + "\"}"
		return nil, errors.New(jsonResp)
	}
	
	res := PO{}
	json.Unmarshal(poAsBytes, &res)
	if res.TransID == transId{
		fmt.Println("PO found with id : " + transId)
		fmt.Println(res);
		res.SellerName = args[1]
		res.BuyerName = args[2]
		res.ExpectedDeliveryDate = args[3]
		res.PO_status = args[4]
		res.PO_date = args[5]
		res.ItemID = args[6]
		res.Item_name = args[7]
		res.Item_quantity, err = strconv.Atoi(args[8])
		if err != nil {
			return nil, errors.New("Expecting integer value for asset holding")
		}
	}
	
	//build the PO json string manually
	order := 	`{`+
			`"transId": "` + res.TransID + `" , `+
			`"sellerName": "` + res.SellerName + `" , `+
			`"buyerName": "` + res.BuyerName + `" , `+
			`"ExpectedDeliveryDate": "` + res.ExpectedDeliveryDate + `" , `+ 
			`"PO_status": "` + res.PO_status + `" , `+ 
			`"PO_date": "` + res.PO_date + `" , `+ 
			`"id": "` + res.ItemID + `" , `+ 
			`"name": "` + res.Item_name + `" , `+ 
			`"quantity": "` +  strconv.Itoa(res.Item_quantity) + `" `+ 
			`}`
	err = stub.PutState(transId, []byte(order))									//store PO with id as key
	if err != nil {
		return nil, err
	}
		
	//get the PO index
	/*poIndexAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get PO index")
	}
	var poIndex []string
	json.Unmarshal(poIndexAsBytes, &poIndex)							//un stringify it aka JSON.parse()
	
	//append
	poIndex = append(poIndex, transId)									//add PO transID to index list
	fmt.Println("! PO index: ", poIndex)
	jsonAsBytes, _ := json.Marshal(poIndex)
	err = stub.PutState(POIndexStr, jsonAsBytes)						//store name of PO

	fmt.Println("- end create PO")*/
	return nil, nil
}

// ============================================================================================================================
// Init PO - create a new PO, store into chaincode state
// ============================================================================================================================
func (t *ManagePO) create_po(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	//var po PO
	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 9 {
		return nil, errors.New("Incorrect number of arguments. Expecting 9")
	}
	//input sanitation
	fmt.Println("- start create_po")
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
		return nil, errors.New("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return nil, errors.New("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return nil, errors.New("7th argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return nil, errors.New("8th argument must be a non-empty string")
	}
	if len(args[8]) <= 0 {
		return nil, errors.New("9th argument must be a non-empty string")
	}
	transId := args[0]
	sellerName := args[1]
	buyerName := args[2]
	expectedDeliveryDate := args[3]
	po_status := args[4]
	po_date := args[5]
	item_id := args[6]
	item_name := args[7]
	item_quantity, err := strconv.Atoi(args[8])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	
	poAsBytes, err := stub.GetState(transId)
	if err != nil {
		return nil, errors.New("Failed to get PO transID")
	}
	res := PO{}
	json.Unmarshal(poAsBytes, &res)
	if res.TransID == transId{
		fmt.Println("This PO arleady exists: " + transId)
		fmt.Println(res);
		return nil, errors.New("This PO arleady exists")				//all stop a marble by this name exists
	}
	
	//build the PO json string manually
	order := 	`{`+
			`"transId": "` + transId + `" , `+
			`"sellerName": "` + sellerName + `" , `+
			`"buyerName": "` + buyerName + `" , `+
			`"ExpectedDeliveryDate": "` + expectedDeliveryDate + `" , `+ 
			`"PO_status": "` + po_status + `" , `+ 
			`"PO_date": "` + po_date + `" , `+ 
			`"id": "` + item_id + `" , `+ 
			`"name": "` + item_name + `" , `+ 
			`"quantity": "` +  strconv.Itoa(item_quantity) + `" `+ 
			`}`
	err = stub.PutState(transId, []byte(order))									//store PO with id as key
	if err != nil {
		return nil, err
	}
		
	//get the PO index
	poIndexAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get PO index")
	}
	var poIndex []string
	json.Unmarshal(poIndexAsBytes, &poIndex)							//un stringify it aka JSON.parse()
	
	//append
	poIndex = append(poIndex, transId)									//add PO transID to index list
	fmt.Println("! PO index: ", poIndex)
	jsonAsBytes, _ := json.Marshal(poIndex)
	err = stub.PutState(POIndexStr, jsonAsBytes)						//store name of PO

	fmt.Println("- end create PO")
	return nil, nil
}
