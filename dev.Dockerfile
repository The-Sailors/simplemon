#Docker file to be used for development purposes
ARG GO_VERSION
FROM golang:${GO_VERSION}-bullseye

RUN apt-get update && \
    apt-get install -y make

WORKDIR /app/src
COPY . .

#build the golang binary
RUN make simplemon 

EXPOSE 8080

#run the binary
CMD ["go", "run", "./cmd/simplemon/"]
