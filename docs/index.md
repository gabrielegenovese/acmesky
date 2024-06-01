# ACMESky Documentation

This is the project's documentation of **Ingegneria del Software Orientata ai Servizi** for the AA 2023/2024 University of Bologna's Master Degree. The team members are:

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

    mkdocs.yml    # The configuration file.
    docker-compose.yaml
    docs/
        index.md  # The documentation homepage.
        ...       # Other markdown pages, images and other files.
    resources/
        ...
    src/
        ACMESkyDB/
        ACMESkyNCC/
        ACMESkyService/
        ACMESkyWebSite/
        BankService/
        FlightCompanyDB/
        FlightCompanyService/
        GeoDistance/
        Prontogram/

> TODO: controllare

## How to start

To start all the infrastructure just use the following command in the root directory and wait for a bit!

```shell
docker compose up -d
```
