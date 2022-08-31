#!/bin/bash -e

# Update ca-certificates first to prevent using outdated certificates for the
# podman repository.
# https://github.com/containers/podman/issues/8533#issuecomment-944873690
apt-get update
apt-get install -y ca-certificates

# add podman repository
echo "deb https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_$(lsb_release -rs)/ /" | tee /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list
curl -L https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_$(lsb_release -rs)/Release.key | apt-key add -

# Ignore apt-get update errors to avoid failing due to misbehaving repo;
# true errors would fail in the apt-get install phase
apt-get update || true

# install podman for running the test suite
apt-get install -y podman wget build-essential

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

# enable http socket
cat > /etc/systemd/system/podman.service << EOF
[Unit]
Description=Podman API Service
Requires=podman.socket
After=podman.socket
Documentation=man:podman-system-service(1)
StartLimitIntervalSec=0

[Service]
Type=simple
ExecStart=/usr/bin/podman system service

[Install]
WantedBy=multi-user.target
Also=podman.socket
EOF

cat > /etc/systemd/system/podman.socket << EOF
[Unit]
Description=Podman API Socket
Documentation=man:podman-system-service(1)

[Socket]
ListenStream=%t/podman/podman.sock
SocketMode=0660

[Install]
WantedBy=sockets.target
EOF

systemctl daemon-reload
# enable http api
systemctl start podman
