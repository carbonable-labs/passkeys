## Getting Started

Download those utils :

- [ ] [just](https://github.com/casey/just)
- [ ] [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
- [ ] [atlas](https://atlasgo.io/docs)

## Env vars

- DATABASE_URL
- RPC_URL
- RPC_API_KEY
- RELYING_PARTY_ORIGINS
- RELYING_PARTY_NAME
- RELYING_PARTY_ID

## Migrations

Schema is defined in `sql/schema.sql` and migrations are defined in `sql/migrations`.

To create a new migration, edit `sql/schema.sql` then run :

```bash
just migrate_diff migration_name
```

To apply migration, run :

```bash
just migrate
```

Whenever you modified files in `sql` run :

```bash
sqlc generate
```
