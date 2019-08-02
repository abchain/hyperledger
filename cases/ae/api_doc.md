FORMAT:1A

# Group HTTP Example

## SUCCESS [/api/v1/success]

### 成功响应示例 [GET]

- result 为业务相关的数据，详见各个接口的定义

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result" : "response object"
            }

## ERROR [/api/v1/error]

### 失败响应示例 [GET]

- error
    - code: 错误码
    - message: 错误信息
    - data: 额外的错误相关数据

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "error" : {
                    "code" : -100,
                    "message" : "error message",
                    "data" : "data object"
                }
            }

# Group Account

## 账号管理 [/api/v1/account/{accountID}/{index}]

### 创建账号 [POST]

- 请求参数说明

    - accountID: 新创建账号的账号 ID
        - 账号 ID 在本地需要保证唯一性
        - 账号 ID 仅用于方便本地操作，该信息并未在机器间同步
    - \[index\]: 使用已经创建的账号创建子账号
        - 账号 ID 之前必须已经创建
        - 创建的子账号将会被记录在本地，并以 accountID : index 的形式显示 

- 响应参数说明

    - result: 成功创建的账号地址

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

            accountID=account01

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result": "ATVmjywxE9maiZY92vySKfupRiu3tg0G-Q"
            }

### 查询账号地址 [GET]

- 响应参数说明

    - result: 账号 ID 对应的账号地址

+ Parameters

    + accountID: `account01` (string, required) - 账号 ID

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result": "ATVmjywxE9maiZY92vySKfupRiu3tg0G-Q"
            }

### 修改账号别名 [PATCH]

- 请求参数说明

    - newAccountID: 新账号 ID

- 响应参数说明

    - result: 账号地址

+ Parameters

    + accountID: `account01` (string, required) - 原账号 ID

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

            newAccountID=account02

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result": "ATVmjywxE9maiZY92vySKfupRiu3tg0G-Q"
            }

### 删除账号 [DELETE]

- 响应参数说明

    - result: 成功删除的账号地址

+ Parameters

    + accountID: `account01` (string, required) - 账号 ID

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result": "ATVmjywxE9maiZY92vySKfupRiu3tg0G-Q"
            }

### 列出本地所有账号 [GET]

- GET 响应参数说明

    - result: 所有账号信息
        - key 为账号 ID
        - value 为账号地址

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result": {
                    "account01" : "ATVmjywxE9maiZY92vySKfupRiu3tg0G-Q",
                    "account02" : "AWVQUuLC66BMT71kJeY11wDnDdbtltZNUA",
                    "account03" : "AfN2Wq9ISsClm8wuqmgxt92oHra72YzvHA"
                }
            }

### 获取子账号地址 [GET]

- 响应参数说明

    - result: 子账号地址
    - error
        - code 为 -100 时，表示该索引值的子账号无效，需更换索引值重新获取

+ Parameters

    + accountID: `account01` (string, required) - 账号 ID
    + index: `100` (number, required) - 子账号索引，从 1 开始计数

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result": "AWVQUuLC66BMT71kJeY11wDnDdbtltZNUA"
            }

## 密钥管理 [/api/v1/privkey/{accountID}]

### 导入密钥 [POST]

- 请求参数说明

    - accountID: 导入密钥后创建的账号 ID
    - privkey: 密钥

- 响应参数说明

    - result: 成功导入的账号地址

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

            accountID=account01&privkey=tKo3QrjyPfzTHJkbQ_ALANnLVxavKt77h3GICqZ2q38=

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result": "ATVmjywxE9maiZY92vySKfupRiu3tg0G-Q"
            }

### 导出密钥 [GET]

- 响应参数说明

    - result: 账号密钥

+ Parameters

    + accountID: `account01` (string, required) - 账号 ID

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result": "MCUCAQECIDtF+PyxKKkGz5sMRUYMZ4kYtTk+W5EwOfoGyo3ZJIaP"
            }

# Group Registrar

## 申请登记账号相关 [/api/v1/registrar/{accountID}]

### 申请登记账号 [POST]

