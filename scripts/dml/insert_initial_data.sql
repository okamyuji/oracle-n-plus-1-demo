-- Oracle N+1問題デモ用DMLスクリプト
-- 初期データ投入

-- ============================================
-- 部署マスターデータ投入
-- ============================================

INSERT INTO departments (department_id, department_name, location) VALUES
(seq_departments.NEXTVAL, '営業部', '東京');

INSERT INTO departments (department_id, department_name, location) VALUES
(seq_departments.NEXTVAL, '開発部', '大阪');

INSERT INTO departments (department_id, department_name, location) VALUES
(seq_departments.NEXTVAL, '人事部', '東京');

INSERT INTO departments (department_id, department_name, location) VALUES
(seq_departments.NEXTVAL, '総務部', '名古屋');

INSERT INTO departments (department_id, department_name, location) VALUES
(seq_departments.NEXTVAL, 'マーケティング部', '福岡');

INSERT INTO departments (department_id, department_name, location) VALUES
(seq_departments.NEXTVAL, '経理部', '東京');

INSERT INTO departments (department_id, department_name, location) VALUES
(seq_departments.NEXTVAL, '法務部', '東京');

INSERT INTO departments (department_id, department_name, location) VALUES
(seq_departments.NEXTVAL, 'IT企画部', '大阪');

-- ============================================
-- 社員データ投入
-- ============================================

-- 営業部（department_id = 1）
INSERT INTO employees (employee_id, first_name, last_name, email, department_id, salary, hire_date) VALUES
(seq_employees.NEXTVAL, '太郎', '田中', 'tanaka.taro@company.com', 1, 5000000, TO_DATE('2020-04-01', 'YYYY-MM-DD'));

INSERT INTO employees (employee_id, first_name, last_name, email, department_id, salary, hire_date) VALUES
(seq_employees.NEXTVAL, '花子', '佐藤', 'sato.hanako@company.com', 1, 4800000, TO_DATE('2021-04-01', 'YYYY-MM-DD'));

INSERT INTO employees (employee_id, first_name, last_name, email, department_id, salary, hire_date) VALUES
(seq_employees.NEXTVAL, '一郎', '鈴木', 'suzuki.ichiro@company.com', 1, 5200000, TO_DATE('2019-04-01', 'YYYY-MM-DD'));

-- 開発部（department_id = 2）
INSERT INTO employees (employee_id, first_name, last_name, email, department_id, salary, hire_date) VALUES
(seq_employees.NEXTVAL, '次郎', '高橋', 'takahashi.jiro@company.com', 2, 6000000, TO_DATE('2018-04-01', 'YYYY-MM-DD'));

INSERT INTO employees (employee_id, first_name, last_name, email, department_id, salary, hire_date) VALUES
(seq_employees.NEXTVAL, '美咲', '山田', 'yamada.misaki@company.com', 2, 5800000, TO_DATE('2020-10-01', 'YYYY-MM-DD'));

INSERT INTO employees (employee_id, first_name, last_name, email, department_id, salary, hire_date) VALUES
(seq_employees.NEXTVAL, '健太', '中村', 'nakamura.kenta@company.com', 2, 5500000, TO_DATE('2022-04-01', 'YYYY-MM-DD'));

-- 人事部（department_id = 3）
INSERT INTO employees (employee_id, first_name, last_name, email, department_id, salary, hire_date) VALUES
(seq_employees.NEXTVAL, '雅子', '伊藤', 'ito.masako@company.com', 3, 4500000, TO_DATE('2017-04-01', 'YYYY-MM-DD'));

INSERT INTO employees (employee_id, first_name, last_name, email, department_id, salary, hire_date) VALUES
(seq_employees.NEXTVAL, '正樹', '渡辺', 'watanabe.masaki@company.com', 3, 4700000, TO_DATE('2019-10-01', 'YYYY-MM-DD'));

-- 総務部（department_id = 4）
INSERT INTO employees (employee_id, first_name, last_name, email, department_id, salary, hire_date) VALUES
(seq_employees.NEXTVAL, '真理', '小林', 'kobayashi.mari@company.com', 4, 4200000, TO_DATE('2021-04-01', 'YYYY-MM-DD'));

-- マーケティング部（department_id = 5）
INSERT INTO employees (employee_id, first_name, last_name, email, department_id, salary, hire_date) VALUES
(seq_employees.NEXTVAL, '大輝', '加藤', 'kato.daiki@company.com', 5, 5300000, TO_DATE('2020-04-01', 'YYYY-MM-DD'));

-- ============================================
-- 受注データ投入
-- ============================================

-- 受注1
INSERT INTO orders (order_id, customer_id, customer_name, order_date, total_amount, status) VALUES
(seq_orders.NEXTVAL, 1001, '株式会社ABC商事', TO_DATE('2024-01-15', 'YYYY-MM-DD'), 250000, 'COMPLETED');

