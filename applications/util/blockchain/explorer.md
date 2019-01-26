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

### 转账 \[TOKEN.TRANSFER\], \[MTOKEN.TRANSFER\]

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
+ token: token名，对于主token，这一项不出现

### 分配 \[TOKEN.ASSIGN\], \[MTOKEN.ASSIGN\]

一个分配事务例子如下

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

### 批处理方法

批处理方法的方法名由chaincode定义（例如在示例chaincode中一个批处理方法名为"batch"），其Detai对象为批方法执行的所有事务列表，列表每个元素均包含Method和Detail方法，处理程序可以用处理单个事务的逻辑处理列表中每个元素

一个批处理事务的例子如下

```
    {
        "Height": "0",
        "TxID": "4D65822107FCFD52157D0BD188A97B60",
        "Chaincode": "local",
        "Method": "batch",
        "CreatedFlag": false,
        "ChaincodeModule": "AtomicEnergy_v1",
        "Nonce": "179564CF68A50DFFFBDD6A31D5A3A324C2B30722",
        "Detail": [
            {
                "Method": "MTOKEN.INIT",
                "Detail": "EOSA: Unknown message"
            },
            {
                "Method": "MTOKEN.ASSIGN",
                "Detail": {
                    "amount": "100000000000000000000000000",
                    "to": "AS6mOHOU3PKLpmHXinFgxK-3GhDb9YuC2g",
                    "token": "EOSA"
                }
            }
        ]
    }
```

列表的第二个元素是一个可以被解析的MTOKEN.ASSIGN方法

## 事务错误

任何提交的事务无论是否正确地执行都会返回其事务ID，但是错误的事务将不会被包含在区块中，通常地（取决于区块链平台的实现）执行此事务的区块将提供一个事件告知事务的执行错误

例如本地链模块（模拟YA-fabric的行为）在区块中返回错误执行的事务的事件例子如下

```
{
    "jsonrpc": "2.0",
    "result": {
        "Height": "1",
        "Hash": "Local",
        "TimeStamp": "2019-01-25 17:54:34.9503012 +0800 CST m=+6.546807301",
        "TxEvents": [
            {
                "TxID": "78629A0F5F3F164F157D0EC263FD3FEC",
                "Chaincode": "local",
                "Name": "INVOKEERROR",
                "Status": 1,
                "Detail": "Local invoke error: No enough balance"
            }
        ]
    }
}
```

事务执行错误的事件的Status成员值总是非0（正确执行事务的事件Status为0），其Name和Detail成员通常包含了执行错误的原因描述

