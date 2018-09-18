<template>
    <div class="members">
        <Table border :columns="columns" :data="list"></Table>
    </div>
</template>

<script>
export default {
    data(){
        return {
            list:[],
            columns:[{
                        title: 'ID',
                        key: 'ID'
                    },
                    {
                        title: 'Name',
                        key: 'name'
                    },
                    {
                        title: 'Role',
                        key: 'role'
                    },
                    {
                        title: 'Status',
                        key: 'status',
                        render: (h, params) => {
                            return h("Button", {
                                props:{
                                    type:params.row.status == 'healthy' ? 'success' :'error',
                                    size:'small'
                                }
                            }, params.row.status)
                        }
                    },
                    {
                        title: 'DB Size',
                        key: 'db_size'
                    },
                    {
                        title: 'PeerURLs',
                        key: 'peerURLs'
                    },
                    {
                        title: 'ClientURLs',
                        key: 'clientURLs'
                    }]
        }
    },
    mounted(){
        this.getList();
    },
    methods:{
        getList(){
            this.$http.get('/members').then(resp=>{
                if(resp.data.err == ''){
                    this.list = resp.data.result;
                }else{
                    this.$Message.error(resp.data.err);
                }
            });
        }
    }
}
</script>

<style scoped>

</style>
