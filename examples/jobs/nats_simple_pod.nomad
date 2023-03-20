# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

//
// This job runs a pair of nats-server and prometheus-nats-exporter tasks
//
// The exporter tasks runs in the network namespace of the server task
// and can thus access the localhost http monitoring port without exposing it
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

        args = [
          "--config",
          "/local/nats-server.conf"
        ]

        // the "server" task must define the complete network
        // environment, so we will also pre-define the exporter
        // port mapping here
        ports = [
          "server",
          "exporter"
        ]
      }
    }

    task "exporter" {

      driver = "podman"

      // ensure to run the exporter _after_ the server
      // so that we can join the network namespace
      lifecycle {
        hook    = "poststart"
        sidecar = "true"
      }

      config {
        image = "docker://natsio/prometheus-nats-exporter:0.7.0"

        // here we join the servers network namespace
        network_mode = "task:server"

        // ... in order to access the private monitoring port
        args = [
          "-varz",
          "http://localhost:8222"
        ]
      }
    }

  }
}



