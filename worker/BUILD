load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["queue.go", "worker.go"],
    deps = [
        "//bazel:go_default_library",
        "//proto/service/worker:go_default_library",
        "//proto/workspace:go_default_library",
        "@org_golang_x_net//context:go_default_library",
    ],
    visibility = ["//visibility:public"],
)
