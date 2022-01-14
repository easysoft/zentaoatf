<template>
  <div>
    <a-select
        ref="select"
        v-model:value="currProject.path"
        :bordered="true"
        style="width: 160px"
        @change="selectProject"
    >
      <a-select-option v-for="item in projects" :key="item.id" :value="item.path">{{ item.name }}</a-select-option>
      <a-select-option key="" value="">新建</a-select-option>
    </a-select>

    <dir-selection
        :visible="selectionFormVisible"
        :onCancel="selectionCancel"
        :onSubmit="selectionSubmit"
    />
  </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, Ref, ref} from "vue";
import {useStore} from "vuex";

import {ProjectData} from "@/store/project";
import {Execution} from "@/views/execution/data";
import {Props} from "ant-design-vue/lib/form/useForm";
import {message} from "ant-design-vue";
import DirSelection from "@/views/component/file/DirSelection.vue";

interface RightTopProject {
  projects: ComputedRef<any[]>;
  currProject: ComputedRef;

  selectProject: (value: string) => void;
  selectionFormVisible: Ref<boolean>;
  setSelectionFormVisible:  (val: boolean) => void;
  selectionSubmit: (parentDir: string) => Promise<void>;
  selectionCancel: () => void;
}

export default defineComponent({
  name: 'RightTopProjectSelection',
  components: {DirSelection},
  setup(): RightTopProject {
    const store = useStore<{ project: ProjectData }>();

    const projects = computed<any[]>(() => store.state.project.projects);
    const currProject = computed<any>(() => store.state.project.currProject);

    store.dispatch('project/fetchProject', '');

    onMounted(() => {
      console.log('onMounted')
    })

    const selectProject = (value): void => {
      console.log('selectProject', value)

      if (value === '') {
        setSelectionFormVisible(true)
      } else {
        store.dispatch('project/fetchProject', value);
      }
    }

    const selectionFormVisible = ref<boolean>(false);
    const setSelectionFormVisible = (val: boolean) => {
      selectionFormVisible.value = val;
    };

    const selectionSubmit = async (parentDir: string) => {
      console.log('selectionSubmit', parentDir)

      await store.dispatch('project/fetchProject', parentDir);

      setSelectionFormVisible(false);
    }
    const selectionCancel = () => {
      store.dispatch('project/fetchProject', currProject.value.path);
      setSelectionFormVisible(false);
    }

    return {
      selectProject,
      projects,
      currProject,
      selectionFormVisible,
      setSelectionFormVisible,
      selectionSubmit,
      selectionCancel,
    }
  }
})
</script>