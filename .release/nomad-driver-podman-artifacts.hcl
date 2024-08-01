schema = 1
artifacts {
  zip = [
    "nomad-driver-podman_${version}_linux_amd64.zip",
    "nomad-driver-podman_${version}_linux_arm64.zip",
  ]
  rpm = [
    "nomad-driver-podman-${version_linux}-1.aarch64.rpm",
    "nomad-driver-podman-${version_linux}-1.x86_64.rpm",
  ]
  deb = [
    "nomad-driver-podman_${version_linux}-1_amd64.deb",
    "nomad-driver-podman_${version_linux}-1_arm64.deb",
  ]
}
