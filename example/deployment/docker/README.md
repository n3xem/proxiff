# Proxiff with Nginx Mirror - Docker Environment

このディレクトリには、Nginx mirrorモジュールを使ったProxiffのDocker環境が含まれています。

## アーキテクチャ

```
Client
  ↓
Nginx (Port 8000)
  ├─> Production Server (レスポンスをクライアントに返す)
  └─> Proxiff (mirror、レスポンスは無視)
        ├─> Newer Server
        └─> Current Server
              ↓
        差分検出とログ出力
```

### サービス構成

- **nginx**: リバースプロキシ + トラフィックミラーリング (Port 8000)
- **production**: 本番サーバー（クライアントにレスポンスを返す）
- **proxiff**: 差分検出プロキシ
- **newer**: 新バージョンのサーバー
- **current**: 現行バージョンのサーバー

## 使い方

### 前提条件

- Docker
- Docker Compose

### クイックスタート

```bash
# テストスクリプトを実行（ビルド、起動、テスト、クリーンアップまで自動）
./test-integration.sh
```

### 手動操作

#### 1. 起動

```bash
docker-compose up -d --build
```

#### 2. ログの確認

```bash
# 全サービスのログ
docker-compose logs -f

# Proxiffのログのみ
docker-compose logs -f proxiff

# Nginxのログのみ
docker-compose logs -f nginx
```

#### 3. リクエストの送信

```bash
# 基本的なリクエスト
curl http://localhost:8000/

# /api/users
curl http://localhost:8000/api/users

# /api/status
curl http://localhost:8000/api/status

# ヘルスチェック
curl http://localhost:8000/health
```

#### 4. 差分の確認

Proxiffのログを確認：

```bash
docker-compose logs proxiff | grep "Difference detected"
```

#### 5. 停止とクリーンアップ

```bash
docker-compose down -v
```

## 動作の仕組み

1. **クライアント → Nginx**
   - クライアントからNginxにリクエストが送信される

2. **Nginx → Production & Proxiff**
   - Nginxは本番サーバーにリクエストを転送
   - 同時に、`mirror`ディレクティブでProxiffにもリクエストをコピー

3. **Production → Client**
   - 本番サーバーのレスポンスがクライアントに返される

4. **Proxiff → Newer & Current**
   - ProxiffはミラーリングされたリクエストをNewerとCurrentの両方に転送
   - 両方のレスポンスを比較
   - 差分があればログに記録

5. **ミラーリングの特性**
   - Proxiffのレスポンスは無視される（クライアントには返らない）
   - Proxiffの処理がタイムアウトしても本番には影響しない

## カスタマイズ

### ポート変更

`docker-compose.yml`でNginxのポートを変更：

```yaml
nginx:
  ports:
    - "9000:80"  # 8000から9000に変更
```

### Nginx設定のカスタマイズ

`../nginx/nginx.conf`を編集して、以下をカスタマイズできます：

- タイムアウト設定
- ログフォーマット
- プロキシヘッダー
- ミラーリングの条件

### 独自のComparatorプラグインを使用

`docker-compose.yml`のproxiffサービスで、プラグインを使用：

```yaml
proxiff:
  build:
    context: ../../..
    dockerfile: Dockerfile
  command: [
    "./proxiff",
    "-newer", "http://newer:8080",
    "-current", "http://current:8080",
    "-port", "8080",
    "-plugin", "./plugin-status-only"
  ]
```

## トラブルシューティング

### サービスが起動しない

```bash
# サービスの状態を確認
docker-compose ps

# ログを確認
docker-compose logs

# 特定のサービスのログ
docker-compose logs <service-name>
```

### ヘルスチェックが失敗する

```bash
# ヘルスチェックの状態を確認
docker-compose ps

# コンテナに入って確認
docker-compose exec proxiff sh
wget -O- http://localhost:8080/
```

### ポートが既に使用されている

`docker-compose.yml`のポート設定を変更するか、既存のプロセスを停止：

```bash
# ポート8000を使用しているプロセスを確認
lsof -i :8000

# または
netstat -an | grep 8000
```

## 本番環境への適用

本番環境で使用する場合の推奨事項：

1. **リソース制限**
   ```yaml
   proxiff:
     deploy:
       resources:
         limits:
           cpus: '0.5'
           memory: 512M
   ```

2. **ログの永続化**
   ```yaml
   proxiff:
     volumes:
       - ./logs:/var/log/proxiff
   ```

3. **環境変数での設定**
   ```yaml
   proxiff:
     environment:
       - NEWER_URL=http://newer:8080
       - CURRENT_URL=http://current:8080
   ```

4. **モニタリング**
   - Prometheusメトリクスの追加
   - ログ集約（Fluentd、Logstash等）
   - アラート設定

## 関連ドキュメント

- [Nginx Mirror Module Documentation](https://nginx.org/en/docs/http/ngx_http_mirror_module.html)
- [Proxiff README](../../../README-ja.md)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
