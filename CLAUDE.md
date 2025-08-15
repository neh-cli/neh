# CLAUDE.md

このファイルは、このリポジトリでコードを扱う際のClaude Code (claude.ai/code) へのガイダンスを提供します。

## 概要

NehはGoで書かれたCLIアプリケーションで、WebSocket接続を通じて大規模言語モデルとのリアルタイム対話を提供します。コマンドラインインターフェースにはCobraフレームワークを使用し、AI処理のためにバックエンドサーバーと通信します。

## 開発コマンド

### ビルドとテスト
```bash
# テストを実行してバイナリをビルド
make build

# テストのみ実行
make test
# または
go test -v ./...

# テストなしでビルド
go build -o build/bin/neh

# ビルド成果物をクリーン
make clean

# goモジュールを整理
make tidy
```

### 個別テストの実行
```bash
# 特定パッケージのテストを実行
go test -v ./cmd/
go test -v ./cmd/shared/

# 特定のテスト関数を実行
go test -v -run TestFunctionName ./cmd/
```

## アーキテクチャ

### コアコンポーネント

1. **コマンド構造**: `cmd/root.go`にルートコマンドがあり、個別のサブコマンドは`cmd/`ディレクトリ下の別々のファイルに配置されたCobraフレームワークを使用。

2. **WebSocket通信**: すべてのAI対話は`cmd/shared/utils.go`で管理されるWebSocket接続を通じて行われます。フローは以下の通り：
   - サーバーへのWebSocket接続を確立
   - ActionCableチャンネル（LargeLanguageModelQueryChannel）にサブスクライブ
   - クエリとUUIDを含むHTTPリクエストを送信
   - シーケンス番号で順序付けされたストリーミングレスポンスをWebSocket経由で受信

3. **サーバーエンドポイント**: 
   - 本番環境: `wss://yoryo-app.onrender.com/cable`
   - 開発環境: `NEH_WORKING_ON_LOCALHOST`と`NEH_SERVER_ENDPOINT_DEVELOPMENT`環境変数で設定

### 主要な環境変数

- `NEH_PERSONAL_ACCESS_TOKEN`: APIアクセスに必要な認証トークン
- `NEH_DEBUG`: "t"に設定してデバッグ出力を有効化
- `NEH_WORKING_ON_LOCALHOST`: ローカル開発エンドポイントに切り替え
- `NEH_SERVER_ENDPOINT_DEVELOPMENT`: カスタム開発サーバーエンドポイント

### 設定

アプリは言語設定（`lang`フィールド）などの設定を`~/.config/neh/config.yml`から読み込みます。

### コマンドタイプ

- **クエリコマンド** (`o`, `c`): オプションのクリップボード内容と共にユーザークエリを送信
- **処理コマンド** (`explain`, `fix`, `refactor`, `refine`): クリップボード内容に特定のAI操作を適用
- **翻訳** (`t`): 設定された言語に基づいてクリップボード内容を翻訳
- **ユーティリティコマンド** (`clip`, `todo`, `decache`, `status`, `version`, `config`)

### メッセージフロー

1. コマンドが認証ヘッダー付きでWebSocket接続を初期化
2. 一意のUUIDでActionCableチャンネルにサブスクライブ
3. メッセージペイロードを含むHTTP POSTリクエストを送信
4. シーケンス番号で順序付けされたストリーミングレスポンスを受信
5. 異なるメッセージタイプを処理: output、error、worker_done

### テストアプローチ

テストはソースファイルと同じ場所に配置されます（例：`o_test.go`、`decache_test.go`）。プロジェクトは標準のGoテストフレームワークを使用しています。