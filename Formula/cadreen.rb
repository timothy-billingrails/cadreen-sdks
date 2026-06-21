class Cadreen < Formula
  desc "Intelligence infrastructure for developers"
  homepage "https://accomplishanything.today"
  version "0.2.1"
  license "UNLICENSED"

  on_macos do
    on_arm do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_darwin_arm64"
      sha256 "580b5ba99038cf607189b015499be74d3c50a2349c57394cbcd2c88e7b20a2df"
    end
    on_intel do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_darwin_amd64"
      sha256 "f14e80fab4835e61a658a7219965353e9389ea9e0801900b8244ee9fb8c1ce49"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_linux_arm64"
      sha256 "da1722f295a9a73d3a5435ce92790b58fe2f7d273dde666763554b66db27aa0c"
    end
    on_intel do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_linux_amd64"
      sha256 "ccd03c2430dbd85c5694ec994e5e44cb63d10124aa6cd48097a4917e96f95d39"
    end
  end

  def install
    bin.install "cadreen"
    generate_completions_from_executable(bin/"cadreen", "completion")
  end

  test do
    assert_match "cadreen", shell_output("#{bin}/cadreen version")
  end
end
