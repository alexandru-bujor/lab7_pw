package stats

import (
    "MegaMobileBack/pkg/db"
    "time"
)

type Service struct{}

func NewService() *Service { return &Service{} }

// TopProduct aggregates best sellers
type TopProduct struct {
    ProductTitle string  `json:"product_title"`
    Quantity     int     `json:"quantity"`
    Revenue      float64 `json:"revenue"`
}

// Summary aggregates dashboard stats
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
    TopProducts      []TopProduct `json:"top_products"`
}

// GetSummary fetches aggregated statistics
func (s *Service) GetSummary() (*Summary, error) {
    summary := &Summary{}

    // Revenue & profit
    db.DB.Table("orders").Select("COALESCE(SUM(sell_price),0) as total_revenue, COALESCE(SUM(profit),0) as total_profit").Scan(summary)

    // Orders counts by status
    db.DB.Table("orders").Select("COUNT(*)").Scan(&summary.TotalOrders)
    db.DB.Table("orders").Where("status = ?", "pending").Count(&summary.PendingOrders)
    db.DB.Table("orders").Where("status = ?", "processing").Count(&summary.ProcessingOrders)
    db.DB.Table("orders").Where("status = ?", "completed").Count(&summary.CompletedOrders)

    // Products
    db.DB.Table("products").Select("COUNT(*)").Scan(&summary.TotalProducts)

    // Customers
    db.DB.Table("clients").Select("COUNT(*)").Scan(&summary.TotalCustomers)

    // Inventory
    db.DB.Table("inventory_items").Where("status = ?", "in_stock").Count(&summary.InventoryInStock)
    db.DB.Table("inventory_items").Where("status = ?", "sold").Count(&summary.InventorySold)

    // Top products by revenue
    db.DB.Table("orders").
        Select("product_title, SUM(quantity) as quantity, SUM(sell_price) as revenue").
        Group("product_title").
        Order("revenue DESC").
        Limit(5).
        Scan(&summary.TopProducts)

    return summary, nil
}

// DailyStats represents statistics for a specific day
type DailyStats struct {
    Date              string  `json:"date"`
    ServiceRequests   int64   `json:"service_requests"`
    LombardRequests   int64   `json:"lombard_requests"`
    Orders            int64   `json:"orders"`
    Revenue           float64 `json:"revenue"`
    NewCustomers      int64   `json:"new_customers"`
}

// GetDailyStats fetches statistics for today and yesterday
func (s *Service) GetDailyStats() (today *DailyStats, yesterday *DailyStats, err error) {
    today = &DailyStats{}
    yesterday = &DailyStats{}

    // Get today's date (start and end of day)
    now := time.Now()
    todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    todayEnd := todayStart.Add(24 * time.Hour)
    
    // Get yesterday's date
    yesterdayStart := todayStart.Add(-24 * time.Hour)
    yesterdayEnd := todayStart

    // Today's stats
    db.DB.Table("service_requests").
        Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).
        Count(&today.ServiceRequests)
    
    db.DB.Table("lombard_requests").
        Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).
        Count(&today.LombardRequests)
    
    db.DB.Table("orders").
        Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).
        Count(&today.Orders)
    
    db.DB.Table("orders").
        Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).
        Select("COALESCE(SUM(sell_price), 0)").
        Scan(&today.Revenue)
    
    db.DB.Table("clients").
        Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).
        Count(&today.NewCustomers)
    
    today.Date = todayStart.Format("2006-01-02")

    // Yesterday's stats
    db.DB.Table("service_requests").
        Where("created_at >= ? AND created_at < ?", yesterdayStart, yesterdayEnd).
        Count(&yesterday.ServiceRequests)
    
    db.DB.Table("lombard_requests").
        Where("created_at >= ? AND created_at < ?", yesterdayStart, yesterdayEnd).
        Count(&yesterday.LombardRequests)
    
    db.DB.Table("orders").
        Where("created_at >= ? AND created_at < ?", yesterdayStart, yesterdayEnd).
        Count(&yesterday.Orders)
    
    db.DB.Table("orders").
        Where("created_at >= ? AND created_at < ?", yesterdayStart, yesterdayEnd).
        Select("COALESCE(SUM(sell_price), 0)").
        Scan(&yesterday.Revenue)
    
    db.DB.Table("clients").
        Where("created_at >= ? AND created_at < ?", yesterdayStart, yesterdayEnd).
        Count(&yesterday.NewCustomers)
    
    yesterday.Date = yesterdayStart.Format("2006-01-02")

    return today, yesterday, nil
}