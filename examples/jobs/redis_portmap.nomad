job "redis" {
  datacenters = ["dc1"]
  type        = "service"

  group "cache" {
    network {
      port "redis" { to = 6379 }
    }
    task "redis" {
      driver = "podman"

      env {
        foo = "bar"
      }
      config {
        image = "docker://redis"
        ports = ["redis"]
      }
    }
  }
}


