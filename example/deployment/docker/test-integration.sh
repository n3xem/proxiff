#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "========================================="
echo "Proxiff Integration Test with Nginx Mirror"
echo "========================================="
echo ""

# カラー定義
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# クリーンアップ関数
cleanup() {
    echo ""
    echo "${YELLOW}Cleaning up...${NC}"
    docker compose down -v
    echo "${GREEN}Cleanup completed${NC}"
}

# エラー時にクリーンアップ
trap cleanup EXIT

# 既存のコンテナを停止
echo "${YELLOW}Stopping existing containers...${NC}"
docker compose down -v 2>/dev/null || true
echo ""

# コンテナをビルドして起動
echo "${YELLOW}Building and starting containers...${NC}"
docker compose up -d --build

# ヘルスチェック
echo ""
echo "${YELLOW}Waiting for services to be healthy...${NC}"

# 最大120秒待機
MAX_WAIT=120
ELAPSED=0

while [ $ELAPSED -lt $MAX_WAIT ]; do
    HEALTHY=$(docker compose ps | grep -c "healthy" || true)
    TOTAL=$(docker compose ps | grep -c "Up" || true)

    echo "  Healthy services: $HEALTHY / $TOTAL"

    if [ "$HEALTHY" -eq 5 ]; then
        echo "${GREEN}All services are healthy!${NC}"
        break
    fi

    sleep 2
    ELAPSED=$((ELAPSED + 2))
done

if [ $ELAPSED -ge $MAX_WAIT ]; then
    echo "${RED}Services failed to become healthy within ${MAX_WAIT} seconds${NC}"
    docker compose ps
    docker compose logs
    exit 1
fi

echo ""
echo "${YELLOW}Running integration tests...${NC}"
echo ""

# テスト1: 基本的なリクエスト
echo "Test 1: Basic request to nginx (should return production response)"
RESPONSE=$(curl -s http://localhost:8000/)
echo "Response: $RESPONSE"

if echo "$RESPONSE" | grep -q '"version":"production"'; then
    echo "${GREEN}✓ Test 1 passed: Production response received${NC}"
else
    echo "${RED}✗ Test 1 failed: Expected production response${NC}"
    exit 1
fi

sleep 1
echo ""

# テスト2: /api/users エンドポイント
echo "Test 2: Request to /api/users"
RESPONSE=$(curl -s http://localhost:8000/api/users)
echo "Response: $RESPONSE"

if echo "$RESPONSE" | grep -q '"version":"production"'; then
    echo "${GREEN}✓ Test 2 passed: Users endpoint returned production response${NC}"
else
    echo "${RED}✗ Test 2 failed: Expected production response${NC}"
    exit 1
fi

sleep 1
echo ""

# テスト3: /api/status エンドポイント
echo "Test 3: Request to /api/status"
RESPONSE=$(curl -s http://localhost:8000/api/status)
echo "Response: $RESPONSE"

if echo "$RESPONSE" | grep -q '"version":"production"'; then
    echo "${GREEN}✓ Test 3 passed: Status endpoint returned production response${NC}"
else
    echo "${RED}✗ Test 3 failed: Expected production response${NC}"
    exit 1
fi

echo ""
echo "${YELLOW}Waiting for proxiff to process mirrored requests...${NC}"
sleep 3

# テスト4: Proxiffのログを確認
echo ""
echo "Test 4: Check proxiff logs for difference detection"
echo "${YELLOW}Proxiff logs:${NC}"
docker compose logs proxiff | tail -20

DIFF_COUNT=$(docker compose logs proxiff | grep -c "response difference detected" || true)

if [ "$DIFF_COUNT" -gt 0 ]; then
    echo ""
    echo "${GREEN}✓ Test 4 passed: Proxiff detected $DIFF_COUNT differences${NC}"
else
    echo ""
    echo "${RED}✗ Test 4 failed: No differences detected in proxiff logs${NC}"
    exit 1
fi

# テスト5: ヘルスチェックエンドポイント
echo ""
echo "Test 5: Health check endpoint"
HEALTH_RESPONSE=$(curl -s http://localhost:8000/health)

if echo "$HEALTH_RESPONSE" | grep -q "healthy"; then
    echo "${GREEN}✓ Test 5 passed: Health check endpoint works${NC}"
else
    echo "${RED}✗ Test 5 failed: Health check failed${NC}"
    exit 1
fi

# サービス状態の表示
echo ""
echo "${YELLOW}Final service status:${NC}"
docker compose ps

echo ""
echo "========================================="
echo "${GREEN}All integration tests passed!${NC}"
echo "========================================="
echo ""
echo "Architecture:"
echo "  Client -> Nginx (port 8000)"
echo "            ├─> Production Server (returns response to client)"
echo "            └─> Proxiff (mirror, response ignored)"
echo "                  ├─> Newer Server"
echo "                  └─> Current Server"
echo ""
echo "Proxiff compares responses from Newer and Current servers"
echo "and logs any differences detected."
echo ""
