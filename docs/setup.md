# Setup Guide

## 1. Initialize Virtual Machines

| Machine   | IP Address        | RAM | Disk | User/Password   | Domain                |
|-----------|-------------------|-----|------|-----------------|-----------------------|
| runner    | 192.168.226.130   | 2GB | 20GB | ltkem/password  | runner.local.work     |
| gitlab    | 192.168.226.131   | 2GB | 30GB | ltkem/password  | gitlab.local.work     |
| database  | 192.168.226.132   | 2GB | 20GB | ltkem/password  | database.local.work   |
| jenkins   | 192.168.226.133   | 2GB | 20GB | ltkem/password  | jenkins.local.work    |
| dev       | 192.168.226.134   | 2GB | 20GB | ltkem/password  | dev.local.work        |

## 2. Update `/etc/hosts`

Add the following lines to your `/etc/hosts` file on all machines to map hostnames to IP addresses:

```text
192.168.226.130 runner.local.work
192.168.226.131 gitlab.local.work
192.168.226.132 database.local.work
192.168.226.133 jenkins.local.work
192.168.226.134 dev.local.work
```

## 3. Copy SSH Public Key

To enable passwordless SSH access, copy your public key to each machine:

```bash
ssh-copy-id -i ~/.ssh/id_ed25519.pub ltkem@192.168.226.130
```
Alternatively:
```bash
cat ~/.ssh/id_ed25519.pub | ssh ltkem@192.168.226.130 "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"
```

Repeat for each machine as needed.

## 4. Configure Network with Netplan

Example Netplan configuration (`/etc/netplan/50-cloud-init.yaml`):

```yaml
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

Apply the configuration:

```bash
sudo netplan apply
```

## 5. Run Ansible Playbooks

Use the following Ansible playbooks to automate setup tasks. Replace `hosts.ini` and playbook names as appropriate.

- Change hostnames:
  ```bash
  ansible-playbook -i hosts.ini change_hostname.yml -vv --ask-vault-pass --ask-become-pass
  ```

- Install Docker:
  ```bash
  ansible-playbook -i hosts.ini install_docker.yml -vv --ask-vault-pass --ask-become-pass
  ```

- Update `/etc/hosts` on all machines:
  ```bash
  ansible-playbook -i hosts.ini add_hosts.yml -vv --ask-vault-pass --ask-become-pass
  ```

- Install Jenkins:
  ```bash
  ansible-playbook -i hosts.ini install_jenkins.yml -vv --ask-vault-pass --ask-become-pass
  ```

- Start Jenkins agent:
  ```bash
  ansible-playbook -i hosts.ini start_jenkins_agent.yml -vv --ask-vault-pass --ask-become-pass
  ```

---

**Note:**  
- Replace IP addresses, usernames, and passwords as needed for your environment.
- Ensure you have the necessary privileges and Ansible Vault passwords to run the playbooks.
