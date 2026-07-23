-- Create a new database
create database hoc_golang;

-- Drop a database
drop database testing;

-- Create a new schema
create schema school;

-- Drop a schema
drop schema school cascade;

-- One - One
-- Create users table
create table if not  exists users (
	user_id serial primary key,
	name varchar(50) not null,
	email varchar(100) unique not null
);

-- Create profiles table
create table if not exists profiles (
	profile_id serial primary key,
	user_id int unique not null,
	phone varchar(10),
	address varchar(100),
	constraint fk_user foreign key (user_id) references users(user_id) on delete cascade
);

-- Drop table (dangerours)
drop table if exists profiles;
drop table if exists users;

-- One - Many
-- Create categories table
create table if not exists categories (
	category_id serial primary key,
	name varchar(50) not null
);

-- Create products table
create table if not exists products (
	product_id serial primary key,
	category_id int not null,
	name varchar(100) not null,
	price int not null check (price > 0),
	image varchar(255),
	status int not null check (status in (1,2)),
	constraint fk_category foreign key (category_id) references categories (category_id) on delete restrict
);

-- Drop table (dangerours)
drop table if exists products;
drop table if exists categories;


-- Many - Many
-- Create students table
create table if not exists students (
	student_id serial primary key,
	name varchar(50) not null
);

-- Create courses table
create table if not exists courses (
	course_id serial primary key,
	name varchar(50) not null
);

-- Create students_courses table
create table if not exists students_courses (
	student_id int not null,
	course_id int not null,
	primary key (student_id, course_id),
	constraint fk_student foreign key (student_id) references students(student_id) on delete cascade,
	constraint fk_course foreign key (course_id) references courses(course_id) on delete cascade
)

-- Drop table (dangerours)
drop table if exists students_courses;
drop table if exists courses;
drop table if exists students;


------------------------- Các cấu truy vấn (SQL) hay sử dụng -----------------------------------
-- Thêm dữ liệu: INSERT INTO table (col1, col2) VALUES (val1, val2)
insert into users (name, email) values ('Vu Quoc Tuan', 'contact.quoctuan@gmail.com');
insert into users (name, email) values ('Toney Teo', 'contact.teo@gmail.com');
insert into users (name, email) values ('Le Van Tung', 'contact.tung@gmail.com');

insert into profiles (user_id, phone, address) values (1, '0901234567', '123 Cach mang thang 8');

insert into categories (name) values ('Dien thoai'), ('Laptop');

insert into products (category_id, name, price, image, status) values
(3, 'iPhone 18 Pro Max', 10, 'images/iphone-18-pro-max.jpg', 1),
(4, 'iPhone 17 Pro Max', 30000000, 'images/iphone-17-pro-max.jpg', 1);

-- Cập nhật dữ liệu: UPDATE table SET col1 = value1, col2 = val2 WHERE condition
update users set email = 'tuan@quoctuan.com', name = 'Mr.Tuan' where user_id = 1;
update profiles set phone = '0906784312' where user_id = 2;

-- Xóa dữ liệu: DELETE FROM table WHERE condition
delete from users where user_id = 3;
delete from products;
delete from categories where category_id = 1;

-- Lấy dữ liệu: SELECT * FROM table WHERE condition ORDER BY col [DESC/ASC] LIMIT ... OFFSET ...
select * from products;
select name, price from products;
select count(*) as total_rows from products;

select * from products where price >= 400000 and price <= 1000000;

select * from products order by price desc;
select * from products order by price asc;

select * from products limit 3 offset 4;

select name, price from products 
where price >= 400000 and price <= 30000000
order by price desc
limit 3;

select category_id, count(*) from products
group by category_id
having count(*) > 2;


------------------------- JOIN - Kết hợp dữ liệu từ nhiều bảng -----------------------------------
-- INNER JOIN: chỉ lấy các dòng khớp ở cả 2 bảng (1-1: users - profiles)
select u.name, u.email, p.phone, p.address
from users u
inner join profiles p on p.user_id = u.user_id;

-- LEFT JOIN: lấy toàn bộ dòng bên trái (users), dòng nào không có profile thì cột profile = NULL
select u.name, p.phone
from users u
left join profiles p on p.user_id = u.user_id;

-- JOIN 1-nhiều: mỗi product kèm tên category tương ứng
select pr.name as product_name, c.name as category_name, pr.price
from products pr
inner join categories c on c.category_id = pr.category_id;

-- JOIN nhiều-nhiều qua bảng trung gian: student nào học course nào
select s.name as student_name, c.name as course_name
from students_courses sc
inner join students s on s.student_id = sc.student_id
inner join courses c on c.course_id = sc.course_id;

