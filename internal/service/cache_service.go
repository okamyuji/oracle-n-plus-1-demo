package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"oracle-n-plus-1-demo/config"
	"oracle-n-plus-1-demo/internal/cache"

	"github.com/redis/go-redis/v9"
)

// CacheResult - キャッシュ性能測定結果
type CacheResult struct {
	Method        string        `json:"method"`
	ExecutionTime time.Duration `json:"execution_time"`
	MemoryUsage   int64         `json:"memory_usage_bytes"`
	HitRate       float64       `json:"hit_rate"`
	Description   string        `json:"description"`
}

// CacheService - キャッシュ性能比較サービス
type CacheService struct {
	db                  *sql.DB
	redisClient         *redis.Client
	config              *config.Config
	results             []CacheResult
	performanceAnalyzer *cache.PerformanceAnalyzer
	bufferCache         *cache.OracleBufferCache
	resultCache         *cache.OracleResultCache
}

// NewCacheService - キャッシュサービスのコンストラクタ
func NewCacheService(db *sql.DB, cfg *config.Config) *CacheService {
	// Redis接続を試行（失敗してもサービスは動作する）
	var redisClient *redis.Client
	if cfg.RedisHost != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		})

		// 接続テスト
		ctx := context.Background()
		if err := redisClient.Ping(ctx).Err(); err != nil {
			fmt.Printf("Redis接続に失敗しました（キャッシュ比較はスキップされます）: %v\n", err)
			redisClient = nil
		}
	}

	return &CacheService{
		db:                  db,
		redisClient:         redisClient,
		config:              cfg,
		results:             make([]CacheResult, 0),
		performanceAnalyzer: cache.NewPerformanceAnalyzer(db),
		bufferCache:         cache.NewOracleBufferCache(db),
		resultCache:         cache.NewOracleResultCache(db),
	}
}

// TestOracleInternalCache - Oracle内蔵キャッシュのテスト
func (c *CacheService) TestOracleInternalCache(runs int) error {
	fmt.Println("=== Oracle内蔵キャッシュ詳細性能テスト ===")
	fmt.Printf("新しい詳細分析エンジンを使用（%d回実行）\n\n", runs)

	// 1. 包括的性能分析の実行
	analysisResults, err := c.performanceAnalyzer.PerformComprehensiveAnalysis(runs)
	if err != nil {
		fmt.Printf("包括的分析でエラーが発生しましたが、従来の分析を実行します: %v\n", err)
		// フォールバック：従来の分析を実行
		return c.testOracleInternalCacheFallback(runs)
	}

	// 2. 分析結果の統合
	c.integrateAnalysisResults(analysisResults)

	fmt.Println("\n✅ Oracle内蔵キャッシュ詳細分析が完了しました")
	fmt.Println("より詳細な分析結果は包括的レポートで確認できます")

	return nil
}

// testOracleInternalCacheFallback - 従来のOracle内蔵キャッシュテスト（フォールバック用）
func (c *CacheService) testOracleInternalCacheFallback(runs int) error {
	fmt.Println("=== Oracle内蔵キャッシュ基本性能テスト（フォールバック）===")

	// 1. Database Buffer Cacheテスト
	if err := c.testDatabaseBufferCache(runs); err != nil {
		return fmt.Errorf("database buffer cacheテストでエラー: %w", err)
	}

	// 2. Result Cacheテスト
	if err := c.testResultCache(runs); err != nil {
		return fmt.Errorf("result cacheテストでエラー: %w", err)
	}

	// 3. PL/SQL Function Result Cacheテスト
	if err := c.testPLSQLFunctionCache(runs); err != nil {
		return fmt.Errorf("pl/sql function cacheテストでエラー: %w", err)
	}

	return nil
}

