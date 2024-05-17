<template>
    <lay-container>
        <lay-row space="20">
            <lay-col md="12">
                <lay-row  space="10">
                    <lay-col md="18">
                        <lay-input v-model="port" type="number"></lay-input>
                    </lay-col>
                    <lay-col md="4">
                        <lay-ripple>
                            <lay-button :type="status?'danger':'primary'" @click="start" >{{ status?"停止":"启动" }}</lay-button>
                        </lay-ripple>
                    </lay-col>
                </lay-row>
                <div class="msg">
                    <div v-for="msg in messages" style="margin-top:5px;">
                        <div class="msg-user">{{getShortName(msg.name)}} : </div>
                        <div class="msg-data">{{msg.msg}}</div>
                    </div>
                </div>
            </lay-col>
            <lay-col md="12" class="right">
                <div>
                    <h3>在线用户</h3>
                    <div class="online">
                        <div v-for="item in users" :key="item" :style="sendUsers.includes(item)?selected:unselected"  @click="changeSendUsers(item)">{{getShortName(item)}}</div>
                    </div>
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
import {reactive, ref} from 'vue'
    import {Start, Stop, Status, SendMsg, All} from "../../wailsjs/go/main/App.js";
    import {EventsOn, LogPrint} from "../../wailsjs/runtime/runtime.js";
    const port=ref(9527)
    const status=ref(false)
    const sendData=ref("")
    const users=ref([])
    const sendUsers=ref([])
    const messages=ref([])
    const selected=ref({border:"1px solid red"})
    const unselected=ref({border:"1px solid #d1d1d1"})
    const debounce=(fn,delay)=>{
        let timer=null
        return ()=>{
            if(timer) clearTimeout(timer)
            timer=setTimeout(()=>{
                fn()
            },delay)
        }

    }
    EventsOn("server_start",()=>{
        status.value=true
    })
    EventsOn("server_stop",()=>{
        status.value=false
    })
    EventsOn("ws_conn",()=>{
         All().then(all=>{
             users.value=all
         })
    })
    EventsOn("ws_message",(name,msg)=>{
        messages.value.push({name,msg})
    })
    EventsOn("ws_close",()=>{
        All().then(all=>{
            users.value=all
        })
    })
    const  start= ()=>{
        debounce(()=>{
            if(status.value){
                LogPrint("stop")
                Stop().then(()=>{
                    sendUsers.value=[]
                })
            }else{
                LogPrint("start")
                Start(parseInt(port.value))
            }
        },200)()


    }
    const send=()=>{
        let data=sendData.value.trim()
        if(data!==""){
            SendMsg(sendUsers.value,data).then(()=>{
                messages.value.push({name:"服务器",msg:data})
            })
        }
    }
    const changeSendUsers=(name)=>{
        let i=sendUsers.value.indexOf(name);
        if(i<0){
            //不存在
            sendUsers.value.push(name)
        }else{
            sendUsers.value.splice(i,1)
        }
    }

    const getShortName=(name)=>{
        return name.substring(0,5)
    }
</script>

<style scoped>
.right{
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    height: 90vh;

}
.msg{
    height: calc(85vh - 50px );
    padding: 8px;
    width: 100%;
    overflow: auto;
}
.online{
    height: 300px;
    overflow: auto;
    display: flex;
    flex-flow: row wrap;

}
.online div{
    height: 25px;
    line-height: 25px;
    width: 100px;
    text-align: center;
    border: 1px solid #d1d1d1;
    margin:20px;
    cursor: pointer;
}
.msg-user{
    color: #1e9fff;
}
.msg-data{
    background-color: #45984c;
    border-radius: 5px;
    padding: 10px 5px;
    margin: 5px;
    color: #ffffff;

}
</style>