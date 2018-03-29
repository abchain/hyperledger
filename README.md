## Abchain's Hyperledger Utilities

A package for building chaincodes which are worked in hyperledger-fabric, made your work can immigrate between 0.6, 1.x and YA-fabric smoothly

一个可以方便地构造chaincode并平滑地在YA-fabric，fabric 1.x，fabric 0.6之间迁移的工具集

Currently this project include:

项目当前的内容包括：

* A [framework](https://github.com/abchain/hyperledger/tree/master/chaincode/lib) for building the immigratable chaincode

  用于构建可移植的chaincode的框架

* chaincode modules which can be used for many common business such as transfer of assets, controlling of accounts etc

  实际业务常用的chaincode模块实现，包括资产转移，账户管理等

* An [utility](https://github.com/abchain/hyperledger/tree/master/tx) helping to build the transactions sent to chaincode

  协助构建chaincode事务（transaction）的工具集

* A cryptography account scheme similar to which in the popular crytocurrency (bitcoin or ethereum) and can be easily applied to chaincode

  一个和常见加密货币（比特币和以太坊）类似，可以简单地使用的密码学账户方案

* Works helping to build a REST-API client for the construction and dispatching of transactions

  用于构建和发送事务的REST风格的客户端实现

* Scripts for deployment of YA-fabric

  一些用于部署YA-fabric的脚本集

## License <a name="license"></a>
This project uses the [Apache License Version 2.0](LICENSE) software license.
