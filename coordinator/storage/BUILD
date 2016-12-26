load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["storage.go", "in_memory.go"],
    deps = [
        "//proto/storage:go_default_library",
    ],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    size = "small",
    srcs = glob(["*_test.go"]),
    deps = [
        "@com_github_stretchr_testify//assert:go_default_library",
    ],
    library = ":go_default_library",
)

