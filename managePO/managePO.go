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
	var buyerName, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ID of the var to query")
	}

	buyerName = args[0]
	valAsbytes, err := stub.GetState(buyerName)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + buyerName + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
//  getPO_bySeller - display PO details for a specific Seller from chaincode state
// ============================================================================================================================
func (t *ManagePO) getPO_bySeller(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var sellerName, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ID of the var to query")
	}

	sellerName = args[0]
	valAsbytes, err := stub.GetState(sellerName)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + sellerName + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
//  get_AllPO- display details of all PO from chaincode state
// ============================================================================================================================
func (t *ManagePO) get_AllPO(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1 argument")
	}
	poAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get PO index")
	}
	var errResp string
	var poJson []string
	fmt.Println("poAsBytes")
	fmt.Println(poAsBytes)
	var poIndex []string
	json.Unmarshal(poAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	fmt.Println("poIndex")
	fmt.Println(poIndex)
	for i,val := range poIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all PO")
		poAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		json.Unmarshal(poAsBytes, &poJson)
		fmt.Println("************* individual PO" )
		fmt.Println(poJson) 
		/*fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all PO")
		poIndex = append(poIndex[:i], poIndex[i+1:]...)	
		fmt.Println("poIndex in for loop: ")
		fmt.Println(poIndex)
		for x:= range poIndex{											//debug prints...
			fmt.Println(string(x) + " - " + poIndex[x])
		}*/
	}
	return poAsBytes, nil													//send it onward
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
	fmt.Println("poAsBytes in create_po")
	fmt.Println(poAsBytes);
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
	fmt.Println("poIndexAsBytes in create_po")
	fmt.Println(poIndexAsBytes);
	var poIndex []string
	json.Unmarshal(poIndexAsBytes, &poIndex)							//un stringify it aka JSON.parse()
	
	//append
	poIndex = append(poIndex, transId)									//add PO transID to index list
	fmt.Println("poIndex in create_po after append")
	fmt.Println(poIndex);
	fmt.Println("! PO index: ", poIndex)
	jsonAsBytes, _ := json.Marshal(poIndex)
	err = stub.PutState(POIndexStr, jsonAsBytes)						//store name of PO
	fmt.Println("jsonAsBytes in create_po after PutState(POIndexStr ")
	fmt.Println(jsonAsBytes);
	fmt.Println("- end create PO")
	return nil, nil
}

// ============================================================================================================================
// Set User Permission on PO
// ============================================================================================================================
/*func (t *ManagePO) set_user(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	
	//   0       1
	// "name", "bob"
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	fmt.Println("- start set user")
	fmt.Println(args[0] + " - " + args[1])
	marbleAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}
	res := Marble{}
	json.Unmarshal(marbleAsBytes, &res)										//un stringify it aka JSON.parse()
	res.User = args[1]														//change the user
	
	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[0], jsonAsBytes)								//rewrite the marble with id as key
	if err != nil {
		return nil, err
	}
	
	fmt.Println("- end set user")
	return nil, nil
}*/

// ============================================================================================================================
// Open Trade - create an open trade for a marble you want with marbles you have 
// ============================================================================================================================
/*func (t *ManagePO) open_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var will_size int
	var trade_away Description
	
	//	0        1      2     3      4      5       6
	//["bob", "blue", "16", "red", "16"] *"blue", "35*
	if len(args) < 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting like 5?")
	}
	if len(args)%2 == 0{
		return nil, errors.New("Incorrect number of arguments. Expecting an odd number")
	}

	size1, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}

	open := AnOpenTrade{}
	open.User = args[0]
	open.Timestamp = makeTimestamp()											//use timestamp as an ID
	open.Want.Color = args[1]
	open.Want.Size =  size1
	fmt.Println("- start open trade")
	jsonAsBytes, _ := json.Marshal(open)
	err = stub.PutState("_debug1", jsonAsBytes)

	for i:=3; i < len(args); i++ {												//create and append each willing trade
		will_size, err = strconv.Atoi(args[i + 1])
		if err != nil {
			msg := "is not a numeric string " + args[i + 1]
			fmt.Println(msg)
			return nil, errors.New(msg)
		}
		
		trade_away = Description{}
		trade_away.Color = args[i]
		trade_away.Size =  will_size
		fmt.Println("! created trade_away: " + args[i])
		jsonAsBytes, _ = json.Marshal(trade_away)
		err = stub.PutState("_debug2", jsonAsBytes)
		
		open.Willing = append(open.Willing, trade_away)
		fmt.Println("! appended willing to open")
		i++;
	}
	
	//get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)										//un stringify it aka JSON.parse()
	
	trades.OpenTrades = append(trades.OpenTrades, open);						//append to open trades
	fmt.Println("! appended open to trades")
	jsonAsBytes, _ = json.Marshal(trades)
	err = stub.PutState(openTradesStr, jsonAsBytes)								//rewrite open orders
	if err != nil {
		return nil, err
	}
	fmt.Println("- end open trade")
	return nil, nil
}*/

