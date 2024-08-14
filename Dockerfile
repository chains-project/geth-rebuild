ARG UBUNTU_DIST=""

ARG REFERENCE_LOC=/bin/geth-reference
ARG REPRODUCE_LOC=/bin/geth-reproduce

FROM ubuntu:${UBUNTU_DIST} as builder

ARG REFERENCE_LOC
ARG REPRODUCE_LOC

# Artifact spec
ARG OS="" 
ARG ARCH=""
ARG GETH_VERSION=""
ARG GETH_COMMIT=""
ARG SHORT_COMMIT=""

# Toolchain spec
ARG GO_VERSION=""
ARG BUILD_CMD=""
ARG TOOLCHAIN_DEPS=""

# Environment spec
ARG URL_SLUG=""
ARG GOARM=""
ARG ELF_TARGET=""
ARG UTIL_DEPS=""
# For cross compilation, CGO must initially be disabled to avoid error
# CGO is enabled later on in the geth build script
ARG CGO_ENABLED=""

# Install packages
RUN apt-get update && apt-get install -yq --no-install-recommends --force-yes \
    ${TOOLCHAIN_DEPS} \
    ${UTIL_DEPS}

#RUN ln -s /usr/include/asm-generic /usr/include/asm

# Install Go
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin


# Fetch reference binary and strip symbols + build ids
ENV BIN_DIR=${URL_SLUG}
ENV TAR_DIR="${BIN_DIR}.tar.gz"
ENV REF_URL="https://gethstore.blob.core.windows.net/builds/${TAR_DIR}"
RUN wget ${REF_URL} && \ 
    tar -xvf ${TAR_DIR} && \
    cd ${BIN_DIR} && \
    strip --input-target=${ELF_TARGET} --remove-section .note.go.buildid --remove-section .note.gnu.build-id geth && \
    mv geth ${REFERENCE_LOC}


# Copy local geth repo
ENV GETH_SRC_DIR=./tmp/go-ethereum
ENV GETH_DIR=/go-ethereum
COPY ${GETH_SRC_DIR} ${GETH_DIR} 

# Rebuild the reference binary
WORKDIR ${GETH_DIR}
RUN git fetch && git checkout -b geth-reproduce ${GETH_COMMIT} && \
    ${BUILD_CMD} ./cmd/geth


# Strip symbols and build ids
WORKDIR ${GETH_DIR}/build/bin
RUN strip --input-target=${ELF_TARGET} --remove-section .note.go.buildid --remove-section .note.gnu.build-id geth && \
    mv geth ${REPRODUCE_LOC}


# Second stage build for compact final image
FROM alpine:latest

# Copy produced binaries
ARG REFERENCE_LOC
ARG REPRODUCE_LOC

COPY --from=builder ${REFERENCE_LOC} ${REFERENCE_LOC}
COPY --from=builder ${REPRODUCE_LOC} ${REPRODUCE_LOC}

# Get binary comparison script 
#TODO send in as ARG
ENV SCRIPT_SRC=./internal/scripts/compare_binary_SHA.sh 
ENV COMPARE_SCRIPT=/bin/compare_binary_SHA.sh 

COPY ${SCRIPT_SRC} ${COMPARE_SCRIPT}
RUN chmod +x ${COMPARE_SCRIPT}

# Run binary verification/comparison on run
CMD ["/bin/compare_binary_SHA.sh", "/bin/geth-reference", "/bin/geth-reproduce"]