class Cadreen < Formula
  desc "Intelligence infrastructure for developers"
  homepage "https://accomplishanything.today/infra"
  version "0.4.2"
  license "UNLICENSED"

  on_macos do
    on_arm do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_darwin_arm64"
      sha256 "9c5a7b9f94c08a5ab86ed92a6513a327bbc135cfa6bec759e17515a1332cc6f5"
    end
    on_intel do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_darwin_amd64"
      sha256 "5b346fed65b9a5fbfeaed95fe5ef651b6ae52f94281ffb8264bbe1486ba1ffdf"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_linux_arm64"
      sha256 "a9e183290c5387e6dd41fcbb21624dd0c66babf2b153192f0124b83e4c0480c9"
    end
    on_intel do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_linux_amd64"
      sha256 "ddc2e2047e8345d48431e811f591a470c865d80d67c01ede74bc2c1292855d40"
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
