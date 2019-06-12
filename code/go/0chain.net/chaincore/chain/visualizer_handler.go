package chain

import (
	"encoding/json"
	"math"
	"net/http"
	"sort"

	"0chain.net/chaincore/block"

	"0chain.net/chaincore/node"
)

type bNode struct {
	ID                 string  `json:"id"`
	PrevID             string  `json:"prev_id"`
	Round              int64   `json:"round"`
	Rank               int     `json:"rank"`
	GeneratorID        int     `json:"generator_id"`
	GeneratorName      string  `json:"generator_name"`
	ChainWeight        float64 `json:"chain_weight"`
	Verifications      int     `json:"verifications"`
	Verified           bool    `json:"verified"`
	VerificationFailed bool    `json:"verification_failed"`
	Notarized          bool    `json:"notarized"`
	Finalized          bool    `json:"finalized"`
	X                  int     `json:"x"`
	Y                  int     `json:"y"`
	Size               int     `json:"size"`
}

// GetActiveSetMinerIndex Get the miner's SetIndex if found in ActiveSet. Else, return -1
func (c *Chain) GetActiveSetMinerIndex(roundNum int64, bgNode *node.GNode) int {
	n := c.GetActivesetMinerForRound(roundNum, bgNode)

	if n != nil {
		return n.SetIndex
	}
	return -1
}

//WIPBlockChainHandler - all the blocks in the memory useful to visualize and debug
func (c *Chain) WIPBlockChainHandler(w http.ResponseWriter, r *http.Request) {
	bl := c.getBlocks()
	var minr int64 = math.MaxInt64
	var maxr int64
	for _, b := range bl {
		if b.Round < minr {
			minr = b.Round
		}
		if b.Round > maxr {
			maxr = b.Round
		}
	}
	if minr < maxr-12 {
		minr = maxr - 12
	}
	if minr <= 0 {
		minr = 1
	}
	sort.SliceStable(bl, func(i, j int) bool {
		if bl[i].Round == bl[j].Round {
			return bl[i].RoundRank < bl[j].RoundRank
		}
		return bl[i].Round < bl[j].Round
	})
	finzalizedBlocks := make(map[string]bool)
	for fb := c.LatestFinalizedBlock; fb != nil; fb = fb.PrevBlock {
		finzalizedBlocks[fb.Hash] = true
	}
	bNodes := make([]*bNode, 0, len(bl))
	radius := 3
	padding := 5
	DXR := c.NumGenerators*radius + padding
	DYR := DXR
	for _, b := range bl {
		if b.Round < minr {
			continue
		}
		miner := node.GetNode(b.MinerID)
		x := int(b.Round - minr)
		y := c.GetActiveSetMinerIndex(b.Round, miner)
		_, finalized := finzalizedBlocks[b.Hash]
		bNd := &bNode{
			ID:                 b.Hash,
			PrevID:             b.PrevHash,
			Round:              b.Round,
			Rank:               b.RoundRank,
			GeneratorID:        y,
			GeneratorName:      miner.Description,
			ChainWeight:        b.ChainWeight,
			Verifications:      len(b.VerificationTickets),
			Verified:           b.GetVerificationStatus() != block.VerificationPending,
			VerificationFailed: b.GetVerificationStatus() == block.VerificationFailed,
			Notarized:          b.IsBlockNotarized(),
			Finalized:          finalized,
			X:                  x*DXR*2 + DXR,
			Y:                  y*DYR*2 + DYR,
			Size:               6 * (c.NumGenerators - b.RoundRank),
		}
		bNodes = append(bNodes, bNd)
	}
	//TODO: make CORS more restrictive
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bNodes)
}
