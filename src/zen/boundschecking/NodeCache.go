package boundschecking

type NodeCache struct {
	elements map[int32][]NormalizedNode
}

func NewNodeCache() *NodeCache {
	return &NodeCache{
		make(map[int32][]NormalizedNode),
	}
}

func (cache *NodeCache) GetNodeSingleton(node NormalizedNode) NormalizedNode {
	hashCode := node.GetHashCode()
	nodesWithHash := cache.elements[hashCode]

	for _, nodeCheck := range nodesWithHash {
		if node.Compare(nodeCheck) == 0 {
			return nodeCheck
		}
	}

	cache.elements[hashCode] = append(nodesWithHash, node)

	return node
}
