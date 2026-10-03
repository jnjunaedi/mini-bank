# 🏦 Mini-Bank Core System (Microservices)

This project is an implementation of a **Core Banking** system based on a **Microservices** architecture using the **Go (Golang)** programming language. This system is designed to handle basic banking transactions safely, consistently, and has test coverage [index].

---

## 🛠️ Tech Stack & Pustaka Utama

This system is built using a modern technology ecosystem:

*   **Primary Programming Language**: ![Go](https://shields.io) **Golang v1.23** (Clean, efficient, and prioritizes pure performance).
*   **Communication Protocol**: ![gRPC](https://shields.io) **Protocol Buffers (Proto3)** for super fast internal gRPC data contracts [index].
*   **Data Storage (Database)**: ![MySQL](https://shields.io) **MySQL 8.0** as an ACID-compliant transaction storage.
*   **Containerization**: ![Docker](https://shields.io) **Docker & Docker Compose** for application environment standardization from local to production [index].
*   **Testing Library**: **`Stretchr/Testify`** (`assert` & `mock`) for the testing architecture [index].
*   **Database Driver**: **`go-sql-driver/mysql`** for interaction with the database engine.

---

## 🏗️ Arsitektur Sistem & Teknis

This application is split into two main services (*microservices*) that communicate efficiently using **gRPC**:

1. **`account-service`**: Responsible for customer profile management, secure PIN verification, and account status checking [index].
2. **`transaction-service`**: Fully responsible for financial business logic such as Cash Deposit (`CashDeposit`), Cash Withdrawal (`CashWithdrawal`), Inter-Account Transfer (`Transfer`), up to Bulk Payroll Processing (`ProcessBatchGaji`) [index, index].

### Fitur Unggulan Arsitektur:
* **gRPC Communication**: Super fast communication between services using internal gRPC protocol, separating the external HTTP gateway from the internal connection [index].
* **Strict Unit Testing**: Code testing using layered *Mocking* method (`testify/mock`) to simulate network failure and *database* isolation [index].
* **Docker Multi-Stage Builds**: Smart container configuration that only carries mature binary files, cutting *image* size down to only ~15MB for server RAM efficiency [index].

---

## 🧪 Strategi Unit Testing (Robust & Clean)

* **Database & Network Isolation**: Using dual mock robots to cut physical dependencies to MySQL and external gRPC [index].
* **Efficient Array Testing**: In the *Bulk Payroll* feature, validation is carried out smartly through the structural size of the array and checking sample elements to ensure memory performance integrity [index].

### Cara Menjalankan Unit Test (Lokal)
Enter one of the *service* folders, then run the Go native command:
```bash
go test -v ./...
```

---

## 🐳 Konfigurasi Kontainerisasi (Docker Compose)

* **Startup Order (`depends_on`)**: Ensures the *database* is running before the Go application attempts to connect to the network [index].
* **Dynamic Network Discovery**: Inter-container gRPC communication connects automatically purely using the *service* stage name (`account-service:50051`) without static IP dependency [index].

### Berkas Orkestra (`docker-compose.yml`)
```yaml
services:
  minibank-db:
    image: mysql:8.0
    container_name: minibank-mysql-db
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: minibank_db
    ports:
      - "3307:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    restart: always

  account-service:
    build:
      context: ./account-service
      dockerfile: Dockerfile
    container_name: account-microservice
    ports:
      - "8080:8080"
      - "50051:50051"
```

---

## 📈 Pengembangan Masa Depan (Roadmap)
- [ ] Distributed tracing integration using OpenTelemetry.
- [ ] Implementation of asynchronous data synchronization patterns using Apache Kafka / RabbitMQ.
- [ ] Automatic Database Migration implementation when the Docker container first goes up.
