
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
  config.vm.provision "shell", inline: <<-SHELL
    apt-get update
    apt-get install -y unzip gcc runc podman
    source /home/vagrant/.bashrc
    # Install golang-1.15.6
    if [ ! -f "/usr/local/go/bin/go" ]; then
      curl -s -L -o go1.15.6.linux-amd64.tar.gz https://dl.google.com/go/go1.15.6.linux-amd64.tar.gz
      sudo tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz
      rm -f go1.15.6.linux-amd64.tar.gz
    fi
    # Install nomad-1.0.0
    if [ ! -f "/usr/bin/nomad" ]; then
      wget --quiet https://releases.hashicorp.com/nomad/1.0.0/nomad_1.0.0_linux_amd64.zip
      unzip nomad_1.0.0_linux_amd64.zip -d /usr/bin
      rm -f nomad_1.0.0_linux_amd64.zip
    fi
    # Run setup
    cd /home/vagrant/go/src/nomad-driver-containerd/vagrant
    ./setup.sh
  SHELL
end
