# Makefile for PDF生成CLIツール
# このMakefileは、PDF生成ツールのビルドと実行を簡単にするためのものです

# 変数の定義
BINARY_NAME=pdf-generator
MAIN_FILE=main.go

# デフォルトターゲット: アプリケーションを実行する
.PHONY: run
run:
	@echo "🚀 PDF生成ツールを実行しています..."
	go run $(MAIN_FILE)

# ビルドターゲット: 実行可能ファイルをビルドする
.PHONY: build
build:
	@echo "🔨 アプリケーションをビルドしています..."
	go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ ビルド完了: $(BINARY_NAME)"

# 依存関係のインストール
.PHONY: deps
deps:
	@echo "📦 依存関係をインストールしています..."
	go mod download
	go mod tidy
	@echo "✅ 依存関係のインストール完了"

# クリーンアップ: 生成されたファイルを削除
.PHONY: clean
clean:
	@echo "🧹 クリーンアップ中..."
	rm -f $(BINARY_NAME)
	rm -rf PDF/*.pdf
	@echo "✅ クリーンアップ完了"

# テスト実行
.PHONY: test
test:
	@echo "🧪 テストを実行しています..."
	go test ./...

# ヘルプ表示
.PHONY: help
help:
	@echo "利用可能なコマンド:"
	@echo "  make run     - アプリケーションを実行する"
	@echo "  make build   - 実行可能ファイルをビルドする"
	@echo "  make deps    - 依存関係をインストールする"
	@echo "  make clean   - 生成されたファイルを削除する"
	@echo "  make test    - テストを実行する"
	@echo "  make help    - このヘルプを表示する"

