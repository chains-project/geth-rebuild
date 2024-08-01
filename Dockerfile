ARG UBUNTU_DIST=""

FROM ubuntu:${UBUNTU_DIST} as builder


ARG GETH_SRC_DIR="./tmp/go-ethereum"

# artifact spec
ARG OS="" 
ARG ARCH=""
ARG GETH_VERSION=""
ARG GETH_COMMIT=""
ARG SHORT_COMMIT=""

# toolchain spec
ARG GO_VERSION=""
ARG CC=""
ARG BUILD_CMD=""
ARG TOOLCHAIN_DEPS=""

# environment spec
ARG GOOS=""
ARG GOARCH=""
ARG GOARM=""
ARG ELF_TARGET=""
ARG UTIL_DEPS=""

ENV CGO_ENABLED=1
ENV PATH=/usr/local:/usr/bin:/usr/local/go/bin:$PATH


RUN dpkg --add-architecture armel && apt-get update && apt-get install -yq --no-install-recommends --force-yes \
    ${UTIL_DEPS} \
    ${TOOLCHAIN_DEPS}

# Install Go
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz


ENV CROSS_COMPILE=arm-linux-gnueabi-
ENV LD_LIBRARY_PATH=/usr/arm-linux-gnueabi/lib/:/usr/arm-linux-gnueabi/usr/lib/
#test
#RUN ln -s /usr/include/asm-generic /usr/include/asm

# ENV qemu-arm -L /usr/arm-linux-gnueabi/
# RUN unset LD_LIBRARY_PATH

# RUN cp /usr/arm-linux-gnueabi/lib/ld-linux.so.3 /lib
# RUN cp /usr/arm-linux-gnueabi/lib/libgcc_s.so.1 /lib
# RUN cp /usr/arm-linux-gnueabi/lib/libc.so.6 /lib
# RUN cp /usr/arm-linux-gnueabi/lib/libresolv.so.2 /lib
# RUN cp /usr/arm-linux-gnueabi/lib/libpthread.so.0 /lib

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
ENV GETH_DEST_DIR=/go-ethereum
COPY ${GETH_SRC_DIR} ${GETH_DEST_DIR} 

RUN cd ${GETH_DEST_DIR} && git fetch && git checkout -b geth-reproduce ${GETH_COMMIT} && \
     ${BUILD_CMD} ./cmd/geth


RUN cd ${GETH_DEST_DIR}/build/bin && \
    strip --input-target=${ELF_TARGET} --remove-section .note.go.buildid --remove-section .note.gnu.build-id geth && \
    mv geth ${REPRODUCE_DEST}

# Second stage build for compact final image
FROM alpine:latest

COPY --from=builder /bin/geth-reference /bin/geth-reference
COPY --from=builder /bin/geth-reproduce /bin/geth-reproduce