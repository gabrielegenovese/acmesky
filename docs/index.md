# ACMESky Documentation

This is the project's documentation of the **Service Oriented Software Engineering**'s course (Ingegneria del Software Orientata ai Servizi) for the AA 2023/2024 University of Bologna's Master Degree. The team members are:

- Luca Bassi
- Gabriele Genovese
- Jacopo Rimediotti

## Description of the specifications

ACMESky is a fictional company that offers multiple services for buying and managing client's return flights. A _Service Oriented Architecture_ (SOA) should be designed and implemented to support ACMESky's activities.

ACMESky offers a service that allows customers to specify, through a web portal, their interest in round-trip air transfers within a defined period and at a cost below a certain set limit.

ACMESky daily queries airlines to obtain quotes for flights of interest to its customers.

ACMESky also receives last-minute offers from airlines, which are sent upon activation without a preset schedule.

When ACMESky finds a flight compatible with a customer's request, it prepares an offer.

The offer is sent to the customer via the messaging app Prontogram. If interested, the customer has 24 hours to connect to the ACMESky web portal to confirm the offer, specifying the code received via Prontogram.

During confirmation, the customer must also proceed with the payment. ACMESky relies on a banking service provider for payment management: ACMESky redirects the customer to the provider's site and then awaits the provider's confirmation message that the payment has been made.

If the flight cost exceeds 1000 euros, ACMESky offers the customer a free transfer service to/from the airport if it is within 30 kilometers of their home.

In this case, ACMESky uses various chauffeur rental companies with which it has commercial agreements. The chosen company is the one that has a location closest to the customer's home. ACMESky sends a request to this company to book a transfer that departs two hours before the scheduled flight departure time.

## Project layout

This is the structure of the project's folder.

    docs/                           # Documentation contents
    src/                            # All microservices codebase
    resources/
        uml/                        # UML source files
        bpmn/                       # BPMN source files
        bpmn-chors/                 # BPMN Choreography source files
    docker-compose-[service].yaml   # All docker compose files
    requirements.txt                # Documentation requirements
    mkdocs.yml                      # Documentation config file

## Deploy using Docker compose

To start all the infrastructure just use the following commands in the root directory and wait for a bit!

If running locally as demo you must run the following docker-compose command first as setup:
This will create a fake shared network which allow all other container to comunicate each to others.

```bash
docker-compose -f docker-compose-shared.yaml up -d
```

Deploy all platforms compose

```bash
docker-compose \
    -f docker-compose-camunda.yaml \
    -f docker-compose-ACMESky.yaml \
    -f docker-compose-FlightCompany.yaml \
    -f docker-compose-Prontogram.yaml \
    -f docker-compose-GeoDistance.yaml \
    -f docker-compose-Bank.yaml \
    up -d
```

## Docs

The documentation is generated with [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/).
