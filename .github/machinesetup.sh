#!/bin/bash -e

# add podman repository
echo "deb http://ppa.launchpad.net/projectatomic/ppa/ubuntu $(lsb_release -cs) main" > /etc/apt/sources.list.d/podman.list
apt-key adv --recv-keys --keyserver keyserver.ubuntu.com 018BA5AD9DF57A4448F0E6CF8BECF1637AD8C79D

# Ignore apt-get update errors to avoid failing due to misbehaving repo;
# true errors would fail in the apt-get install phase
apt-get update || true

# install podman for running the test suite
apt-get install -y podman wget ca-certificates

# get catatonit (to check podman --init switch)
cd /tmp
wget https://github.com/openSUSE/catatonit/releases/download/v0.1.4/catatonit.x86_64
mkdir -p /usr/libexec/podman
mv catatonit* /usr/libexec/podman/catatonit
chmod +x /usr/libexec/podman/catatonit

echo "====== Installed podman:"
# ensure to remember the used version when checking a build log
podman info

echo "====== Podman version:"
podman version

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
