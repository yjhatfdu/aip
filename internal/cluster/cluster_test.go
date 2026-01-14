package cluster

import "testing"

func TestClusterInvalidBands(t *testing.T) {
	_, err := ClusterSigs([]SigInfo{{Sig: "a", Count: 1}}, Params{
		Bands:    4,
		BandBits: 8,
	})
	if err == nil {
		t.Fatal("expected error for invalid bands*band-bits")
	}
}

func TestClusterAllUnionWithHighThreshold(t *testing.T) {
	items := []SigInfo{
		{Sig: "error one", Count: 3, Sample: "a"},
		{Sig: "error two", Count: 1, Sample: "b"},
	}
	out, err := ClusterSigs(items, Params{
		Threshold:  64,
		Bands:      8,
		BandBits:   8,
		MinCluster: 1,
		Samples:    1,
	})
	if err != nil {
		t.Fatalf("ClusterSigs error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(out))
	}
	if out[0].Count != 4 {
		t.Fatalf("count mismatch: %d", out[0].Count)
	}
	if out[0].Repr != "error one" {
		t.Fatalf("repr mismatch: %q", out[0].Repr)
	}
	if len(out[0].Samples) != 1 {
		t.Fatalf("samples mismatch: %#v", out[0].Samples)
	}
}

func TestHamming(t *testing.T) {
	h := Simhash("abc", 1)
	if got := Hamming(h, h); got != 0 {
		t.Fatalf("hamming mismatch: %d", got)
	}
}
