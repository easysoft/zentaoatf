<template>
    <div className="indexlayout-main-conent">
        <a-card :bordered="false">

             <div v-for="(item, index) in svgIcons" class="list" :key="index" >
                <a-popover>
                    <template #content>
                        &lt;icon-svg type="{{item}}" /&gt;
                    </template>
                    <div>
                        <icon-svg :type="item || ''" style="font-size: 30px" />
                        <span>{{item}}</span>
                    </div>
                </a-popover>
              </div>

              <a-divider />

                <a-list>
                    <template #header><h2>{{t('page.icon.svg.remark.title')}}</h2></template>
                    <a-list-item> 组件位置： @/components/IconSvg</a-list-item>
                    <a-list-item> 创建原因：方便自定义使用svg图标 </a-list-item>
                </a-list>
                <a-list>
                    <template #header><h2>使用方法：</h2></template>
                    <a-list-item>
                        1、下载或制作svg文件，存放到 <a-tag>@/assets/iconsvg</a-tag>
                        目录下，自己可以对此目录下svg进行删减。
                    </a-list-item>
                    <a-list-item>
                        2、项目会根据 <a-tag>@/assets/iconsvg/svgo.yml</a-tag>
                        配置自动压缩精简svg，也可以独立运行 <a-tag>yarn svgo</a-tag> 或
                        <a-tag>npm run svgo</a-tag>压缩精简svg
                    </a-list-item>
                    <a-list-item>3、使用Demo：</a-list-item>
                    <a-list-item>import IconSvg from '@/components/IconSvg';</a-list-item>
                    <a-list-item>
                        &lt;IconSvg type="svg文件名" class="" style=""/&gt;
                    </a-list-item>
                </a-list>

          
        </a-card>
    </div>
</template>
<script lang="ts">
import { defineComponent } from "vue";
import { useI18n } from "vue-i18n";
import IconSvg from "@/components/IconSvg";

const requireAll = (requireContext: any/* __WebpackModuleApi.RequireContext */) =>
  requireContext.keys();
const svgIcons = requireAll(
  require.context('../../../../assets/iconsvg', false, /\.svg$/),
).map(i => {
  const item = i.match(/\.\/(.*)\.svg/);
  return item && item[1];
});


export default defineComponent({
    components: {
        IconSvg
    },
    setup() {
        const { t } = useI18n();

        return {
            t,
            svgIcons
        }
    }
})
</script>
<style lang="less" scoped>
.list {
  padding: 10px 20px;
  width: 100px;
  height: 100px;
  float: left;
  text-align: center;
  font-size: 30px;
  overflow: hidden;
  span {
    display: block;
    font-size: 16px;
    margin-top: 10px;
  }
}
</style>