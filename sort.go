package main

// see
// SortMultiKeys example
// https://pkg.go.dev/sort#example-package-SortMultiKeys

import (
	"sort"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type lessFunc func(p1, p2 *ec2.Instance) bool

func OrderBy(lessFunc ...lessFunc) *MultiSorter {
	return &MultiSorter{
		less: lessFunc,
	}
}

type MultiSorter struct {
	instances []*ec2.Instance
	less      []lessFunc
}

func (ms *MultiSorter) Sort(instances []*ec2.Instance) {
	ms.instances = instances
	sort.Sort(ms)
}

func (ms *MultiSorter) Len() int { return len(ms.instances) }
func (ms *MultiSorter) Swap(i, j int) {
	ms.instances[i], ms.instances[j] = ms.instances[j], ms.instances[i]
}
func (ms *MultiSorter) Less(i, j int) bool {
	p, q := ms.instances[i], ms.instances[j]
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}