// integrateAnalysisResults - 分析結果をCacheServiceに統合
func (c *CacheService) integrateAnalysisResults(results *cache.AnalysisResults) {
	// Buffer Cache結果の統合
	if results.OracleBufferMetrics != nil {
		c.results = append(c.results, CacheResult{
			Method:        "Oracle_Buffer_Cache_Advanced",
			ExecutionTime: results.OracleBufferMetrics.TestExecutionTime,
			MemoryUsage:   results.OracleBufferMetrics.TotalSizeBytes,
			HitRate:       results.OracleBufferMetrics.HitRatio,
			Description:   "Oracle Database Buffer Cache（高度分析）",
		})
	}

	// Result Cache結果の統合
	if results.OracleResultMetrics != nil {
		c.results = append(c.results, CacheResult{
			Method:        "Oracle_Result_Cache_Advanced",
			ExecutionTime: results.OracleResultMetrics.TestExecutionTime,
			MemoryUsage:   results.OracleResultMetrics.MemoryUsage,
			HitRate:       results.OracleResultMetrics.HitRatio,
			Description:   "Oracle Server Result Cache（高度分析）",
		})
	}

	// 総合効率性メトリクスの統合
	if results.PerformanceComparison != nil && results.PerformanceComparison.EfficiencyMetrics != nil {
		c.results = append(c.results, CacheResult{
			Method:        "Oracle_Integrated_Cache",
			ExecutionTime: (results.OracleBufferMetrics.TestExecutionTime + results.OracleResultMetrics.TestExecutionTime) / 2,
			MemoryUsage:   results.OracleBufferMetrics.TotalSizeBytes + results.OracleResultMetrics.MemoryUsage,
			HitRate:       results.PerformanceComparison.EfficiencyMetrics.OverallCacheEfficiency,
			Description:   fmt.Sprintf("Oracle統合キャッシュ（総合効率%.1f%%）", results.PerformanceComparison.EfficiencyMetrics.OverallCacheEfficiency),
		})
	}
}

// testDatabaseBufferCache - Database Buffer Cacheの性能テスト
func (c *CacheService) testDatabaseBufferCache(runs int) error {
	fmt.Println("\n--- Database Buffer Cache テスト ---")

	var totalDuration time.Duration
	var hitCount int64

	for i := 0; i < runs; i++ {
		start := time.Now()

		// 複数回同じデータにアクセスしてBuffer Cacheの効果を測定
		query := `
			SELECT o.order_id, o.customer_id, o.total_amount,
			       od.detail_id, od.product_id, od.quantity
			FROM orders o
			JOIN order_details od ON o.order_id = od.order_id
			WHERE o.order_date >= SYSDATE - 7
			AND ROWNUM <= 100`

		rows, err := c.db.Query(query)
		if err != nil {
			return fmt.Errorf("database buffer cacheクエリでエラー: %w", err)
		}

		var count int
		for rows.Next() {
			var orderID, customerID, detailID, productID int64
			var totalAmount float64
			var quantity int

			err := rows.Scan(&orderID, &customerID, &totalAmount, &detailID, &productID, &quantity)
			if err != nil {
				if cerr := rows.Close(); cerr != nil {
					fmt.Printf("rows.Close() failed: %v\n", cerr)
				}
				return fmt.Errorf("スキャンエラー: %w", err)
			}
			count++
		}
		if cerr := rows.Close(); cerr != nil {
			fmt.Printf("rows.Close() failed: %v\n", cerr)
		}

		duration := time.Since(start)
		totalDuration += duration

		if i > 0 && duration < 50*time.Millisecond { // 2回目以降で高速ならキャッシュヒット
			hitCount++
		}

		if i == 0 {
			fmt.Printf("初回実行時間: %v (キャッシュなし)\n", duration)
		} else if i < 3 {
			fmt.Printf("%d回目実行時間: %v\n", i+1, duration)
		}
	}

	avgDuration := totalDuration / time.Duration(runs)
	hitRate := float64(hitCount) / float64(runs-1) * 100

	c.results = append(c.results, CacheResult{
		Method:        "Oracle_Buffer_Cache",
		ExecutionTime: avgDuration,
		MemoryUsage:   0, // Buffer Cacheのサイズは別途取得
		HitRate:       hitRate,
		Description:   "Oracle Database Buffer Cache（データブロックキャッシュ）",
	})

	fmt.Printf("平均実行時間: %v\n", avgDuration)
	fmt.Printf("推定キャッシュヒット率: %.1f%%\n", hitRate)

	// Buffer Cache統計の取得
	if err := c.getBufferCacheStats(); err != nil {
		fmt.Printf("Buffer Cache統計取得でエラー: %v\n", err)
	}

	return nil
}

