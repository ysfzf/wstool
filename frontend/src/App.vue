<template>
    <div class="content" :style="currentHeight">
        <lay-tab v-model="index" type="brief">
            <lay-tab-item title="客户端" v-if="true" id="1">
                 <client />
            </lay-tab-item>
            <lay-tab-item title="服务端" id="2">
                <service />
            </lay-tab-item>
        </lay-tab>
    </div>

</template>

<script setup>
import Service from "./components/Service.vue";
import {ref} from 'vue'
import {LogPrint, ScreenGetAll, WindowGetSize} from "../wailsjs/runtime/runtime.js";
import Client from "./components/Client.vue";

const index=ref("1")
const currentHeight=ref({height:800})
window.onresize=()=>{
    WindowGetSize().then(r=>{
        LogPrint(r.w+" "+r.h)
        currentHeight.value={height:r.h}
    })
}
</script>
<style>

</style>
