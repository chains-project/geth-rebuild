ARG UBUNTU_DIST="focal"

FROM ubuntu:${UBUNTU_DIST} as builder

ARG GETH_DIR=""
ARG GO_VERSION=""
ARG GETH_VERSION=""
ARG OS="" 
ARG ARCH=""
ARG GETH_COMMIT=""
ARG SHORT_COMMIT=""
ARG BUILD_CMD=""
ARG ARM_V=""
ARG C_COMPILER=""
ARG PACKAGES=""
ARG ELF_TARGET="elf64-x86-64"

RUN apt-get update && apt-get install -yq --no-install-recommends --force-yes \
    ${PACKAGES} \
    ${C_COMPILER}


#ln -s /usr/include/asm-generic /usr/include/asm
# TODO need to check that gcc version is the same

# Install Go
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
ENV GOARM=$ARM_V
ENV CGO_ENABLED=1

# Fetch reference binary and strip symbols + build ids
ENV BIN_DIR="geth-${OS}-${ARCH}-${GETH_VERSION}-${SHORT_COMMIT}"
ENV TAR_DIR="${BIN_DIR}.tar.gz"
ENV REF_URL="https://gethstore.blob.core.windows.net/builds/${TAR_DIR}"

RUN wget ${REF_URL} && \ 
    tar -xvf ${TAR_DIR} && \
    cd ${BIN_DIR} && \
    mv geth /bin/geth-reference && \
    strip --input-target=${ELF_TARGET} --remove-section .note.go.buildid --remove-section .note.gnu.build-id /bin/geth-reference

# Copy geth repo and rebuild reference binary
# TODO decide if should clone again
COPY ./tmp/go-ethereum /go-ethereum 

RUN echo $GOARM
RUN cd go-ethereum && git fetch && git checkout -b geth-reproduce ${GETH_COMMIT} && \
    ${BUILD_CMD} ./cmd/geth

RUN mv go-ethereum/build/bin/geth /bin/geth-reproduce && \
    strip --input-target=${ELF_TARGET} --remove-section .note.go.buildid --remove-section .note.gnu.build-id /bin/geth-reproduce


FROM alpine:latest

COPY --from=builder /bin/geth-reference /bin/geth-reference
COPY --from=builder /bin/geth-reproduce /bin/geth-reproduce

# TODO compare here??