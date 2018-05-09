FROM golang:1.8.3 as builder
LABEL builder=true
COPY src /go/src/
RUN set -x && \
	cd /go/src/dynatrace/cmd/apdex && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/apdex .

FROM scratch
LABEL maintainer="kin.wai.koo@dynatrace.com"
LABEL builder=false
COPY --from=builder /go/bin/apdex /

# we need to copy the certificates over because we're connecting over SSL
COPY --from=builder /etc/ssl /etc/ssl

EXPOSE 8080

ENTRYPOINT ["/apdex"]

