FROM jolielang/jolie:1.11.2-alpine as build

ARG SERVICE_FILENAME=./main.ol
ARG INPUTPORT_FILENAME=ProntogramService
ARG PRONTOGRAM_SERVICE_HOST=localhost:3000

COPY ./backend/src /backend
COPY ./backend/rest_template.json /backend/

WORKDIR /backend

RUN jolie2openapi $SERVICE_FILENAME ${INPUTPORT_FILENAME} ${PRONTOGRAM_SERVICE_HOST} .
RUN mv ./$INPUTPORT_FILENAME.json ./openapi.json

FROM jolielang/leonardo
 
ENV LEONARDO_WWW=/static

COPY ./docs/swagger ${LEONARDO_WWW}
COPY --from=build /backend/openapi.json $LEONARDO_WWW/

WORKDIR /leonardo

EXPOSE 8080