-- RIGHT JOIN: lấy toàn bộ dòng bên phải (products), dòng category nào không tồn tại thì cột category = NULL
select pr.name as product_name, c.name as category_name
from categories c
right join products pr on pr.category_id = c.category_id;

-- FULL OUTER JOIN: lấy toàn bộ dòng ở cả 2 bảng, khớp được thì ghép, không thì để NULL
select c.name as category_name, pr.name as product_name
from categories c
full outer join products pr on pr.category_id = c.category_id;


------------------------- Hàm tổng hợp (Aggregate functions) -----------------------------------
select count(*) as total_products from products;
select avg(price) as avg_price from products;
select min(price) as min_price, max(price) as max_price from products;

-- Kết hợp với GROUP BY: giá trung bình theo từng category
select category_id, avg(price) as avg_price from products
group by category_id;


------------------------- DISTINCT - Loại bỏ giá trị trùng lặp -----------------------------------
-- Danh sách category_id đang có sản phẩm (không trùng lặp)
select distinct category_id from products;

-- DISTINCT trên nhiều cột: các cặp (category_id, status) duy nhất
select distinct category_id, status from products;


------------------------- Subquery - Truy vấn lồng nhau -----------------------------------
-- Các sản phẩm thuộc category có tên 'Laptop'
select * from products
where category_id = (select category_id from categories where name = 'Laptop');

-- Users chưa có profile (subquery trong NOT IN)
select * from users
where user_id not in (select user_id from profiles);

-- Sản phẩm có giá cao hơn giá trung bình
select * from products
where price > (select avg(price) from products);


------------------------- Transaction - Giao dịch -----------------------------------
-- Gom nhiều câu lệnh thành 1 khối "tất cả hoặc không gì cả" (atomic)
begin;

update products set price = price - 500000 where product_id = 1;
update products set price = price + 500000 where product_id = 2;

commit; -- lưu thay đổi thật sự vào DB
-- rollback; -- nếu có lỗi thì gọi rollback để huỷ toàn bộ thay đổi trong transaction

-- SAVEPOINT: rollback về 1 điểm giữa chừng thay vì huỷ toàn bộ transaction
begin;
update products set price = 100 where product_id = 1;
savepoint sp1;
update products set price = 200 where product_id = 2; -- giả sử bước này sai
rollback to savepoint sp1; -- chỉ huỷ bước sau savepoint, giữ lại bước trước
commit;

-- Isolation level (mức cô lập) - hạn chế các transaction chạy song song "thấy" nhau thế nào
-- READ COMMITTED (mặc định ở Postgres) | REPEATABLE READ | SERIALIZABLE
begin transaction isolation level repeatable read;
select * from products where product_id = 1;
commit;


------------------------- UPSERT - INSERT hoặc UPDATE nếu đã tồn tại -----------------------------------
-- Cần có unique/primary key để Postgres biết dòng nào là "trùng" (ở đây là email)
insert into users (name, email)
values ('Vu Quoc Tuan', 'contact.quoctuan@gmail.com')
on conflict (email) do update
set name = excluded.name;

-- Nếu trùng thì bỏ qua, không làm gì cả
insert into categories (name) values ('Laptop')
on conflict (name) do nothing;


------------------------- Window functions - Hàm cửa sổ -----------------------------------
-- Khác GROUP BY: window function KHÔNG gom dòng lại, mỗi dòng vẫn giữ nguyên, chỉ thêm cột tính toán

-- Đánh số thứ tự sản phẩm theo giá giảm dần, trong từng category
select
	product_id,
	category_id,
	name,
	price,
	row_number() over (partition by category_id order by price desc) as row_num,
	rank() over (partition by category_id order by price desc) as price_rank
from products;

-- Tổng giá trị toàn bộ category ngay trên từng dòng sản phẩm (running context, không gom dòng)
select
	product_id,
	name,
	category_id,
	price,
	sum(price) over (partition by category_id) as total_category_price
from products;


------------------------- CTE (WITH) - Truy vấn lồng nhau dễ đọc hơn subquery -----------------------------------
-- WITH đặt tên cho 1 truy vấn con, dùng lại được nhiều lần trong câu truy vấn chính
with expensive_products as (
	select * from products where price > 1000000
)
select category_id, count(*) from expensive_products
group by category_id;

-- Recursive CTE: dùng cho dữ liệu phân cấp (cây thư mục, sơ đồ tổ chức...)
-- Ví dụ: bảng employees tự tham chiếu chính nó qua manager_id
create table if not exists employees (
	employee_id serial primary key,
	name varchar(50) not null,
	manager_id int references employees(employee_id)
);

