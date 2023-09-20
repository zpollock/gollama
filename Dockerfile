FROM golang:1.21.0

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy \
    && apt-get update \
        && apt-get install -y \
            nmap \
            vim  \ 
            less \
    && make -C chat/llama.cpp
    