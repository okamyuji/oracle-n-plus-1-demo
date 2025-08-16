#!/usr/bin/env bash

# Oracle N+1問題デモ用大量ダミーデータ生成スクリプト
# 負荷テスト用の大量データを生成

set -e

# スクリプト設定
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# 設定変数
ORACLE_USER="${ORACLE_USER:-testuser}"
ORACLE_PASSWORD="${ORACLE_PASSWORD:-testpass}"
ORACLE_HOST="${ORACLE_HOST:-localhost}"
ORACLE_PORT="${ORACLE_PORT:-1521}"
ORACLE_SERVICE="${ORACLE_SERVICE:-XEPDB1}"

# データ生成量設定
DEPARTMENTS_COUNT=20        # 部署数
EMPLOYEES_PER_DEPT=50      # 部署あたりの社員数
ORDERS_COUNT=1000          # 受注数
DETAILS_PER_ORDER=5        # 受注あたりの明細数（平均）

echo "============================================"
echo "Oracle N+1問題デモ用大量データ生成開始"
echo "============================================"
echo "接続先: ${ORACLE_HOST}:${ORACLE_PORT}/${ORACLE_SERVICE}"
echo "ユーザー: ${ORACLE_USER}"
echo ""
echo "生成データ量:"
echo "  - 部署数: ${DEPARTMENTS_COUNT}"
echo "  - 社員数: $((DEPARTMENTS_COUNT * EMPLOYEES_PER_DEPT))"
echo "  - 受注数: ${ORDERS_COUNT}"
echo "  - 明細数: $((ORDERS_COUNT * DETAILS_PER_ORDER))"
echo "============================================"
echo ""

# SQLPlusコマンドチェック
if ! command -v sqlplus &> /dev/null; then
    echo "エラー: sqlplus コマンドが見つかりません。Oracle Clientをインストールしてください。"
    exit 1
fi

# 接続確認
echo "データベース接続確認中..."
if ! echo "SELECT 1 FROM DUAL;" | sqlplus -S "${ORACLE_USER}/${ORACLE_PASSWORD}@${ORACLE_HOST}:${ORACLE_PORT}/${ORACLE_SERVICE}" > /dev/null 2>&1; then
    echo "エラー: データベースに接続できません。接続情報を確認してください。"
    exit 1
fi
echo "接続OK"
echo ""

# 一時SQLファイル生成
TEMP_SQL="${SCRIPT_DIR}/load_test_data.sql"

echo "大量データ生成SQLを作成中..."

cat > "$TEMP_SQL" << 'EOF'
-- 大量ダミーデータ生成スクリプト

SET SERVEROUTPUT ON;
SET TIMING ON;

