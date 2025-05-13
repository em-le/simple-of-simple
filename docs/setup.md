# Guide

1. Init Virtual Machines

| Machine    | IP              | RAM | Disk | User/Password     | Domain                |
|------------|-----------------|-----|------|-------------------|-----------------------|
| runner     | 192.168.226.130 | 2GB | 20GB | ltkem/password    | runner.local.work     |
| gitlab     | 192.168.226.131 | 2GB | 30GB | ltkem/password    | gitlab.local.work     |
| database   | 192.168.226.132 | 2GB | 20GB | ltkem/password    | database.local.work   |
| jenkins    | 192.168.226.133 | 2GB | 20GB | ltkem/password    | jenkins.local.work    |
| dev        | 192.168.226.134 | 2GB | 20GB | ltkem/password    | dev.local.work        |

```/etc/hosts
192.168.226.130 runner.local.work
192.168.226.131 gitlab.local.work
192.168.226.132 database.local.work
192.168.226.133 jenkins.local.work
192.168.226.134 dev.local.work
```

Copying Public Key Using SSH
```bash
ssh-copy-id -i ~/.ssh/id_ed25519.pub ltkem@192.168.226.130
```
```bash
cat ~/.ssh/id_ed25519.pub | ssh ltkem@192.168.226.130 "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"
```

NetPlan config:
```/etc/netplan/50-cloud-init.yaml 
network:
  version: 2
  ethernets:
    ens160:
      dhcp4: no
      addresses: [192.168.226.132/24]
      gateway4: 192.168.226.2
      nameservers:
        addresses: [8.8.8.8, 8.8.4.4]
```
```bash
netplan apply
```

```bash
ansible-playbook -i hosts.ini change_hostname.yml --ask-vault-pass --ask-become-pass
```

```bash
ansible-playbook -i hosts.ini install_docker.yml -vvv --ask-vault-pass --ask-become-pass
```

```bash
ansible-playbook -i hosts.ini add_hosts.yml --ask-vault-pass --ask-become-pass
```