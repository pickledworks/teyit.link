# First, let's build the frontend
FROM node:alpine AS clientassetsbuilder

WORKDIR /app
COPY package.json package.json ./
RUN npm install

COPY client client
COPY public public

RUN npm run build

# Second, we need to build the backend application
# Copied and modified from https://gist.github.com/pierreprinetti/8a9a0f4602dc4e0b95ad685edddfb5ae
FROM golang:1.11 AS builder

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR $GOPATH/src/gitlab.com/nod/teyit/link
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

# We use the statik library to embed static files into the backend binary
RUN go get github.com/rakyll/statik

# Copy needed files
COPY database database
COPY handlers handlers
COPY utils utils
# Get the public assets from the frontend builder container
COPY --from=clientassetsbuilder /app/public public
COPY main.go main.go

# Generate a binary with our public assets bundled in
RUN go generate main.go

# Copy the code from the host and compile it
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

# Creating an empty container
FROM scratch

# Get the trusted CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy our executable from the Go builder container
COPY --from=builder /app ./

# We need the template files
COPY views views

ENTRYPOINT ["./app"]

# Let our proxy know which port to attach inside the container
EXPOSE 8080
