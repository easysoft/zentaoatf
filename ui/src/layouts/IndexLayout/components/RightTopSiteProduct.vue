<template>
  <div>
    <!-- zentao site selection -->
    <a-dropdown
        :dropdownMatchSelectWidth="false"
        class="dropdown-list">

      <a class="t-link-btn" @click.prevent>
        <span class="name">{{currSite.name}}</span>
        <span class="icon2"><icon-svg type="down"></icon-svg></span>
      </a>
      <template #overlay>
        <a-menu class="menu">
          <template v-for="item in sites" :key="item.path">
            <a-menu-item v-if="currSite.path !== item.path">
                <div class="line">
                  <div class="t-link name" @click="selectSite(item)">{{ item.name }}</div>
                </div>
            </a-menu-item>
          </template>
        </a-menu>
      </template>
    </a-dropdown>

    <!-- zentao product selection -->
    <a-dropdown
        :dropdownMatchSelectWidth="false"
        class="dropdown-list">

      <a class="t-link-btn" @click.prevent>
        <span class="name">{{currProduct.name}}</span>
        <span class="icon2"><icon-svg type="down"></icon-svg></span>
      </a>
      <template #overlay>
        <a-menu class="menu">
          <template v-for="item in workspaces" :key="item.path">
            <a-menu-item v-if="currProduct.path !== item.path">
              <div class="line">
                <div class="t-link name" @click="selectProduct(item)">{{ item.name }}</div>
              </div>
            </a-menu-item>
          </template>
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, Ref, ref, watch} from "vue";
import {useRouter} from "vue-router";
import {useStore} from "vuex";
import IconSvg from "@/components/IconSvg/index";
import SiteProductCreateForm from "@/views/component/workspace/create.vue";
import {hideMenu} from "@/utils/dom";
import {useI18n} from "vue-i18n";
import {ZentaoData} from "@/store/zentao";

interface RightTopSiteProduct {
  t: (key: string | number) => string;
  sites: ComputedRef<any[]>;
  products: ComputedRef<any[]>;
  currSite: Ref
  currProduct: Ref

  selectSite: (item) => void;
  selectProduct: (item) => void;
}

export default defineComponent({
  name: 'RightTopSiteProduct',
  components: {IconSvg},
  setup(): RightTopSiteProduct {
    const { t } = useI18n();

    const router = useRouter();
    const store = useStore<{ zentao: ZentaoData }>();

    const sites = computed<any[]>(() => store.state.zentao.sites);
    const products = computed<any>(() => store.state.zentao.products);

    const currSite = computed<any>(() => store.state.zentao.currSite);
    const currProduct = computed<any>(() => store.state.zentao.currProduct);

    store.dispatch('zentao/fetchSitesAndProducts')

    onMounted(() => {
      console.log('onMounted')
    })

    const selectSite = (item): void => {
      console.log('selectSite', item)
    }
    const selectProduct = (item): void => {
      console.log('selectProduct', item)
    }

    return {
      t,
      selectSite,
      selectProduct,
      sites,
      products,
      currSite,
      currProduct,
    }
  }
})
</script>

<style lang="less">
.create-link {
  padding: 14px 10px;
  width: 150px;
  cursor: pointer;
  text-align: right;
}
.dropdown-list {
  display: inline-block;
  margin-right: 26px;
  padding-top: 13px;
  font-size: 15px !important;

  .name {
    margin-right: 5px;
  }
  .icon2 {
    .svg-icon {
      vertical-align: -3px !important;
    }
  }
}

.menu {
  .ant-dropdown-menu-item {
    cursor: default;
    .ant-dropdown-menu-title-content {
      cursor: default;
      .line {
        display: flex;
        .name {
          flex: 1;
          margin-top: 3px;
          font-size: 16px;
        }
        .space {
          width: 20px;
        }
        .icon {
          width: 15px;
          font-size: 16px;
          line-height: 28px;
        }
      }

    }
  }

}
</style>