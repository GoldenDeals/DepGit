# Database Migrations

This directory contains SQL migration files that are automatically applied to the database when the application starts.

## Migration File Format

- Migration files should be named in the format `NNN_description.sql` where `NNN` is a three-digit number starting from 001.
- Migrations are applied in order based on the file name.
- Each migration file should contain valid SQL statements.
- Statements are separated by semicolons (;).

## Creating a New Migration

To create a new migration:

1. Create a new SQL file in this directory with the next sequential number.
2. Add the necessary SQL statements to the file.
3. The migration will be automatically applied the next time the application starts.

## Example

```sql
-- 002_add_email_index.sql
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
```

## Notes

- Migrations are applied in a transaction - if any statement fails, the entire migration is rolled back.
- Once a migration is applied, it is recorded in the `migrations` table and will not be applied again.
- Always test your migrations in a development environment before deploying to production. 