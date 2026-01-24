package servicecategory

 import (
 	"database/sql/driver"
 	"encoding/json"
 	"fmt"
 )

 type MultilingualName struct {
 	RO string `json:"ro"`
 	RU string `json:"ru"`
 	EN string `json:"en"`
 }

 // Implement driver.Valuer for GORM/SQL
 func (m MultilingualName) Value() (driver.Value, error) {
 	return json.Marshal(m)
 }

 // Implement sql.Scanner for GORM/SQL
 func (m *MultilingualName) Scan(value interface{}) error {
 	switch v := value.(type) {
 	case []byte:
 		return json.Unmarshal(v, m)
 	case string:
 		return json.Unmarshal([]byte(v), m)
 	default:
 		return fmt.Errorf("unsupported type: %T", value)
 	}
 }

 type ServiceCategory struct {
 	ID   uint            `gorm:"primaryKey" json:"id"`
 	Name MultilingualName `gorm:"type:json" json:"name"`
 }

 func (ServiceCategory) TableName() string {
 	return "service_categories"
 }
