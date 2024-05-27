FROM jolielang/jolie:1.11.2-alpine

WORKDIR /backend

EXPOSE 8000

EXPOSE 8001

ENV NCC_BACKEND_API=http://ncc_backend

COPY . .

RUN ./generateWSDL.sh

CMD [ "jolie", "main.ol" ]
