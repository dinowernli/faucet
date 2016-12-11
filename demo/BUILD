load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "demo",
    srcs = ["demo.go"],
    deps = [
        "//coordinator:go_default_library",
        "//proto/config:go_default_library",
        "//proto/service/worker:go_default_library",
        "//worker:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
