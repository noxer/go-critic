package lint

import "go/ast"

func init() {
	addChecker(&defaultCaseOrderChecker{}, attrExperimental, attrSyntaxOnly)
}

type defaultCaseOrderChecker struct {
	checkerBase
}

func (c *defaultCaseOrderChecker) InitDocumentation(d *Documentation) {
	d.Summary = "Detects when default case in switch isn't on 1st or last position"
	d.Before = `
switch {
case x > y:
	// ...
default: // <- not the best position
	// ...
case x == 10:
	// ...
}`
	d.After = `
switch {
case x > y:
	// ...
case x == 10:
	// ...
default: // <- last case (could also be the first one)
	// ...
}`
}

func (c *defaultCaseOrderChecker) VisitStmt(stmt ast.Stmt) {
	swtch, ok := stmt.(*ast.SwitchStmt)
	if !ok {
		return
	}
	for i, stmt := range swtch.Body.List {
		caseStmt, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}
		// is `default` case
		if caseStmt.List == nil {
			if i != 0 && i != len(swtch.Body.List)-1 {
				c.warn(caseStmt)
			}
		}
	}
}

func (c *defaultCaseOrderChecker) warn(cause *ast.CaseClause) {
	c.ctx.Warn(cause, "consider to make `default` case as first or as last case")
}
