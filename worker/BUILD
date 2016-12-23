load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

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

go_test(
    name = "go_default_test",
    size = "small",
    srcs = glob(["*_test.go"]),
    deps = [
        "@org_github_stretchr_testify//assert:go_default_library",
    ],
    library = ":go_default_library",
)
