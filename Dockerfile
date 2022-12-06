FROM golang:1.17-alpine AS staging

COPY . /workspace
WORKDIR /workspace
ENV GO111MODULE=on
ENV TZ=Asia/Taipei
RUN apk add bash tzdata make && make

ENTRYPOINT ["/workspace/collector"]