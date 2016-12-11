load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "worker",
    srcs = ["worker.go"],
    deps = [
        "//demo:go_default_library",
        "//proto/config:go_default_library",
        "//proto/service/worker:go_default_library",
    ],
)
