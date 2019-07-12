#!/usr/bin/env bash

#hack to fix terminal after game exits

go run cmd/snake.go
/bin/stty sane
