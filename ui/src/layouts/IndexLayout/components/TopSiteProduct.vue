<template>
  <div v-if="sites.length > 0">
    <!-- zentao site selection -->
    <a-dropdown
        :dropdownMatchSelectWidth="false"
        class="dropdown-list">

      <a class="t-link-btn" @click.prevent>
        <span class="name">{{ currSite.name }}</span>
        <span class="icon2"><icon-svg type="down"></icon-svg></span>
      </a>
      <template #overlay>
        <a-menu v-if="sites.length > 1" class="menu">
          <template v-for="item in sites" :key="item.id">
            <a-menu-item v-if="currSite.id !== item.id">
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
        v-if="products.length > 0"
        :dropdownMatchSelectWidth="false"
        class="dropdown-list">

      <a class="t-link-btn" @click.prevent>
        <span class="name">{{currProduct.name ? currProduct.name : '无'}}</span>
        <span class="icon2"><icon-svg type="down"></icon-svg></span>
      </a>
      <template #overlay>
        <a-menu v-if="products.length > 1" class="menu">
          <template v-for="item in products" :key="item.id">
            <a-menu-item v-if="currProduct.id !== item.id">
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
import {notification} from 'ant-design-vue';
import {useStore} from "vuex";
import IconSvg from "@/components/IconSvg/index";
import {useI18n} from "vue-i18n";
import {ZentaoData} from "@/store/zentao";
import {getInitStatus} from "@/utils/cache";

interface TopSiteProductSetupData {
  t: (key: string | number) => string;
  sites: ComputedRef<any[]>;
  products: ComputedRef<any[]>;
  currSite: Ref
  currProduct: Ref

  selectSite: (item) => void;
  selectProduct: (item) => void;
}

export default defineComponent({
  name: 'TopSiteProduct',
  components: {IconSvg},
  setup(): TopSiteProductSetupData {
    const { t } = useI18n();
    const router = useRouter();

    const store = useStore<{ Zentao: ZentaoData }>();

    const sites = computed<any[]>(() => store.state.Zentao.sites);
    const products = computed<any>(() => store.state.Zentao.products);

    const currSite = computed<any>(() => store.state.Zentao.currSite);
    const currProduct = computed<any>(() => store.state.Zentao.currProduct);

    store.dispatch('Zentao/fetchSitesAndProduct', {}).then((payload) => {
      showZentaoMsg(payload)
    })

    watch(currSite, ()=> {
      console.log(`watch currSite id = ${currSite.value.id}`)

      if (currSite.value.id <= 0) {
        getInitStatus().then((initStatus) => {
          console.log('initStatus', initStatus)
          if (!initStatus) {
            router.push(`/site/list`)
          }
        })
      }
    })

    onMounted(() => {
      console.log('onMounted')
    })

    const showZentaoMsg = (payload): void => {
      if (payload.zentaoErr) {
        notification.error({
          message: t('zentao_request_failed_title'),
          description: t('zentao_request_failed_desc'),
          duration: null,
        });
      }
    }

    const selectSite = (site): void => {
      console.log('selectSite', site)
      store.dispatch('Zentao/fetchSitesAndProduct', {currSiteId: site.id}).then((payload) => {
        showZentaoMsg(payload)
      })
    }
    const selectProduct = (product): void => {
      console.log('selectProduct', product)
      store.dispatch('Zentao/fetchSitesAndProduct', {currProductId: product.id}).then((payload) => {
        showZentaoMsg(payload)
      })
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