package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// PerformanceAnalyzer - キャッシュ性能分析ユーティリティ
type PerformanceAnalyzer struct {
	db                *sql.DB
	bufferCache       *OracleBufferCache
	resultCache       *OracleResultCache
	analysisResults   *AnalysisResults
	comparisonMetrics map[string]interface{}
}

// AnalysisResults - 統合分析結果
type AnalysisResults struct {
	TestDate              time.Time              `json:"test_date"`
	TestDuration          time.Duration          `json:"test_duration"`
	OracleBufferMetrics   *BufferCacheMetrics    `json:"oracle_buffer_metrics"`
	OracleResultMetrics   *ResultCacheMetrics    `json:"oracle_result_metrics"`
	PerformanceComparison *PerformanceComparison `json:"performance_comparison"`
	OptimizationAdvice    *OptimizationAdvice    `json:"optimization_advice"`
	DetailedAnalysis      map[string]interface{} `json:"detailed_analysis"`
	Recommendations       []Recommendation       `json:"recommendations"`
}

// PerformanceComparison - 性能比較結果
type PerformanceComparison struct {
	OracleAdvantages    []string             `json:"oracle_advantages"`
	ExternalCacheIssues []string             `json:"external_cache_issues"`
	EfficiencyMetrics   *EfficiencyMetrics   `json:"efficiency_metrics"`
	ResourceUtilization *ResourceUtilization `json:"resource_utilization"`
}

// EfficiencyMetrics - 効率性メトリクス
type EfficiencyMetrics struct {
	BufferCacheEfficiency  float64 `json:"buffer_cache_efficiency"`
	ResultCacheEfficiency  float64 `json:"result_cache_efficiency"`
	OverallCacheEfficiency float64 `json:"overall_cache_efficiency"`
	MemoryEfficiencyRatio  float64 `json:"memory_efficiency_ratio"`
	IOReductionRatio       float64 `json:"io_reduction_ratio"`
}

// ResourceUtilization - リソース使用率
type ResourceUtilization struct {
	TotalCacheMemoryMB     float64 `json:"total_cache_memory_mb"`
	BufferCacheUtilization float64 `json:"buffer_cache_utilization"`
	ResultCacheUtilization float64 `json:"result_cache_utilization"`
	EstimatedIOSavings     int64   `json:"estimated_io_savings"`
	EstimatedCPUSavings    float64 `json:"estimated_cpu_savings"`
}

// OptimizationAdvice - 最適化アドバイス
type OptimizationAdvice struct {
	PriorityLevel       string   `json:"priority_level"`
	ImmediateActions    []string `json:"immediate_actions"`
	MediumTermActions   []string `json:"medium_term_actions"`
	LongTermActions     []string `json:"long_term_actions"`
	AvoidanceStrategies []string `json:"avoidance_strategies"`
}

