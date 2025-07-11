package run

import (
	"context"
	"os"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/plugin"
	"github.com/ohkrab/krab/fmtx"
	"github.com/ohkrab/krab/plugins"
	"github.com/ohkrab/krab/tpls"
)

type Runner struct {
	fs       *config.Filesystem
	tpls     *tpls.Templates
	registry *plugins.Registry
	migrator *Migrator
}

type Command any

// type CommandHandler interface {
// 	HandleMigrateAudit(audited *Audited)
// }

type CommandMigrateFixUp struct {
	Command

	Driver  string
	Set     string
	Version string
	Comment string
}

type CommandMigrateFixDown struct {
	Command

	Driver  string
	Set     string
	Version string
    Comment string
}

type CommandMigrateStatus struct {
	Command

	Driver string
	Set    string
}

type CommandMigrateUp struct {
	Command

	Driver string
	Set    string
}

type CommandMigrateDown struct {
	Command

	Driver  string
	Set     string
	Version string
}

type CommandMigrateAudit struct {
	Command

	Driver   string
	Set      string
	N        uint
	FullView bool
}

func New(
	fs *config.Filesystem,
	tpls *tpls.Templates,
	registry *plugins.Registry,
) *Runner {
	migrator := NewMigrator(fs)

	return &Runner{
		fs:       fs,
		tpls:     tpls,
		registry: registry,
		migrator: migrator,
	}
}

func (r *Runner) MustConfig() *config.Config {
	parser := config.NewParser(r.fs)
	parsed, err := parser.LoadAndParse()
	if err != nil {
		fmtx.WriteError(err.Error())
		os.Exit(1)
	}

	builder := NewBuilder(r.fs, parsed, r.registry)
	cfg, errs := builder.BuildConfig()
	if errs.HasErrors() {
		for _, err := range errs.Errors {
			fmtx.WriteError(err.Error())
		}
		os.Exit(1)
	}

	return cfg
}

func (r *Runner) MustDriver(registry *plugins.Registry, cfg *config.Config, name string) plugin.DriverInstance {
	definedDriver, ok := cfg.Drivers[name]
	if !ok {
		fmtx.WriteError("argument error: Driver not defined in config (metadata.name): %s", name)
		os.Exit(1)
	}
	driver, err := registry.Get(definedDriver.Spec.Driver)
	if err != nil {
		fmtx.WriteError(err.Error())
		os.Exit(1)
	}
	return plugin.DriverInstance{
		Driver: driver,
		Config: definedDriver,
	}
}

func (r *Runner) MustMigrationSet(cfg *config.Config, name string) *config.MigrationSet {
	set, ok := cfg.MigrationSets[name]
	if !ok {
		fmtx.WriteError("argument error: MigrationSet not found: %s", name)
		os.Exit(1)
	}
	return set
}

func (r *Runner) ExecuteMigrateFixUp(ctx context.Context, cmd *CommandMigrateFixUp) (*MigrateFixUpResult, error) {
	cfg := r.MustConfig()
	driver := r.MustDriver(r.registry, cfg, cmd.Driver)
	set := r.MustMigrationSet(cfg, cmd.Set)

	return r.migrator.MigrateFixUp(ctx, cfg, MigrateFixUpOptions{
		Driver:  driver,
		Set:     set,
		Version: cmd.Version,
        Comment: cmd.Comment,
	})
}

func (r *Runner) ExecuteMigrateFixDown(ctx context.Context, cmd *CommandMigrateFixDown) (*MigrateFixDownResult, error) {
	cfg := r.MustConfig()
	driver := r.MustDriver(r.registry, cfg, cmd.Driver)
	set := r.MustMigrationSet(cfg, cmd.Set)

	return r.migrator.MigrateFixDown(ctx, cfg, MigrateFixDownOptions{
		Driver:  driver,
		Set:     set,
		Version: cmd.Version,
        Comment: cmd.Comment,
	})
}

func (r *Runner) ExecuteMigrateUp(ctx context.Context, cmd *CommandMigrateUp) (*MigrateUpResult, error) {
	cfg := r.MustConfig()
	driver := r.MustDriver(r.registry, cfg, cmd.Driver)
	set := r.MustMigrationSet(cfg, cmd.Set)

	return r.migrator.MigrateUp(ctx, cfg, MigrateUpOptions{
		Driver: driver,
		Set:    set,
	})
}

func (r *Runner) ExecuteMigrateDown(ctx context.Context, cmd *CommandMigrateDown) (*MigrateDownResult, error) {
	cfg := r.MustConfig()
	driver := r.MustDriver(r.registry, cfg, cmd.Driver)
	set := r.MustMigrationSet(cfg, cmd.Set)

	return r.migrator.MigrateDown(ctx, cfg, MigrateDownOptions{
		Driver:  driver,
		Set:     set,
		Version: cmd.Version,
	})
}

func (r *Runner) ExecuteMigrateStatus(ctx context.Context, cmd *CommandMigrateStatus) (*MigrateStatusResult, error) {
	cfg := r.MustConfig()
	driver := r.MustDriver(r.registry, cfg, cmd.Driver)
	set := r.MustMigrationSet(cfg, cmd.Set)

	return r.migrator.MigrateStatus(ctx, cfg, MigrateStatusOptions{
		Driver: driver,
		Set:    set,
	})
}

func (r *Runner) ExecuteMigrateAudit(ctx context.Context, cmd *CommandMigrateAudit) (*MigrateAuditResult, error) {
	cfg := r.MustConfig()
	driver := r.MustDriver(r.registry, cfg, cmd.Driver)
	set := r.MustMigrationSet(cfg, cmd.Set)

	return r.migrator.MigrateAudit(ctx, cfg, MigrateAuditOptions{
		Driver:      driver,
		Set:         set,
		FilterLastN: cmd.N,
	})
}
