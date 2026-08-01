package turso

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rcovery/go-stock-control/internal/part"
	"github.com/rcovery/go-stock-control/internal/part/errs"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(DB *sql.DB) *Repository {
	return &Repository{
		DB,
	}
}

const partColumns = `
	id, name, category, current_stock, minimum_stock,
	average_daily_sales, lead_time_days, unit_cost, criticality_level
`

func scanPart(row interface{ Scan(...any) error }) (part.Part, error) {
	var p part.Part

	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Category,
		&p.CurrentStock,
		&p.MinimumStock,
		&p.AverageDailySales,
		&p.LeadTimeDays,
		&p.UnitCost,
		&p.CriticalityLevel,
	)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (r *Repository) List(ctx context.Context) ([]part.Part, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT `+partColumns+`
		FROM parts
		ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("part list: %w", err)
	}
	defer rows.Close()

	var parts []part.Part
	for rows.Next() {
		p, scanErr := scanPart(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("part list: %w", scanErr)
		}
		parts = append(parts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("part list: %w", err)
	}

	return parts, nil
}

func (r *Repository) ListByCategory(ctx context.Context, category string) ([]part.Part, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT `+partColumns+`
		FROM parts
		WHERE category = ?
		ORDER BY name
	`, category)
	if err != nil {
		return nil, fmt.Errorf("part list by category: %w", err)
	}
	defer rows.Close()

	var parts []part.Part
	for rows.Next() {
		p, scanErr := scanPart(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("part list by category: %w", scanErr)
		}
		parts = append(parts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("part list by category: %w", err)
	}

	return parts, nil
}

func (r *Repository) Create(ctx context.Context, p part.Part) error {
	_, insertionErr := r.DB.ExecContext(ctx, `
		INSERT INTO parts
		(id, name, category, current_stock, minimum_stock,
		 average_daily_sales, lead_time_days, unit_cost, criticality_level)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, p.ID, p.Name, p.Category, p.CurrentStock, p.MinimumStock,
		p.AverageDailySales, p.LeadTimeDays, p.UnitCost, p.CriticalityLevel,
	)
	if insertionErr != nil {
		return errs.NotCreatedErr.New(insertionErr.Error())
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, p part.Part) error {
	result, execErr := r.DB.ExecContext(ctx, `
		UPDATE parts
		SET name = ?,
			category = ?,
			current_stock = ?,
			minimum_stock = ?,
			average_daily_sales = ?,
			lead_time_days = ?,
			unit_cost = ?,
			criticality_level = ?
		WHERE id = ?
	`, p.Name, p.Category, p.CurrentStock, p.MinimumStock,
		p.AverageDailySales, p.LeadTimeDays, p.UnitCost, p.CriticalityLevel,
		p.ID,
	)
	if execErr != nil {
		return fmt.Errorf("part update: %w", execErr)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("part update rows affected: %w", err)
	}
	if rows == 0 {
		return errs.NotFoundError.New(fmt.Sprintf("part %q not found", p.ID))
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id part.ID) error {
	result, execErr := r.DB.ExecContext(ctx, `
		DELETE FROM parts
		WHERE id = ?
	`, id)
	if execErr != nil {
		return fmt.Errorf("part delete: %w", execErr)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("part delete rows affected: %w", err)
	}
	if rows == 0 {
		return errs.NotFoundError.New(fmt.Sprintf("part %q not found", id))
	}

	return nil
}
