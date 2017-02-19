#!/bin/sh

bazel build //client && $FAUCET_ROOT/bazel-bin/client/client