- 请求参数说明

    - accountID | account : 待登记的账号 ID 或账号
    - publicKey : 申请登记的公钥字符串，以16进制表示，如果指定账号，此参数被忽略
    - [usage] : 所登记账号的作用，应当是用逗号分割的字符串，当前保留

- 响应参数说明

    - result: Fabric transaction ID

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

            accountID=account01

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result": "ec239f5e-06ff-497a-96e2-d3ee9d266867"
            }

### 查询登记结果 [GET]

- 响应参数说明

    - result: 审批结果
        - 1 : 等待审批
        - 2 : 审批通过
        - 3 : 审批被拒绝

+ Parameters

    + accountID: `account01` (string, required) - 待查询登记结果的账号 ID

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "status": "0",
                "result": 2
            }

## 提交登记审批结果 [/api/v1/registrar/audit]

### 提交登记审批结果 [POST]

- 说明

    - 该请求只有在管理员节点上才能执行成功

- 请求参数说明

    - address: 待审批的账号地址
        - 一次请求可携带多个 address 参数
    - pass
        - true: 审批通过
        - fasle: 审批不通过

- 响应参数说明

    - result: Fabric transaction ID

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

            pass=true&address=AfN2Wq9ISsClm8wuqmgxt92oHra72YzvHA&address=AfN2Wq9ISsfadsfafsbsVDF

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result": "ec239f5e-06ff-497a-96e2-d3ee9d266867"
            }



# Group Assign

## token分配量初始化 [/api/v1/assign/init]

### 执行初始化 [POST]

- 说明
    - 此事务只执行一次，初始化chaincode中记录的token信息
    - 此事务的执行要求特权账号

- 请求参数说明

    - total: 设置总token数量，之后可以使用assign方法进行分配

- 响应参数说明

    - result: 无，显示“OK”

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

            total=100000

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "status": "0",
                "result": "ok"
            }

## token分配 [/api/v1/assign]

### 执行分配 [POST]

- 说明
    - 将当前未分配的币余额分配到特定的账号
    - 此事务的执行要求特权账号

- 请求参数说明

    - to: 受付人地址
    - amount: 转账金额
    - nonce: (可选)此次分配的唯一标识，相同 nonce 值的事务在一定时间（1小时）内不会被重复收入区块

- 响应参数说明

    - result: 转账事务 ID (fundID)，可以在fund方法中进行查询

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

            to=AfN2Wq9ISsClm8wuqmgxt92oHra72YzvHA&amount=100000

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "status": "0",
                "result": "ec239f5e06ff497a96e2d3ee9d266867"
            }

### 查询全局状态 [GET]

- 响应参数说明

    - result: token数量分配信息
        - total: 总token的数量
        - unassign: 未分配的token数量

 
+ Parameters

    + fundID: `ec239f5e06ff497a96e2d3ee9d266867` (string, required) - 转账事务 ID，POST 请求的响应中提取

+ Response 200 (application/json;charset=utf-8)

    + Body
            {
                "jsonrpc": "2.0",
                "result": {
                    "total": "1000000000000000000000000",
                    "unassign": "600000000000000000000000"
                }
            }

# Group Fund

## 转账相关 [/api/v1/fund/{fundID}]

### 转账 [POST]

- 请求参数说明

    - accountID | account: 支付人的账号ID或账号地址
        - 账号地址必须是已经记录在本地的地址，可以是根账号或者子账号
    - \[index\]: 使用 accountID 的子账号
    - from: 支付人地址，如果指定账号，此参数被忽略
    - to: 受付人地址
    - amount: 转账金额
    - nonce: (可选)此次转账的唯一标识，相同 nonce 值的事务在一定时间（1小时）内不会被重复收入区块

- 响应参数说明

    - txID: 事务 ID 
    - Nonce: 事务使用的Nonc，以16进制表示。如果用户指定一个字符串为Nonce，将在后面用()提示
    - Data：每个转账事务唯一的标识（fundNonce）。具有相同唯一标识的转账事务只会被执行一次。

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

            accountID=account01&index=0&to=AfN2Wq9ISsClm8wuqmgxt92oHra72YzvHA&amount=100000&nonce=fdsaf12313

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "status": "0",
                "result": {
                    "txID": "365A858149C6E2D115B6F8BF1F76165C",
                    "Nonce": "31313131 (1111)",
                    "Data": "E1917A90F64FABE16A92E2DD0C2046B4F75C5D18FB0882136FD9027460B541DF"
                }
            }


