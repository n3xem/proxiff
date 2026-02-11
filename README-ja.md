# Proxiff

Goで書かれたDiffy風のHTTPプロキシツールで、2つの異なるサーバー（newer と current）からのレスポンスを比較します。

[English README](README.md)

## 機能

- HTTPリクエストを2つの異なるサーバー（newer と current）に転送
- レスポンスを比較して差分をログに記録
- クライアントにはcurrentサーバーのレスポンスを返却

## インストール

```bash
go build -o proxiff ./cmd/proxiff
```

## 使い方

基本的な使い方：

```bash
./proxiff -newer http://localhost:8082 -current http://localhost:8081 -port 8080
```

オプション：
- `-newer`: newerサーバーのURL（必須）
- `-current`: currentサーバーのURL（必須）
- `-port`: リスニングポート（デフォルト: 8080）

## 使用例

1. 2つのサンプルサーバーを起動：

```bash
# ターミナル1: Currentサーバー
go run ./example/servers -port 8081 -version current

# ターミナル2: Newerサーバー
go run ./example/servers -port 8082 -version newer
```

2. proxiffを起動：

```bash
# ターミナル3: Proxiffプロキシ
go run ./cmd/proxiff -newer http://localhost:8082 -current http://localhost:8081 -port 8080
```

3. プロキシにリクエストを送信：

```bash
# currentサーバーのレスポンスが返されますが
# newerサーバーとの差分がログに記録されます
curl http://localhost:8080/
curl http://localhost:8080/api/users
curl http://localhost:8080/api/status
```

## 差分ログの出力例

Proxiffは[google/go-cmp](https://github.com/google/go-cmp)ライブラリを使用して、2つのサーバーのレスポンスを比較します。差分が検出された場合、わかりやすい形式でログに出力されます。

### フィールド追加の差分（/api/users）

newerバージョンで`email`フィールドが追加された場合：

```
2025/11/23 15:07:51 Difference detected:   &comparator.Response{
  	StatusCode: 200,
  	Headers: http.Header{
- 		"Content-Length": {"78"},
+ 		"Content-Length": {"130"},
  		"Content-Type":   {"application/json"},
  		"Date":           {"Sun, 23 Nov 2025 15:07:51 GMT"},
  	},
  	Body: bytes.Join({
  		`{"users":[{"`,
- 		`id":1,"name":"Alice"},{`,
+ 		`email":"alice@example.com","id":1,"name":"Alice"},{"email":"bob@`,
+ 		`example.com",`,
  		`"id":2,"name":"Bob"}],"version":"`,
- 		"current",
+ 		"newer",
  		"\"}\n",
  	}, ""),
  }
```

### ステータスコードの差分（/api/status）

currentサーバーがHTTP 200を返し、newerサーバーがHTTP 201を返す場合：

```
2025/11/23 15:07:53 Difference detected:   &comparator.Response{
- 	StatusCode: 200,
+ 	StatusCode: 201,
  	Headers: http.Header{
- 		"Content-Length": {"36"},
+ 		"Content-Length": {"34"},
  		"Content-Type":   {"application/json"},
  		"Date":           {"Sun, 23 Nov 2025 15:07:53 GMT"},
  	},
  	Body: bytes.Join({
  		`{"status":"ok","version":"`,
- 		"current",
+ 		"newer",
  		"\"}\n",
  	}, ""),
  }
```

### 差分がない場合

レスポンスが完全に一致する場合：

```
2025/11/23 15:04:00 Responses match
```

### ログの記号の意味

- `-` 記号：currentサーバーの内容（削除された部分）
- `+` 記号：newerサーバーの内容（追加された部分）

## アーキテクチャ

Proxiffは`SimpleComparator`（[google/go-cmp](https://github.com/google/go-cmp)ベース）を使用してレスポンスを比較します。ステータスコード、ヘッダー、レスポンスボディを比較対象とします。比較ロジックは`Comparator`インターフェースの背後にあるため、必要に応じてカスタム実装に差し替え可能です。

## テスト

全てのテストを実行：

```bash
go test ./... -v
```

特定のパッケージのテストを実行：

```bash
go test ./comparator/... -v
go test ./proxy/... -v
```

## プロジェクト構成

```
proxiff/
├── cmd/
│   └── proxiff/        # メインCLIアプリケーション
│       └── main.go
├── comparator/         # 比較ロジック
│   ├── comparator.go   # インターフェース定義
│   ├── simple.go       # SimpleComparator実装
│   └── simple_test.go  # テスト
├── proxy/              # プロキシのコア機能
│   ├── proxy.go
│   └── proxy_test.go
└── example/
    ├── deployment/     # サンプルデプロイメント構成
    │   ├── docker/     # Docker Composeを使ったNginx連携例
    │   └── nginx/      # Nginx設定サンプル
    └── servers/        # テスト用サンプルサーバー
        └── main.go
```

## ユースケース

1. **カナリアデプロイメント**: 本番環境とカナリア版を比較
2. **マイグレーションテスト**: 新実装が旧実装と同じ動作をすることを確認
3. **A/Bテスト**: 異なるアルゴリズム実装を比較
4. **リグレッションテスト**: APIレスポンスの予期しない変更を検出

## デプロイメント例（Nginx Mirrorモジュールを使用）

`example/deployment/` ディレクトリには、Nginxの[mirrorモジュール](https://nginx.org/en/docs/http/ngx_http_mirror_module.html)と組み合わせて本番トラフィックに影響を与えずに新旧バージョンの比較を行う方法のサンプル構成が含まれています。

### アーキテクチャ

```
クライアント
  ↓
Nginx
  ├─> 本番サーバー（レスポンスをクライアントに返す）
  └─> Proxiff（mirror、レスポンスは無視）
        ├─> Newer Server
        └─> Current Server
              ↓
        差分検出とログ出力
```

### Docker Composeサンプル環境

`example/deployment/docker/` ディレクトリには、すぐに使えるサンプル環境が用意されています：

```bash
cd example/deployment/docker

# サンプル環境を起動
docker compose up -d

# ログを確認
docker compose logs -f proxiff

# テストリクエストを送信
curl http://localhost:8000/api/users

# クリーンアップ
docker compose down -v
```

### 統合テストの実行

```bash
cd example/deployment/docker
./test-integration.sh
```

セットアップのカスタマイズ方法などの詳細は [example/deployment/docker/README.md](example/deployment/docker/README.md) を参照してください。

### メリット

1. **本番への影響なし**: ミラーリングされたトラフィックのレスポンスは無視されるため、本番に影響しません
2. **リアルなトラフィック**: 実際の本番トラフィックを使って新旧バージョンを比較できます
3. **タイムアウト分離**: Proxiffがタイムアウトしても本番サービスは影響を受けません
4. **段階的な検証**: 本番にデプロイする前に新バージョンの動作を検証できます

## ライセンス

MIT
