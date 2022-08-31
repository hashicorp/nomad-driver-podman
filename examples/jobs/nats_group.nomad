//
// This job runs a pair of nats-server and prometheus-nats-exporter tasks
//
// A isolated network stack is defined at group level and both
// tasks will join it.
// 
// The exporter tasks can thus access the servers http monitoring port without
// exposing it.
//
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

      // define a shared network namespace for all tasks in this group
      mode = "bridge"

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

        ports = [
          "server"
        ]
      }
    }

    task "exporter" {

      driver = "podman"

      config {
        image = "docker://natsio/prometheus-nats-exporter:0.7.0"

        // NOTE: this example does not bind the servers monitoring port to localhost due
        //       nomad issue #10014
        args = [
          "-varz",
          "http://${NOMAD_IP_server}:8222"
        ]

        ports = [
          "exporter"
        ]

      }

    }

  }
}



