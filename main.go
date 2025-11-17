package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/signintech/gopdf"
)

// main関数: アプリケーションのエントリーポイント
// 1. cli.txtファイルを読み込む
// 2. テキスト内容をPDFに変換する
// 3. PDFディレクトリに出力する
func main() {
	// 入力ファイル名を定義
	inputFile := "cli.txt"
	outputDir := "PDF"

	// ステップ1: cli.txtファイルの存在確認と読み込み
	fmt.Printf("📄 ファイル '%s' を読み込んでいます...\n", inputFile)
	content, err := readTextFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ エラー: %v\n", err)
		os.Exit(1)
	}

	// ステップ2: PDFディレクトリの作成（存在しない場合）
	if err := ensureDirectory(outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "❌ エラー: %v\n", err)
		os.Exit(1)
	}

	// ステップ3: PDFファイルの生成
	outputFile := filepath.Join(outputDir, "output.pdf")
	fmt.Printf("📝 PDFファイル '%s' を生成しています...\n", outputFile)
	if err := generatePDF(content, outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "❌ エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ PDFファイルが正常に生成されました: %s\n", outputFile)
}

// readTextFile: テキストファイルを読み込む関数
// 引数: filename - 読み込むファイルのパス
// 戻り値: ファイルの内容とエラー
func readTextFile(filename string) (string, error) {
	// os.ReadFileを使用してファイル全体を一度に読み込む
	// UTF-8エンコーディングで日本語も正しく読み込める
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("ファイル '%s' の読み込みに失敗しました: %w", filename, err)
	}
	return string(data), nil
}

// ensureDirectory: ディレクトリが存在するか確認し、存在しない場合は作成する
// 引数: dirname - 確認/作成するディレクトリのパス
// 戻り値: エラー（存在するか作成成功時はnil）
func ensureDirectory(dirname string) error {
	// os.Statでディレクトリの存在確認
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		// ディレクトリが存在しない場合は作成
		// 0755は読み取り・実行・書き込み権限を設定
		if err := os.MkdirAll(dirname, 0755); err != nil {
			return fmt.Errorf("ディレクトリ '%s' の作成に失敗しました: %w", dirname, err)
		}
		fmt.Printf("📁 ディレクトリ '%s' を作成しました\n", dirname)
	}
	return nil
}

// generatePDF: テキスト内容からPDFファイルを生成する関数
// 引数: content - PDFに含めるテキスト内容
//      outputPath - 出力するPDFファイルのパス
// 戻り値: エラー
func generatePDF(content string, outputPath string) error {
	// gopdfライブラリを使用してPDFオブジェクトを作成
	// gopdfは日本語フォントにより良い対応をしています
	pdf := gopdf.GoPdf{}

	// A4サイズのページを設定（210mm x 297mm）
	// gopdfの単位はポイント（pt）のため、mmからptへ変換
	pageWidth := mmToPt(210.0)
	pageHeight := mmToPt(297.0)
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: pageWidth, H: pageHeight}})

	// 日本語フォントの設定
	// gopdfで日本語を表示するには、TTF形式のフォントファイルが必要です
	fontPath := findJapaneseFont()
	fontName := "JapaneseFont"
	fontAdded := false

	if fontPath != "" {
		// フォントファイルが存在する場合、TTFフォントとして追加
		// AddTTFFontの引数: フォントファイルのパス, フォント名
		fmt.Printf("📝 フォントファイル '%s' を使用します\n", fontPath)
		
		// フォントを追加
		err := pdf.AddTTFFont(fontName, fontPath)
		if err != nil {
			fmt.Printf("⚠️  警告: フォント '%s' の追加に失敗しました: %v\n", fontPath, err)
			fmt.Println("   デフォルトフォントを使用します。")
		} else {
			// フォントを設定
			err = pdf.SetFont(fontName, "", 12)
			if err != nil {
				fmt.Printf("⚠️  警告: フォント '%s' の設定に失敗しました: %v\n", fontName, err)
			} else {
				fmt.Printf("✅ フォント '%s' を正常に追加しました\n", filepath.Base(fontPath))
				fontAdded = true
			}
		}
	}

	// フォントが追加されなかった場合、デフォルトフォントを使用
	if !fontAdded {
		if fontPath == "" {
			fmt.Println("⚠️  警告: 日本語フォントファイルが見つかりません。デフォルトフォントを使用します。")
			fmt.Println("   日本語を正しく表示するには、TTF形式の日本語フォントファイルが必要です。")
		}
		// gopdfのデフォルトフォントを設定
		err := pdf.SetFont("helvetica", "", 12)
		if err != nil {
			return fmt.Errorf("デフォルトフォントの設定に失敗しました: %w", err)
		}
	}

	// 最初のページを追加
	pdf.AddPage()

	// テキストを改行で分割して処理
	lines := strings.Split(content, "\n")

	// ページのマージンを設定（上下左右20mm）
	margin := mmToPt(20.0)
	bottomMargin := mmToPt(20.0)

	// フォントサイズに基づいて行の高さを計算（6mm）
	lineHeight := mmToPt(6.0)
	blankLineGap := mmToPt(5.0)

	// 現在のY座標（縦位置）を設定
	y := margin

	// 各行をPDFに追加
	for _, line := range lines {
		// 空行の場合は少しスペースを追加
		if strings.TrimSpace(line) == "" {
			y += blankLineGap
			continue
		}

		// ページの下端に近づいた場合は新しいページを追加
		// マージンと余白（20mm）を考慮して改ページを判定
		if y > pageHeight-bottomMargin-lineHeight {
			pdf.AddPage()
			y = margin
		}

		// テキストを出力
		// SetX, SetYで位置を設定し、Textメソッドで直接テキストを出力
		// Textメソッドは位置を自動更新しないため、位置管理が簡単で正確
		pdf.SetXY(margin, y)

		// Textメソッドでテキストを出力
		// 引数: テキスト内容
		// SetX, SetYで設定した位置にテキストを出力
		err := pdf.Text(line)
		if err != nil {
			return fmt.Errorf("テキストの出力に失敗しました: %w", err)
		}

		// 次の行の位置を計算（行間: 6mm）
		y += lineHeight
	}

	// PDFファイルを出力
	// WritePdfメソッドでファイルに保存
	if err := pdf.WritePdf(outputPath); err != nil {
		return fmt.Errorf("PDFファイルの生成に失敗しました: %w", err)
	}

	return nil
}

