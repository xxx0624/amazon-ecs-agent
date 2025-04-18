//go:build unit
// +build unit

// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package engine

import (
	"testing"

	"github.com/aws/amazon-ecs-agent/agent/utils"
	commonutils "github.com/aws/amazon-ecs-agent/ecs-agent/utils"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/stretchr/testify/assert"
)

func getTestHostResourceManager(cpu int32, mem int32, ports []string, portsUdp []string, gpuIDs []string) *HostResourceManager {
	hostResources := make(map[string]types.Resource)
	hostResources["CPU"] = types.Resource{
		Name:         utils.Strptr("CPU"),
		Type:         utils.Strptr("INTEGER"),
		IntegerValue: cpu,
	}

	hostResources["MEMORY"] = types.Resource{
		Name:         utils.Strptr("MEMORY"),
		Type:         utils.Strptr("INTEGER"),
		IntegerValue: mem,
	}

	hostResources["PORTS_TCP"] = types.Resource{
		Name:           utils.Strptr("PORTS_TCP"),
		Type:           utils.Strptr("STRINGSET"),
		StringSetValue: ports,
	}

	hostResources["PORTS_UDP"] = types.Resource{
		Name:           utils.Strptr("PORTS_UDP"),
		Type:           utils.Strptr("STRINGSET"),
		StringSetValue: portsUdp,
	}

	hostResources["GPU"] = types.Resource{
		Name:           utils.Strptr("GPU"),
		Type:           utils.Strptr("STRINGSET"),
		StringSetValue: gpuIDs,
	}

	hostResourceManager := NewHostResourceManager(hostResources)

	return &hostResourceManager
}

func getTestTaskResourceMap(cpu int32, mem int32, ports []string, portsUdp []string, gpuIDs []string) map[string]types.Resource {
	taskResources := make(map[string]types.Resource)
	taskResources["CPU"] = types.Resource{
		Name:         utils.Strptr("CPU"),
		Type:         utils.Strptr("INTEGER"),
		IntegerValue: cpu,
	}

	taskResources["MEMORY"] = types.Resource{
		Name:         utils.Strptr("MEMORY"),
		Type:         utils.Strptr("INTEGER"),
		IntegerValue: mem,
	}

	taskResources["PORTS_TCP"] = types.Resource{
		Name:           utils.Strptr("PORTS_TCP"),
		Type:           utils.Strptr("STRINGSET"),
		StringSetValue: ports,
	}

	taskResources["PORTS_UDP"] = types.Resource{
		Name:           utils.Strptr("PORTS_UDP"),
		Type:           utils.Strptr("STRINGSET"),
		StringSetValue: portsUdp,
	}

	taskResources["GPU"] = types.Resource{
		Name:           utils.Strptr("GPU"),
		Type:           utils.Strptr("STRINGSET"),
		StringSetValue: gpuIDs,
	}

	return taskResources
}

func TestHostResourceConsumeSuccess(t *testing.T) {
	hostResourcePort1 := "22"
	hostResourcePort2 := "1000"
	gpuIDs := []string{"gpu1", "gpu2", "gpu3", "gpu4"}
	h := getTestHostResourceManager(int32(2048), int32(2048), []string{hostResourcePort1}, []string{hostResourcePort2}, gpuIDs)

	testTaskArn := "arn:aws:ecs:us-east-1:<aws_account_id>:task/cluster-name/11111"
	taskPort1 := "23"
	taskPort2 := "1001"
	taskGpuId1 := "gpu2"
	taskGpuId2 := "gpu3"
	taskResources := getTestTaskResourceMap(int32(512), int32(768), []string{taskPort1}, []string{taskPort2}, []string{taskGpuId1, taskGpuId2})

	consumed, _ := h.consume(testTaskArn, taskResources)
	assert.Equal(t, consumed, true, "Incorrect consumed status")
	assert.Equal(t, h.consumedResource["CPU"].IntegerValue, int32(512), "Incorrect cpu resource accounting during consume")
	assert.Equal(t, h.consumedResource["MEMORY"].IntegerValue, int32(768), "Incorrect memory resource accounting during consume")
	assert.Equal(t, h.consumedResource["PORTS_TCP"].StringSetValue[0], "22", "Incorrect port resource accounting during consume")
	assert.Equal(t, h.consumedResource["PORTS_TCP"].StringSetValue[1], "23", "Incorrect port resource accounting during consume")
	assert.Equal(t, len(h.consumedResource["PORTS_TCP"].StringSetValue), 2, "Incorrect port resource accounting during consume")
	assert.Equal(t, h.consumedResource["PORTS_UDP"].StringSetValue[0], "1000", "Incorrect udp port resource accounting during consume")
	assert.Equal(t, h.consumedResource["PORTS_UDP"].StringSetValue[1], "1001", "Incorrect udp port resource accounting during consume")
	assert.Equal(t, len(h.consumedResource["PORTS_UDP"].StringSetValue), 2, "Incorrect port resource accounting during consume")
	assert.Equal(t, h.consumedResource["GPU"].StringSetValue[0], "gpu2", "Incorrect gpu resource accounting during consume")
	assert.Equal(t, h.consumedResource["GPU"].StringSetValue[1], "gpu3", "Incorrect gpu resource accounting during consume")
	assert.Equal(t, len(h.consumedResource["GPU"].StringSetValue), 2, "Incorrect gpu resource accounting during consume")
}

