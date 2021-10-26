# Initial stage: download modules
FROM golang:1.17 as modules

ADD go.mod go.sum /app/
RUN cd /app && go mod download

# Intermediate stage: Build the binary
FROM golang:1.17 as builder

COPY --from=modules /go/pkg /go/pkg

RUN useradd -u 10001 shortener
RUN mkdir -p /shortener
ADD . /shortener
WORKDIR /shortener

RUN make build

# Final stage: Run the binary
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
USER shortener

COPY --from=builder /shortener/build/shortener /shortener
COPY --from=builder /shortener/web/ /web/

CMD ["/shortener"]
