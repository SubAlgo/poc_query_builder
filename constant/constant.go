package constant

type SQLOperation string

const (
	SQLOperation_Eq      SQLOperation = "="
	SQLOperation_Gt      SQLOperation = ">"
	SQLOperation_Gte     SQLOperation = ">="
	SQLOperation_Lt      SQLOperation = "<"
	SQLOperation_Lte     SQLOperation = "<="
	SQLOperation_Ne      SQLOperation = "<>"
	SQLOperation_In      SQLOperation = "IN"
	SQLOperation_Between SQLOperation = "BETWEEN"
	SQLOperation_Contain SQLOperation = "ILIKE"

	StringContainValuePattern = "%%%s%%"
	StringInValuePattern      = "(%%%s%%)"
)

type SQLWhereOperation string

const (
	SQLWhereOperation_NULL SQLWhereOperation = " WHERE"
	SQLWhereOperation_AND  SQLWhereOperation = " AND"
	SQLWhereOperation_OR   SQLWhereOperation = " OR"
	SQLWhereOperation_NOT  SQLWhereOperation = " NOT"
)
