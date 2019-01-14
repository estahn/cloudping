FROM scratch

ARG BUILD_DATE
ARG VCS_REF
ARG VERSION

COPY cloudping /

ENTRYPOINT [ "/cloudping" ]

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="cloudping" \
      org.label-schema.description="cloudping identifies the regions geographically closest and returns them in order of lowest to highest `response time`." \
      org.label-schema.url="https://enricostahn.com/" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/estahn/cloudping" \
      org.label-schema.vendor="Enrico Stahn" \
      org.label-schema.version=$VERSION \
      org.label-schema.schema-version="1.0" \
      org.label-schema.docker.cmd="docker run -it --rm estahn/cloudping"