- 事件

    每个成功执行的转账事务会在对应的区块产生一个名为TRANSFERTOKEN的事件，其detail是转账唯一标识（即返回值中的Data）

### 查询转账结果 [GET]

- 说明
    - fund 事务未上链时，通过本接口查询到的 result 为空

- 响应参数说明

    - result: 转账事务详细信息
        - state: 转账结果
            - 1: 转账成功
            - 2: 转账失败
        - error: 转账失败但是事务上链时，具体的失败原因
        - from: 支付人地址
        - to: 受付人地址
        - amount：转账金额
        - \[external\]: 表示转账由外部chaincode执行
        - time: 事务上链时间
 
+ Parameters

    + fundID: `ec239f5e06ff497a96e2d3ee9d266867` (string, required) - 转账事务 ID，POST 请求的响应中提取

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc": "2.0",
                "result": {
                    "txID" : "fafdafasfdafdafdafda"
                    "state": 1,
                    "error": "",
                    "from": "CAESFH-nFRqmcVEGQ31enumWC-PrfLiuGgTFYBnHIgh_pxUapnFRBg==",
                    "to": "CAESFFA3jJnD6W3tgf6664j780gIKkvuGgQOgdxdIghQN4yZw-lt7Q==",
                    "amount": 1000,
                }
            }


# 多Token支持

上述方法可以在Group名前添加token路径\[token.<Token Name>\]，此时方法将操作token Name对应的token，参数和输出格式不变，例如对名为VCTX的token执行转账，方法URL为：

- **\[POST\] /api/v1/token.VCTX/fund**


*Token名只能由字母和数字构成，长度为4-16字节之间*

## Token创建 [/api/v1/adv/create]

### 创建一个token，token将首先分配给创建者地址，再由创建者进行后续的分发

- 请求参数说明

    - accountID | account: 创建者的账号ID或账号地址
        - 账号地址必须是已经记录在本地的地址，可以是根账号或者子账号
    - \[index\]: 使用 accountID 的子账号
    - to: token分配到的地址，如果指定账号，此参数被忽略
    - total: 创建的token的总量
    - name: token名

- 响应参数说明

    - txID: 创建事务 ID

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

            total=100000000000000000000000&to=AZT1uordAf_qOrTeOPjfYkl4eSUUwEJT7Q&name=VCT1

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc": "2.0",
                "result": {
                    "txID": "2019-01-21 11:28:56.7305051 +0800 CST m=+640.183320601",
                    "Nonce": "0VWmW8+oPRFzn8dYKIRiwr3eW4w=",
                    "Data": [
                        {
                            "fundNonce": "6nFcIGIEh9B_MUbkG4TytnP3_7ynAjedOMURhxpmZkQ=",
                        }
                    ]
                }
            }

# Group subscription

## 分润 [/api/v1/subscription/{subscriptionAddr}]

### 注册分润协议 [POST]

- 请求参数说明

    - accountID | account: 分润账户的账号 ID 或账号地址
        - 账号地址必须是已经记录在本地的地址，可以是根账号或者子账号
    - \[index\]: 使用 accountID 的子账号
    - initiator: 发起协议的地址，如果指定账号，则此参数被忽略
    - contract: 分润方案，是\[ 地址 : 权重 \] 形式的字符串，代表接收分润的地址和对应的分润比例，可以包含多个contract字段

- 响应参数说明

    - result
        - subscriptionAddress: 分润协议入账的地址，注意此地址不是注册用的分润账户的地址，可以使用此地址查询到对应的分润协议

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

            accountID=account01&contract=AfN2Wq9ISsClm8wuqmgxt92oHra72YzvHA:50&contract=jJnD6W3tgf6664j780gIKkvuGgQOgdxdI:50

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "status": "0",
                "result": {
                    "subscriptionAddress" : "tgf6664j780gIKkvuGgQOgdxdIghQN4",
                }
            }

### 查询分润协议 [GET]

- 说明
    - subscription 事务未上链时，通过本接口查询到的 result 为空

