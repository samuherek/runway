# Runway

The runway application for people to see how they are doing financialy into the future. :)

## To Run

### Go

```bash
air
```

### Templ

```bash
templ generate --watch
```

### TailwindCSS

```bash
cd tailwind
npm run watch-css
```

## Envs

```
NONE SO FAR
```


### DB

Using `migrate` with necessary ENV variable `export POSTGRESQL_URL=postgres://postgres:password@localhost:5432/runway?sslmode=disable`

- Drop database `migrate --database ${POSTGRESQL_URL} -path db/migrations drop` 
- Migrate up `migrate --database ${POSTGRESQL_URL} -path db/migrations up` 
- Create migration `migrate create -ext sql -dir db/migrations -seq create_users_table`
