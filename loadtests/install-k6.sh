#!/bin/sh
go install go.k6.io/xk6/cmd/xk6@latest
xk6 build --with github.com/avitalique/xk6-file@latest