- 响应参数说明

    - result: 分润协议详细信息
        - address: 分润账户地址，即注册分润协议时返回的 subscriptionAddress
        - shares: 分润账户中累积的总金额
        - balance: 分润账户中当前的余额
        - contract: 分润协议，对象名为协议中每个分润的地址，值为如下参数
            - weight: 协议规定的分润比例，用 0 - 1 之间的小数表示
            - shares: 分润地址已提取的总金额
            - availiable: 分润地址当前仍可提取的金额
                - 由于分数计算的原因，协议所有账户的 availiable 金额之和与 balance 可能会有个位数的差异
 
+ Parameters

    + fundID: `ec239f5e06ff497a96e2d3ee9d266867` (string, required) - 分润事务 ID，POST 请求的响应中提取

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc": "2.0",
                "result": {
                    "address" : "tgf6664j780gIKkvuGgQOgdxdIghQN4"
                    "shares": 3000,
                    "balance": 2000,
                    "contract": {
                        "AfN2Wq9ISsClm8wuqmgxt92oHra72YzvHA":{
                            "weight": 50,
                            "shares": 0,
                            "availiable": 1500
                        }
                        "jJnD6W3tgf6664j780gIKkvuGgQOgdxdI":{
                            "weight": 50,
                            "shares": 1000,
                            "availiable": 500
                        }
                    }
                }
            }

## 利润提取 [/api/v1/subscription/redeem/{subscriptionAddr}]

### 提取利润 [POST]

- 说明
    - 利润提取将在资产 chaincode 上产生一个对应的转账事务，将分润账户中的资金转移到分润协议中对应的账户内

- 请求参数说明

    - accountID | account: 账号 ID 或账号地址
        - 账号地址必须是已经记录在本地的地址，可以是根账号或者子账号
        - 账号 ID 对应的地址必须包含在分润协议中
    - \[index\]: 使用 accountID 的子账号
    - \[amount\]: 希望从分润账户中提取的金额，默认提取当前所有可用的数目
    - \[to\]: 一个或多个执行分润的地址，必须包含在分润协议中。所有的分润地址都会提取amount中指定的金额或者可提取的最大金额
    
        **注意系统当前无法为多个分润地址同时生成凭据，在需要凭据的情况下，只能使用data/方法产生待签名的事务，并交给各个地址进行签名**

- 响应参数说明

    - result
        - 转账事务 ID (fundID)


# Group RawTransaction 

## 生成一个待签名的事务 [/api/v1/data]

- 说明

     此路径下可以连接上述业务API中的任何路径，结果将产生一个对应的待签名事务内容和需签名的hash值，调用者可以使用自己的私钥签名此hash并提交已签名的事务，例如 \[POST\] /api/v1/data/fund 将生成一个待签名的转账事务

- 响应参数说明

    - raw: 生成的待签名事务
    - hash: 此事务需要签名的hash值，以十六进制数表示
    - promise: 此事务如果是调用（invoke），提供和实际调用时相同的返回值；如果是查询（query），仅显示返回值中包含的数据内容，而不包含实际的值

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result": {
                    "raw":"I::TOKEN.TRANSFER:ChoKB0FCQ0hBSU4SD0F0b21pY0VuZXJneV92MRIGCPPk8OEFGhRrumYaeGyTASFlvIj4UyA1NTgckw==:CgsIRZUWFAFISgAAABIWChT9hccdqdkYsNsFR5nG+3qAMCdWnhoWChQSbaukOOqE58Q8L1ajIA7WXjcbOw==",
                    "hash":"03BD91127B5FED4EC9C0F71A516944880558E5EFC71520A38607189EC302251E",
                    "promise": {
                        "txID": "pending",
                        "Nonce": "H-5R9kjK42HSFuA1_h4CqY_8IfBdEAU2aE1FWE79gVA=",
                        "Data": "a7pmGnhskwEhZbyI+FMgNTU4HJM="
                    }
                }
            }

## 转换公钥到地址 [POST /api/v1/account/frompublickey]

- 请求参数说明

    - pubkeybuffer: 需要转换为地址的公钥内容，以16进制数表示，格式可以参考 application/util/signhelper 中的node.js示例