func TestHostResourceConsumeFail(t *testing.T) {
	hostResourcePort1 := "22"
	hostResourcePort2 := "1000"
	gpuIDs := []string{"gpu1", "gpu2", "gpu3", "gpu4"}
	h := getTestHostResourceManager(int32(2048), int32(2048), []string{hostResourcePort1}, []string{hostResourcePort2}, gpuIDs)

	testTaskArn := "arn:aws:ecs:us-east-1:<aws_account_id>:task/cluster-name/11111"
	taskPort1 := "22"
	taskPort2 := "1001"
	taskGpuId1 := "gpu2"
	taskGpuId2 := "gpu3"
	taskResources := getTestTaskResourceMap(int32(512), int32(768), []string{taskPort1}, []string{taskPort2}, []string{taskGpuId1, taskGpuId2})

	consumed, _ := h.consume(testTaskArn, taskResources)
	assert.Equal(t, consumed, false, "Incorrect consumed status")
	assert.Equal(t, h.consumedResource["CPU"].IntegerValue, int32(0), "Incorrect cpu resource accounting during consume")
	assert.Equal(t, h.consumedResource["MEMORY"].IntegerValue, int32(0), "Incorrect memory resource accounting during consume")
	assert.Equal(t, h.consumedResource["PORTS_TCP"].StringSetValue[0], "22", "Incorrect port resource accounting during consume")
	assert.Equal(t, len(h.consumedResource["PORTS_TCP"].StringSetValue), 1, "Incorrect port resource accounting during consume")
	assert.Equal(t, h.consumedResource["PORTS_UDP"].StringSetValue[0], "1000", "Incorrect udp port resource accounting during consume")
	assert.Equal(t, len(h.consumedResource["PORTS_UDP"].StringSetValue), 1, "Incorrect port resource accounting during consume")
	assert.Equal(t, len(h.consumedResource["GPU"].StringSetValue), 0, "Incorrect gpu resource accounting during consume")
}

func TestHostResourceRelease(t *testing.T) {
	hostResourcePort1 := "22"
	hostResourcePort2 := "1000"
	gpuIDs := []string{"gpu1", "gpu2", "gpu3", "gpu4"}
	h := getTestHostResourceManager(int32(2048), int32(2048), []string{hostResourcePort1}, []string{hostResourcePort2}, gpuIDs)

	testTaskArn := "arn:aws:ecs:us-east-1:<aws_account_id>:task/cluster-name/11111"
	taskPort1 := "23"
	taskPort2 := "1001"
	taskGpuId1 := "gpu2"
	taskGpuId2 := "gpu3"
	taskResources := getTestTaskResourceMap(int32(512), int32(768), []string{taskPort1}, []string{taskPort2}, []string{taskGpuId1, taskGpuId2})

	h.consume(testTaskArn, taskResources)
	h.release(testTaskArn, taskResources)

	assert.Equal(t, h.consumedResource["CPU"].IntegerValue, int32(0), "Incorrect cpu resource accounting during release")
	assert.Equal(t, h.consumedResource["MEMORY"].IntegerValue, int32(0), "Incorrect memory resource accounting during release")
	assert.Equal(t, h.consumedResource["PORTS_TCP"].StringSetValue[0], "22", "Incorrect port resource accounting during release")
	assert.Equal(t, len(h.consumedResource["PORTS_TCP"].StringSetValue), 1, "Incorrect port resource accounting during release")
	assert.Equal(t, h.consumedResource["PORTS_UDP"].StringSetValue[0], "1000", "Incorrect udp port resource accounting during release")
	assert.Equal(t, len(h.consumedResource["PORTS_UDP"].StringSetValue), 1, "Incorrect udp port resource accounting during release")
	assert.Equal(t, len(h.consumedResource["GPU"].StringSetValue), 0, "Incorrect gpu resource accounting during release")
}

