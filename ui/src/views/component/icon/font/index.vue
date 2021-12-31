<template>
    <div className="indexlayout-main-conent">
        <a-card :bordered="false">

                <a-alert
                    message="使用此组件需要配置 ‘@/config/settings.ts’ 中 ‘iconfontUrl’ 参数。"
                    type="warning"
                />

                <a-list>
                    <template #header><h2>说明：</h2></template>
                    <a-list-item> 组件位置： @/components/IconFont</a-list-item>
                    <a-list-item> 创建原因：方便统一管理使用 iconfont.cn 上 js 资源图标 </a-list-item>
                </a-list>
                <a-list>
                    <template #header><h2>使用方法：</h2></template>
                    <a-list-item>
                        1、iconfont.cn 上生成 js 资源。
                    </a-list-item>
                    <a-list-item>
                        2、@/config/settings.ts 文件中配置 iconfont.cn 生成的js文件地址。
                    </a-list-item>                   
                    <a-list-item>3、使用Demo：</a-list-item>
                    <a-list-item>import IconFont from '@/components/IconFont';</a-list-item>
                    <a-list-item>
                        &lt;IconFont type="svg文件名" class="" style=""/&gt;
                    </a-list-item>
                </a-list>

          
        </a-card>
    </div>
</template>