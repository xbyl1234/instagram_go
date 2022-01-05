// let FBSharedFramework = Module.getBaseAddress("FBSharedFramework")
// console.log(`FBSharedFramework : ${FBSharedFramework}`)
//
// //Instagram 190.0
// var left = FBSharedFramework.add(0x15A8B4);
// console.log(`before: ${hexdump(left, {length: 8, ansi: true})}`);
// let maxPatchSize = 64;
// Memory.patchCode(left, maxPatchSize, function (code) {
//     let cw = new Arm64Writer(code, {pc: left});
//     cw.putBytes([0x80, 0x08, 0x00, 0x54]); //b.eq #0x150
//     cw.flush()
// });
// console.log(`before: ${hexdump(left, {length: 8, ansi: true})}`);
//
// function printBacktrace(context) {
//     console.log('called from:\n')
//     var bt = Thread.backtrace(context, Backtracer.ACCURATE)
//     for (var i = 0; i < bt.length; i++) {
//         console.log(bt[i] - FBSharedFramework + "  " + DebugSymbol.fromAddress(bt[i]))
//     }
// }
var oldCode = 0

function pass_ins_sslpinning() {
    var module = Process.findModuleByName("FBSharedFramework")
    //190.0
    //191.0
    //0x81 -> 0x80
    var pattern1 = "9f 06 00 71 ?? ?? ?? ?? e8 ?? 02 91 00 81 00 91"
    //192.0
    //193.0
    //0x21 -> 0x20
    //202.0
    //0x41 -> 0x40
    console.log("pattern1")
    var result = Memory.scanSync(module.base, module.size, pattern1)
    if (result.length === 1) {
        var addr = result[0].address.add(4)
        console.log("find at: " + addr)
        if (oldCode === 0) {
            oldCode = addr.readU8()
            console.log("code is " + oldCode)
        }

        console.log(`before: ${hexdump(addr, {length: 8, ansi: true})}`);

        Memory.patchCode(addr, 64, function (code) {
            let cw = new Arm64Writer(code, {pc: addr});
            cw.putBytes([oldCode - 1,
                addr.add(1).readU8(),
                addr.add(2).readU8(),
                addr.add(3).readU8()]);

            cw.flush()

        });

        console.log(`after: ${hexdump(addr, {length: 8, ansi: true})}`);
    } else {
        console.log("pattern error len:" + result.length)
    }
}

// pass_ins_sslpinning()
// frida -U -n Instagram -l agent/inst.js --no-pause