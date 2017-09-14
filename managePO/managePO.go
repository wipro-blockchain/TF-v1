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

"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ManagePO example simple Chaincode implementation
type ManagePO struct {
}

var POIndexStr = "_POindex"				//name for the key/value that will store a list of all known PO

type PO struct{							// Attributes of a PO 
	TransID string `json:"transId"`					
	SellerName string `json:"sellerName"`
	BuyerName string `json:"buyerName"`					
	ExpectedDeliveryDate string `json:"expectedDeliveryDate"`
	PO_status string `json:"po_status"`
	PO_date string `json:"po_date"`
	ItemId string `json:"item_id"`
	Item_name string `json:"item_name"`
	Item_quantity string `json:"item_quantity"`
	Price string `json:"price"`
	Buyer_sign string `json:"buyer_sign"`
	Seller_sign string `json:"seller_sign"`
	Seller_Remarks string `json:"seller_remarks"`
}
// ============================================================================================================================
// Main - start the chaincode for PO management
// ============================================================================================================================
func main() {			
	err := shim.Start(new(ManagePO))
	if err != nil {
		fmt.Printf("Error starting PO management chaincode: %s", err)
	}
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManagePO) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var msg string
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting ' ' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// Initialize the chaincode
	msg = args[0]
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
	tosend := "{ \"message\" : \"ManagePO chaincode is deployed successfully.\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
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
	} else if function == "create_po" {											//create a new PO
		return t.create_po(stub, args)
	}else if function == "delete_po" {									// delete a PO
		return t.delete_po(stub, args)
	}else if function == "update_po" {									//update a PO
		return t.update_po(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)
	errMsg := "{ \"message\" : \"Received unknown function invocation\", \"code\" : \"503\"}"
	err := stub.SetEvent("errEvent", []byte(errMsg))
	if err != nil {
		return nil, err
	} 
	return nil, nil			//error
}
// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *ManagePO) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getPO_byID" {													//Read a PO by transId
		return t.getPO_byID(stub, args)
	} else if function == "getPO_byBuyer" {													//Read a PO by Buyer's name
		return t.getPO_byBuyer(stub, args)
	} else if function == "getPO_bySeller" {													//Read a PO by Seller's name
		return t.getPO_bySeller(stub, args)
	} else if function == "get_AllPO" {													//Read all POs
		return t.get_AllPO(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error
	errMsg := "{ \"message\" : \"Received unknown function query\", \"code\" : \"503\"}"
	err := stub.SetEvent("errEvent", []byte(errMsg))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// getPO_byID - get PO details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManagePO) getPO_byID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var transId string
	var err error
	fmt.Println("start getPO_byID")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'transId' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set transId
	transId = args[0]
	valAsbytes, err := stub.GetState(transId)									//get the transId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \""+ transId + " not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	//fmt.Print("valAsbytes : ")
	//fmt.Println(valAsbytes)
	fmt.Println("end getPO_byID")
	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
//  getPO_byBuyer - get PO details by buyer's name from chaincode state
// ============================================================================================================================
func (t *ManagePO) getPO_byBuyer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, buyerName, errResp string
	var poIndex []string
	var valIndex PO
	fmt.Println("start getPO_byBuyer")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'buyerName' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set buyer's name
	buyerName = args[0]
	//fmt.Println("buyerName" + buyerName)
	poAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get PO index string")
	}
	//fmt.Print("poAsBytes : ")
	//fmt.Println(poAsBytes)
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	fmt.Print("poIndex : ")
	fmt.Println(poIndex)
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	jsonResp = "{"
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getPO_byBuyer")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.BuyerName == buyerName{
			fmt.Println("Buyer found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			//fmt.Println("jsonResp inside if")
			//fmt.Println(jsonResp)
			if i < len(poIndex)-1 {
				jsonResp = jsonResp + ","
			}
		} 
	}
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end getPO_byBuyer")
	return []byte(jsonResp), nil											//send it onward
}

