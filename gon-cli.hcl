source = ["dist/loophole-cli-macos_darwin_amd64/loophole"]
bundle_id = "cloud.loophole.cli"

apple_id {
  username = "@env:AC_USERNAME"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Main Development GmbH"
}

zip {
  output_path = "./dist/loophole-cli_macos.zip"
}

dmg {
  output_path = "./dist/loophole-cli_macos.dmg"
  volume_name = "Loophole"
}