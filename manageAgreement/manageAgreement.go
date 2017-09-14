/*
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

// ManageAgreement example simple Chaincode implementation
type ManageAgreement struct {
}

var AgreementIndexStr = "_Agreementindex"				//name for the key/value that will store a list of all known Agreement
var FraudListIndexStr = "_FraudListIndexStr"

type Agreement struct{							// Attributes of a Agreement 
	AgreementID string `json:"agreementId"`	
	TransID string `json:"transId"`
	Agreement_status string `json:"agreement_status"`
	BuyerName string `json:"buyer_name"`					
	SellerName string `json:"seller_name"`
	ShipperName string `json:"shipper_name"`
	BB_name string `json:"bb_name"`
	SB_name string `json:"sb_name"`	
	PortAuthName string `json:"agreementPortAuth_name"`
	AgreementCU_date string `json:"agreementCU_date"`
	ItemId string `json:"item_id"`
	Item_name string `json:"item_name"`
	Item_quantity string `json:"item_quantity"`
	Total_Value string `json:"total_value"`
	Delivery_date string `json:"delivery_date"`
	ExtraCharges string `json:"extraCharges"`
	Shipper_fees string `json:"shipper_fees"`
	DocumentName string `json:"document_name"`
	DocumentURL string `json:"document_url"`
	TC_Text string `json:"tc_text "`
	Buyer_sign string `json:"buyer_sign"`
	BuyerBank_sign string `json:"buyerBank_sign"`
	Seller_sign string `json:"seller_sign"`
	SellerBank_sign string `json:"sellerBank_sign "`
	Industry string `json:"industry "`
	GoodsPrice string `json:"goodsPrice "`
}
type Fraud_list struct{
	FraudID string `json:"fraudId"`	
	FraudName string `json:"fraudName"`
}
// ============================================================================================================================
// Main - start the chaincode for Agreement management
// ============================================================================================================================
func main() {			
	err := shim.Start(new(ManageAgreement))
	if err != nil {
		fmt.Printf("Error starting Agreement management chaincode: %s", err)
	}
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManageAgreement) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var msg string
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"Intial_Value\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// Initialize the chaincode
	msg = args[0]
	fmt.Println("ManageAgreement chaincode is deployed successfully.");
	
	// Write the state to the ledger
	err = stub.PutState("abc", []byte(msg))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(AgreementIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	var frauds []string
	fraudListAsBytes, _ := json.Marshal(frauds)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(FraudListIndexStr, fraudListAsBytes)
	if err != nil {
		return nil, err
	}

	tosend := "{ \"message\" : \"ManageAgreement chaincode is deployed successfully.\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// Run - Our entry agreementint for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *ManageAgreement) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}
// ============================================================================================================================
// Invoke - Our entry agreementint for Invocations
// ============================================================================================================================
func (t *ManageAgreement) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "create_agreement" {											//create a new Agreement
		return t.create_agreement(stub, args)
	}else if function == "delete_agreement" {									// delete an Agreement
		return t.delete_agreement(stub, args)
	}else if function == "update_agreement" {									//update an Agreement
		return t.update_agreement(stub, args)
	}else if function == "update_fraud_list" {									//update an Agreement
		return t.update_fraud_list(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error
	errMsg := "{ \"message\" : \"Received unknown function invocation\", \"code\" : \"503\"}"
	err := stub.SetEvent("errEvent", []byte(errMsg))
	if err != nil {
		return nil, err
	} 
	return nil, nil	
}
// ============================================================================================================================
// Query - Our entry agreementint for Queries
// ============================================================================================================================
func (t *ManageAgreement) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getAgreement_byID" {													//Read a Agreement by AgreementID
		return t.getAgreement_byID(stub, args)
	} else if function == "getAgreement_byBuyer" {													//Read a Agreement by Buyer
		return t.getAgreement_byBuyer(stub, args)
	} else if function == "getAgreement_bySeller" {													//Read a Agreement by Seller
		return t.getAgreement_bySeller(stub, args)
	} else if function == "get_AllAgreement" {													//Read all Agreements
		return t.get_AllAgreement(stub, args)
	}else if function == "getAgreement_byShipper" {													//Read a Agreement by Shipper
		return t.getAgreement_byShipper(stub, args)
	} else if function == "getAgreement_byBuyerBank" {													//Read a Agreement by Buyer bank
		return t.getAgreement_byBuyerBank(stub, args)
	} else if function == "getAgreement_bySellerBank" {													//Read a Agreement by Seller bank
		return t.getAgreement_bySellerBank(stub, args)
	}else if function == "getAgreement_byPortAuthority" {													//Read a Agreement by Port Authority
		return t.getAgreement_byPortAuthority(stub, args)
	}else if function == "get_fraud_list" {													//Read a Agreement by Port Authority
		return t.get_fraud_list(stub, args)
	}else if function == "getApprovalStatus" {													//Read a Agreement by Port Authority
		return t.getApprovalStatus(stub, args)
	}else if function == "get_fraud_details" {													//Read a Agreement by Port Authority
		return t.get_fraud_details(stub, args[0])
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
// getAgreement_byID - get Agreement details for a specific AgreementID from chaincode state
// ============================================================================================================================
func (t *ManageAgreement) getAgreement_byID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var agreementId string
	var err error
	fmt.Println("start getAgreement_byID")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"AgreementID\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set agreementId
	agreementId = args[0]
	valAsbytes, err := stub.GetState(agreementId)									//get the agreementId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \""+ agreementId + " not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Print("valAsbytes : ")
	fmt.Println(valAsbytes)
	fmt.Println("end getAgreement_byID")
	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
//  getAgreement_byBuyer - get Agreement details by buyer's name from chaincode state
// ============================================================================================================================
func (t *ManageAgreement) getAgreement_byBuyer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, buyer_name, errResp string
	var agreementIndex []string
	var valIndex Agreement
	fmt.Println("start getAgreement_byBuyer")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"Buyer_Name\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set buyer's name
	buyer_name = args[0]
	fmt.Println("buyer_name : " + buyer_name)
	agreementAsBytes, err := stub.GetState(AgreementIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index string")
	}
	fmt.Print("agreementAsBytes : ")
	fmt.Println(agreementAsBytes)
	json.Unmarshal(agreementAsBytes, &agreementIndex)								//un stringify it aka JSON.parse()
	fmt.Print("agreementIndex : ")
	fmt.Println(agreementIndex)
	fmt.Println("len(agreementIndex) : ")
	fmt.Println(len(agreementIndex))
	jsonResp = "{"
	for i,val := range agreementIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getAgreement_byBuyer")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.BuyerName == buyer_name{
			fmt.Println("Buyer found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			if i < len(agreementIndex)-1 {
				jsonResp = jsonResp + ","
			}
		} else{
			errMsg := "{ \"message\" : \""+ buyer_name+ " Not Found.\", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		}
	}
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getAgreement_byBuyer")
	return []byte(jsonResp), nil											//send it onward
}

// ============================================================================================================================
//  getApprovalStatus - get approval details of an Agreement for a specific user from chaincode state
// ============================================================================================================================
func (t *ManageAgreement) getApprovalStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var user , result string
	var agreementIndex Agreement
	fmt.Println("Fetching Agreements")
	var err error
	if len(args) != 2 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"User\" and \" agreementID\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set user and agreementID
	user = args[0]
	agreementId := args[1]
	agreementAsBytes, err := stub.GetState(agreementId)
	if err != nil {
		errMsg := "{ \"message\" : \""+ agreementId + " not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Print("agreementAsBytes : ")
	fmt.Println(agreementAsBytes)
	json.Unmarshal(agreementAsBytes, &agreementIndex)								//un stringify it aka JSON.parse()
	if agreementIndex.SellerName == user{
		fmt.Println("Seller found")
		result = "{" + "\""+ "agreementId" + "\": \"" + agreementId + "\", \""+ "Seller_sign" + "\":\"" + string(agreementIndex.Seller_sign) + "\"}"
		fmt.Println("result: "+ result)
		if string(agreementIndex.Seller_sign) == "false"{
			// call update agreement and set seller_sign "true"
		}
	}else if agreementIndex.BuyerName == user{
		fmt.Println("Buyer found")
		fmt.Print(string(agreementIndex.Agreement_status));
		result = "{" + "\""+ "agreementId" + "\": \"" + agreementId + "\", \""+ "Buyer_sign" + "\":\"" + string(agreementIndex.Buyer_sign) + "\"}"
		fmt.Println("result: "+ result)
	}else if agreementIndex.BB_name == user{
		fmt.Println("Buyer Bank found")
		result = "{" + "\""+ "agreementId" + "\": \"" + agreementId + "\", \""+ "BuyerBank_sign" + "\":\"" + string(agreementIndex.BuyerBank_sign) + "\"}"
		fmt.Println("result: "+ result)
		if string(agreementIndex.Seller_sign) == "true"{
			// check for conditions and call update agreement and set buyer bank sign "true"
		}
	}else if agreementIndex.SB_name == user{
		fmt.Println("Seller Bank found")
		result = "{" + "\""+ "agreementId" + "\": \"" + agreementId + "\", \""+ "SellerBank_sign" + "\":\"" + string(agreementIndex.SellerBank_sign) + "\"}"
		fmt.Println("result: "+ result)
	}else{
		errMsg := "{ \"message\" : \""+ user+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("Fetched Approval Status")
	return []byte(result), nil											//send it onward
}


// ============================================================================================================================
//  getAgreement_bySeller - get Agreement details for a specific Seller from chaincode state
// ============================================================================================================================
func (t *ManageAgreement) getAgreement_bySeller(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, seller_name, errResp string
	var agreementIndex []string
	var valIndex Agreement
	fmt.Println("start getAgreement_bySeller")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"Seller_Name\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set seller name
	seller_name = args[0]
	fmt.Println("seller_name: " + seller_name)
	agreementAsBytes, err := stub.GetState(AgreementIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index")
	}
	fmt.Print("agreementAsBytes : ")
	fmt.Println(agreementAsBytes)
	json.Unmarshal(agreementAsBytes, &agreementIndex)								//un stringify it aka JSON.parse()
	fmt.Print("agreementIndex : ")
	fmt.Println(agreementIndex)
	fmt.Println("len(agreementIndex) : ")
	fmt.Println(len(agreementIndex))
	jsonResp = "{"
	for i,val := range agreementIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getting sellerName")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.SellerName == seller_name{
			fmt.Println("Seller found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			if i < len(agreementIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}else{
			errMsg := "{ \"message\" : \""+ seller_name+ " Not Found.\", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		}
		
	}
	
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getAgreement_bySeller")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
//  get_AllAgreement- get details of all Agreement from chaincode state
// ============================================================================================================================
func (t *ManageAgreement) get_AllAgreement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var agreementIndex []string
	fmt.Println("start get_AllAgreement")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \" \" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	agreementAsBytes, err := stub.GetState(AgreementIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index")
	}
	fmt.Print("agreementAsBytes : ")
	fmt.Println(agreementAsBytes)
	json.Unmarshal(agreementAsBytes, &agreementIndex)								//un stringify it aka JSON.parse()
	fmt.Print("agreementIndex : ")
	fmt.Println(agreementIndex)
	jsonResp = "{"
	for i,val := range agreementIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Agreement")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i < len(agreementIndex)-1 {
			jsonResp = jsonResp + ","
		}
	}
	fmt.Println("len(agreementIndex) : ")
	fmt.Println(len(agreementIndex))
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end get_AllAgreement")
	return []byte(jsonResp), nil
											//send it onward
}

// ============================================================================================================================
//  getAgreement_byShipper - get Agreement details for a specific Shipper from chaincode state
// ============================================================================================================================
func (t *ManageAgreement) getAgreement_byShipper(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, shipper_name, errResp string
	var agreementIndex []string
	var valIndex Agreement
	fmt.Println("start getAgreement_byShipper")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"Shipper_Name\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set Shipper name
	shipper_name = args[0]
	fmt.Println("shipper_name: " + shipper_name)
	agreementAsBytes, err := stub.GetState(AgreementIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index")
	}
	fmt.Print("agreementAsBytes : ")
	fmt.Println(agreementAsBytes)
	json.Unmarshal(agreementAsBytes, &agreementIndex)								//un stringify it aka JSON.parse()
	fmt.Print("agreementIndex : ")
	fmt.Println(agreementIndex)
	fmt.Println("len(agreementIndex) : ")
	fmt.Println(len(agreementIndex))
	jsonResp = "{"
	for i,val := range agreementIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getting sellerName")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.ShipperName == shipper_name{
			fmt.Println("Shipper found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			if i < len(agreementIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}else{
			errMsg := "{ \"message\" : \""+ shipper_name+ " Not Found.\", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		}
		
	}
	
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getAgreement_byShipper")
	return []byte(jsonResp), nil											//send it onward
}

// ============================================================================================================================
//  getAgreement_byBuyerBank - get Agreement details for a specific Buyer bank from chaincode state
// ============================================================================================================================
func (t *ManageAgreement) getAgreement_byBuyerBank(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, bb_name, errResp string
	var agreementIndex []string
	var valIndex Agreement
	fmt.Println("start getAgreement_byBuyerBank")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"Buyer_Bank_Name\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set Buyer Bank
	bb_name = args[0]
	fmt.Println("bb_name: " + bb_name)
	agreementAsBytes, err := stub.GetState(AgreementIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index")
	}
	fmt.Print("agreementAsBytes : ")
	fmt.Println(agreementAsBytes)
	json.Unmarshal(agreementAsBytes, &agreementIndex)								//un stringify it aka JSON.parse()
	fmt.Print("agreementIndex : ")
	fmt.Println(agreementIndex)
	fmt.Println("len(agreementIndex) : ")
	fmt.Println(len(agreementIndex))
	jsonResp = "{"
	for i,val := range agreementIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getting sellerName")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.BB_name == bb_name{
			fmt.Println("Buyer Bank found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			if i < len(agreementIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}else{
			errMsg := "{ \"message\" : \""+ bb_name + " Not Found.\", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		}
		
	}
	
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getAgreement_byBuyerBank")
	return []byte(jsonResp), nil											//send it onward
}

// ============================================================================================================================
//  getAgreement_bySellerBank - get Agreement details for a specific Seller bank from chaincode state
// ============================================================================================================================
func (t *ManageAgreement) getAgreement_bySellerBank(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, sb_name, errResp string
	var agreementIndex []string
	var valIndex Agreement
	fmt.Println("start getAgreement_bySellerBank")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"Seller_Bank_Name\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set seller bank 
	sb_name = args[0]
	fmt.Println("sb_name: " + sb_name)
	agreementAsBytes, err := stub.GetState(AgreementIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index")
	}
	fmt.Print("agreementAsBytes : ")
	fmt.Println(agreementAsBytes)
	json.Unmarshal(agreementAsBytes, &agreementIndex)								//un stringify it aka JSON.parse()
	fmt.Print("agreementIndex : ")
	fmt.Println(agreementIndex)
	fmt.Println("len(agreementIndex) : ")
	fmt.Println(len(agreementIndex))
	jsonResp = "{"
	for i,val := range agreementIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getting sellerName")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.SB_name == sb_name{
			fmt.Println("Seller found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			if i < len(agreementIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}else{
			errMsg := "{ \"message\" : \""+ sb_name + " Not Found.\", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		}
		
	}
	
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getAgreement_bySellerBank")
	return []byte(jsonResp), nil											//send it onward
}

// ============================================================================================================================
//  getAgreement_byPortAuthority - get Agreement details for a specific Port Authority from chaincode state
// ============================================================================================================================
func (t *ManageAgreement) getAgreement_byPortAuthority(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, agreementPortAuth_name, errResp string
	var agreementIndex []string
	var valIndex Agreement
	fmt.Println("start getAgreement_byPortAuthority")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"Port_Authority_Name\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set Port authority name
	agreementPortAuth_name = args[0]
	fmt.Println("agreementPortAuth_name: " + agreementPortAuth_name)
	agreementAsBytes, err := stub.GetState(AgreementIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index")
	}
	fmt.Print("agreementAsBytes : ")
	fmt.Println(agreementAsBytes)
	json.Unmarshal(agreementAsBytes, &agreementIndex)								//un stringify it aka JSON.parse()
	fmt.Print("agreementIndex : ")
	fmt.Println(agreementIndex)
	fmt.Println("len(agreementIndex) : ")
	fmt.Println(len(agreementIndex))
	jsonResp = "{"
	for i,val := range agreementIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getting sellerName")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.PortAuthName == agreementPortAuth_name{
			fmt.Println("Seller found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			if i < len(agreementIndex)-1 {
				jsonResp = jsonResp + ","
			}
		} else{
			errMsg := "{ \"message\" : \""+ agreementPortAuth_name + " Not Found.\", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		}
		
	}
	
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getAgreement_byPortAuthority")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
//  get_fraud_details - get Fraud details by fraud's name from chaincode state
// ============================================================================================================================
func (t *ManageAgreement) get_fraud_details(stub shim.ChaincodeStubInterface, args string) ([]byte, error) {
	var jsonResp, fraud_name, errResp string
	var fraudListIndex []string
	var valIndex Fraud_list
	fmt.Println("Fetching Fraud details.")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"Fraud_Name\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set fraud's name
	fraud_name = args
	fmt.Println("fraud_name : " + fraud_name)
	fraudListAsBytes, err := stub.GetState(FraudListIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Fraud List index string")
	}
	fmt.Print("fraudListAsBytes : ")
	fmt.Println(fraudListAsBytes)
	json.Unmarshal(fraudListAsBytes, &fraudListIndex)								//un stringify it aka JSON.parse()
	fmt.Print("fraudListIndex : ")
	fmt.Println(fraudListIndex)
	fmt.Println("len(fraudListIndex) : ")
	fmt.Println(len(fraudListIndex))
	jsonResp = "{"
	for i,val := range fraudListIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for get_fraud_details()")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		json.Unmarshal(valueAsBytes, &valIndex)
		fmt.Print("valIndex: ")
		fmt.Print(valIndex)
		if valIndex.FraudName == fraud_name{
			fmt.Println("Fraud Name found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			if i < len(fraudListIndex)-1 {
				jsonResp = jsonResp + ","
			}
		} else{
			errMsg := "{ \"message\" : \""+ fraud_name+ " Not Found.\", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		}
	}
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("Fetched Fraud details.")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
//  get_fraud_list - get Fraud list from chaincode state
// ============================================================================================================================
func (t *ManageAgreement) get_fraud_list(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var fraudListIndex []string
	fmt.Println("Fetching Fraud list.")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \" \" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fraudListAsBytes, err := stub.GetState(FraudListIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Fraud List index")
	}
	fmt.Print("fraudListAsBytes : ")
	fmt.Println(fraudListAsBytes)
	json.Unmarshal(fraudListAsBytes, &fraudListIndex)								//un stringify it aka JSON.parse()
	fmt.Print("fraudListIndex : ")
	fmt.Println(fraudListIndex)
	jsonResp = "{"
	for i,val := range fraudListIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for Fetching Fraud List")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i < len(fraudListIndex)-1 {
			jsonResp = jsonResp + ","
		}
	}
	fmt.Println("len(fraudListIndex) : ")
	fmt.Println(len(fraudListIndex))
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("Fetched Fraud list.")
	return []byte(jsonResp), nil
}
// ============================================================================================================================
// Delete - remove a Agreement from chain
// ============================================================================================================================
func (t *ManageAgreement) delete_agreement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"agreementID\" as an argument.\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set agreementId
	agreementId := args[0]
	err := stub.DelState(agreementId)													//remove the Agreement from chaincode
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to delete state\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	//get the Agreement index
	agreementAsBytes, err := stub.GetState(AgreementIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index")
	}
	fmt.Println("agreementAsBytes in delete agreement")
	fmt.Println(agreementAsBytes);
	var agreementIndex []string
	json.Unmarshal(agreementAsBytes, &agreementIndex)								//un stringify it aka JSON.parse()
	fmt.Println("agreementIndex in delete agreement")
	fmt.Println(agreementIndex);
	//remove agreement from index
	for i,val := range agreementIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + agreementId)
		if val == agreementId{															//find the correct Agreement
			fmt.Println("found Agreement with matching agreementId")
			agreementIndex = append(agreementIndex[:i], agreementIndex[i+1:]...)			//remove it
			for x:= range agreementIndex{											//debug prints...
				fmt.Println(string(x) + " - " + agreementIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(agreementIndex)									//save new index
	err = stub.PutState(AgreementIndexStr, jsonAsBytes)
	tosend := "{ \"agreementID\" : \""+agreementId+"\", \"message\" : \"Agreement deleted succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// Write - update Agreement into chaincode state
// ============================================================================================================================
func (t *ManageAgreement) update_agreement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error
	fmt.Println("start update_agreement")
	if len(args) != 26{
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 26 arguments.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set agreementId
	agreementId := args[0]
	agreementAsBytes, err := stub.GetState(agreementId)									//get the Agreement for the specified agreementId from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + agreementId + "\"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Print("agreementAsBytes in update agreement")
	fmt.Println(agreementAsBytes);
	res := Agreement{}
	json.Unmarshal(agreementAsBytes, &res)

	if res.AgreementID == agreementId{
		fmt.Println("Agreement found with agreementId : " + agreementId)
		fmt.Println(res);
		
		res.TransID = args[1]
		res.Agreement_status = args[2]
		res.BuyerName = args[3]
		res.SellerName = args[4]
		res.ShipperName = args[5]
		res.BB_name = args[6]
		res.SB_name	= args[7]
		res.PortAuthName = args[8]
		res.AgreementCU_date = args[9]
		res.ItemId = args[10]
		res.Item_name = args[11]
		res.Item_quantity = args[12]
		res.Total_Value = args[13]
		res.Delivery_date = args[14]
		res.ExtraCharges = args[15]
		res.Shipper_fees = args[16]
		res.DocumentName = args[17]
		res.DocumentURL = args[18]
		res.TC_Text = args[19]
		res.Buyer_sign = args[20]
		res.BuyerBank_sign = args[21]
		res.Seller_sign = args[22]
		res.SellerBank_sign = args[23]
		res.Industry = args[24]
		res.GoodsPrice = args[25]
		
		totalValue,err := strconv.Atoi(res.Total_Value)
		if err != nil {
			return nil, errors.New("Error while converting string 'total_value' to int ")
		}

		// Auto Approval
		/*if (totalValue <= 10000 || res.Industry == "Books" || res.Industry == "Mobiles & Tablets"){
			res.BuyerBank_sign = "true";
			if (res.Industry == "Books" || res.Industry == "Mobiles & Tablets"){
				res.SellerBank_sign = "true";
			}
		}*/
		if (totalValue <= 10000 && (res.Industry == "Books" || res.Industry == "Mobiles & Tablets")){
			res.BuyerBank_sign = "true";
		}
		if (res.Industry == "Books" || res.Industry == "Mobiles & Tablets"){
			res.SellerBank_sign = "true";
		}
		if(res.BuyerBank_sign == "true" && res.Seller_sign == "false" && res.SellerBank_sign == "false"){
			res.Agreement_status = "Approved By Buyer Bank"
		}
		if(res.BuyerBank_sign == "true" && res.Seller_sign == "true" && res.SellerBank_sign == "false"){
			res.Agreement_status = "Approved By Seller"
		}
		if(res.BuyerBank_sign == "true" && res.Seller_sign == "true" && res.SellerBank_sign == "true"){
			res.Agreement_status = "Approved By Seller Bank"
		}
		
	}else{
		errMsg := "{ \"message\" : \""+ agreementId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	//build the Agreement json string manually
	input := 	`{`+
		`"agreementId": "` + res.AgreementID + `" , `+
		`"transId": "` + res.TransID + `" , `+ 
		`"agreement_status": "` + res.Agreement_status + `" , `+ 
		`"buyer_name": "` + res.BuyerName + `" , `+
		`"seller_name": "` + res.SellerName + `" , `+
		`"shipper_name": "` + res.ShipperName + `" , `+ 
		`"bb_name": "` + res.BB_name + `" , `+ 
		`"sb_name": "` + res.SB_name + `" , `+ 
		`"agreementPortAuth_name": "` + res.PortAuthName + `" , `+ 
		`"agreementCU_date": "` + res.AgreementCU_date + `" , `+ 
		`"item_id": "` + res.ItemId + `" , `+ 
		`"item_name": "` + res.Item_name + `" , `+ 
		`"item_quantity": "` + res.Item_quantity + `" , `+ 
		`"total_value": "` + res.Total_Value + `" , `+ 
		`"delivery_date": "` + res.Delivery_date + `" , `+ 
		`"extraCharges": "` + res.ExtraCharges + `" , `+ 
		`"shipper_fees" : "` + res.Shipper_fees + `", `+ 
		`"document_name": "` + res.DocumentName + `" , `+ 
		`"document_url": "` + res.DocumentURL + `" , `+ 
		`"tc_text" : "` + res.TC_Text + `" , `+ 
		`"buyer_sign": "` + res.Buyer_sign + `" , `+ 
		`"buyerBank_sign": "` + res.BuyerBank_sign + `" , `+ 
		`"seller_sign": "` + res.Seller_sign + `" , `+ 
		`"sellerBank_sign" : "` + res.SellerBank_sign + `" , `+ 
		`"industry" : "` + res.Industry + `" , `+
		`"goodsPrice" : "` + res.GoodsPrice + `" `+ 
		`}`
	err = stub.PutState(agreementId, []byte(input))									//store Agreement with id as key
	if err != nil {
		return nil, err
	}
	tosend := "{ \"agreementID\" : \""+agreementId+"\", \"message\" : \"Agreement updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	fmt.Println("end update_agreement")
	return nil, nil
}
// ============================================================================================================================
// create Agreement - create a new Agreement, store into chaincode state
// ============================================================================================================================
func (t *ManageAgreement) create_agreement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 26 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 26 arguments.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("start create_agreement")
	
		agreementId := args[0]
		transId := args[1]
		agreement_status := args[2]
		buyer_name := args[3]
		seller_name := args[4]
		shipper_name := args[5]
		bb_name := args[6]
		sb_name	:= args[7]
		agreementPortAuth_name := args[8]
		agreementCU_date := args[9]
		item_id := args[10]
		item_name := args[11]
		item_quantity:= args[12]
		total_value := args[13]
		delivery_date := args[14]
		extraCharges := args[15]
		shipper_fees := args[16]
		document_name := args[17]
		document_url := args[18]
		tc_text := args[19]
		buyer_sign := args[20]
		buyerBank_sign := args[21]
		seller_sign := args[22]
		sellerBank_sign := args[23]
		industry := args[24]
		goodsPrice := args[25]
		
		fmt.Println("Checking fraud list...");

		buyer, err:= t.get_fraud_details(stub, buyer_name)
		if buyer != nil{
			errMsg := "{ \"Agreement ID\" : \""+agreementId+"\", \"message\" : \"Buyer name exists in Fraud list. So, Agreement auto-rejected by System. \", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		} else if err != nil{
			errMsg := "{ \"message\" : \"Error while checking for Buyer in Fraud list. \", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		}
		seller, err := t.get_fraud_details(stub, seller_name)
		if seller != nil{
			errMsg := "{ \"Agreement ID\" : \""+agreementId+"\", \"message\" : \"Seller name exists in Fraud list. So, Agreement auto-rejected by System. \", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		}else if err != nil{
			errMsg := "{ \"message\" : \"Error while checking for Seller in Fraud list. \", \"code\" : \"503\"}"
			err = stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
			return nil, nil
		}
		fmt.Println("Checked fraud list successfully.");

		agreementAsBytes, err := stub.GetState(agreementId)
		if err != nil {
			return nil, errors.New("Failed to get Agreement ID")
		}
		fmt.Print("agreementAsBytes: ")
		fmt.Println(agreementAsBytes)
		res := Agreement{}
		json.Unmarshal(agreementAsBytes, &res)
		fmt.Print("res: ")
		fmt.Println(res)
		if res.AgreementID == agreementId{
			fmt.Println("This Agreement already exists: " + agreementId)
			errMsg := "{ \"message\" : \"This Agreement already exists.\", \"code\" : \"503\"}"
			err := stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
		return nil, nil				//all stop a Agreement by this name exists
	}
	
	//build the Agreement json string manually
	input := 	`{`+
		`"agreementId": "` + agreementId + `" , `+
		`"transId": "` + transId + `" , `+ 
		`"agreement_status": "` + agreement_status + `" , `+ 
		`"buyer_name": "` + buyer_name + `" , `+
		`"seller_name": "` + seller_name + `" , `+
		`"shipper_name": "` + shipper_name + `" , `+ 
		`"bb_name": "` + bb_name + `" , `+ 
		`"sb_name": "` + sb_name + `" , `+ 	
		`"agreementPortAuth_name": "` + agreementPortAuth_name + `" , `+ 
		`"agreementCU_date": "` + agreementCU_date + `" , `+ 
		`"item_id": "` + item_id + `" , `+ 
		`"item_name": "` + item_name + `" , `+ 
		`"item_quantity": "` + item_quantity + `" , `+ 
		`"total_value": "` + total_value + `" , `+ 
		`"delivery_date": "` + delivery_date + `" , `+ 
		`"extraCharges": "` + extraCharges + `" , `+ 
		`"shipper_fees": "` + shipper_fees + `", `+ 
		`"document_name": "` + document_name + `" , `+ 
		`"document_url": "` + document_url + `" , `+ 
		`"tc_text": "` + tc_text + `" , `+ 
		`"buyer_sign": "` + buyer_sign + `" , `+ 
		`"buyerBank_sign": "` + buyerBank_sign + `" , `+ 
		`"seller_sign": "` + seller_sign + `" , `+ 
		`"sellerBank_sign": "` + sellerBank_sign + `", `+ 
		`"industry": "` + industry + `" , `+
		`"goodsPrice": "` + goodsPrice + `" `+
		`}`
		fmt.Println("input: " + input)
		fmt.Print("input in bytes array: ")
		fmt.Println([]byte(input))
	err = stub.PutState(agreementId, []byte(input))									//store Agreement with agreementId as key
	if err != nil {
		return nil, err
	}
	//get the Agreement index
	agreementIndexAsBytes, err := stub.GetState(AgreementIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Agreement index")
	}
	var agreementIndex []string
	fmt.Print("agreementIndexAsBytes: ")
	fmt.Println(agreementIndexAsBytes)
	
	json.Unmarshal(agreementIndexAsBytes, &agreementIndex)							//un stringify it aka JSON.parse()
	fmt.Print("agreementIndex after unmarshal..before append: ")
	fmt.Println(agreementIndex)
	//append
	agreementIndex = append(agreementIndex, agreementId)									//add Agreement transID to index list
	fmt.Println("! Agreement index after appending agreementId: ", agreementIndex)
	jsonAsBytes, _ := json.Marshal(agreementIndex)
	fmt.Print("jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(AgreementIndexStr, jsonAsBytes)						//store name of Agreement
	if err != nil {
		return nil, err
	}

	tosend := "{ \"agreementID\" : \""+agreementId+"\", \"message\" : \"Agreement created succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	}

	fmt.Println("end create_agreement")
	return nil, nil
}
// ============================================================================================================================
// create Fraud_list - add an entry in the farus list, store into chaincode state
// ============================================================================================================================
func (t *ManageAgreement) update_fraud_list(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// update_fraud_list("fraudID","fraudName")
	var err error
	if len(args) != 2 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 2 arguments.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("Updating Fraud list.")
	
	fraudId := args[0]
	fraudName := args[1]

	fraudListAsBytes, err := stub.GetState(fraudId)
	if err != nil {
		return nil, errors.New("Failed to get fraudID")
	}
	fmt.Print("fraudListAsBytes: ")
	fmt.Println(fraudListAsBytes)
	res := Fraud_list{}
	json.Unmarshal(fraudListAsBytes, &res)
	fmt.Print("res: ")
	fmt.Println(res)
	if res.FraudID == fraudId{
		fmt.Println("This Fraud Name already exists: " + fraudId)
		errMsg := "{ \"message\" : \"This Fraud Name already exists.\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil				//all stop a Fraud List by this name exists
	}
	
	//build the fraud json string manually
	input := 	`{`+
		`"fraudId": "` + fraudId + `" , `+
		`"fraudName": "` + fraudName  + `" `+ 
		`}`
		fmt.Println("input: " + input)
		fmt.Print("input in bytes array: ")
		fmt.Println([]byte(input))
	err = stub.PutState(fraudId, []byte(input))									//store Fraud with fraudId as key
	if err != nil {
		return nil, err
	}
	//get the Fraud List index
	fraudListIndexAsBytes, err := stub.GetState(FraudListIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Fraud List index")
	}
	var fraudListIndex []string
	fmt.Print("fraudListIndexAsBytes: ")
	fmt.Println(fraudListIndexAsBytes)
	
	json.Unmarshal(fraudListIndexAsBytes, &fraudListIndex)							//un stringify it aka JSON.parse()
	//append
	fraudListIndex = append(fraudListIndex, fraudId)									//add fraudId to index list
	fmt.Println("! fraud List index after appending fraudId: ", fraudListIndex)
	jsonAsBytes, _ := json.Marshal(fraudListIndex)
	fmt.Print("jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(FraudListIndexStr, jsonAsBytes)						//store name of Agreement
	if err != nil {
		return nil, err
	}
	tosend := "{ \"Fraud ID\" : \""+fraudId+"\", \"message\" : \"Fraud ID added succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("Fraud list updated successfully.")
	return nil, nil
}
/*func (t *ManageAgreement) approve_agreement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	/*var jsonResp , str string
	var err error
	fmt.Println("start approve_agreement")
	if len(args) != 3{
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 3 arguments.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set agreementId
	agreementId := args[0]
	bb_name := args[1]
	sb_name := args[2]
	//sign := args[2]
	agreementAsBytes, err := stub.GetState(agreementId)									//get the Agreement for the specified agreementId from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + agreementId + "\"}"
		return nil, errors.New(jsonResp)
	}
	res := Agreement{}
	json.Unmarshal(agreementAsBytes, &res)
	if res.AgreementID == agreementId{
		fmt.Println("Agreement found with agreementId : " + agreementId)
		fmt.Println("res: "+res);
		if res.BB_name == bb_name {
			totalValue,err := strconv.Atoi(res.Total_Value)
			if err != nil {
				return nil, errors.New("Error while converting string 'Total_Value' to int ")
			}
			if (totalValue <= 10000 && res.Industry == "Books" || res.Industry == "Mobiles & Tablets"){
					res.BuyerBank_sign = "true";
					str = "Buyer Bank Signature"
			}
		} else if res.SB_name == sb_name {
			if (res.Industry == "Books" || res.Industry == "Mobiles & Tablets"){
					res.SellerBank_sign = "true";
					str = "Seller Bank Signature"
			}
		}
	}else{
		errMsg := "{ \"message\" : \""+ agreementId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}*/
	
	//build the Agreement json string manually
	/*input := 	`{`+
		`"agreementId": "` + res.AgreementID + `" , `+
		`"transId": "` + res.TransID + `" , `+ 
		`"agreement_status": "` + res.Agreement_status + `" , `+ 
		`"buyer_name": "` + res.BuyerName + `" , `+
		`"seller_name": "` + res.SellerName + `" , `+
		`"shipper_name": "` + res.ShipperName + `" , `+ 
		`"bb_name": "` + res.BB_name + `" , `+ 
		`"sb_name": "` + res.SB_name + `" , `+ 
		`"agreementPortAuth_name": "` + res.PortAuthName + `" , `+ 
		`"agreementCU_date": "` + res.AgreementCU_date + `" , `+ 
		`"item_id": "` + res.ItemId + `" , `+ 
		`"item_name": "` + res.Item_name + `" , `+ 
		`"item_quantity": "` + res.Item_quantity + `" , `+ 
		`"total_value": "` + res.Total_Value + `" , `+ 
		`"delivery_date": "` + res.Delivery_date + `" , `+ 
		`"extraCharges": "` + res.ExtraCharges + `" , `+ 
		`"shipper_fees": "` + res.Shipper_fees + `" , `+ 
		`"document_name": "` + res.DocumentName + `" , `+ 
		`"document_url": "` + res.DocumentURL + `" , `+ 
		`"tc_text" : "` + res.TC_Text + `" , `+ 
		`"buyer_sign": "` + res.Buyer_sign + `" , `+ 
		`"buyerBank_sign": "` + res.BuyerBank_sign + `" , `+ 
		`"seller_sign": "` + res.Seller_sign + `" , `+ 
		`"sellerBank_sign" : "` + res.SellerBank_sign + `" , `+ 
		`"industry" : "` + res.Industry + `" , `+ 
		`"goodsPrice" : "` + res.GoodsPrice + `" `+ 
		`}`
	fmt.Println("input: "+input)
	err = stub.PutState(agreementId, []byte(input))									//store Agreement with id as key
	if err != nil {
		return nil, err
	}*/
	/*tosend := "{ \"agreementID\" : \""+agreementId+"\", \"message\" : \"" + str + " updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	fmt.Println("end approve_agreement")
	return nil, nil
}*/
