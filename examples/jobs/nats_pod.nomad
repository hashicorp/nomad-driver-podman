# Copyright IBM Corp. 2019, 2025
# SPDX-License-Identifier: MPL-2.0

//
// This job runs a pair of nats-server and prometheus-nats-exporter tasks
// in the network namespace defined by a "pause" container.
//
// Because of that, all tasks will share a common network namespace and
// thus the exporter can access the localhost http monitoring port without exposing it
//
// A curl http://<your-ip>:7777/metrics hits the exporter which in turn grabs
// values from the nats-server
//
job "nats" {

  datacenters = ["dc1"]
  type        = "service"

  group "nats" {

    // Expose the default port for nats clients and 
    // also the exporters http endpoint
    network {
      port "server" { static = 4222 }
      port "exporter" { static = 7777 }
    }

    task "pod" {
      driver = "podman"

      // run the pause container before
      // the main workload and other sidecars
      lifecycle {
        hook    = "prestart"
        sidecar = "true"
      }

      config {
        image = "docker://k8s.gcr.io/pause:3.1"

        // the "pod" task must define the complete network
        // port mapping here
        ports = [
          "server",
          "exporter"
        ]
      }
    }

    task "server" {
      driver = "podman"

      // no "lifecycle" task: this is the main workload

      // server configuration file
      template {
        change_mode = "noop"
        destination = "local/nats-server.conf"
        data        = file("./templates/nats-server.conf.tpl")
      }

      config {
        image = "docker://nats:2.2.6"

        // here we join the pods network namespace
        network_mode = "task:pod"

        args = [
          "--config",
          "/local/nats-server.conf"
        ]

      }
    }

    task "exporter" {

      driver = "podman"

      // ensure to start the exporter _after_ the server
      lifecycle {
        hook    = "poststart"
        sidecar = "true"
      }

      config {
        image = "docker://natsio/prometheus-nats-exporter:0.7.0"

        // here we join the pods network namespace
        network_mode = "task:pod"

        // ... in order to access the private monitoring port
        args = [
          "-varz",
          "http://localhost:8222"
        ]
      }
    }

  }
}



