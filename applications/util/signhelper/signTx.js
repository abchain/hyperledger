const util = require('ethereumjs-util')
const readline = require('readline')

let secret = Buffer.from('0C28FCA386C7A227600B2FE50B7CAE11EC86D3BF1FBE471BE89827E19D72AA1D', 'hex')

const prompt = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

let msg = new Promise((resolve, reject) => {
        
    prompt.once('SIGINT', () => {reject('User reject')})
    prompt.question('Input hash:', (msg) => {

        try{
            resolve(msg)
        }catch(e){
            reject(e)
        }  
    });

})

Promise.resolve(msg).then(msg => {
    let msgb = Buffer.from(msg, 'hex')
    console.log('Sign this hash:', msgb)
    return {s: util.ecsign(msgb, secret), hash: msgb}})
.then(
    sign =>{
        let {s:{r, s, v}, hash} = sign
        console.log('Signature [r,s,v]:', r.toString('hex'), s.toString('hex'), v)

        //this is what we should do to grasp data from sign and send it to SDK
        let pkbuf = util.ecrecover(hash, v, r, s)
        console.log('public key buffer:', pkbuf.toString('hex'))
    }
).catch(console.log).then(() => {prompt.close()})