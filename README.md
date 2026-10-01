# Elelibrary

Search Google Books and mark favorites. Made with ❤️, Go and Vue

## Run with Docker

Requires Docker with Compose.

```sh
cp .env.example .env   # then set GOOGLE_BOOKS_API_KEY
docker compose up --build
```

Open http://localhost:8000.

The database schema in `be/migrations/` is applied only when the `db` volume is first created. After adding a migration, reset your local data with `docker compose down -v`.

## Run without Docker

Backend (needs a Postgres with `be/migrations/*.sql` applied):

```sh
cd be
cp .env.example .env   # then fill in
set -a; source .env; set +a
go run ./cmd/api
```

Frontend (proxies `/api` to `localhost:8080`):

```sh
cd fe
npm install
npm run dev
```
