#!/usr/bin/env bash

#hack to fix terminal after game exits

go run snake.go
/bin/stty sane
