package utils

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"strings"
	"time"
)

const (
	// UserEvent is an event created by a user
	UserEvent = "userEvent"
)

type StoreEvent struct {
	EventType string      `json:"EventType"`
	Recipient int64       `json:"Recipient"`
	Obj       interface{} `json:"Obj"`
}

type IStore interface {
	Connect(host string, password string, db int)
	Subscribe(channelNames ...string)
	Set(key string, content interface{})
	Keys(search string) (keys []string)
	SetString(key string, content string)
	Delete(key string)
	Publish(channel string, content interface{})
	PublishEvent(eventType string, recipient int64, obj interface{})
	PublishString(channel string, content string)
	ReceiveMessage() *redis.Message
	Ping() string
	Close()
	Get(key string, obj interface{}) (e error)
}

// Store is the redis store
type Store struct {
	Client      *redis.Client
	pubSub      *redis.PubSub
	isConnected bool
}

// NewStore creates a new store
func NewStore(host string, password string, db int) *Store {
	store := Store{}
	store.Connect(host, password, db)
	return &store
}

// Connect connects to the redis store
func (s *Store) Connect(host string, password string, db int) {

	log.Printf("Connecting to store @ %s", host)
	if s.isConnected == true {
		panic("Store is already connected...")
	}
	s.Client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db, // use default DB
	})

	s.isConnected = true
}

// Subscribe subscribes to a channel on redis
func (s *Store) Subscribe(channelNames ...string) {

	s.pubSub = s.Client.PSubscribe(channelNames...)

	_, err := s.pubSub.ReceiveTimeout(time.Second)
	if err != nil {
		panic(err)
	}
}

// Set creates an entry to redis with an object as its content
func (s *Store) Set(key string, content interface{}) {
	log.Printf("Store.Set: %s", key)
	contentString, _ := json.Marshal(content)
	if e := s.Client.Set(key, contentString, 0).Err(); e != nil {
		panic(e)
	}
}

// SetString creates an entry to redis with a string as its content
func (s *Store) SetString(key string, content string) {
	if e := s.Client.Set(key, content, 0).Err(); e != nil {
		panic(e)
	}
}

// Delete delets from the store
func (s *Store) Delete(key string) {
	if e := s.Client.Del(key).Err(); e != nil {
		panic(e)
	}
	log.Printf("Store.Delete: %s", key)
}

// Publish publishes to a channel
func (s *Store) Publish(channel string, content interface{}) {
	contentString, _ := json.Marshal(content)
	if e := s.Client.Publish(channel, contentString).Err(); e != nil {
		panic(e)
	}
}

func (s *Store) PublishEvent(eventType string, recipient int64, obj interface{}) {
	event := &StoreEvent{
		EventType: eventType,
		Recipient: recipient,
		Obj:       obj,
	}
	s.Publish(UserEvent, event)
}

// PublishString publishes a string to a channel
func (s *Store) PublishString(channel string, content string) {
	if e := s.Client.Publish(channel, content).Err(); e != nil {
		panic(e)
	}
}

// ReceiveMessage receives a message from PubSub
func (s *Store) ReceiveMessage() *redis.Message {
	message, e := s.pubSub.ReceiveMessage()
	if e != nil {
		panic(e)
	}
	return message
}

// Ping sends a ping command to the redis store and returns the response
func (s *Store) Ping() string {

	var pong string

	var e error

	if pong, e = s.Client.Ping().Result(); e != nil {
		panic("There was an error connecting to the redis server")
	}

	return pong
}

// Close closes the connection to redis
func (s *Store) Close() {

	if !s.isConnected == false {
		panic("Cannot close. Store is not connected...")
	}

	log.Printf("Closing Store Connection")
	s.Client.Close()
}

// Get gets an object from the store
func (s *Store) Get(key string, obj interface{}) (e error) {

	log.Printf("Store.Get: %s", key)
	var val string

	if val, e = s.Client.Get(key).Result(); e != nil {
		return
	}

	decoder := json.NewDecoder(strings.NewReader(val))
	e = decoder.Decode(obj)
	return
}

// Keys returns a slice of keys based on a search string
func (s *Store) Keys(search string) (keys []string) {
	var e error
	if keys, e = s.Client.Keys(search).Result(); e != nil {
		keys = []string{}
	}

	return
}
