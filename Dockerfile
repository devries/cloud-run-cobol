FROM golang:1.13 as gobuild
ADD invoke.go /src/invoke.go
RUN cd /src && CGO_ENABLED=0 GOOS=linux go build -o invoke invoke.go

FROM ubuntu:18.04
RUN apt-get update -y -q
RUN apt-get install -y -q open-cobol

ADD hello.cobol /app/hello.cobol

RUN cd /app && cobc -x -o hw hello.cobol
COPY --from=gobuild /src/invoke /app/invoke

WORKDIR /app
CMD ["/app/invoke"]
