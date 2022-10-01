package storage

import (
	"cardWithWords/internal/pkg/data/words/russian"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Words interface {
	GetCard(quantity int) (string, error)
}

type database struct {
	db *badger.DB
}

// GetAccessToWords opens or creates a new Badger database and check words in it
func GetAccessToWords(path string) (Words, error) {
	var (
		d   = new(database)
		err error
	)

	d.db, err = badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, fmt.Errorf("[GetAccessToWords] couldn't open database in path %s: %v", path, err)
	}

	var quantity int
	if err = d.db.View(
		func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = false
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Rewind(); it.Valid(); it.Next() {
				quantity++
			}

			return nil
		},
	); err != nil {
		return nil, fmt.Errorf("[GetAccessToWords] couldn't get quantity of words: %v", err)
	}

	if quantity != len(russian.Russian) {
		// drop inconsistent data
		if err = d.db.DropAll(); err != nil {
			return nil, fmt.Errorf("[GetAccessToWords] couldn't drop database: %v", err)
		}

		// insert words
		if err = d.db.Update(
			func(txn *badger.Txn) error {
				var i int
				for _, word := range russian.Russian {
					if err := txn.Set([]byte(strconv.Itoa(i)), []byte(word)); err != nil {
						return fmt.Errorf("[GetAccessToWords] couldn't set word: %v", err)
					}
					i++
				}

				return nil
			},
		); err != nil {
			return nil, fmt.Errorf("[GetAccessToWords] couldn't update database: %v", err)
		}
	}

	return d, nil
}

func (d *database) GetCard(quantity int) (string, error) {
	var (
		card string
		err  error
	)

	if err = d.db.View(
		func(txn *badger.Txn) error {
			for i := 0; i < quantity; i++ {
				key := strconv.Itoa(random(0, len(russian.Russian)-1))
				item, err := txn.Get([]byte(key))
				if err != nil {
					return err
				}

				valCopy, err := item.ValueCopy(nil)
				if err != nil {
					return fmt.Errorf("[GetCard] couldn't get value for key %s: %v", key, err)
				}

				card += strings.ToUpper(string(valCopy)) + "\n"
			}
			return nil
		},
	); err != nil {
		return "", fmt.Errorf("[GetCard] couldn't get words for card: %v", err)
	}

	return card, nil
}

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
