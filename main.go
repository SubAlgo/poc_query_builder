package main

import (
	"fmt"
	"queryBuilder/constant"
	"strings"
)

func main() {
	q := QueryBuilder{}

	q.Select = `
	select 
	 	c.id, cb.business_id, 
	 	dv.id, dv.document_path, dv.document_type, dv.status, dv.expiry_date, 
	 	dv.file_type,
	 	dv.created_at, dv.created_by, dv.updated_at, dv.updated_by
	`

	q.Table = ` customer c
		inner join customer_business cb on c.id = cb.customer_id
		inner join business b ON cb.business_id = b.id
		left join document_validation dv on b.id  = dv.owner_id`

	q.
		Where(Exception{Column: "c.id", Operation: constant.SQLOperation_Eq, Value: "xxx-123"}).
		Where(Exception{Column: "c.age", Operation: constant.SQLOperation_In, Value: []interface{}{"33", "34"}})

	q.Or(Exception{Column: "dv.expiry", Operation: constant.SQLOperation_Between, Value: []interface{}{"2023-08-22", "2023-08-23"}}).
		SetPagination(2, 5)

	x, y := q.ToSQL()
	fmt.Println(x)
	fmt.Println(y...)
}

type QueryBuilder struct {
	SQL          string
	Args         []interface{}
	Select       string
	Table        string
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
	q.Builder.WriteString(q.Select)

	if strings.Contains(strings.ToLower(q.Table), "from") {
		fmt.Fprintf(&q.Builder, " %s", q.Table)
	} else {
		fmt.Fprintf(&q.Builder, " from %s", q.Table)
	}

	// set where clause
	for _, ex := range q.Exceptions {
		l := len(ex.runningParam)
		if ex.Operation == constant.SQLOperation_Between {
			fmt.Fprintf(&q.Builder, "%s (%s %s $%d AND $%d)", ex.WhereOperation, ex.Column, ex.Operation, ex.runningParam[0], ex.runningParam[1])
		} else {
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

func (q *QueryBuilder) Where(except Exception) (tx *QueryBuilder) {
	n := len(q.Exceptions)
	if n == 0 {
		except.WhereOperation = constant.SQLWhereOperation_NULL
	} else {
		except.WhereOperation = constant.SQLWhereOperation_AND
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
