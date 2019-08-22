## TxDumper 

Help dumping out all of the transactions within a chain (currently it just raise the client implement for YA-fabric) to a bunch of text lines written by JSON objects, which respect the **ctormsg** format in *peer chaincode invoke/query* command used in fabric (0.6 or YA-fabric)

Notice YA-fabric has provided a new *peer* implement which can read ctormsg arguments from stdin so you can tunnel them for dumping and replaying txs to another chain

Use --help command to see usage and avaliable options