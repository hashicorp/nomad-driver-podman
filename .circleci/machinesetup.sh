#!/bin/bash -e

# add podman repository
echo "deb http://ppa.launchpad.net/projectatomic/ppa/ubuntu $(lsb_release -cs) main" > /etc/apt/sources.list.d/podman.list
apt-key adv --recv-keys --keyserver keyserver.ubuntu.com 018BA5AD9DF57A4448F0E6CF8BECF1637AD8C79D

# Ignore apt-get update errors to avoid failing due to misbehaving repo;
# true errors would fail in the apt-get install phase
apt-get update || true

# install podman for running the test suite
apt-get install -y podman wget ca-certificates

# enable varlink socket (not included in ubuntu package)
cat > /etc/systemd/system/io.podman.service << EOF
[Unit]
Description=Podman Remote API Service
Requires=io.podman.socket
After=io.podman.socket
Documentation=man:podman-varlink(1)

[Service]
Type=simple
ExecStart=/usr/bin/podman varlink unix:%t/podman/io.podman --timeout=60000
TimeoutStopSec=30
KillMode=process

[Install]
WantedBy=multi-user.target
Also=io.podman.socket
EOF

cat > /etc/systemd/system/io.podman.socket << EOF
[Unit]
Description=Podman Remote API Socket
Documentation=man:podman-varlink(1)

[Socket]
ListenStream=%t/podman/io.podman
SocketMode=0600

[Install]
WantedBy=sockets.targett
EOF

systemctl daemon-reload
systemctl start io.podman


# remove default circleci go
cd /usr/local
rm -rf go

# setup go 1.12.x, instead
wget -O- https://storage.googleapis.com/golang/go1.12.13.linux-amd64.tar.gz| tar xfz -
ln -s /usr/local/go/bin/go /usr/bin/go