set shell := ["cmd.exe", "/c"]

build:
    set GOOS=linux
    set GOARCH=arm
    set GOARM=5
    go build -o out/gotemp

