package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// SOAP envelope for request
type RequestMsg struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	SoapEnv string   `xml:"xmlns:soapenv,attr"`
	ReqNS   string   `xml:"xmlns:req,attr"`
	Body    struct {
		RequestMsg struct {
			Request string `xml:",cdata"`
		} `xml:"req:RequestMsg"`
	} `xml:"soapenv:Body"`
}

// Request structure
type Request struct {
	XMLName  xml.Name `xml:"Request"`
	KeyOwner string   `xml:"KeyOwner"`
	Identity struct {
		Caller struct {
			CallerType   string `xml:"CallerType"`
			ThirdPartyID string `xml:"ThirdPartyID"`
			Password     string `xml:"Password"`
			ResultURL    string `xml:"ResultURL"`
		} `xml:"Caller"`
		Initiator struct {
			IdentifierType     string `xml:"IdentifierType"`
			Identifier         string `xml:"Identifier"`
			SecurityCredential string `xml:"SecurityCredential"`
			ShortCode          string `xml:"ShortCode"`
		} `xml:"Initiator"`
		ReceiverParty struct {
			IdentifierType string `xml:"IdentifierType"`
			Identifier     string `xml:"Identifier"`
		} `xml:"ReceiverParty"`
	} `xml:"Identity"`
	Transaction struct {
		CommandID  string `xml:"CommandID"`
		Timestamp  string `xml:"Timestamp"`
		Parameters struct {
			Parameter []struct {
				Key   string `xml:"Key"`
				Value string `xml:"Value"`
			} `xml:"Parameter"`
		} `xml:"Parameters"`
		ReferenceData struct {
			ReferenceItem struct {
				Key   string `xml:"Key"`
				Value string `xml:"Value"`
			} `xml:"ReferenceItem"`
		} `xml:"ReferenceData"`
	} `xml:"Transaction"`
}

// Response represents the response structure
type Response struct {
	XMLName                  xml.Name `xml:"Response"`
	ResponseCode             string   `xml:"ResponseCode"`
	ResponseDesc             string   `xml:"ResponseDesc"`
	OriginatorConversationID string   `xml:"OriginatorConversationID"`
	ConversationID           string   `xml:"ConversationID"`
	ServiceStatus            string   `xml:"ServiceStatus"`
}

func main() {
	request := createRequest()
	response, err := sendRequest("http://localhost:8080/b2c", request)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	fmt.Printf("Response: %+v\n", response)
}

func createRequest() *RequestMsg {
	request := &RequestMsg{
		SoapEnv: "http://schemas.xmlsoap.org/soap/envelope/",
		ReqNS:   "http://cps.huawei.com/cpsinterface/request",
	}

	innerRequest := &Request{
		KeyOwner: "1",
		Identity: struct {
			Caller struct {
				CallerType   string `xml:"CallerType"`
				ThirdPartyID string `xml:"ThirdPartyID"`
				Password     string `xml:"Password"`
				ResultURL    string `xml:"ResultURL"`
			} `xml:"Caller"`
			Initiator struct {
				IdentifierType     string `xml:"IdentifierType"`
				Identifier         string `xml:"Identifier"`
				SecurityCredential string `xml:"SecurityCredential"`
				ShortCode          string `xml:"ShortCode"`
			} `xml:"Initiator"`
			ReceiverParty struct {
				IdentifierType string `xml:"IdentifierType"`
				Identifier     string `xml:"Identifier"`
			} `xml:"ReceiverParty"`
		}{
			Caller: struct {
				CallerType   string `xml:"CallerType"`
				ThirdPartyID string `xml:"ThirdPartyID"`
				Password     string `xml:"Password"`
				ResultURL    string `xml:"ResultURL"`
			}{
				CallerType:   "2",
				ThirdPartyID: "POS_Broker",
				Password:     "B1YNY8GylVo=",
				ResultURL:    "http://10.71.109.150:8888/mockResultBinding",
			},
			Initiator: struct {
				IdentifierType     string `xml:"IdentifierType"`
				Identifier         string `xml:"Identifier"`
				SecurityCredential string `xml:"SecurityCredential"`
				ShortCode          string `xml:"ShortCode"`
			}{
				IdentifierType:     "11",
				Identifier:         "aaa",
				SecurityCredential: "iUoiP9iGVwE=",
				ShortCode:          "1008",
			},
			ReceiverParty: struct {
				IdentifierType string `xml:"IdentifierType"`
				Identifier     string `xml:"Identifier"`
			}{
				IdentifierType: "1",
				Identifier:     "8613500000000",
			},
		},
		Transaction: struct {
			CommandID  string `xml:"CommandID"`
			Timestamp  string `xml:"Timestamp"`
			Parameters struct {
				Parameter []struct {
					Key   string `xml:"Key"`
					Value string `xml:"Value"`
				} `xml:"Parameter"`
			} `xml:"Parameters"`
			ReferenceData struct {
				ReferenceItem struct {
					Key   string `xml:"Key"`
					Value string `xml:"Value"`
				} `xml:"ReferenceItem"`
			} `xml:"ReferenceData"`
		}{
			CommandID: "InitTrans_2304",
			Timestamp: time.Now().Format("20060102150405"),
			Parameters: struct {
				Parameter []struct {
					Key   string `xml:"Key"`
					Value string `xml:"Value"`
				} `xml:"Parameter"`
			}{
				Parameter: []struct {
					Key   string `xml:"Key"`
					Value string `xml:"Value"`
				}{
					{Key: "Amount", Value: "10.00"},
					{Key: "Currency", Value: "RON"},
					{Key: "ReasonType", Value: "Pay for Individual B2C_VDF_Demo"},
					{Key: "Remark", Value: "For buy good refund, the original transaction number can be filled in remark."},
				},
			},
			ReferenceData: struct {
				ReferenceItem struct {
					Key   string `xml:"Key"`
					Value string `xml:"Value"`
				} `xml:"ReferenceItem"`
			}{
				ReferenceItem: struct {
					Key   string `xml:"Key"`
					Value string `xml:"Value"`
				}{
					Key:   "POSDeviceID",
					Value: "POS234789",
				},
			},
		},
	}

	innerRequestXML, err := xml.MarshalIndent(innerRequest, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling inner request: %v", err)
	}

	request.Body.RequestMsg.Request = string(innerRequestXML)

	return request
}

func sendRequest(url string, request *RequestMsg) (*Response, error) {
	requestXML, err := xml.MarshalIndent(request, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %v", err)
	}

	resp, err := http.Post(url, "application/xml", bytes.NewBuffer(requestXML))
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	fmt.Println("Response XML:", string(body))

	var response Response
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return &response, nil
}
