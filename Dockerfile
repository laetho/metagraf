FROM gcr.io/distroless/static:nonroot
COPY . /
USER nonroot:nonroot

ENTRYPOINT ["/mg"]