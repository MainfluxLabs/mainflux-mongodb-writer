###
# Mainflux MongoDB Writer Dockerfile
###

FROM golang:alpine
MAINTAINER Mainflux

ENV MONGO_HOST mongo
ENV MONGO_PORT 27017

ENV NATS_HOST nats
ENV NATS_PORT 4222

###
# Install
###
# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/mainflux/mainflux-mongodb-writer
RUN cd /go/src/github.com/mainflux/mainflux-mongodb-writer && go install

###
# Run main command with dockerize
###
CMD mainflux-mongodb-writer -m $MONGO_HOST -n NATS_HOST
