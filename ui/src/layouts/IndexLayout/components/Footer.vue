<template>
    <div class="indexlayout-right-footer">
      <div class="footer-links">
        <a href="https://ztf.im/ztf.html" target="_blank">
          {{t('download')}}
        </a>
        <a href="https://ztf.im/book/ztf/ztf-about-26.html" target="_blank">
          {{t('help')}}
        </a>
        <a href="https://github.com/easysoft/ztf" target="_blank">
          {{t('src')}}
        </a>
      </div>
      <div>Copyright © 2019-2022 青岛易软天创, All Rights Reserved</div>
    </div>
</template>
<script lang="ts">
import { defineComponent } from "vue";
export default defineComponent({
    name: 'RightFooter',
    components: {
    }
})
</script>
<style lang="less" scoped>
.indexlayout-right-footer {
  height: 32px;
  line-height: 32px;
  overflow: hidden;
  text-align: center;
  font-size: 14px;
  color: #808695;
  .footer-links {
    margin-bottom: 10px;
    a {
      margin-right: 30px;
      color: #808695;
      text-decoration: none;
      &:last-child {
        margin-right: 0;
      }
      &:hover {
        color: #515a6e;
      }
    }
  }
}
</style>