package datastructures

import (
	"testing"
	"zen/test"
)

func TestAddRemove(t *testing.T) {
	var set = new(BitSet32)

	test.Assert(t, !set.Has(0), "Does not have at start")
	test.Assert(t, set.size == 0, "Should start empty")
	set.AddToSet(0)
	test.Assert(t, set.size == 1, "Should increase size")
	test.Assert(t, set.Has(0), "Have the value after")
	set.AddToSet(0)
	test.Assert(t, set.size == 1, "Should not double add")
	test.Assert(t, set.Has(0), "Have the value after")
	set.RemoveFromSet(0)
	test.Assert(t, set.size == 0, "Should decrement count")
	test.Assert(t, !set.Has(0), "Can remove")
	set.RemoveFromSet(0)
	test.Assert(t, set.size == 0, "Should not double remove")

	set.AddToSet(31)
	test.Assert(t, set.size == 1, "Should count max value")
	test.Assert(t, set.Has(31), "Should add max value")
}

func checkSubSet(t *testing.T, set *BitSet32, subsetSize uint32, expectedSize uint32) {
	var subSetCount = uint32(0)
	var alreadyHas = make(map[uint32]bool)

	var fullSet = set.ToInt32()

	set.ForEachSubSet(subsetSize, func(subSet *BitSet32) {
		subSetCount = subSetCount + 1
		var asInt = subSet.ToInt32()

		if (fullSet | asInt) != fullSet {
			t.Errorf("Bit inside subset that was not part of original set %b", asInt)
		}

		if alreadyHas[asInt] {
			t.Errorf("Duplicate subset %b", asInt)
		} else {
			alreadyHas[asInt] = true
		}
	})

	t.Logf("%d", subSetCount)
	test.Assert(t, subSetCount == expectedSize, "Expected number of subsets")
}

func TestEachSubset(t *testing.T) {
	var set = BitSet32FromData(uint32(0b10101))

	test.Assert(t, set.size == 3, "Size is correct")
	test.Assert(t, set.Has(0), "Has first bit")
	test.Assert(t, !set.Has(1), "Doesn't second bit")
	test.Assert(t, set.Has(2), "Has third bit")
	test.Assert(t, set.Has(4), "Has fifth bit")

	var firstSubSetIteration []uint32 = nil

	set.ForEachSubSet(2, func(subSet *BitSet32) {
		firstSubSetIteration = append(firstSubSetIteration, subSet.ToInt32())
	})

	test.Assert(t, len(firstSubSetIteration) == 3, "Correct number of sub sets")
	test.Assert(t, firstSubSetIteration[0] == 0b10100, "First subset entry")
	test.Assert(t, firstSubSetIteration[1] == 0b10001, "Second subset entry")
	test.Assert(t, firstSubSetIteration[2] == 0b00101, "Third subset entry")

	var secondSubSet []uint32 = nil

	set.ForEachSubSet(1, func(subSet *BitSet32) {
		secondSubSet = append(secondSubSet, subSet.ToInt32())
	})

	test.Assert(t, len(secondSubSet) == 3, "Correct number of sub sets for second set")
}

func TestManySubSets(t *testing.T) {
	checkSubSet(t, BitSet32FromData(0b110011), 1, 4)
	checkSubSet(t, BitSet32FromData(0b110011), 2, 6)
	checkSubSet(t, BitSet32FromData(0b110011), 3, 4)

	checkSubSet(t, BitSet32FromData(0b1100110011), 3, 20)
}
