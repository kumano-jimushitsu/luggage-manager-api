FROM golang:1.17

WORKDIR /api

COPY ./config ./config
COPY ./server ./server

RUN cd server && go mod download

RUN cd server && go build -o . ./main/server.go ./main/routes.go ./main/db.go

EXPOSE 8080

CMD ["./server/server"]