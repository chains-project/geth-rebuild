ARG UBUNTU_DIST=""

FROM ubuntu:${UBUNTU_DIST} as builder

#ARG GETH_DIR=""
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
ARG ELF_TARGET=""

RUN apt-get update && apt-get install -yq --no-install-recommends \
    ${PACKAGES} \
    ${C_COMPILER}


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
ENV REFERENCE_DEST=/bin/geth-reference
ENV REPRODUCE_DEST=/bin/geth-reproduce

RUN wget ${REF_URL} && \ 
    tar -xvf ${TAR_DIR} && \
    cd ${BIN_DIR} && \
    strip --input-target=${ELF_TARGET} --remove-section .note.go.buildid --remove-section .note.gnu.build-id geth && \
    mv geth ${REFERENCE_DEST}

# Copy geth repo and rebuild reference binary
# TODO decide if should clone again
ENV GETH_DIR=/go-ethereum
COPY ./tmp/go-ethereum ${GETH_DIR} 

RUN cd ${GETH_DIR} && git fetch && git checkout -b geth-reproduce ${GETH_COMMIT} && \
    ${BUILD_CMD} ./cmd/geth

RUN cd ${GETH_DIR}/build/bin && \
    strip --input-target=${ELF_TARGET} --remove-section .note.go.buildid --remove-section .note.gnu.build-id geth && \
    mv geth ${REPRODUCE_DEST}


FROM alpine:latest

COPY --from=builder /bin/geth-reference /bin/geth-reference
COPY --from=builder /bin/geth-reproduce /bin/geth-reproduce

# TODO compare here??