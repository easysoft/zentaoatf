<template>
  <div v-if="!currProject.path">
    <a-empty :image="simpleImage" :description="t('pls_create_project')"/>
  </div>

  <div v-if="currProject.path" class="indexlayout-main-conent">
    <a-card :bordered="false">
      <template #title>
        {{ t('test_exec') }}
      </template>
      <template #extra>
        <div class="opt">
          <template v-if="currProject.type === 'func'">
            <a-button @click="execCase" type="primary" class="exec-button">
              <span class="exec-icon"><icon-svg type="exec"></icon-svg></span>
              <span class="exec-text">{{ t('exec') }}{{ t('case') }}</span>
            </a-button>

            <a-dropdown>
              <a-button type="primary" class="exec-button">
                <span class="button-text">
                  <span class="exec-icon"><icon-svg type="exec"></icon-svg></span>
                  <span class="exec-text">{{ t('exec') }}</span>
                </span>
                <icon-svg type="down"></icon-svg>
              </a-button>

              <template #overlay>
                <a-menu class="menu">
                  <a-menu-item @click="execModule" class="t-link">
                    <span class="t-link">{{ t('module') }}</span>
                  </a-menu-item>
                  <a-menu-item @click="execSuite" class="t-link">
                    <span class="t-link">{{ t('suite') }}</span>
                  </a-menu-item>
                  <a-menu-item @click="execTask" class="t-link">
                    <span class="t-link">{{ t('task') }}</span>
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>

          </template>

          <template v-if="currProject.type === 'unit'">
            <a-button @click="execUnit" type="primary">{{ t('execute_unit_or_automated') }}</a-button>
          </template>
        </div>
      </template>

      <div>
        <a-table
            row-key="seq"
            :columns="columns"
            :data-source="models"
            :loading="loading"
            :pagination="false"
        >
          <template #seq="{ text }">
            {{ text }}
          </template>
          <template #execBy="{ record }">
            {{ execBy(record) }}
          </template>
          <template #startTime="{ record }">
            <span v-if="record.startTime">{{ momentTime(record.startTime) }}</span>
          </template>
          <template #duration="{ record }">
            {{ record.duration }}秒
          </template>
          <template #result="{ record }">
            <span class="t-pass t-status">
              {{ record.pass }}&nbsp;
              <icon-svg type="pass"></icon-svg>&nbsp;
              ({{ percent(record.pass, record.total) }})
            </span>
            <span class="t-fail t-status">
              {{ record.fail }}&nbsp;
              <icon-svg type="fail"></icon-svg>&nbsp;
              ({{ percent(record.fail, record.total) }})
            </span>
            <span class="t-skip t-status">
              {{ record.skip }}&nbsp;
              <icon-svg type="skip"></icon-svg>&nbsp;
              ({{ percent(record.skip, record.total) }})
            </span>
          </template>
          <template #action="{ record }">
            <a-button @click="() => viewResult(record)" type="link" size="small">{{ t('view') }}</a-button>
            <a-button @click="() => deleteExec(record)" type="link" size="small"
                      :loading="deleteLoading.includes(record.seq)">{{ t('delete') }}
            </a-button>
          </template>

        </a-table>
      </div>
    </a-card>
  </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, ref, Ref, watch} from "vue";
import {Execution} from '../data.d';
import {useStore} from "vuex";

import {Empty, Form, message, Modal} from "ant-design-vue";
import {StateType} from "../store";
import {useRouter} from "vue-router";
import {momentUnixDef, percentDef} from "@/utils/datetime";
import {execByDef} from "@/utils/testing";
import {ProjectData} from "@/store/project";
import {hideMenu} from "@/utils/dom";
import throttle from "lodash.throttle";
import {useI18n} from "vue-i18n";
import IconSvg from "@/components/IconSvg/index";

const useForm = Form.useForm;

interface ListExecSetupData {
  t: (key: string | number) => string;
  currProject: ComputedRef;

