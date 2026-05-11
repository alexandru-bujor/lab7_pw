package stats

import (
	"MegaMobileBack/pkg/db"
	"sync"
	"time"
)

type Service struct{}

func NewService() *Service { return &Service{} }

type TopProduct struct {
	ProductTitle string  `json:"product_title"`
	Quantity     int     `json:"quantity"`
	Revenue      float64 `json:"revenue"`
}

type Summary struct {
	TotalRevenue     float64      `json:"total_revenue"`
	TotalProfit      float64      `json:"total_profit"`
	TotalOrders      int64        `json:"total_orders"`
	PendingOrders    int64        `json:"pending_orders"`
	ProcessingOrders int64        `json:"processing_orders"`
	CompletedOrders  int64        `json:"completed_orders"`
	TotalProducts    int64        `json:"total_products"`
	TotalCustomers   int64        `json:"total_customers"`
	InventoryInStock int64        `json:"inventory_in_stock"`
	InventorySold    int64        `json:"inventory_sold"`
	TopProducts      []TopProduct `gorm:"-" json:"top_products"`
}

func (s *Service) GetSummary() (*Summary, error) {
	summary := &Summary{}

	type queryTask struct {
		fn func() error
	}

	tasks := []queryTask{
		{func() error {
			return db.DB.Table("orders").Select("COALESCE(SUM(sell_price),0) as total_revenue, COALESCE(SUM(profit),0) as total_profit").Scan(summary).Error
		}},
		{func() error { return db.DB.Table("orders").Count(&summary.TotalOrders).Error }},
		{func() error {
			return db.DB.Table("orders").Where("status = ?", "pending").Count(&summary.PendingOrders).Error
		}},
		{func() error {
			return db.DB.Table("orders").Where("status = ?", "processing").Count(&summary.ProcessingOrders).Error
		}},
		{func() error {
			return db.DB.Table("orders").Where("status = ?", "completed").Count(&summary.CompletedOrders).Error
		}},
		{func() error { return db.DB.Table("products").Count(&summary.TotalProducts).Error }},
		{func() error { return db.DB.Table("clients").Count(&summary.TotalCustomers).Error }},
		{func() error {
			return db.DB.Table("inventory_items").Where("status = ?", "in_stock").Count(&summary.InventoryInStock).Error
		}},
		{func() error {
			return db.DB.Table("inventory_items").Where("status = ?", "sold").Count(&summary.InventorySold).Error
		}},
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		go func(t queryTask) {
			defer wg.Done()
			if err := t.fn(); err != nil {
				errCh <- err
			}
		}(task)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	err := db.DB.Table("orders").
		Select("product_title, SUM(quantity) as quantity, SUM(sell_price) as revenue").
		Group("product_title").
		Order("revenue DESC").
		Limit(5).
		Scan(&summary.TopProducts).Error

	if err != nil {
		return nil, err
	}

	return summary, nil
}

type DailyStats struct {
	Date            string  `json:"date"`
	ServiceRequests int64   `json:"service_requests"`
	LombardRequests int64   `json:"lombard_requests"`
	Orders          int64   `json:"orders"`
	Revenue         float64 `json:"revenue"`
	NewCustomers    int64   `json:"new_customers"`
}

func (s *Service) GetDailyStats() (today *DailyStats, yesterday *DailyStats, err error) {
	today = &DailyStats{}
	yesterday = &DailyStats{}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)
	yesterdayStart := todayStart.Add(-24 * time.Hour)
	yesterdayEnd := todayStart

	var wg sync.WaitGroup
	wg.Add(10)

	go func() {
		defer wg.Done()
		db.DB.Table("service_requests").Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).Count(&today.ServiceRequests)
	}()
	go func() {
		defer wg.Done()
		db.DB.Table("lombard_requests").Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).Count(&today.LombardRequests)
	}()
	go func() {
		defer wg.Done()
		db.DB.Table("orders").Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).Count(&today.Orders)
	}()
	go func() {
		defer wg.Done()
		db.DB.Table("orders").Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).Select("COALESCE(SUM(sell_price), 0)").Scan(&today.Revenue)
	}()
	go func() {
		defer wg.Done()
		db.DB.Table("clients").Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).Count(&today.NewCustomers)
	}()
	go func() {
		defer wg.Done()
		db.DB.Table("service_requests").Where("created_at >= ? AND created_at < ?", yesterdayStart, yesterdayEnd).Count(&yesterday.ServiceRequests)
	}()
	go func() {
		defer wg.Done()
		db.DB.Table("lombard_requests").Where("created_at >= ? AND created_at < ?", yesterdayStart, yesterdayEnd).Count(&yesterday.LombardRequests)
	}()
	go func() {
		defer wg.Done()
		db.DB.Table("orders").Where("created_at >= ? AND created_at < ?", yesterdayStart, yesterdayEnd).Count(&yesterday.Orders)
	}()
	go func() {
		defer wg.Done()
		db.DB.Table("orders").Where("created_at >= ? AND created_at < ?", yesterdayStart, yesterdayEnd).Select("COALESCE(SUM(sell_price), 0)").Scan(&yesterday.Revenue)
	}()
	go func() {
		defer wg.Done()
		db.DB.Table("clients").Where("created_at >= ? AND created_at < ?", yesterdayStart, yesterdayEnd).Count(&yesterday.NewCustomers)
	}()

	wg.Wait()

	today.Date = todayStart.Format("2006-01-02")
	yesterday.Date = yesterdayStart.Format("2006-01-02")

	return today, yesterday, nil
}

