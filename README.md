# txx
Simple library to proxy request to database in order to make calls transactional.

Have you ever thought, that Go lacks the Java's `@Transactional` kludge in it? 
Well, this package solves that problem. It puts the transaction into the context 
and requires you to pass this context everywhere, wherever you want to have an 
opportunity to be run in transaction.

Obviously, that implies that only contextual database methods are supported.

**Warning:** despite that this package is quite handy, I strongly do not recommend
overuse it. Better to organize your business logic to be idempotent and your
database access services to make transactions begin/end moments as clear, as 
possible.

## Example

```go
db := txx.NewTxWrapper(sqlx.NewDb())

var dest string

err := db.Tx(context.Background(), nil, func(ctx context.Context) error {
    if err := db.GetContext(ctx, &dest, `SELECT something FROM somewhere`); err != nil {
		return fmt.Errorf("failed to get something: %w", err)
    }
	return nil
})
if err != nil {
	return fmt.Errorf("transaction failed: %w", err)
}
```

During the calls to the `Runner` interface methods, the `TxWrapper` checks the
context for having a transaction in it, and, if it has, it hijacks the call and
passes it to the transaction in the context, instead of the database driver.

## Transactional business logic
You may try to use this thing to make your business logic transactional.

Let's say, you have several object stores:
```go
type Object1Store struct{db *txx.TxWrapper}

func (s *Object1Store) PerformSomeOperation(ctx context.Context, ...) error {}

type Object2Store struct{db *txx.txWrapper}

func (s *Object2Store) PerformSomeOtherOperation(ctx context.Context, ...) error {}
```

And you want to run `Object1Store.PerformSomeOperation` and `Object2Store.PerformSomeOtherOperation`
within the transaction (for instance, in case when you have several instances of the service and
using the `sync.Lock` and other synchronization primitives becomes impossible), then, you
may use the `TransactionalFactory` for that purpose:

```go
type TransactionalFactory struct {
	*Object1Store
	*Object2Store
	tx *txx.TxWrapper
}

func (f *TransactionalFactory) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return f.tx.Tx(ctx, nil, fn)
}
```
