package main

import (
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func TestAdd(t *testing.T) {
	core := resource.NewQuantity(1, resource.DecimalSI)
	v := v1.ResourceList{
		v1.ResourceCPU: *core,
	}

	require.Equal(t, int64(1), v.Cpu().Value())

	r := ResourceStat{}
	r.Limit = add(r.Limit, v)
	r.Limit = add(r.Limit, v)
	require.Equal(t, int64(2), r.Limit.Cpu().Value())
}
