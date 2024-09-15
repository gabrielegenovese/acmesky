# ACMEsky

Project for the Service Oriented Software Engineering (Ingegneria del Software Orientata ai Servizi) AA 2023/2024's course.

## Deploy using Docker compose

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
    -f docker-compose-NCC.yaml \
    -f docker-compose-Bank.yaml \
    -f docker-compose-workers.yaml \
    up -d
```

## Docs

The documentation can be viewd [here](https://acmesky-expuss2000-d511ae4437280f0ac43140f188bc34fa077888571f00.gitlab.io/) and it's generated with [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/).
