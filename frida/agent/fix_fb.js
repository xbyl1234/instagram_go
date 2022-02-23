let FBSharedFramework = Module.getBaseAddress("FBSharedFramework")
console.log(`FBSharedFramework : ${FBSharedFramework}`)

var verifyWithMetrics = FBSharedFramework.add(0x224ECC);
Interceptor.attach(verifyWithMetrics, {
    onEnter: function (args) {
        console.log("on verifyWithMetrics")
    },
    onLeave: function (ret) {
        console.log("on verifyWithMetrics exit")
        return 1
    }
});

// Cheat FBLiger protection settings to disable fizz and SSL cache

var resetSettings = [
    "persistentSSLCacheEnabled",
    "crossDomainSSLCacheEnabled",
    "fizzEnabled",
    "fizzPersistentCacheEnabled",
    "quicFizzEarlyDataEnabled",
    "fizzEarlyDataEnabled",
    "enableFizzCertCompression"
];

function cheatFBLigerSettings() {
    var resolver = new ApiResolver('objc');
    for (var i = 0; i < resetSettings.length; i++) {
        var matches = resolver.enumerateMatchesSync("-[FBLigerConfig " + resetSettings[i] + "]");
        if (matches.length < 1) {
            console.log("[w] Failed to reset " + resetSettings[i] + ", address not found!");
            continue;
        }
        Interceptor.attach(matches[0]["address"], {
            onLeave: function (retval) {
                console.log("[i] -[FBLigerConfig *] called!");
                retval.replace(0);

            }
        });
        console.log("[i] -[FBLIgerConfig " + resetSettings[i] + "] reset!")
    }
}

cheatFBLigerSettings();
// frida -U -n Instagram -l agent/fix_fb.js --no-pause


// Cheat cerificate verification callcbacks from boringssl and FBSharedFramework
// function cheatCallbacks() {
//     var SSL_CTX_sess_set_new_cb_addr = DebugSymbol.findFunctionsNamed("SSL_CTX_sess_set_new_cb");
//     var SSL_CTX_set_cert_verify_callback_addr = DebugSymbol.findFunctionsNamed("SSL_CTX_set_cert_verify_callback");
//     var SSL_CTX_set_cert_verify_result_callback_addr = DebugSymbol.findFunctionsNamed("SSL_CTX_set_cert_verify_result_callback");
//     var SSL_CTX_set_verify_addr = DebugSymbol.findFunctionsNamed("SSL_CTX_set_verify");
//     var SSL_set_verify_addr = DebugSymbol.findFunctionsNamed("SSL_set_verify");
//     var SSL_set_cert_cb_addr = DebugSymbol.findFunctionsNamed("SSL_set_cert_cb");
//     var SSL_CTX_set_cert_cb_addr = DebugSymbol.findFunctionsNamed("SSL_CTX_set_cert_cb");
//     var X509_STORE_CTX_set_verify_cb_addr = DebugSymbol.findFunctionsNamed("X509_STORE_CTX_set_verify_cb");
//
//
//     for(var i = 0; i < SSL_CTX_set_cert_verify_callback_addr.length; i++) {
//         Interceptor.replace(SSL_CTX_set_cert_verify_callback_addr[i], new NativeCallback(function () {
//             console.log("[i] SSL_CTX_set_cert_verify_callback(...) called!");
//             return;
//         }, 'void', []));
//     }
//     console.log("[i] SSL_CTX_set_cert_verify_callback(...) hooked!");
//
//     for(var i = 0; i < SSL_CTX_set_cert_verify_result_callback_addr.length; i++) {
//         Interceptor.replace(SSL_CTX_set_cert_verify_result_callback_addr[i], new NativeCallback(function () {
//             console.log("[i] SSL_CTX_set_cert_verify_result_callback(...) called!");
//             return;
//         }, 'void', []));
//     }
//     console.log("[i] SSL_CTX_set_cert_verify_result_callback(...) hooked!");
//
//     for(var i = 0; i < SSL_CTX_set_verify_addr.length; i++) {
//         Interceptor.replace(SSL_CTX_set_verify_addr[i], new NativeCallback(function () {
//             console.log("[i] SSL_CTX_set_verify(...) called!");
//             return;
//         }, 'void', []));
//     }
//     console.log("[i] SSL_CTX_set_verify(...) hooked!");
//
//     for(var i = 0; i < SSL_set_verify_addr.length; i++) {
//         Interceptor.replace(SSL_set_verify_addr[i], new NativeCallback(function () {
//             console.log("[i] SSL_set_verify(...) called!");
//             return;
//         }, 'void', []));
//     }
//     console.log("[i] SSL_set_verify(...) hooked!");
//
//     for(var i = 0; i < SSL_set_cert_cb_addr.length; i++) {
//         Interceptor.replace(SSL_set_cert_cb_addr[i], new NativeCallback(function () {
//             console.log("[i] SSL_set_cert_cb(...) called!");
//             return;
//         }, 'void', []));
//     }
//     console.log("[i] SSL_set_cert_cb(...) hooked!");
//
//     for(var i = 0; i < SSL_CTX_set_cert_cb_addr.length; i++) {
//         Interceptor.replace(SSL_CTX_set_cert_cb_addr[i], new NativeCallback(function () {
//             console.log("[i] SSL_CTX_set_cert_cb(...) called!");
//             return;
//         }, 'void', []));
//     }
//     console.log("[i] SSL_CTX_set_cert_cb(...) hooked!");
//
//     for(var i = 0; i < X509_STORE_CTX_set_verify_cb_addr.length; i++) {
//         Interceptor.replace(X509_STORE_CTX_set_verify_cb_addr[i], new NativeCallback(function () {
//             console.log("[i] X509_STORE_CTX_set_verify_cb(...) called!");
//             return;
//         }, 'void', []));
//     }
//     console.log("[i] X509_STORE_CTX_set_verify_cb(...) hooked!");
// }
// cheatCallbacks();
//
// // Cheat SecTrustEvaluate, just in case :)
//
// function cheatSecTrustEvaluate() {
//     var SecTrustEvaluate_prt = Module.findExportByName("Security", "SecTrustEvaluate");
//     if (SecTrustEvaluate_prt == null) {
//         console.log("[e] Security!SecTrustEvaluate(...) not found!");
//         return;
//     }
//     var SecTrustEvaluate = new NativeFunction(SecTrustEvaluate_prt, "int", ["pointer", "pointer"]);
//     Interceptor.replace(SecTrustEvaluate_prt, new NativeCallback(function(trust, result) {
//         console.log("[i] SecTrustEvaluate(...) called!");
//         var osstatus = SecTrustEvaluate(trust, result);
//         Memory.writeU8(result, 1);
//         return 0;
//     }, "int", ["pointer", "pointer"]));
//     console.log("[i] SecTrustEvaluate(...) hooked!");
// }
// cheatSecTrustEvaluate();