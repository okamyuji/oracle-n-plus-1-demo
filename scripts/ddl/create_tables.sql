-- Oracle N+1問題デモ用DDLスクリプト
-- 受注管理システムと社員管理システムのテーブル作成

-- ============================================
-- 部署マスターテーブル
-- ============================================
CREATE TABLE departments (
    department_id NUMBER(10) PRIMARY KEY,
    department_name VARCHAR2(100) NOT NULL,
    location VARCHAR2(100),
    created_at DATE DEFAULT SYSDATE,
    updated_at DATE DEFAULT SYSDATE
);

-- 部署名にインデックス作成（検索で使用）
CREATE INDEX idx_departments_name ON departments(department_name);

-- ============================================
-- 社員テーブル
-- ============================================
CREATE TABLE employees (
    employee_id NUMBER(10) PRIMARY KEY,
    first_name VARCHAR2(50) NOT NULL,
    last_name VARCHAR2(50) NOT NULL,
    email VARCHAR2(200) UNIQUE NOT NULL,
    department_id NUMBER(10),
    salary NUMBER(10,2),
    hire_date DATE DEFAULT SYSDATE,
    created_at DATE DEFAULT SYSDATE,
    updated_at DATE DEFAULT SYSDATE,
    CONSTRAINT fk_employees_department 
        FOREIGN KEY (department_id) REFERENCES departments(department_id)
);

-- 部署IDにインデックス作成（JOINで使用）
CREATE INDEX idx_employees_department_id ON employees(department_id);
-- 社員名にインデックス作成（検索で使用）
CREATE INDEX idx_employees_name ON employees(last_name, first_name);

-- ============================================
-- 受注テーブル
-- ============================================
CREATE TABLE orders (
    order_id NUMBER(10) PRIMARY KEY,
    customer_id NUMBER(10) NOT NULL,
    customer_name VARCHAR2(100) NOT NULL,
    order_date DATE DEFAULT SYSDATE,
    total_amount NUMBER(12,2) DEFAULT 0,
    status VARCHAR2(20) DEFAULT 'PENDING',
    created_at DATE DEFAULT SYSDATE,
    updated_at DATE DEFAULT SYSDATE
);

-- 顧客IDにインデックス作成
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
-- 受注日にインデックス作成（期間検索で使用）
CREATE INDEX idx_orders_order_date ON orders(order_date);
-- ステータスにインデックス作成
CREATE INDEX idx_orders_status ON orders(status);

-- ============================================
-- 受注明細テーブル
-- ============================================
CREATE TABLE order_details (
    detail_id NUMBER(10) PRIMARY KEY,
    order_id NUMBER(10) NOT NULL,
    product_id NUMBER(10) NOT NULL,
    product_name VARCHAR2(200) NOT NULL,
    quantity NUMBER(8) NOT NULL,
    unit_price NUMBER(10,2) NOT NULL,
    line_amount NUMBER(12,2) GENERATED ALWAYS AS (quantity * unit_price),
    created_at DATE DEFAULT SYSDATE,
    updated_at DATE DEFAULT SYSDATE,
    CONSTRAINT fk_order_details_order 
        FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
);

-- 受注IDにインデックス作成（N+1問題対策で重要）
CREATE INDEX idx_order_details_order_id ON order_details(order_id);
-- 商品IDにインデックス作成
CREATE INDEX idx_order_details_product_id ON order_details(product_id);

-- ============================================
-- シーケンス作成
-- ============================================

-- 部署ID用シーケンス
CREATE SEQUENCE seq_departments
    START WITH 1
    INCREMENT BY 1
    NOCACHE;

-- 社員ID用シーケンス
CREATE SEQUENCE seq_employees
    START WITH 1
    INCREMENT BY 1
    NOCACHE;

-- 受注ID用シーケンス
CREATE SEQUENCE seq_orders
    START WITH 1
    INCREMENT BY 1
    NOCACHE;

-- 受注明細ID用シーケンス
CREATE SEQUENCE seq_order_details
    START WITH 1
    INCREMENT BY 1
    NOCACHE;

-- ============================================
-- 統計情報収集
-- ============================================

-- 統計情報を収集してオプティマイザの判断を向上させる
EXEC DBMS_STATS.GATHER_TABLE_STATS(USER, 'DEPARTMENTS');
EXEC DBMS_STATS.GATHER_TABLE_STATS(USER, 'EMPLOYEES');
EXEC DBMS_STATS.GATHER_TABLE_STATS(USER, 'ORDERS');
EXEC DBMS_STATS.GATHER_TABLE_STATS(USER, 'ORDER_DETAILS');

COMMIT;
