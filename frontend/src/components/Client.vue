<template>
    <lay-container>
        <lay-row space="20">
            <lay-col md="12">
                <div class="msg">
                    <div v-for="msg in messages" style="margin-top:5px;">
                        <div class="msg-user">{{ msg.name }} :</div>
                        <div class="msg-data">{{ msg.msg }}</div>
                    </div>
                </div>
            </lay-col>
            <lay-col md="12" class="right">
                <div>
                    <lay-row space="10">
                        <lay-col md="16">
                            <div class="ws">ws://</div>
                            <lay-input v-model="addr" class="addr"></lay-input>
                        </lay-col>
                        <lay-col md="4">
                            <lay-ripple>
                                <lay-button :type="status?'danger':'primary'" @click="start">{{
                                        status ? "断开" : "连接"
                                    }}
                                </lay-button>
                            </lay-ripple>
                        </lay-col>
                    </lay-row>
                </div>
                <div>
                    <lay-textarea placeholder="请输入发送的内容" v-model="sendData" style="resize: none"></lay-textarea>
                    <div style="margin-top:10px;text-align: right;">
                        <lay-button type="normal" @click="send">发送</lay-button>
                    </div>
                </div>
               
            </lay-col>
        </lay-row>
    </lay-container>
</template>
<script setup>
import {ref} from 'vue'
import {LogPrint} from "../../wailsjs/runtime/runtime.js";
import {ShowErrDialog} from "../../wailsjs/go/main/App.js";

const addr = ref("127.0.0.1:9527/ws")
const status = ref(false)
const sendData = ref("")
const messages = ref([])
let ws;
const start = () => {
    if(ws===undefined){
        ws = new WebSocket("ws://" + addr.value)
        ws.onopen = () =>{
            status.value=true
            ws.send("hello"); //将消息发送到服务端
        }
        ws.onmessage =  (e) =>{
            //当客户端收到服务端发来的消息时，触发onmessage事件，参数e.data包含server传递过来的数据
            LogPrint(e.data);
            messages.value.unshift({name: "服务器", msg: e.data})
        }
        ws.onclose = (e) =>{
            //当客户端收到服务端发送的关闭连接请求时，触发onclose事件
            LogPrint("close");
            status.value=false
            ws=undefined
        }
        ws.onerror = (e) =>{
            LogPrint("error")
            LogPrint(e.message)
        }
    }else {
        LogPrint("123")
        ws.close()
        ws=undefined

    }


}
const send = () => {
    let data = sendData.value.trim()
    if (data !== "") {
        if(ws!==undefined){
            ws.send(data)
            messages.value.unshift({name: "我", msg: data})
        }else{
            ShowErrDialog("还没有连接到服务端")
        }
    }else{
        ShowErrDialog("消息内容不能为空")
    }
}

</script>

<style scoped>
.ws {
    width: 30px;
    padding: 11px 0;
    position: absolute;
    left: 15px;
}

.addr {
    padding-left: 30px;
}

.right {
    display: flex;
    flex-direction: column;
    /* justify-content: space-between; */
    height: 90vh;

}

.msg {
    height: calc(90vh - 50px);
    padding: 8px;
    width: 100%;
    overflow: auto;
}

.online {
    height: 300px;
    overflow: auto;
    display: flex;
    flex-flow: row wrap;

}

.online div {
    height: 25px;
    line-height: 25px;
    width: 100px;
    text-align: center;
    border: 1px solid #d1d1d1;
    margin: 20px;
    cursor: pointer;
}

.msg-user {
    color: #1e9fff;
}

.msg-data {
    background-color: #45984c;
    border-radius: 5px;
    padding: 10px 5px;
    margin: 5px;
    color: #ffffff;

}
</style>