func TestConsumable(t *testing.T) {
	testCases := []struct {
		name                       string
		cpu                        int32
		mem                        int32
		ports                      []uint16
		portsUdp                   []uint16
		gpus                       []string
		canBeConsumed              bool
		expectedFailedResourceKeys []string
	}{
		{
			name:                       "consumable",
			cpu:                        int32(1024),
			mem:                        int32(1024),
			ports:                      []uint16{25},
			portsUdp:                   []uint16{1003},
			gpus:                       []string{"gpu1", "gpu2"},
			canBeConsumed:              true,
			expectedFailedResourceKeys: nil,
		},
		{
			name:                       "cpu not consumable",
			cpu:                        int32(2500),
			mem:                        int32(1024),
			ports:                      []uint16{},
			portsUdp:                   []uint16{},
			gpus:                       []string{},
			canBeConsumed:              false,
			expectedFailedResourceKeys: []string{"CPU"},
		},
		{
			name:                       "memory not consumable",
			cpu:                        int32(1024),
			mem:                        int32(2500),
			ports:                      []uint16{},
			portsUdp:                   []uint16{},
			gpus:                       []string{},
			canBeConsumed:              false,
			expectedFailedResourceKeys: []string{"MEMORY"},
		},
		{
			name:                       "tcp ports not consumable",
			cpu:                        int32(1024),
			mem:                        int32(1024),
			ports:                      []uint16{22},
			portsUdp:                   []uint16{},
			gpus:                       []string{},
			canBeConsumed:              false,
			expectedFailedResourceKeys: []string{"PORTS_TCP"},
		},
		{
			name:                       "udp ports not consumable",
			cpu:                        int32(1024),
			mem:                        int32(1024),
			ports:                      []uint16{},
			portsUdp:                   []uint16{1000},
			gpus:                       []string{},
			canBeConsumed:              false,
			expectedFailedResourceKeys: []string{"PORTS_UDP"},
		},
		{
			name:                       "multiple resources not consumable - cpu and udp ports",
			cpu:                        int32(2500),
			mem:                        int32(1024),
			ports:                      []uint16{},
			portsUdp:                   []uint16{1000},
			gpus:                       []string{},
			canBeConsumed:              false,
			expectedFailedResourceKeys: []string{"CPU", "PORTS_UDP"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hostResourcePort1 := "22"
			hostResourcePort2 := "1000"
			gpuIDs := []string{"gpu1", "gpu2", "gpu3", "gpu4"}
			h := getTestHostResourceManager(int32(2048), int32(2048), []string{hostResourcePort1},
				[]string{hostResourcePort2}, gpuIDs)

			resources := getTestTaskResourceMap(tc.cpu, tc.mem, commonutils.Uint16SliceToStringSlice(tc.ports),
				commonutils.Uint16SliceToStringSlice(tc.portsUdp), tc.gpus)
			canBeConsumed, failedResourceKeys, err := h.consumable(resources)
			assert.Equal(t, tc.canBeConsumed, canBeConsumed,
				"Error in checking if resources can be successfully consumed")
			assert.Equal(t, nil, err,
				"Error in checking if resources can be successfully consumed, error returned from consumable")
			assert.ElementsMatch(t, tc.expectedFailedResourceKeys, failedResourceKeys)
		})
	}
}

func TestResourceHealthTrue(t *testing.T) {
	hostResourcePort1 := "22"
	hostResourcePort2 := "1000"
	gpuIDs := []string{"gpu1", "gpu2", "gpu3", "gpu4"}
	h := getTestHostResourceManager(int32(2048), int32(2048), []string{hostResourcePort1}, []string{hostResourcePort2}, gpuIDs)

	resources := getTestTaskResourceMap(1024, 1024, commonutils.Uint16SliceToStringSlice([]uint16{22}), commonutils.Uint16SliceToStringSlice([]uint16{1000}), []string{"gpu1", "gpu2"})
	err := h.checkResourcesHealth(resources)
	assert.NoError(t, err, "Error in checking healthy resource map status")
}

// Verify Resource health status checks gpu status properly from valid pool of gpus and returns error
func TestResourceHealthGPUFalse(t *testing.T) {
	hostResourcePort1 := "22"
	hostResourcePort2 := "1000"
	gpuIDs := []string{"gpu1", "gpu2", "gpu3", "gpu4"}
	h := getTestHostResourceManager(int32(2048), int32(2048), []string{hostResourcePort1}, []string{hostResourcePort2}, gpuIDs)

	resources := getTestTaskResourceMap(1024, 1024, commonutils.Uint16SliceToStringSlice([]uint16{22}), commonutils.Uint16SliceToStringSlice([]uint16{1000}), []string{"gpu1", "gpu5"})
	err := h.checkResourcesHealth(resources)
	assert.Error(t, err, "Error in checking unhealthy resource map status")
}
