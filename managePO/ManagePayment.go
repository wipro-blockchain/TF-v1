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

// ManagePayment example simple Chaincode implementation
type ManagePayment struct {
}

var PaymentIndexStr = "_PaymentIndex"				//name for the key/value that will store a list of all known marbles
//var openTradesStr = "_opentrades"				//name for the key/value that will store all open trades

type Payment struct{
	PaymentId string `json:"paymentId"`					//the fieldtags are needed to keep case from bouncing around
	AgreementId string `json:"agreemetId"`
	BuyerName string `json:"buyerName"`					//the fieldtags are needed to keep case from bouncing around
	SellerName string `json:"sellerName"`
	BuyerAccount string `json:"buyerAccount"`
	SellerAccount string `json:"sellerAccount"`
	AmountTransferred string `json:"amountTransferred"`
	PaymentStatus string `json:"paymentStatus"`
	PaymentCreateDate string `json:"paymentCreateDate"`
	PaymentUpdateDate string `json:"paymentUpdateDate"`
	PaymentDeadlineDate string `json:"paymentDeadlineDate"`
}


// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(ManagePayment))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManagePayment) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var msg string
	//var Payment_list allPayment

	var err error
	if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	// Initialize the chaincode
	
	msg = args[0]
	fmt.Println("ManagePayment chaincode is deployed with the message : "+ msg)
	/*Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}*/

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(msg))				//making a test var "abc" to read/write check right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(PaymentIndexStr, jsonAsBytes)
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
func (t *ManagePayment) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *ManagePayment) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "createPayment" {											//writes a value to the chaincode state
		return t.createPayment(stub, args)
	}else if function == "deletePayment" {									//create a new marble
		return t.deletePayment(stub, args)
	}else if function == "updatePayment" {									//create a new trade order
		return t.updatePayment(stub, args)
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
func (t *ManagePayment) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getPaymentByID" {													//read a variable
		return t.getPaymentByID(stub, args)
	} else if function == "getPaymentByBuyer" {													//read a variable
		return t.getPaymentByBuyer(stub, args)
	} else if function == "getPaymentBySeller" {													//read a variable
		return t.getPaymentBySeller(stub, args)
	} else if function == "getAllPayment" {													//read a variable
		return t.getAllPayment(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// getPaymentByID - display Payment details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManagePayment) getPaymentByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var paymentId, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ID of the var to query")
	}

	paymentId = args[0]
	valAsbytes, err := stub.GetState(paymentId)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + paymentId + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
//  getPaymentByBuyer - get Payment details by buyer name from chaincode state
// ============================================================================================================================
func (t *ManagePayment) getPaymentByBuyer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var errResp, buyerName string
	var paymentIndex, valIndex []string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	buyerName = args[0]
	paymentAsBytes, err := stub.GetState(PaymentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Payment index")
	}
	json.Unmarshal(paymentAsBytes, &paymentIndex)								//un stringify it aka JSON.parse()
	for i,val := range paymentIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Payment")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		json.Unmarshal(valueAsBytes, &valIndex)
		for j,value := range valIndex{
			fmt.Println(strconv.Itoa(j) + " - looking at " + value + " for all Payment")
			if value == buyerName {
				return []byte(valIndex[j]),nil
			}
		}
	}
	return nil, nil													//send it onward
}
// ============================================================================================================================
//  getPaymentBySeller - display Payment details for a specific Seller from chaincode state
// ============================================================================================================================
func (t *ManagePayment) getPaymentBySeller(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var errResp, sellerName string
	var paymentIndex, valIndex []string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	sellerName = args[0]
	paymentAsBytes, err := stub.GetState(PaymentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Payment index")
	}
	json.Unmarshal(paymentAsBytes, &paymentIndex)								//un stringify it aka JSON.parse()
	//jsonResp = "{"
	for i,val := range paymentIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Payment")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		json.Unmarshal(valueAsBytes, &valIndex)
		for j,value := range valIndex{
			fmt.Println(strconv.Itoa(j) + " - looking at " + value + " for all Payment")
			if value == sellerName {
				return []byte(valIndex[j]),nil
			}
		}
	}
	return nil, nil													//send it onward
}
// ============================================================================================================================
//  getAllPayment- display details of all Payment from chaincode state
// ============================================================================================================================
func (t *ManagePayment) getAllPayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var paymentIndex []string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	paymentAsBytes, err := stub.GetState(PaymentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Payment index")
	}
	json.Unmarshal(paymentAsBytes, &paymentIndex)								//un stringify it aka JSON.parse()
	jsonResp = "{"
	for i,val := range paymentIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Payment")
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
func (t *ManagePayment) deletePayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	paymentId := args[0]
	err := stub.DelState(paymentId)													//remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the marble index
	paymentAsBytes, err := stub.GetState(PaymentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get marble index")
	}
	fmt.Println("paymentAsBytes in delete payment")
	fmt.Println(paymentAsBytes);
	var paymentIndex []string
	json.Unmarshal(paymentAsBytes, &paymentIndex)								//un stringify it aka JSON.parse()
	fmt.Println("paymentIndex in delete payment")
	fmt.Println(paymentIndex);
	//remove marble from index
	for i,val := range paymentIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + paymentId)
		if val == paymentId{															//find the correct marble
			fmt.Println("found marble")
			paymentIndex = append(paymentIndex[:i], paymentIndex[i+1:]...)			//remove it
			for x:= range paymentIndex{											//debug prints...
				fmt.Println(string(x) + " - " + paymentIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(paymentIndex)									//save new index
	err = stub.PutState(PaymentIndexStr, jsonAsBytes)
	return nil, nil
}

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *ManagePayment) updatePayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error
	fmt.Println("running updatePayment()")

	if len(args) != 9 {
		return nil, errors.New("Incorrect number of arguments. Expecting 9.")
	}
	paymentId := args[0]
	paymentAsBytes, err := stub.GetState(paymentId)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + paymentId + "\"}"
		return nil, errors.New(jsonResp)
	}
	
	res := Payment{}
	json.Unmarshal(paymentAsBytes, &res)
	if res.PaymentId == paymentId{
		fmt.Println("Payment found with id : " + paymentId)
		fmt.Println(res);
		res.PaymentId = args[0]
		res.AgreementId = args[1]
		res.BuyerName = args[2]
		res.SellerName = args[3]
		res.BuyerAccount = args[4]
		res.SellerAccount = args[5]
		res.AmountTransferred = args[6]
		res.PaymentStatus = args[7]
		res.PaymentCreateDate = args[8]
		res.PaymentUpdateDate = args[9]
		res.PaymentDeadlineDate = args[10]
	}
	
	//build the Payment json string manually
	order := `{`+
			`"PaymentId" : "` + res.PaymentId   + `", `+
			`"AgreementId" : "` + res.AgreementId   + `", `+
			`"BuyerName" : "` + res.BuyerName   + `", `+
			`"SellerName" : "` + res.SellerName   + `", `+
			`"BuyerAccount" : "` + res.BuyerAccount   + `", `+
			`"SellerAccount" : "` + res.SellerAccount   + `", `+
			`"AmountTransferred" : "` + res.AmountTransferred   + `", `+
			`"PaymentStatus" : "` + res.PaymentStatus   + `", `+
			`"PaymentCreateDate" : "` + res.PaymentCreateDate   + `", `+
			`"PaymentUpdateDate" : "` + res.PaymentUpdateDate   + `", `+
			`"PaymentDeadlineDate " : "` + res.PaymentDeadlineDate   + `", `+
			`}`

	err = stub.PutState(paymentId, []byte(order))									//store Payment with id as key
	if err != nil {
		return nil, err
	}
		
	//get the Payment index
	/*paymentIndexAsBytes, err := stub.GetState(PaymentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Payment index")
	}
	var paymentIndex []string
	json.Unmarshal(paymentIndexAsBytes, &paymentIndex)							//un stringify it aka JSON.parse()
	
	//append
	paymentIndex = append(paymentIndex, paymentId)									//add Payment paymentId to index list
	fmt.Println("! Payment index: ", paymentIndex)
	jsonAsBytes, _ := json.Marshal(paymentIndex)
	err = stub.PutState(PaymentIndexStr, jsonAsBytes)						//store name of Payment

	fmt.Println("- end create Payment")*/
	return nil, nil
}