type DailyHistoricalStats struct {
	Date         string  `json:"date"`
	Orders       int64   `json:"orders"`
	Revenue      float64 `json:"revenue"`
	NewCustomers int64   `json:"new_customers"`
}

func (s *Service) GetHistoricalDailyStats(days int) ([]DailyHistoricalStats, error) {
	var stats []DailyHistoricalStats

	now := time.Now()
	statsMap := make(map[string]*DailyHistoricalStats)

	// Go back 'days-1' to include today and the past days to equal total 'days'
	startDate := now.AddDate(0, 0, -days+1)
	startDateStr := startDate.Format("2006-01-02")

	for i := 0; i < days; i++ {
		dateStr := startDate.AddDate(0, 0, i).Format("2006-01-02")
		statsMap[dateStr] = &DailyHistoricalStats{
			Date:         dateStr,
			Orders:       0,
			Revenue:      0,
			NewCustomers: 0,
		}
	}

	var orderStats []struct {
		Date    string
		Orders  int64
		Revenue float64
	}

	err := db.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(id) as orders,
			COALESCE(SUM(sell_price), 0) as revenue
		FROM orders
		WHERE created_at >= ? AND status != 'cancelled'
		GROUP BY DATE(created_at)
	`, startDateStr).Scan(&orderStats).Error

	if err != nil {
		return nil, err
	}

	for _, os := range orderStats {
		if stat, ok := statsMap[os.Date]; ok {
			stat.Orders = os.Orders
			stat.Revenue = os.Revenue
		}
	}

	var clientStats []struct {
		Date         string
		NewCustomers int64
	}

	err = db.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(id) as new_customers
		FROM clients
		WHERE created_at >= ?
		GROUP BY DATE(created_at)
	`, startDateStr).Scan(&clientStats).Error

	if err != nil {
		return nil, err
	}

	for _, cs := range clientStats {
		if stat, ok := statsMap[cs.Date]; ok {
			stat.NewCustomers = cs.NewCustomers
		}
	}

	for i := 0; i < days; i++ {
		dateStr := startDate.AddDate(0, 0, i).Format("2006-01-02")
		stats = append(stats, *statsMap[dateStr])
	}

	return stats, nil
}
