package main

const Version = "0.20.4"

// We want to replace this variable at build time with "-ldflags -X main.GitSHA=xxx", where const is not supported.
var GitSHA = ""
