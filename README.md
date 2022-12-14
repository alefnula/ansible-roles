# K3S Utils

## Setup project

Create a `vars/secrets.yml` file with the same variables as in the `vars/secrets.yaml.template`.

```bash
# Encrypt the file
ansible-vault encrypt vars/secrets.yml

# Create a .passwd file and put your encryption password in it.
# You can also choose not to do this and enter the ansible vault password
# every time you issue a command.
touch .passwd
vim .passwd
```

## Setup Virtual Machines

### Setup multipass

Set multipass local driver to `qemu`:

```
multipass set local.driver=qemu
```

### Running the instances

```
mage vm:init
mage vm:setup
mage k3s:up
mage longhor:up
mage postgres:up
```

## Teardown Virtual Machines

```
mage postgres:down   # optional
mage longhorn:down   # optional
mage k3s:down        # optional
mage vm:destroy
```
