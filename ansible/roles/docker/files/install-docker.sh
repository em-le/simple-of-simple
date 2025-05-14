#!/usr/bin/env bash

if command -v docker &> /dev/null; then
  echo "Docker already installed, Skip..."
  exit 0
fi

echo "Installing Docker..."
set -e

apt-get update
apt-get -qq install ca-certificates curl gnupg lsb-release

mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update
apt-get -qq install docker-ce docker-ce-cli docker-ce-rootless-extras containerd.io docker-buildx-plugin docker-compose-plugin cgroupfs-mount libltdl7 pigz

if [ ! -f "/bin/docker-compose" ]; then
  echo 'docker compose "$@"' | tee /bin/docker-compose
  chmod +x /bin/docker-compose
fi

docker run hello-world
