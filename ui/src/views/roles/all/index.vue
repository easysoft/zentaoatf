<template>
    <div className="indexlayout-main-conent">
        <a-card :bordered="false">

                <a-alert message="此栏目为权限功能做演示！" type="warning" />


                <a-descriptions
                    title="以下有三个账号,密码：123456，你可以登录不同的账号查看此栏目下的区别："
                    layout="vertical"
                    bordered
                    style="margin-top: 20px"
                    size="small"
                    >
                    <a-descriptions-item label="admin">
                        超级管理员（拥有所有权限）
                    </a-descriptions-item>
                    <a-descriptions-item label="user">临时用户账号</a-descriptions-item>
                    <a-descriptions-item label="test">临时测试账号</a-descriptions-item>
                </a-descriptions>

                <a-alert
                    message="注意：退出账号后登录其他账号，请刷新在查看效果。"
                    type="error"
                    style="margin-top: 20px"
                />

                <a-descriptions
                    title="此页面所有用户都可以查看，以下只做操作按钮权限演示，页面权限请登录不同账号查看左侧栏目对应菜单区别。"
                    layout="vertical"
                    bordered
                    :column="1"
                    style="margin-top: 20px"
                    size="small"
                    >
                    <a-descriptions-item label="不做验证的操作按钮">
                        <a-button type="primary">编辑</a-button>
                        <a-button type="primary" danger> 删除</a-button>
                    </a-descriptions-item>

                    <a-divider><h3>Permission 组件使用方法：</h3></a-divider>

                    <a-descriptions-item label="user账号可操作按钮">
                        <Permission roles="user">
                            <template #otherwise>无权操作，此参数可赋值为空！</template>
                            <a-button type="primary">编辑</a-button>
                            <a-button type="primary" danger>
                                删除
                            </a-button>
                        </Permission>
                    </a-descriptions-item>
                    <a-descriptions-item label="test账号可操作按钮">
                        <Permission roles="test">
                            <template #otherwise>无权操作，此参数可赋值为空！</template>
                            <a-button type="primary">编辑</a-button>
                            <a-button type="primary" danger>
                                删除
                            </a-button>
                        </Permission>
                    </a-descriptions-item>

                    <a-divider><h3>v-permission 指令使用方法：</h3></a-divider>

                    <a-descriptions-item label="user账号可操作按钮">
                            <a-button v-permission="'user'" type="primary">编辑</a-button>
                            <a-button v-permission="'user'" type="primary" danger>删除</a-button>
                    </a-descriptions-item>
                    <a-descriptions-item label="test账号可操作按钮">
                            <a-button v-permission="'test'" type="primary">编辑</a-button>
                            <a-button v-permission="'test'" type="primary" danger>删除 </a-button>
                    </a-descriptions-item>


                </a-descriptions>

          
        </a-card>
    </div>
</template>
<script lang="ts">
import { defineComponent } from "vue";
import Permission from "@/components/Permission/index.vue";
import vPermission from '@/directives/permission';
export default defineComponent({
    name: 'RolesAll',
    components: {
        Permission
    },
    directives: {
        permission: vPermission
    }
})
</script>