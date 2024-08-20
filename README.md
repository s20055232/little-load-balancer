# Simple Load Balancer

> 一個簡單的負載平衡器(load balancer)

## 定義

一個基本的負載平衡器(load balancer)應該滿足以下幾個核心功能:

- 流量分配: 能夠將進入的網路流量均勻地分配到多個後端伺服器。
- 健康檢查: 定期檢查後端伺服器的狀態,確保只將流量導向健康的伺服器。
- 會話持久性: 能夠將來自同一客戶端的請求始終導向同一台後端伺服器,以維持會話狀態。
- 協議支援: 至少支援常見的網路協議,如HTTP、HTTPS和TCP。
- 可擴展性: 能夠輕鬆地添加或移除後端伺服器,以應對流量變化。
- 基本的故障轉移: 當某個後端伺服器失效時,能夠自動將流量重新分配到其他健康的伺服器。
- 簡單的監控和報告: 提供基本的統計資訊,如伺服器狀態、流量分配情況等。

## 使用

將專案 clone 下來之後，先對 Load Balancer 進行設定，Load Balancer 所導流的 address 是透過 `config.toml` 設定的。

設定後執行

```go
go run .
```

這時候可以看到服務運行並 listen port 8000，接著啟動不同 Shell 連續執行範例用的 service 於 `config.toml` 中設定的 address

```go
go run ./simple_server/simple_server.go -port=13241
go run ./simple_server/simple_server.go -port=13242
...
```

會看到 Simple Load Balancer 正常輸出 Ping 到不同 service 的 log。

接著使用瀏覽器輸入 `http://127.0.0.1:8000` ，應該會看到 `Hello World` 出現在畫面上

## 實現

- [x] 使用 Round-Robin 演算法進行流量分配
- [x] 背景運行的 Health Check 確認服務的健康狀態
- [x] 將流量分配到正常響應的服務
- [x] 當死去的服務重新啟動時會將該服務重新放回可以導流的對象中
- [ ] 回傳被導流的服務的 response
- [ ] 簡單的監控和報告

## 注意

這個 Load Balancer 不考慮：

- 不同協議，目前只考慮 HTTP
- 持久性，Session 是沒有紀錄的
- 可擴展性
