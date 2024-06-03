# Choreographies

The following choreographies formalise all possible interactions that may occur in the system with ACMESky. Interactions that do not include ACMESky Service as the sender or recipient are ignored, as they are external and their implementation is unknown to the company.

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

$\text{reqFlight}::=(\text{reqFlightInfo} : ACM \to FC_i)\;;\;(\text{resFlightInfo} : FC_i \to ACM)$

where $i$ can be any flight company and the operations are:

- `reqFlightInfo`: request fights for the user;
- `resFlightInfo`: respond fights' information to ACMESky.

### Registering a user's interest

> ACMESky should register if a user is interested in an air route.

$\text{regUser}::=(\\ \; (\text{registerFlightInterest}: USR_j \to ACM) \;; \\ \;((\text{resConfirm}: ACM \to USR_j) \,+ (\text{resError}: ACM \to USR_j))\\ )^*$

> TODO: controllare che avvenga veramente così

where $j$ can be any users who want to register an interest for a flight, then ACMESky can reply with a confirmation or an error. The operations are:

- `registerFlightInterest`: request from te user to register a new interest;
- `resConfirm`: confim response;
- `resError`: error response.

### Receiving and notify last-minute offers

> A flight company notify ACMESky that a last-minute offer is available. ACMESky, throught Prontogram, should notify every user interested. The user can decide to accept or ignore an offer. If a user accept an offer, then ACMESky notify the flight company.

$\text{lastMinOffer}::=(\\
\;( \text{recvOffer}: FC_i \to ACM )\;;\\
\;( \text{notifyPG}: ACM \to PG )\;;\\
\;( \text{notifyUsr}: PG \to USR_j )\;;\\
\;( 1 +\\
\;\;(\\
\;\;\;( \text{acceptOffer} : USR_j \to ACM)\;;\\
\;\;\;( \text{notifyFlightCompany} : ACM \to FC_i)\\
\;\;)\\
\;)\\
)^*$

> TODO: controllare che avvenga veramente così

This choreography is valid for each flight company $i$ and for each user $n$ that want to know about last-minute offers. The operations are:

- `recvOffer`: the flight company place a new offer;
- `notifyPG`: ACMESky use Prontogram API to start the notification process;
- `notifyUsr`: Prontogram notify every user;
- `acceptOffer`: the user insert the code offer in ACMESky;
- `notifyFlightCompany`: the corresponding flight company is nortified.

### Buying a ticket

> The user can buy a ticket, paying throughth the Bank Service. Then ACMESky check if the flight cost exceeds 1000 euros to offer free transfert service if the airport is 30 km away from the accommodation. ACMESky also book the nearest NCC company to the airport.

$\text{buyTicket}::=(\\
\;(\text{wantToBuy}: USR_j \to ACM)\;;\\
\;(\text{requestPayment}: ACM \to BK)\;;\\
\;(\text{resPaymentData}: BK \to ACM)\;;\\
\;(\text{resPaymentData}: ACM \to USR_j)\;;\\
\;(\text{payReceipt}: USR_j \to BK)\;;\\
\;(\\
\;\;(\\
\;\;\;(\text{paymentFailed}: BK \to ACM)\;;\\
\;\;\;(\text{sendError}: ACM \to USR_j)\\
\;\;)\\
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
\;\;\;(\text{sendTicketAndData}: ACM \to USR_j)\\
\;\;)\\
\;)\\
)^*$

> TODO: controllare che avvenga veramente così

where $USR_j$ is any user, $FC_i$ is any flight company and $NCC_k$ is any rental service. The opertations are:

- `wantToBuy`: a user insert a code to buy a ticket
- `requestPayment`: ACMESky requests the generation of a new payment practice
- `resPaymentData`: the bank creates the payment practice and return the informations
- `resPaymentData`: ACMESky returns the information to the user and ask to pay the bill
- `payReceipt`: the user pay with its bank details
- `paymentFailed`: the payment has failed and ACMESky is notified
- `sendError`: ACMESky notify the user of the failed operation
- `paymentOk`: the payment was succesful and the user is notified
- `bookTicket`: ACMESky asks the flight company to book the ticket
- `sendTicketData`:  the flight comapany sends the ticket informations to ACMESky
- `calcGeoDistance`: ACMESky asks the distance of two points to the GeoDistance service
- `resDistance`: the GeoDistance service responds to ACMESky
- `bookTransport`: ACMESky books a cabin for the user
- `resBookTransport`: the rental service responds with details to ACMESky
- `sendTicketAndData`: ACMESky sends the ticket and all the optional informations to the user

