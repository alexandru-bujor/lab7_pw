# Database migrations

Run these in order when deploying or upgrading the database. The app also uses GORM AutoMigrate on startup; these SQL files are for:

- Manual or CI deployment
- Data backfills (e.g. set_all_orders_on_book.sql)
- Ensuring schema matches the Go models

## Order

1. **20251229_create_service_categories.sql** – Creates `service_categories` table.
2. **20260127_001_schema_updates.sql** – Adds columns and tables: `categories.parent_id`, `products.old_price`, `settings`, `category_templates`, `accounting_type` / `is_visible_in_secondary` on orders, inventory_items, clients, products, services, lombard_requests, and `users.permissions`.
3. **set_all_products_on_book.sql** – Backfill: set existing products to `accounting_type = 'on_book'` where null/empty.
4. **set_all_orders_on_book.sql** – Backfill: set existing orders to `accounting_type = 'on_book'` where null/empty.

## Notes

- If a column or table already exists (e.g. from GORM AutoMigrate), you can ignore the corresponding "Duplicate column" or "Table already exists" error for that statement.
- Run with your MySQL client, e.g. `mysql -u user -p dbname < migrations/20260127_001_schema_updates.sql`
