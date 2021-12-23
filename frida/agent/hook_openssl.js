// EVP_EncryptUpdate
//B3AC8
//EVP_EncryptFinal_ex
//B4104
//EVP_DecryptUpdate
//B1B54
//EVP_DecryptFinal_ex
//B25B0
// Convert a hex string to a byte array
function hexToBytes(hex) {
    for (var bytes = [], c = 0; c < hex.length; c += 2)
        bytes.push(parseInt(hex.substr(c, 2), 16));
    return bytes;
}

// Convert a ASCII string to a hex string
function stringToHex(str) {
    return str.split("").map(function (c) {
        return ("0" + c.charCodeAt(0).toString(16)).slice(-2);
    }).join("");
}

// Convert a hex string to a ASCII string
function hexToString(hexStr) {
    var hex = hexStr.toString();//force conversion
    var str = '';
    for (var i = 0; i < hex.length; i += 2)
        str += String.fromCharCode(parseInt(hex.substr(i, 2), 16));
    return str;
}

// Convert a byte array to a hex string
function bytesToHex(pointer, len) {
    // console.log("point: " + pointer + " len: " + len)
    if (len === 0 || pointer.toInt32() === 0) {
        return ""
    }
    pointer = new NativePointer(pointer)
    for (var hex = [], i = 0; i < len; i++) {
        // console.log("+++++++++" + pointer.add(i).readU8())
        var ch = pointer.add(i).readU8()
        hex.push((ch >>> 4).toString(16));
        hex.push((ch & 0xF).toString(16));
    }
    return hex.join("");
}

let FBSharedFramework = Module.getBaseAddress("FBSharedFramework")
console.log(`FBSharedFramework : ${FBSharedFramework}`)

function printBacktrace(context) {
    console.log('called from:\n')
    var bt = Thread.backtrace(context, Backtracer.ACCURATE)
    for (var i = 0; i < bt.length; i++) {
        console.log(bt[i] - FBSharedFramework + "  " + DebugSymbol.fromAddress(bt[i]))
    }
}


function HookOpensslEvp() {
    var EVP_EncryptUpdate = FBSharedFramework.add(0xB3AC8);
    var EVP_EncryptFinal_ex = FBSharedFramework.add(0xB4104);
    var EVP_DecryptUpdate = FBSharedFramework.add(0xB1B54);
    var EVP_DecryptFinal_ex = FBSharedFramework.add(0xB25B0);


    //encode
    Interceptor.attach(EVP_EncryptUpdate, {
        onEnter: function (args) {
            // console.log("on EVP_EncryptUpdate")
            // printBacktrace(this.context)
            console.log("context: " + args[0] + " send1: " + bytesToHex(args[3], args[4].toInt32()))
        },
        onLeave: function (ret) {
            // console.log("on EVP_EncryptUpdate exit")

        }
    });

    Interceptor.attach(EVP_EncryptFinal_ex, {
        onEnter: function (args) {
            // console.log("on EVP_EncryptFinal_ex")
            console.log("context: " + args[0] + " send2: " + bytesToHex(args[1], args[2].readInt()))
        },
        onLeave: function (ret) {
            // console.log("on EVP_EncryptFinal_ex exit")
        }
    });

    //decode
    Interceptor.attach(EVP_DecryptUpdate, {
        onEnter: function (args) {
            // console.log("on EVP_DecryptUpdate")
            // console.log('RegisterNatives called from:\n' + Thread.backtrace(this.context, Backtracer.ACCURATE).map(DebugSymbol.fromAddress).join('\n') + '\n');

            this.a0 = args[0]
            this.a1 = args[1]
            this.a2 = args[2]
        },
        onLeave: function (ret) {
            if (this.a1.toInt32() !== 0) {
                console.log("context: " + this.a0 + " recv1: " + bytesToHex(this.a1, this.a2.readInt()))
            }
            // console.log("on EVP_DecryptUpdate exit")
        }
    });

    Interceptor.attach(EVP_DecryptFinal_ex, {
        onEnter: function (args) {
            // console.log("on EVP_DecryptFinal_ex")
            this.a0 = args[0]
            this.a1 = args[1]
            this.a2 = args[2]
        },
        onLeave: function (ret) {
            console.log("context: " + this.a0 + " recv2: " + bytesToHex(this.a1, this.a2.readInt()))
            // console.log("on EVP_DecryptFinal_ex exit")
        }
    });
}

