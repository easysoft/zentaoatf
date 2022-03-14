<template>
  <div ref="topMenuCon" class="top-menu" style="width: 100%">
    <template v-for="(item, key) in menuData">
      <a-link
          :key="key"
          v-if="!item.hidden"
          :to="item.path"
          :class="{'active': belongTopMenu === item.path }"
          :id="pathToId(item.path)"
          class="top-menu-li"
      >
        {{t(item.title)}}
      </a-link>
    </template>
  </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, PropType, Ref, ref, watch} from "vue";
import {useRouter} from "vue-router";
import {useStore} from "vuex";
import IconSvg from "@/components/IconSvg/index";
import {ProjectData} from "@/store/project";
import ProjectCreateForm from "@/views/component/project/create.vue";
import {createProject} from "@/services/project";
import {hideMenu} from "@/utils/dom";
import {useI18n} from "vue-i18n";
import { DownOutlined } from '@ant-design/icons-vue';
import {RoutesDataItem} from "@/utils/routes";

interface RightTopMenuSetupData {
  t: (key: string | number) => string;
  pathToId: (val) => void
}

export default defineComponent({
  name: 'RightTopProject',
  components: {},
  props: {
    menuData: {
      type: Array as PropType<RoutesDataItem[]>,
      default: () => {
        return [];
      }
    },
    belongTopMenu: {
      type: String,
      default: ''
    },
  },

  setup(): RightTopMenuSetupData {
    const { t } = useI18n();

    const pathToId = (path) => {
      return path.replaceAll('/', 'menu-')
    }

    return {
      t,
      pathToId,
    }
  }
})
</script>

<style lang="less" scoped>
@import '../../../assets/css/global.less';

.top-menu {
  height: @headerHeight;
  line-height: @headerHeight;
  flex: 1;
  overflow: hidden;
  overflow-x: auto;

  .top-menu-li {
    display: inline-block;
    padding: 0 15px;
    height: @headerHeight;
    text-decoration: none;
    color: #FFFFFF;
    font-size: 15px;
    border-bottom: solid 3px transparent;
    &:hover,
    &.active {
      background-color: @menu-dark-bg-active;
      color: @menu-dark-highlight-color;
    }
  }

  .breadcrumb {
    line-height: @headerHeight;
  }
}
</style>