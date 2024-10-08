# UML

In this section, the SOA modeling for the ACMESky organization is presented using UML diagrams with the TinySOA profile. These diagrams aim to highlight the capabilities accessible through the system and the interfaces that expose them, both externally and internally, for each service within the SOA. Specifically, there are three types of services:

- **Task (or Process)**: Exposes capabilities achieved through internal processes, potentially involving human participation. These are closely related to the domain problem.
- **Entity**: Represents a single activity, usually automated (e.g., saving a record in a database).
- **Utility**: Similar to Tasks but not directly tied to the domain problem.

## Requesting flights

![uml req flights](assets/uml/req-flights.png "req flights")

## Registering a user's interest

![uml reg interest](assets/uml/reg-inter.png "reg interest")

## Receiving and notify last-minute offers

![uml last min](assets/uml/last-min.png "last min")

## Buying a ticket

![uml buy ticket](assets/uml/buy-tickets.png "buy ticket")
