# syntax=docker/dockerfile:1.3

ARG UBUNTU_VERSION

FROM --platform=$BUILDPLATFORM ubuntu:${UBUNTU_VERSION} AS base
ENV DEBIAN_FRONTEND=noninteractive \
	GOPATH=/go \
	PATH=/go/bin:/usr/local/go/bin:$PATH \
	CC=clang \
	CXX=clang++

FROM --platform=$BUILDPLATFORM base AS golang
ARG GOLANG_VERSION
RUN set -ex && \
	apt-get update && \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		curl \
	&& \
	curl -fsS "https://storage.googleapis.com/golang/go${GOLANG_VERSION}.linux-amd64.tar.gz" | tar -xzf - -C /usr/local && \
	go version && \
	go env && \
	mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH" && \
	\
	apt-get purge -y curl

FROM --platform=$BUILDPLATFORM golang AS llvm
ARG LLVM_VERSION
RUN set -eux && \
	case "${LLVM_VERSION}" in \
		(4|5|6) \
			APT_LLVM_VERSION=${LLVM_VERSION}.0 \
			LIBGCC_VERSION=7 \
			;; \
		*) \
			APT_LLVM_VERSION=${LLVM_VERSION} \
			LIBGCC_VERSION=10 \
			;; \
	esac && \
	\
	apt-get update && \
	apt-get install -y --no-install-recommends \
		git \
		libc6-dev \
		libgcc-${LIBGCC_VERSION}-dev \
		make \
		pkg-config \
		\
		libllvm${APT_LLVM_VERSION} \
		llvm-${APT_LLVM_VERSION} \
		llvm-${APT_LLVM_VERSION}-dev \
		llvm-${APT_LLVM_VERSION}-runtime \
		clang-${APT_LLVM_VERSION} \
		libclang-common-${APT_LLVM_VERSION}-dev \
		libclang-${APT_LLVM_VERSION}-dev \
		libclang1-${APT_LLVM_VERSION} \
	&& \
	ln -s /usr/bin/clang-${APT_LLVM_VERSION} /usr/bin/clang && \
	ln -s /usr/bin/clang++-${APT_LLVM_VERSION} /usr/bin/clang++ && \
	ln -s /usr/bin/llvm-config-${APT_LLVM_VERSION} /usr/bin/llvm-config && \
	\
	rm -rf \
		 /var/cache/debconf/* \
		 /var/lib/apt/lists/* \
		 /var/log/* \
		 /tmp/* \
		 /var/tmp/*

FROM --platform=$BUILDPLATFORM llvm AS gen
WORKDIR /go/src/github.com/go-clang/gen
COPY go.mod go.sum .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache \
	go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache \
	set -eux && \
	make build/${LLVM_VERSION}
