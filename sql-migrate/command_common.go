package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"time"

	"github.com/nicolai86/sql-migrate"
)

func GenerateMigrationTemplate(name string) error {
	env, err := GetEnvironment()
	if err != nil {
		return fmt.Errorf("Could not parse config: %s", err)
	}

	source := migrate.FileMigrationSource{
		Dir: env.Dir,
	}

	template, err := filepath.Abs("../templates/migration.sql")
	if err != nil {
		return fmt.Errorf("Could not read template file: %s", err)
	}
	b, err := ioutil.ReadFile(template)
	if err != nil {
		return fmt.Errorf("Could not read template file: %s", err)
	}

	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s.sql", timestamp, name)
	file := path.Join(source.Dir, fileName)
	err = ioutil.WriteFile(file, b, 0644)
	if err != nil {
		return fmt.Errorf("Could not write template file: %s", err)
	}
	ui.Output(fmt.Sprintf("Created new Migration: %s", file))

	return nil
}

func ApplyMigrations(dir migrate.MigrationDirection, dryrun bool, limit int) error {
	env, err := GetEnvironment()
	if err != nil {
		return fmt.Errorf("Could not parse config: %s", err)
	}

	db, dialect, err := GetConnection(env)
	if err != nil {
		return err
	}

	source := migrate.FileMigrationSource{
		Dir: env.Dir,
	}

	if dryrun {
		migrations, err := migrate.PlanMigration(db, dialect, source, dir, limit)
		if err != nil {
			return fmt.Errorf("Cannot plan migration: %s", err)
		}

		for _, m := range migrations {
			PrintMigration(m, dir)
		}
	} else {
		n, err := migrate.ExecMax(db, dialect, source, dir, limit)
		if err != nil {
			return fmt.Errorf("Migration failed: %s", err)
		}

		if n == 1 {
			ui.Output("Applied 1 migration")
		} else {
			ui.Output(fmt.Sprintf("Applied %d migrations", n))
		}
	}

	return nil
}

func PrintMigration(m *migrate.PlannedMigration, dir migrate.MigrationDirection) {
	if dir == migrate.Up {
		ui.Output(fmt.Sprintf("==> Would apply migration %s (up)", m.Id))
		for _, q := range m.Up {
			ui.Output(q)
		}
	} else if dir == migrate.Down {
		ui.Output(fmt.Sprintf("==> Would apply migration %s (down)", m.Id))
		for _, q := range m.Down {
			ui.Output(q)
		}
	} else {
		panic("Not reached")
	}
}
