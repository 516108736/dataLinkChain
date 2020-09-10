# 数据上链

## 介绍
数据上链服务，主要是为了在[QuarkChain](http://devnet.quarkchain.io/)链上存放数据   
支持以下操作  
- 数据上链  
- 数据查询
    
**注意**: 要求 [Go 1.13+](https://golang.org/dl/)


## 服务启动

配置说明:   
- PrivateKey:QKC的钱包账户私钥
    + 主网: 请到[这里](https://mainnet.quarkchain.io/wallet)申请并查看QKC余额 
    + 测试网：请到[这里](https://devnet.quarkchain.io/wallet)申请并查看QKC余额
- Host:链对外提供的接口
    + 主网：  "http://34.222.230.172:38391"
    + 测试网："http://50.112.62.65:38391" 
- NetWorkID:链的networkid
    +主网：1
    +测试网：255    
    
供测试的配置：
-   主网
       +   ```bash
           {
             "PrivateKey": "0xdeb9010341b0aad25898017552177bd3fc88a9114a74316db871234b6f7eaa9f",
             "Host": "http://34.222.230.172:38391",
             "NetWorkID":1
           }
           ```  
-   测试网
       +    ```bash
            {
              "PrivateKey": "0x5a546e1cab61bda07dd9d964af917bd4bfb24bcb067c86ef479225772bce0053",
              "Host": "http://50.112.62.65:38391",
              "NetWorkID":255
            } 
            ```  

           

启动方式
```bash
# Clone the repository
mkdir -p $GOPATH/src/github.com/516108736
cd $GOPATH/src/github.com/516108736
git clone https://github.com/516108736/dataLinkChain
# 填写localConfig.json
go run main.go
```

## 调用方式
###数据上链
    curl -X POST -H 'content-type: application/json' --data '{"序号":123,"材料名称":"螺栓","材料类别":"土建","出库数量":666777888999000}' http://*.*.*.*:8080/
###数据查询
    curl http://*.*.*.*:8080/?txHash=0x3d3ec0c6cf6d716bb507585b090fc8da0248513e27d928559ec18ed35655767900000000