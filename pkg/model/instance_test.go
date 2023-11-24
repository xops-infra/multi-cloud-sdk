package model

import (
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	instances := []*Instance{
		{
			Name:      tea.String("test"),
			PrivateIP: []*string{tea.String("10.1.1.1")},
			PublicIP:  []*string{tea.String("1.1.1.1")},
			Profile:   "tencent",
			Status:    InstanceStatusPending,
		},
		{
			Name:      tea.String("test1"),
			PrivateIP: []*string{tea.String("10.1.1.2")},
			PublicIP:  []*string{tea.String("1.1.1.1")},
			Status:    InstanceStatusRunning,
			Profile:   "aws",
		},
	}

	input := InstanceQueryInput{}

	// test ip
	input.Ip = tea.String("1.1.1.1")
	newInstances := input.Filter(instances)
	assert.Equal(t, 2, len(newInstances))

	input.Name = tea.String("test")
	newInstances = input.Filter(instances)
	assert.Equal(t, 1, len(newInstances))

	input.Profile = tea.String("aws")
	newInstances = input.Filter(instances)
	assert.Equal(t, 0, len(newInstances))

	input.Profile = tea.String("tencent")
	newInstances = input.Filter(instances)
	assert.Equal(t, 1, len(newInstances))

	input.Status = InstanceStatusRunning.TString()
	newInstances = input.Filter(instances)
	assert.Equal(t, 0, len(newInstances))

	input.Status = InstanceStatusPending.TString()
	newInstances = input.Filter(instances)
	assert.Equal(t, 1, len(newInstances))
}
