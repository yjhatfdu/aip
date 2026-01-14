package cluster

import (
	"errors"
	"fmt"
	"sort"
)

type SigInfo struct {
	Sig      string
	Count    int
	FirstTS  string
	LastTS   string
	SampleTS string
	Sample   string
	Hash     uint64
}

type Params struct {
	Threshold  int
	Bands      int
	BandBits   int
	MinCluster int
	Samples    int
}

func ClusterSigs(items []SigInfo, params Params) ([]Cluster, error) {
	if params.Bands <= 0 || params.BandBits <= 0 {
		return nil, errors.New("bands and band-bits must be > 0")
	}
	if params.Bands*params.BandBits != 64 {
		return nil, fmt.Errorf("bands*band-bits must equal 64 (got %d)", params.Bands*params.BandBits)
	}
	if params.MinCluster <= 0 {
		params.MinCluster = 1
	}
	if params.Samples < 0 {
		params.Samples = 0
	}
	if params.Threshold < 0 {
		params.Threshold = 0
	}

	for i := range items {
		items[i].Hash = Simhash(items[i].Sig, items[i].Count)
	}

	buckets := make(map[bandKey][]int)
	for i, item := range items {
		for b := 0; b < params.Bands; b++ {
			buckets[bandKey{band: b, val: bandValue(item.Hash, b, params.BandBits)}] = append(buckets[bandKey{band: b, val: bandValue(item.Hash, b, params.BandBits)}], i)
		}
	}

	uf := newUnionFind(len(items))
	for _, idxs := range buckets {
		if len(idxs) < 2 {
			continue
		}
		for i := 0; i < len(idxs); i++ {
			for j := i + 1; j < len(idxs); j++ {
				a := idxs[i]
				b := idxs[j]
				if Hamming(items[a].Hash, items[b].Hash) <= params.Threshold {
					uf.union(a, b)
				}
			}
		}
	}

	type clusterAgg struct {
		Cluster
		reprCount int
	}
	clusters := map[int]*clusterAgg{}
	for i, item := range items {
		root := uf.find(i)
		c, ok := clusters[root]
		if !ok {
			c = &clusterAgg{}
			clusters[root] = c
		}
		c.Count += item.Count
		if c.Repr == "" || item.Count > c.reprCount || (item.Count == c.reprCount && item.Sig < c.Repr) {
			c.Repr = item.Sig
			c.reprCount = item.Count
		}
		if item.FirstTS != "" {
			if c.FirstTS == "" || item.FirstTS < c.FirstTS {
				c.FirstTS = item.FirstTS
			}
		}
		if item.LastTS != "" {
			if c.LastTS == "" || item.LastTS > c.LastTS {
				c.LastTS = item.LastTS
			}
		}
		if params.Samples > 0 && len(c.Samples) < params.Samples && item.Sample != "" {
			c.Samples = append(c.Samples, Sample{TS: item.SampleTS, Raw: item.Sample})
		}
	}

	out := make([]Cluster, 0, len(clusters))
	for _, c := range clusters {
		if c.Count >= params.MinCluster {
			out = append(out, c.Cluster)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count == out[j].Count {
			return out[i].Repr < out[j].Repr
		}
		return out[i].Count > out[j].Count
	})
	return out, nil
}

type bandKey struct {
	band int
	val  uint64
}

func bandValue(h uint64, band int, bits int) uint64 {
	shift := band * bits
	mask := uint64(1<<bits) - 1
	return (h >> shift) & mask
}
