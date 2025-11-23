package run

import (
	"context"
	"fmt"

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
	logger   *fmtx.Logger
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
	logger *fmtx.Logger,
) *Runner {
	migrator := NewMigrator(fs, logger)

	return &Runner{
		fs:       fs,
		tpls:     tpls,
		registry: registry,
		migrator: migrator,
		logger:   logger,
	}
}

func (r *Runner) UseConfig() (*config.Config, error) {
	parser := config.NewParser(r.fs, r.logger)
	parsed, err := parser.LoadAndParse()
	if err != nil {
		return nil, err
	}

	builder := NewBuilder(r.fs, parsed, r.registry)
	cfg, errs := builder.BuildConfig()
	if errs.HasErrors() {
		for _, err := range errs.Errors {
			r.logger.WriteError(err.Error())
		}
		return nil, fmt.Errorf("invalid configuration")
	}

	return cfg, nil
}

func (r *Runner) UseDriver(registry *plugins.Registry, cfg *config.Config, name string) (*plugin.DriverInstance, error) {
	definedDriver, ok := cfg.Drivers[name]
	if !ok {
		return nil, fmt.Errorf("argument error: Driver not defined in config (metadata.name): %s", name)
	}
	driver, err := registry.Get(definedDriver.Spec.Driver)
	if err != nil {
		return nil, err
	}
	return &plugin.DriverInstance{
		Driver: driver,
		Config: definedDriver,
	}, nil
}

func (r *Runner) UseMigrationSet(cfg *config.Config, name string) (*config.MigrationSet, error) {
	set, ok := cfg.MigrationSets[name]
	if !ok {
		return nil, fmt.Errorf("argument error: MigrationSet not found: %s", name)
	}
	return set, nil
}

func (r *Runner) ExecuteMigrateFixUp(ctx context.Context, cmd *CommandMigrateFixUp) (*MigrateFixUpResult, error) {
	cfg, err := r.UseConfig()
	if err != nil {
		return nil, err
	}
	driver, err := r.UseDriver(r.registry, cfg, cmd.Driver)
    if err != nil {
        return nil, err
    }
	set, err := r.UseMigrationSet(cfg, cmd.Set)
    if err != nil {
        return nil, err
    }

	return r.migrator.MigrateFixUp(ctx, cfg, MigrateFixUpOptions{
		Driver:  driver,
		Set:     set,
		Version: cmd.Version,
		Comment: cmd.Comment,
	})
}

func (r *Runner) ExecuteMigrateFixDown(ctx context.Context, cmd *CommandMigrateFixDown) (*MigrateFixDownResult, error) {
	cfg, err := r.UseConfig()
	if err != nil {
		return nil, err
	}
	driver, err := r.UseDriver(r.registry, cfg, cmd.Driver)
    if err != nil {
        return nil, err
    }
	set, err := r.UseMigrationSet(cfg, cmd.Set)
    if err != nil {
        return nil, err
    }

	return r.migrator.MigrateFixDown(ctx, cfg, MigrateFixDownOptions{
		Driver:  driver,
		Set:     set,
		Version: cmd.Version,
		Comment: cmd.Comment,
	})
}

func (r *Runner) ExecuteMigrateUp(ctx context.Context, cmd *CommandMigrateUp) (*MigrateUpResult, error) {
	cfg, err := r.UseConfig()
	if err != nil {
		return nil, err
	}
	driver, err := r.UseDriver(r.registry, cfg, cmd.Driver)
    if err != nil {
        return nil, err
    }
	set, err := r.UseMigrationSet(cfg, cmd.Set)
    if err != nil {
        return nil, err
    }

	return r.migrator.MigrateUp(ctx, cfg, MigrateUpOptions{
		Driver: driver,
		Set:    set,
	})
}

func (r *Runner) ExecuteMigrateDown(ctx context.Context, cmd *CommandMigrateDown) (*MigrateDownResult, error) {
	cfg, err := r.UseConfig()
	if err != nil {
		return nil, err
	}
	driver, err := r.UseDriver(r.registry, cfg, cmd.Driver)
    if err != nil {
        return nil, err
    }
	set, err := r.UseMigrationSet(cfg, cmd.Set)
    if err != nil {
        return nil, err
    }

	return r.migrator.MigrateDown(ctx, cfg, MigrateDownOptions{
		Driver:  driver,
		Set:     set,
		Version: cmd.Version,
	})
}

func (r *Runner) ExecuteMigrateStatus(ctx context.Context, cmd *CommandMigrateStatus) (*MigrateStatusResult, error) {
	cfg, err := r.UseConfig()
	if err != nil {
		return nil, err
	}
	driver, err := r.UseDriver(r.registry, cfg, cmd.Driver)
    if err != nil {
        return nil, err
    }
	set, err := r.UseMigrationSet(cfg, cmd.Set)
    if err != nil {
        return nil, err
    }

	return r.migrator.MigrateStatus(ctx, cfg, MigrateStatusOptions{
		Driver: driver,
		Set:    set,
	})
}

func (r *Runner) ExecuteMigrateAudit(ctx context.Context, cmd *CommandMigrateAudit) (*MigrateAuditResult, error) {
	cfg, err := r.UseConfig()
	if err != nil {
		return nil, err
	}
	driver, err := r.UseDriver(r.registry, cfg, cmd.Driver)
    if err != nil {
        return nil, err
    }
	set, err := r.UseMigrationSet(cfg, cmd.Set)
    if err != nil {
        return nil, err
    }

	return r.migrator.MigrateAudit(ctx, cfg, MigrateAuditOptions{
		Driver:      driver,
		Set:         set,
		FilterLastN: cmd.N,
	})
}
