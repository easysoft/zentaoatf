<template>
  <div class="indexlayout-main-conent">
    <a-card :bordered="false">
      <template #title>
        <div>
          <a-button type="primary" @click="() => record()" id="com-deeptest-record" class="act-btn">录制</a-button>
          <a-button @click="() => playback()" class="act-btn">播放</a-button>
        </div>
      </template>
      <template #extra>
        <a-button type="link" @click="() => back()">返回</a-button>
      </template>

      <div class="script">
        <div class="title">测试步骤</div>
        <div class="desc">
          <div v-for="(step, index) in script.steps" :key="index" class="step">
            <div class="cmd">{{step.action}} {{step.selector}}  {{step.value}}</div>
            <div class="capture" style="border: 3px cornflowerblue;"><img :src="step.image"></div>
          </div>
        </div>
      </div>

    </a-card>
  </div>
</template>

<script lang="ts">
import {defineComponent, onMounted, onBeforeUnmount, getCurrentInstance, ComputedRef, Ref, ref, reactive} from "vue";
import { useStore } from 'vuex';
import {StateType as ListStateType} from "@/views/script/store";
import {useRouter} from "vue-router";
import {getToken} from "@/utils/localToken";
import {ScriptItem, StepItem} from "@/views/script/data";

interface DesignScriptPageSetupData {
  script: ScriptItem;

  loading: Ref<boolean>;
  getScript:  (current: number) => Promise<void>;
  back: () => void;
}

export default defineComponent({
    name: 'ScriptViewPage',
    setup(): DesignScriptPageSetupData {
      const router = useRouter();
      const store = useStore<{ ListScript: ListStateType}>();

      const script = reactive<ScriptItem>({steps: []})
      const loading = ref<boolean>(true);

      const id = +router.currentRoute.value.params.id
      console.log('id', id)
      const getScript = async (id: number): Promise<void> => {
        loading.value = true;
        // await store.dispatch('ListScript/getScript', {
        // });
        loading.value = false;
      }

      const back = ():void =>  {
        router.push(`/~/script/list`)
      }

      let init = true;
      let wsMsg = reactive({in: '', out: ''});

      let room: string | null = ''
      getToken().then((token) => {
        room = token
      })

      onMounted(() => {
        getScript(1);
      });

      return {
        script,
        loading,
        getScript,
        back,
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
