package shard

import (
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"time"
)

var (
	// id size is the size of a full ID (in bits)
	// Note this is out of a maximum of 64 bits
	shortIDSize int64 = 64
	// DateSize is 41 bits
	// 64 - 23 = 41
	shortDateSize int64 = 41
	// shortShardSize is 10 bits?
	// 41 - 10 = 31
	shortShardSize int64 = 13
	// local int is 31 bits
	shortOurEpoch                = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 100000
	shortAutoIncrementSize int64 = 10
	// Number of ids we can generate in a millisecond
	shortAutIncrementCount int64 = 1024
	// length of string to shard against
	shortShardStrLength = 3
)

type ShortSharder struct {
	sequence       int64
	numberOfShards int64
	currentShardID int64
}

// NewRoundRobin returns a new ShortShardID based on an incrementing sequence integer mod the number of shards
// Allows for an equal distribution of entries across shards
func (s *ShortSharder) NewRoundRobin() *ShortShardID {
	shard := NewShortShardID(
		(s.sequence % s.numberOfShards),
		s.sequence,
	)
	s.sequence++
	return shard
}

// NewFromSubID returns a new ShortShardID based on some other value (subID) that is taken into account
// when generating the destination shard
// Allows for distribution of entries across shards to be dependent on some other value
// Usage:
// 		// This would group blog entries into shards based on createdByUserID
// 		blogEntryID := s.NewFromSubID(createdByUserID)
// 		blogEntryID2 := s.NewFromSubID(createdByUserID) // Same shard as line above
func (s *ShortSharder) NewFromSubID(subID int64) *ShortShardID {
	shard := NewShortShardID(
		subID%s.numberOfShards,
		s.sequence,
	)
	s.sequence++
	return shard
}

// GetShardFromSubID returns the shard number without incrementing the internal sequence
func (s *ShortSharder) GetShardFromSubID(subID int64) int64 {
	return subID % s.numberOfShards
}

func buildShortShardFromString(numberOfShards int64, str string) (shard int64, e error) {
	h := fnv.New32a()

	if len(str) == 0 {
		e = errors.New("Shard string cannot be empty")
		return
	}

	shardStrLengthInternal := shortShardStrLength
	if len(str) < shortShardStrLength {
		shardStrLengthInternal = len(str)
	}
	h.Write([]byte(str[0:shardStrLengthInternal]))
	shard = int64(h.Sum32()) % numberOfShards
	return
}

// NewFromString mods the first N characters of string str against the number of shards
// Throws an error if str has a length of zero
func (s *ShortSharder) NewFromString(str string) (shardID *ShortShardID, e error) {

	var shard int64
	if shard, e = buildShardFromString(s.numberOfShards, str); e != nil {
		return
	}

	shardID = NewShortShardID(
		shard,
		s.sequence,
	)

	s.sequence++

	return
}

func (s *ShortSharder) NumberOfShards() int64 {
	return s.numberOfShards
}

func (s *ShortSharder) CurrentShardID() int64 {
	return s.currentShardID
}

func NewShortSharder(numberOfShards int64) *ShortSharder {
	return &ShortSharder{
		numberOfShards: numberOfShards,
	}
}

type ShortShardID struct {
	id        int64
	shard     int64
	timestamp int64
	sequence  int64
}

func NewShortShardIDFromID(id int64) *ShortShardID {
	s := &ShortShardID{
		id: id,
	}
	s.timestamp = s.Timestamp()
	s.sequence = s.Sequence()
	s.shard = s.Shard()
	return s
}

func NewShortShardID(shard, sequence int64) *ShortShardID {

	nowMillis := time.Now().UnixNano() / 100000
	timestamp := nowMillis - shortOurEpoch
	return &ShortShardID{
		shard:     shard,
		timestamp: timestamp,
		sequence:  sequence % shortAutIncrementCount,
	}
}

func (s *ShortShardID) Shard() int64 {
	return s.ID() % 1000
	// return s.ID() >> (shortIDSize - shortDateSize - shortShardSize) & (int64(math.Pow(2, float64(shortIDSize-shortDateSize-shortShardSize))) - 1)
}

func (s *ShortShardID) Timestamp() int64 {
	return s.ID() / 1000
	// return s.ID() >> (shortIDSize - shortDateSize)
}

func (s *ShortShardID) Sequence() int64 {
	return 0
	// return s.ID() >> (shortIDSize - shortDateSize - shortShardSize - shortAutoIncrementSize) & (int64(math.Pow(2, float64((shortAutoIncrementSize))) - 1))
}

// ID returns the id of the shard
func (s *ShortShardID) ID() int64 {
	if s.id == 0 {

		// tsBefore := (time.Now().UnixNano() / 100000) - shortOurEpoch
		// fmt.Println("tsBefore", tsBefore)
		// fmt.Println("shardBefore", s.shard)
		newIDString := fmt.Sprintf("%d%03d", s.timestamp, s.shard)
		s.id, _ = strconv.ParseInt(newIDString, 10, 64)

		// ts := newID / 1000
		// newShard := newID % 1000
		// fmt.Println("ts", ts)
		// fmt.Println("newShard", newShard)

		// fmt.Println("ID()", newID)
		// fmt.Println("ID()", s.shard, s.timestamp, s.sequence)
		// str := fmt.Sprintf("%04d%d%04d", s.shard, s.timestamp, s.sequence)
		// fmt.Println(str)
		// s.id, _ = strconv.ParseInt(str, 10, 64)
		// Bitwise inclusive OR
		// datePart * 2, 23 times
		// 2^23 * datePart

		// id := s.timestamp << (shortIDSize - shortDateSize) // 64 - 41

		// 2^10 * id
		// id |= s.shard << (shortIDSize - shortDateSize - shortShardSize) // (64 - 41 - 13)
		// id |= s.sequence
		// s.id = id
	}

	return s.id
}

// hash does an FNV1a hash of the string
func (s *ShortShardID) Hash() uint32 {

	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%d", s.ID())))
	return h.Sum32()

}

func (s *ShortShardID) String() string {
	return fmt.Sprintf("ID: %d, DatePart: %d, Shard: %d, Sequence: %d, Hash: %d", s.ID(), s.Timestamp(), s.Shard(), s.Sequence(), s.Hash())
}
