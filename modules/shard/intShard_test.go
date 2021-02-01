package number

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIntID(t *testing.T) {
	// shard := hash("1.2.3.4") % 4096
	// fmt.Printf("Shard: %d\n", shard)
	// BigInt(8) 9223372036854775807
	srcID := int64(241294492511762325)
	shardID := ParseIntShardID(srcID)
	assert.Equal(t, srcID, shardID.SrcID)
	assert.Equal(t, int64(3429), shardID.ShardID)
	assert.Equal(t, int64(1), shardID.TypeID)
	assert.Equal(t, int64(7075733), shardID.LocalID)

	// idObj := buildIntID(origID.ShardID, origID.TypeID, origID.LocalID)
	// fmt.Printf("NewID: %d\n", idObj.SrcID)
}

func TestBuildIntShardID(t *testing.T) {

	srcID := int64(241294492511762325)
	shardID := BuildIntShardID(int64(3429), int64(1), int64(7075733))

	assert.Equal(t, srcID, shardID.SrcID)
	assert.Equal(t, int64(3429), shardID.ShardID)
	assert.Equal(t, int64(1), shardID.TypeID)
	assert.Equal(t, int64(7075733), shardID.LocalID)

}