// testResultCache - Result Cacheの性能テスト
func (c *CacheService) testResultCache(runs int) error {
	fmt.Println("\n--- Oracle Result Cache テスト ---")

	// Result Cacheヒント付きクエリ
	query := `
		SELECT /*+ RESULT_CACHE */
		       customer_id, COUNT(*) as order_count,
		       SUM(total_amount) as total_sales
		FROM orders
		WHERE order_date >= SYSDATE - 30
		GROUP BY customer_id
		ORDER BY total_sales DESC`

	var totalDuration time.Duration

	for i := 0; i < runs; i++ {
		start := time.Now()

		rows, err := c.db.Query(query)
		if err != nil {
			return fmt.Errorf("result cacheクエリでエラー: %w", err)
		}

		var count int
		for rows.Next() {
			var customerID int64
			var orderCount int
			var totalSales float64

			err := rows.Scan(&customerID, &orderCount, &totalSales)
			if err != nil {
				if cerr := rows.Close(); cerr != nil {
					fmt.Printf("rows.Close() failed: %v\n", cerr)
				}
				return fmt.Errorf("スキャンエラー: %w", err)
			}
			count++
		}
		if cerr := rows.Close(); cerr != nil {
			fmt.Printf("rows.Close() failed: %v\n", cerr)
		}

		duration := time.Since(start)
		totalDuration += duration

		if i == 0 {
			fmt.Printf("初回実行時間: %v (キャッシュなし)\n", duration)
		} else if i < 3 {
			fmt.Printf("%d回目実行時間: %v\n", i+1, duration)
		}
	}

	avgDuration := totalDuration / time.Duration(runs)

	c.results = append(c.results, CacheResult{
		Method:        "Oracle_Result_Cache",
		ExecutionTime: avgDuration,
		MemoryUsage:   0,
		HitRate:       0, // Result Cache統計から後で取得
		Description:   "Oracle Server Result Cache（クエリ結果キャッシュ）",
	})

	fmt.Printf("平均実行時間: %v\n", avgDuration)

	// Result Cache統計の取得
	if err := c.getResultCacheStats(); err != nil {
		fmt.Printf("Result Cache統計取得でエラー: %v\n", err)
	}

	return nil
}

// testPLSQLFunctionCache - PL/SQL Function Result Cacheの性能テスト
func (c *CacheService) testPLSQLFunctionCache(runs int) error {
	fmt.Println("\n--- PL/SQL Function Result Cache テスト ---")

	// Function Result Cache付きファンクションを作成
	createFunctionSQL := `
		CREATE OR REPLACE FUNCTION get_customer_order_summary(p_customer_id NUMBER)
		RETURN VARCHAR2
		RESULT_CACHE RELIES_ON (orders)
		IS
			l_summary VARCHAR2(1000);
		BEGIN
			SELECT 'Orders: ' || COUNT(*) || ', Total: $' || ROUND(SUM(total_amount), 2)
			INTO l_summary
			FROM orders
			WHERE customer_id = p_customer_id
			AND order_date >= SYSDATE - 90;
			
			RETURN l_summary;
		EXCEPTION
			WHEN NO_DATA_FOUND THEN
				RETURN 'No orders found';
		END;`

	_, err := c.db.Exec(createFunctionSQL)
	if err != nil {
		fmt.Printf("PL/SQL関数作成でエラー（スキップ）: %v\n", err)
		return nil
	}

	var totalDuration time.Duration

	for i := 0; i < runs; i++ {
		start := time.Now()

		// 複数の顧客IDで関数を呼び出し
		for customerID := 1; customerID <= 10; customerID++ {
			var result string
			err := c.db.QueryRow("SELECT get_customer_order_summary(?) FROM DUAL", customerID).Scan(&result)
			if err != nil {
				continue // エラーは無視して続行
			}
		}

		duration := time.Since(start)
		totalDuration += duration

		if i == 0 {
			fmt.Printf("初回実行時間: %v (キャッシュなし)\n", duration)
		} else if i < 3 {
			fmt.Printf("%d回目実行時間: %v\n", i+1, duration)
		}
	}

	avgDuration := totalDuration / time.Duration(runs)

	c.results = append(c.results, CacheResult{
		Method:        "Oracle_Function_Cache",
		ExecutionTime: avgDuration,
		MemoryUsage:   0,
		HitRate:       0,
		Description:   "Oracle PL/SQL Function Result Cache",
	})

	fmt.Printf("平均実行時間: %v\n", avgDuration)

	return nil
}

