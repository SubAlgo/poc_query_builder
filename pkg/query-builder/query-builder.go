package querybuilder

import (
	"fmt"
	"queryBuilder/pkg/constant"
	"strings"
)

type QueryBuilder struct {
	SQL          string
	Args         []interface{}
	selects      string
	table        string
	WhereClause  string
	Exceptions   []Exception // where clause
	limit        int
	offset       int
	Builder      strings.Builder
	RunningParam int
}

type Exception struct {
	WhereOperation constant.SQLWhereOperation
	Column         string
	Operation      constant.SQLOperation
	Value          interface{}
	runningParam   []int
}

func (q *QueryBuilder) ToSQL() (string, []interface{}) {
	q.Builder = strings.Builder{}
	q.Builder.WriteString(q.selects)

	if strings.Contains(strings.ToLower(q.table), "from") {
		fmt.Fprintf(&q.Builder, " %s", q.table)
	} else {
		fmt.Fprintf(&q.Builder, " from %s", q.table)
	}

	// set where clause
	for _, ex := range q.Exceptions {
		l := len(ex.runningParam)

		switch ex.Operation {
		case constant.SQLOperation_Between:
			fmt.Fprintf(&q.Builder, "%s (%s %s $%d AND $%d)", ex.WhereOperation, ex.Column, ex.Operation, ex.runningParam[0], ex.runningParam[1])
		case constant.SQLOperation_In:
			fmt.Fprintf(&q.Builder, "%s (%s %s", ex.WhereOperation, ex.Column, ex.Operation)
			for i := 0; i < l; i++ {
				if l == 1 {
					fmt.Fprintf(&q.Builder, " ($%d)", ex.runningParam[i])
				} else if i == 0 {
					fmt.Fprintf(&q.Builder, " ($%d, ", ex.runningParam[i])
				} else if i == l-1 {
					fmt.Fprintf(&q.Builder, "$%d)", ex.runningParam[i])
				} else {
					fmt.Fprintf(&q.Builder, "$%d, ", ex.runningParam[i])
				}
			}
			q.Builder.WriteString(")")
		case constant.SQLOperation_Ne_NULL:
			fmt.Fprintf(&q.Builder, "%s (%s %s)", ex.WhereOperation, ex.Column, ex.Operation)
		case constant.SQLOperation_Contain:
			fmt.Fprintf(&q.Builder, "%s (%s %s $%d)", ex.WhereOperation, ex.Column, ex.Operation, ex.runningParam[0])
		default:
			if l == 1 {
				fmt.Fprintf(&q.Builder, "%s (%s %s $%d)", ex.WhereOperation, ex.Column, ex.Operation, ex.runningParam[0])
			} else if l > 1 {
				fmt.Fprintf(&q.Builder, "%s (%s %s ", ex.WhereOperation, ex.Column, ex.Operation)

				for i := 0; i < l; i++ {
					if i == l-1 {
						fmt.Fprintf(&q.Builder, "$%d", ex.runningParam[i])
					} else {
						fmt.Fprintf(&q.Builder, "$%d, ", ex.runningParam[i])
					}
				}
				q.Builder.WriteString(")")
			}
		}
	}

	if q.offset > 0 {
		fmt.Fprintf(&q.Builder, " OFFSET %d", q.offset)
	}

	if q.limit > 0 {
		fmt.Fprintf(&q.Builder, " LIMIT %d", q.limit)
	}

	return q.Builder.String(), q.Args
}

func (q *QueryBuilder) Select(query string) (tx *QueryBuilder) {
	q.selects = query
	return q
}

func (q *QueryBuilder) Table(query string) (tx *QueryBuilder) {
	q.table = query
	return q
}

func (q *QueryBuilder) Where(except Exception) (tx *QueryBuilder) {
	n := len(q.Exceptions)
	if n == 0 {
		except.WhereOperation = constant.SQLWhereOperation_NULL
	} else {
		except.WhereOperation = constant.SQLWhereOperation_AND
	}

	if except.Operation == constant.SQLOperation_NULL || except.Operation == constant.SQLOperation_Ne_NULL {
		q.Exceptions = append(q.Exceptions, except)
		return q
	}

	var countValue int
	switch v := except.Value.(type) {
	case []string:
		countValue = len(v)
	case []int:
		countValue = len(v)
	case []interface{}:
		countValue = len(v)
	default:
		countValue = 1
	}

	// set running param $1, $2
	for i := 0; i < countValue; i++ {
		q.RunningParam += 1
		except.runningParam = append(except.runningParam, q.RunningParam)
	}
	q.Exceptions = append(q.Exceptions, except)

	q.addArgs(except.Value)
	return q
}

func (q *QueryBuilder) Or(except Exception) (tx *QueryBuilder) {
	n := len(q.Exceptions)
	if n == 0 {
		except.WhereOperation = constant.SQLWhereOperation_NULL
	} else {
		except.WhereOperation = constant.SQLWhereOperation_OR
	}

	if except.Operation == constant.SQLOperation_NULL || except.Operation == constant.SQLOperation_Ne_NULL {
		q.Exceptions = append(q.Exceptions, except)
		return q
	}

	var countValue int
	switch v := except.Value.(type) {
	case []string:
		countValue = len(v)
	case []int:
		countValue = len(v)
	case []interface{}:
		countValue = len(v)
	default:
		countValue = 1
	}

	// set running param $1, $2
	for i := 0; i < countValue; i++ {
		q.RunningParam += 1
		except.runningParam = append(except.runningParam, q.RunningParam)
	}
	q.Exceptions = append(q.Exceptions, except)

	q.addArgs(except.Value)
	return q
}

func (q *QueryBuilder) Limit(limit int) (tx *QueryBuilder) {
	q.limit = limit
	return q
}

func (q *QueryBuilder) Offset(offset int) (tx *QueryBuilder) {
	q.offset = offset
	return q
}

func (q *QueryBuilder) addArgs(arg interface{}) (tx *QueryBuilder) {
	switch v := arg.(type) {
	case []string:
		for _, arg := range v {
			q.Args = append(q.Args, arg)
		}
	case []int:
		for _, arg := range v {
			q.Args = append(q.Args, arg)
		}
	case []interface{}:
		q.Args = append(q.Args, v...)
	default:
		q.Args = append(q.Args, v)
	}
	return q
}

func (q *QueryBuilder) SetPagination(pageNo, pageSize int) (tx *QueryBuilder) {
	if pageNo <= 0 || pageSize <= 0 {
		q.offset = 0
		q.limit = 10
	}

	q.offset = (pageNo - 1) * pageSize
	q.limit = pageSize
	return q
}
