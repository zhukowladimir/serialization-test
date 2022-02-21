FROM golang:1.17

WORKDIR /serialization-test

COPY main.go ./
COPY go.* ./
COPY schema.avsc ./
COPY proto_stuff/models/test.pb.go ./proto_stuff/models/
COPY report/report_template.xlsx ./report/

RUN mkdir files
RUN go run main.go
