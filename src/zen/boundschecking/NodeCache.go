package boundschecking

type NodeCache struct {
	elements map[int32][]NormalizedNode
}

func (cache *NodeCache) GetNodeSingleton(node NormalizedNode) NormalizedNode {
	hashCode := node.GetHashCode()
	nodesWithHash := cache.elements[hashCode]

	for _, nodeCheck := range nodesWithHash {
		if node.IsEqual(nodeCheck) {
			return nodeCheck
		}
	}

	cache.elements[hashCode] = append(nodesWithHash, node)

	return node
}