DECLARE
    -- 日本語の姓・名リスト
    TYPE name_array IS VARRAY(50) OF VARCHAR2(20);
    surnames name_array := name_array(
        '田中', '佐藤', '鈴木', '高橋', '渡辺', '伊藤', '山田', '中村', '小林', '加藤',
        '吉田', '山口', '松本', '井上', '木村', '林', '清水', '山崎', '森', '池田',
        '橋本', '斎藤', '柴田', '酒井', '藤井', '山本', '近藤', '今井', '石川', '坂本',
        '福田', '太田', '西村', '三浦', '谷口', '小川', '前田', '岡田', '後藤', '長谷川',
        '石井', '村上', '小野', '中川', '原田', '青木', '竹内', '金子', '和田', '中島'
    );
    
    given_names name_array := name_array(
        '太郎', '花子', '一郎', '美咲', '健太', '雅子', '次郎', '真理', '大輝', '由美',
        '正樹', '麻衣', '康夫', '智子', '秀樹', '恵子', '博之', '直美', '光男', '裕子',
        '和也', '亜紀', '哲也', '典子', '義男', '京子', '隆', '美穂', '修', '美香',
        '明', '理恵', '誠', '千恵', '勇', '奈美', '薫', '彩', '徹', '愛',
        '聡', '舞', '豊', '麻美', '武', '香織', '進', '友美', '昭', '純子'
    );

    -- 部署名リスト
    TYPE dept_array IS VARRAY(30) OF VARCHAR2(50);
    dept_names dept_array := dept_array(
        '営業部', '開発部', '人事部', '総務部', 'マーケティング部', '経理部', '法務部', 'IT企画部',
        '技術部', '品質管理部', '生産管理部', '物流部', '広報部', '企画部', '海外事業部',
        '新規事業部', 'データサイエンス部', 'セキュリティ部', 'クラウド事業部', 'モバイル事業部',
        'AI研究部', 'コンサルティング部', 'カスタマーサポート部', '内部監査部', 'リスク管理部',
        '事業企画部', 'システム運用部', 'インフラ部', 'デザイン部', 'プロダクト企画部'
    );

    -- 地域リスト
    TYPE location_array IS VARRAY(20) OF VARCHAR2(20);
    locations location_array := location_array(
        '東京', '大阪', '名古屋', '横浜', '福岡', '札幌', '仙台', '広島', '京都', '神戸',
        '千葉', 'さいたま', '静岡', '新潟', '浜松', '熊本', '鹿児島', '長崎', '大分', '宮崎'
    );

    -- 商品名リスト
    TYPE product_array IS VARRAY(100) OF VARCHAR2(100);
    products product_array := product_array(
        'ノートパソコン Type-A', 'ノートパソコン Type-B', 'デスクトップPC', 'タブレット',
        'スマートフォン', 'モニター 24インチ', 'モニター 27インチ', '4Kモニター',
        'ワイヤレスマウス', '有線マウス', 'ゲーミングマウス', 'キーボード', 'ゲーミングキーボード',
        'ワイヤレスキーボード', 'ヘッドセット', 'スピーカー', 'Webカメラ', 'プリンター',
        '複合機', 'スキャナー', '外付けHDD', 'SSD', 'USBメモリ', 'SD カード',
        'USB-Cハブ', 'USB-Aハブ', '電源アダプター', 'モバイルバッテリー', 'ワイヤレス充電器',
        'LANケーブル', 'HDMIケーブル', 'USB-Cケーブル', 'ディスプレイポートケーブル',
        'ルーター', 'スイッチングハブ', 'アクセスポイント', 'サーバー', 'NAS',
        'UPS', 'ファイアウォール', 'ロードバランサー', 'プロジェクター', 'ホワイトボード',
        'オフィスチェア', 'デスク', '書類ケース', 'シュレッダー', '電話機',
        'ソフトウェアライセンス', 'アンチウイルス', 'Office365', 'Adobe Creative',
        'サポート契約', '保守契約', '技術コンサルティング', 'クラウドサービス', 'データベースライセンス'
    );

    -- 会社名リスト
    TYPE company_array IS VARRAY(50) OF VARCHAR2(100);
    companies company_array := company_array(
        '株式会社ABC商事', '有限会社XYZ販売', '株式会社DEF企画', '合同会社GHI物産',
        '株式会社JKL工業', '有限会社MNO貿易', '株式会社PQR技術', '合資会社STU建設',
        '株式会社VWX情報', '有限会社YZA商会', '株式会社BCD製造', '合同会社EFG流通',
        '株式会社HIJ開発', '有限会社KLM企画', '株式会社NOP商事', '合資会社QRS工業',
        '株式会社TUV技術', '有限会社WXY物産', '株式会社ZAB情報', '合同会社CDE商会',
        '株式会社FGH製造', '有限会社IJK流通', '株式会社LMN開発', '合資会社OPQ企画',
        '株式会社RST商事', '有限会社UVW工業', '株式会社XYZ技術', '合同会社ABC物産',
        '株式会社DEF情報', '有限会社GHI商会', '株式会社JKL製造', '合資会社MNO流通',
        '株式会社PQR開発', '有限会社STU企画', '株式会社VWX商事', '合同会社YZA工業',
        '株式会社BCD技術', '有限会社EFG物産', '株式会社HIJ情報', '合資会社KLM商会',
        '株式会社NOP製造', '有限会社QRS流通', '株式会社TUV開発', '合同会社WXY企画',
        '株式会社ZAB商事', '有限会社CDE工業', '株式会社FGH技術', '合資会社IJK物産',
        '株式会社LMN情報', '有限会社OPQ商会'
    );

    v_counter NUMBER := 0;
    v_dept_id NUMBER;
    v_emp_id NUMBER;
    v_order_id NUMBER;
    v_detail_id NUMBER;
    v_customer_id NUMBER;
    v_dept_name VARCHAR2(100);
    v_location VARCHAR2(50);
    v_first_name VARCHAR2(50);
    v_last_name VARCHAR2(50);
    v_company_name VARCHAR2(200);
    v_product_name VARCHAR2(200);

