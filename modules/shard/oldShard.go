package shard

const (
	// id size is the size of a full ID (in bits)
	// Note this is out of a maximum of 64 bits
	// idSize int64 = 62
	// shardSize is 16 bits
	// 62 - 16 = 46
	// shardSize int64 = 16
	// typeSize is 10 bits
	// 46 - 10 = 36
	typeSize int64 = 10
	// localIDSize is 36 bits
	// 36 - 36 = 0
	localIDSize int64 = 36
	// numberOfShards is the number of available shards
	numberOfShards int64 = 16
)

// IntShardID is a numeric unique identifier that includes shard information
type IntShardID struct {
	SrcID   int64
	ShardID int64
	TypeID  int64
	LocalID int64
}

// StringShardID is a string unique identifier that includes shard information
type StringShardID struct {
	SrcID   string
	ShardID int64
	TypeID  int64
	LocalID int64
}
