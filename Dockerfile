FROM golang as builder

WORKDIR /go/src/github.com/BruceLEO1969/ses-exporter
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ses-exporter .

FROM alpine

COPY --from=builder /go/src/github.com/BruceLEO1969/ses-exporter/ses-exporter /bin
RUN apk --update add --no-cache ca-certificates

EXPOSE 9435

CMD ["ses-exporter"]
