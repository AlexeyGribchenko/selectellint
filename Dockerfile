FROM golang:1.24.0-alpine AS builder

RUN apk add --no-cache gcc musl-dev git

WORKDIR /workspace

COPY . .

WORKDIR /workspace/plugin

RUN go mod download

RUN git clone https://github.com/golangci/golangci-lint.git /golangci-lint-src
WORKDIR /golangci-lint-src

RUN git checkout v2.5.0

RUN CGO_ENABLED=1 go build -o golangci-lint ./cmd/golangci-lint

WORKDIR /workspace

RUN CGO_ENABLED=1 go build -buildmode=plugin -o /workspace/selectellint.so ./plugin

FROM golang:1.24.0-alpine

RUN apk add --no-cache gcc musl-dev

COPY --from=builder /golangci-lint-src/golangci-lint /usr/local/bin/golangci-lint

COPY --from=builder /workspace/selectellint.so /usr/local/lib/selectellint.so

COPY --from=builder /workspace/.golangci.yaml /workspace/.golangci.yaml

WORKDIR /app

CMD ["golangci-lint", "version"]