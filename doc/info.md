# toychain

a simple blockchain for fun

## Topic

1. Block (header, body)
1. Address
1. Transaction
1. Persistence
1. Network
1. Native token
1. Client SDK
1. Consensus (pos)

## Features

## 場景

1. 如何控制出塊時間每1秒一次?
1. 如何避免雙花攻擊 (double spending)?
1. Eth Nouce: https://segmentfault.com/a/1190000022628999

## 技術考量點

1. 帳戶模型: UTXO vs Account Model, 目前是使用 Account model
1. 共識模式型 (Pos vs Pow) => 這邊選擇 POS
1. GAS 模型

## Demo

1. 創建 `創世紀` 區塊，並建立 coinbase 100 萬到 God Address
1. God Address 轉帳 1000 到 `Account1` 地址
1. 查詢 `Account1` 的餘額
1. 從 `Account1` 轉帳到 `Account2`
1. 查詢 TX
1. 查詢 `Account1` 和 `Account2` 的餘額
1. 查詢區塊練高度

## TestAccount

1. GOD

- `Seed`: e34f9a593a9332b97395a60319fbf7186ca7f59168a5d73cdf203bfcdc02b0e1
- `Address`: LteoBMKhjHjV14rEcLF154CPs9BY6JYqexDwkMQc2cGEWDvAv
- `PrivateKey`: e34f9a593a9332b97395a60319fbf7186ca7f59168a5d73cdf203bfcdc02b0e12d29f680b4136828cd53abe747e6b023033d414d260fb64307d4ff8f3bf746e3

1. TestAccount1

- `Seed`: 12eddeb86bc3a03907decae3b63d597bfc0e97979520ae48560b4cb19dc1a823
- `Address`: 2vii4rSRYLgUP56q8jYUP6wQpNevqBDveLTrk1qB86gPEeHp3L
- `PrivateKey`: 12eddeb86bc3a03907decae3b63d597bfc0e97979520ae48560b4cb19dc1a823fdac012d28da54010f28c39089958315c4b3da8de15c022deed9c0ced1aa4471

1. TestAccount2

- `Seed`: fe49f00d1dc2d2670db543f57194dfae90fc603916cf9be5622d3b26483ee3cc
- `Address`: RyAc7Si4rWWuV3TxXu2mqa9jqg1RSbhnRjKWzJ4PiXf6nqSNK
- `PrivateKey`: fe49f00d1dc2d2670db543f57194dfae90fc603916cf9be5622d3b26483ee3cc38b1993bb1f51f58ba1c30e6861fa558647d4ce57e92000cc7882b9482de0949

## Reference

1. https://ithelp.ithome.com.tw/articles/10216297
1. https://www.chainnews.com/zh-hant/articles/547932547984.htm (Pos)