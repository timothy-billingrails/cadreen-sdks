class Cadreen < Formula
  desc "Intelligence infrastructure for developers"
  homepage "https://accomplishanything.today"
  version "0.1.0"
  license "UNLICENSED"

  on_macos do
    on_arm do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_darwin_arm64"
      sha256 "8f596a58fa95928d7e9a29d2a56eafde3e578c1a74c85d56d277d517a093bb57"
    end
    on_intel do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_darwin_amd64"
      sha256 "4a673c5eb913433f3b0fb7fc7c19bfefda14fb093dfc40f3140c565774d61a22"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_linux_arm64"
      sha256 "da50951b2ddeefe87342db8e4642bbc56d378f88810d16e03d47ea36ad655df7"
    end
    on_intel do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_linux_amd64"
      sha256 "630e9e9b9929bda40853712d420808c2526b2a0435e4a774bfaa5b4359e9e264"
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
