package shard

import (
	"fmt"
	"hash/fnv"
	"math"
)

// ParseIntShardID parses an intID into a ShardID object
func ParseIntShardID(id int64) IntShardID {

	// fmt.Println(0xF)
	// os.Exit(0)

	shardSpacer := idSize - shardSize
	typeSpacer := shardSpacer - typeSize
	localIDSpacer := typeSpacer - localIDSize

	shardAnd := int64(math.Pow(2, float64(shardSize))) - 1
	typeAnd := int64(math.Pow(2, float64(typeSize))) - 1
	localIDAnd := int64(math.Pow(2, float64(localIDSize))) - 1

	fmt.Printf("%10s | %40s | %20d | %15d | %20d\n", "OrigID", "id", 0, 0, id)
	fmt.Printf("%10s | %40s | %20d | %15d | %20d\n", "ShardID", "(id>>shardSpacer)&0xFFFF", id>>shardSpacer, shardAnd, (id>>shardSpacer)&shardAnd)
	fmt.Printf("%10s | %40s | %20d | %15d | %20d\n", "TypeID", "(id>>typeSpacer)&0x3FF", id>>typeSpacer, typeAnd, (id>>typeSpacer)&typeAnd)
	fmt.Printf("%10s | %40s | %20d | %15d | %20d\n", "LocalID", "(id>>localIDSpacer)&0xFFFFFFFFF", id>>localIDSpacer, localIDAnd, (id>>localIDSpacer)&localIDAnd)

	return IntShardID{
		SrcID:   id,
		ShardID: (id >> shardSpacer) & 0xFFFF,        // 2^16 - 1 == 65535  == 16^4 - 1
		TypeID:  (id >> typeSpacer) & 0x3FF,          // 2^10 - 1 == 1023
		LocalID: (id >> localIDSpacer) & 0xFFFFFFFFF, // 2^36 - 1 == 68719476735 == 16^9 - 1
	}
}

// BuildIntShardID builds an integer ID
func BuildIntShardID(shardID, typeID, localID int64) IntShardID {

	shardSpacer := idSize - shardSize
	typeSpacer := shardSpacer - typeSize
	localIDSpacer := typeSpacer - localIDSize

	return IntShardID{
		SrcID:   (shardID << shardSpacer) | (typeID << typeSpacer) | (localID << localIDSpacer),
		ShardID: shardID,
		TypeID:  typeID,
		LocalID: localID,
	}
}

// hash does an FNV1a hash of the string
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