// ============================================================================================================================
// Init Payment- create a new Payment, store into chaincode state
// ============================================================================================================================
func (t *ManagePayment) createPayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	//var payment Payment
	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 9 {
		return nil, errors.New("Incorrect number of arguments. Expecting 9")
	}
	//input sanitation
	fmt.Println("- start createPayment")
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

	paymentId := args[0]
	agreementId := args[1]
	buyerName := args[2]
	sellerName := args[3]
	buyerAccount := args[4]
	sellerAccount := args[5]
	amountTransferred := args[6]
	paymentStatus := args[7]
	paymentCreateDate := args[8]
	paymentUpdateDate := args[9]
	paymentDeadlineDate := args[10]

	paymentAsBytes, err := stub.GetState(paymentId)
	if err != nil {
		return nil, errors.New("Failed to get Payment paymentId")
	}
	res := Payment{}
	json.Unmarshal(paymentAsBytes, &res)
	if res.PaymentId == paymentId{
		fmt.Println("This Payment arleady exists: " + paymentId)
		fmt.Println(res);
		return nil, errors.New("This Payment arleady exists")				//all stop a marble by this name exists
	}
	
	//build the Payment json string manually
	
	order := `{`+
			`"PaymentId" : "` + paymentId   + `", `+
			`"AgreementId" : "` + agreementId   + `", `+
			`"BuyerName" : "` + buyerName   + `", `+
			`"SellerName" : "` + sellerName   + `", `+
			`"BuyerAccount" : "` + buyerAccount   + `", `+
			`"SellerAccount" : "` + sellerAccount   + `", `+
			`"AmountTransferred" : "` + amountTransferred   + `", `+
			`"PaymentStatus" : "` + paymentStatus   + `", `+
			`"PaymentCreateDate" : "` + paymentCreateDate   + `", `+
			`"PaymentUpdateDate" : "` + paymentUpdateDate   + `", `+
			`"PaymentDeadlineDate " : "` + paymentDeadlineDate   + `", `+
			`}`

	err = stub.PutState(paymentId, []byte(order))									//store Payment with id as key
	if err != nil {
		return nil, err
	}
		
	//get the Payment index
	paymentIndexAsBytes, err := stub.GetState(PaymentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Payment index")
	}
	var paymentIndex []string
	json.Unmarshal(paymentIndexAsBytes, &paymentIndex)							//un stringify it aka JSON.parse()
	
	//append
	paymentIndex = append(paymentIndex, paymentId)									//add Payment paymentId to index list
	fmt.Println("! Payment index: ", paymentIndex)
	jsonAsBytes, _ := json.Marshal(paymentIndex)
	err = stub.PutState(PaymentIndexStr, jsonAsBytes)						//store name of Payment

	fmt.Println("- end create Payment")
	return nil, nil
}
