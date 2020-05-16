package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"time"
)

type Item struct {
	Price  float64 `bson:"p"`
	Amount float64 `bson:"a"`
}

type OrderBook struct {
	Symbol    string
	Timestamp time.Time `bson:"t"`
	Asks      []Item    `bson:"a"`
	Bids      []Item    `bson:"b"`
}

type Store struct {
	client   *mongo.Client
	database *mongo.Database

	obBuffer []*OrderBook
}

// NewStore new store
// uri: "mongodb://localhost:27017"
func NewStore(uri string, databaseName string) *Store {
	clientOptions := options.Client().ApplyURI(uri)
	// # 压缩通信(string或[]string), snappy(3.4) | zlib(3.6) | zstd(4.2)
	//    compressors: "snappy"
	clientOptions = clientOptions.SetCompressors([]string{"snappy"})
	var client *mongo.Client
	var err error
	if client, err = mongo.Connect(
		context.TODO(), clientOptions); err != nil {
		return nil
	}

	database := client.Database(databaseName)
	s := &Store{
		client:   client,
		database: database,
	}
	return s
}

func (s *Store) Close() error {
	s.SyncBuffer(true)
	return s.client.Disconnect(context.TODO())
}

func (s *Store) Insert(ob *OrderBook) error {
	s.obBuffer = append(s.obBuffer, ob)
	s.SyncBuffer(false)
	return nil
}

type item [2]float64

type ob struct {
	Timestamp int64  `bson:"t"`
	Asks      []item `bson:"a"`
	Bids      []item `bson:"b"`
}

func (s *Store) SyncBuffer(force bool) {
	if !force && len(s.obBuffer) < 100000 {
		return
	}
	if len(s.obBuffer) == 0 {
		return
	}

	type data struct {
		docs []interface{}
	}

	m := map[string]*data{}

	for _, v := range s.obBuffer {
		if d, ok := m[v.Symbol]; ok {
			d.docs = append(d.docs, s.convert(v))
		} else {
			d := &data{}
			d.docs = append(d.docs, s.convert(v))
			m[v.Symbol] = d
		}
	}

	for symbol, d := range m {
		s.SaveMany(exchange, symbol, d.docs)
	}

	s.obBuffer = nil
}

func (s *Store) convert(ob *OrderBook) (result ob) {
	result.Timestamp = ob.Timestamp.UnixNano() / int64(time.Millisecond)
	for _, v := range ob.Asks {
		result.Asks = append(result.Asks, item{
			v.Price,
			v.Amount,
		})
	}
	for _, v := range ob.Bids {
		result.Bids = append(result.Bids, item{
			v.Price,
			v.Amount,
		})
	}
	return
}

func (s *Store) SaveMany(exchange string, symbol string, docs []interface{}) (err error) {
	//c := s.database.Collection("deribit:" + symbol)
	var c *mongo.Collection
	c, err = s.GetCollection(exchange, symbol, false)
	if err != nil {
		return
	}
	_, err = c.InsertMany(context.TODO(), docs)
	return
}

func (s *Store) GetCollection(exchange string, symbol string, index bool) (c *mongo.Collection, err error) {
	c = s.database.Collection(exchange + ":" + symbol) // deribit:
	//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	//defer cancel()
	ctx := context.TODO()
	if index {
		var cur *mongo.Cursor
		cur, err = c.Indexes().List(ctx)
		if err != nil {
			return
		}
		key := bson.M{}
		exists := false
		for cur.Next(ctx) {
			err = cur.Decode(&key)
			if err != nil {
				log.Fatal(err)
			}
			//log.Println(key)
			if key["name"] == "t_1" {
				exists = true
			}
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		cur.Close(ctx)

		//log.Printf("%v", exists)

		// (DuplicateKey) E11000 duplicate key error collection: tick_db.deribit:BTC-PERPETUAL index: t_1 dup key: { : 1569888005038 }
		// map[key:map[_id:1] name:_id_ ns:tick_db.deribit:BTC-PERPETUAL v:2]
		if !exists {
			keys := bsonx.Doc{bsonx.Elem{Key: "t", Value: bsonx.Int64(1)}}
			indexOpts := options.Index() //.SetUnique(true)
			_, err := c.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: keys, Options: indexOpts})
			log.Printf("Index created")
			if err != nil {
				log.Fatal(err)
			}
		}

		if true {
			keys := bsonx.Doc{bsonx.Elem{Key: "t", Value: bsonx.Int64(1)}, bsonx.Elem{
				Key:   "_id",
				Value: bsonx.Int64(1),
			}}
			indexOpts := options.Index() //.SetUnique(true)
			_, err := c.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: keys, Options: indexOpts})
			log.Printf("Index created")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return
}