  columns: any;
  models: ComputedRef<Execution[]>;
  loading: Ref<boolean>;
  list: () => void
  viewResult: (item) => void;

  deleteLoading: Ref<string[]>;
  deleteExec: (item) => void;

  execCase: () => void;
  execModule: () => void;
  execSuite: () => void;
  execTask: () => void;
  execUnit: () => void;

  execBy: (item) => string;
  momentTime: (tm) => string;
  percent: (numb, total) => string;
  simpleImage: any
}

export default defineComponent({
  name: 'ExecListPage',
  components: {
    IconSvg,
  },
  setup(): ListExecSetupData {
    const {t} = useI18n();

    const projectStore = useStore<{ project: ProjectData }>();
    const currProject = computed<any>(() => projectStore.state.project.currProject);

    const execBy = execByDef
    const momentTime = momentUnixDef
    const percent = percentDef

    const columns = [
      {
        title: t('no'),
        dataIndex: 'seq',
      },
      {
        title: t('exec_type'),
        dataIndex: 'execBy',
        slots: {customRender: 'execBy'},
      },
      {
        title: t('start_time'),
        dataIndex: 'startTime',
        slots: {customRender: 'startTime'},
      },
      {
        title: t('duration'),
        dataIndex: 'duration',
        slots: {customRender: 'duration'},
      },
      {
        title: t('result'),
        dataIndex: 'result',
        slots: {customRender: 'result'},
      },
      {
        title: t('opt'),
        key: 'action',
        width: 260,
        slots: {customRender: 'action'},
      },
    ];

    const router = useRouter();
    const store = useStore<{ History: StateType }>();

    const models = computed<any[]>(() => store.state.History.items);
    const loading = ref<boolean>(true);
    const list = throttle(async () => {
      loading.value = true;
      await store.dispatch('History/list', {});
      loading.value = false;
    }, 600)
    list();

    watch(currProject, (newProject, oldVal) => {
      console.log('watch currProject', newProject)
      list()
    }, {deep: true})

    onMounted(() => {
      console.log('onMounted')
      hideMenu(currProject.value) // jump from not available page for unittest
    })

    // 查看
    const viewResult = (item) => {
      router.push(`/exec/history/${item.testType}/${item.seq}`)
    }

    // 删除
    const deleteLoading = ref<string[]>([]);
    const deleteExec = (item) => {
      Modal.confirm({
        title: t('confirm_to_delete_result'),
        okText: t('confirm'),
        cancelText: t('cancel'),
        onOk: async () => {
          deleteLoading.value = [item.seq];
          const res: boolean = await store.dispatch('History/delete', item.seq);
          if (res === true) {
            message.success(t('delete_success'));
            await list();
          }
          deleteLoading.value = [];
        }
      });
    }

    const execCase = () => {
      console.log("execCase")
      router.push(`/exec/run/case/-/-`)
    }
    const execModule = () => {
      console.log("execModule")
      router.push(`/exec/run/module/0/0/-/-`)
    }
    const execSuite = () => {
      console.log("execSuite")
      router.push(`/exec/run/suite/0/0/-/-`)
    }
    const execTask = () => {
      console.log("execSuite")
      router.push(`/exec/run/task/0/0/-/-`)
    }

    const execUnit = () => {
      console.log("execUnit")
      router.push(`/exec/run/unit`)
    }

    return {
      t,
      currProject,

      columns,
      models,
      loading,
      list,

      viewResult,
      deleteLoading,
      deleteExec,

      execCase,
      execModule,
      execSuite,
      execTask,
      execUnit,

      execBy,
      momentTime,
      percent,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,
    }
  }

})
</script>

<style lang="less" scoped>
.exec-button {
  padding-left: 23px;
  .exec-icon {
    display: inline-block;
    margin-right: 5px;
  }
  .button-text {
    display: inline-block;
    margin-right: 6px;
  }
}

</style>