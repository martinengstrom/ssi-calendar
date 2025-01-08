package storage

import (
  "log"
  "bytes"
  "encoding/gob"
  "time"
  "ssi-calendar/client"
  badger "github.com/dgraph-io/badger/v4"
)

type Storage struct {
  Badger *badger.DB
}

type Settings struct {
  Token string
  RefreshToken string
  Expiry string
}

type Event struct {
}

func NewStorage() *Storage {
  db, err := badger.Open(badger.DefaultOptions("data"))
  if err != nil {
    log.Fatal(err)
  }
  return &Storage{Badger: db}
}

func (c *Storage) Close() {
  c.Badger.Close()
}

func (c *Storage) Serialize(value interface{}) ([]byte, error) {
  var buf bytes.Buffer
  enc := gob.NewEncoder(&buf)
  err := enc.Encode(value)
  if err != nil {
    return nil, err
  }
  return buf.Bytes(), nil
}

func (c *Storage) Deserialize(data []byte, value interface{}) error {
  buf := bytes.NewBuffer(data)
  dec := gob.NewDecoder(buf)
  return dec.Decode(value)
}

/*
  TODO:
  Refactored all of this shit to create some abstraction layers between storing specific data and storing/fetching generic data
  basically we need generic methods for writing, reading and readAll that are typ agnostic
*/

func (c *Storage) UpdateEvent(event client.EventDetails) {
  existingEvent := c.GetEvent(event.Id)
  eventUpdated := !event.IsEqualTo(existingEvent)

  if eventUpdated {
    log.Println("Event (" + event.Name + ") was updated")
    event.UpdatedAt = time.Now()
    err := c.Badger.Update(func(txn *badger.Txn) error {
      data, err := c.Serialize(event)
      if err != nil {
        return err
      }

      return txn.Set([]byte(event.Id), data)
    })

    if err != nil {
      log.Fatalf("Failed to write data: %v", err)
    }
  }
}

func (c *Storage) GetEvent(id string) client.EventDetails {
  var retrievedEvent client.EventDetails
  err := c.Badger.View(func(txn *badger.Txn) error {
    item, err := txn.Get([]byte(id))
    if err != nil {
      return err
    }

    return item.Value(func(val []byte) error {
      return c.Deserialize(val, &retrievedEvent)
    })
  })

  if err != nil {
    log.Printf("Failed to read event: %v", err)
  }

  return retrievedEvent
}

func (c* Storage) GetEvents() []client.EventDetails {
  var retrievedEvent client.EventDetails
  var retrievedEvents []client.EventDetails
  err := c.Badger.View(func(txn *badger.Txn) error {
    opts := badger.DefaultIteratorOptions
    opts.PrefetchValues = true
    it := txn.NewIterator(opts)
    defer it.Close()

    for it.Rewind(); it.Valid(); it.Next() {
      item := it.Item()
      //key := item.Key()

      err := item.Value(func(val []byte) error {
        e := c.Deserialize(val, &retrievedEvent)
        retrievedEvents = append(retrievedEvents, retrievedEvent)
        return e
      })
      if err != nil {
        return err
      }
    }

    return nil
  })

  if err != nil {
    log.Fatal(err)
  }

  return retrievedEvents;
}
