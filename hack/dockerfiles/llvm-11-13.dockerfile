# syntax=docker/dockerfile:1.3

FROM --platform=$BUILDPLATFORM ubuntu:21.10 AS base
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
	apt-get update && \
	apt-get install -y --no-install-recommends \
		git \
		libc6-dev \
		libgcc-10-dev \
		make \
		pkg-config \
		\
		libllvm${LLVM_VERSION} \
		llvm-${LLVM_VERSION} \
		llvm-${LLVM_VERSION}-dev \
		llvm-${LLVM_VERSION}-runtime \
		clang-${LLVM_VERSION} \
		libclang-common-${LLVM_VERSION}-dev \
		libclang-${LLVM_VERSION}-dev \
		libclang1-${LLVM_VERSION} \
	&& \
	ln -s /usr/bin/clang-${LLVM_VERSION} /usr/bin/clang && \
	ln -s /usr/bin/clang++-${LLVM_VERSION} /usr/bin/clang++ && \
	ln -s /usr/bin/llvm-config-${LLVM_VERSION} /usr/bin/llvm-config && \
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
