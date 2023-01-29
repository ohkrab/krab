package krab

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
	"github.com/wzshiming/ctc"
)

// CmdTestRun returns migration status information.
type CmdTestRun struct {
	Connection krabdb.Connection
	Suite      *TestSuite
	Registry   *CmdRegistry
}

// ResponseTestRun json
type ResponseTestRun struct {
}

func (c *CmdTestRun) Addr() krabhcl.Addr { return c.Suite.Addr() }

func (c *CmdTestRun) Name() []string { return append([]string{"test"}, c.Suite.Addr().Labels...) }

func (c *CmdTestRun) HttpMethod() string { return http.MethodPost }

func (c *CmdTestRun) Do(ctx context.Context, o CmdOpts) (interface{}, error) {
	var result ResponseTestRun

	for _, do := range c.Suite.Before.Dos {
		for _, migrate := range do.Migrate {
			addr, err := krabhcl.Expression{Expr: migrate.SetExpr}.Addr()
			if err != nil {
				return nil, fmt.Errorf("Failed to parse MigrationSet reference: %w", err)
			}

			for _, cmd := range c.Registry.Commands {
				if addr.Equal(cmd.Addr()) {
					if cmd.Name()[1] == migrate.Type {
						inputs := InputsFromCtyInputs(do.CtyInputs)
						migrateInputs := InputsFromCtyInputs(migrate.CtyInputs)
						inputs.Merge(migrateInputs)
						result, err := cmd.Do(ctx, CmdOpts{NamedInputs: inputs})
						if err != nil {
							return nil, fmt.Errorf("Failed to execute before hook: %w", err)
						}
						respUp, ok := result.([]ResponseMigrateUp)
						if ok {
							for _, migration := range respUp {
								fmt.Println(ctc.ForegroundYellow, "UP  ", migration.Success, migration.Version, migration.Name, ctc.Reset)
							}
						}
						respDown, ok := result.([]ResponseMigrateDown)
						if ok {
							for _, migration := range respDown {
								fmt.Println(ctc.ForegroundYellow, "DOWN", migration.Success, migration.Version, migration.Name, ctc.Reset)
							}
						}
					}
				}
			}
		}
	}

	err := c.Connection.Get(func(db krabdb.DB) error {
		resp, err := c.run(ctx, db, o.NamedInputs)
		result = resp
		return err
	})

	return result, err
}

func (c *CmdTestRun) run(ctx context.Context, db krabdb.DB, inputs NamedInputs) (ResponseTestRun, error) {
	result := ResponseTestRun{}

	for _, testCase := range c.Suite.Tests {
		fmt.Println(ctc.ForegroundBlue, testCase.Name, ctc.Reset)
		for _, it := range testCase.Its {
			fmt.Println("  ", ctc.ForegroundBlue, it.Comment, ctc.Reset)

			// apply SET parameters
			if testCase.Set != nil {
				sb := &strings.Builder{}
				testCase.Set.ToSQL(sb)
				_, err := db.ExecContext(ctx, sb.String())
				if err != nil {
					panic(fmt.Errorf("SET parameters not set: %w", err))
				}
			}

			// execute `do` from `it` and collect results for further expectations
			// fmt.Println("  ", it.Do.SQL)
			queryResult := []map[string]interface{}{}
			typeResult := map[string]string{}
			rows, capturedErr := db.QueryContext(ctx, it.Do.SQL)
			capturedErrConsumed := false

			if capturedErr == nil {
				defer rows.Close()

				types, _ := rows.ColumnTypes()
				for _, colType := range types {
					typeResult[colType.Name()] = colType.DatabaseTypeName()
				}
				for rows.Next() {
					row := map[string]interface{}{}
					rows.MapScan(row)
					queryResult = append(queryResult, row)
				}
			} else {
				// if query errored we can populate result map for easier processing later
				queryResult = append(queryResult, map[string]interface{}{
					"error": capturedErr.Error(),
				})
			}

			for _, asserts := range it.RowsAsserts {
				for _, expect := range asserts.Expectations {
					if expect.Subject == "error" {
						capturedErrConsumed = true
						fmt.Println(ctc.ForegroundGreen, *expect.Contains, ctc.ForegroundYellow, capturedErr.Error(), ctc.Reset)
					}
					if expect.Equal == nil {
						fmt.Println("    ", ctc.ForegroundGreen, expect.Subject, "= null", ctc.Reset)
					} else {
						fmt.Println("    ", ctc.ForegroundGreen, expect.Subject, "=", *expect.Equal, ctc.Reset)
					}
				}
			}

			if !capturedErrConsumed && capturedErr != nil {
				panic(fmt.Errorf("query `%s` failed to execute: %w", it.Do.SQL, capturedErr))
			}

			for _, rowAssert := range it.RowAsserts {
				rowAssert.EachRow(func(i int64) {
					query := &strings.Builder{}
					query.WriteString("SELECT key, value FROM ( VALUES\n")

					rowData := queryResult[i]
					for i, expect := range rowAssert.Expectations {
						value, ok := rowData[expect.Subject]
						if i != 0 {
							query.WriteString("\n,")
						}
						query.WriteString(fmt.Sprintf("(%s, %s)", krabdb.QuoteIdent(expect.Subject), krabdb.Quote(value)))
						fmt.Println(value)
						if !ok {
							panic("Expectation defined for missing column")
						}

						color := ctc.ForegroundGreen
						success := makeAssertion(expect, value)
						if !success {
							color = ctc.ForegroundRed
						}
						fmt.Println("    ", color, expect.Subject, value, *expect.Equal, ctc.Reset)
					}
					query.WriteString("\n) AS t(key, value)")
					fmt.Println("    ", ctc.ForegroundGreen, query.String(), ctc.Reset)
				})
			}
		}
	}

	// _, err := db.ExecContext(ctx, sql)

	return result, nil
}

func makeAssertion(expect *Expect, value interface{}) bool {
	if expect.Equal == nil {
		return value == nil
	}
	// fmt.Printf("// expect: %T, value: %T\n", *expect.Equal, value)

	return *expect.Equal == value
}
