# E2E tests (Playwright)

Run against a **running** backend on `http://localhost:8080`.

## 1. Start backend

```bash
# From repo root
docker-compose up -d db
go run .   # or: copy .env.example to .env and set DB_* so DB connects
```

## 2. Run E2E

```bash
cd e2e
npm install
npx playwright test
```

Optional: use another base URL:

```bash
BASE_URL=http://localhost:8080 npx playwright test
```

## Tests

- **api.spec.ts** – `/ping`, GraphiQL GET, `hello` query, `createUser` mutation
- **exploratory.spec.ts** – GraphiQL UI loads
