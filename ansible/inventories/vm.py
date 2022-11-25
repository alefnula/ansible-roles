#!/usr/bin/env python3

import json
import getpass
import argparse
import subprocess


class VM:
    NODES = ["green", "blue", "red", "black"]

    def __init__(self, vm: dict):
        self.name = vm["name"]
        self.ipv4 = vm["ipv4"]

    @property
    def ip(self):
        if len(self.ipv4) > 0:
            return self.ipv4[0]
        return "127.0.0.1"

    @property
    def hostname(self) -> str:
        return f"{self.name}.test"

    @property
    def meta(self) -> dict:
        return {
            "ansible_host": self.ip,
            "ansible_user": "ubuntu",
            "color": self.name,
        }

    @classmethod
    def get_all(cls) -> list["VM"]:
        all_vms = json.loads(
            subprocess.check_output(
                [
                    "multipass",
                    "list",
                    "--format",
                    "json",
                ]
            ).decode("utf-8")
        )["list"]
        return [cls(vm) for vm in all_vms if vm["name"] in cls.NODES]


class Inventory:
    LOCALHOST_VARS = {
        "ansible_user": getpass.getuser(),
        "ansible_host": "127.0.0.1",
        "ansible_connection": "local",
        "ansible_python_interpreter": "/usr/local/bin/python3",
    }

    def __init__(self, vms: list[VM]):
        self.vms = vms

    @property
    def vars(self):
        return {
            "ansible_user": "ubuntu",
        }

    @staticmethod
    def _print(data: dict):
        print(json.dumps(data, indent=4))

    @staticmethod
    def _meta(vms: list[VM]) -> dict:
        return {"hostvars": {vm.hostname: vm.meta for vm in vms}}

    def list(self) -> None:
        green = []
        nodes = []
        for vm in self.vms:
            if vm.name == "green":
                green.append(vm)
            nodes.append(vm)

        hosts = {
            "all": {
                "hosts": [vm.hostname for vm in self.vms],
                "vars": self.vars,
            },
            "green": {
                "hosts": [vm.hostname for vm in green],
                "vars": self.vars,
            },
            "nodes": {
                "hosts": [vm.hostname for vm in nodes],
                "vars": self.vars,
            },
            "localhost": {
                "hosts": ["127.0.0.1"],
                "vars": self.LOCALHOST_VARS,
            },
        }
        hosts["_meta"] = self._meta(self.vms)
        self._print(hosts)

    def host(self, host: str) -> None:
        host = None
        for vm in self.vms:
            if host in (vm.hostname, vm.ip):
                host = vm
                break
        self._print(self._meta([] if host is None else [host]))


def create_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser()
    parser.add_argument("--list", action="store_true")
    parser.add_argument("--host", action="store")
    return parser


def main():
    parser = create_parser()
    args = parser.parse_args()
    inventory = Inventory(VM.get_all())
    if args.list:
        inventory.list()
    elif args.host:
        inventory.host(args.host)


if __name__ == "__main__":
    main()
