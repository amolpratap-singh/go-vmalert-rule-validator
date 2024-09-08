#GO Builder


#First Layer Compilation
FROM ${GO_BUILDER_OS}:latest as builder


RUN mkdir -p /usr/app/{bin,src,pkg}

ENV GOPATH=/usr/app
ENV PATH=${PATH}:${GOPATH}/bin

#COPY docker/vmalert-init-container/src /usr/app/src

WORKDIR ${GOPATH}/
#Build logdata_easyjson
RUN go mod init vmalert \
 && go mod tidy \
 && go install main.go