BEGIN
    DBMS_OUTPUT.PUT_LINE('大量ダミーデータ生成開始...');
    
    -- 1. 部署データ生成
    DBMS_OUTPUT.PUT_LINE('部署データ生成中...');
    FOR i IN 1..20 LOOP
        v_dept_name := dept_names(MOD(i-1, dept_names.COUNT) + 1);
        v_location := locations(MOD(i-1, locations.COUNT) + 1);
        IF i > dept_names.COUNT THEN
            v_dept_name := v_dept_name || '_' || TO_CHAR(CEIL(i / dept_names.COUNT));
        END IF;
        
        INSERT INTO departments (
            department_id, 
            department_name, 
            location,
            created_at,
            updated_at
        ) VALUES (
            seq_departments.NEXTVAL,
            v_dept_name,
            v_location,
            SYSDATE - DBMS_RANDOM.VALUE(0, 365),
            SYSDATE
        );
        
        v_counter := v_counter + 1;
        IF MOD(v_counter, 5) = 0 THEN
            COMMIT;
        END IF;
    END LOOP;
    
    DBMS_OUTPUT.PUT_LINE('部署データ生成完了: ' || v_counter || '件');
    v_counter := 0;
    
    -- 2. 社員データ生成
    DBMS_OUTPUT.PUT_LINE('社員データ生成中...');
    FOR dept IN (SELECT department_id FROM departments WHERE department_id > 8) LOOP
        FOR i IN 1..50 LOOP
            v_first_name := given_names(MOD(v_counter + 13, given_names.COUNT) + 1);
            v_last_name := surnames(MOD(v_counter, surnames.COUNT) + 1);
            
            INSERT INTO employees (
                employee_id,
                first_name,
                last_name,
                email,
                department_id,
                salary,
                hire_date,
                created_at,
                updated_at
            ) VALUES (
                seq_employees.NEXTVAL,
                v_first_name,
                v_last_name,
                'employee' || TO_CHAR(v_counter + 1000) || '@company.com',
                dept.department_id,
                ROUND(DBMS_RANDOM.VALUE(3000000, 8000000), -4),
                SYSDATE - DBMS_RANDOM.VALUE(30, 2500),
                SYSDATE - DBMS_RANDOM.VALUE(0, 30),
                SYSDATE
            );
            
            v_counter := v_counter + 1;
            IF MOD(v_counter, 100) = 0 THEN
                COMMIT;
                DBMS_OUTPUT.PUT_LINE('社員データ生成中... ' || v_counter || '件完了');
            END IF;
        END LOOP;
    END LOOP;
    
    DBMS_OUTPUT.PUT_LINE('社員データ生成完了: ' || v_counter || '件');
    v_counter := 0;
    
    -- 3. 受注データ生成
    DBMS_OUTPUT.PUT_LINE('受注データ生成中...');
    FOR i IN 1..1000 LOOP
        v_customer_id := ROUND(DBMS_RANDOM.VALUE(1001, 1050));
        v_company_name := companies(MOD(v_customer_id - 1001, companies.COUNT) + 1);
        
        INSERT INTO orders (
            order_id,
            customer_id,
            customer_name,
            order_date,
            total_amount,
            status,
            created_at,
            updated_at
        ) VALUES (
            seq_orders.NEXTVAL,
            v_customer_id,
            v_company_name,
            SYSDATE - DBMS_RANDOM.VALUE(0, 365),
            ROUND(DBMS_RANDOM.VALUE(50000, 1000000), -3),
            CASE MOD(i, 4)
                WHEN 0 THEN 'COMPLETED'
                WHEN 1 THEN 'PENDING'
                WHEN 2 THEN 'PROCESSING'
                ELSE 'CANCELLED'
            END,
            SYSDATE - DBMS_RANDOM.VALUE(0, 30),
            SYSDATE
        );
        
        v_counter := v_counter + 1;
        IF MOD(v_counter, 100) = 0 THEN
            COMMIT;
            DBMS_OUTPUT.PUT_LINE('受注データ生成中... ' || v_counter || '件完了');
        END IF;
    END LOOP;
    
    DBMS_OUTPUT.PUT_LINE('受注データ生成完了: ' || v_counter || '件');
    v_counter := 0;
    
    -- 4. 受注明細データ生成
    DBMS_OUTPUT.PUT_LINE('受注明細データ生成中...');
    FOR ord IN (SELECT order_id FROM orders WHERE order_id > 5) LOOP
        -- 各受注に3-7個の明細を追加
        FOR i IN 1..ROUND(DBMS_RANDOM.VALUE(3, 7)) LOOP
            v_product_name := products(MOD(v_counter, products.COUNT) + 1);
            
            INSERT INTO order_details (
                detail_id,
                order_id,
                product_id,
                product_name,
                quantity,
                unit_price,
                created_at,
                updated_at
            ) VALUES (
                seq_order_details.NEXTVAL,
                ord.order_id,
                2000 + MOD(v_counter, 100) + 1,
                v_product_name,
                ROUND(DBMS_RANDOM.VALUE(1, 10)),
                ROUND(DBMS_RANDOM.VALUE(1000, 100000), -2),
                SYSDATE - DBMS_RANDOM.VALUE(0, 30),
                SYSDATE
            );
            
            v_counter := v_counter + 1;
            IF MOD(v_counter, 500) = 0 THEN
                COMMIT;
                DBMS_OUTPUT.PUT_LINE('受注明細データ生成中... ' || v_counter || '件完了');
            END IF;
        END LOOP;
    END LOOP;
    
    DBMS_OUTPUT.PUT_LINE('受注明細データ生成完了: ' || v_counter || '件');
    
    COMMIT;
    
    -- 統計情報更新
    DBMS_OUTPUT.PUT_LINE('統計情報更新中...');
    DBMS_STATS.GATHER_TABLE_STATS(USER, 'DEPARTMENTS');
    DBMS_STATS.GATHER_TABLE_STATS(USER, 'EMPLOYEES');
    DBMS_STATS.GATHER_TABLE_STATS(USER, 'ORDERS');
    DBMS_STATS.GATHER_TABLE_STATS(USER, 'ORDER_DETAILS');
    
    DBMS_OUTPUT.PUT_LINE('大量ダミーデータ生成完了！');
    
    -- データ件数確認
    DBMS_OUTPUT.PUT_LINE('=== データ件数確認 ===');
    FOR rec IN (
        SELECT 'DEPARTMENTS' as table_name, COUNT(*) as cnt FROM departments
        UNION ALL
        SELECT 'EMPLOYEES', COUNT(*) FROM employees
        UNION ALL
        SELECT 'ORDERS', COUNT(*) FROM orders
        UNION ALL
        SELECT 'ORDER_DETAILS', COUNT(*) FROM order_details
    ) LOOP
        DBMS_OUTPUT.PUT_LINE(rec.table_name || ': ' || rec.cnt || '件');
    END LOOP;
    
EXCEPTION
    WHEN OTHERS THEN
        DBMS_OUTPUT.PUT_LINE('エラー発生: ' || SQLERRM);
        ROLLBACK;
        RAISE;
END;
/
EOF

echo "大量データ生成を実行中..."
echo "（この処理には数分かかる場合があります）"
echo ""

# SQLPlusでSQL実行
sqlplus -S "${ORACLE_USER}/${ORACLE_PASSWORD}@${ORACLE_HOST}:${ORACLE_PORT}/${ORACLE_SERVICE}" @"$TEMP_SQL"

if [ $? -eq 0 ]; then
    echo ""
    echo "============================================"
    echo "大量ダミーデータ生成完了"
    echo "============================================"
    echo ""
    echo "生成されたデータでN+1問題のデモを実行できます。"
    echo "Goアプリケーションを実行してパフォーマンスを比較してください。"
else
    echo ""
    echo "エラー: データ生成に失敗しました。"
    exit 1
fi

# 一時ファイル削除
rm -f "$TEMP_SQL"

echo ""
echo "スクリプト完了"
