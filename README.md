# 1.1.1 Individual B2C payment

The Individual B2C Payment feature allows an organization to pay for individual customers one by one. Typical applications include paying salaries to many individual customers, distributing government relief, allowances, and subsidies, providing rewards or bonuses as merchants for promotion, and paying interest as banks for micro-saving. This feature is a supplement to the Bulk B2C Payment feature.

## 1.1.1.1 Initiate an Individual B2C Payment Transaction (Organization Operator or SP Operator)
The Mobile Money system supports the Individual B2C Payment service that allows an organization to pay for individual customers. The recipient can be a registered or unregistered customer. A withdrawal voucher will be sent if the recipient is an unregistered customer. The recipient (unregistered customer) can withdraw money at an agent using the voucher. This interface is invoked by a third party of the Mobile Money system to initiate an individual B2C payment transaction.

### Request Parameter Specification

|field | Format |
|------|--------|
|Command ID  |	Follow InitTrans_{ServiceCode} format|
|Initiator	| Operator ID of an Organization Operator/PIN/ShortCode|
|            |  User Name of an Organization Operator/Password/ShortCode |
|            |  User Name of an SP Operator/Password|
|PrimaryParty|	Organization ShortCode. **NOTE** It is mandatory when the initiator is an SP operator. Otherwise the value remains empty.|
|ReceiverParty|	Customer MSISDN|

### Specific parameters in Request/TransactionRequest/Parameters
|Parameter	|Data Type|	Mandatory or Optional|	Description	|Example|
|-----------|---------|----------------------|--------------|--------|
|Amount|	xs:string	|Mandatory|	Indicates the payment amount. The precision ranges from the minimum currency unit to 10^9 * minimum currency unit. For example, for CNY, the value of this parameter ranges from 1 to 1000000000.00.	|10.00|
|Currency|	xs:string|	Mandatory|	Indicates the currency type. It is specified by using the Char 3 Code defined in the ISO 4217 standard.|	CNY|
|ReasonType	|xs:string|	Mandatory|	Indicates the reason type for a service or transaction. This field can be set to the reason type alias. Rules: If this field is not configured, the system queries the reason type to be applied to the current transaction based on configurations made by the system operator. If a request specifies a reason type, the reason type in the request is directly used. However, the reason type must belong to the requested transaction type. This field is mandatory if a reason type must be manually specified as required by the requested service.|	Pay for Individual B2C_VDF_Demo. **NOTE** The values of this parameter is customizable.|
