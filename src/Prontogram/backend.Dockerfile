FROM jolielang/jolie:1.11.2-alpine

COPY ./backend/src /backend

EXPOSE ${PRONTOGRAM_SERVICE_PORT}

WORKDIR /backend

ENTRYPOINT [ "jolie", "main.ol" ]