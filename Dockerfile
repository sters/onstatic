ARG GO_VERSION=1.17

FROM dockercore/golang-cross

ARG GORELEASER_VERSION=0.183.0
ARG GORELEASER_DOWNLOAD_FILE=goreleaser_Linux_x86_64.tar.gz
ARG GORELEASER_DOWNLOAD_URL=https://github.com/goreleaser/goreleaser/releases/download/v${GORELEASER_VERSION}/${GORELEASER_DOWNLOAD_FILE}
ENV GO111MODULE=on

RUN  wget ${GORELEASER_DOWNLOAD_URL}; \
    tar -xzf $GORELEASER_DOWNLOAD_FILE -C /usr/bin/ goreleaser; \
    rm $GORELEASER_DOWNLOAD_FILE;

CMD ["goreleaser", "-v"]
