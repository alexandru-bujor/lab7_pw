package inventory

// InventorySection represents a stock section/location
type InventorySection struct {
    ID   int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    Name string `gorm:"column:name;unique;not null" json:"name"`
}

func (InventorySection) TableName() string { return "inventory_sections" }
