# Client port of 4222 on all interfaces
port: 4222

# Allow ownly localhost access to the monitoring pot
http: localhost:8222

# This is for clustering multiple servers together.
cluster {
  # It is recommended to set a cluster name
  name: "my_cluster"

  # Route connections to be received on any interface on port 6222
  port: 6222

}
