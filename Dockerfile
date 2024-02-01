# syntax=docker/dockerfile:1
#FROM gcr.io/wordless-412000/github.com/rwbutts/wordless_vue_dist:latest AS vue-app-dist

# Build the application from source
FROM golang:1.21 AS build-stage

ARG TAG_NAME='0.0.0'

WORKDIR /app
COPY go.* ./

COPY /main/* /app/main/
COPY /words/* /app/words/
#COPY /wwwroot/ /wwwroot/
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.VERSION=${TAG_NAME}" -o /wordless ./main

# Deploy the application binary into a lean image
FROM scratch

WORKDIR /

COPY --from=build-stage /wordless /wordless

