FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64
ENV GOPROXY https://goproxy.cn,direct


WORKDIR /build
COPY . .
RUN sh -c "[ -f go.mod ]" || exit
COPY /etc /app/etc
RUN go build -ldflags="-s -w" -tags no_k8s -o /app/user ./user.go


FROM alpine

ENV TZ Asia/Shanghai
ENV CONFIG_FILE=etc/user.yaml

WORKDIR /app
COPY --from=builder /app/user /app/user
COPY --from=builder /app/etc /app/etc

EXPOSE 8081

CMD ["sh", "-c", "./user -f ${CONFIG_FILE}"]