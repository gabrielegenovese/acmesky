# Choreographies

The following choreographies formalise all possible interactions that may occur in the system with ACMESky. Interactions that do not include ACMESky Service as the sender or recipient are ignoerd, as they are external and their implementation is unknown to the company.

## Naming convention

The following abbreviations will be used:

- $ACM$: the ACMESky Service
- $FC_i$: a Fight Company Service, where $i \in \{1...N\}$
- $PG$: the Prontogram Service
- $BK$: the Bank Service
- $GD$: the GeoDistance Service
- $USR_j$: a user, where $j \in \{1...M\}$
- $NCC_k$: the NCC Service, where $k \in \{1...P\}$

where $N$, $M$ and $P$ are natural numbers.

## Choreography

All the following choreographies constitute a single parallel choreography. They are divided for readability.

### Requesting flights

> ACMESky should ask to every flight company what are the available flights.

$(\text{reqFlightInfo} : ACM \to FC_i)\;;\;(\text{resFlightInfo} : FC_i \to ACM)$

where $i$ can be any flight company and the operations are:

- _reqFlightInfo_: request fights for the user;
- _resFlightInfo_: respond fights' information to ACMESky.

### Registering a user's interest

> ACMESky should register if a user is interested in an air route.

$(\\ \; (\text{registerFlightInterest}: USR_j \to ACM) \;; \\ \;((\text{resConfirm}: ACM \to USR_j) \,+ (\text{resError}: ACM \to USR_j))\\ )^*$

> TODO: controllare che avvenga veramente così

where $j$ can be any users who want to register an interest for a flight, then ACMESky can reply with a confirmation or an error. The operations are:

- _registerFlightInterest_: request from te user to register a new interest;
- _resConfirm_: confim response;
- _resError_: error response.

### Receiving and notify last-minute offers

> A flight company notify ACMESky that a last-minute offer is available. ACMESky, throught Prontogram, should notify every user interested. The user can decide to accept or ignore an offer. If a user accept an offer, then ACMESky notify the flight company.

$(\\
\;( \text{recvOffer}: FC_i \to ACM )\;;\\
\;( \text{notifyPG}: ACM \to PG )\;;\\
\;( \text{notifyUsr}: PG \to USR_j )\;;\\
\;( 1 +\\
\;\;(\\
\;\;\;( \text{acceptOffer} : USR_n \to ACM)\;;\\
\;\;\;( \text{notifyFlightCompany} : ACM \to FC_i)\\
\;\;)\\
\;)\\
)^*$

> TODO: controllare che avvenga veramente così

This choreography is valid for each flight company $i$ and for each user $n$ that want to know about last-minute offers. The operations are:

- _recvOffer_: the flight company place a new offer;
- _notifyPG_: ACMESky use Prontogram API to start the notification process;
- _notifyUsr_: Prontogram notify every user;
- _acceptOffer_: the user insert the code offer in ACMESky;
- _notifyFlightCompany_: the corresponding flight company is nortified.

### Buying a ticket

> The user can buy a ticket, paying throughth the Bank Service. Then ACMESky check if the flight cost exceeds 1000 euros to offer free transfert service if the airport is 30 km away from the accommodation.

$(\\
\;(\text{wantToBuy}: USR_j \to ACM)\;;\\
\;(\text{requestPayment}: ACM \to BK)\;;\\
\;(\text{resPaymentData}: BK \to ACM)\;;\\
\;(\text{resDataToUser}: ACM \to USR_j)\;;\\
\;(\text{payReceipt}: USR_j \to BK)\;;\\
\;(\\
\;\;(\text{paymentFailed}: BK \to USR_j)\\
\;\;+\\
\;\;(\\
\;\;\;(\text{paymentOk}: BK \to ACM)\;;\\
\;\;\;(\text{bookTicket}: ACM \to FC_i)\;;\\
\;\;\;(\text{sendTicketData}: FC_i \to ACM)\;;\\
\;\;\;(\\
\;\;\;\; 1 + \\
\;\;\;\;(\\
\;\;\;\;\;(\text{calcGeoDistance}: ACM \to GD)\;;\\
\;\;\;\;\;(\text{resDistance}: GD \to ACM)\;;\\
\;\;\;\;\;(\\
\;\;\;\;\;\;1 + \\
\;\;\;\;\;\;(\\
\;\;\;\;\;\;\;(\\
\;\;\;\;\;\;\;\;(\text{calcGeoDistance}: ACM \to GD)\;;\\
\;\;\;\;\;\;\;\;(\text{resDistance}: GD \to ACM)\\
\;\;\;\;\;\;\;)^*\;;\\
\;\;\;\;\;\;\;(\text{bookTransport}: ACM \to NCC_k)\;;\\
\;\;\;\;\;\;\;(\text{resBookTransport}: NCC_k \to ACM)\;;\\
\;\;\;\;\;\;)\\
\;\;\;\;\;)\\
\;\;\;\;)\\
\;\;\;)\;;\\
\;\;\;(\text{sendTicketAndData}: ACM \to USR_j)\;;\\
\;\;)\\
\;)\\
)^*$

where $USR_j$ is any user, $FC_i$ is any flight company and $NCC_k$ is any rental service. The opertations are:

- _wantToBuy_:
- _requestPayment_:
- _resPaymentData_:
- _resDataToUser_:
- _payReceipt_:
- _paymentFailed_:
- _paymentOk_:
- _bookTicket_:
- _sendTicketData_:
- _calcGeoDistance_:
- _resDistance_:
- _bookTransport_:
- _resBookTransport_:
- _sendTicketAndData_:

## Verifying connectedness

### Requesting flights

### Registering a user's interest

### Receiving and notify last-minute offers

### Buying a ticket

## Projections

### $ACM$

### $FC_i$

### $PG$

### $BK$

### $GD$

### $USR_j$
