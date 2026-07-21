class Cadreen < Formula
  desc "Intelligence infrastructure for developers"
  homepage "https://accomplishanything.today/infra"
  version "0.4.2"
  license "UNLICENSED"

  on_macos do
    on_arm do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_darwin_arm64"
      sha256 "1bb75dc7ce83b82767629adb8087d654b2a4abfcbdafc8e542507be57aeb9c1a"
    end
    on_intel do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_darwin_amd64"
      sha256 "974ecb60e419c2f04ba2992dc03e21062f339c066d4953d2c9b283b6ea977fe9"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_linux_arm64"
      sha256 "a2753aef2437b1e618415a2063934bc5d1eebf1bab39fa07b7b21a977fd656b2"
    end
    on_intel do
      url "https://github.com/timothy-billingrails/cadreen-sdks/releases/download/cli-v#{VERSION}/cadreen_linux_amd64"
      sha256 "151ea83c81a15c95e68ca385aed3996507005e272ec5732306f2246fd183f5a1"
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
