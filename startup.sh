#!/bin/bash

echo "Starting dbproxy server..."
go run dbproxy/server/cmd/main.go &
DBPROXY_PID=$!

echo "Starting account server..."
go run account/server/cmd/main.go &
ACCOUNT_PID=$!

echo "Starting file server..."
go run file/cmd/main.go &
FILE_PID=$!

echo "Starting transfer service..."
go run transfer/cmd/main.go &
TRANSFER_PID=$!

echo "Starting gateway server..."
go run gateway/cmd/main.go &
GATEWAY_PID=$!

# 等待所有后台进程完成
wait $DBPROXY_PID
wait $ACCOUNT_PID
wait $FILE_PID
wait $TRANSFER_PID
wait $GATEWAY_PID