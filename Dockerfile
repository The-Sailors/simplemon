#build stage
ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine3.17 AS builder
RUN apk add --no-cache make
WORKDIR /app/src
COPY . .
RUN make simplemon
ENV PORT=8080

#final stage
FROM alpine:3.17.2
WORKDIR /app/src
RUN apk --no-cache add ca-certificates
# copy the generated binary from the build stage
COPY --from=builder /app/src/cmd/simplemon/simplemon /app/src/cmd/simplemon
ENV PORT=8080
EXPOSE $PORT
CMD ["app/src/cmd/simplemon/simplemon"]
