load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "coordinator",
    srcs = ["coordinator.go"],
    deps = [
        "//proto/config:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
