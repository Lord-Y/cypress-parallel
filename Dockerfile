ARG DOCKER_IMAGE_NODEJS="node:14-alpine"
ARG DOCKER_IMAGE_GOLANG="golang:1.17-alpine"

FROM ${DOCKER_IMAGE_NODEJS} as buildernode
RUN mkdir /app
WORKDIR /app
COPY ui .
RUN rm -rf .git .env*
RUN apk update && apk upgrade
RUN apk add --no-cache git curl && rm -rf /var/cache/apk/*
RUN npm install
RUN npm run build

FROM ${DOCKER_IMAGE_GOLANG} as buildergo
RUN apk update && apk add --no-cache git ca-certificates
RUN adduser -D -g '' appuser
COPY . $GOPATH/src/
WORKDIR $GOPATH/src/
RUN rm -rf ui/dist routers/ui/dist
COPY --from=buildernode /app/dist ui/dist
COPY --from=buildernode /app/dist routers/ui/dist
RUN rm -rf $GOPATH/pkg/* $GOPATH/src/go.sum $GOPATH/.git /var/cache/apk/*
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN:$PATH
ENV GO111MODULE=on
RUN go env -w GOFLAGS=-mod=mod
RUN go env
RUN go mod download
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -a -tags netgo -ldflags '-w -extldflags "-static"' .

FROM scratch
# Import from builder.
COPY --from=buildergo /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=buildergo /etc/passwd /etc/passwd
COPY --from=buildergo /go/bin/cypress-parallel /go/bin/
USER appuser
ENTRYPOINT ["/go/bin/cypress-parallel"]
