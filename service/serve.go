package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type RequestMsg struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		RequestMsg struct {
			Request string `xml:",cdata"`
		} `xml:"RequestMsg"`
	} `xml:"Body"`
}
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
	http.HandleFunc("/b2c", handleB2CRequest)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func handleB2CRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received request:\n%s\n", string(body))

	var requestMsg RequestMsg
	err = xml.Unmarshal(body, &requestMsg)
	if err != nil {
		log.Printf("Error parsing XML: %v", err)
		http.Error(w, "Error parsing XML", http.StatusBadRequest)
		return
	}

	var request Request
	err = xml.Unmarshal([]byte(requestMsg.Body.RequestMsg.Request), &request)
	if err != nil {
		log.Printf("Error parsing inner XML: %v", err)
		log.Printf("Inner XML content:\n%s\n", requestMsg.Body.RequestMsg.Request)
		http.Error(w, "Error parsing inner XML", http.StatusBadRequest)
		return
	}

	log.Printf("Received B2C payment request for receiver: %s", request.Identity.ReceiverParty.Identifier)

	response := &Response{
		ResponseCode:             "0",
		ResponseDesc:             "Process service request successfully.",
		OriginatorConversationID: fmt.Sprintf("S_X%s", time.Now().Format("20060102150405")),
		ConversationID:           fmt.Sprintf("AG_%s", time.Now().Format("20060102T150405")),
		ServiceStatus:            "0",
	}

	responseXML, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error generating response: %v", err)
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(responseXML)
}
