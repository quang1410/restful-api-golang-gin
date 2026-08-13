-- Seed data cho bảng users: 100 user active + 10 user đã soft delete.
-- Tất cả email đều thuộc domain @example.com nên có thể xoá sạch dữ liệu seed
-- mà không đụng tới dữ liệu thật (xem make seed-clean).
--
-- Password gốc của mọi user: Password@123 (bcrypt cost 10).

BEGIN;

-- Xoá dữ liệu seed cũ để chạy lại nhiều lần vẫn cho kết quả giống nhau
DELETE FROM users WHERE user_email LIKE '%@example.com';

-- 100 user active -----------------------------------------------------------
-- user_created_at được đảo thứ tự so với user_id (bước nhảy 254 ngày mod 365)
-- để order_by=user_id và order_by=user_created_at cho ra 2 thứ tự khác nhau.
INSERT INTO users (
    user_email,
    user_password,
    user_fullname,
    user_age,
    user_status,
    user_level,
    user_created_at,
    user_updated_at
)
SELECT
    lower(n.last_name) || '.' || lower(n.first_name) || n.i || '@example.com',
    '$2a$10$eFREkT9072NiTo5P3lngB.S1ZFwvhAslH1AHk25eEUlq9uMqeHbPi',
    n.last_name || ' ' || n.mid_name || ' ' || n.first_name,
    18 + (n.i * 7) % 45,
    CASE WHEN n.i % 10 = 0 THEN 3 WHEN n.i % 5 = 0 THEN 2 ELSE 1 END,
    CASE WHEN n.i % 25 = 0 THEN 1 WHEN n.i % 5 = 0 THEN 2 ELSE 3 END,
    now() - (((n.i * 254) % 365) || ' days')::interval - ((n.i * 13) || ' minutes')::interval,
    now() - (((n.i * 254) % 365) || ' days')::interval - ((n.i * 13) || ' minutes')::interval
FROM (
    SELECT
        s.i,
        (ARRAY['Nguyen','Tran','Le','Pham','Hoang','Vu','Dang','Bui','Do','Ngo'])[(s.i % 10) + 1]              AS last_name,
        (ARRAY['Van','Thi','Minh','Quoc','Hai','Thanh','Ngoc','Anh','Duy','Phuong'])[((s.i * 3) % 10) + 1]     AS mid_name,
        (ARRAY['An','Binh','Cuong','Dung','Hoa','Khanh','Linh','Nam','Quan','Tuan',
               'Yen','Bao','Chi','Giang','Hieu'])[((s.i * 7) % 15) + 1]                                        AS first_name
    FROM generate_series(1, 100) AS s(i)
) AS n;

-- 10 user đã soft delete (để test restore / trash / filter deleted) ----------
INSERT INTO users (
    user_email,
    user_password,
    user_fullname,
    user_age,
    user_status,
    user_level,
    user_created_at,
    user_updated_at,
    user_deleted_at
)
SELECT
    'deleted.user' || s.i || '@example.com',
    '$2a$10$eFREkT9072NiTo5P3lngB.S1ZFwvhAslH1AHk25eEUlq9uMqeHbPi',
    'Deleted User ' || s.i,
    20 + s.i,
    2,
    3,
    now() - ((s.i * 11) || ' days')::interval,
    now() - ((s.i * 11) || ' days')::interval,
    now() - ((s.i * 2) || ' days')::interval
FROM generate_series(1, 10) AS s(i);

COMMIT;