function logp() {
    console.log(`before: ${hexdump(FBSharedFramework.add(0x02A9B38), {length: 8, ansi: true})}`);
}

logp()

function HookOpenssl() {
    console.log("HookOpenssl")
    var SSL_SESSION_new = FBSharedFramework.add(0x02A9B38);
    Interceptor.attach(SSL_SESSION_new, {
        onEnter: function (args) {
            // console.log("on SSL_SESSION_new")
            console.log("SSL_SESSION_new  ")
            // console.log('RegisterNatives called from:\n' + Thread.backtrace(this.context, Backtracer.ACCURATE).map(DebugSymbol.fromAddress).join('\n') + '\n');
        },
        onLeave: function (ret) {
            // console.log("on SSL_SESSION_new exit")
            //
        }
    });

    var d2i_SSL_SESSION = FBSharedFramework.add(0xC72D00);
    Interceptor.attach(d2i_SSL_SESSION, {
        onEnter: function (args) {
            // console.log("on d2i_SSL_SESSION")
            console.log("d2i_SSL_SESSION  ")
            // console.log('RegisterNatives called from:\n' + Thread.backtrace(this.context, Backtracer.ACCURATE).map(DebugSymbol.fromAddress).join('\n') + '\n');
        },
        onLeave: function (ret) {
            // console.log("on d2i_SSL_SESSION exit")
            //
        }
    });

    var getSessionFromCacheData = FBSharedFramework.add(0x948D10);
    Interceptor.attach(getSessionFromCacheData, {
        onEnter: function (args) {
            // console.log("on getSessionFromCacheData")
            console.log("getSessionFromCacheData send: ")
            // console.log('RegisterNatives called from:\n' + Thread.backtrace(this.context, Backtracer.ACCURATE).map(DebugSymbol.fromAddress).join('\n') + '\n');
        },
        onLeave: function (ret) {
            // console.log("on getSessionFromCacheData exit")
            //
        }
    });

    var ssl_new = FBSharedFramework.add(0x2A5EC0);
    Interceptor.attach(ssl_new, {
        onEnter: function (args) {
            // console.log("on ssl_new")
            console.log("ssl_new send: ")
            // console.log('RegisterNatives called from:\n' + Thread.backtrace(this.context, Backtracer.ACCURATE).map(DebugSymbol.fromAddress).join('\n') + '\n');
        },
        onLeave: function (ret) {
            // console.log("on ssl_new exit")

        }
    });

    var ssl_write = FBSharedFramework.add(0x2B72E4);
    Interceptor.attach(ssl_write, {
        onEnter: function (args) {
            // console.log("on ssl_write")
            console.log("ssl_write send: ")
            // console.log('RegisterNatives called from:\n' + Thread.backtrace(this.context, Backtracer.ACCURATE).map(DebugSymbol.fromAddress).join('\n') + '\n');
        },
        onLeave: function (ret) {
            // console.log("on ssl_write exit")

        }
    });

    var ssl_connect = FBSharedFramework.add(0x2A7D5C);
    Interceptor.attach(ssl_connect, {
        onEnter: function (args) {
            // console.log("on ssl_connect")
            console.log("ssl_connect : ")
            console.log('RegisterNatives called from:\n' + Thread.backtrace(this.context, Backtracer.ACCURATE).map(DebugSymbol.fromAddress).join('\n') + '\n');
        },
        onLeave: function (ret) {
            // console.log("on ssl_connect exit")

        }
    });

    var asslsock_connectSuccess = FBSharedFramework.add(0x0099D70C);
    Interceptor.attach(asslsock_connectSuccess, {
        onEnter: function (args) {
            // console.log("on ssl_connect")
            console.log("asslsock_connectSuccess : ")
        },
        onLeave: function (ret) {
            // console.log("on ssl_connect exit")
            //
        }
    });
}

HookOpensslEvp()
HookOpenssl()
// frida -U -n Instagram -l agent/hook_openssl.js --no-pause

















