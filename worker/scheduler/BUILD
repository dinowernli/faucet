load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["scheduler.go"],
    deps = [
        "//bazel:go_default_library",
        "//proto/service/worker:go_default_library",
        "//proto/workspace:go_default_library",
        "//worker/checkout:go_default_library",
        "@com_github_sirupsen_logrus//:go_default_library",
    ],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    size = "small",
    srcs = glob(["*_test.go"]),
    deps = [
        "@com_github_stretchr_testify//assert:go_default_library",
        "@com_github_stretchr_testify//mock:go_default_library",
    ],
    library = ":go_default_library",
)

