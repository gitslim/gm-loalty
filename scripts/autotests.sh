#!/bin/env bash

go build -buildvcs=false -o ./cmd/gophermart/gophermart ./cmd/gophermart && \
          gophermarttest \
            -test.v -test.run=^TestGophermart$ \
            -gophermart-binary-path=cmd/gophermart/gophermart \
            -gophermart-host=localhost \
            -gophermart-port=8080 \
            -gophermart-database-uri="postgresql://postgres:postgres@localhost/praktikum?sslmode=disable" \
            -accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
            -accrual-host=localhost \
            -accrual-port=$(random unused-port) \
            -accrual-database-uri="postgresql://postgres:postgres@localhost/praktikum?sslmode=disable"
