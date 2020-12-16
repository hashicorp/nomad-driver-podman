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
        image = "docker://redis"
      }
    }
  }
}

