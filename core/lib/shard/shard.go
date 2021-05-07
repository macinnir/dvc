package shard

import (
	"errors"
	"fmt"
	"hash/fnv"
	"math"
	"time"
)

var (
	// id size is the size of a full ID (in bits)
	// Note this is out of a maximum of 64 bits
	idSize int64 = 64
	// DateSize is 41 bits
	// 64 - 23 = 41
	dateSize int64 = 41
	// shardSize is 10 bits?
	// 41 - 10 = 31
	shardSize int64 = 13
	// local int is 31 bits
	ourEpoch                = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / int64(time.Millisecond)
	autoIncrementSize int64 = 10
	// Number of ids we can generate in a millisecond
	autoIncrementCount int64 = 1024
	// length of string to shard against
	shardStrLength = 3
)

type Sharder struct {
	sequence       int64
	numberOfShards int64
	currentShardID int64
}

// NewRoundRobin returns a new ShardID based on an incrementing sequence integer mod the number of shards
// Allows for an equal distribution of entries across shards
func (s *Sharder) NewRoundRobin() *ShardID {
	shard := NewShardID(
		(s.sequence % s.numberOfShards),
		s.sequence,
	)
	s.sequence++
	return shard
}

// NewFromSubID returns a new ShardID based on some other value (subID) that is taken into account
// when generating the destination shard
// Allows for distribution of entries across shards to be dependent on some other value
// Usage:
// 		// This would group blog entries into shards based on createdByUserID
// 		blogEntryID := s.NewFromSubID(createdByUserID)
// 		blogEntryID2 := s.NewFromSubID(createdByUserID) // Same shard as line above
func (s *Sharder) NewFromSubID(subID int64) *ShardID {
	shard := NewShardID(
		subID%s.numberOfShards,
		s.sequence,
	)
	s.sequence++
	return shard
}

func buildShardFromString(numberOfShards int64, str string) (shard int64, e error) {
	h := fnv.New32a()

	if len(str) == 0 {
		e = errors.New("Shard string cannot be empty")
		return
	}

	shardStrLengthInternal := shardStrLength
	if len(str) < shardStrLength {
		shardStrLengthInternal = len(str)
	}
	h.Write([]byte(str[0:shardStrLengthInternal]))
	shard = int64(h.Sum32()) % numberOfShards
	return
}

// NewFromString mods the first N characters of string str against the number of shards
// Throws an error if str has a length of zero
func (s *Sharder) NewFromString(str string) (shardID *ShardID, e error) {

	var shard int64
	if shard, e = buildShardFromString(s.numberOfShards, str); e != nil {
		return
	}

	shardID = NewShardID(
		shard,
		s.sequence,
	)

	s.sequence++

	return
}

func (s *Sharder) NumberOfShards() int64 {
	return s.numberOfShards
}

func (s *Sharder) CurrentShardID() int64 {
	return s.currentShardID
}

func NewSharder(numberOfShards int64) *Sharder {
	return &Sharder{
		numberOfShards: numberOfShards,
	}
}

type ShardID struct {
	id        int64
	shard     int64
	timestamp int64
	sequence  int64
}

func NewShardIDFromID(id int64) *ShardID {
	s := &ShardID{
		id: id,
	}
	s.timestamp = s.Timestamp()
	s.sequence = s.Sequence()
	s.shard = s.Shard()
	return s
}

func NewShardID(shard, sequence int64) *ShardID {

	nowMillis := time.Now().UnixNano() / int64(time.Millisecond)
	timestamp := nowMillis - ourEpoch
	return &ShardID{
		shard:     shard,
		timestamp: timestamp,
		sequence:  sequence % autoIncrementCount,
	}
}

func (s *ShardID) Shard() int64 {
	return s.ID() >> (idSize - dateSize - shardSize) & (int64(math.Pow(2, float64(idSize-dateSize-shardSize))) - 1)
}

func (s *ShardID) Timestamp() int64 {
	return s.ID() >> (idSize - dateSize)
}

func (s *ShardID) Sequence() int64 {
	return s.ID() >> (idSize - dateSize - shardSize - autoIncrementSize) & (int64(math.Pow(2, float64((autoIncrementSize))) - 1))
}

// ID returns the id of the shard
func (s *ShardID) ID() int64 {
	if s.id == 0 {
		// Bitwise inclusive OR
		// datePart * 2, 23 times
		// 2^23 * datePart
		id := s.timestamp << (idSize - dateSize) // 64 - 41
		// 2^10 * id
		id |= s.shard << (idSize - dateSize - shardSize) // (64 - 41 - 13)
		id |= s.sequence
		s.id = id
	}

	return s.id
}

// hash does an FNV1a hash of the string
func (s *ShardID) Hash() uint32 {

	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%d", s.ID())))
	return h.Sum32()

}

func (s *ShardID) String() string {
	return fmt.Sprintf("ID: %d, DatePart: %d, Shard: %d, Sequence: %d, Hash: %d", s.ID(), s.Timestamp(), s.Shard(), s.Sequence(), s.Hash())
}
