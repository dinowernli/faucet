workspace(name = "me_dinowernli_faucet")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "b7759f01d29c075db177f688ffb4464aad2b8fbb7017f89a1d3819ce07f1d584",
    strip_prefix = "rules_go-0.3.1",
    url = "https://github.com/bazelbuild/rules_go/archive/0.3.1.tar.gz",
)

load("@io_bazel_rules_go//go:def.bzl", "go_repositories", "new_go_repository")
load("@io_bazel_rules_go//proto:go_proto_library.bzl", "go_proto_repositories")

go_repositories()

go_proto_repositories()

http_archive(
    name = "io_bazel",
    sha256 = "8e4646898fa9298422e69767752680d34cbf21bcae01c401b11aa74fcdb0ef66",
    strip_prefix = "bazel-0.4.1",
    url = "https://github.com/bazelbuild/bazel/archive/0.4.1.tar.gz",
)

new_go_repository(
    name = "org_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    commit = "f390dcf405f7b83c997eac1b06768bb9f44dec18",
)

new_go_repository(
    name = "com_github_pmezard_go_difflib",
    importpath = "github.com/pmezard/go-difflib",
    commit = "792786c7400a136282c1664665ae0a8db921c6c2",
)

new_go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    commit = "346938d642f2ec3594ed81d874461961cd0faa76",
)
