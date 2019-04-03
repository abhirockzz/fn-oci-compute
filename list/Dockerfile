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
ARG PRIVATE_KEY_NAME
COPY $PRIVATE_KEY_NAME /function/$PRIVATE_KEY_NAME
# OCI_PRIVATE_KEY_FILE_NAME is used as environment variable in the function code. Altered name to avoid confusion
ENV OCI_PRIVATE_KEY_FILE_NAME=${PRIVATE_KEY_NAME}
ENTRYPOINT ["./func"]
