FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY *.go ./

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /exporter .

FROM gcr.io/distroless/static-debian13:nonroot

ENV CORE_HOST="" \
    CORE_PORT=9851

COPY --from=build /exporter /exporter

EXPOSE 80

HEALTHCHECK --interval=60s --start-period=5s CMD ["/exporter", "healthcheck"]

ENTRYPOINT ["/exporter"]
