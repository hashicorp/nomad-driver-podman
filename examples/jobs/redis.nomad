job "redis" {
  datacenters = ["dc1"]
  type        = "service"

  group "cache" {
    task "redis" {
      driver = "podman"

      env {
        foo = "bar"
      }
      config {
        port_map {
          redis = 6379
        }

        image = "docker://redis"
      }
      resources {
        network {
          port "redis" {}
        }
      }
    }
  }
}