// TestExternalCache - 外部キャッシュ（Redis）のテスト
func (c *CacheService) TestExternalCache(runs int) error {
	if c.redisClient == nil {
		fmt.Println("=== 外部キャッシュ（Redis）テスト ===")
		fmt.Println("Redis接続が利用できないため、外部キャッシュテストをスキップします。")
		return nil
	}

	fmt.Println("=== 外部キャッシュ（Redis）性能テスト ===")

	if err := c.testRedisCache(runs); err != nil {
		return fmt.Errorf("redisキャッシュテストでエラー: %w", err)
	}

	return nil
}

// testRedisCache - Redisキャッシュの性能テスト
func (c *CacheService) testRedisCache(runs int) error {
	fmt.Println("\n--- Redis外部キャッシュ テスト ---")

	ctx := context.Background()
	var totalDuration time.Duration
	var hitCount int64

	// テストデータの準備
	testQuery := `
		SELECT o.order_id, o.customer_id, o.total_amount,
		       od.detail_id, od.product_id, od.quantity
		FROM orders o
		JOIN order_details od ON o.order_id = od.order_id
		WHERE o.order_date >= SYSDATE - 7
		AND ROWNUM <= 100`

	for i := 0; i < runs; i++ {
		start := time.Now()

		cacheKey := "orders_with_details_last_7_days"

		// Redisからキャッシュ取得を試行
		cachedData, err := c.redisClient.Get(ctx, cacheKey).Result()
		if err == redis.Nil {
			// キャッシュミス：データベースから取得してキャッシュに保存
			rows, err := c.db.Query(testQuery)
			if err != nil {
				return fmt.Errorf("データベースクエリでエラー: %w", err)
			}

			var results []map[string]interface{}
			for rows.Next() {
				var orderID, customerID, detailID, productID int64
				var totalAmount float64
				var quantity int

				err := rows.Scan(&orderID, &customerID, &totalAmount, &detailID, &productID, &quantity)
				if err != nil {
					if cerr := rows.Close(); cerr != nil {
						fmt.Printf("rows.Close() failed: %v\n", cerr)
					}
					return fmt.Errorf("スキャンエラー: %w", err)
				}

				result := map[string]interface{}{
					"order_id":     orderID,
					"customer_id":  customerID,
					"total_amount": totalAmount,
					"detail_id":    detailID,
					"product_id":   productID,
					"quantity":     quantity,
				}
				results = append(results, result)
			}
			if cerr := rows.Close(); cerr != nil {
				fmt.Printf("rows.Close() failed: %v\n", cerr)
			}

			// Redisにキャッシュ
			jsonData, err := json.Marshal(results)
			if err != nil {
				return fmt.Errorf("JSON変換エラー: %w", err)
			}

			err = c.redisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute).Err()
			if err != nil {
				return fmt.Errorf("redisキャッシュ保存エラー: %w", err)
			}

			if i == 0 {
				fmt.Printf("初回実行時間: %v (データベース + キャッシュ保存)\n", time.Since(start))
			}
		} else if err != nil {
			return fmt.Errorf("redisアクセスエラー: %w", err)
		} else {
			// キャッシュヒット：Redisからデータを取得
			var results []map[string]interface{}
			err := json.Unmarshal([]byte(cachedData), &results)
			if err != nil {
				return fmt.Errorf("JSON解析エラー: %w", err)
			}

			hitCount++
			if i < 3 {
				fmt.Printf("%d回目実行時間: %v (キャッシュヒット)\n", i+1, time.Since(start))
			}
		}

		duration := time.Since(start)
		totalDuration += duration
	}

	avgDuration := totalDuration / time.Duration(runs)
	hitRate := float64(hitCount) / float64(runs) * 100

	c.results = append(c.results, CacheResult{
		Method:        "Redis_External_Cache",
		ExecutionTime: avgDuration,
		MemoryUsage:   0, // Redis使用量は別途取得
		HitRate:       hitRate,
		Description:   "Redis外部キャッシュ（JSONシリアライゼーション）",
	})

	fmt.Printf("平均実行時間: %v\n", avgDuration)
	fmt.Printf("キャッシュヒット率: %.1f%%\n", hitRate)

	// Redis使用量の取得
	if err := c.getRedisMemoryUsage(); err != nil {
		fmt.Printf("Redis使用量取得でエラー: %v\n", err)
	}

	return nil
}

