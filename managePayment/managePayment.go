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

var PaymentIndexStr = "_PaymentIndex"	//name for the key/value that will store a list of all known payments

var AccountIndexStr = "_AccountIndex"	//name for the key/value that will store a list of all known accounts
var BuyerAccountNumber = "965832147012"
var SellerAccountNumber = "741258963512"

type Payment struct{
	PaymentID string `json:"paymentId"`					//the fieldtags are needed to keep case from bouncing around
	AgreementID string `json:"agreementId"`
	BuyerName string `json:"buyerName"`					//the fieldtags are needed to keep case from bouncing around
	SellerName string `json:"sellerName"`
	BuyerAccount string `json:"buyerAccount"`
	SellerAccount string `json:"sellerAccount"`
	AmountTransferred string `json:"amountTransferred"`
	PaymentStatus string `json:"paymentStatus"`
	PaymentCUDate string `json:"paymentCUDate"`
	PaymentDeadlineDate string `json:"paymentDeadlineDate"`
	BuyerBank_sign string `json:"buyerBank_sign"`
	BB_name string `json:"bb_name"`
	SB_name string `json:"sb_name"`
}

type AccountInfo struct{
	BuyerAccountNumber string `json:"buyerAccountNumber"`
	BuyerAccountBalance string `json:"buyerAccountBalance"`
	SellerAccountNumber string `json:"sellerAccountNumber"`
	SellerAccountBalance string `json:"sellerAccountBalance"`
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
	var balance string
	var err error

	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"Initial_Value\" as an argument.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// Initialize the chaincode
	
	balance = args[0]
	fmt.Println("ManagePayment chaincode is deployed successfully.")

	//build the Account json string manually
	account := `{`+
		`"buyerAccountNumber" : "` +  BuyerAccountNumber  + `", `+
		`"buyerAccountBalance" : "` + balance   + `", `+
		`"sellerAccountNumber" : "` +  SellerAccountNumber  + `", `+
		`"sellerAccountBalance" : "` + balance   + `"`+
		`}`

	err = stub.PutState(AccountIndexStr, []byte(account))			//store Account with id as key
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(PaymentIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	tosend := "{ \"message\" : \"ManagePayment chaincode is deployed successfully.\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
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
	}else if function == "deletePayment" {									//create a new payment
		return t.deletePayment(stub, args)
	}else if function == "updatePayment" {									//create a new trade order
		return t.updatePayment(stub, args)
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
	} else if function == "getAccountDetails" {													//read a variable
		return t.getAccountDetails(stub, args)
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
// getPaymentByID - display Payment details for a specific ID from chaincode state
// ============================================================================================================================
func (t *ManagePayment) getPaymentByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var paymentId string
	var err error
	fmt.Println("start getPaymentByID")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"paymentID\" as an argument.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set paymentId
	paymentId = args[0]
	valAsbytes, err := stub.GetState(paymentId)									//get the var from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \""+ paymentId + " not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Print("valAsbytes : ")
	fmt.Println(valAsbytes)
	fmt.Println("end getPaymentByID")
	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
//  getPaymentByBuyer - get Payment details by buyer name from chaincode state
// ============================================================================================================================
func (t *ManagePayment) getPaymentByBuyer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var errResp, jsonResp, buyerName string
	var paymentIndex []string
	var valIndex Payment
	var err error
	fmt.Println("start getPaymentByBuyer")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"Buyer_Name\" arguments.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	// set buyer name
	buyerName = args[0]
	fmt.Println("buyerName : " + buyerName)
	paymentAsBytes, err := stub.GetState(PaymentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Payment index")
	}
	fmt.Print("paymentAsBytes : ")
	fmt.Println(paymentAsBytes)
	json.Unmarshal(paymentAsBytes, &paymentIndex)								//un stringify it aka JSON.parse()
	fmt.Print("paymentIndex : ")
	fmt.Println(paymentIndex)
	fmt.Println("len(paymentIndex) : ")
	fmt.Println(len(paymentIndex))
	jsonResp = "{"
	for i,val := range paymentIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getPaymentByBuyer")
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
		if valIndex.BuyerName == buyerName{
			fmt.Println("Buyer found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			if i < len(paymentIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}else{
			errMsg := "{ \"message\" : \""+ buyerName+ " Not Found.\", \"code\" : \"503\"}"
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
	fmt.Println("end getPaymentByBuyer")
	return []byte(jsonResp), nil													//send it onward
}
// ============================================================================================================================
//  getPaymentBySeller - display Payment details for a specific Seller from chaincode state
// ============================================================================================================================
func (t *ManagePayment) getPaymentBySeller(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var errResp, sellerName, jsonResp string
	var paymentIndex []string
	var valIndex Payment
	var err error
	fmt.Println("start getPaymentBySeller")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"Seller_Name\" arguments.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set seller name
	sellerName = args[0]
	fmt.Println("sellerName: " + sellerName)
	paymentAsBytes, err := stub.GetState(PaymentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Payment index")
	}
	fmt.Print("paymentAsBytes : ")
	fmt.Println(paymentAsBytes)
	json.Unmarshal(paymentAsBytes, &paymentIndex)								//un stringify it aka JSON.parse()
	fmt.Print("paymentIndex : ")
	fmt.Println(paymentIndex)
	fmt.Println("len(paymentIndex) : ")
	fmt.Println(len(paymentIndex))
	jsonResp = "{"
	for i,val := range paymentIndex{
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
		if valIndex.SellerName == sellerName{
			fmt.Println("Seller found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			if i < len(paymentIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}else{
			errMsg := "{ \"message\" : \""+ sellerName+ " Not Found.\", \"code\" : \"503\"}"
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
	fmt.Println("end getPaymentBySeller")

	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
//  getAllPayment- display details of all Payment from chaincode state
// ============================================================================================================================
func (t *ManagePayment) getAllPayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var paymentIndex []string
	var err error
	fmt.Println("start getAllPayment")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \" \" arguments.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	paymentAsBytes, err := stub.GetState(PaymentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Payment index")
	}
	fmt.Print("paymentAsBytes : ")
	fmt.Println(paymentAsBytes)
	json.Unmarshal(paymentAsBytes, &paymentIndex)								//un stringify it aka JSON.parse()
	fmt.Print("paymentIndex : ")
	fmt.Println(paymentIndex)
	jsonResp = "{"
	for i,val := range paymentIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Payment")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i < len(paymentIndex)-1 {
			jsonResp = jsonResp + ","
		}
	}
	fmt.Println("len(paymentIndex) : ")
	fmt.Println(len(paymentIndex))
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("end getAllPayment")
	return []byte(jsonResp), nil
											//send it onward
}
// ============================================================================================================================
//  getAccountDetails - get account details from chaincode
// ============================================================================================================================
func (t *ManagePayment) getAccountDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("start getAccountDetails")
	
	accountAsBytes, err := stub.GetState(AccountIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Account index")
	}
	fmt.Print("accountAsBytes : ")
	fmt.Println(accountAsBytes)
	
	fmt.Println("end getAccountDetails")
	return accountAsBytes, nil													//send it onward
}
// ============================================================================================================================
//  updateBalance - updates balance for accounts into chaincode
// ============================================================================================================================
func (t *ManagePayment) updateBalance(stub shim.ChaincodeStubInterface, transferAmt string) ([]byte, error) {
	
	var accountIndex AccountInfo
	var err error
	fmt.Println("start updateBalance")
	
	transferAmount, _ := strconv.ParseFloat(transferAmt, 64)
	fmt.Println("In updateBalance with transferAmount::")
	fmt.Println(transferAmount)

	//Get State
	accountAsBytes, err := stub.GetState(AccountIndexStr)			//get the var from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \"Error while fetching Accounts.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	json.Unmarshal(accountAsBytes, &accountIndex)
	accountBuyerBal, _ := strconv.ParseFloat(accountIndex.BuyerAccountBalance, 64)
	accountSellerBal, _ := strconv.ParseFloat(accountIndex.SellerAccountBalance, 64)
	buyerAccountBalance	:= accountBuyerBal - transferAmount
	sellerAccountBalance := accountSellerBal + transferAmount
	fmt.Println("In updateBalance after balance calculation::")
	fmt.Println(buyerAccountBalance)
	fmt.Println(sellerAccountBalance)
	
	//build the Account json string manually
	account := `{`+
		`"buyerAccountNumber" : "` +  BuyerAccountNumber  + `", `+
		`"buyerAccountBalance" : "` + strconv.FormatFloat(buyerAccountBalance, 'f', 2, 64)   + `", `+
		`"sellerAccountNumber" : "` +  SellerAccountNumber  + `", `+
		`"sellerAccountBalance" : "` + strconv.FormatFloat(sellerAccountBalance, 'f', 2, 64)   + `" `+
		`}`
	fmt.Println("In updateBalance account to commit::" + account)

	err = stub.PutState(AccountIndexStr, []byte(account))			//store Account with id as key
	if err != nil {
		return nil, err
	}
	
	fmt.Println("end updateBalance")
	return nil, nil													//send it onward
}
// ============================================================================================================================
// Delete - remove a Payment from state
// ============================================================================================================================
func (t *ManagePayment) deletePayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"paymentID\" arguments.\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set paymentId
	paymentId := args[0]
	err := stub.DelState(paymentId)													//remove the key from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to delete state\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	//get the payment index
	paymentAsBytes, err := stub.GetState(PaymentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Payment index")
	}
	fmt.Println("paymentAsBytes in delete payment")
	fmt.Println(paymentAsBytes);
	var paymentIndex []string
	json.Unmarshal(paymentAsBytes, &paymentIndex)								//un stringify it aka JSON.parse()
	fmt.Println("paymentIndex in delete payment")
	fmt.Println(paymentIndex);
	//remove payment from index
	for i,val := range paymentIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + paymentId)
		if val == paymentId{															//find the correct payment
			fmt.Println("found payment")
			paymentIndex = append(paymentIndex[:i], paymentIndex[i+1:]...)			//remove it
			for x:= range paymentIndex{											//debug prints...
				fmt.Println(string(x) + " - " + paymentIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(paymentIndex)									//save new index
	err = stub.PutState(PaymentIndexStr, jsonAsBytes)
	tosend := "{ \"paymentID\" : \""+paymentId+"\", \"message\" : \"Payment deleted succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}

// ============================================================================================================================
// Write - update Payment into chaincode state
// ============================================================================================================================
func (t *ManagePayment) updatePayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error
	fmt.Println("running updatePayment()")

	if len(args) != 13 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 13 arguments.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	//set paymentId
	paymentId := args[0]
	paymentAsBytes, err := stub.GetState(paymentId)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + paymentId + "\"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Print("paymentAsBytes in update payment")
	fmt.Println(paymentAsBytes);
	res := Payment{}
	json.Unmarshal(paymentAsBytes, &res)
	if res.PaymentID == paymentId{
		fmt.Println("Payment found with id : " + paymentId)
		fmt.Println(res);

		res.AgreementID = args[1]
		res.BuyerName = args[2]
		res.SellerName = args[3]
		res.BuyerAccount = args[4]
		res.SellerAccount = args[5]
		res.AmountTransferred = args[6]
		res.PaymentCUDate = args[7]
		res.PaymentStatus = args[8]
		res.PaymentDeadlineDate = args[9]
		res.BuyerBank_sign = args[10]
		res.BB_name = args[11]
		res.SB_name = args[12]
	}else{
		errMsg := "{ \"message\" : \""+ paymentId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	//build the Payment json string manually
	order := `{`+
		`"paymentId" : "` + res.PaymentID   + `", `+
		`"agreementId" : "` + res.AgreementID   + `", `+
		`"buyerName" : "` + res.BuyerName   + `", `+
		`"sellerName" : "` + res.SellerName   + `", `+
		`"buyerAccount" : "` + res.BuyerAccount   + `", `+
		`"sellerAccount" : "` + res.SellerAccount   + `", `+
		`"amountTransferred" : "` + res.AmountTransferred   + `", `+
		`"paymentCUDate" : "` + res.PaymentCUDate   + `", `+
		`"paymentStatus" : "` + res.PaymentStatus   + `", `+
		`"paymentDeadlineDate" : "` + res.PaymentDeadlineDate   + `", `+
		`"buyerBank_sign" : "` + res.BuyerBank_sign   + `", `+
		`"bb_name" : "` + res.BB_name   + `", `+
		`"sb_name" : "` + res.SB_name   + `"`+
		`}`

	if res.BuyerBank_sign == "true"{
		fmt.Println("Buyer Bank sign is true with amount to be transferred :: " + res.AmountTransferred)
		t.updateBalance(stub, res.AmountTransferred)
	}

	err = stub.PutState(paymentId, []byte(order))									//store Payment with id as key
	if err != nil {
		return nil, err
	}

	tosend := "{ \"paymentID\" : \""+paymentId+"\", \"message\" : \"Payment updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 	
	fmt.Println("end updatePayment()")
	return nil, nil
}

// ============================================================================================================================
// Init Payment- create a new Payment, store into chaincode state
// ============================================================================================================================
func (t *ManagePayment) createPayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 13 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 13 arguments.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	//input sanitation
	fmt.Println("- start createPayment")
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
	}
	if len(args[9]) <= 0 {
		return nil, errors.New("10th argument must be a non-empty string")
	}
*/
	paymentId := args[0]
	agreementId := args[1]
	buyerName := args[2]
	sellerName := args[3]
	buyerAccount := BuyerAccountNumber
	sellerAccount := SellerAccountNumber
	amountTransferred := args[4]
	paymentCUDate := args[5]
	paymentStatus := args[6]
	paymentDeadlineDate := args[7]
	buyerBank_sign := args[8]
	bb_name := args[9]
	sb_name := args[10]

	paymentAsBytes, err := stub.GetState(paymentId)
	if err != nil {
		return nil, errors.New("Failed to get Payment paymentId")
	}
	fmt.Print("paymentAsBytes: ")
	fmt.Println(paymentAsBytes)
	res := Payment{}
	json.Unmarshal(paymentAsBytes, &res)
	fmt.Print("res: ")
	fmt.Println(res)
	if res.PaymentID == paymentId{
		fmt.Println("This Payment arleady exists: " + paymentId)
		errMsg := "{ \"message\" : \"This Payment arleady exists.\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil				//all stop a payment by this name exists
	}
	
	//build the Payment json string manually
	
	order := `{`+
		`"paymentId" : "` + paymentId   + `", `+
		`"agreementId" : "` + agreementId   + `", `+
		`"buyerName" : "` + buyerName   + `", `+
		`"sellerName" : "` + sellerName   + `", `+
		`"buyerAccount" : "` + buyerAccount   + `", `+
		`"sellerAccount" : "` + sellerAccount   + `", `+
		`"amountTransferred" : "` + amountTransferred   + `", `+
		`"paymentCUDate" : "` + paymentCUDate   + `", `+
		`"paymentStatus" : "` + paymentStatus   + `", `+
		`"paymentDeadlineDate" : "` + paymentDeadlineDate   + `", `+
		`"buyerBank_sign" : "` + buyerBank_sign   + `", `+
		`"bb_name" : "` + bb_name   + `", `+
		`"sb_name" : "` + sb_name   + `"`+
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
	fmt.Print("paymentIndexAsBytes: ")
	fmt.Println(paymentIndexAsBytes)
	
	json.Unmarshal(paymentIndexAsBytes, &paymentIndex)							//un stringify it aka JSON.parse()
	fmt.Print("paymentIndexAsBytes after unmarshal..before append: ")
	fmt.Println(paymentIndexAsBytes)
	
	//append
	paymentIndex = append(paymentIndex, paymentId)									//add Payment paymentId to index list
	fmt.Println("! Payment index: ", paymentIndex)
	jsonAsBytes, _ := json.Marshal(paymentIndex)
	fmt.Print("jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(PaymentIndexStr, jsonAsBytes)						//store name of Payment
	if err != nil {
		return nil, err
	}

	tosend := "{ \"paymentID\" : \""+paymentId+"\", \"message\" : \"Payment created succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 

	fmt.Println("end createPayment()")
	return nil, nil
}
