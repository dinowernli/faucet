#!/bin/sh

bazel build //demo && $FAUCET_ROOT/bazel-bin/demo/demo
