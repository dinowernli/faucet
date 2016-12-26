#!/bin/sh

bazel build client && ./bazel-bin/client/client
