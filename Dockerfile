FROM golang:1.24-alpine AS build

RUN apk add --no-cache make git build-base bash

ENV PATH=$PATH:/go/bin
ADD . /src/oxia-nightly

RUN cd /src/oxia-nightly \
    && make

FROM alpine:3.21

RUN apk add --no-cache bash bash-completion jq
RUN apk upgrade --no-cache

RUN mkdir /oxia-nightly
WORKDIR /oxia-nightly

COPY --from=build /src/oxia-nightly/bin/night /oxia-nightly/bin/night
ENV PATH=$PATH:/oxia-nightly/bin

RUN oxia completion bash > ~/.bashrc

CMD [ "/bin/bash" ]
