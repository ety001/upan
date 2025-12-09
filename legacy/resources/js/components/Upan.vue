<template>
    <div>
        <div class="card">
            <div class="card-header">获取文件</div>
            <div class="card-body">
                <form>
                    <div class="form-group">
                        <label for="code">下载码</label>
                        <input type="text" v-model="code" class="form-control">
                    </div>
                    <button type="button" @click="getFile" class="btn btn-primary">下载</button>
                </form>
            </div>
        </div>

        <div class="card" style="margin-top: 20px;">
            <div class="card-header">上传文件</div>
            <div class="card-body">
                <form method="post" action="/upload" enctype="multipart/form-data">
                    <div class="form-group">
                        <input type="file" class="form-control-file" name="o">
                        <small id="o" class="form-text text-muted">文件大小限制: {{ fileMaxSize }}MB, 保留时间: {{ fileExpireTime }}小时</small>
                    </div>
                    <button type="submit" class="btn btn-primary">上传</button>
                </form>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        data() {
            return {
                code: null,
            };
        },
        props: ['fileMaxSize', 'fileExpireTime'],
        mounted() {
            console.log('Component mounted.')
        },
        methods: {
            getFile() {
                if (this.code) {
                    window.location.href = '/file/' + this.code;
                }
            },
        },
    }
</script>
