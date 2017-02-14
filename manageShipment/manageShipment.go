/*/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license shipments.  See the NOTICE file
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

// ManageShipment example simple Chaincode implementation
type ManageShipment struct {
}

var ShipmentIndexStr = "_Shipmentindex"				//name for the key/value that will store a list of all known Shipment

type Shipment struct{							// Attributes of a Shipment 
	ShipmentID string `json:"shipmentId"`	
	TransID string `json:"transId"`
	AgreementID string `json:"agreementId"`
	Shipment_status string `json:"shipment_status"`
	Source string `json:"source"`					
	Destination string `json:"destination"`
	ActualDelivery_date string `json:"actualDelivery_date"`
	Shipment_date string `json:"shipment_date"`
	ShipperName string `json:"shipper_name"`
	
}
// ============================================================================================================================
// Main - start the chaincode for Shipment management
// ============================================================================================================================
func main() {			
	err := shim.Start(new(ManageShipment))
	if err != nil {
		fmt.Printf("Error starting Shipment management chaincode: %s", err)
	}
}
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *ManageShipment) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
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
	fmt.Println("ManageShipment chaincode is deployed successfully.");
	
	// Write the state to the ledger
	err = stub.PutState("abc", []byte(msg))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(ShipmentIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	tosend := "{ \"message\" : \"ManageShipment chaincode is deployed successfully.\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// Run - Our entry shipmentint for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *ManageShipment) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}
// ============================================================================================================================
// Invoke - Our entry shipmentint for Invocations
// ============================================================================================================================
func (t *ManageShipment) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "create_shipment" {											//create a new Shipment
		return t.create_shipment(stub, args)
	}else if function == "delete_shipment" {									// delete an Shipment
		return t.delete_shipment(stub, args)
	}else if function == "update_shipment" {									//update an Shipment
		return t.update_shipment(stub, args)
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
// Query - Our entry shipmentint for Queries
// ============================================================================================================================
func (t *ManageShipment) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getShipment_byID" {													//Read a Shipment by ShipmentID
		return t.getShipment_byID(stub, args)
	} else if function == "getShipment_byStatus" {													//Read a Shipment by Shipper
		return t.getShipment_byStatus(stub, args)
	} else if function == "get_AllShipment" {													//Read all Shipments
		return t.get_AllShipment(stub, args)
	}else if function == "getShipment_byShipper" {													//Read a Shipment by Shipper
		return t.getShipment_byShipper(stub, args)
	}

	fmt.Println("query did not find func: " + function)						//error
	errMsg := "{ \"message\" : \"Received unknown function"+ function + "\", \"code\" : \"503\"}"
	err := stub.SetEvent("errEvent", []byte(errMsg))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// getShipment_byID - get Shipment details for a specific ShipmentID from chaincode state
// ============================================================================================================================
func (t *ManageShipment) getShipment_byID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var shipmentId string
	var err error
	fmt.Println("Fetching Shipment by shipmentID")
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"ShipmentID\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set shipmentId
	shipmentId = args[0]
	valAsbytes, err := stub.GetState(shipmentId)									//get the shipmentId from chaincode state
	if err != nil {
		errMsg := "{ \"message\" : \""+ shipmentId + " not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Print("valAsbytes : ")
	fmt.Println(valAsbytes)
	fmt.Println("Fetched Shipment by shipmentID")
	return valAsbytes, nil													//send it onward
}
// ============================================================================================================================
//  getShipment_byShipper - get Shipment details by Shipper from chaincode state
// ============================================================================================================================
func (t *ManageShipment) getShipment_byShipper(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, shipper_name, errResp string
	var shipmentIndex []string
	var valIndex Shipment
	fmt.Println("Fetching Shipment by Shipper")
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
	fmt.Println("shipper_name : " + shipper_name)
	shipmentAsBytes, err := stub.GetState(ShipmentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Shipment index string")
	}
	fmt.Print("shipmentAsBytes : ")
	fmt.Println(shipmentAsBytes)
	json.Unmarshal(shipmentAsBytes, &shipmentIndex)								//un stringify it aka JSON.parse()
	fmt.Print("shipmentIndex : ")
	fmt.Println(shipmentIndex)
	fmt.Println("len(shipmentIndex) : ")
	fmt.Println(len(shipmentIndex))
	jsonResp = "{"
	for i,val := range shipmentIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for getShipment_byShipper")
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
			if i < len(shipmentIndex)-1 {
				jsonResp = jsonResp + ","
			}
		} else{
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
	fmt.Println("Fetched Shipment by Shipper")
	return []byte(jsonResp), nil											//send it onward
}

// ============================================================================================================================
//  getShipment_byStatus - get Shipment details for a specific Seller from chaincode state
// ============================================================================================================================
func (t *ManageShipment) getShipment_byStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, shipment_status, errResp string
	var shipmentIndex []string
	var valIndex Shipment
	fmt.Println("Fetching Shipment by Shipment Status")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"Shipment_status\" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set shipment_status
	shipment_status = args[0]
	fmt.Println("shipment_status: " + shipment_status)
	shipmentAsBytes, err := stub.GetState(ShipmentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Shipment index")
	}
	fmt.Print("shipmentAsBytes : ")
	fmt.Println(shipmentAsBytes)
	json.Unmarshal(shipmentAsBytes, &shipmentIndex)								//un stringify it aka JSON.parse()
	fmt.Print("shipmentIndex : ")
	fmt.Println(shipmentIndex)
	fmt.Println("len(shipmentIndex) : ")
	fmt.Println(len(shipmentIndex))
	jsonResp = "{"
	for i,val := range shipmentIndex{
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
		if valIndex.Shipment_status == shipment_status{
			fmt.Println("Shipment found")
			jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
			fmt.Println("jsonResp inside if")
			fmt.Println(jsonResp)
			if i < len(shipmentIndex)-1 {
				jsonResp = jsonResp + ","
			}
		}else{
			errMsg := "{ \"message\" : \""+ shipment_status+ " Not Found.\", \"code\" : \"503\"}"
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
	fmt.Println("Fetched Shipment by Shipment Status")
	return []byte(jsonResp), nil											//send it onward
}
// ============================================================================================================================
//  get_AllShipment- get details of all Shipment from chaincode state
// ============================================================================================================================
func (t *ManageShipment) get_AllShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp, errResp string
	var shipmentIndex []string
	fmt.Println("Fetching All Shipments")
	var err error
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \" \" as an argument\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	shipmentAsBytes, err := stub.GetState(ShipmentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Shipment index")
	}
	fmt.Print("shipmentAsBytes : ")
	fmt.Println(shipmentAsBytes)
	json.Unmarshal(shipmentAsBytes, &shipmentIndex)								//un stringify it aka JSON.parse()
	fmt.Print("shipmentIndex : ")
	fmt.Println(shipmentIndex)
	jsonResp = "{"
	for i,val := range shipmentIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for all Shipment")
		valueAsBytes, err := stub.GetState(val)
		if err != nil {
			errResp = "{\"Error\":\"Failed to get state for " + val + "\"}"
			return nil, errors.New(errResp)
		}
		fmt.Print("valueAsBytes : ")
		fmt.Println(valueAsBytes)
		jsonResp = jsonResp + "\""+ val + "\":" + string(valueAsBytes[:])
		if i < len(shipmentIndex)-1 {
			jsonResp = jsonResp + ","
		}
	}
	fmt.Println("len(shipmentIndex) : ")
	fmt.Println(len(shipmentIndex))
	jsonResp = jsonResp + "}"
	fmt.Println("jsonResp : " + jsonResp)
	fmt.Print("jsonResp in bytes : ")
	fmt.Println([]byte(jsonResp))
	fmt.Println("Fetched All Shipments")
	return []byte(jsonResp), nil
											//send it onward
}
// ============================================================================================================================
// delete_shipment - remove a Shipment from chain
// ============================================================================================================================
func (t *ManageShipment) delete_shipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting \"shipmentID\" as an argument.\", \"code\" : \"503\"}"
		err := stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set shipmentId
	shipmentId := args[0]
	err := stub.DelState(shipmentId)													//remove the Shipment from chaincode
	if err != nil {
		errMsg := "{ \"message\" : \"Failed to delete state\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}

	//get the Shipment index
	shipmentAsBytes, err := stub.GetState(ShipmentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Shipment index")
	}
	fmt.Println("shipmentAsBytes in delete shipment")
	fmt.Println(shipmentAsBytes);
	var shipmentIndex []string
	json.Unmarshal(shipmentAsBytes, &shipmentIndex)								//un stringify it aka JSON.parse()
	fmt.Println("shipmentIndex in delete shipment")
	fmt.Println(shipmentIndex);
	//remove shipment from index
	for i,val := range shipmentIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + shipmentId)
		if val == shipmentId{															//find the correct Shipment
			fmt.Println("Found Shipment with matching shipmentId")
			shipmentIndex = append(shipmentIndex[:i], shipmentIndex[i+1:]...)			//remove it
			for x:= range shipmentIndex{											//debug prints...
				fmt.Println(string(x) + " - " + shipmentIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(shipmentIndex)									//save new index
	err = stub.PutState(ShipmentIndexStr, jsonAsBytes)
	tosend := "{ \"shipmentID\" : \""+shipmentId+"\", \"message\" : \"Shipment deleted succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	return nil, nil
}
// ============================================================================================================================
// update_shipment - update Shipment into chaincode state
// ============================================================================================================================
func (t *ManageShipment) update_shipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string
	var err error
	fmt.Println("Updating Shipment")
	if len(args) != 9{
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 9 arguments.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	// set shipmentId
	shipmentId := args[0]
	shipmentAsBytes, err := stub.GetState(shipmentId)									//get the Shipment for the specified shipmentId from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + shipmentId + "\"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Print("shipmentAsBytes in update shipment")
	fmt.Println(shipmentAsBytes);
	res := Shipment{}
	json.Unmarshal(shipmentAsBytes, &res)
	if res.ShipmentID == shipmentId{
		fmt.Println("Shipment found with shipmentId : " + shipmentId)
		fmt.Println(res);

		res.TransID = args[1]
		res.AgreementID = args[2]
		res.Shipment_status = args[3]
		res.Source = args[4]
		res.Destination = args[5]
		res.ActualDelivery_date = args[6]
		res.Shipment_date = args[7]
		res.ShipperName	= args[8]
		
	}else{
		errMsg := "{ \"message\" : \""+ shipmentId+ " Not Found.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	
	//build the Shipment json string manually
	input := 	`{`+
		`"shipmentId": "` + res.ShipmentID + `" , `+
		`"transId": "` + res.TransID + `" , `+ 
		`"agreementId": "` + res.AgreementID + `" , `+ 
		`"shipment_status": "` + res.Shipment_status + `" , `+ 
		`"source": "` + res.Source + `" , `+
		`"destination": "` + res.Destination + `" , `+
		`"actualDelivery_date": "` + res.ActualDelivery_date + `" , `+ 
		`"shipment_date": "` + res.Shipment_date + `" , `+ 
		`"shipper_name": "` + res.ShipperName + `" `+ 
		`}`
	err = stub.PutState(shipmentId, []byte(input))									//store Shipment with id as key
	if err != nil {
		return nil, err
	}
	tosend := "{ \"shipmentID\" : \""+shipmentId+"\", \"message\" : \"Shipment updated succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	} 
	fmt.Println("Updated Shipment succcessfully.")
	return nil, nil
}
// ============================================================================================================================
// create Shipment - create a new Shipment, store into chaincode state
// ============================================================================================================================
func (t *ManageShipment) create_shipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 9{
		errMsg := "{ \"message\" : \"Incorrect number of arguments. Expecting 9 arguments.\", \"code\" : \"503\"}"
		err = stub.SetEvent("errEvent", []byte(errMsg))
		if err != nil {
			return nil, err
		} 
		return nil, nil
	}
	fmt.Println("Creating Shipment")
		
		shipmentId := args[0]
		transId := args[1]
		agreementId := args[2]
		shipment_status := args[3]
		source := args[4]
		destination := args[5]
		actualDelivery_date := args[6]
		shipment_date := args[7]
		shipper_name	:= args[8]
		
		shipmentAsBytes, err := stub.GetState(shipmentId)
		if err != nil {
			return nil, errors.New("Failed to get Shipment ID")
		}
		fmt.Print("shipmentAsBytes: ")
		fmt.Println(shipmentAsBytes)
		res := Shipment{}
		json.Unmarshal(shipmentAsBytes, &res)
		fmt.Print("res: ")
		fmt.Println(res)
		if res.ShipmentID == shipmentId{
			fmt.Println("This Shipment already exists: " + shipmentId)
			errMsg := "{ \"message\" : \""+ shipmentId+" already exists.\", \"code\" : \"503\"}"
			err := stub.SetEvent("errEvent", []byte(errMsg))
			if err != nil {
				return nil, err
			} 
		return nil, nil				//all stop a Shipment by this name exists
	}
	
	//build the Shipment json string manually
	input := 	`{`+
		`"shipmentId": "` + shipmentId + `" , `+
		`"transId": "` + transId + `" , `+ 
		`"agreementId": "` + agreementId + `" , `+ 
		`"shipment_status": "` + shipment_status + `" , `+ 
		`"source": "` + source + `" , `+
		`"destination": "` + destination + `" , `+
		`"actualDelivery_date": "` + actualDelivery_date + `" , `+ 
		`"shipment_date": "` + shipment_date + `" , `+ 
		`"shipper_name": "` + shipper_name + `" `+ 
		`}`
		fmt.Println("input: " + input)
		fmt.Print("input in bytes array: ")
		fmt.Println([]byte(input))
	err = stub.PutState(shipmentId, []byte(input))									//store Shipment with shipmentId as key
	if err != nil {
		return nil, err
	}
	//get the Shipment index
	shipmentIndexAsBytes, err := stub.GetState(ShipmentIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get Shipment index")
	}
	var shipmentIndex []string
	fmt.Print("shipmentIndexAsBytes: ")
	fmt.Println(shipmentIndexAsBytes)
	
	json.Unmarshal(shipmentIndexAsBytes, &shipmentIndex)							//un stringify it aka JSON.parse()
	fmt.Print("shipmentIndex after unmarshal..before append: ")
	fmt.Println(shipmentIndex)
	//append
	shipmentIndex = append(shipmentIndex, shipmentId)									//add Shipment transID to index list
	fmt.Println("! Shipment index after appending shipmentId: ", shipmentIndex)
	jsonAsBytes, _ := json.Marshal(shipmentIndex)
	fmt.Print("jsonAsBytes: ")
	fmt.Println(jsonAsBytes)
	err = stub.PutState(ShipmentIndexStr, jsonAsBytes)						//store name of Shipment
	if err != nil {
		return nil, err
	}
	tosend := "{ \"shipmentID\" : \""+shipmentId+"\", \"message\" : \"Shipment created succcessfully\", \"code\" : \"200\"}"
	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
	}
	
	fmt.Println("Shipment created succcessfully.")
	return nil, nil
}
