FROM fnproject/go:dev as build-stage
WORKDIR /function
RUN go get -u github.com/golang/dep/cmd/dep
ADD . /go/src/func/
RUN cd /go/src/func/ && dep ensure
RUN cd /go/src/func/ && go build -o func

FROM fnproject/go
WORKDIR /function
COPY --from=build-stage /go/src/func/func /function/
# Added OCI private key for OCI Go SDK authentication
COPY oci_private_key.pem /function/oci_private_key.pem
ENTRYPOINT ["./func"]
