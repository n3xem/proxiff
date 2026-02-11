# Proxiff - クイックスタートガイド

このガイドは、5分でProxiffを使い始めるのに役立ちます。

[English Quick Start Guide](QUICKSTART.md)

## Proxiffとは？

Proxiffは、2つの異なるサーバー（newer と current）からのレスポンスを比較し、差分をログに記録しながら、クライアントにはcurrentサーバーのレスポンスを返すDiffy風のHTTPプロキシです。

## クイックスタート

### 1. デモを実行

Proxiffの動作を確認する最も簡単な方法：

```bash
# デモスクリプトを実行可能にする（まだの場合）
chmod +x example/run-demo.sh

# デモを実行
./example/run-demo.sh
```

これにより以下が起動します：
- ポート8081でCurrentサーバー
- ポート8082でNewerサーバー
- ポート8080でProxiffプロキシ

### 2. テストする

別のターミナルで、以下のリクエストを試してください：

```bash
# 基本エンドポイント - バージョンの違いが表示されます
curl http://localhost:8080/

# ユーザーエンドポイント - newer版にはemailフィールドがあります
curl http://localhost:8080/api/users

# ステータスエンドポイント - 異なるステータスコード（200 vs 201）
curl http://localhost:8080/api/status
```

### 3. ログを確認

Proxiffプロキシは検出した全ての差分をログに記録します：

```
2025/11/23 14:35:35 Difference detected: Body differs: newer="...", current="..."
2025/11/23 14:35:50 Difference detected: Body differs: newer="...", current="..."
2025/11/23 14:35:51 Difference detected: Status code differs: newer=201, current=200
```

## 手動セットアップ

手動で実行したい場合：

### 1. ビルド

```bash
go build -o proxiff ./cmd/proxiff
go build -o sample-server ./example/servers
```

### 2. サーバーを起動

```bash
# ターミナル1: Current版
./sample-server -port 8081 -version current

# ターミナル2: Newer版
./sample-server -port 8082 -version newer

# ターミナル3: Proxiff
./proxiff -newer http://localhost:8082 -current http://localhost:8081 -port 8080
```

### 3. リクエストを送信

```bash
curl http://localhost:8080/
```

## 独自のサーバーで使用

サンプルサーバーを実際のサービスに置き換えます：

```bash
./proxiff \
  -newer http://your-new-service.example.com \
  -current http://your-current-service.example.com \
  -port 8080
```

その後、クライアントを現在のサービスではなく `http://localhost:8080` に向けます。

## テスト

全てのテストを実行：

```bash
go test ./... -v
```

## 次のステップ

- 詳細なドキュメントについては [README-ja.md](README-ja.md) をお読みください
- 比較ロジックの実装は [comparator/simple.go](comparator/simple.go) を参照してください

## 一般的なユースケース

1. **カナリアデプロイメント**: 本番環境とカナリア版を比較
2. **マイグレーションテスト**: 新実装が旧実装の動作と一致することを保証
3. **A/Bテスト**: 異なるアルゴリズム実装を比較
4. **リグレッションテスト**: APIレスポンスの予期しない変更を検出

## トラブルシューティング

### ポートが既に使用されている

もしポート8080が既に使用されている場合は、`-port`フラグで別のポートを指定してください：

```bash
./proxiff -newer http://localhost:8082 -current http://localhost:8081 -port 9000
```

### サーバーが起動しない

サーバーが既に実行されているか確認してください：

```bash
# macOS/Linux
lsof -i :8080
lsof -i :8081
lsof -i :8082

# プロセスを停止
pkill -f "proxiff"
pkill -f "sample-server"
```

## より詳しく

Proxiffは完全にTDD（テスト駆動開発）で実装されています。コードを読むことで、Go言語での良いプラクティスとテスト手法を学ぶことができます。

- `comparator/` - インターフェースの設計とテスト
- `proxy/` - HTTPプロキシの実装とテスト
- `example/` - 実践的な使用例

ハッピーコーディング！
