#!/bin/bash

set -e

main() {
  installGeth
  runGethnet

  printf -- "\033[34mStarting main.go...\033[0m\n"
  buildAndRunSubscriber
  sleep 2
  printf -- "\033[34mmain.go is running.\033[0m\n"

  printps
  printf -- "\033[34mSleeping for 5...\033[0m\n"
  sleep 5
  printf -- "\033[34mSimultaneously kill all child processes. Check ps afterwards..\033[0m\n"
}

exit_handler() {
  # Clear all signal handlers to prevent handler loop
  trap - 1 2 3 15
  # Kill all child subprocesses
  kill -- -$$ || true # Unsubscription won't end
}

trap "exit_handler" EXIT SIGTERM SIGINT

buildAndRunSubscriber() {
  DIR=`mktemp -d`
  go build -o "${DIR}/blockingbuild" main.go
  ${DIR}/blockingbuild &
}

installGeth() {
  printf -- "\033[34mInstalling geth 1.8.27...\033[0m\n"
  ethpkg=github.com/ethereum/go-ethereum
  ethpath=$GOPATH/src/$ethpkg
  if [ -d "$ethpath" ]; then
    pushd "$ethpath" >/dev/null
    git checkout master &>/dev/null
    go get -d -u $ethpkg
  else
    go get -d $ethpkg
    pushd "$ethpath" >/dev/null
  fi
  git checkout v1.8.27 2>/dev/null
  popd >/dev/null
  go install $ethpkg/cmd/geth
}

waitForResponse ()
{
  printf -- "\033[34mWaiting for $1.\033[0m\n"
  sleepCount=0
  while [ "$sleepCount" -le "300" ] && ! curl -s "$1" >/dev/null; do
      sleep 1
      sleepCount=$((sleepCount+1))
  done

  if [ "$sleepCount" -gt "300" ]; then
    printf -- "\033[31mTimed out waiting for $1 (waited 300s).\033[0m\n"
    exit 1
  fi
  printf -- "\033[34mService on $1 is ready.\033[0m\n"
}

runGethnet() {
  printf -- "\033[34mStarting geth...\033[0m\n"
  ./gethnet &
  waitForResponse http://127.0.0.1:18545
  printf -- "\033[34mGeth is running.\033[0m\n"
}

printps() {
  printf -- "\033[34mProcesses...\033[0m\n"
  ps
}

main