// DisplayCacheComparison - キャッシュ比較結果を表示
func (c *CacheService) DisplayCacheComparison() error {
	if len(c.results) == 0 {
		return fmt.Errorf("比較結果がありません")
	}

	fmt.Println("\n=== キャッシュ性能比較結果 ===")

	// 結果を表形式で表示
	fmt.Printf("%-25s | %-15s | %-10s | %s\n", "キャッシュ方式", "平均実行時間", "ヒット率", "説明")
	fmt.Println(strings.Repeat("-", 80))

	for _, result := range c.results {
		hitRateStr := fmt.Sprintf("%.1f%%", result.HitRate)
		if result.HitRate == 0 {
			hitRateStr = "N/A"
		}

		fmt.Printf("%-25s | %-15v | %-10s | %s\n",
			result.Method,
			result.ExecutionTime,
			hitRateStr,
			result.Description)
	}

	// 性能分析
	c.analyzePerformance()

	return nil
}

// analyzePerformance - 性能分析とOracle内蔵キャッシュの優位性を説明
func (c *CacheService) analyzePerformance() {
	fmt.Println("\n=== 性能分析結果 ===")

	oracleResults := make([]CacheResult, 0)
	var redisResult *CacheResult

	for _, result := range c.results {
		if strings.HasPrefix(result.Method, "Oracle_") {
			oracleResults = append(oracleResults, result)
		} else if result.Method == "Redis_External_Cache" {
			redisResult = &result
		}
	}

	fmt.Println("\n1. Oracle内蔵キャッシュの優位性:")
	fmt.Println("   ✓ データの移動が不要（メモリ効率）")
	fmt.Println("   ✓ シリアライゼーション/デシリアライゼーション不要")
	fmt.Println("   ✓ ネットワークI/Oなし")
	fmt.Println("   ✓ 自動的なキャッシュ無効化とデータ整合性")
	fmt.Println("   ✓ 複数レベルのキャッシュ（Buffer Cache + Result Cache + Function Cache）")

	if redisResult != nil && len(oracleResults) > 0 {
		fmt.Println("\n2. 外部キャッシュ（Redis）の課題:")
		fmt.Println("   ✗ ネットワーク通信のオーバーヘッド")
		fmt.Println("   ✗ JSONシリアライゼーション/デシリアライゼーションのコスト")
		fmt.Println("   ✗ データ整合性管理の複雑さ")
		fmt.Println("   ✗ 追加のインフラストラクチャとメンテナンス")
		fmt.Println("   ✗ メモリの二重使用（Oracle + Redis）")

		// 最速のOracle結果と比較
		var fastestOracle CacheResult
		for i, result := range oracleResults {
			if i == 0 || result.ExecutionTime < fastestOracle.ExecutionTime {
				fastestOracle = result
			}
		}

		if fastestOracle.ExecutionTime < redisResult.ExecutionTime {
			improvement := float64(redisResult.ExecutionTime.Nanoseconds()) / float64(fastestOracle.ExecutionTime.Nanoseconds())
			fmt.Printf("\n3. 性能比較結果:\n")
			fmt.Printf("   Oracle内蔵キャッシュ(%s): %v\n", fastestOracle.Method, fastestOracle.ExecutionTime)
			fmt.Printf("   Redis外部キャッシュ: %v\n", redisResult.ExecutionTime)
			fmt.Printf("   Oracle内蔵キャッシュが%.1fx高速\n", improvement)
		}
	}

	fmt.Println("\n4. 推奨事項:")
	fmt.Println("   → Oracle環境ではDatabase固有のキャッシュメカニズムを最大限活用する")
	fmt.Println("   → 外部キャッシュは以下の場合のみ検討:")
	fmt.Println("     - マイクロサービス間でのデータ共有")
	fmt.Println("     - 外部APIからの取得データ")
	fmt.Println("     - Oracleでカバーできない計算集約的な結果")
	fmt.Println("   → N+1問題はSQL設計の改善で根本的に解決する")
}

