# start from golang image based on alpine-3.8
FROM golang:1.12-alpine3.10 AS dev-build

# add our cgo dependencies
RUN apk add --update --no-cache ca-certificates cmake make g++ openssl-dev git curl pkgconfig
# clone seabolt-1.7.0 source code
RUN git clone -b 1.7 https://github.com/neo4j-drivers/seabolt.git /seabolt
# invoke cmake build and install artifacts - default location is /usr/local
WORKDIR /seabolt/build
# CMAKE_INSTALL_LIBDIR=lib is a hack where we override default lib64 to lib to workaround a defect
# in our generated pkg-config file 
RUN cmake -D CMAKE_BUILD_TYPE=Release -D CMAKE_INSTALL_LIBDIR=lib .. && cmake --build . --target install

# install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/larien/family-tree
ADD . /go/src/github.com/larien/family-tree

ENV GO111MODULES=on

# Install Dependencies
RUN go mod download
RUN go get "github.com/neo4j/neo4j-go-driver/neo4j"
RUN go get "github.com/gin-contrib/cors"
RUN go get "github.com/gin-gonic/gin"

CMD ["go", "run", "main.go"]