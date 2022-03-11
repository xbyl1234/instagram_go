function hook() {
    let oldXHR = window.XMLHttpRequest;


    function ajaxEventTrigger(event) {
        var ajaxEvent = new CustomEvent(event, {
            detail: this
        })
        // console.log("this", this)
        // console.log("url", this.responseURL)
        // console.log("u", this.responseURL === "https://www.barotem.com/auth/getSms")

        if (this.responseURL === "https://www.barotem.com/auth/getSms") {
            console.log("sms", this.sms)
        }
        window.dispatchEvent(ajaxEvent)
    }

    function newXHR() {
        let realXHR = new oldXHR();
        let oldSendFun = realXHR.send;

        realXHR.send = function (body) {

            this.sms = body
            oldSendFun.call(realXHR, body)
        }

        realXHR.addEventListener('loadend', function () {
            ajaxEventTrigger.call(this, 'ajaxLoadEnd')
        }, false)

        return realXHR;
    }

    window.XMLHttpRequest = newXHR;
}

hook()
