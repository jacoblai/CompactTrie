package CompactTrie

import "sort"

type TrieNode struct {
	FeatureID  uint32
	ChildCount uint16
	ChildIndex uint32
	IsLeaf     bool
}

type CompactTrie struct {
	FeatureDict map[string]uint32
	Nodes       []TrieNode
	Index       []uint32
}

func NewCompactTrie(keys []string) *CompactTrie {
	trie := &CompactTrie{
		FeatureDict: make(map[string]uint32),
	}
	trie.BuildTrie(keys)
	trie.CompressTrie()
	return trie
}

func (trie *CompactTrie) Lookup(key string) bool {
	node := &trie.Nodes[0]
	pos := 0
	remain := len(key)

	for remain > 0 {
		feature := trie.GetLongestMatch(key[pos:])
		featureID := trie.GetFeatureID(feature)

		if featureID != node.FeatureID {
			break
		}

		if node.IsLeaf && remain == len(feature) {
			return true
		}

		pos += len(feature)
		remain -= len(feature)

		if remain > 0 {
			found := sort.Search(int(node.ChildCount), func(i int) bool {
				return trie.Nodes[node.ChildIndex+uint32(i)].FeatureID >= trie.GetFeatureID(string(key[pos]))
			})

			if found == int(node.ChildCount) {
				break
			}

			node = &trie.Nodes[node.ChildIndex+uint32(found)]
		}
	}

	return false
}

func (trie *CompactTrie) BuildTrie(keys []string) {
	trie.Nodes = make([]TrieNode, 1)

	for _, key := range keys {
		node := &trie.Nodes[0]
		pos := 0
		remain := len(key)

		for remain > 0 {
			feature := trie.GetLongestMatch(key[pos:])
			node = trie.InsertNode(node, feature)

			pos += len(feature)
			remain -= len(feature)
		}

		node.IsLeaf = true
	}
}

func (trie *CompactTrie) CompressTrie() {
	var queue []*TrieNode
	queue = append(queue, &trie.Nodes[0])

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		trie.CompressPath(node)
		trie.MergePrefix(node)

		for i := uint32(0); i < uint32(node.ChildCount); i++ {
			queue = append(queue, &trie.Nodes[node.ChildIndex+i])
		}
	}

	trie.Index = make([]uint32, 0, len(trie.Nodes))

	for i := range trie.Nodes {
		if trie.Nodes[i].IsLeaf {
			trie.Index = append(trie.Index, uint32(i))
		}
	}
}

func (trie *CompactTrie) InsertNode(node *TrieNode, feature string) *TrieNode {
	featureID := trie.GetFeatureID(feature)

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		child := &trie.Nodes[node.ChildIndex+i]
		if child.FeatureID == featureID {
			return child
		}
	}

	trie.Nodes = append(trie.Nodes, TrieNode{
		FeatureID: featureID,
	})
	newNode := &trie.Nodes[len(trie.Nodes)-1]

	if node.ChildCount == 0 {
		node.ChildIndex = uint32(len(trie.Nodes) - 1)
	}

	node.ChildCount++
	return newNode
}

func (trie *CompactTrie) NewNode(feature string) *TrieNode {
	trie.Nodes = append(trie.Nodes, TrieNode{
		FeatureID: trie.GetFeatureID(feature),
	})
	return &trie.Nodes[len(trie.Nodes)-1]
}

func (trie *CompactTrie) CompressPath(node *TrieNode) {
	if !node.IsLeaf && node.ChildCount == 1 {
		child := &trie.Nodes[node.ChildIndex]

		if !child.IsLeaf {
			node.FeatureID = trie.MergeFeature(node.FeatureID, child.FeatureID)
			node.ChildIndex = child.ChildIndex
			node.ChildCount = child.ChildCount
		}
	}
}

func (trie *CompactTrie) MergePrefix(node *TrieNode) {
	prefixMap := make(map[uint32][]*TrieNode)

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		child := &trie.Nodes[node.ChildIndex+i]
		prefixMap[child.FeatureID] = append(prefixMap[child.FeatureID], child)
	}

	for _, nodes := range prefixMap {
		if len(nodes) > 1 {
			mergedNode := nodes[0]

			for i := 1; i < len(nodes); i++ {
				child := nodes[i]
				mergedNode.ChildIndex = trie.MergeIndex(mergedNode.ChildIndex,
					mergedNode.ChildCount, child.ChildIndex, child.ChildCount)
				mergedNode.ChildCount += child.ChildCount
			}
		}
	}

	node.ChildCount = uint16(len(prefixMap))

	if node.ChildCount > 1 {
		newIndex := uint32(len(trie.Nodes))

		for _, nodes := range prefixMap {
			trie.Nodes = append(trie.Nodes, *nodes[0])
		}

		node.ChildIndex = newIndex
	}
}

func (trie *CompactTrie) GetFeatureID(feature string) uint32 {
	if id, ok := trie.FeatureDict[feature]; ok {
		return id
	}

	id := uint32(len(trie.FeatureDict))
	trie.FeatureDict[feature] = id
	return id
}

func (trie *CompactTrie) MergeFeature(feature1, feature2 uint32) uint32 {
	f1 := trie.GetFeature(feature1)
	f2 := trie.GetFeature(feature2)
	newFeature := f1 + f2
	return trie.GetFeatureID(newFeature)
}

func (trie *CompactTrie) GetFeature(featureID uint32) string {
	for f, id := range trie.FeatureDict {
		if id == featureID {
			return f
		}
	}

	panic("invalid feature id")
}

func (trie *CompactTrie) GetLongestMatch(s string) string {
	var longestMatch string

	for i := 1; i <= len(s); i++ {
		if _, ok := trie.FeatureDict[s[:i]]; ok {
			longestMatch = s[:i]
		}
	}

	return longestMatch
}

func (trie *CompactTrie) MergeIndex(index1 uint32, count1 uint16, index2 uint32, count2 uint16) uint32 {
	newIndex := uint32(len(trie.Nodes))

	for i := uint32(0); i < uint32(count1); i++ {
		trie.Nodes = append(trie.Nodes, trie.Nodes[index1+i])
	}

	for i := uint32(0); i < uint32(count2); i++ {
		trie.Nodes = append(trie.Nodes, trie.Nodes[index2+i])
	}

	return newIndex
}
