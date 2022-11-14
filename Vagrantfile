
# Specify minimum Vagrant version and Vagrant API version
Vagrant.require_version ">= 2.2.0"
VAGRANTFILE_API_VERSION = "2"

# Create box
Vagrant.configure("2") do |config|
  config.vm.define "podman-linux"
  config.vm.box = "ubuntu/focal64"
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
    p.privileged = false
    p.inline = <<-SHELL

      go_version=$(cat /home/vagrant/nomad-driver-podman/.go-version)

      if [ ! -f /usr/local/go/bin/go ]; then
        curl -sSL -o go.tgz "https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz"
        sudo tar -C /usr/local -xzf go.tgz
        rm -f go.tgz
        export PATH=/usr/local/go/bin:$PATH
        echo "export PATH=/usr/local/go/bin:\$PATH" >> $HOME/.bash_profile
      fi

      if ! command -v unzip >/dev/null; then
        sudo apt-get install -y unzip
      fi

      nomad_version=1.3.1
      if [ ! -f /usr/bin/nomad ]; then
        curl -o nomad.zip -sSL https://releases.hashicorp.com/nomad/${nomad_version}/nomad_${nomad_version}_linux_amd64.zip
        sudo unzip nomad.zip -d /usr/bin
        rm -f nomad.zip
      fi

      . $HOME/.bash_profile

      if ! command -v gcc >/dev/null; then
        sudo apt-get install -y gcc
      fi

      cd nomad-driver-podman
      make deps
    SHELL
  end
end
