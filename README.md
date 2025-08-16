# Oracle N+1問題デモンストレーション - Go実装

Oracle Databaseを使用したN+1問題の実例とその解決策を学習するためのGoアプリケーションです。業務アプリ開発でよく発生するN+1問題の原因、影響、そして効果的な解決策をOracle環境で実践的にデモンストレーションします。

## 目次

- [概要](#概要)
- [N+1問題とは](#n1問題とは)
- [プロジェクト構成](#プロジェクト構成)
- [セットアップ](#セットアップ)
- [使用方法](#使用方法)
- [実装内容](#実装内容)
- [パフォーマンス比較](#パフォーマンス比較)
- [技術的詳細](#技術的詳細)
- [Oracle固有の最適化](#oracle固有の最適化)

## 概要

このプロジェクトは、Oracle Databaseを使用した業務アプリ開発において頻繁に発生するN+1問題を、実際のコード例とパフォーマンス測定を通して学習できるデモアプリケーションです。

### 特徴

- **実践的なコード例**: 問題のあるコードと修正されたコードの両方を提供
- **パフォーマンス測定**: 実際の実行時間を比較測定
- **Oracle最適化**: Oracle Database固有の最適化手法を活用
- **包括的な解決策**: JOIN、IN句、バッチ処理など複数のアプローチを実装
- **キャッシュ性能比較**: Oracle内蔵キャッシュ vs Redis外部キャッシュの性能分析
- **大量データ生成**: 実際の業務環境を模擬した大量ダミーデータでのテスト

## N+1問題とは

N+1問題は、データベースアクセスにおける代表的なパフォーマンス問題です：

1. **初回クエリ（1回）**: 親データを取得
2. **追加クエリ（N回）**: 各親データに対して子データを個別取得

結果として「1 + N回」のクエリが実行され、パフォーマンスが大幅に劣化します。

### 問題の例

```go
// ❌ N+1問題のあるコード
orders := getOrders()              // 1回のクエリ
for _, order := range orders {     // N回のループ
    details := getOrderDetails(order.ID)  // N回のクエリ
    // 処理...
}
```

### 解決策

```go
// ✅ 最適化されたコード
ordersWithDetails := getOrdersWithDetailsJoin()  // 1回のクエリ
for _, orderWithDetails := range ordersWithDetails {
    // 処理...
}
```

## プロジェクト構成

```shell
oracle-n-plus-1-demo/
├── cmd/
│   └── main.go                # メインアプリケーション
├── go.mod                     # Go modules設定
├── go.sum                     # 依存関係のチェックサム
├── env.example                # 環境変数のサンプル
├── linter.sh                  # リンター実行スクリプト
├── README.md                  # このファイル
├── config/
│   └── config.go              # 設定管理とDB接続
├── internal/
│   ├── cache/                 # キャッシュ機能実装
│   │   ├── cache_analyzer.go   # キャッシュ性能分析
│   │   ├── oracle_buffer_cache.go # Buffer Cache実装
│   │   └── oracle_result_cache.go # Result Cache実装
│   └── service/
│       ├── cache_service.go    # キャッシュサービス
│       └── demo_service.go     # デモサービス
├── models/
│   └── models.go              # データモデル定義
├── repository/
│   ├── repository_problem.go  # N+1問題のあるリポジトリ
│   └── repository_optimized.go # 最適化されたリポジトリ
└── scripts/
    ├── ddl/
    │   └── create_tables.sql   # テーブル作成DDL
    ├── dml/
    │   └── insert_initial_data.sql # 初期データDML
    └── load_test_data.sh       # 大量ダミーデータ生成スクリプト
```

## セットアップ

### 前提条件

- Go 1.19以上
- Oracle Database 12c以上
- Oracle Client（SQL*Plus推奨）

### 1. プロジェクトのクローン

```bash
git clone https://github.com/okamyuji/oracle-n-plus-1-demo
cd oracle-n-plus-1-demo
```

### 2. 依存関係のインストール

```bash
go mod tidy
```

### 3. データベースのセットアップ

```bash
# DDLでテーブル作成
sqlplus username/password@hostname:1521/service_name @scripts/ddl/create_tables.sql

# 初期データの投入
sqlplus username/password@hostname:1521/service_name @scripts/dml/insert_initial_data.sql

# 大量のダミーデータ生成（負荷テスト用）
export ORACLE_USER=your_username
export ORACLE_PASSWORD=your_password
export ORACLE_HOST=your_host
export ORACLE_SERVICE=your_service
./scripts/load_test_data.sh
```

### 4. 環境設定

`.env`ファイルを作成し、Oracle接続情報を設定：

```bash
cp .env.example .env
```

`.env`ファイルを編集：

```env
DB_HOST=localhost
DB_PORT=1521
DB_SERVICE_NAME=ORCLPDB1
DB_USERNAME=your_username
DB_PASSWORD=your_password
```

## 使用方法

### 基本的な実行

```bash
# 全てのパフォーマンステストを実行
go run cmd/main.go

# 統計情報とサンプルデータを表示
go run cmd/main.go -stats -sample

# 過去7日間の受注データでテスト
go run cmd/main.go -days=7
```

### オプション

- `-days=30`: 取得する受注データの日数（デフォルト: 30日）
- `-sample`: サンプルデータを表示
- `-stats`: データベース統計情報を表示
- `-order-only`: 受注データのパフォーマンステストのみ実行
- `-employee-only`: 社員データのパフォーマンステストのみ実行
- `--cache-only`: キャッシュ性能比較テストのみ実行
- `-help`: ヘルプを表示

### 使用例

```bash
# 受注データのみテスト、統計情報表示
go run cmd/main.go -order-only -stats

# 社員データのみテスト
go run cmd/main.go -employee-only

# キャッシュ性能比較テストのみ実行
go run cmd/main.go --cache-only

# 詳細情報付きで全テスト実行
go run cmd/main.go -days=7 -sample -stats
```

## 実装内容

### 1. 問題のあるアプローチ（N+1問題）

**受注管理システムでの例:**

```go
func (r *ProblemOrderRepository) GetOrdersWithDetails(days int) ([]models.OrderWithDetails, error) {
    // 1. 受注一覧を取得（1回のクエリ）
    orders, err := r.GetOrdersByDays(days)
    if err != nil {
        return nil, err
    }

    var result []models.OrderWithDetails
    // 2. 各受注ごとに明細を取得（N回のクエリ）
    for _, order := range orders {
        details, err := r.GetDetailsByOrderID(order.OrderID)  // N+1問題発生！
        if err != nil {
            return nil, err
        }
        result = append(result, models.OrderWithDetails{
            Order:   order,
            Details: details,
        })
    }

    return result, nil
}
```

### 2. 最適化されたアプローチ

#### 解決策1: JOINを使用した一括取得

```go
func (r *OptimizedOrderRepository) GetOrdersWithDetailsJoin(days int) ([]models.OrderWithDetails, error) {
    query := `
        SELECT 
            o.order_id, o.customer_id, o.order_date, o.total_amount,
            od.detail_id, od.product_id, od.quantity, od.unit_price
        FROM orders o
        LEFT JOIN order_details od ON o.order_id = od.order_id
        WHERE o.order_date >= SYSDATE - :1
        ORDER BY o.order_id, od.detail_id`

    // 1回のクエリで全データを取得
    rows, err := r.db.Query(query, days)
    // ... 結果の処理
}
```

#### 解決策2: IN句を使用したバッチ取得

```go
func (r *OptimizedOrderRepository) GetOrdersWithDetailsBatch(days int) ([]models.OrderWithDetails, error) {
    // 1. 受注一覧を取得
    orders, err := r.GetOrdersByDays(days)
    if err != nil {
        return nil, err
    }

    // 2. 受注IDを抽出
    orderIDs := extractOrderIDs(orders)

    // 3. IN句で明細を一括取得
    allDetails, err := r.GetDetailsByOrderIDs(orderIDs)
    if err != nil {
        return nil, err
    }

    // 4. メモリ上でグルーピング
    return groupOrdersWithDetails(orders, allDetails), nil
}
```

### 3. パフォーマンス測定機能

各アプローチの実行時間を測定し、改善効果を定量的に評価：

```go
type PerformanceResult struct {
    Method        string        `json:"method"`
    ExecutionTime time.Duration `json:"execution_time"`
    RecordCount   int           `json:"record_count"`
    Description   string        `json:"description"`
}
```

## パフォーマンス比較

### 最新の実測結果（大量データでのテスト）

**受注データ（83件）でのN+1問題テスト:**

| アプローチ | 実行時間 | SQL実行回数 | 改善率 |
|-----------|----------|-------------|---------|
| N+1問題あり | 32.16ms | 84回 | 基準 |
| JOIN最適化 | 2.97ms | 1回 | **10.8x高速** |
| バッチ取得 | 4.94ms | 2回 | **6.5x高速** |

**社員データ（1,010件）でのN+1問題テスト:**

| アプローチ | 実行時間 | SQL実行回数 | 改善率 |
|-----------|----------|-------------|---------|
| N+1問題あり | 275.27ms | 1,011回 | 基準 |
| JOIN最適化 | 3.24ms | 1回 | **85.0x高速** |
| バッチ取得 | 2.68ms | 2回 | **102.6x高速** |

### キャッシュ性能比較結果

**Oracle内蔵キャッシュ vs Redis外部キャッシュ:**

| キャッシュ方式 | 平均実行時間 | ヒット率 | 特徴 |
|---|---|---|---|
| **Redis外部キャッシュ** | **522.9µs** | 90.0% | 最高速度（キャッシュヒット時） |
| **Oracle Result Cache** | **852.2µs** | N/A | SQL結果キャッシュ |
| **Oracle Buffer Cache** | **1.18ms** | 100.0% | データブロックキャッシュ |
| **Oracle Function Cache** | **5.03ms** | N/A | PL/SQL関数キャッシュ |

### データ量による影響の特徴

データ量が増加するにつれて、N+1問題の影響は指数関数的に悪化します：

- **小規模（83件）**: 約10x～12xの改善
- **中規模（1,010件）**: 約85x～103xの改善  
- **大規模（10,000件以上）**: 推定200x～500xの改善

**重要なポイント:**

- データ量が12倍増えると、N+1問題の実行時間は約8.6倍に増加
- 解決策の効果はデータ量の増加に比例して顕著になる
- Oracle内蔵キャッシュは運用負荷が最小で実用的

## 技術的詳細

### 使用技術

- **言語**: Go 1.19+
- **ORMなし**: database/sqlパッケージを直接使用
- **Oracle Driver**: [sijms/go-ora](https://github.com/sijms/go-ora)
- **設定管理**: 環境変数 + .envファイル

### データベース設計

**テーブル構成:**

1. **orders（受注）**
   - order_id (PK)
   - customer_id  
   - order_date
   - total_amount

2. **order_details（受注明細）**
   - detail_id (PK)
   - order_id (FK)
   - product_id
   - quantity
   - unit_price

3. **employees（社員）**
   - employee_id (PK)
   - first_name, last_name
   - email
   - department_id (FK)
   - hire_date, salary

4. **departments（部署）**
   - department_id (PK)
   - department_name
   - location

### インデックス戦略

パフォーマンス最適化のため、以下のインデックスを作成：

```sql
-- 外部キー用インデックス
CREATE INDEX idx_order_details_order_id ON order_details(order_id);
CREATE INDEX idx_employees_department_id ON employees(department_id);

-- 検索条件用インデックス  
CREATE INDEX idx_orders_order_date ON orders(order_date);

-- 複合インデックス
CREATE INDEX idx_orders_date_customer ON orders(order_date, customer_id);
```

## Oracle固有の最適化

### 1. Result Cacheの活用

```sql
-- クエリ結果のキャッシュ
SELECT /*+ RESULT_CACHE */
       customer_id, COUNT(*) as order_count
FROM orders 
WHERE order_date >= SYSDATE - 30
GROUP BY customer_id;
```

### 2. Database Buffer Cacheの最適化

```sql
-- バッファキャッシュのヒット率確認
SELECT ROUND((1 - (phy.value / (cur.value + con.value))) * 100, 2) as buffer_hit_ratio
FROM V$SYSSTAT phy, V$SYSSTAT cur, V$SYSSTAT con
WHERE phy.name = 'physical reads cache'
  AND cur.name = 'db block gets from cache'  
  AND con.name = 'consistent gets from cache';
```

### 3. PL/SQL Function Result Cache

```sql
CREATE OR REPLACE FUNCTION get_department_name(p_dept_id NUMBER)
RETURN VARCHAR2
RESULT_CACHE RELIES_ON (departments)
IS
    l_name VARCHAR2(100);
BEGIN
    SELECT department_name INTO l_name
    FROM departments
    WHERE department_id = p_dept_id;
    
    RETURN l_name;
END;
```

### 4. Client Result Cacheの設定

```go
// Go側での設定
props := map[string]string{
    "oracle.jdbc.enableResultSetCache": "true",
    "oracle.jdbc.resultSetCacheSize":   "10485760", // 10MB
}
```

## 開発・運用での注意点

### 1. 開発時の品質管理

```bash
# リンター実行（開発時必須）
./linter.sh

# アプリケーションのビルド確認
go build -o oracle-n-plus-1-demo cmd/main.go

# 単体テスト実行
go test ./...
```

### 2. 監視ポイント

- **SQLトレース**: 実行されるSQL文の監視
- **実行計画**: インデックスの使用状況確認
- **バッファキャッシュ**: ヒット率の監視
- **待機イベント**: I/O待機の分析

### 3. コードレビューでのチェック項目

```go
// ❌ 避けるべきパターン
for _, item := range items {
    relatedData := repository.GetByID(item.ID)  // ループ内DBアクセス
    // 処理...
}

// ✅ 推奨パターン  
ids := extractIDs(items)
allRelatedData := repository.GetByIDs(ids)     // 一括取得
relationMap := createRelationMap(allRelatedData)
for _, item := range items {
    relatedData := relationMap[item.ID]
    // 処理...
}
```

### 4. 設計時の考慮事項

- **データアクセスパターンの明確化**
- **バッチサイズの適切な設定**
- **キャッシュ戦略の計画**
- **インデックス設計の最適化**

## まとめ

このデモアプリケーションを通して、以下のことが学習できます：

1. **N+1問題の実際の影響**: パフォーマンス劣化の定量的な把握
2. **効果的な解決策**: JOIN、IN句、バッチ処理の適切な使い分け
3. **Oracle固有の最適化**: Database固有の機能を活用した高速化
4. **実装のベストプラクティス**: 保守性を保ちながらパフォーマンスを向上させる方法

N+1問題は適切な対策により大幅なパフォーマンス改善が可能です。このプロジェクトを参考に、実際の業務アプリ開発でも効果的な最適化を行ってください。

## ライセンス

MIT License

## 貢献

プルリクエストやイシューの報告を歓迎します。改善提案や追加機能のアイデアがあれば、お気軽にご連絡ください。
