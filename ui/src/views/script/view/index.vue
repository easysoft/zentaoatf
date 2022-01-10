<template>
  <div class="indexlayout-main-conent">
    <a-card :bordered="false">
      <template #title>
        <div>
          <a-button type="primary" @click="() => save()" id="com-deeptest-record" class="act-btn">保存</a-button>
        </div>
      </template>
      <template #extra>
      </template>

      <div class="script">
        {{script.content}}
      </div>

    </a-card>
  </div>
</template>

<script lang="ts">
import {
  defineComponent,
  onMounted,
  computed
} from "vue";
import { useStore } from 'vuex';
import {ScriptData} from "@/views/script/store";
import {useRouter} from "vue-router";
import {Script} from "@/views/script/data";

interface DesignScriptPageSetupData {
  script: any;
  save: () => void;
}

export default defineComponent({
    name: 'ScriptViewPage',
    setup(): DesignScriptPageSetupData {
      const router = useRouter();
      const store = useStore<{Script: ScriptData}>();

      const storeScript = useStore<{ script: ScriptData }>();
      const script = computed<Partial<Script>>(() => storeScript.state.script.detail);

      const save = ():void =>  {
        console.log('save')
      }

      onMounted(() => {
        console.log('onMounted')
      });

      return {
        script,
        save,
      }
    }
})
</script>

<style lang="less" scoped>
  .act-btn {
    margin-right: 20px;
  }

  .script {
    .title {
      font-weight: bolder;
    }
    .desc {
      .step {
        display: flex;
        .cmd {
          flex: 1;
        }
        .capture {
          width: 600px;

          img {
            height: 50px;
          }
        }
      }
    }
  }
</style>
