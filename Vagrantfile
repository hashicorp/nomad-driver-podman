
# Specify minimum Vagrant version and Vagrant API version
Vagrant.require_version ">= 1.6.0"
VAGRANTFILE_API_VERSION = "2"

# Create box
Vagrant.configure("2") do |config|
  config.vm.define "podman-linux"
  config.vm.box = "hashicorp/bionic64"
  config.vm.synced_folder ".", "/home/vagrant/nomad-driver-podman"
  config.ssh.extra_args = ["-t", "cd /home/vagrant/nomad-driver-podman; bash --login"]
  config.vm.network "forwarded_port", guest: 4646, host: 4646, host_ip: "127.0.0.1"
  config.vm.provider "virtualbox" do |vb|
      vb.name = "podman-linux"
      vb.cpus = 2
      vb.memory = 2048
  end

  config.vm.provision "shell" do |p|
    p.name = "machinesetup.sh"
    p.path = ".github/machinesetup.sh"
  end

  config.vm.provision "shell" do |p|
    p.name = "Install go"
    p.privileged = false
    p.inline = <<-SHELL
      go_version=1.14.12
      if [ ! -f /usr/local/go/bin/go ]; then
        curl -sSL -o go.tgz "https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz"
        sudo tar -C /usr/local -xzf go.tgz
        rm -f go.tgz
        echo "export PATH=/usr/local/go/bin:\$PATH" >> $HOME/.bashprofile
      fi

      . $HOME/.bashprofile

      if ! command -v gcc >/dev/null; then
        apt-get install -y gcc
      fi

      cd nomad-driver-podman
      make deps
    SHELL
  end
end
