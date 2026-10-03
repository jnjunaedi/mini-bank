# 🏦 Mini-Bank Core System (Microservices)

Proyek ini adalah implementasi sistem *Core Banking* berbasis arsitektur **Microservices** menggunakan bahasa pemrograman **Go (Golang)**. Sistem ini dirancang untuk menangani transaksi perbankan dasar secara aman, konsisten, dan memiliki cakupan pengujian (*test coverage*) [index].

---

## 🛠️ Tech Stack & Pustaka Utama

Sistem ini dirakit menggunakan ekosistem teknologi modern:

*   **Bahasa Pemrograman Utama**: ![Go](https://shields.io) **Golang v1.23** (Bersih, efisien, dan mengutamakan performa murni).
*   **Protokol Komunikasi**: ![gRPC](https://shields.io) **Protocol Buffers (Proto3)** untuk kontrak data gRPC internal yang super cepat [index].
*   **Penyimpanan Data (Database)**: ![MySQL](https://shields.io) **MySQL 8.0** sebagai penyimpanan transaksi yang bersifat ACID-compliant.
*   **Kontainerisasi**: ![Docker](https://shields.io) **Docker & Docker Compose** untuk standarisasi lingkungan aplikasi dari lokal hingga produksi [index].
*   **Pustaka Pengujian (Testing)**: **`Stretchr/Testify`** (`assert` & `mock`) untuk arsitektur pengujian [index].
*   **Driver Database**: **`go-sql-driver/mysql`** untuk interaksi dengan *database engine*.

---


## 🏗️ Arsitektur Sistem & Teknis

Aplikasi ini dipecah menjadi dua layanan utama (*microservices*) yang saling berkomunikasi secara efisien menggunakan **gRPC**:

1. **`account-service`**: Bertanggung jawab atas manajemen profil nasabah, verifikasi PIN keamanan, dan pengecekan status rekening [index].
2. **`transaction-service`**: Bertanggung jawab penuh atas logika bisnis keuangan seperti Setor Tunai (`CashDeposit`), Tarik Tunai (`CashWithdrawal`), Transfer (`Transfer`) Antar Rekening, hingga Pemrosesan Gaji Massal (`ProcessBatchGaji`) [index, index].

### Fitur Unggulan Arsitektur:
* **gRPC Communication**: Komunikasi antar-service super cepat menggunakan protokol gRPC internal, memisahkan gerbang HTTP luar dari koneksi internal [index].
* **Strict Unit Testing**: Pengujian kode menggunakan metode *Mocking* berlapis (`testify/mock`) untuk mensimulasikan kegagalan jaringan dan isolasi *database* [index].
* **Docker Multi-Stage Builds**: Konfigurasi kontainer cerdas yang hanya membawa file biner matang, memangkas ukuran *image* hingga hanya ~15MB demi efisiensi RAM server [index].

---

## 🧪 Strategi Unit Testing (Robust & Clean)

* **Isolasi Database & Jaringan**: Menggunakan robot tiruan ganda untuk memotong dependensi fisik ke MySQL dan gRPC luar [index].
* **Pengujian Array yang Efisien**: Pada fitur *Batch Gaji Massal*, validasi dilakukan secara cerdas melalui ukuran struktural array dan pengecekan sampel elemen guna memastikan integritas performa memori [index].

### Cara Menjalankan Unit Test (Lokal)
Masuk ke salah satu folder *service*, lalu jalankan perintah bawaan Go:
```bash
go test -v ./...
```

---

## 🐳 Konfigurasi Kontainerisasi (Docker Compose)

* **Urutan Startup (`depends_on`)**: Memastikan *database* menyala sebelum aplikasi Go mencoba menghubungkan jaringan [index].
* **Dynamic Network Discovery**: Komunikasi gRPC antar-kontainer otomatis terhubung murni menggunakan nama panggung *service* (`account-service:50051`) tanpa ketergantungan IP statis [index].

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
- [ ] Integrasi distributed tracing menggunakan OpenTelemetry.
- [ ] Penerapan pola sinkronisasi data asinkron menggunakan Apache Kafka / RabbitMQ.
- [ ] Implementasi Database Migration otomatis saat kontainer Docker pertama kali naik.
