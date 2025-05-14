#!/bin/bash
set -e

if [ -f /etc/gitlab/gitlab.rb ]; then
    echo "Gitlab already installed, Skip..."
    exit 0
fi

apt update
apt install -y curl openssh-server ca-certificates tzdata perl
curl -s https://packages.gitlab.com/install/repositories/gitlab/gitlab-ee/script.deb.sh | bash
apt install -y gitlab-ee=17.9.8-ee.0
gitlab-ctl reconfigure
