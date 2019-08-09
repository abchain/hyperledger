FROM abchain/fabric:base AS coderepo

RUN go get -u github.com/op/go-logging
RUN go get -u github.com/spf13/viper
RUN go get -u github.com/gocraft/web

RUN rm -rf $(find ${GOPATH}/src -name .git -type d)

FROM abchain/fabric
COPY --from=coderepo ${GOPATH}/src ${GOPATH}/src/
RUN go get -d hyperledger.abchain.org/cases/ae

