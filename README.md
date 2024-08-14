## SOAP API Money Transfer Demo
This is a simple demo of whats i thought is described in the B2C document.
hypothetical scenario like withdrawing money from a betting account through an agent. the withdrawee shows the successfull response (like a voucher) and the agent gives him money.

1) **client.go**: Contains the client-side code that constructs a SOAP request, sends it to the server, and handles the response.
   
3) **server.go**: Contains the server-side code that listens for incoming SOAP requests and unmarshall the request for processing. But no processing is done as its done to simply demonstarte whats B2C_payment_api.doc described. it just simply responds with the result.

further improvements are to be done.
