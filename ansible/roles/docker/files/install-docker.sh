#!/usr/bin/env bash

echo "Installing Docker..."
set -e

sudo apt-get update
sudo apt-get -qq install ca-certificates curl gnupg lsb-release

sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update
sudo apt-get -qq install docker-ce docker-ce-cli docker-ce-rootless-extras containerd.io docker-buildx-plugin docker-compose-plugin cgroupfs-mount libltdl7 pigz

if [ ! -f "/bin/docker-compose" ]; then
  echo 'docker compose "$@"' | sudo tee /bin/docker-compose
  sudo chmod +x /bin/docker-compose
fi

sudo docker run hello-world