// findJapaneseFont: 日本語フォントファイルを検索する関数
// 戻り値: フォントファイルのパス（見つからない場合は空文字列）
func findJapaneseFont() string {
	// 一般的な日本語フォントファイル名のリスト
	// 優先順位: NotoSansJP-Regular > NotoSansJP-VariableFont > その他のNotoSansJP > ZenOldMincho
	fontNames := []string{
		// Noto Sans JP Regular（最優先 - 推奨）
		"static/NotoSansJP-Regular.ttf",
		"NotoSansJP-Regular.ttf",
		// Noto Sans JP Variable Font（可変フォント）
		"NotoSansJP-VariableFont_wght.ttf",
		// Noto Sans JP その他のスタイル
		"static/NotoSansJP-Medium.ttf",
		"static/NotoSansJP-Light.ttf",
		"static/NotoSansJP-Bold.ttf",
		"static/NotoSansJP-SemiBold.ttf",
		"static/NotoSansJP-ExtraLight.ttf",
		"static/NotoSansJP-ExtraBold.ttf",
		"static/NotoSansJP-Thin.ttf",
		"static/NotoSansJP-Black.ttf",
		"NotoSansJP-Medium.ttf",
		"NotoSansJP-Light.ttf",
		"NotoSansJP-Bold.ttf",
		"NotoSansJP-SemiBold.ttf",
		// Noto Sans CJK
		"NotoSansCJK-Regular.ttf",
		"NotoSansCJK.ttf",
		"NotoSansJP.ttf",
		"NotoSans-Regular.ttf",
		// ZenOldMincho
		"ZenOldMincho-Regular.ttf",
		"ZenOldMincho-Medium.ttf",
		"ZenOldMincho-SemiBold.ttf",
		"ZenOldMincho-Bold.ttf",
		"ZenOldMincho-Black.ttf",
	}

	// 検索するディレクトリのリスト
	searchDirs := []string{
		"./font",                     // プロジェクト内のfontディレクトリ（最優先）
		"./fonts",                    // プロジェクト内のfontsディレクトリ
		os.Getenv("HOME") + "/Library/Fonts", // macOSのユーザーフォントディレクトリ
		"/Library/Fonts",             // macOSのシステムフォントディレクトリ
		"/System/Library/Fonts",      // macOSのシステムフォントディレクトリ
	}

	// 各ディレクトリとフォント名の組み合わせを確認
	for _, dir := range searchDirs {
		for _, fontName := range fontNames {
			fontPath := filepath.Join(dir, fontName)
			// ファイルが存在するか確認
			if _, err := os.Stat(fontPath); err == nil {
				return fontPath
			}
		}
	}

	// フォントファイルが見つからない場合は空文字列を返す
	return ""
}

// mmToPt: ミリメートルをポイント（pt）に変換するヘルパー関数
func mmToPt(mm float64) float64 {
	return mm * 72.0 / 25.4
}
