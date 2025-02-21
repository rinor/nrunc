// Copyright (c) 2023-2024, Nubificus LTD
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hypervisors

import (
	"strings"
	"syscall"

	"github.com/rinor/nrunc/pkg/unikontainers/unikernels"
)

const (
	CloudHypervisorVmm    VmmType = "cloud-hypervisor"
	CloudHypervisorBinary string  = "cloud-hypervisor"
	CHJsonFilename        string  = "ch.json"
)

type CloudHypervisor struct {
	binaryPath string
	binary     string
}

func (vmm *CloudHypervisor) Stop(_ string) error {
	return nil
}

func (vmm *CloudHypervisor) Ok() error {
	return nil
}
func (vmm *CloudHypervisor) Path() string {
	return vmm.binaryPath
}

func (vmm *CloudHypervisor) Execve(args ExecArgs, ukernel unikernels.Unikernel) error {
	vmmString := string(CloudHypervisorVmm)
	vmmMem := bytesToStringMB(args.MemSizeB)
	cmdString := vmm.binaryPath // + " --api-socket /tmp/ch_" + args.TapDevice + ".socket"
	cmdString += " --memory size=" + vmmMem + "M"
	cmdString += " --cpus boot=1"
	cmdString += " --rng src=/dev/urandom"
	cmdString += " --console off"
	cmdString += " --serial tty"

	if args.TapDevice != "" {
		netcli := ukernel.MonitorNetCli(vmmString)
		if netcli == "" {
			netcli += "--net tap=" + args.TapDevice + ",mac=" + args.GuestMAC + ",ip=,mask="
		}
		cmdString += " " + netcli
	}
	if args.BlockDevice != "" {
		blockCli := ukernel.MonitorBlockCli(vmmString)
		if blockCli == "" {
			blockCli += "--disk path=" + args.BlockDevice
		}
		cmdString += " " + blockCli

		kernel, err := ukernel.KernelFromBlock(args.BlockDevice, args.UnikernelPath)
		if err != nil {
			return err
		}
		if kernel != "" && kernel != args.UnikernelPath {
			args.UnikernelPath = kernel
		}
	}
	if args.UnikernelPath != "" {
		cmdString += " --kernel " + args.UnikernelPath
	}
	if args.InitrdPath != "" {
		cmdString += " --initramfs " + args.InitrdPath
	}
	cmdString = appendNonEmpty(cmdString, " ", ukernel.MonitorCli(vmmString))
	exArgs := strings.Split(cmdString, " ")
	if args.Command != "" {
		exArgs = append(exArgs, "--cmdline", args.Command)
	}
	vmmLog.WithField("cloud-hypervisor command", exArgs).Info("Ready to execve cloud-hypervisor")

	return syscall.Exec(vmm.Path(), exArgs, args.Environment) //nolint: gosec
}
