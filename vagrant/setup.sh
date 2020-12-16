#!/bin/bash -e

cd /tmp
wget https://github.com/openSUSE/catatonit/releases/download/v0.1.4/catatonit.x86_64
mkdir -p /usr/libexec/podman
mv catatonit* /usr/libexec/podman/catatonit
chmod +x /usr/libexec/podman/catatonit

# enable http socket (not included in ubuntu package)
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
