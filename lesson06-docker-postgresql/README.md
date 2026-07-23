# Lesson 06 — Docker & PostgreSQL

Chạy PostgreSQL bằng Docker Compose (xem [docker-compose.yml](docker-compose.yml)).

**Nội dung khác:** [Quan hệ (relationship) trong SQL](SQL-RELATIONSHIPS.md)

## Các lệnh Docker thường dùng

| CMD | Ý nghĩa |
|---|---|
| `docker-compose up -d` | Lệnh để khởi động và tạo container (nếu chưa có) từ file cấu hình.<br>`-d` viết tắt của `--detach` → nghĩa là chạy ngầm (ẩn) |
| `docker-compose down` | Dừng và xóa toàn bộ container nhưng sẽ không xóa volume và image |
| `docker ps` | Xem những container đang chạy<br>(optional: `-a` để xem toàn bộ container) |
| `docker image` | Liệt kê toàn bộ image trên docker |
| `docker rmi <image id \| name>` | Xóa image cụ thể |
| `docker stop <container id \| name>` | Dừng 1 container cụ thể |
| `docker restart <container id \| name>` | Khởi động lại 1 container cụ thể |
| `docker rm <container id \| name>` | Xóa container (cần phải dừng trước khi xóa) |
| `docker volume ls` | Liệt kê toàn bộ volume trên host |
| `docker volume rm <volume name>` | Xóa volume (⚠️ nguy hiểm: sẽ mất toàn bộ dữ liệu nếu đang lưu trong volume) |

## Kết nối tới database

| Thông tin | Giá trị |
|---|---|
| Host | `localhost` |
| Port | `5432` |
| Database | `master-golang` |
| User | `root` |
| Password | xem `POSTGRES_PASSWORD` trong `docker-compose.yml` |

```bash
# Vào psql bên trong container
docker exec -it postgres-db psql -U root -d master-golang
```

## Import và Export database

**Export database** (dump toàn bộ database ra file `.sql` trên máy host):

```bash
docker exec -i postgres-db pg_dump -U root -d master-golang > ./backupdb-master-golang.sql
```

**Import database** (nạp lại file `.sql` vào database trong container):

```bash
docker exec -i postgres-db psql -U root -d master-golang < ./backupdb-master-golang.sql
```

> `-i` giữ stdin mở để `>` / `<` chuyển dữ liệu giữa host và container.
> Ở đây không dùng `-t` vì dữ liệu không đi qua terminal.
