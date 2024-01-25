package bbolt

import (
	"context"
	"encoding/json"
	"github.com/devdammit/shekel/cmd/unit/internal/entities"
	port "github.com/devdammit/shekel/cmd/unit/internal/ports/repositories"
	"github.com/devdammit/shekel/internal/resources"
	"go.etcd.io/bbolt"
	"sort"
)

var InvoicesBucket = []byte("invoices")

type InvoicesRepository struct {
	db   *resources.Bolt
	data map[uint64]entities.Invoice
}

func NewInvoicesRepository(db *resources.Bolt) *InvoicesRepository {
	return &InvoicesRepository{
		db:   db,
		data: make(map[uint64]entities.Invoice),
	}
}

func (r *InvoicesRepository) GetName() string {
	return "invoices"
}

func (r *InvoicesRepository) Start() error {
	err := r.db.Update(func(tx *bbolt.Tx) error {
		root := tx.Bucket(resources.BoltRootBucket)

		_, err := root.CreateBucketIfNotExists(InvoicesBucket)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return r.db.View(func(tx *bbolt.Tx) error {
		invoices := tx.Bucket(resources.BoltRootBucket).Bucket(InvoicesBucket)

		err = invoices.ForEach(func(k, v []byte) error {
			var invoice entities.Invoice

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

func (r *InvoicesRepository) CreateTx(tx *bbolt.Tx, invoice entities.Invoice) (*entities.Invoice, error) {
	invoices := tx.Bucket(resources.BoltRootBucket).Bucket(InvoicesBucket)

	id, _ := invoices.NextSequence()

	invoice.ID = id

	r.data[id] = invoice

	data, err := json.Marshal(invoice)
	if err != nil {
		return nil, err
	}

	err = invoices.Put(resources.Itob(int(invoice.ID)), data)
	if err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (r *InvoicesRepository) FindByDates(_ context.Context, req port.FindByDatesRequest) ([]entities.Invoice, error) {
	var invoices []entities.Invoice
	var counter uint64

	for i, invoice := range r.data {
		if invoice.Date.After(req.StartedAt.Time) && invoice.Date.Before(req.EndedAt.Time) {
			if req.Limit != nil && counter >= *req.Limit {
				break
			}
			if req.Offset != nil && i < *req.Offset {
				continue
			}

			invoices = append(invoices, invoice)

			counter++
		}
	}

	sort.Slice(invoices, func(i, j int) bool {
		if req.OrderBy != nil && *req.OrderBy == port.OrderByDateDesc {
			return invoices[i].Date.After(invoices[j].Date.Time)
		} else if req.OrderBy != nil && *req.OrderBy == port.OrderByDateAsc {
			return invoices[i].Date.Before(invoices[j].Date.Time)
		}

		return invoices[i].Date.After(invoices[j].Date.Time)
	})

	return invoices, nil
}

func (r *InvoicesRepository) GetByID(_ context.Context, id uint64) (*entities.Invoice, error) {
	invoice, ok := r.data[id]
	if !ok {
		return nil, port.ErrNotFound
	}

	return &invoice, nil
}

func (r *InvoicesRepository) GetByTemplateID(_ context.Context, id uint64) ([]entities.Invoice, error) {
	var invoices []entities.Invoice

	for _, invoice := range r.data {
		if invoice.Template == nil {
			continue
		}

		if invoice.Template.ID == id {
			invoices = append(invoices, invoice)
		}
	}

	return invoices, nil
}

func (r *InvoicesRepository) DeleteTx(tx *bbolt.Tx, id uint64) error {
	invoices := tx.Bucket(resources.BoltRootBucket).Bucket(InvoicesBucket)

	err := invoices.Delete(resources.Itob(int(id)))
	if err != nil {
		return err
	}

	delete(r.data, id)

	return nil
}

func (r *InvoicesRepository) UpdateTx(tx *bbolt.Tx, invoice entities.Invoice) (*entities.Invoice, error) {
	invoices := tx.Bucket(resources.BoltRootBucket).Bucket(InvoicesBucket)

	data, err := json.Marshal(invoice)
	if err != nil {
		return nil, err
	}

	err = invoices.Put(resources.Itob(int(invoice.ID)), data)
	if err != nil {
		return nil, err
	}

	r.data[invoice.ID] = invoice

	return &invoice, nil
}
