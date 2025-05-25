# 💰 Wallet Backend

The **wallet-backend** is a modular Go-based backend for a digital wallet service. It supports:

- Wallet management and transactions
- Webhook handling
- Redis-based caching
- PostgreSQL for persistent storage
- Loki & Promtail for logging and observability

Built with **Go**, uses **GORM** for ORM, and integrates with **Docker** for easy deployment.

---

## 🚧 Prerequisites

Before you run this project, make sure:

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) is installed and running.
- Ports `5432` (PostgreSQL) and `6379` (Redis) are **not in use** by other processes on your machine.

> 💡 To check if the ports are in use:

### On **Windows (PowerShell)**

```powershell
netstat -ano | findstr :5432
netstat -ano | findstr :6379
```

### On **Mac/Linux**

```bash
lsof -i :5432
lsof -i :6379
```

---

## 🧰 How to Run

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/wallet-backend.git
cd wallet-backend
```

### 2. Build and run with Docker

```bash
docker-compose up --build
```

This will:

- Start the Go application
- Launch PostgreSQL on port `5432`
- Launch Redis on port `6379`
- Launch Loki + Promtail for logs

---

## ⚙️ Environment Variables

The backend uses the following environment variables, configured inside `docker-compose.yml`:

| Variable     | Description                    | Example                                                       |
|--------------|--------------------------------|---------------------------------------------------------------|
| `DB_URL`     | PostgreSQL connection string   | `postgres://user:pass@postgres:5432/dbname?sslmode=disable`  |
| `REDIS_ADDR` | Redis host and port            | `redis:6379`                                                  |

---

## 🛠 Features

- 🔐 Auth-ready user model
- 💸 Wallet + Transaction models
- 🧠 Redis-based session or cache layer
- 📬 Webhook listener and logger
- 📊 Loki and Promtail integration for logs
- 🧪 Auto-migration via GORM

---

## 🧹 Cleanup

To stop and remove containers:

```bash
docker-compose down
```

To remove volumes (Postgres data, etc):

```bash
docker-compose down -v
```

---

## 🐞 Troubleshooting

- **Database errors**: Ensure no app is already using port `5432`.
- **Redis errors**: Ensure Redis is not already installed locally and bound to `6379`.
- **UUID errors**: Make sure `uuid-ossp` extension is added (you can modify your Docker `init.sql` to include `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

---


## 📁 Project Structure

```
.
├── auth
│   ├── handlers
│   │   └── auth_handler.go
│   ├── middleware
│   │   └── auth.go
│   ├── models
│   │   ├── response.go
│   │   └── user.go
│   ├── routes
│   │   └── auth_routes.go
│   └── services
│       ├── auth.go
│       └── session.go
├── cmd
│   └── main.go
├── docker-compose.yml
├── Dockerfile
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── logs  [error opening dir]
├── promtail-config.yaml  [error opening dir]
├── README.md
├── users_simulated.csv
├── utils
│   └── db.go
├── wallet
│   ├── handlers
│   │   ├── simulation_handler.go
│   │   └── wallet_handler.go
│   ├── models
│   │   ├── fee.go
│   │   ├── simulate_option.go
│   │   └── transaction.go
│   ├── routes
│   │   └── transaction_routes.go
│   ├── services
│   │   ├── simulation_services.go
│   │   └── wallet_services.go
│   └── utils
│       └── sim_status.go
├── wallet-backend
└── webhook
    ├── handlers
    │   └── webhook_handler.go
    ├── interfaces
    │   └── interface_webservcies.go
    ├── middleware
    │   └── hmac_validate.go
    ├── mocks
    │   └── mock_webservice.go
    ├── models
    │   └── webhook_model.go
    ├── processor
    ├── queue
    ├── routes
    │   └── webhook_route.go
    ├── services
    │   └── webhook_services.go
    ├── test
    │   ├── main.go
    │   └── payload.json
    └── utils
        └── jsonb.go
```

---

## ✅ After Setup

Once the containers are up and running, visit:

👉 **[http://127.0.0.1:8080/swagger/](http://127.0.0.1:8080/swagger/)**

to explore the **interactive API documentation** via Swagger UI.

---



## 📜 License

MIT License. See [LICENSE](./LICENSE).

---

## 👨‍💻 Maintainer

Nazrawi Gedion  
Feel free to open issues or PRs!