// ============================================================================================================================
// Perform Trade - close an open trade and move ownership
// ============================================================================================================================
/*func (t *ManagePO) perform_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	
	//	0		1					2					3				4					5
	//[data.id, data.closer.user, data.closer.name, data.opener.user, data.opener.color, data.opener.size]
	if len(args) < 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}
	
	fmt.Println("- start close trade")
	timestamp, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return nil, errors.New("1st argument must be a numeric string")
	}
	
	size, err := strconv.Atoi(args[5])
	if err != nil {
		return nil, errors.New("6th argument must be a numeric string")
	}
	
	//get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)															//un stringify it aka JSON.parse()
	
	for i := range trades.OpenTrades{																//look for the trade
		fmt.Println("looking at " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10) + " for " + strconv.FormatInt(timestamp, 10))
		if trades.OpenTrades[i].Timestamp == timestamp{
			fmt.Println("found the trade");
			
			
			marbleAsBytes, err := stub.GetState(args[2])
			if err != nil {
				return nil, errors.New("Failed to get thing")
			}
			closersMarble := Marble{}
			json.Unmarshal(marbleAsBytes, &closersMarble)											//un stringify it aka JSON.parse()
			
			//verify if marble meets trade requirements
			if closersMarble.Color != trades.OpenTrades[i].Want.Color || closersMarble.Size != trades.OpenTrades[i].Want.Size {
				msg := "marble in input does not meet trade requriements"
				fmt.Println(msg)
				return nil, errors.New(msg)
			}
			
			marble, e := findMarble4Trade(stub, trades.OpenTrades[i].User, args[4], size)			//find a marble that is suitable from opener
			if(e == nil){
				fmt.Println("! no errors, proceeding")

				t.set_user(stub, []string{args[2], trades.OpenTrades[i].User})						//change owner of selected marble, closer -> opener
				t.set_user(stub, []string{marble.Name, args[1]})									//change owner of selected marble, opener -> closer
			
				trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)		//remove trade
				jsonAsBytes, _ := json.Marshal(trades)
				err = stub.PutState(openTradesStr, jsonAsBytes)										//rewrite open orders
				if err != nil {
					return nil, err
				}
			}
		}
	}
	fmt.Println("- end close trade")
	return nil, nil
}*/

// ============================================================================================================================
// findMarble4Trade - look for a matching marble that this user owns and return it
// ============================================================================================================================
/*func findMarble4Trade(stub shim.ChaincodeStubInterface, user string, color string, size int )(m Marble, err error){
	var fail Marble;
	fmt.Println("- start find marble 4 trade")
	fmt.Println("looking for " + user + ", " + color + ", " + strconv.Itoa(size));

	//get the marble index
	marblesAsBytes, err := stub.GetState(POIndexStr)
	if err != nil {
		return fail, errors.New("Failed to get marble index")
	}
	var poIndex []string
	json.Unmarshal(marblesAsBytes, &poIndex)								//un stringify it aka JSON.parse()
	
	for i:= range poIndex{													//iter through all the marbles
		//fmt.Println("looking @ marble name: " + poIndex[i]);

		marbleAsBytes, err := stub.GetState(poIndex[i])						//grab this marble
		if err != nil {
			return fail, errors.New("Failed to get marble")
		}
		res := Marble{}
		json.Unmarshal(marbleAsBytes, &res)										//un stringify it aka JSON.parse()
		//fmt.Println("looking @ " + res.User + ", " + res.Color + ", " + strconv.Itoa(res.Size));
		
		//check for user && color && size
		if strings.ToLower(res.User) == strings.ToLower(user) && strings.ToLower(res.Color) == strings.ToLower(color) && res.Size == size{
			fmt.Println("found a marble: " + res.Name)
			fmt.Println("! end find marble 4 trade")
			return res, nil
		}
	}
	
	fmt.Println("- end find marble 4 trade - error")
	return fail, errors.New("Did not find marble to use in this trade")
}
*/
// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
/*func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}*/

