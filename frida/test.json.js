import util from 'ethereumjs-util'

const text = "62309d1c9bd7c26cd07cdbda"
const msg = Buffer.from(text)
const sig = "0xd8e3c1cde5c053bf916825eb9cd4531006d596495004be1ad1ceab4555e3949344757dcc7b33fec101ee4965ffc7164657b5b93ad60955fb60f76cbfb56ca24b01"
const res = util.fromRpcSig(sig)
const prefix = Buffer.from("\x19Ethereum Signed Message:\n");
const prefixedMsg = util.keccak256(
    Buffer.concat([prefix, Buffer.from(String(msg.length)), msg])
);

var a = Buffer.concat([prefix, Buffer.from(String(msg.length)), msg]);

console.log(Buffer.concat([prefix, Buffer.from(String(msg.length)), msg]).toString())
const pubKey = util.ecrecover(prefixedMsg, res.v, res.r, res.s);
const addrBuf = util.pubToAddress(pubKey);
const addr = util.bufferToHex(addrBuf);

console.log(addr)

// 0xe33bc853bc59d7b0bb725a5fd86ce166e7fff7be