load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["loader.go"],
    deps = [
        "//proto/config:go_default_library",
        "@com_github_golang_protobuf//jsonpb:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    size = "small",
    srcs = ["loader_test.go"],
    deps = [
        "@org_github_stretchr_testify//assert:go_default_library",
    ],
    library = ":go_default_library",
)