## Verifying connectedness

> TODO: controllare

We will analyse the connectedness of the separeted choreographies because there are no codition for parallel composition.

### Requesting flights

This choreography is connected for the sequence because the receiver in `reqFlightInfo` is equal to the sender in `resFlightInfo`. This choreography is connected for the **sender** pattern.

### Registering a user's interest

This choreography is connected because the receiver in `registerFlightInterest` is equal to the sender in both `resConfirm` and `resError`. This choreography is connected for the **sender** pattern.

### Receiving and notify last-minute offers

This choreography is connected because the receiver in `recvOffer` is equal to the sender in `notifyPG`, the receiver in `notifyPG` is equal to the sender in `notifyUsr`, the receiver in `notifyUsr` is equal to the sender in `acceptOffer` and the receiver in `acceptOffer` is equal to the sender in `notifyFlightCompany`. This choreography is connected for the **sender** pattern.

### Buying a ticket

This choreography is connected because:

- the receiver in `wantToBuy` is equal to the sender in `requestPayment`; 
- the receiver in `requestPayment` is equal to the sender in `resPaymentData`; 
- the receiver in `resPaymentData` is equal to the sender in `payReceipt`; 
- the receiver in `payReceipt` is equal in both operation's senders (`paymentFailed` and `paymentOk`); 
- the receiver in `paymentFailed` is equal to the sender in `sendError`; 
- the receiver in `sendError` is equal to the sender in `sendTicketAndData`; 
- the receiver in `paymentOk` is equal to the sender in `bookTicket`; 
- the receiver in `bookTicket` is equal to the sender in `sendTicketData`; 
- the receiver in `sendTicketData` is equal to the sender in `calcGeoDistance`; 
- the receiver in `sendTicketData` is equal to the sender in `sendTicketAndData`; 
- the receiver in `calcGeoDistance` is equal to the sender in `resDistance`; 
- the receiver in `resDistance` is equal to the sender in `bookTransport`; 
- the receiver in `resDistance` is equal to the sender in `sendTicketAndData`; 
- the receiver in `resBookTransport` is equal to the sender in `sendTicketAndData`.

This choreography is connected for the **sender** pattern $(b=c)$, so the entire choreography is connected.

## Projections

Choreographies where the participant taken in consideration is not present will be ignored.

### $ACM$

$\text{proj}(\text{reqFlight}, ACM) = (\text{reqFlightInfo}@FC_i\;;\;\overline{\text{resFlightInfo}}@FC_i)$

$\text{proj}(\text{regUser}, ACM) = (\overline{\text{registerFlightInterest}}@USR_j\;;\;(\text{resConfirm}@USR_j\;+\; \text{resError}@USR_j))$

$\text{proj}(\text{lastMinOffer}, ACM) = ((\text{recvOffer}@FC_i);(\overline{\text{notifyPG}}@PG); 1;( 1 +((\text{acceptOffer}@USR_j);(\overline{\text{notifyFlightCompany}}@FC_i)))) \\
= (\text{recvOffer}@FC_i)\;;\;(\overline{\text{notifyPG}}@PG)\;;\;(\;1\;+\;((\text{acceptOffer}@USR_j)\;;\;(\overline{\text{notifyFlightCompany}}@FC_i)))$

