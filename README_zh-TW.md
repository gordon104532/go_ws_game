# go_ws_game
用 Golang 建立一個 WebSocket 服務，預期在 Google Cloud Run 上執行，最終只用 API 實作即可。

# 功能
- 多人遊戲: 玩家可以互相競爭，進行有趣的遊戲。
- 即時互動: WebSocket 連線讓玩家可以即時更新和互動。
- 排行榜: 追蹤玩家的進度和排名。
- API 驅動: 所有遊戲邏輯和互動都透過一個定義良好的 API 管理。
# 工作項目
## 後端
-遊戲邏輯: 實作核心遊戲機制，包括問題生成、答案驗證和分數計算。
-WebSocket 服務器: 處理 WebSocket 連線，管理玩家互動，並廣播遊戲更新。
-API 端點: 提供玩家登入、問題擷取、答案提交、排行榜存取和其他與遊戲相關操作的端點。
-資料儲存: 使用資料庫（例如 Cloud SQL）儲存遊戲資料、玩家資訊和排行榜排名。
## 前端
-使用者介面: 提供友善的介面讓玩家與遊戲互動。
-WebSocket 客戶端: 連線到 WebSocket 服務器，並處理即時通訊。
-API 整合: 與後端 API 通訊，擷取遊戲資料並提交玩家動作。
## 測試
-單元測試: 徹底測試個別元件和函式。
-整合測試: 驗證不同元件之間的互動，以及遊戲的整體功能。
-端到端測試: 模擬真實世界情境，確保遊戲按預期運作。
## 部署
-GitHub Actions 觸發 build Cloud Run: 自動化建置和部署流程（尚未實作）。
## 方案一
多多認識波普貓 大家認識的波普貓

答題遊戲，答完題可以出一題。 排行榜。 所有操作都透過 API 完成。

## API 端點
- 登入: `POST /login?username=apple` - 允許玩家使用使用者名稱登入。
- 取題目: `GET /quiz` - 擷取新的測驗問題。
- 答題: `POST /quiz` - 提交測驗問題的答案。
- 新增題目: `PUT /quiz` - 新增新的測驗問題。
- 重置題目: `GET /quiz/init` - 將測驗重置為初始狀態。
- 列出題目: `GET /quiz/list` - 列出所有可用的測驗問題。
- 彩蛋檢查: `GET /easter_egg` - 檢查是否有彩蛋。
- 健康檢查: `GET /health_check` - 檢查服務的健康狀況。
- 舉報: `POST /report` - 允許玩家舉報問題或疑慮。
## 議題
- 題目、答案、答題比率 怎麼存？
- 排行榜怎麼存？怎樣算對？
## 流程
1. 進入網頁
2. 輸入暱稱 打 API 設定暱稱（`POST /login?username=apple`）
3. 開始按鈕
4. 打 API 取得題目 （`GET /quiz`）
5. 顯示題目在上方 並有 A B 兩個選項
6. 選擇選項 打 API 送出 （`PUT /quiz`）
7. 顯示結果
8. 重複第四步 直到回傳為空
9. 顯示排行榜 （`GET /leader_board`）
## 備註
- 浮動 IP 怎麼打到？
- 用 Docker 掛載 Volume 做持久化？
- 0502 待實作功能
  - 舉報
  - 留言
  - 察看結果 不需要再回答一次
- 0507 待實作功能
  - ACME challenge:
    - apk add certbot
    - certbot certonly --webroot （程式運行中使用）
    - 將 domain 的 80 port 導向服務。
    - 將 static 導向憑證資料夾。
    - 注意每個帳戶、每個域名、每小時只能嘗試 5 次
  - 證書轉換
    - 將導出的 IIS 格式的证书轉換成 pem 和 key 檔案。
    - fullchain.pem 是 "crt" 檔案。
    - privkey.pem 是 "key" 檔案。
    - cp /etc/letsencrypt/live/gordon-the-owl.tw/fullchain.pem /go/src/.ssl
    - cp /etc/letsencrypt/live/gordon-the-owl.tw/privkey.pem /go/src/.ssl

## 起服務
- 本地起 HTTP 8080 (供本地前端使用)
- Docker 起 HTTPS 443
- 共用 Redis
