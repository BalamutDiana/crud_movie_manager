FROM golang:1.19-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o crud_movie_manager cmd/main.go

CMD ["./crud_movie_manager"]