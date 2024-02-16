package memcached

import (
	"bytes"
	"encoding/gob"
	// "os"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/matdexir/ping-ping/models"
)

type Client struct {
	client *memcache.Client
}

func NewMemcached() (*Client, error) {
	connStr := "127.0.0.1:8002"
	client := memcache.New(connStr)

	if err := client.Ping(); err != nil {
		return nil, err
	}
	client.Timeout = 100 * time.Millisecond
	client.MaxIdleConns = 100

	return &Client{
		client: client,
	}, nil
}

func (c *Client) GetPosts(pconst string) (models.QueryCache, error) {
	item, err := c.client.Get(pconst)
	if err != nil {
		return models.QueryCache{}, err
	}

	b := bytes.NewReader(item.Value)
	var queryItems models.QueryCache

	if err := gob.NewDecoder(b).Decode(&queryItems); err != nil {
		return models.QueryCache{}, err
	}
	return queryItems, nil
}

func (c *Client) SetPosts(qc models.QueryCache) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(qc); err != nil {
		return err
	}

	return c.client.Set(&memcache.Item{
		Key:        qc.Parameters,
		Value:      b.Bytes(),
		Expiration: int32(time.Now().Add(25 * time.Second).Unix()),
	})

}

func (c *Client) Close() error {
	if err := c.client.Close(); err != nil {
		return err
	}
	return nil
}
