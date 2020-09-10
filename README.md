# 数据上链

## 介绍
数据上链服务，主要是为了在QuarkChain链上存放数据   


支持以下网络
-   主网区块浏览器: https://mainnet.quarkchain.io/
-   测试网区块浏览器: https://devnet.quarkchain.io/
   
支持以下操作  
- 数据上链  
- 数据查询
    
**注意**: 要求 [Go 1.13+](https://golang.org/dl/)


## 服务启动

### 启动方式
```bash
# Clone the repository
mkdir -p $GOPATH/src/github.com/516108736
cd $GOPATH/src/github.com/516108736
git clone https://github.com/516108736/dataLinkChain
go build

# 申请私钥:到指定的区块链浏览器页面申请账户

# 申请QKC:
    主网：联系QKC官方人员获取
    测试网：https://devnet.quarkchain.io/faucet

# 查看指定网络的可用host
    打开上述区块浏览器，右上角可用看到部分节点的IP(启动时需"--host=***"来指明)

# 将私钥进行加密处理
./dataLinkChain -type=encrypt --private=申请到的私钥 --password=qkc

# 上面命令会生成新的私钥
./dataLinkChain --private=上个命令加密后的字符串 --password=qkc --host="http://34.222.230.172:38391"

```

## 调用方式

### 数据上链
    注意：
        1.数据格式必须为json格式
        2.目前数目长度最大为87926(byte) 
        
    curl -X POST -H 'content-type: application/json' --data '{"序号":123,"材料名称":"螺栓","材料类别":"土建","出库数量":666777888999000}' http://*.*.*.*:8080/
    

### 数据查询
    
    curl http://*.*.*.*:8080/?txHash=0x3d3ec0c6cf6d716bb507585b090fc8da0248513e27d928559ec18ed35655767900000000
    
