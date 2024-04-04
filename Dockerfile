FROM golang:1.21 AS BUILDER

WORKDIR $GOPATH/src/smallest-golang/app/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

FROM gcr.io/distroless/static-debian11

COPY --from=BUILDER /main .

ENTRYPOINT ["/main"]