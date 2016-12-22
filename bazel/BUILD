load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["bazel.go"],
    deps = [
        "//proto/workspace:go_default_library",
    ],
    visibility = ["//visibility:public"],
)
