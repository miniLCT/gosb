package constraints

type (
	// Call is Function Call
	Call func()

	// Consumer is one input argument and no returns Function
	Consumer[T any] func(T)

	// BiConsumer is two input arguments and no returns Function
	BiConsumer[T, U any] func(T, U)

	// UFunc is Unary Function
	UFunc[P, R any] func(P) R

	// BiFunc is Binary Function
	BiFunc[P1, P2, R any] func(P1, P2) R

	// Cmp is Compare Function
	Cmp[T any] BiFunc[T, T, int]

	// Eql is Equal Function
	Eql[T any] BiFunc[T, T, bool]

	// Less is Less Function
	Less[T any] BiFunc[T, T, bool]
)
