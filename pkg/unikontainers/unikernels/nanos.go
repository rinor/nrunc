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

package unikernels

import (
	"github.com/nanovms/ops/fs"
)

const (
	NanosUnikernel string = "nanos"
	NanosKernel    string = "kernel"
)

type nanos struct {
	cmdline string
	net     struct {
		address string
		gateway string
		mask    string
	}
}

func (n *nanos) Init(data UnikernelParams) error {
	if data.EthDeviceIP != "" {
		n.net.address = "en1.ipaddr=" + data.EthDeviceIP
	}
	if data.EthDeviceGateway != "" {
		n.net.gateway = "en1.gateway=" + data.EthDeviceGateway
	}
	if data.EthDeviceMask != "" {
		n.net.mask = "en1.netmask=" + data.EthDeviceMask
	}
	if data.CmdLine != "" {
		n.cmdline = data.CmdLine
	}

	return nil
}

func (n *nanos) CommandString() (string, error) {
	var command string
	if n.net.address != "" {
		command += n.net.address
	}
	if n.net.gateway != "" {
		command += " " + n.net.gateway
	}
	if n.net.mask != "" {
		command += " " + n.net.mask
	}
	if n.cmdline != "" {
		command += " " + n.cmdline
	}

	return command, nil
}

func (n *nanos) SupportsBlock() bool {
	return true
}

func (n *nanos) SupportsFS(_ string) bool {
	return false
}

func (n *nanos) MonitorNetCli(_ string) string {
	return ""
}

func (n *nanos) MonitorBlockCli(_ string) string {
	return ""
}

func (n *nanos) MonitorCli(_ string) string {
	return ""
}

func (n *nanos) KernelFromBlock(imagePath string, kernelDst string) (string, error) {
	if kernelDst == "" {
		kernelDst = NanosUnikernel + "/" + NanosKernel
	}

	bootfs, err := fs.NewReaderBootFS(imagePath)
	if err != nil {
		return "", err
	}
	defer bootfs.Close()

	err = bootfs.CopyFile(NanosKernel, kernelDst, false)
	if err != nil {
		return "", err
	}

	return kernelDst, nil
}

func newNanos() *nanos {
	return &nanos{}
}
