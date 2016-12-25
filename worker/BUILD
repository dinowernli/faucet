load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["worker.go"],
    deps = [
        "//bazel:go_default_library",
        "//proto/service/worker:go_default_library",
        "//proto/workspace:go_default_library",
        "//worker/scheduler:go_default_library",
        "@com_github_sirupsen_logrus//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_x_net//context:go_default_library",
    ],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    size = "small",
    srcs = glob(["*_test.go"]),
    deps = [
        "@org_github_stretchr_testify//assert:go_default_library",
    ],
    library = ":go_default_library",
)
