package expense

import (
	"context"
	"fmt"

	"financial-backend/internal/entities"
	"financial-backend/internal/models"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, expense *entities.Expense) error {
	return r.db.WithContext(ctx).Create(expense).Error
}

func (r *repository) Update(ctx context.Context, expense *entities.Expense) error {
	return r.db.WithContext(ctx).Save(expense).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Expense{}).Error
}

func (r *repository) Get(ctx context.Context, id string) (*entities.Expense, error) {
	var expense entities.Expense
	if err := r.db.WithContext(ctx).First(&expense, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar despesa: %v", err)
	}
	return &expense, nil
}

func (r *repository) List(ctx context.Context, description, expenseType, category, budgetId, recurrecy, method string, page models.PageRequest) (expenses []*entities.Expense, count int64, err error) {
	query := r.db.WithContext(ctx)

	if description != "" {
		query = query.Where("description like ?", "%"+description+"%")
	}

	if expenseType != "" {
		query = query.Where("type = ?", expenseType)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if budgetId != "" {
		query = query.Where("budget_id = ?", budgetId)
	}

	if recurrecy != "" {
		query = query.Where("recurrecy = ?", recurrecy)
	}

	if method != "" {
		query = query.Where("method = ?", method)
	}

	if recurrecy != "" {
		query = query.Where("recurrecy = ?", recurrecy)
	}

	if recurrecy != "" {
		query = query.Where("recurrecy = ?", recurrecy)
	}

	if err := query.Preload("Budget").Offset(page.Offset()).Limit(int(page.Limit)).Find(&expenses).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao listar despesas: %v", err)
	}

	if err := query.Model(&entities.Expense{}).Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao contar Despesas: %v", err)
	}

	return expenses, count, nil

}

func (r *repository) GetExpensesWithoutMovimentInMonth(ctx context.Context) (expenses []*entities.Expense, err error) {
	query := `with expense as (select *
                 from expenses
                 where method != 'credit_card'
                   and type = 'recurring'
                   and (end_date is null or end_date >= current_date)
                   and budget_id is not null)
select e.*,
       b.description AS "budget__description",
       b.amount      AS "budget__amount",
       b.end_date    AS "budget__end_date",
       b.created_at  AS "budget__created_at",
       b.updated_at  AS "budget__updated_at"
from budget_movements bm
         right join expense e on (bm.origin = e.id and bm.month = cast(to_char(now() :: date, 'MM') as numeric) and
                                  bm.year = cast(to_char(now() :: date, 'YYYY') as numeric))
         join public.budgets b on e.budget_id = b.id
where bm.id is null
and e.budget_id is not null
`

	if err := r.db.WithContext(ctx).Raw(query).Preload("Budget").Find(&expenses).Error; err != nil {
		return make([]*entities.Expense, 0), err
	}

	return
}
func (r *repository) SummaryByMonth(ctx context.Context, month int, year int) (float64, error) {
	query := "select sum(amount) as total_expense from expenses  where start_date >= ? and (end_date <= ? or end_date is null)"

	startDate := fmt.Sprintf("%d-%02d-01", year, month)
	endDate := fmt.Sprintf("%d-%02d-01", year, month+1)

	var amount float64
	if err := r.db.WithContext(ctx).Raw(query, startDate, endDate).Scan(&amount).Error; err != nil {
		return 0, err
	}
	return amount, nil
}
