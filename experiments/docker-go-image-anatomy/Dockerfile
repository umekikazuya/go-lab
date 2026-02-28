FROM golang:1.26

WORKDIR /work
COPY inspect.sh .
COPY go.mod .
COPY inspect_test.go .

CMD ["bash", "/work/inspect.sh"]
