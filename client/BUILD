load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "client",
    srcs = ["client.go"],
    deps = [
        "//proto/client:go_default_library",
        "//proto/service/coordinator:go_default_library",
        "//proto/workspace:go_default_library",
        "@com_github_sirupsen_logrus//:go_default_library",
        "@com_github_golang_protobuf//jsonpb:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_x_net//context:go_default_library",

    ],
)