-- 受注2
INSERT INTO orders (order_id, customer_id, customer_name, order_date, total_amount, status) VALUES
(seq_orders.NEXTVAL, 1002, '有限会社XYZ販売', TO_DATE('2024-01-20', 'YYYY-MM-DD'), 180000, 'COMPLETED');

-- 受注3
INSERT INTO orders (order_id, customer_id, customer_name, order_date, total_amount, status) VALUES
(seq_orders.NEXTVAL, 1003, '株式会社DEF企画', TO_DATE('2024-02-01', 'YYYY-MM-DD'), 320000, 'PENDING');

-- 受注4
INSERT INTO orders (order_id, customer_id, customer_name, order_date, total_amount, status) VALUES
(seq_orders.NEXTVAL, 1001, '株式会社ABC商事', TO_DATE('2024-02-10', 'YYYY-MM-DD'), 150000, 'PROCESSING');

-- 受注5
INSERT INTO orders (order_id, customer_id, customer_name, order_date, total_amount, status) VALUES
(seq_orders.NEXTVAL, 1004, '合同会社GHI物産', TO_DATE('2024-02-15', 'YYYY-MM-DD'), 420000, 'COMPLETED');

-- ============================================
-- 受注明細データ投入
-- ============================================

-- 受注1の明細
INSERT INTO order_details (detail_id, order_id, product_id, product_name, quantity, unit_price) VALUES
(seq_order_details.NEXTVAL, 1, 2001, 'ノートパソコン Type-A', 5, 80000);

INSERT INTO order_details (detail_id, order_id, product_id, product_name, quantity, unit_price) VALUES
(seq_order_details.NEXTVAL, 1, 2002, 'ワイヤレスマウス', 10, 3000);

INSERT INTO order_details (detail_id, order_id, product_id, product_name, quantity, unit_price) VALUES
(seq_order_details.NEXTVAL, 1, 2003, 'USB-Cハブ', 5, 4000);

-- 受注2の明細
INSERT INTO order_details (detail_id, order_id, product_id, product_name, quantity, unit_price) VALUES
(seq_order_details.NEXTVAL, 2, 2004, 'モニター 24インチ', 3, 45000);

INSERT INTO order_details (detail_id, order_id, product_id, product_name, quantity, unit_price) VALUES
(seq_order_details.NEXTVAL, 2, 2005, 'キーボード', 3, 15000);

-- 受注3の明細
INSERT INTO order_details (detail_id, order_id, product_id, product_name, quantity, unit_price) VALUES
(seq_order_details.NEXTVAL, 3, 2001, 'ノートパソコン Type-A', 2, 80000);

INSERT INTO order_details (detail_id, order_id, product_id, product_name, quantity, unit_price) VALUES
(seq_order_details.NEXTVAL, 3, 2006, 'プリンター', 1, 120000);

INSERT INTO order_details (detail_id, order_id, product_id, product_name, quantity, unit_price) VALUES
(seq_order_details.NEXTVAL, 3, 2007, 'スキャナー', 2, 30000);

-- 受注4の明細
INSERT INTO order_details (detail_id, order_id, product_id, product_name, quantity, unit_price) VALUES
(seq_order_details.NEXTVAL, 4, 2002, 'ワイヤレスマウス', 20, 3000);

INSERT INTO order_details (detail_id, order_id, product_id, product_name, quantity, unit_price) VALUES
(seq_order_details.NEXTVAL, 4, 2008, 'Webカメラ', 15, 6000);

-- 受注5の明細
INSERT INTO order_details (detail_id, order_id, product_id, product_name, quantity, unit_price) VALUES
(seq_order_details.NEXTVAL, 5, 2009, 'サーバー', 1, 300000);

INSERT INTO order_details (detail_id, order_id, product_id, product_name, quantity, unit_price) VALUES
(seq_order_details.NEXTVAL, 5, 2010, 'ネットワーク機器', 2, 60000);

-- ============================================
-- 統計情報の更新
-- ============================================
EXEC DBMS_STATS.GATHER_TABLE_STATS(USER, 'DEPARTMENTS');
EXEC DBMS_STATS.GATHER_TABLE_STATS(USER, 'EMPLOYEES');
EXEC DBMS_STATS.GATHER_TABLE_STATS(USER, 'ORDERS');
EXEC DBMS_STATS.GATHER_TABLE_STATS(USER, 'ORDER_DETAILS');

COMMIT;

-- ============================================
-- データ確認用クエリ
-- ============================================

-- 部署別社員数
SELECT 
    d.department_name,
    COUNT(e.employee_id) as employee_count
FROM departments d
LEFT JOIN employees e ON d.department_id = e.department_id
GROUP BY d.department_id, d.department_name
ORDER BY d.department_id;

-- 受注別明細数
SELECT 
    o.order_id,
    o.customer_name,
    o.order_date,
    COUNT(od.detail_id) as detail_count,
    SUM(od.line_amount) as calculated_total
FROM orders o
LEFT JOIN order_details od ON o.order_id = od.order_id
GROUP BY o.order_id, o.customer_name, o.order_date
ORDER BY o.order_id;
