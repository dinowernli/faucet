load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "client",
    srcs = ["client.go"],
    deps = [
        "//proto/service/coordinator:go_default_library",
        "@com_github_sirupsen_logrus//:go_default_library",
    ],
)
