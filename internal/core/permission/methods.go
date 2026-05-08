package permission

import "context"

func (c *Checker) CanRead(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionRead)
}

func (c *Checker) CanWrite(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionWrite)
}

func (c *Checker) CanCreate(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionCreate)
}

func (c *Checker) CanDelete(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionDelete)
}

func (c *Checker) CanSubmit(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionSubmit)
}

func (c *Checker) CanCancel(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionCancel)
}

func (c *Checker) CanAmend(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionAmend)
}

func (c *Checker) CanPrint(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionPrint)
}

func (c *Checker) CanEmail(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionEmail)
}

func (c *Checker) CanExport(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionExport)
}

func (c *Checker) CanImport(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionImport)
}

func (c *Checker) CanShare(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionShare)
}

func (c *Checker) CanReport(ctx context.Context, doctype string, roles []string) (bool, error) {
	return c.Can(ctx, doctype, roles, ActionReport)
}

func (c *Checker) CanUserRead(ctx context.Context, user string, doctype string) (bool, error) {
	return c.CanUser(ctx, user, doctype, ActionRead)
}

func (c *Checker) CanUserCreate(ctx context.Context, user string, doctype string) (bool, error) {
	return c.CanUser(ctx, user, doctype, ActionCreate)
}

func (c *Checker) CanUserWrite(ctx context.Context, user string, doctype string) (bool, error) {
	return c.CanUser(ctx, user, doctype, ActionWrite)
}

func (c *Checker) CanUserDelete(ctx context.Context, user string, doctype string) (bool, error) {
	return c.CanUser(ctx, user, doctype, ActionDelete)
}
