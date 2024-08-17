# UML

> TODO: controlla e modifica

In this section, the SOA modeling for the ACMESky organization is presented using UML diagrams with the TinySOA profile. These diagrams aim to highlight the capabilities accessible through the system and the interfaces that expose them, both externally and internally, for each service within the SOA. Specifically, there are three types of services:

- **Task (or Process)**: Exposes capabilities achieved through internal processes, potentially involving human participation. These are closely related to the domain problem.
- **Entity**: Represents a single activity, usually automated (e.g., saving a record in a database).
- **Utility**: Similar to Tasks but not directly tied to the domain problem.

## Requesting flights

![uml req flights](assets/uml/req-flights.png "req flights")

For ACMESky, the identified capabilities are FlightAvailability and GetFlightAvailabilityOperation. These capabilities are exposed through two interfaces, enabling the system to capture and provide information on available flights based on the user's request.

Additionally, the BPMN diagram identified another capability for the user, named SetFlightRequestOperation. This capability allows the user to receive feedback on the success or failure of their flight availability request.

## Registering a user's interest

![uml reg interest](assets/uml/reg-inter.png "reg interest")

For ACMESky, the identified capabilities are Interest and GetInterestOperation. These capabilities are exposed through two interfaces, enabling the system to capture a user's interest in a trip and store it for future reference.

Additionally, the BPMN diagram revealed another capability for the user, called SetInterestOperation. This capability allows the user to receive feedback on whether their interest registration was successful or not.

## Receiving and notify last-minute offers

![uml last min](assets/uml/last-min.png "last min")

For ACMESky, the identified capabilities are LastMinuteOffers and GetLastMinuteOffersOperation. These capabilities are exposed through two interfaces, allowing the system to capture and deliver last-minute offers from flight companies to users.

Additionally, the BPMN diagram identified another capability for the user, named SetOfferSubscriptionOperation. This capability enables the user to receive feedback on the success or failure of their subscription to last-minute offers.

## Buying a ticket

![uml buy ticket](assets/uml/buy-tickets.png "buy ticket")

> TODO: description
