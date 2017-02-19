load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["coordinator.go"],
    deps = [
        "//config:go_default_library",
        "//coordinator/storage:go_default_library",
        "//proto/config:go_default_library",
        "//proto/service/coordinator:go_default_library",
        "//proto/service/worker:go_default_library",
        "//proto/storage:go_default_library",
        "@com_github_sirupsen_logrus//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_x_net//context:go_default_library",
    ],
    visibility = ["//visibility:public"],
)
