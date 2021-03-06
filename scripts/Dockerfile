FROM golang:alpine AS builder

RUN apk add build-base

ENV GOPATH=/gopath
RUN mkdir -p /gopath/src/github.com/backd-io/backd

RUN mkdir -p /app

WORKDIR /app

COPY .. /gopath/src/github.com/backd-io/backd/

#RUN cd /gopath/src/github.com/backd-io/backd && go test ./...

RUN sh -x \
    && for ARTIFACT in admin auth backd functions objects sessions ; \
    do \
        if [ ! -f /app/$ARTIFACT ]; then \
            cd /gopath/src/github.com/backd-io/backd/cmd/$ARTIFACT && \
            CGO_ENABLED=0 go build -o /app/$ARTIFACT . ; \
        fi ; \
    done

# admin
FROM alpine:latest

ARG ARTIFACT
ARG BUILD_DATE
ARG VCS_URL
ARG VCS_REF
ARG VERSION
ARG API_PORT
ARG METRICS_PORT

LABEL org.label-schema.build-date=${BUILD_DATE}                      \
      org.label-schema.name="backd - ${ARTIFACT}"                    \
      org.label-schema.description="backd - Application Backend"     \
      org.label-schema.url="http://backd.io/"                        \
      org.label-schema.vcs-ref=${VCS_REF}                            \
      org.label-schema.vcs-url=https://${VCS_URL}                    \
      org.label-schema.vendor="antoniofernandezvara+backd@gmail.com" \
      org.label-schema.version=${VERSION}                            \
      org.label-schema.schema-version="1.0"                          \
      maintainer="antoniofernandezvara+backd@gmail.com"

RUN apk add --update --no-cache ca-certificates

COPY --from=builder /app/${ARTIFACT} /artifact 

EXPOSE ${API_PORT}
EXPOSE ${METRICS_PORT}

ENTRYPOINT ["/artifact"]