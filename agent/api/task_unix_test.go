// +build !windows

// Copyright 2014-2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package api

import (
	"testing"

	"github.com/aws/amazon-ecs-agent/agent/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyVolumeName1                  = "empty-volume-1"
	emptyVolumeContainerPath1         = "/my/empty-volume-1"
	expectedEmptyVolumeGeneratedPath1 = "/ecs-empty-volume/" + emptyVolumeName1

	emptyVolumeName2                  = "empty-volume-2"
	emptyVolumeContainerPath2         = "/my/empty-volume-2"
	expectedEmptyVolumeGeneratedPath2 = "/ecs-empty-volume/" + emptyVolumeName2

	expectedEmptyVolumeContainerImage = "amazon/ecs-emptyvolume-base"
	expectedEmptyVolumeContainerTag   = "autogenerated"
	expectedEmptyVolumeContainerCmd   = "not-applicable"
)

func TestAddNetworkResourceProvisioningDependencyNop(t *testing.T) {
	testTask := &Task{
		Containers: []*Container{
			{
				Name: "c1",
			},
		},
	}
	testTask.addNetworkResourceProvisioningDependency(nil)
	assert.Equal(t, 1, len(testTask.Containers))
}

func TestAddNetworkResourceProvisioningDependencyWithENI(t *testing.T) {
	testTask := &Task{
		ENI: &ENI{},
		Containers: []*Container{
			{
				Name: "c1",
			},
		},
	}
	cfg := &config.Config{
		PauseContainerImageName: "pause-container-image-name",
		PauseContainerTag:       "pause-container-tag",
	}
	testTask.addNetworkResourceProvisioningDependency(cfg)
	assert.Equal(t, 2, len(testTask.Containers),
		"addNetworkResourceProvisioningDependency should add another container")
	pauseContainer, ok := testTask.ContainerByName(PauseContainerName)
	require.True(t, ok, "Expected to find pause container")
	assert.Equal(t, ContainerCNIPause, pauseContainer.Type, "pause container should have correct type")
	assert.True(t, pauseContainer.Essential, "pause container should be essential")
	assert.Equal(t, cfg.PauseContainerImageName+":"+cfg.PauseContainerTag, pauseContainer.Image,
		"pause container should use configured image")
}
