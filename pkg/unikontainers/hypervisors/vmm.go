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
	"errors"
	"fmt"
	"os/exec"

	"github.com/rinor/nrunc/pkg/unikontainers/unikernels"
	"github.com/sirupsen/logrus"
)

const DefaultMemory uint64 = 256 // The default memory for every hypervisor: 256 MB

// ExecArgs holds the data required by Execve to start the VMM
// FIXME: add extra fields if required by additional VMM's
type ExecArgs struct {
	Container     string   // The container ID
	UnikernelPath string   // The path of the unikernel inside rootfs
	TapDevice     string   // The TAP device name
	BlockDevice   string   // The block device path
	InitrdPath    string   // The path to the initrd of the unikernel
	Command       string   // The unikernel's command line
	IPAddress     string   // The IP address of the TAP device
	GuestMAC      string   // The MAC address of the guest network device
	Seccomp       bool     // Enable or disable seccomp filters for the VMM
	MemSizeB      uint64   // The size of the memory provided to the VM in bytes
	Environment   []string // Environment
}

type VmmType string

var ErrVMMNotInstalled = errors.New("vmm not found")
var vmmLog = logrus.WithField("subsystem", "hypervisors")

type VMM interface {
	Execve(args ExecArgs, ukernel unikernels.Unikernel) error
	Stop(t string) error
	Path() string
	Ok() error
}

func NewVMM(vmmType VmmType) (vmm VMM, err error) {
	defer func() {
		if err != nil {
			vmmLog.Error(err.Error())
		}
	}()
	switch vmmType {
	case QemuVmm:
		vmmPath, err := exec.LookPath(QemuBinary + cpuArch())
		if err != nil {
			return nil, ErrVMMNotInstalled
		}
		return &Qemu{binary: QemuBinary, binaryPath: vmmPath}, nil
	case FirecrackerVmm:
		vmmPath, err := exec.LookPath(FirecrackerBinary)
		if err != nil {
			return nil, ErrVMMNotInstalled
		}
		return &Firecracker{binary: FirecrackerBinary, binaryPath: vmmPath}, nil
	case CloudHypervisorVmm:
		vmmPath, err := exec.LookPath(CloudHypervisorBinary)
		if err != nil {
			return nil, ErrVMMNotInstalled
		}
		return &CloudHypervisor{binary: CloudHypervisorBinary, binaryPath: vmmPath}, nil
	default:
		return nil, fmt.Errorf("vmm \"%s\" is not supported", vmmType)
	}
}
