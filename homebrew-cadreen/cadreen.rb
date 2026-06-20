class Cadreen < Formula
  desc "Intelligence infrastructure for developers"
  homepage "https://accomplishanything.today"
  version "0.1.0"
  license "UNLICENSED"

  on_macos do
    on_arm do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_darwin_arm64"
      sha256 "PLACEHOLDER_ARM64_DARWIN_SHA256"
    end
    on_intel do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_darwin_amd64"
      sha256 "PLACEHOLDER_AMD64_DARWIN_SHA256"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_linux_arm64"
      sha256 "PLACEHOLDER_ARM64_LINUX_SHA256"
    end
    on_intel do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_linux_amd64"
      sha256 "PLACEHOLDER_AMD64_LINUX_SHA256"
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