$\text{proj}(\text{buyTicket}, ACM) = (\\
\;(\text{wantToBuy}@USR_j)\;;\\
\;(\overline{\text{requestPayment}@BK})\;;\\
\;(\text{resPaymentData}@BK)\;;\\
\;(\overline{\text{resPaymentData}@USR_j})\;;\\
\;1\;;\\
\;(\\
\;\;(\\
\;\;\;(\text{paymentFailed}@BK)\;;\\
\;\;\;(\overline{\text{sendError}}@USR_j)\\
\;\;)\\
\;\;+\\
\;\;(\\
\;\;\;(\text{paymentOk}@BK)\;;\\
\;\;\;(\overline{\text{bookTicket}@FC_i})\;;\\
\;\;\;(\text{sendTicketData}@FC_i)\;;\\
\;\;\;(\\
\;\;\;\; 1 + \\
\;\;\;\;(\\
\;\;\;\;\;(\overline{\text{calcGeoDistance}@GD})\;;\\
\;\;\;\;\;(\text{resDistance}@GD)\;;\\
\;\;\;\;\;(\\
\;\;\;\;\;\;1 + \\
\;\;\;\;\;\;(\\
\;\;\;\;\;\;\;(\\
\;\;\;\;\;\;\;\;(\overline{\text{calcGeoDistance}@GD})\;;\\
\;\;\;\;\;\;\;\;(\text{resDistance}@GD)\\
\;\;\;\;\;\;\;)^*\;;\\
\;\;\;\;\;\;\;(\overline{\text{bookTransport}@NCC_k})\;;\\
\;\;\;\;\;\;\;(\text{resBookTransport}@NCC_k)\;;\\
\;\;\;\;\;\;)\\
\;\;\;\;\;)\\
\;\;\;\;)\\
\;\;\;)\;;\\
\;\;\;(\overline{\text{sendTicketAndData}@USR_j})\\
\;\;)\\
\;)\\
)$

### $FC_i$

$\text{proj}(\text{reqFlight}, FC_i) = (\overline{\text{reqFlightInfo}}@ACM\;;\;\text{resFlightInfo}@ACM)$