insert into employees (name, manager_id) values
('CEO Nam', null),
('Manager Hoa', 1),
('Staff Long', 2),
('Staff Mai', 2);

-- Lấy toàn bộ cấp dưới (trực tiếp + gián tiếp) của "CEO Nam"
with recursive subordinates as (
	select employee_id, name, manager_id from employees where name = 'CEO Nam'
	union all
	select e.employee_id, e.name, e.manager_id
	from employees e
	inner join subordinates s on e.manager_id = s.employee_id
)
select * from subordinates;

drop table if exists employees;


------------------------- Index & EXPLAIN ANALYZE - Tối ưu truy vấn -----------------------------------
-- Tạo index giúp tăng tốc độ tìm kiếm/lọc/sắp xếp trên cột đó (đánh đổi: chậm hơn khi insert/update)
create index if not exists idx_products_category_id on products (category_id);
create unique index if not exists idx_users_email on users (email); -- unique constraint cũng tự tạo index

-- EXPLAIN ANALYZE: xem Postgres THẬT SỰ chạy câu truy vấn như thế nào (có dùng index không, mất bao lâu)
explain analyze
select * from products where category_id = 1;

drop index if exists idx_products_category_id;


------------------------- Locking - Khoá dòng khi cập nhật đồng thời -----------------------------------
-- FOR UPDATE: khoá (các) dòng được SELECT lại, transaction khác phải chờ tới khi transaction này COMMIT/ROLLBACK
-- Thường dùng khi đọc số dư/tồn kho rồi update, để tránh 2 request cùng đọc 1 giá trị cũ rồi ghi đè nhau (race condition)
begin;

select * from products where product_id = 1 for update;
update products set price = price - 100000 where product_id = 1;

commit;


------------------------- JSON / JSONB - Dữ liệu dạng JSON (đặc trưng của Postgres) -----------------------------------
alter table products add column if not exists metadata jsonb;

update products set metadata = '{"color": "black", "warranty_months": 12, "tags": ["hot", "new"]}'
where product_id = 1;

-- -> trả về JSON, ->> trả về text
select metadata -> 'color' as color_json, metadata ->> 'color' as color_text from products;

-- Lọc theo giá trị trong JSON
select * from products where metadata ->> 'color' = 'black';

-- @> kiểm tra JSON có "chứa" 1 phần tử/giá trị hay không (containment)
select * from products where metadata @> '{"warranty_months": 12}';

alter table products drop column if exists metadata;


------------------------- EXISTS vs IN -----------------------------------
-- IN: so sánh với 1 danh sách giá trị cụ thể (subquery trả về nhiều dòng, 1 cột)
select * from products
where category_id in (select category_id from categories where name in ('Laptop', 'Dien thoai'));

-- EXISTS: chỉ kiểm tra CÓ tồn tại dòng nào khớp hay không, thường nhanh hơn IN với bảng lớn
select * from categories c
where exists (select 1 from products p where p.category_id = c.category_id);


------------------------- UNION / INTERSECT / EXCEPT - Kết hợp kết quả của nhiều truy vấn -----------------------------------
-- UNION: gộp kết quả, tự loại bỏ trùng lặp (UNION ALL thì giữ nguyên, không loại trùng -> nhanh hơn)
select name from users
union
select name from students;

select name from users
union all
select name from students;

-- INTERSECT: chỉ lấy dòng xuất hiện ở CẢ 2 truy vấn
-- select name from table_a intersect select name from table_b;

-- EXCEPT: lấy dòng có ở truy vấn 1 nhưng KHÔNG có ở truy vấn 2
-- select name from table_a except select name from table_b;


------------------------- CASE WHEN - Rẽ nhánh điều kiện trong SELECT -----------------------------------
select
	name,
	price,
	case
		when price < 1000000 then 'Rẻ'
		when price between 1000000 and 10000000 then 'Trung bình'
		else 'Đắt'
	end as price_range
from products;


------------------------- Hàm xử lý ngày giờ & chuỗi thường dùng -----------------------------------
select now();                                  -- thời điểm hiện tại (timestamp)
select current_date;                           -- ngày hiện tại
select now() - interval '7 days';              -- 7 ngày trước
select age(now(), '2020-01-01'::date);         -- khoảng cách thời gian giữa 2 mốc

select upper(name), lower(email) from users;   -- viết hoa / viết thường
select concat(name, ' <', email, '>') as display_name from users; -- nối chuỗi
select length(name) from users;                -- độ dài chuỗi
select substring(email from 1 for 3) from users; -- cắt chuỗi (3 ký tự đầu)
select trim('  hello  ');                      -- xoá khoảng trắng 2 đầu