- 响应参数说明

    - result: 公钥对应的地址

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

            pubkeybuffer=EC:01,d0de0aaeaefad02b8bdc8a01a1b8b11c696bd3d66a2c5f10780d95b7df42645cd85228a6fb29940e858e7e55842ae2bd115d1ed7cc0e82d934e929c97648cb0a

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc": "2.0",
                "result": "ARJtq6Q46oTnxDwvVqMgDtZeNxs7Ybt81A"
            }


## 提交事务 [POST /api/v1/rawtransaction]]

- 请求参数说明

    - tx: 提交的事务内容，编码方案和生成待签名事务时相同
    - \[sig\]: 编码为字符串的签名值，生成的格式可以参考 application/util/signhelper 中的node.js示例
        - 一次请求可携带0 ~ 多个 sig 参数
    - \[accountID | account\]: 账号 ID 或账号地址，系统会使用此账号向事务中追加一个签名
        - 账号地址必须是已经记录在本地的地址，可以是根账号或者子账号
    - \[index\]: 使用 accountID 的子账号        

- 响应参数说明

    - result: 提交的事务 ID
        - 当前API不支持提交query类型的事务

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

            tx=I::TOKEN.TRANSFER:ChoKB0FCQ0hBSU4SD0F0b21pY0VuZXJneV92MRIGCKKw8OEFGhRiXies8Zp97ktRv1lyR4mZtZV8Vw==:CgoVLQLH4Ur2gAAAEhYKFJT1uordAf/qOrTeOPjfYkl4eSUUGhYKFBJtq6Q46oTnxDwvVqMgDtZeNxs7&sig=EC:01,d0de0aaeaefad02b8bdc8a01a1b8b11c696bd3d66a2c5f10780d95b7df42645cd85228a6fb29940e858e7e55842ae2bd115d1ed7cc0e82d934e929c97648cb0a,5f27d831cfe37e7542a1a5d9c687d935f0fd10dc60c2605be7a07ae26b77e22e23ebcbeed6ca7a1c9873009bc060ece0930d3013221efc87e9a4b1b1bb654b6c:

## 执行签名 [POST /api/v1/signature]]

- 请求参数说明

    - hash: 需要签名的hash，用16进制表示
    - accountID | account: 账号 ID 或账号地址
        - 账号地址必须是已经记录在本地的地址，可以是根账号或者子账号
    - \[index\]: 使用 accountID 的子账号

- 响应参数说明

    - result: 生成的签名数据，编码为可在事务提交中使用的格式

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

        hash=A0C248E68881CA5D00B4CCCD0CDD3CEF0747674CC13F6559825C9393FF8089ED
        &accountID=aaa

+ Response 200 (application/json;charset=utf-8)

    + Body

        {
            "jsonrpc": "2.0",
            "result": "EC:02,D7B150F1A79153F8CF8755E42154B58F194D9A6A0E7805A1BBAA528107DA25AC18466D511E2706D8B4E1D00B3C0533D6637ED3D050547B01FD24D3C9E9F6D673,EBAD9654E5A5BA6D56EBEC4DC2B15444AE6ACF49CF2BD4844E09EF269EEA65E34A6D0600C92BAC7FE97465F37B5A88E2E233ADBE54BD547CF01872DB3DE454A9:"
        }

# Group Blockchain

## 获取区块链基础信息 [/api/v1/chain]

### 获取区块链基础信息 [GET]

- 响应参数说明

    - result

        - height: 当前区块高度
        - currentBlockHash: 当前最新区块 hash
        - previousBlockHash: 前一区块 hash

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc" : "2.0",
                "result": {
                    "height": 4,
                    "currentBlockHash": "pbYJVyyQVX8JdO+zjJUYWP1Z8RvnRtyPy1mwJ7jdzPF47Be9WF+RVVjlCZJhOq+EGVFfqci+t7i/FkVWfAEvnw==",
                    "previousBlockHash": "kqFCSGw1Z7WYE3hGJe5Gyj4IHmPP2XcvkBHOtBcuVry57cIcFe0cWNtE8H3dCbXnKWjZSXBYh3+8+KJYQfUUGA=="
                }
            }

