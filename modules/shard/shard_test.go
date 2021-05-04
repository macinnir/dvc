package shard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSharder_NewRoundRobin(t *testing.T) {

	// TODO: Insert the record first?
	// var autoIncrement int64 = 5001

	// var userID int64 = 123456
	var shardCount int64 = 10

	sharder := NewSharder(shardCount)

	for k := 0; k < 20; k++ {
		shard := sharder.NewRoundRobin()
		t.Logf("#%d: %d -> Shard: %d", shard.Sequence(), int64(k)%sharder.NumberOfShards(), shard.Shard())
		assert.Equal(t, int64(k), shard.Sequence())
		assert.Equal(t, int64(k)%sharder.NumberOfShards(), shard.Shard())
	}
}

func TestSharder_NewFromSubID(t *testing.T) {
	var shardCount int64 = 10
	sharder := NewSharder(shardCount)

	userIDs := []struct {
		userID  int64
		shardID int64
	}{
		{1, 1},
		{2, 2},
		{5, 5},
		{10, 0},
		{11, 1},
		{12, 2},
		{15, 5},
	}

	for k := range userIDs {
		shard := sharder.NewFromSubID(userIDs[k].userID)
		assert.Equal(t, userIDs[k].shardID, shard.Shard())
	}
}
