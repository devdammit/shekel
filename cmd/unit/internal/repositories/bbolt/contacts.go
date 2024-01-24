package bbolt

import (
	"context"
	"encoding/json"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/repositories"
	"github.com/devdammit/shekel/internal/resources"
	"github.com/devdammit/shekel/pkg/log"
	"github.com/devdammit/shekel/pkg/pointer"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"go.etcd.io/bbolt"
	"sync"
	"time"
)

var ContactsBucket = []byte("contacts")

type ContactsRepository struct {
	db   *resources.Bolt
	data map[uint64]entities.Contact

	sync.RWMutex
}

func NewContactsRepository(db *resources.Bolt) *ContactsRepository {
	return &ContactsRepository{
		db:   db,
		data: make(map[uint64]entities.Contact),
	}
}

func (r *ContactsRepository) GetName() string {
	return "contacts"
}

func (r *ContactsRepository) Start() error {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		root := tx.Bucket(resources.BoltRootBucket)

		_, err := root.CreateBucketIfNotExists(ContactsBucket)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return r.db.View(func(tx *bbolt.Tx) error {
		contacts := tx.Bucket(resources.BoltRootBucket).Bucket(ContactsBucket)

		err = contacts.ForEach(func(k, v []byte) error {
			var contact entities.Contact

			err = json.Unmarshal(v, &contact)
			if err != nil {
				return err
			}

			r.data[contact.ID] = contact

			return nil
		})
		if err != nil {
			return err
		}

		log.Info("contacts loaded", log.Int("count", len(r.data)))

		return nil
	})
}

func (r *ContactsRepository) Create(_ context.Context, contact entities.Contact) (*entities.Contact, error) {
	for _, entity := range r.data {
		if entity.Name == contact.Name && entity.Text == contact.Text {
			return nil, port.ErrAlreadyExists
		}
	}

	err := r.db.Update(func(tx *bbolt.Tx) error {
		contacts := tx.Bucket(resources.BoltRootBucket).Bucket(ContactsBucket)

		contact.ID, _ = contacts.NextSequence()
		contact.CreatedAt = datetime.NewDateTime(time.Now())
		contact.UpdatedAt = datetime.NewDateTime(time.Now())

		data, err := json.Marshal(contact)
		if err != nil {
			return err
		}

		return contacts.Put(resources.Itob(int(contact.ID)), data)
	})

	if err != nil {
		return nil, err
	}

	r.Lock()
	defer r.Unlock()

	r.data[contact.ID] = contact

	return &contact, nil
}

func (r *ContactsRepository) SetQRCode(_ context.Context, contactID uint64, qrCode entities.QRCode) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.data[contactID]; !ok {
		return port.ErrNotFound
	}

	contact := r.data[contactID]
	contact.QRCodes = append(contact.QRCodes, qrCode)
	contact.UpdatedAt = datetime.NewDateTime(time.Now())

	r.data[contactID] = contact

	err := r.db.Update(func(tx *bbolt.Tx) error {
		contacts := tx.Bucket(resources.BoltRootBucket).Bucket(ContactsBucket)

		data, err := json.Marshal(contact)
		if err != nil {
			return err
		}

		return contacts.Put(resources.Itob(int(contact.ID)), data)
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *ContactsRepository) RemoveQRCode(_ context.Context, contactID uint64, bankName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.data[contactID]; !ok {
		return port.ErrNotFound
	}

	contact := r.data[contactID]
	contact.UpdatedAt = datetime.NewDateTime(time.Now())

	for i, qrCode := range contact.QRCodes {
		if qrCode.BankName == bankName {
			contact.QRCodes = append(contact.QRCodes[:i], contact.QRCodes[i+1:]...)
			break
		}
	}

	r.data[contactID] = contact

	err := r.db.Update(func(tx *bbolt.Tx) error {
		contacts := tx.Bucket(resources.BoltRootBucket).Bucket(ContactsBucket)

		data, err := json.Marshal(contact)
		if err != nil {
			return err
		}

		return contacts.Put(resources.Itob(int(contact.ID)), data)
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *ContactsRepository) GetAll(_ context.Context, withDeleted *bool) ([]entities.Contact, error) {
	r.RLock()
	defer r.RUnlock()

	var contacts []entities.Contact

	for _, contact := range r.data {
		if contact.DeletedAt != nil && !*withDeleted {
			continue
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (r *ContactsRepository) GetByID(_ context.Context, ID uint64) (*entities.Contact, error) {
	r.RLock()
	defer r.RUnlock()

	if _, ok := r.data[ID]; !ok {
		return nil, port.ErrNotFound
	}

	contact := r.data[ID]

	return &contact, nil
}

func (r *ContactsRepository) Update(_ context.Context, contact entities.Contact) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.data[contact.ID]; !ok {
		return port.ErrNotFound
	}

	contact.UpdatedAt = datetime.NewDateTime(time.Now())

	r.data[contact.ID] = contact

	return r.db.Update(func(tx *bbolt.Tx) error {
		contacts := tx.Bucket(resources.BoltRootBucket).Bucket(ContactsBucket)

		data, err := json.Marshal(contact)
		if err != nil {
			return err
		}

		return contacts.Put(resources.Itob(int(contact.ID)), data)
	})
}

func (r *ContactsRepository) Remove(ctx context.Context, ID uint64) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.data[ID]; !ok {
		return port.ErrNotFound
	}

	contact := r.data[ID]
	contact.DeletedAt = pointer.Ptr(datetime.NewDateTime(time.Now()))

	r.data[ID] = contact

	err := r.db.Update(func(tx *bbolt.Tx) error {
		contacts := tx.Bucket(resources.BoltRootBucket).Bucket(ContactsBucket)

		data, err := json.Marshal(contact)
		if err != nil {
			return err
		}

		return contacts.Put(resources.Itob(int(contact.ID)), data)
	})

	if err != nil {
		return err
	}

	return nil
}
