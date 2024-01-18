package increment_period

import (
	"context"
	"fmt"
	"github.com/devdammit/shekel/cmd/unit/internal/use-cases/increment-period/units"
	"github.com/devdammit/shekel/pkg/log"
	"regexp"
	"runtime/debug"
	"time"
)

type Service struct {
	chain []units.Unit
}

var panicRegex = regexp.MustCompile(`runtime/panic\.go[ \S]+\s[ \S]+\s+([\S]+\.go:\d+)`)

func formatRecover(r interface{}, message string) string {
	var res string

	stack := string(debug.Stack())
	results := panicRegex.FindAllStringSubmatch(stack, -1)

	if len(results) > 0 {
		res = fmt.Sprintf("%s panic: %s in %s", message, r, results[0][1])
	} else {
		res = fmt.Sprintf("%s panic: %s stack: %s", message, r, debug.Stack())
	}

	return res
}

func NewService(
	periodsRepository units.PeriodsRepository,
	invoicesRepository units.InvoicesRepository,
	transactionsRepository units.TransactionsRepository,
	accountsRepository units.AccountsRepository,
	accountsService units.AccountService,
) *Service {
	chain := []func() units.Unit{
		func() units.Unit { return units.NewPeriodExtender(periodsRepository) },
		func() units.Unit { return units.NewInvoicesChecker(invoicesRepository) },
		func() units.Unit { return units.NewTransactionsExtender(transactionsRepository) },
		func() units.Unit { return units.NewBalanceCalculator(accountsService) },
		func() units.Unit { return units.NewResultSaver(accountsRepository, periodsRepository) },
	}

	service := &Service{
		chain: make([]units.Unit, len(chain)),
	}

	for i, factory := range chain {
		unit := factory()

		service.chain[i] = unit
	}

	return service
}

func (s *Service) Execute(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			err := units.ErrorUnitPanic(formatRecover(r, "unit handle got panic"))
			log.With(log.Err(err)).Error("error in chain")
		}
	}()

	start := time.Now()
	payload := units.NewPayload()

	request := units.Request{}

	for _, unit := range s.chain {
		_, err := unit.Handle(ctx, &request, payload)
		if err != nil {
			log.WithContext(ctx).With(log.Err(err), log.String("unit", unit.GetName())).Error("err has stopped the chain")
		}

		return err
	}

	log.WithContext(ctx).With(log.Duration("duration", time.Since(start))).Info("period closed")

	return nil
}
