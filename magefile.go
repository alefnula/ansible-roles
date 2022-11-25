//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

/// ! Ansible helpers

var VaultArgs []string

func init() {
	// Set colors
	os.Setenv("MAGEFILE_ENABLE_COLOR", "true")

	if _, err := os.Stat(".passwd"); err == nil {
		VaultArgs = []string{"--vault-password-file", ".passwd"}
	} else {
		VaultArgs = []string{"--ask-vault-password"}
	}
}

type Ansible struct {
	Tags string
}

// Run Ansible
func (a Ansible) Run() error {
	// extraVarsBytes, _ := json.Marshal(map[string]string{
	// 	"cluster": a.Cluster,
	// })

	// Create arguments
	args := []string{
		fmt.Sprintf("ansible/main.yml"),
		"--inventory",
		fmt.Sprintf("ansible/inventories/vm.py"),
		// "--extra-vars",
		// string(extraVarsBytes),
	}

	args = append(args, VaultArgs...)

	if a.Tags != "" {
		args = append(args, "--tags", a.Tags)
	}

	return sh.RunV("ansible-playbook", args...)
}

// ! Help

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Help

// : Show help.
func Help() error {
	return sh.RunV("mage", "-l")
}

// ! VM

type VM mg.Namespace

// : Bootstrap development environment.
func (VM) Init() error { return Ansible{Tags: "init"}.Run() }

// : Starts VM cluster.
func (VM) Up() error { return sh.RunV("multipass", "start", "--all") }

// : Stops VM cluster.
func (VM) Down() error { return sh.RunV("multipass", "stop", "--all") }

// : Prints VM status.
func (VM) Status() error { return sh.RunV("multipass", "list") }

// : Destroy VM cluster.
func (VM) Destroy() error {
	if err := sh.RunV("multipass", "stop", "--all"); err != nil {
		return err
	}
	if err := sh.RunV("multipass", "delete", "--all"); err != nil {
		return err
	}
	return sh.RunV("multipass", "purge")
}

// ! Setup

// : Setup cluster.
func (VM) Setup() error { return Ansible{Tags: "setup"}.Run() }

// ! K3s
type K3s mg.Namespace

// : K3S setup.
func (K3s) Up() error { return Ansible{Tags: "k3s-up"}.Run() }

// : K3S teardown.
func (K3s) Down() error { return Ansible{Tags: "k3s-down"}.Run() }

// ! Longhorn
type Longhorn mg.Namespace

// : Longhorn setup.
func (Longhorn) Up() error { return Ansible{Tags: "longhorn-up"}.Run() }

// : Longhorn teardown.
func (Longhorn) Down() error { return Ansible{Tags: "longhorn-down"}.Run() }

// ! Postgres
type Postgres mg.Namespace

// : Longhorn setup.
func (Postgres) Up() error { return Ansible{Tags: "postgres-up"}.Run() }

// : Longhorn teardown.
func (Postgres) Down() error { return Ansible{Tags: "postgres-down"}.Run() }
