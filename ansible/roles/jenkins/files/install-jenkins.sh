#!/bin/bash

if systemctl list-units --type=service | grep -q jenkins; then
  echo "Jenkins service exists, Skip..."
  exit 0
fi

apt update
apt install fontconfig openjdk-21-jre -y
wget -O /etc/apt/keyrings/jenkins-keyring.asc \
  https://pkg.jenkins.io/debian-stable/jenkins.io-2023.key
echo "deb [signed-by=/etc/apt/keyrings/jenkins-keyring.asc]" \
  https://pkg.jenkins.io/debian-stable binary/ | tee \
  /etc/apt/sources.list.d/jenkins.list > /dev/null
apt install jenkins -y
systemctl start jenkins
systemctl enable jenkins
