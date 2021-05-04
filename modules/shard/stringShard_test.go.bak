package number

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStringID(t *testing.T) {
	// shard := hash("1.2.3.4") % 4096
	// fmt.Printf("Shard: %d\n", shard)
	// BigInt(8) 9223372036854775807
	// srcID := int64(241294492511762325)
	shardID := GetShardIDFromStringID("64c4b7ec-7b88-4323-a3f5-57386c5e9bc2")
	fmt.Println(shardID)
	// assert.Equal(t, srcID, shardID.SrcID)
	// assert.Equal(t, int64(3429), shardID.ShardID)
	// assert.Equal(t, int64(1), shardID.TypeID)
	// assert.Equal(t, int64(7075733), shardID.LocalID)

	// idObj := buildIntID(origID.ShardID, origID.TypeID, origID.LocalID)
	// fmt.Printf("NewID: %d\n", idObj.SrcID)
}

func TestBuildStringShardID(t *testing.T) {

	srcID := int64(241294492511762325)
	shardID := BuildIntShardID(int64(3429), int64(1), int64(7075733))

	assert.Equal(t, srcID, shardID.SrcID)
	assert.Equal(t, int64(3429), shardID.ShardID)
	assert.Equal(t, int64(1), shardID.TypeID)
	assert.Equal(t, int64(7075733), shardID.LocalID)

}
