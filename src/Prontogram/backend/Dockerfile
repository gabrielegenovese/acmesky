FROM jolielang/jolie:1.11.2-alpine

WORKDIR /app

COPY ./src/* ./

EXPOSE ${PRONTOGRAM_SERVICE_PORT}

ENTRYPOINT [ "jolie", "main.ol" ]
