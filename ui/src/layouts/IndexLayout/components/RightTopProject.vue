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

    <project-create-form
      :visible="formVisible"
      :onCancel="cancel"
      :onSubmit="createProject"
    />
  </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, Ref, ref} from "vue";
import {useStore} from "vuex";

import {ProjectData} from "@/store/project";
import ProjectCreateForm from "@/views/component/project/Create.vue";
import {createProject} from "@/services/project";

interface RightTopProject {
  projects: ComputedRef<any[]>;
  currProject: ComputedRef;

  selectProject: (value: string) => void;
  formVisible: Ref<boolean>;
  setFormVisible:  (val: boolean) => void;
  createProject: (project: any) => Promise<void>;
  cancel: () => void;
}

export default defineComponent({
  name: 'RightTopProject',
  components: {ProjectCreateForm},
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
        setFormVisible(true)
      } else {
        store.dispatch('project/fetchProject', value)
      }
    }

    const formVisible = ref<boolean>(false);
    const setFormVisible = (val: boolean) => {
      formVisible.value = val;
    };

    const createProject = async (project: any) => {
      console.log('createProject', project)
      createProject(project).then(() => {
        store.dispatch('project/fetchProject', project.path);
        setFormVisible(false);
      }).catch(err => { console.log('') })
    }
    const cancel = () => {
      store.dispatch('project/fetchProject', currProject.value.path);
      setFormVisible(false);
    }

    return {
      selectProject,
      projects,
      currProject,
      formVisible,
      setFormVisible,
      createProject,
      cancel,
    }
  }
})
</script>