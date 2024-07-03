ARG UBUNTU_VERSION="focal"

FROM ubuntu:${UBUNTU_VERSION} as builder

ARG GETH_DIR=""
ARG GO_VERSION=""
ARG GETH_VERSION=""
ARG OS="" 
ARG ARCH=""
ARG GETH_COMMIT=""
ARG SHORT_COMMIT=""
ARG BUILD_CMD=""
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
COPY ./tmp/go-ethereum /go-ethereum

RUN echo ${BUILD_CMD}
RUN cd go-ethereum && git fetch && git checkout -b geth-reproduce ${GETH_COMMIT} && \
    CGO_ENABLED=1 ${BUILD_CMD} ./cmd/geth

RUN mv go-ethereum/build/bin/geth /bin/geth-reproduce && \
    strip --input-target=${ELF_TARGET} --remove-section .note.go.buildid --remove-section .note.gnu.build-id /bin/geth-reproduce


FROM alpine:latest

COPY --from=builder /bin/geth-reference /bin/geth-reference
COPY --from=builder /bin/geth-reproduce /bin/geth-reproduce

# TODO compare here??