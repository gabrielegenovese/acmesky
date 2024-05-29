# Choreographies

The following choreographies model all possible interactions that may occur in the ACMESky system. Interactions that do not include ACMESky Service as the sender or recipient are ignoerd, as they are external and their implementation is unknown to the company.

## Naming convention

The following abbreviations will be used:

- $ACM$: the ACMESky Service
- $FC_i$: a Fight Company Service, where $i \in \{1...N\}$
- $PG$: the Prontogram Service
- $BK$: the Bank Service
- $GD$: the GeoDistance Service
- $USR_j$: a user, where $j \in \{1...M\}$
- $NCC_k$: the NCC Service

where $N$ and $M$ are natural numbers.

## Choreography

All the following choreographies constitute a single parallel choreography. They are divided for readability.

### Receiving and notify last-minute offers

A flight company notify ACMESky that a last-minute offer is aviable. ACMESky notify Prontogram. Prontogram sends a message to the user. Now, the use can decide to accept or ignore an offer. If a user accept an offer, then ACMESky notify the flight company.

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

This choreography is valid for each flight company $i$ and for each user $n$ that want to know about last-minute offers.

### Registering a user's interest

$(\\ \; (\text{registerFlightInterest }: USR_j \to ACM) \;; \\ \;((\text{resConfirm}: ACM \to USR_j) \,+ (\text{resError}: ACM \to USR_j))\\ )^*$

Any user $j$ can register an interest for a flight, then ACMESky can reply with a confirmation or an error.

### Buying a ticket

The user ask to buy an offer. 

$(\\
\;( \text{wantToBuy}: USR_j \to ACM )\;;\\
\;(\text{requestPayment}: ACM \to BK )\;;\\
\;(\text{resPaymentData}: BK \to ACM )\;;\\
\;(\text{resDataToUser}: ACM \to USR_j )\;;\\
\;(\text{payReceipt}: USR_j \to BK )\;;\\
\;(\\
\;\;(\text{paymentFailed}: BK \to USR_j )\\
\;\;+\\
\;\;(\\
\;\;\;(\text{paymentOk}: BK \to ACM )\;;\\
\;\;\;(\text{bookTicket}: ACM \to FC_i )\;;\\
\;\;\;(\\
\;\;\;\; 1 + \\
\;\;\;\;(\\
\;\;\;\;\;(\text{calcGeoDistance}: ACM \to GD)\;;\\
\;\;\;\;\;(\text{resDistance}: GD \to ACM )\;;\\
\;\;\;\;\;(\\
\;\;\;\;\;\;1 + \\
\;\;\;\;\;\;(\\
\;\;\;\;\;\;\;(\\
\;\;\;\;\;\;\;\;(\text{calcGeoDistance}: ACM \to GD)\;;\\
\;\;\;\;\;\;\;\;(\text{resDistance}: GD \to ACM )\\
\;\;\;\;\;\;\;)^*\;;\\
\;\;\;\;\;\;\;(\text{bookTransport}: ACM \to NCC_k)\;;\\
\;\;\;\;\;\;\;(\text{resBookTransport}: NCC_k \to ACM)\;;\\
\;\;\;\;\;\;)\\
\;\;\;\;\;)\\
\;\;\;\;)\\
\;\;\;)\;;\\
\;\;\;(\text{sendTicketData}: FC_i \to ACM )\;;\\
\;\;\;(\text{sendTicket}: ACM \to USR_j )\;;\\
\;\;)\\
\;)\\
)^*$

## Verifying connectedness

## Projections

### $ACM$

### $FC_i$

### $PG$

### $BK$

### $GD$

### $USR_j$
