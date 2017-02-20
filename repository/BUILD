load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["repository.go"],
    deps = [
        "//proto/workspace:go_default_library",
        "@com_github_sirupsen_logrus//:go_default_library",
    ],
    visibility = ["//visibility:public"],
)
