const { ccclass, property } = cc._decorator;

@ccclass
export default class MainScene extends cc.Component {

    @property(cc.EditBox)
    editBox: cc.EditBox = null;

    private isConnect: boolean = false;
    private ws = null;

    clickSend() {
        if (!this.isConnect) {
            console.log("no connect");
            return;
        }
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            let text = this.editBox.string;
            let binaryData = binaryStringToBuffer(text)
            console.log(binaryData); // 输出：Uint8Array(13) [72, 101, 108, 108, 111, 44, 32, 119, 111, 114, 108, 100, 33]
            this.ws.send(binaryData);
        }
        else {
            console.log("WebSocket instance wasn't ready...");
        }
    }

    clickConnect() {
        this.ws = new WebSocket("ws://127.0.0.1:80");
        this.ws.onopen = (event) => {
            console.log("Send Text WS was opened.");
            this.isConnect = true;
        };
        this.ws.onmessage = (event) => {
            console.log("回包: " + event.data);
        };
        this.ws.onerror = (event) => {
            console.log("Send Text fired an error");
        };
        this.ws.onclose = (event) => {
            console.log("WebSocket instance closed.");
        };
    }
}

function binaryStringToBuffer(binaryString: string): ArrayBuffer {
    const buffer = new ArrayBuffer(binaryString.length / 8);
    const view = new DataView(buffer);
    for (let i = 0; i < binaryString.length; i += 8) {
        const byte = binaryString.substr(i, 8);
        view.setUint8(i / 8, parseInt(byte, 2));
    }
    return buffer;
}
