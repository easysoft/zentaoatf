<template>
    <a-spin :spinning="loading" size="large">
        <a-card :title="t('page.home.articlechartcard.card-title')" class="homeBoxCard">
            <template #extra><a-tag color="cyan">{{t('page.home.text-day')}}</a-tag></template>
            <div class="num">{{visitData.num.toLocaleString()}}</div>
            <div class="height40">
            <div class="articleText">
                <span>
                {{t('page.home.text-daycompare')}}
                {{Math.abs(visitData.day)}}%
                    <CaretUpOutlined v-if="visitData.day > 0" class="colored4014" />
                    <CaretDownOutlined v-else class="color19be6b" />
                </span>
                <span className="margin-l10">
                {{t('page.home.text-weekcompare')}}
                {{Math.abs(visitData.week)}}%
                    <CaretUpOutlined v-if="visitData.week > 0" class="colored4014" />
                    <CaretDownOutlined v-else class="color19be6b" />
                </span>
            </div>
            </div>
            <a-divider />
            <a-row>
            <a-col :span="12">{{t('page.home.text-total')}}</a-col>
            <a-col class="text-align-right" :span="12">
                {{visitData.total.toLocaleString()}}
            </a-col>
            </a-row>
        </a-card>
    </a-spin>
</template>
<script lang="ts">
import { computed, ComputedRef, defineComponent, onMounted, Ref, ref } from "vue";
import { useStore } from "vuex";
import { useI18n } from "vue-i18n";
import { CaretUpOutlined, CaretDownOutlined} from '@ant-design/icons-vue';
import { ArticleChartDataType } from "../../data.d";
import { StateType as HomeStateType } from "../../store";

interface ArticleChartCardSetupData {
    t: (key: string | number) => string;
    loading: Ref<boolean>;
    visitData: ComputedRef<ArticleChartDataType>;
}

export default defineComponent({
    name: 'ArticleChartCard',
    components: {
        CaretUpOutlined,
        CaretDownOutlined
    },
    setup(): ArticleChartCardSetupData {
        const store = useStore<{ Home: HomeStateType}>();
        const { t } = useI18n();

        // 数据
        const visitData = computed<ArticleChartDataType>(()=> store.state.Home.articleChartData);
        // 读取数据 func
        const loading = ref<boolean>(true);
        const getData = async () => {
            loading.value = true;
            await store.dispatch('Home/queryArticleChartData');
            loading.value = false;
        }

        onMounted(()=> {
           getData();
        })

        return {
            t,
            loading,
            visitData
        }
    }
})
</script>
<style lang="less" scoped>
.homeBoxCard {
  margin-bottom: 24px;
  ::v-deep(.ant-card-head) {
    padding-left: 12px;
    padding-right: 12px;
  }
  ::v-deep(.ant-card-body) {
    padding: 12px;
  }
  ::v-deep(.ant-divider-horizontal) {
    margin: 8px 0;
  }
  .num {
    font-size: 30px;
    color: #515a6e;
  }
  .height40 {
    height: 40px;
  }
  .articleText {
    padding-top: 20px;
  }
  .color19be6b {
    color: #19be6b;
  }
  .colored4014 {
    color: #ed4014;
  }
}
</style>