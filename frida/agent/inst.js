let FBSharedFramework = Module.getBaseAddress("FBSharedFramework")
console.log(`FBSharedFramework : ${FBSharedFramework}`)



var left = FBSharedFramework.add(0x15A8B4);
console.log(`before: ${hexdump(left, {length: 8, ansi: true})}`);
let maxPatchSize = 64;
Memory.patchCode(left, maxPatchSize, function (code) {
    let cw = new Arm64Writer(code, {pc: left});
    cw.putBytes([0x80, 0x08, 0x00, 0x54]); //b.eq #0x150
    cw.flush()
});
console.log(`before: ${hexdump(left, {length: 8, ansi: true})}`);

function printBacktrace(context) {
    console.log('called from:\n')
    var bt = Thread.backtrace(context, Backtracer.ACCURATE)
    for (var i = 0; i < bt.length; i++) {
        console.log(bt[i] - FBSharedFramework + "  " + DebugSymbol.fromAddress(bt[i]))
    }
}




// var verifyWithMetrics = FBSharedFramework.add(0x15AF8C);
// Interceptor.attach(verifyWithMetrics, {
//     onEnter: function (args) {
//         console.log("on verifyWithMetrics")
//     },
//     onLeave: function (ret) {
//         console.log("on verifyWithMetrics exit")
//         this.w20 = 1
//         return 1
//     }
// });

// var OpenSSLerror = FBSharedFramework.add(0X10635B8);
// Interceptor.attach(OpenSSLerror, {
//     onEnter: function (args) {
//         console.log("on OpenSSLerror")
//         printBacktrace(this.context)
//     },
//     onLeave: function (ret) {
//         console.log("on OpenSSLerror exit")
//         return 1
//     }
// });


// frida -U -n Instagram -l agent/inst.js --no-pause


// var Openssl_Write = FBSharedFramework.add(0x2B72E4);
// Interceptor.attach(Openssl_Write, {
//     onEnter: function (args) {
//         console.log("on Openssl_Write")
//         console.log(hexdump(args[1], {
//             length: args[2]
//         }))
//     },
//     onLeave: function (ret) {
//         console.log("on Openssl_Write exit")
//         return 1
//     }
// });
//
//
// var BIO_Write = FBSharedFramework.add(0x21BA64);
// Interceptor.attach(BIO_Write, {
//     onEnter: function (args) {
//         console.log("on BIO_Write")
//         console.log('RegisterNatives called from:\n' + Thread.backtrace(this.context, Backtracer.ACCURATE).map(DebugSymbol.fromAddress).join('\n') + '\n');
//         console.log(hexdump(args[1], {
//             length: args[2].toInt32()
//         }))
//     },
//     onLeave: function (ret) {
//         console.log("on BIO_Write exit")
//         return 1
//     }
// });
//
//
// Interceptor.attach(Module.findExportByName(null, "write"), {
//     onEnter: function (args) {
//         console.log("on write")
//         console.log('RegisterNatives called from:\n' + Thread.backtrace(this.context, Backtracer.ACCURATE).map(DebugSymbol.fromAddress).join('\n') + '\n');
//     },
//     onLeave: function (ret) {
//         console.log("on write exit")
//         return 1
//     }
// });
//
//
// Interceptor.attach(Module.findExportByName(null, "send"), {
//     onEnter: function (args) {
//         console.log("on send")
//         console.log('RegisterNatives called from:\n' + Thread.backtrace(this.context, Backtracer.ACCURATE).map(DebugSymbol.fromAddress).join('\n') + '\n');
//     },
//     onLeave: function (ret) {
//         console.log("on send exit")
//         return 1
//     }
// });


// var mds =Process.enumerateModules()
// for (var index=0; index< mds.length;index++){
//     console.log(mds[index].name)
// }
//
// console.log(Process.findModuleByName('libc.so'))
//
//
// Interceptor.attach(Module.findExportByName(null, "send"), {
//     onEnter: function (args) {
//         console.log("on send")
//         console.log('RegisterNatives called from:\n' + Thread.backtrace(this.context, Backtracer.ACCURATE).map(DebugSymbol.fromAddress).join('\n') + '\n');
//     },
//     onLeave: function (ret) {
//         console.log("on send exit")
//         return 1
//     }
// });
////
// var CheckSsl = FBSharedFramework.add(0x220A54);
// var CheckSsl_f = FBSharedFramework.add(0x21FFAC);
//
//
// Interceptor.attach(CheckSsl, {
//     onEnter: function (args) {
//         console.log("on check")
//     },
//     onLeave: function (ret) {
//         console.log("on check exit")
//         return 1
//     }
// });
//
// Interceptor.attach(CheckSsl_f, {
//     onEnter: function (args) {
//         console.log("on check_f")
//     },
//     onLeave: function (ret) {
//         console.log("on check_f exit")
//         return 1
//     }
// });


// b.ne #0x150
// 81 0A 00 54

//

// Cheat verifyWithMetrics from proxygen