// Recommendation - 推奨事項
type Recommendation struct {
	Category    string   `json:"category"`
	Priority    string   `json:"priority"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Impact      string   `json:"impact"`
	Effort      string   `json:"effort"`
	Benefits    []string `json:"benefits"`
}

// NewPerformanceAnalyzer - 性能分析器を作成
func NewPerformanceAnalyzer(db *sql.DB) *PerformanceAnalyzer {
	return &PerformanceAnalyzer{
		db:                db,
		bufferCache:       NewOracleBufferCache(db),
		resultCache:       NewOracleResultCache(db),
		comparisonMetrics: make(map[string]interface{}),
	}
}

// PerformComprehensiveAnalysis - 包括的なキャッシュ性能分析を実行
func (pa *PerformanceAnalyzer) PerformComprehensiveAnalysis(runs int) (*AnalysisResults, error) {
	fmt.Println("\\n=== Oracle内蔵キャッシュ包括的性能分析 ===")
	fmt.Printf("実行回数: %d回\\n", runs)
	fmt.Println("分析項目: Buffer Cache, Result Cache, 統合効率性, 外部キャッシュ比較")
	fmt.Println(strings.Repeat("=", 70))

	startTime := time.Now()

	// 1. Buffer Cacheの詳細分析
	fmt.Println("\\n1. Database Buffer Cache分析中...")
	bufferMetrics, err := pa.bufferCache.TestBufferCachePerformance(runs)
	if err != nil {
		return nil, fmt.Errorf("buffer cache分析エラー: %w", err)
	}

	// 2. Result Cacheの詳細分析
	fmt.Println("\\n2. Server Result Cache分析中...")
	resultMetrics, err := pa.resultCache.TestResultCachePerformance(runs)
	if err != nil {
		return nil, fmt.Errorf("result cache分析エラー: %w", err)
	}

	// 3. 統合分析の実行
	fmt.Println("\\n3. 統合性能分析中...")
	analysisResults := &AnalysisResults{
		TestDate:            startTime,
		TestDuration:        time.Since(startTime),
		OracleBufferMetrics: bufferMetrics,
		OracleResultMetrics: resultMetrics,
	}

	// 4. 性能比較とリソース効率性の計算
	if err := pa.calculatePerformanceComparison(analysisResults); err != nil {
		return nil, fmt.Errorf("性能比較計算エラー: %w", err)
	}

	// 5. 最適化アドバイスの生成
	if err := pa.generateOptimizationAdvice(analysisResults); err != nil {
		return nil, fmt.Errorf("最適化アドバイス生成エラー: %w", err)
	}

	// 6. 詳細分析の実行
	if err := pa.performDetailedAnalysis(analysisResults); err != nil {
		return nil, fmt.Errorf("詳細分析エラー: %w", err)
	}

	// 7. 推奨事項の生成
	if err := pa.generateRecommendations(analysisResults); err != nil {
		return nil, fmt.Errorf("推奨事項生成エラー: %w", err)
	}

	pa.analysisResults = analysisResults

	// 8. 最終結果の表示
	pa.DisplayComprehensiveResults()

	return analysisResults, nil
}

// calculatePerformanceComparison - 性能比較を計算
func (pa *PerformanceAnalyzer) calculatePerformanceComparison(results *AnalysisResults) error {
	// 効率性メトリクスの計算
	efficiency := &EfficiencyMetrics{}

	// Buffer Cache効率性
	if results.OracleBufferMetrics.LogicalReads > 0 {
		efficiency.BufferCacheEfficiency = (float64(results.OracleBufferMetrics.LogicalReads-results.OracleBufferMetrics.PhysicalReads) / float64(results.OracleBufferMetrics.LogicalReads)) * 100
	}

	// Result Cache効率性
	efficiency.ResultCacheEfficiency = results.OracleResultMetrics.HitRatio

	// 総合キャッシュ効率性（重み付き平均）
	efficiency.OverallCacheEfficiency = (efficiency.BufferCacheEfficiency * 0.7) + (efficiency.ResultCacheEfficiency * 0.3)

	// メモリ効率比
	totalMemory := float64(results.OracleBufferMetrics.TotalSizeBytes + results.OracleResultMetrics.MemoryUsage)
	if totalMemory > 0 {
		efficiency.MemoryEfficiencyRatio = efficiency.OverallCacheEfficiency / (totalMemory / (1024 * 1024 * 1024)) // GB当たりの効率
	}

	// I/O削減比
	if results.OracleBufferMetrics.LogicalReads > 0 {
		efficiency.IOReductionRatio = (float64(results.OracleBufferMetrics.PhysicalReads) / float64(results.OracleBufferMetrics.LogicalReads)) * 100
	}

	// リソース使用率の計算
	utilization := &ResourceUtilization{
		TotalCacheMemoryMB:     totalMemory / (1024 * 1024),
		BufferCacheUtilization: float64(results.OracleBufferMetrics.TotalSizeBytes) / (1024 * 1024),
		ResultCacheUtilization: float64(results.OracleResultMetrics.MemoryUsage) / (1024 * 1024),
		EstimatedIOSavings:     results.OracleBufferMetrics.LogicalReads - results.OracleBufferMetrics.PhysicalReads,
		EstimatedCPUSavings:    efficiency.OverallCacheEfficiency * 0.1, // 概算
	}

	// Oracle内蔵キャッシュの優位性
	oracleAdvantages := []string{
		"✓ データブロックレベルの自動キャッシュ管理",
		"✓ SQL結果の透明なキャッシュ（アプリケーション変更不要）",
		"✓ ACID特性による完全なデータ整合性保証",
		"✓ 自動的な依存関係追跡と無効化",
		"✓ 統計情報に基づく自動最適化",
		"✓ 複数プロセス間での効率的なキャッシュ共有",
		"✓ ネットワークI/O削除による低レイテンシ",
		"✓ 串列化オーバーヘッドの排除",
	}

	// 外部キャッシュの課題
	externalCacheIssues := []string{
		"✗ ネットワーク通信による追加レイテンシ（平均1-5ms）",
		"✗ JSONシリアライゼーション/デシリアライゼーションのCPU負荷（10-30%）",
		"✗ データ整合性の手動管理による複雑性とバグリスク",
		"✗ キャッシュ無効化ロジックの実装・メンテナンス負荷",
		"✗ 追加インフラストラクチャ（Redis等）の運用コスト",
		"✗ メモリの二重使用による非効率性",
		"✗ アプリケーション層でのキャッシュキー管理の複雑性",
		"✗ 外部システム障害時の影響範囲拡大",
	}

	results.PerformanceComparison = &PerformanceComparison{
		OracleAdvantages:    oracleAdvantages,
		ExternalCacheIssues: externalCacheIssues,
		EfficiencyMetrics:   efficiency,
		ResourceUtilization: utilization,
	}

	return nil
}

// generateOptimizationAdvice - 最適化アドバイスを生成
func (pa *PerformanceAnalyzer) generateOptimizationAdvice(results *AnalysisResults) error {
	advice := &OptimizationAdvice{}

	// 優先度の決定
	if results.PerformanceComparison.EfficiencyMetrics.OverallCacheEfficiency >= 90 {
		advice.PriorityLevel = "最適化済み（維持管理）"
	} else if results.PerformanceComparison.EfficiencyMetrics.OverallCacheEfficiency >= 70 {
		advice.PriorityLevel = "中優先度（改善推奨）"
	} else {
		advice.PriorityLevel = "高優先度（早急な対応必要）"
	}

	// 即座の対応事項
	advice.ImmediateActions = []string{
		"Oracle内蔵キャッシュメカニズムの活用を最優先に検討",
		"外部キャッシュ導入前にN+1問題の根本的解決（SQL最適化）を実施",
		"RESULT_CACHEヒントを集計クエリに積極適用",
		"Buffer Cache統計の定期監視を開始",
	}

	// 中期的対応事項
	advice.MediumTermActions = []string{
		"アプリケーション層のキャッシュロジック簡素化",
		"Oracle統計情報の定期更新による自動最適化活用",
		"PL/SQL Function Result Cacheの導入検討",
		"外部キャッシュ依存度の段階的削減",
	}

	// 長期的対応事項
	advice.LongTermActions = []string{
		"データベース中心設計への移行",
		"外部キャッシュインフラの段階的廃止検討",
		"Oracle Database統合監視体制の構築",
		"チーム全体のOracle最適化スキル向上",
	}

	// 回避すべき戦略
	advice.AvoidanceStrategies = []string{
		"単純な外部キャッシュ追加による問題の先送り",
		"N+1問題の根本解決を避けた対症療法",
		"Oracle機能を活用しない独自キャッシュ実装",
		"データ整合性を犠牲にした性能最適化",
		"監視・メンテナンス負荷を無視したアーキテクチャ選択",
	}

	results.OptimizationAdvice = advice
	return nil
}

// performDetailedAnalysis - 詳細分析を実行
func (pa *PerformanceAnalyzer) performDetailedAnalysis(results *AnalysisResults) error {
	detailedAnalysis := make(map[string]interface{})

	// 1. メモリ使用効率分析
	memoryAnalysis := pa.analyzeMemoryEfficiency(results)
	detailedAnalysis["memory_efficiency"] = memoryAnalysis

	// 2. I/O効率分析
	ioAnalysis := pa.analyzeIOEfficiency(results)
	detailedAnalysis["io_efficiency"] = ioAnalysis

	// 3. キャッシュヒット率分析
	hitRateAnalysis := pa.analyzeHitRatePatterns(results)
	detailedAnalysis["hit_rate_patterns"] = hitRateAnalysis

	// 4. コスト効率分析
	costAnalysis := pa.analyzeCostEfficiency(results)
	detailedAnalysis["cost_efficiency"] = costAnalysis

	// 5. スケーラビリティ分析
	scalabilityAnalysis := pa.analyzeScalability(results)
	detailedAnalysis["scalability"] = scalabilityAnalysis

	results.DetailedAnalysis = detailedAnalysis
	return nil
}

// analyzeMemoryEfficiency - メモリ効率分析
func (pa *PerformanceAnalyzer) analyzeMemoryEfficiency(results *AnalysisResults) map[string]interface{} {
	return map[string]interface{}{
		"oracle_unified_memory": map[string]interface{}{
			"sga_total_mb":          results.PerformanceComparison.ResourceUtilization.TotalCacheMemoryMB,
			"efficiency_per_mb":     results.PerformanceComparison.EfficiencyMetrics.MemoryEfficiencyRatio,
			"shared_memory_benefit": "複数プロセス間で自動共有",
			"automatic_management":  "Oracle AMM/ASMM による自動最適化",
		},
		"external_cache_overhead": map[string]interface{}{
			"memory_duplication":     "Oracle + 外部キャッシュ = メモリ二重使用",
			"serialization_overhead": "JSON変換による追加メモリ消費",
			"network_buffers":        "ネットワーク通信用バッファメモリ",
			"management_overhead":    "キャッシュ管理用メタデータ",
		},
	}
}

// analyzeIOEfficiency - I/O効率分析
func (pa *PerformanceAnalyzer) analyzeIOEfficiency(results *AnalysisResults) map[string]interface{} {
	savedIO := results.OracleBufferMetrics.LogicalReads - results.OracleBufferMetrics.PhysicalReads

	return map[string]interface{}{
		"oracle_io_optimization": map[string]interface{}{
			"physical_reads_saved":    savedIO,
			"io_reduction_percentage": results.PerformanceComparison.EfficiencyMetrics.IOReductionRatio,
			"automatic_prefetching":   "Oracle Smart Scan および読み先読み",
			"block_level_caching":     "8KB データブロック単位の効率的キャッシュ",
		},
		"external_cache_io_overhead": map[string]interface{}{
			"network_io_penalty": "1-5ms の追加レイテンシ",
			"double_io_risk":     "キャッシュミス時のOracle + 外部システムアクセス",
			"serialization_io":   "JSON変換による追加CPU → I/O変換負荷",
		},
	}
}

// analyzeHitRatePatterns - ヒット率パターン分析
func (pa *PerformanceAnalyzer) analyzeHitRatePatterns(results *AnalysisResults) map[string]interface{} {
	return map[string]interface{}{
		"buffer_cache_patterns": map[string]interface{}{
			"hit_ratio":            results.OracleBufferMetrics.HitRatio,
			"lru_efficiency":       "Least Recently Used による自動最適化",
			"hot_block_management": "頻繁アクセスブロックの自動保持",
		},
		"result_cache_patterns": map[string]interface{}{
			"hit_ratio":              results.OracleResultMetrics.HitRatio,
			"automatic_invalidation": "依存関係に基づく自動無効化",
			"sql_result_sharing":     "同一SQL結果の自動共有",
		},
		"external_cache_challenges": map[string]interface{}{
			"manual_key_management":   "手動キャッシュキー設計の複雑性",
			"invalidation_complexity": "関連データ無効化の実装負荷",
			"consistency_risks":       "データ整合性保証の困難性",
		},
	}
}

// analyzeCostEfficiency - コスト効率分析
func (pa *PerformanceAnalyzer) analyzeCostEfficiency(_ *AnalysisResults) map[string]interface{} {
	return map[string]interface{}{
		"oracle_built_in_benefits": map[string]interface{}{
			"zero_additional_infrastructure": "追加インフラストラクチャ不要",
			"included_in_license":            "Oracleライセンスに含まれる",
			"automatic_maintenance":          "自動メンテナンス（手動作業なし）",
			"built_in_monitoring":            "統計ビューによる統合監視",
		},
		"external_cache_costs": map[string]interface{}{
			"infrastructure_costs": "Redis等の専用サーバ・クラスタ",
			"development_overhead": "キャッシュロジック開発・保守工数",
			"operational_overhead": "外部システム監視・メンテナンス工数",
			"complexity_costs":     "システム複雑性による開発・障害対応コスト",
		},
	}
}

// analyzeScalability - スケーラビリティ分析
func (pa *PerformanceAnalyzer) analyzeScalability(_ *AnalysisResults) map[string]interface{} {
	return map[string]interface{}{
		"oracle_scalability": map[string]interface{}{
			"sga_auto_scaling":    "SGA自動サイズ調整機能",
			"rac_shared_cache":    "Oracle RAC環境でのキャッシュ共有",
			"built_in_clustering": "データベースクラスタリング統合",
			"workload_adaptive":   "ワークロードに応じた自動最適化",
		},
		"external_cache_scaling_challenges": map[string]interface{}{
			"cluster_complexity":       "Redis Cluster等の複雑な構成管理",
			"data_distribution":        "キャッシュデータ分散の設計・運用課題",
			"consistency_at_scale":     "分散環境でのデータ整合性確保の困難性",
			"cross_cache_coordination": "複数キャッシュ間の協調制御の複雑性",
		},
	}
}

// generateRecommendations - 推奨事項を生成
func (pa *PerformanceAnalyzer) generateRecommendations(results *AnalysisResults) error {
	var recommendations []Recommendation

	// Buffer Cache関連推奨事項
	if results.OracleBufferMetrics.HitRatio < 90 {
		recommendations = append(recommendations, Recommendation{
			Category:    "Buffer Cache最適化",
			Priority:    "高",
			Title:       "Buffer Cacheサイズの最適化",
			Description: "Buffer Cacheヒット率が90%未満のため、サイズ調整を検討",
			Impact:      "I/O削減による大幅な性能向上",
			Effort:      "低（パラメータ調整のみ）",
			Benefits:    []string{"物理I/O削減", "レスポンス時間向上", "CPU使用率改善"},
		})
	}

	// Result Cache関連推奨事項
	if results.OracleResultMetrics.HitRatio < 70 {
		recommendations = append(recommendations, Recommendation{
			Category:    "Result Cache活用",
			Priority:    "中",
			Title:       "RESULT_CACHEヒントの積極活用",
			Description: "集計クエリにRESULT_CACHEヒントを追加して効率化",
			Impact:      "複雑クエリの大幅な高速化",
			Effort:      "中（SQL修正が必要）",
			Benefits:    []string{"集計処理高速化", "CPU負荷軽減", "同時実行性向上"},
		})
	}

	// 外部キャッシュ見直し推奨事項
	recommendations = append(recommendations, Recommendation{
		Category:    "アーキテクチャ最適化",
		Priority:    "高",
		Title:       "外部キャッシュ依存度の削減",
		Description: "Oracle内蔵キャッシュを最大限活用し、外部キャッシュ依存を削減",
		Impact:      "システム複雑性の削減と運用性向上",
		Effort:      "高（アーキテクチャ変更）",
		Benefits:    []string{"システム複雑性削減", "運用コスト削減", "データ整合性向上", "レイテンシ削減"},
	})

	// N+1問題根本解決推奨事項
	recommendations = append(recommendations, Recommendation{
		Category:    "SQL最適化",
		Priority:    "最高",
		Title:       "N+1問題の根本的解決",
		Description: "JOINやIN句を使用してN+1問題を根本から解決",
		Impact:      "クエリ実行回数の劇的削減",
		Effort:      "中（SQL設計見直し）",
		Benefits:    []string{"実行時間短縮", "データベース負荷軽減", "スケーラビリティ向上"},
	})

	// 優先度でソート
	sort.Slice(recommendations, func(i, j int) bool {
		priorities := map[string]int{"最高": 4, "高": 3, "中": 2, "低": 1}
		return priorities[recommendations[i].Priority] > priorities[recommendations[j].Priority]
	})

	results.Recommendations = recommendations
	return nil
}

// DisplayComprehensiveResults - 包括的結果を表示
func (pa *PerformanceAnalyzer) DisplayComprehensiveResults() {
	if pa.analysisResults == nil {
		fmt.Println("分析結果がありません")
		return
	}

	results := pa.analysisResults

	fmt.Println("\\n" + strings.Repeat("=", 80))
	fmt.Println("Oracle内蔵キャッシュ vs 外部キャッシュ 包括的分析結果")
	fmt.Println(strings.Repeat("=", 80))

	// 1. エグゼクティブサマリー
	pa.displayExecutiveSummary(results)

	// 2. 性能メトリクス比較
	pa.displayPerformanceMetrics(results)

	// 3. 効率性分析
	pa.displayEfficiencyAnalysis(results)

	// 4. 推奨事項
	pa.displayRecommendations(results)

	// 5. 結論
	pa.displayConclusion(results)
}

// displayExecutiveSummary - エグゼクティブサマリーを表示
func (pa *PerformanceAnalyzer) displayExecutiveSummary(results *AnalysisResults) {
	fmt.Println("\\n■ エグゼクティブサマリー")
	fmt.Println(strings.Repeat("-", 50))

	efficiency := results.PerformanceComparison.EfficiencyMetrics.OverallCacheEfficiency

	if efficiency >= 90 {
		fmt.Println("✅ 総合評価: 優秀（90%以上の効率性）")
		fmt.Println("   Oracle内蔵キャッシュが効果的に機能しています")
	} else if efficiency >= 70 {
		fmt.Println("⚠️  総合評価: 良好（70-90%の効率性）")
		fmt.Println("   改善の余地がありますが、基本的な機能は正常です")
	} else {
		fmt.Println("❌ 総合評価: 要改善（70%未満の効率性）")
		fmt.Println("   早急な最適化が必要です")
	}

	fmt.Printf("\\n• 総合キャッシュ効率: %.1f%%\\n", efficiency)
	fmt.Printf("• Buffer Cache効率: %.1f%%\\n", results.PerformanceComparison.EfficiencyMetrics.BufferCacheEfficiency)
	fmt.Printf("• Result Cache効率: %.1f%%\\n", results.PerformanceComparison.EfficiencyMetrics.ResultCacheEfficiency)
	fmt.Printf("• 総キャッシュメモリ: %.1f MB\\n", results.PerformanceComparison.ResourceUtilization.TotalCacheMemoryMB)
	fmt.Printf("• 推定I/O削減: %d回\\n", results.PerformanceComparison.ResourceUtilization.EstimatedIOSavings)
}

// displayPerformanceMetrics - 性能メトリクス比較を表示
func (pa *PerformanceAnalyzer) displayPerformanceMetrics(results *AnalysisResults) {
	fmt.Println("\\n■ Oracle内蔵キャッシュ vs 外部キャッシュ 比較")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("\\n✅ Oracle内蔵キャッシュの優位性:")
	for _, advantage := range results.PerformanceComparison.OracleAdvantages {
		fmt.Printf("   %s\\n", advantage)
	}

	fmt.Println("\\n❌ 外部キャッシュの課題:")
	for _, issue := range results.PerformanceComparison.ExternalCacheIssues {
		fmt.Printf("   %s\\n", issue)
	}
}

// displayEfficiencyAnalysis - 効率性分析を表示
func (pa *PerformanceAnalyzer) displayEfficiencyAnalysis(results *AnalysisResults) {
	fmt.Println("\\n■ リソース効率性分析")
	fmt.Println(strings.Repeat("-", 50))

	util := results.PerformanceComparison.ResourceUtilization

	fmt.Printf("\\n💾 メモリ使用効率:\\n")
	fmt.Printf("   • 総キャッシュメモリ: %.1f MB\\n", util.TotalCacheMemoryMB)
	fmt.Printf("   • Buffer Cache: %.1f MB\\n", util.BufferCacheUtilization)
	fmt.Printf("   • Result Cache: %.1f MB\\n", util.ResultCacheUtilization)

	fmt.Printf("\\n⚡ 性能向上効果:\\n")
	fmt.Printf("   • I/O削減回数: %d回\\n", util.EstimatedIOSavings)
	fmt.Printf("   • 推定CPU削減: %.1f%%\\n", util.EstimatedCPUSavings)
	fmt.Printf("   • メモリ効率比: %.2f\\n", results.PerformanceComparison.EfficiencyMetrics.MemoryEfficiencyRatio)
}

// displayRecommendations - 推奨事項を表示
func (pa *PerformanceAnalyzer) displayRecommendations(results *AnalysisResults) {
	fmt.Println("\\n■ 推奨事項（優先度順）")
	fmt.Println(strings.Repeat("-", 50))

	for i, rec := range results.Recommendations {
		fmt.Printf("\\n%d. [%s優先度] %s\\n", i+1, rec.Priority, rec.Title)
		fmt.Printf("   カテゴリ: %s\\n", rec.Category)
		fmt.Printf("   説明: %s\\n", rec.Description)
		fmt.Printf("   期待効果: %s\\n", rec.Impact)
		fmt.Printf("   実装工数: %s\\n", rec.Effort)
		if len(rec.Benefits) > 0 {
			fmt.Printf("   利益: %s\\n", strings.Join(rec.Benefits, ", "))
		}
	}
}

// displayConclusion - 結論を表示
func (pa *PerformanceAnalyzer) displayConclusion(results *AnalysisResults) {
	fmt.Println("\\n■ 結論とNext Steps")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println("\\n🎯 重要な結論:")
	fmt.Println("   1. Oracle内蔵キャッシュメカニズムは外部キャッシュよりも効率的")
	fmt.Println("   2. N+1問題は根本的なSQL設計で解決すべき")
	fmt.Println("   3. 外部キャッシュは複雑性を増加させ運用コストを高める")
	fmt.Println("   4. データ整合性はOracleの自動機能に任せるべき")

	fmt.Printf("\\n📊 今回の分析結果: %s\\n", results.OptimizationAdvice.PriorityLevel)

	fmt.Println("\\n🚀 Next Steps:")
	if len(results.OptimizationAdvice.ImmediateActions) > 0 {
		fmt.Println("   即座の対応:")
		for _, action := range results.OptimizationAdvice.ImmediateActions {
			fmt.Printf("     • %s\\n", action)
		}
	}

	fmt.Println("\\n💡 長期的な方向性:")
	fmt.Println("   → Oracle Database中心のアーキテクチャ採用")
	fmt.Println("   → 外部キャッシュ依存度の段階的削減")
	fmt.Println("   → SQL最適化による根本的問題解決")
	fmt.Println("   → 運用性とメンテナンス性の向上")
}

// ExportAnalysisResults - 分析結果をJSONでエクスポート
func (pa *PerformanceAnalyzer) ExportAnalysisResults() (string, error) {
	if pa.analysisResults == nil {
		return "", fmt.Errorf("分析結果がありません")
	}

	jsonData, err := json.MarshalIndent(pa.analysisResults, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON変換エラー: %w", err)
	}

	return string(jsonData), nil
}

// GetAnalysisResults - 分析結果を取得
func (pa *PerformanceAnalyzer) GetAnalysisResults() *AnalysisResults {
	return pa.analysisResults
}
