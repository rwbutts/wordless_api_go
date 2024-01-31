# syntax=docker/dockerfile:1
FROM us-docker.pkg.dev/wordless-412000/wordless-app-dist/wordless-app-dist-prod:2.5.1 AS vue-app-dist

# Build the application from source
FROM golang:1.21 AS build-stage

WORKDIR /app
COPY go.* ./

COPY /main/* /app/main/
COPY /words/* /app/words/
#COPY /wwwroot/ /wwwroot/
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o /wordless ./main

# Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /wordless /wordless
COPY --from=vue-app-dist /dist/ /wwwroot/

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/wordless"]

