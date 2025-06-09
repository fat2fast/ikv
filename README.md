# IKV – Inkventory Service

> **Inkventory** là pet‑project Go mô hình **Modular Monolith** (tách module, sẵn sàng thành micro‑service).
>
> • CSDL mặc định **PostgreSQL 16 + GORM**.

---

## 📂 Cây thư mục chính

```text
ikv/
├── app/                       # <–– toàn bộ mã nguồn Go đặt trong đây
│   ├── go.work               # gom các module con
│   ├── shared/               # logger, config, errors, eventbus
│   │   └── go.mod
│   ├── modules/              # mỗi bounded‑context 1 thư mục
│   │   ├── category/
│   │   │   ├── module.go
│   │   │   ├── model/
│   │   │   ├── service/
│   │   │   └── infras/
│   │   │       ├── repository/
│   │   │       └── controller/
│   │   ├── inventory/
│   │   ├── ordering/
│   │   └── customer/
│   └── cmd/
│       ├── api-gateway/      # gom route HTTP/gRPC
│       ├── worker/           # job nền
│       └── migrate/          # migration runner
├── build/
│   └── Dockerfile            # image multi‑stage build
├── docs/
│   └── flow.png              # sơ đồ kiến trúc/png khác
├── vendor/                   # (tuỳ chọn) module vendoring
├── docker-compose.yml        # Postgres, NATS, Jaeger, Adminer, ikv
├── Dockerfile                # root image tiện CI
├── LICENSE
└── README.md                 # bạn đang đọc
```

---

## 🔄 Luồng phụ thuộc trong 1 module

```
 model  ←── service  ←── infras/repository (GORM Postgres)
                        ↑
                   infras/controller (HTTP / gRPC)
```

Không module nào được import trực tiếp code của module khác; chỉ giao tiếp qua **interface** repository/service hoặc eventbus.

## 🔍 Sơ đồ tổng thể (Mermaid)

```mermaid
graph TD
  subgraph "API Gateway (cmd/api-gateway)"
    GIN[Gin Router] --> CAT_HTTP[Category HTTP Controller]
    GRPC[gRPC Server] --> CAT_GRPC[Category gRPC Controller]
  end

  subgraph "Category Module"
    CAT_HTTP --> CAT_SVC[Category Service]
    CAT_GRPC --> CAT_SVC
    CAT_SVC --> CAT_REPO[(PostgreSQL via GORM)]
  end

  classDef repo fill:#fff2,stroke:#333;
  class CAT_REPO repo;
```
![flow](./docs/flow.png)
---

## 🛠️ Thiết lập nhanh

```bash
# 1. Copy sample env
cp env.example .env

# 2. Khởi động stack dev (Postgres, Jaeger, ...)
docker-compose up -d

# 3. Chạy migration & gateway
go run ./app/cmd/migrate
cd app/cmd/api-gateway
go run .
```

DSLN Postgres (.env):

```
DB_DSN=postgres://ikv:ikv123@localhost:5432/ikv?sslmode=disable
```

---

## 🗺️ Roadmap ngắn

| Sprint | Mục tiêu           | Tasks                                                  |
| ------ | ------------------ | ------------------------------------------------------ |
| 0      | Scaffold           | `go.work init`, cấu hình Postgres trong docker‑compose |
| 1      | Category CRUD      | Repo GORM, Service, HTTP `GET/POST /v1/categories`     |
| 2      | Inventory          | Stock CRUD, low‑stock worker                           |
| 3      | Ordering           | TX giảm stock, sự kiện OrderPlaced                     |
| 4      | Observability & CI | Jaeger, Prometheus, GitLab CI                          |
| 5      | Split Inventory    | Tách repo, viết client gRPC                            |

---

## ☑️ Kiểm thử

```bash
go test ./app/...
```

* Unit test: mock repository.
* Integration: docker‑compose + gọi HTTP.

---

## 📜 License

MIT

---

> **Inkventory – quản lý kho sách, sẵn sàng mở rộng.**
