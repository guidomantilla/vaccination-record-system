FROM golang:1.21.1-alpine AS build

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /workspace
COPY . .
RUN go mod download && go build -a -o main.bin .

FROM golang:1.21.1-alpine

WORKDIR /
COPY --from=build /workspace/main.bin /main.bin
COPY --from=build /workspace/resources/ /resources/
EXPOSE 8080
CMD ["/main.bin", "serve"]
