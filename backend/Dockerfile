FROM golang:1.12-alpine3.10 AS dev-build

RUN apk add --update --no-cache ca-certificates cmake make g++ openssl-dev git curl pkgconfig
RUN git clone -b 1.7 https://github.com/neo4j-drivers/seabolt.git /seabolt
WORKDIR /seabolt/build
RUN cmake -D CMAKE_BUILD_TYPE=Release -D CMAKE_INSTALL_LIBDIR=lib .. && cmake --build . --target install

WORKDIR /go/src/github.com/larien/family-tree
ADD . /go/src/github.com/larien/family-tree

ENV GO111MODULES=on

RUN go get "github.com/neo4j/neo4j-go-driver/neo4j"
RUN go get "github.com/gin-contrib/cors"
RUN go get "github.com/gin-gonic/gin"

EXPOSE 8899

CMD ["go", "run", "main.go"]