## 获取指定高度区块信息 [/api/v1/chain/blocks/{height}]

### 获取指定高度区块信息 [GET]

- 响应参数说明

    - result： 区块详细信息

        - stateHash: 当前世界状态的 hash 值
        - currentBlockHash: 当前区块 hash
        - previousBlockHash: 前一个 block 的hash 值
        - transactions: 本区块内收录的所有事务信息
            - 事务的具体字段描述，参见“查询事务信息”接口的说明
        - events: 本区块内事务所触发的事件通知信息

- event type 说明

```
# 登记根地址
REGISTRAR = 1;
# 登记审批
REGISTRAR_AUDIT = 2;

# 转账
FUND = 11;
# 分润注册
SUBSCRIPTION = 101;
# 分润提取
SUBSCRIPTION_REDEEM = 102;

```

- event 具体字段说明
    - 注意，字段值为空或零值时，该字段会自动隐去，所以如果解析时找不到相应字段，可以认为是空或零值

```
# 通用字段
type: event 类型
txID: 所属的事务 ID

# 登记地址相关字段
addr: 申请登记的地址
audit: 审批通过的地址

# 转账相关字段
fundID: 转账事务 ID
from: 支付人地址
to: 受付人地址
amount: 转账金额
payout: 实际支出
time: 发起转账时间
```


+ Parameters
    + height: `AWVQUuLC66BMT71kJeY11wDnDdbtltZNUA` (string, required) - 区块高度

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc": "2.0",
                "result": {
                    "stateHash": "yxeJh3zfkY66EX34MYEfo74wg3GxHe49hJIkGB013cn8o/a+8hoXXE4FRycf8UevalN5uCPHzZSieiu1i570fw==",
                    "currentBlockHash": "99DweQ11joWx/yYAJhLYciJ7+JjUc1aqp/YS8k4JoUr8vjZfSOTLEjQp1K4JusP2t1K/7XLpOkFMlrLzo+YnYQ==",
                    "previousBlockHash": "qEBXsVeenQyi27XzmPAm4cEHYJXOp66bfEPejTFvTV9AxELqRHoBzvc6YjNgHNsqFjbaM0b92MnxjK3rKOB0UA==",
                    "transactions": [
                        {
                            "chaincodeID": "Egh1Y3Rlc3QwMQ==",                            
                            "txid": "e410dffa-638a-472d-99c1-2329be066669",
                            "timestamp": {
                                "seconds": 1514212963,
                                "nanos": 747340544
                            }
                        }
                    ],
                    "events": [
                        {
                            "type": 14,
                            "txID": "e410dffa-638a-472d-99c1-2329be066669",
                            "fundID": "25721a9a-2bf9-4425-a5ee-850da10b50c7",
                            "from": "CAESFGN9u-3zAwwl7BhFm0sHrnTOYYk3GgSjfuoUIghjfbvt8wMMJQ==",
                            "to": "CAESFKt4yGBBQRFz04Iq6uMSa8aCgNQEGgTnrYCsIggL2hP7UNxlsSoBAQ==",
                            "amount": 1000,
                        }
                    ]
                }
            }

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc": "2.0",
                "result": {
                    "stateHash": "QwVlfAqC5Qb4gfRIjN9mcXMOTPwew/t4MGtSMIMYHDu/Qv3wJZYPXFwObR2VPfsydaBf4kX2o5RttXWcZZ+mgg==",
                    "currentBlockHash": "99DweQ11joWx/yYAJhLYciJ7+JjUc1aqp/YS8k4JoUr8vjZfSOTLEjQp1K4JusP2t1K/7XLpOkFMlrLzo+YnYQ==",
                    "previousBlockHash": "AMc07SeXt+AzH/QlajEVyU2TuQ4TUx3nnPDn4+MIbnWrT/ASFI4mORcfoIj1jA2xapn5K9ma9Zm+MvtIr/3kJg==",
                    "transactions": [
                        {
                            "chaincodeID": "Egh1Y3Rlc3QwMQ==",
                            "txid": "63108a90-de16-4a3d-8669-f67f2ac3ecb6",
                            "timestamp": {
                                "seconds": 1514212705,
                                "nanos": 33062677
                            }
                        },
                        {
                            "chaincodeID": "Egh1Y3Rlc3QwMQ==",                            
                            "txid": "cfff165d-878c-4fbb-b5ae-67f5437f6dc5",
                            "timestamp": {
                                "seconds": 1514212707,
                                "nanos": 545489576
                            }
                        }
                    ],
                    "events": [
                        {
                            "type": 1,
                            "txID": "63108a90-de16-4a3d-8669-f67f2ac3ecb6",
                            "addr": "CAESFAvaE_tQ3GWxy2jrkeVSXPhHJA-yGgTGGg0cIggL2hP7UNxlsQ=="
                        },
                        {
                            "type": 1,
                            "txID": "cfff165d-878c-4fbb-b5ae-67f5437f6dc5",
                            "addr": "CAESFGN9u-3zAwwl7BhFm0sHrnTOYYk3GgSjfuoUIghjfbvt8wMMJQ=="
                        }
                    ]
                }
            }


