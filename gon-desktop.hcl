source = ["dist/loophole-desktop-macos_darwin_amd64/loophole-desktop"]
bundle_id = "cloud.loophole.cli"

apple_id {
  username = "@env:AC_USERNAME"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Main Development GmbH"
}

zip {
  output_path = "./dist/loophole-desktop_macos.zip"
}

dmg {
  output_path = "./dist/loophole-desktop_macos.dmg"
  volume_name = "Loophole Desktop"
}