// ============================================================================================================================
//  getPO_bySeller - get PO details for a specific Seller from chaincode state
// ============================================================================================================================
func (t *ManagePO) getPO_bySeller(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, sellerName, errResp string
	var poIndex []string
	var valIndex PO
	fmt.Println("start getPO_bySeller")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'sellerName' as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set seller name
	sellerName = args[0]
	//fmt.Println("buyerName" + sellerName)
	poAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get PO index")
	}
	//fmt.Print("poAsBytes : ")
	//fmt.Println(poAsBytes)
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	fmt.Print("poIndex : ")
	fmt.Println(poIndex)
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	jsonResp = "{"
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getting sellerName")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.SellerName == sellerName{
			fmt.Println("Seller found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			//fmt.Println("jsonResp inside if")
			//fmt.Println(jsonResp)
			if i < len(poIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}
		
	}
	
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end getPO_bySeller")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
//  get_AllPO- get details of all PO from chaincode state
// ============================================================================================================================
func (t *ManagePO) get_AllPO(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var poIndex []string
	fmt.Println("start get_AllPO")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \" \" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	poAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get PO index")
	}
	//fmt.Print("poAsBytes : ")
	//fmt.Println(poAsBytes)
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	//fmt.Print("poIndex : ")
	//fmt.Println(poIndex)
	jsonResp = "{"
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all PO")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		//fmt.Print("valueAsBytes : ")
		//fmt.Println(valueAsBytes)
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i < len(poIndex)-1 {
			jsonResp = jsonResp + ","
		}
	}
	//fmt.Println("len(poIndex) : ")
	//fmt.Println(len(poIndex))
	jsonResp = jsonResp + "}"
	//fmt.Println("jsonResp : " + jsonResp)
	//fmt.Print("jsonResp in bytes : ")
	//fmt.Println([]byte(jsonResp))
	fmt.Println("end get_AllPO")
	return []byte(jsonResp), nil
											//send it onward
}
// ============================================================================================================================
// Delete - remove a PO from chain
// ============================================================================================================================
func (t *ManagePO) delete_po(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 'transId' as an argument\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set transId
	transId := args[0]
	err := stub.DelState(transId)													//remove the PO from chaincode
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to delete state\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	//get the PO index
	poAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to get PO index\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	//fmt.Println("poAsBytes in delete po")
	//fmt.Println(poAsBytes);
	var poIndex []string
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	//fmt.Println("poIndex in delete po")
	//fmt.Println(poIndex);
	//remove marble from index
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + transId)
		if val == transId{															//find the correct PO
			fmt.Println("found PO with matching transId")
			poIndex = append(poIndex[:i], poIndex[i+1:]...)			//remove it
			for x:= range poIndex{											//debug prints...
				fmt.Println(string(x) + " - " + poIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(poIndex)									//save new index
	err = stub.PutState(POIndexStr, jsonAsBytes)

	tosend := "{ \"transID\" : \""+transId+"\", \"message\" : \"PO deleted succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("PO deleted succcessfully")
	return nil, nil
}
// ============================================================================================================================
// Write - update PO into chaincode state
// ============================================================================================================================
func (t *ManagePO) update_po(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("Updating PO")
	if len(args) != 13 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 13\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set transId
	transId := args[0]
	poAsBytes, err := stub.GetState(transId)									//get the PO for the specified transId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to get state for " + transId + "\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	res := PO{}
	json.Unmarshal(poAsBytes, &res)
	if res.TransID == transId{
		fmt.Println("PO found with transId : " + transId)
		res.SellerName = args[1]
		res.BuyerName = args[2]
		res.ExpectedDeliveryDate = args[3]
		res.PO_date = args[4]
		res.PO_status = args[5]
		res.ItemId = args[6]
		res.Item_name = args[7]
		res.Item_quantity = args[8]
		res.Price = args[9]
		res.Buyer_sign = args[10]
		res.Seller_sign = args[11]
		res.Seller_Remarks = args[12]
	}else{
		errMsg := "{ \"message\" : \""+ transId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	//build the PO json string manually
	po_json := 	`{`+
		`"transId": "` + res.TransID + `" , `+
		`"sellerName": "` + res.SellerName + `" , `+
		`"buyerName": "` + res.BuyerName + `" , `+
		`"expectedDeliveryDate": "` + res.ExpectedDeliveryDate + `" , `+ 
		`"po_date": "` + res.PO_date + `" , `+ 
		`"po_status": "` + res.PO_status + `" , `+ 
		`"item_id": "` + res.ItemId + `" , `+ 
		`"item_name": "` + res.Item_name + `" , `+ 
		`"item_quantity": "` +  res.Item_quantity + `", `+ 
		`"price": "` + res.Price + `" , `+ 
		`"buyer_sign": "` + res.Buyer_sign + `" , `+ 
		`"seller_sign": "` + res.Seller_sign + `" , `+ 
		`"seller_remarks": "` +  res.Seller_Remarks + `" `+ 
	`}`
	err = stub.PutState(transId, []byte(po_json))									//store PO with id as key
	if err != nil {
		return nil, err
	}

	tosend := "{ \"transID\" : \""+transId+"\", \"message\" : \"PO updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("PO updated succcessfully")
	return nil, nil
}
// ============================================================================================================================
// create PO - create a new PO, store into chaincode state
// ============================================================================================================================
func (t *ManagePO) create_po(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 12 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 12\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("start create_po")
	/*if len(args[0]) <= 0 {
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
		}*/
		transId := args[0]
		sellerName := args[1]
		buyerName := args[2]
		expectedDeliveryDate := args[3]
		po_date := args[4]
		po_status := args[5]
		item_id := args[6]
		item_name := args[7]
		item_quantity := args[8]
		price := args[9]
		buyer_sign := args[10]
		seller_sign := args[11]
		seller_remarks := "NA"

		poAsBytes, err := stub.GetState(transId)
		if err != nil {
			return nil, errors.New("Failed to get PO transID")
		}
	
		res := PO{}
		json.Unmarshal(poAsBytes, &res)
		if res.TransID == transId{
			errMsg := "{ \"message\" : \"This PO arleady exists\", \"code\" : \"503\"}"
			err := stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
		return nil, nil				//all stop a PO by this name exists
	}
	
	//build the PO json string manually
	po_json := 	`{`+
		`"transId": "` + transId + `" , `+
		`"sellerName": "` + sellerName + `" , `+
		`"buyerName": "` + buyerName + `" , `+
		`"expectedDeliveryDate": "` + expectedDeliveryDate + `" , `+ 
		`"po_date": "` + po_date + `" , `+ 
		`"po_status": "` + po_status + `" , `+ 
		`"item_id": "` + item_id + `" , `+ 
		`"item_name": "` + item_name + `" , `+ 
		`"item_quantity": "` +  item_quantity + `", `+ 
		`"price": "` + price + `" , `+ 
		`"buyer_sign": "` + buyer_sign + `" , `+ 
		`"seller_sign": "` + seller_sign + `" , `+ 
		`"seller_remarks": "` +  seller_remarks + `" `+ 
	`}`
	
	fmt.Print("po_json in bytes array: ")
	fmt.Println([]byte(po_json))
	err = stub.PutState(transId, []byte(po_json))									//store PO with transId as key
	if err != nil {
		return nil, err
	}
	//get the PO index
	poIndexAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get PO index")
	}
	var poIndex []string
	//fmt.Print("poIndexAsBytes: ")
	//fmt.Println(poIndexAsBytes)
	
	json.Unmarshal(poIndexAsBytes, &poIndex)							//un stringify it aka JSON.parse()
	//fmt.Print("poIndex after unmarshal..before append: ")
	//fmt.Println(poIndex)
	//append
	poIndex = append(poIndex, transId)									//add PO transID to index list
	//fmt.Println("! PO index after appending transId: ", poIndex)
	jsonAsBytes, _ := json.Marshal(poIndex)
	//fmt.Print("jsonAsBytes: ")
	//fmt.Println(jsonAsBytes)
	err = stub.PutState(POIndexStr, jsonAsBytes)						//store name of PO
	if err != nil {
		return nil, err
	}

	tosend := "{ \"transID\" : \""+transId+"\", \"message\" : \"PO created succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("end create_po")
	return nil, nil
}