// DisplayMemoryUsageComparison - メモリ使用量比較を表示
func (c *CacheService) DisplayMemoryUsageComparison() error {
	fmt.Println("\n=== メモリ使用量分析 ===")

	// Oracle SGA情報の取得
	if err := c.getOracleSGAInfo(); err != nil {
		fmt.Printf("Oracle SGA情報取得でエラー: %v\n", err)
	}

	fmt.Println("\nメモリ効率性:")
	fmt.Println("Oracle内蔵キャッシュ:")
	fmt.Println("  ✓ SGAで統合管理されたメモリ使用")
	fmt.Println("  ✓ 自動的なメモリ最適化とパージ")
	fmt.Println("  ✓ 複数のアプリケーションで共有")

	if c.redisClient != nil {
		fmt.Println("\nRedis外部キャッシュ:")
		fmt.Println("  ✗ 専用メモリプール必要")
		fmt.Println("  ✗ Oracle + Redis = メモリの二重使用")
		fmt.Println("  ✗ 手動でのメモリ管理とTuning")
	}

	return nil
}

// getBufferCacheStats - Buffer Cache統計を取得
func (c *CacheService) getBufferCacheStats() error {
	query := `
		SELECT ROUND((1 - (phy.value / (cur.value + con.value))) * 100, 2) as buffer_hit_ratio
		FROM V$SYSSTAT phy, V$SYSSTAT cur, V$SYSSTAT con
		WHERE phy.name = 'physical reads cache'
		AND cur.name = 'db block gets from cache'
		AND con.name = 'consistent gets from cache'`

	var hitRatio float64
	err := c.db.QueryRow(query).Scan(&hitRatio)
	if err != nil {
		return err
	}

	fmt.Printf("Database Buffer Cache ヒット率: %.2f%%\n", hitRatio)
	return nil
}

// getResultCacheStats - Result Cache統計を取得
func (c *CacheService) getResultCacheStats() error {
	query := `
		SELECT COUNT(*) as cached_objects,
		       SUM(block_count) as total_blocks
		FROM V$RESULT_CACHE_OBJECTS
		WHERE type = 'Result'`

	var objectCount, blockCount int
	err := c.db.QueryRow(query).Scan(&objectCount, &blockCount)
	if err != nil {
		fmt.Printf("Result Cache統計は利用できません: %v\n", err)
		return nil
	}

	fmt.Printf("Result Cache: %d個のオブジェクト, %dブロック使用\n", objectCount, blockCount)
	return nil
}

// getOracleSGAInfo - Oracle SGA情報を取得
func (c *CacheService) getOracleSGAInfo() error {
	query := `
		SELECT component, current_size/1024/1024 as size_mb
		FROM V$SGA_DYNAMIC_COMPONENTS
		WHERE component IN ('DEFAULT buffer_pool', 'Shared Pool', 'Result Cache')
		ORDER BY current_size DESC`

	rows, err := c.db.Query(query)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			fmt.Printf("rows.Close() failed: %v\n", cerr)
		}
	}()

	fmt.Println("\nOracle SGA構成:")
	for rows.Next() {
		var component string
		var sizeMB float64

		err := rows.Scan(&component, &sizeMB)
		if err != nil {
			continue
		}

		fmt.Printf("  %s: %.1f MB\n", component, sizeMB)
	}

	return nil
}

// getRedisMemoryUsage - Redis使用量を取得
func (c *CacheService) getRedisMemoryUsage() error {
	if c.redisClient == nil {
		return nil
	}

	ctx := context.Background()
	info, err := c.redisClient.Info(ctx, "memory").Result()
	if err != nil {
		return err
	}

	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "used_memory_human:") {
			memUsage := strings.TrimSpace(strings.Split(line, ":")[1])
			fmt.Printf("Redis使用メモリ: %s\n", memUsage)
			break
		}
	}

	return nil
}
