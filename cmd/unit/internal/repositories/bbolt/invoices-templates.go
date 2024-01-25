package bbolt

import (
	"encoding/json"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	"github.com/devdammit/shekel/internal/resources"
	"go.etcd.io/bbolt"
)

var InvoicesTemplatesBucket = []byte("invoices_templates")

type InvoicesTemplatesRepository struct {
	db   *resources.Bolt
	data map[uint64]entities.InvoiceTemplate
}

func NewInvoicesTemplatesRepository(db *resources.Bolt) *InvoicesTemplatesRepository {
	return &InvoicesTemplatesRepository{
		db:   db,
		data: make(map[uint64]entities.InvoiceTemplate),
	}
}

func (r *InvoicesTemplatesRepository) GetName() string {
	return "invoices_templates"
}

func (r *InvoicesTemplatesRepository) Start() error {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		root := tx.Bucket(resources.BoltRootBucket)

		_, err := root.CreateBucketIfNotExists(InvoicesTemplatesBucket)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return r.db.View(func(tx *bbolt.Tx) error {
		invoices := tx.Bucket(resources.BoltRootBucket).Bucket(InvoicesTemplatesBucket)

		err = invoices.ForEach(func(k, v []byte) error {
			var invoice entities.InvoiceTemplate

			err = json.Unmarshal(v, &invoice)
			if err != nil {
				return err
			}

			r.data[invoice.ID] = invoice

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *InvoicesTemplatesRepository) CreateTx(tx *bbolt.Tx, invoice entities.InvoiceTemplate) (*entities.InvoiceTemplate, error) {
	invoices := tx.Bucket(resources.BoltRootBucket).Bucket(InvoicesTemplatesBucket)

	id, err := invoices.NextSequence()
	if err != nil {
		return nil, err
	}

	invoice.ID = id

	data, err := json.Marshal(&invoice)
	if err != nil {
		return nil, err
	}

	err = invoices.Put(resources.Itob(int(id)), data)
	if err != nil {
		return nil, err
	}

	r.data[invoice.ID] = invoice

	return &invoice, nil
}

func (r *InvoicesTemplatesRepository) DeleteTx(tx *bbolt.Tx, id uint64) error {
	invoices := tx.Bucket(resources.BoltRootBucket).Bucket(InvoicesTemplatesBucket)

	err := invoices.Delete(resources.Itob(int(id)))
	if err != nil {
		return err
	}

	delete(r.data, id)

	return nil
}
