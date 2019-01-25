FROM alpine:3.8

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

COPY bin/${ARTIFACT} /artifact 

EXPOSE ${API_PORT}
EXPOSE ${METRICS_PORT}

ENTRYPOINT ["/artifact"]