## 查询事务信息 [/api/v1/chain/transactions/{transactionID}]

### 查询事务信息 [GET]

- 响应参数说明

    - result: 事务信息    
        - chaincodeID: chaincode 名称
        - txid: 事务 ID
        - timestamp: 事务时间戳

+ Parameters
    + transactionID: `ec239f5e-06ff-497a-96e2-d3ee9d266867` (string, required) 事务 ID

+ Response 200 (application/json;charset=utf-8)

    + Body

            {
                "jsonrpc": "2.0",
                "result": {
                    "chaincodeID": "uctest01",
                    "txid": "c730bd57-a449-49ef-add0-eeccb8ecd627",
                    "timestamp": {
                        "seconds": 1514174583,
                        "nanos": 209469311
                    }
                }
            }


## 解析一个事务 [/api/v1/chain/parseTx]

### 解析事务 [POST]

- 请求参数说明
    - tx: 提交的事务内容，通常是生成的待签名事务

- 响应参数说明

    - result: 被解析的事务信息，具有和查询链上事务信息类似的数据结构，其中TxID将被显示为Unknown

+ Request (application/x-www-form-urlencoded;charset=utf-8)

    + Body

        tx=I::MTOKEN.TRANSFER:ChoKB0FCQ0hBSU4SD0F0b21pY0VuZXJneV92MRIGCITtq+IFGhS4tudQ90DdEcgcAt/AFKKXgs6QoQ==:CgRFT1NBGj0KC6VvpbmQGaXIAAAAEhYKFBJtq6Q46oTnxDwvVqMgDtZeNxs7GhYKFC6mOHOU3PKLpmHXinFgxK+3GhDb

+ Response 200 (application/json;charset=utf-8)

    + Body

        {
            "jsonrpc": "2.0",
            "result": {
                "Height": "0",
                "TxID": "Unknown",
                "Chaincode": "Unknown",
                "Method": "MTOKEN.TRANSFER",
                "CreatedFlag": false,
                "ChaincodeModule": "AtomicEnergy_v1",
                "Nonce": "B8B6E750F740DD11C81C02DFC014A29782CE90A1",
                "Detail": {
                    "amount": "200000000000000000000000000",
                    "from": "AS6mOHOU3PKLpmHXinFgxK-3GhDb9YuC2g",
                    "to": "ARJtq6Q46oTnxDwvVqMgDtZeNxs7Ybt81A",
                    "token": "EOSA"
                },
                "TxHash": "B3C9C3C99BC7DFF9466DDD4D494EC120AECAEF0EF6798583C5667D77DD7324DD"
            }
        }


## 查询地址信息 [/api/v1/address/{address}]

### 查询地址信息 [GET]

- 响应参数说明

    - result: 
        - balance: 地址余额        
        - lastFundID: 最近的一次 fund ID

+ Parameters
    + address: `AWVQUuLC66BMT71kJeY11wDnDdbtltZNUA` (string, required) - 查询地址

+ Response 200 (application/json;charset=utf-8)

    + Body


            {
                "jsonrpc": "2.0",
                "result": {
                    "balance": 9000,
                    "lastFundID": "c52cad8b-6aa3-4001-88df-7da9e22e4526"
                }
            }
