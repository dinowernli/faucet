load("@io_bazel_rules_go//proto:go_proto_library.bzl", "go_proto_library")

go_proto_library(
    name = "go_default_library",
    srcs = ["coordinator.proto"],
    deps = [
        "//proto/workspace:go_default_library",
    ],
    visibility = ["//visibility:public"],
    has_services = 1,
)