// ============================================================================================================================
// Remove Open Trade - close an open trade
// ============================================================================================================================
/*func (t *ManagePO) remove_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	
	//	0
	//[data.id]
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	fmt.Println("- start remove trade")
	timestamp, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return nil, errors.New("1st argument must be a numeric string")
	}
	
	//get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)																//un stringify it aka JSON.parse()
	
	for i := range trades.OpenTrades{																	//look for the trade
		//fmt.Println("looking at " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10) + " for " + strconv.FormatInt(timestamp, 10))
		if trades.OpenTrades[i].Timestamp == timestamp{
			fmt.Println("found the trade");
			trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)				//remove this trade
			jsonAsBytes, _ := json.Marshal(trades)
			err = stub.PutState(openTradesStr, jsonAsBytes)												//rewrite open orders
			if err != nil {
				return nil, err
			}
			break
		}
	}
	
	fmt.Println("- end remove trade")
	return nil, nil
}*/

// ============================================================================================================================
// Clean Up Open Trades - make sure open trades are still possible, remove choices that are no longer possible, remove trades that have no valid choices
// ============================================================================================================================
/*func cleanTrades(stub shim.ChaincodeStubInterface)(err error){
	var didWork = false
	fmt.Println("- start clean trades")
	
	//get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)																		//un stringify it aka JSON.parse()
	
	fmt.Println("# trades " + strconv.Itoa(len(trades.OpenTrades)))
	for i:=0; i<len(trades.OpenTrades); {																		//iter over all the known open trades
		fmt.Println(strconv.Itoa(i) + ": looking at trade " + strconv.FormatInt(trades.OpenTrades[i].Timestamp, 10))
		
		fmt.Println("# options " + strconv.Itoa(len(trades.OpenTrades[i].Willing)))
		for x:=0; x<len(trades.OpenTrades[i].Willing); {														//find a marble that is suitable
			fmt.Println("! on next option " + strconv.Itoa(i) + ":" + strconv.Itoa(x))
			_, e := findMarble4Trade(stub, trades.OpenTrades[i].User, trades.OpenTrades[i].Willing[x].Color, trades.OpenTrades[i].Willing[x].Size)
			if(e != nil){
				fmt.Println("! errors with this option, removing option")
				didWork = true
				trades.OpenTrades[i].Willing = append(trades.OpenTrades[i].Willing[:x], trades.OpenTrades[i].Willing[x+1:]...)	//remove this option
				x--;
			}else{
				fmt.Println("! this option is fine")
			}
			
			x++
			fmt.Println("! x:" + strconv.Itoa(x))
			if x >= len(trades.OpenTrades[i].Willing) {														//things might have shifted, recalcuate
				break
			}
		}
		
		if len(trades.OpenTrades[i].Willing) == 0 {
			fmt.Println("! no more options for this trade, removing trade")
			didWork = true
			trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)					//remove this trade
			i--;
		}
		
		i++
		fmt.Println("! i:" + strconv.Itoa(i))
		if i >= len(trades.OpenTrades) {																	//things might have shifted, recalcuate
			break
		}
	}

	if(didWork){
		fmt.Println("! saving open trade changes")
		jsonAsBytes, _ := json.Marshal(trades)
		err = stub.PutState(openTradesStr, jsonAsBytes)														//rewrite open orders
		if err != nil {
			return err
		}
	}else{
		fmt.Println("! all open trades are fine")
	}

	fmt.Println("- end clean trades")
	return nil
}*/
