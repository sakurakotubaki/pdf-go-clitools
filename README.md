# GoでPDF生成するCLIツール

テキストファイルからPDFを生成するCLIツールです。日本語フォントに対応しています。

## 機能

- テキストファイル（cli.txt）を読み込んでPDFに変換
- 日本語フォント対応（UTF-8エンコーディング）
- PDFディレクトリへの自動出力

## アプリケーションの作成手順

### 1. 必要な環境

- Go 1.21以上がインストールされていること
- macOS/Linux環境（Windowsでも動作しますが、フォント設定が異なる場合があります）

### 2. 依存関係のインストール

```bash
# 依存関係をダウンロードして整理
go mod download
go mod tidy

# またはMakefileを使用
make deps
```

### 3. 実行方法

#### 方法1: Makefileを使用（推奨）

```bash
# アプリケーションを実行
make run

# その他のコマンド
make build   # 実行可能ファイルをビルド
make clean   # 生成されたファイルを削除
make help    # ヘルプを表示
```

#### 方法2: go runコマンドを使用

```bash
go run main.go
```

### 4. 出力

- 入力ファイル: `cli.txt`
- 出力先: `PDF/output.pdf`
- PDFディレクトリが存在しない場合は自動的に作成されます

## コードの説明

### main.goの構造

1. **main関数**: アプリケーションのエントリーポイント
   - cli.txtファイルの読み込み
   - PDFディレクトリの作成
   - PDFファイルの生成

2. **readTextFile関数**: テキストファイルを読み込む
   - UTF-8エンコーディングで日本語も正しく読み込める

3. **ensureDirectory関数**: ディレクトリの存在確認と作成
   - PDFディレクトリが存在しない場合は自動的に作成

4. **generatePDF関数**: PDFファイルの生成
   - gopdfライブラリを使用（日本語フォントに優れた対応）
   - 日本語フォント（NotoSansJP）を設定
   - A4サイズのPDFを生成
   - マージンと行間を適切に設定

### 使用しているパッケージ

- `github.com/signintech/gopdf`: PDF生成ライブラリ
  - 日本語フォントに優れた対応
  - Noto Sans JPフォントを正しく表示可能

## 日本語フォントの設定

gopdfで日本語を正しく表示するには、TTF形式の日本語フォントファイルが必要です。

### フォントファイルの取得方法

1. **Google Noto Fontsを使用する（推奨）**
   ```bash
   # Noto Sans CJK（日本語対応）をダウンロード
   # https://fonts.google.com/noto/specimen/Noto+Sans+JP からダウンロード可能
   ```

2. **フォントファイルの配置**
   - プロジェクト内の `fonts/` または `font/` ディレクトリに配置
   - または、macOSの場合は `~/Library/Fonts/` に配置

3. **対応しているフォントファイル名**
   - `NotoSansCJK-Regular.ttf`（推奨）
   - `NotoSansJP-Regular.ttf`（推奨）
   - `NotoSansCJK.ttf`
   - `NotoSansJP.ttf`
   - `NotoSans-Regular.ttf`
   - `ZenOldMincho-Regular.ttf`（検出されますが、互換性を確認してください）

### フォントの互換性について

**重要**: gopdfは日本語フォントに優れた対応をしています。

- ✅ **推奨**: Noto Sans JPフォント（Google Noto Fonts）
  - `NotoSansJP-Regular.ttf`が推奨されます
  - gopdfライブラリは、Noto Sans JPフォントを正しく表示できます
  - 大きなCIDマップを持つフォントにも対応しています

### フォントファイルが見つからない場合

フォントファイルが見つからない場合、アプリケーションは警告を表示してデフォルトフォント（helvetica）を使用します。この場合、日本語は正しく表示されない可能性があります。

### フォントの追加に失敗した場合

フォントファイルが検出されても、フォントの追加に失敗した場合、アプリケーションは警告を表示してデフォルトフォントを使用します。日本語を正しく表示するには、TTF形式の日本語フォントファイルが必要です。

## 注意事項

- cli.txtファイルはUTF-8エンコーディングで保存してください
- macOSのシステムフォント（TTC形式）は直接使用できません。TTF形式のフォントファイルが必要です