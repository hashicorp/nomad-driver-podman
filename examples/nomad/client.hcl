# Increase log verbosity
log_level = "DEBUG"

# Setup data dir
data_dir = "/tmp/podmanclient"

# Give the agent a unique name. Defaults to hostname
name = "podmanclient"

plugin_dir = "/home/vagrant/nomad-driver-podman/examples/plugins"

client {
  enabled = true
  servers = ["127.0.0.1:4647"]
}

plugin "nomad-driver-podman" {
  config {
    volumes {
      enabled      = true
      selinuxlabel = "z"
    }
  }
}

plugin "raw_exec" {
  config {
    enabled = true
  }
}

# different port than server
ports {
  http = 7646
}

