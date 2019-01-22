# fabric链数据浏览

## 基本模式

使用轮询模式，调用api/v1/chain方法可获得当前的区块高度：

```

{
    "jsonrpc": "2.0",
    "result": {
        "Height": 1
    }
}

```

当发现有新块时，从当前的扫描进度开始查询所有新块（api/v1/chain/blocks/\<block number\>）并收集数据，块数据的格式如下：

```

{
    "jsonrpc": "2.0",
    "result": {
        "Height": "0",
        "Hash": "Local",
        "TimeStamp": "2019-01-22 21:21:31.8927643 +0800 CST m=+7.614275001",
        "Transactions": [
            {
                "Height": "0",
                "TxID": "4D65822107FCFD52157C2E4FBB450734",
                "Chaincode": "local",
                "Method": "batch",
                "CreatedFlag": false,
                "ChaincodeModule": "AtomicEnergy_v1",
                "Nonce": "0E846F0C43AA1A912D709995CED515A025444F7C",
                "Detail": {
                    "txs": [
                        {
                            "method": "MTOKEN.INIT",
                            "payload": "CgRFT1NYEgwKChUtAsfhSvaAAAA="
                        },
                        {
                            "method": "MTOKEN.ASSIGN",
                            "payload": "CgRFT1NYEiQKChUtAsfhSvaAAAASFgoUlPW6it0B/+o6tN44+N9iSXh5JRQ="
                        }
                    ]
                }
            }
        ],
        "TxEvents": [
            {
                "TxID": "4D65822107FCFD52157C2E4FBB450734",
                "Chaincode": "local",
                "Name": "TOKENNAME",
                "Status": 0,
                "Detail": "No parser can be found for this transaction/event",
                "Data": "454F5358"
            },
            {
                "TxID": "4D65822107FCFD52157C2E4FBB450734",
                "Chaincode": "local",
                "Name": "TOKENNAME",
                "Status": 0,
                "Detail": "No parser can be found for this transaction/event",
                "Data": "454F5358"
            }
        ]
    }
}

```

+ Transactions: 块内包含的事务列表，数组内每个元素是一个事务。一个处理程序通常应当处理下列成员：

    + Method：事务执行的方法名，处理程序可以使用此方法名识别其需处理的事务
    + Detail: 此对象中包含可读取的内容

+ TxEvents：块内事务的事件列表

## 需处理事务列表

在钱包和token业务中仅需处理转账和分配方法，对应的方法名和可读内容如下：

### 转账 \[TOKEN.TRANSFER\]

一个转账事务例子如下

```
{
    "jsonrpc": "2.0",
    "result": {
        "Height": "2",
        "TxID": "D5104DC76695721D157C339B297E33C4",
        "Chaincode": "local",
        "Method": "TOKEN.TRANSFER",
        "CreatedFlag": false,
        "ChaincodeModule": "AtomicEnergy_v1",
        "Nonce": "91089A2AA19EB32FE1A53A4A0321D530E196FBB9",
        "Detail": {
            "amount": "100000000000000000000000",
            "from": "ARJtq6Q46oTnxDwvVqMgDtZeNxs7Ybt81A",
            "to": "Ad3xQZRYJZr5Bu-o6rqfmIBQqf3JHALeNA"
        }
    }
}

```

其中Detail对象中包含如下成员

+ amount: 转账的数量
+ from, to: 转账源地址和目标地址

### 转账 \[TOKEN.ASSIGN\]

一个转账事务例子如下

```
{
    "jsonrpc": "2.0",
    "result": {
        "Height": "1",
        "TxID": "78629A0F5F3F164F157C339AC4FED458",
        "Chaincode": "local",
        "Method": "TOKEN.ASSIGN",
        "CreatedFlag": false,
        "ChaincodeModule": "AtomicEnergy_v1",
        "Nonce": "3131313131",
        "Detail": {
            "amount": "200000000000000000000000",
            "to": "ARJtq6Q46oTnxDwvVqMgDtZeNxs7Ybt81A"
        }
    }
}

```

其Detail对象成员和转账相同，但没有from地址