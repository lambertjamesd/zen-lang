package datastructures

const DIRTY_SIZE = ^uint32(0)

type BitSet32 struct {
	data uint32
	size uint32
}

func countBits(value uint32) uint32 {
	var size uint32 = 0

	for i := uint32(0); i < 32; i = i + 1 {
		if (uint32(1)<<i)&value != 0 {
			size = size + 1
		}
	}

	return size
}

func BitSet32FromData(data uint32) *BitSet32 {
	return &BitSet32{
		data,
		countBits(data),
	}
}

func (bitSet *BitSet32) AddToSet(value uint32) {
	var mask = uint32(1) << value

	if (bitSet.data & mask) != mask {
		bitSet.data = bitSet.data | mask
		if bitSet.size != DIRTY_SIZE {
			bitSet.size = bitSet.size + 1
		}
	}
}

func (bitSet *BitSet32) RemoveFromSet(value uint32) {
	var mask = uint32(1) << value

	if (bitSet.data & mask) == mask {
		bitSet.data = bitSet.data & ^mask
		if bitSet.size != DIRTY_SIZE {
			bitSet.size = bitSet.size - 1
		}
	}
}

func (bitSet *BitSet32) Union(other *BitSet32) {
	bitSet.size = DIRTY_SIZE
	bitSet.data = bitSet.data | other.data
}

func (bitSet *BitSet32) Intersection(other *BitSet32) {
	bitSet.size = DIRTY_SIZE
	bitSet.data = bitSet.data & other.data
}

func (bitSet *BitSet32) Clear() {
	bitSet.size = uint32(0)
	bitSet.data = uint32(0)
}

func (bitSet *BitSet32) Has(value uint32) bool {
	var mask = uint32(1) << value
	return (bitSet.data & mask) == mask
}

func (bitSet *BitSet32) ToInt32() uint32 {
	return bitSet.data
}

func (bitSet *BitSet32) Size() uint32 {
	if bitSet.size == DIRTY_SIZE {
		bitSet.size = countBits(bitSet.data)
	}

	return bitSet.size
}

func (bitSet *BitSet32) ForEach(callback func(value uint32) bool) {
	var current = uint32(0)
	var data = bitSet.data

	for data != 0 {
		if (data & 1) != 0 {
			if !callback(current) {
				return
			}
		}
		current = current + 1
		data = data >> 1
	}
}

func (bitSet *BitSet32) ForEachSubSet(subsetSize uint32, callback func(set BitSet32)) {
	if subsetSize >= bitSet.Size() {
		callback(*bitSet)
	} else {
		var bitSetCopy = *bitSet
		bitSetCopy.forEachSubSet(subsetSize, 0, bitSet.size, BitSet32{0, 0}, callback)
	}
}

func (bitSet BitSet32) forEachSubSet(
	targetSize uint32,
	currentSearchValue uint32,
	remainingValues uint32,
	currentSet BitSet32,
	callback func(set BitSet32),
) {
	if currentSet.size == targetSize {
		callback(currentSet)
	} else if currentSet.size+remainingValues >= targetSize {
		for !bitSet.Has(currentSearchValue) {
			currentSearchValue = currentSearchValue + 1
		}

		bitSet.forEachSubSet(targetSize, currentSearchValue+1, remainingValues-1, currentSet, callback)

		(&currentSet).AddToSet(currentSearchValue)

		bitSet.forEachSubSet(targetSize, currentSearchValue+1, remainingValues-1, currentSet, callback)
	}
}

func (bitSet *BitSet32) Copy() *BitSet32 {
	return &BitSet32{
		bitSet.data,
		bitSet.size,
	}
}
