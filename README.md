## Go-Ethereum rpc.Client#Unsubscribe blocks indefinitely

when called simultaneously with geth node shutdown:

[![asciicast](https://asciinema.org/a/204548.png)](https://asciinema.org/a/204548)

### Repro steps:

To help emulate the environment needed to reproduce the bug, `repro.sh` bootstraps,
runs and kills `geth` and `main.go`.

1. `./repro.sh`
2. `ps` aftwerwards and ensure the go executable (ending with `blockingbuild`) is still running. Note that it will not respond to a regular `kill` command.
3. Install [dlv](https://github.com/derekparker/delve) (`go get -u github.com/derekparker/delve/cmd/dlv`)
4. Use `dlv attach <pid> <binarypath>` to attach to process
5. List goroutines with `goroutines`
6. Switch to goroutine stuck on go-ethereum `send`, i.e. `goroutine 1`
7. Look at call stack with `stack`

This issue is indeed spurious, so step 2 won't always reproduce, but in practice it
reproduced 75% of the time.
