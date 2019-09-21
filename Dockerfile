FROM abchain/fabric:base AS coderepo

RUN go get -u github.com/op/go-logging
RUN go get -u github.com/spf13/viper
RUN go get -u github.com/gocraft/web

RUN rm -rf $(find ${GOPATH}/src -name .git -type d)

FROM abchain/fabric:base_0.97
COPY --from=coderepo ${GOPATH}/src ${GOPATH}/src/
COPY / ${GOPATH}/src/hyperledger.abchain.org/