$\text{proj}(\text{lastMinOffer}, FC_i) = \\
\;\;(\;(\overline{\text{recvOffer}}@ACM)\;;\;1\;;\;1\;;\;(\;1\;+\;(\;1\;;\;(\text{notifyFlightCompany}@ACM)\;))) \\
=(\overline{\text{recvOffer}}@ACM)\;;\;(\;1\;+\;(\text{notifyFlightCompany}@ACM)$

$\text{proj}(\text{buyTicket}, FC_i) = (\;1\;;\;1\;;\;1\;;\;1\;;\;1\;;\;\;(\;(1;1)\;+\;(\;1\;;\;\\
\;\;(\text{bookTicket}@ACM)\;;\;(\overline{\text{sendTicketData}}@ACM)\;;\\
\;\;(\;1\;+\;(\;1\;;\;1\;;\;(\;1\;+\;\;(\;(\;1\;;\;1\;)^*;\;1;\;1\;;\;))));1))\\
= (\text{bookTicket}@ACM) \;;\;(\overline{\text{sendTicketData}}@ACM)$

### $PG$

$\text{proj}(\text{lastMinOffer}, PG) = \\
\;\;(\;1\;;\;(\text{notifyPG}@ACM)\;;\;(\overline{\text{notifyUsr}}@USR_j)\;;\;(\;1\;+\;(\;1\;;\;1\;))) \\
=(\text{notifyPG}@ACM)\;;\;(\overline{\text{notifyUsr}}@USR_j)$

### $BK$

$\text{proj}(\text{buyTicket}, BK) = (\;1\;;\\
\;(\;\text{requestPayment}@ACM)\;;\;(\overline{\text{resPaymentData}}@ACM)\;;\;1\;;\\
\;(\text{payReceipt}@USR_j)\;;\;(\;(\overline{\text{paymentFailed}}@USR_j)\;+\;(\;(\overline{\text{paymentOk}}@USR_j)\;;\\
1\;;\;1\;;\;(\;(1;1)\;+\;(\;1\;;\;1\;;(\;1\;+\;\;(\;(\;1\;;\;1\;)^*;1;\;\;1\;))));1))\\
= (\text{requestPayment}@ACM)\; ;\; (\overline{\text{resPaymentData}}@ACM)\; ;\\
\;\;\;(\text{payReceipt}@USR_j)\;; \;(\; (\overline{\text{paymentFailed}}@USR_j)\; +\; (\overline{\text{paymentOk}}@USR_j)\;)$

### $GD$

$\text{proj}(\text{buyTicket}, GD) = (\;1\;;\;1\;;\;1\;;\;1\;;\;1\;;\;(\;(1;1)\;+\;(\;1\;;\;1\;;\;1\;;\;(\;1\;+\;(\;\\
\;\;\;(\text{calcGeoDistance}@ACM)\;;\;(\overline{\text{resDistance}}@ACM)\;;\\
\;\;\;(\;1\;+\;\;(\;(\;(\text{calcGeoDistance}@ACM)\;;\;(\overline{\text{resDistance}}@ACM)\;)^*;\\
1;\;\;1\;))));1))\\
= (\;1\;+\;(\text{calcGeoDistance}@ACM)\;;\;(\overline{\text{resDistance}}@ACM)\;;\; (\;1\;+\;(\text{calcGeoDistance}@ACM)\;;\;(\overline{\text{resDistance}}@ACM)))$

### $USR_j$

$\text{proj}(\text{regUser}, USR_j) = (\text{registerFlightInterest}@ACM\;;\;(\overline{\text{resConfirm}}@ACM\;+\; \overline{\text{resError}}@ACM))$

$\text{proj}(\text{lastMinOffer}, USR_j) = \\
\;\;(\;1\;;\;1\;;\;(\text{notifyUsr}@PG)\;;\;(\;1\;+\;(\;(\overline{\text{acceptOffer}}@ACM)\;;\;1\;))) \\
=(\text{notifyUsr}@PG)\;;\;(\;1\;+\;(\overline{\text{acceptOffer}}@ACM)$

$\text{proj}(\text{buyTicket}, USR_j) = (\\
\;(\overline{\text{wantToBuy}}@ACM)\;;\\
\;1\;;\\
\;1\;;\\
\;(\text{resPaymentData}@ACM)\;;\\
\;(\overline{\text{payReceipt}}@BK)\;;\\
\;(\\
\;\;(\\
\;\;\;1\\\
\;\;\;(\text{sendError}@ACM)\\\
\;\;)\\
\;\;+\\
\;\;(\\
\;\;\;1\;;\\
\;\;\;1\;;\\
\;\;\;1\;;\\
\;\;\;(\\
\;\;\;\; 1 + \\
\;\;\;\;(\\
\;\;\;\;\;1\;;\\
\;\;\;\;\;1\;;\\
\;\;\;\;\;(\\
\;\;\;\;\;\;1 + \\
\;\;\;\;\;\;(\\
\;\;\;\;\;\;\;(\\
\;\;\;\;\;\;\;\;1\;;\\
\;\;\;\;\;\;\;\;1\\
\;\;\;\;\;\;\;)^*\;;\\
\;\;\;\;\;\;\;1\;;\\
\;\;\;\;\;\;\;1\;;\\
\;\;\;\;\;\;)\\
\;\;\;\;\;)\\
\;\;\;\;)\\
\;\;\;)\;;\\
\;\;\;(\text{sendTicketAndData}@ACM)\\
\;\;)\\
\;)\\
)\\
= (\overline{\text{wantToBuy}}@ACM)\;;\; (\text{resPaymentData}@ACM)\;;\;\\
 (\overline{\text{payReceipt}}@BK)\;;\; (\;(\text{sendError}@ACM)\;+\; (\text{sendTicketAndData}@ACM)\;)$

### $NCC_k$

$\text{proj}(\text{buyTicket}, NCC_k) = (\;1\;;\;1\;;\;1\;;\;1\;;\;1\;;\\
\;\;(\;(1;1)\;+\;(\;1\;;\;1\;;\;1\;;\;(\;1\;+\;(\;1\;;\;1\;;\;(\;1\;+\;\;(\;(\;1\;;\;1\;)^*;\\
\;\;\;\;(\text{bookTransport}@ACM);\\
\;\;\;\;(\overline{\text{resBookTransport}}@ACM)\;;\\
\;\;))));1)) = (\text{bookTransport}@ACM)\;;\;(\overline{\text{resBookTransport}}@ACM)$