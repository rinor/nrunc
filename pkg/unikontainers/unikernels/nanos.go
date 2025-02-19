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
	"errors"
	"fmt"

	"github.com/nanovms/ops/fs"
)

const NanosUnikernel string = "nanos"
const NanosKernel string = "kernel"

var ErrUndefinedVersion = errors.New("version is undefined, using default version")
var ErrVersionParsing = errors.New("failed to parse provided version, using default version")

type NanosNet struct {
	Address string
	Mask    string
	Gateway string
}

type NanosBlock struct {
	RootFS string
}

type Nanos struct {
	Command string
	Net     NanosNet
	Block   NanosBlock
}

func (n *Nanos) Init(data UnikernelParams) error {
	n.Net.Address = "en1.ipaddr=" + data.EthDeviceIP
	n.Net.Gateway = "en1.gateway=" + data.EthDeviceGateway
	n.Net.Mask = "en1.netmask=" + data.EthDeviceMask
	return nil
}

func (n *Nanos) CommandString() (string, error) {
	return fmt.Sprintf("%s %s %s %s",
		n.Net.Address,
		n.Net.Gateway,
		n.Net.Mask,
		n.Command), nil
}

func (n *Nanos) SupportsBlock() bool {
	return true
}

func (n *Nanos) SupportsFS(_ string) bool {
	return false
}

func (n *Nanos) MonitorNetCli(_ string) string {
	return ""
}

func (n *Nanos) MonitorBlockCli(_ string) string {
	return ""
}

func (n *Nanos) MonitorCli(_ string) string {
	return ""
}

func (n *Nanos) KernelFromBlock(imagePath string, kernelDst string) (string, error) {
	if kernelDst == "" {
		kernelDst = NanosUnikernel + "/" + NanosKernel
	}

	bootfs, err := fs.NewReaderBootFS(imagePath)
	if err != nil {
		return "", err
	}
	defer bootfs.Close()

	return kernelDst, bootfs.CopyFile(NanosKernel, kernelDst, false)
}

func newNanos() *Nanos {
	return &Nanos{}
}
