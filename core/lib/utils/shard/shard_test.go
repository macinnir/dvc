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
		t.Logf("#%d: %d -> Shard: %d ===> %d (@ %d)", shard.Sequence(), int64(k)%sharder.NumberOfShards(), shard.Shard(), shard.ID(), shard.Timestamp())
		assert.Equal(t, int64(k), shard.Sequence())
		assert.Equal(t, int64(k)%sharder.NumberOfShards(), shard.Shard())
	}
}

func TestSharder_NewShardIDFromID(t *testing.T) {
	id := int64(361468114661024787)

	shard := NewShardIDFromID(id)

	assert.Equal(t, int64(19), shard.Sequence())
	assert.Equal(t, int64(9), shard.Shard())
	assert.Equal(t, int64(43090357144), shard.Timestamp())
}

func TestSharder_buildShardFromString(t *testing.T) {

	var shardCount int64 = 10
	s := NewSharder(shardCount)
	tests := []struct {
		value    string
		shard    int64
		hasError bool
	}{
		{"", 0, true},
		{"a", 0, false},
		{"ab", 6, false},
		{"abc", 1, false},
		{"abcd", 1, false},
		{"abcde", 1, false},
		{"abcdef", 1, false},
		{"abcdefg", 1, false},
	}

	for k := range tests {

		shardID, e := s.NewFromString(tests[k].value)

		if tests[k].hasError {
			assert.NotNil(t, e)
			continue
		}

		assert.Nil(t, e)
		assert.Equal(t, tests[k].shard, shardID.shard)
		t.Logf("String: %s == Shard: %d", tests[k].value, shardID.shard)
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

func TestSharder_GetShardFromSubID(t *testing.T) {
	s := NewSharder(10)

	tests := [][]int64{
		{1, 1},
		{5, 5},
		{11, 1},
		{12, 2},
		{100, 0},
	}

	for k := range tests {
		result := s.GetShardFromSubID(tests[k][0])
		assert.Equal(t, tests[k][1], result)
	}
}

func BenchmarkSharder_NewRoundRobin(b *testing.B) {
	var shardCount int64 = 10

	sharder := NewSharder(shardCount)

	for n := 0; n < b.N; n++ {
		sharder.NewRoundRobin()
	}
}

func BenchmarkSharder_NewShardIDFromID(b *testing.B) {
	id := int64(361468114661024787)

	for n := 0; n < b.N; n++ {
		NewShardIDFromID(id)
	}
}

func BenchmarkSharder_buildShardFromString(b *testing.B) {

	var shardCount int64 = 10
	s := NewSharder(shardCount)
	for n := 0; n < b.N; n++ {
		s.NewFromString("abcdefg")
	}
}

func BenchmarkSharder_NewFromSubID(b *testing.B) {
	var shardCount int64 = 10
	sharder := NewSharder(shardCount)
	for n := 0; n < b.N; n++ {
		sharder.NewFromSubID(123456)
	}
}
