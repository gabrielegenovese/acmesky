# ACMEsky

Project for the Ingegneria del Software Orientata ai Servizi AA 2023/2024's course.

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

- Camunda
  - [Operate](http://localhost:8081)
  - [Tasklist](http://localhost:8082)
- ACMESky
  - [Frontend](http://localhost:8080)
  - [API](http://localhost:8090)
- FlightCompany
  - [API](http://localhost:8091)
- Prontogram
  - [Frontend](http://localhost:8079)
  - [API](http://localhost:8092)

## Docs

The documentation is generated with [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/).
