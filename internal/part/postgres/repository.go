package postgres

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
		return nil, errs.NotFoundError.New(fmt.Sprintf("List: %v", err))
	}
	defer rows.Close()

	var parts []part.Part
	for rows.Next() {
		p, scanErr := scanPart(rows)
		if scanErr != nil {
			return nil, errs.NotFoundError.New(fmt.Sprintf("List: %v", scanErr))
		}
		parts = append(parts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.NotFoundError.New(fmt.Sprintf("List: %v", err))
	}

	return parts, nil
}

func (r *Repository) ListByCategory(ctx context.Context, category string) ([]part.Part, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT `+partColumns+`
		FROM parts
		WHERE category = $1
		ORDER BY name
	`, category)
	if err != nil {
		return nil, errs.NotFoundError.New(fmt.Sprintf("ListByCategory: %v", err))
	}
	defer rows.Close()

	var parts []part.Part
	for rows.Next() {
		p, scanErr := scanPart(rows)
		if scanErr != nil {
			return nil, errs.NotFoundError.New(fmt.Sprintf("ListByCategory: %v", scanErr))
		}
		parts = append(parts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.NotFoundError.New(fmt.Sprintf("ListByCategory: %v", err))
	}

	return parts, nil
}

func (r *Repository) Create(ctx context.Context, p part.Part) error {
	_, insertionErr := r.DB.ExecContext(ctx, `
		INSERT INTO parts
		(id, name, category, current_stock, minimum_stock,
		 average_daily_sales, lead_time_days, unit_cost, criticality_level)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, p.ID, p.Name, p.Category, p.CurrentStock, p.MinimumStock,
		p.AverageDailySales, p.LeadTimeDays, p.UnitCost, p.CriticalityLevel,
	)
	if insertionErr != nil {
		return errs.NotCreatedErr.New(insertionErr.Error())
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, p part.Part) error {
	_, updateErr := r.DB.ExecContext(ctx, `
		UPDATE parts
		SET name = $2,
			category = $3,
			current_stock = $4,
			minimum_stock = $5,
			average_daily_sales = $6,
			lead_time_days = $7,
			unit_cost = $8,
			criticality_level = $9
		WHERE id = $1
	`, p.ID, p.Name, p.Category, p.CurrentStock, p.MinimumStock,
		p.AverageDailySales, p.LeadTimeDays, p.UnitCost, p.CriticalityLevel,
	)
	if updateErr != nil {
		return errs.NotFoundError.New(fmt.Sprintf("Update: %v", updateErr))
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id part.ID) error {
	_, deleteErr := r.DB.ExecContext(ctx, `
		DELETE FROM parts
		WHERE id = $1
	`, id)
	if deleteErr != nil {
		return errs.NotFoundError.New(fmt.Sprintf("Delete: %v", deleteErr))
	}

	return nil
}
