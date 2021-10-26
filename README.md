[![Go Report Card](https://goreportcard.com/badge/github.com/phpCoder88/url-shortener-observable)](https://goreportcard.com/report/github.com/phpCoder88/url-shortener-observable)
[![codecov](https://codecov.io/gh/phpCoder88/url-shortener-observable/branch/master/graph/badge.svg)](https://codecov.io/gh/phpCoder88/url-shortener-observable)

# URL shortener

## Launch migrations using docker

Before launching the command below, check DSN for correctness

```bash
docker run -v "$PWD"/migrations:/migrations \
  --network host \
  migrate/migrate:v4.14.1 \
  -path=/migrations/ \
  -database "postgres://shortener:123456789@localhost:5432/shortener?sslmode=disable" \
  up
```
