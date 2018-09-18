<template>
    <div class="seting">
        <Form :model="user" :label-width="50" style="width:300px;">
        <FormItem prop="username">
            <Input type="text" v-model="user.username" placeholder="Username">
                <Icon type="ios-person-outline" slot="prepend"></Icon>
            </Input>
        </FormItem>
        <FormItem prop="password">
            <Input type="password" v-model="user.password" placeholder="Password">
                <Icon type="ios-lock-outline" slot="prepend"></Icon>
            </Input>
        </FormItem>
        <FormItem>
            <Button type="primary" @click="save()">保存</Button>
        </FormItem>
    </Form>
    </div>
</template>

<script>
export default {
    data(){
        return{
            user:{
                username:'',
                password:''
            }
        }
    },
    mounted(){
        let etcdUserStr = localStorage.getItem('etcdUser');
        if(etcdUserStr != null && typeof etcdUserStr != 'undefined'){
            this.user = JSON.parse(etcdUserStr);
        }
        
    },
    methods:{
        save(){
            // console.log(this.user);
            localStorage.setItem('etcdUser', JSON.stringify(this.user));
            // console.log(localStorage.getItem('etcdUser'))
        }
    }
}
</script>

<style scoped>
.seting{
    width: 100%;
    min-height: 500px;
}
</style>
