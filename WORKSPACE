workspace(name = "me_dinowernli_faucet")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "afec53d875013de6cebe0e51943345c587b41263fdff36df5ff651fbf03c1c08",
    strip_prefix = "rules_go-0.4.4",
    url = "https://github.com/bazelbuild/rules_go/archive/0.4.4.tar.gz",
)

load("@io_bazel_rules_go//go:def.bzl", "go_repositories", "new_go_repository")
load("@io_bazel_rules_go//proto:go_proto_library.bzl", "go_proto_repositories")

go_repositories()

go_proto_repositories()

http_archive(
    name = "io_bazel",
    sha256 = "8e4646898fa9298422e69767752680d34cbf21bcae01c401b11aa74fcdb0ef66",
    strip_prefix = "bazel-0.4.4",
    url = "https://github.com/bazelbuild/bazel/archive/0.4.4.tar.gz",
)

new_go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    tag = "v1.1.0",
)

new_go_repository(
    name = "com_github_pmezard_go_difflib",
    importpath = "github.com/pmezard/go-difflib",
    tag = "v1.0.0",
)

new_go_repository(
    name = "com_github_sirupsen_logrus",
    importpath = "github.com/Sirupsen/logrus",
    tag = "v0.11.0",
)

new_go_repository(
    name = "com_github_stretchr_objx",
    importpath = "github.com/stretchr/objx",
    commit = "1a9d0bb9f541897e62256577b352fdbc1fb4fd94",
)

new_go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    tag = "v1.1.3",
)

