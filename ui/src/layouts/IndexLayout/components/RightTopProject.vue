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
      :onSubmit="submitForm"
    />
  </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, Ref, ref, watch} from "vue";
import {useRouter} from "vue-router";
import {useStore} from "vuex";

import {ProjectData} from "@/store/project";
import ProjectCreateForm from "@/views/component/project/create.vue";
import {createProject} from "@/services/project";
import {hideMenu} from "@/utils/dom";

interface RightTopProject {
  projects: ComputedRef<any[]>;
  currProject: ComputedRef;

  selectProject: (value: string) => void;
  formVisible: Ref<boolean>;
  setFormVisible:  (val: boolean) => void;
  submitForm: (project: any) => Promise<void>;
  cancel: () => void;
}

export default defineComponent({
  name: 'RightTopProject',
  components: {ProjectCreateForm},
  setup(): RightTopProject {
    const router = useRouter();
    const store = useStore<{ project: ProjectData }>();

    const projects = computed<any[]>(() => store.state.project.projects);
    const currProject = computed<any>(() => store.state.project.currProject);

    store.dispatch('project/fetchProject', '')

    const switchProject = (newProject, oldProject) => {
      const routerPath = router.currentRoute.value.path
      if ( (oldProject.id && oldProject.id !== newProject.id && routerPath.indexOf('/exec/history/') > -1)
          || (newProject.type === 'unit' && (routerPath === '/sync' || routerPath === '/script/list'))) {
        router.push(`/exec/history`) // will call hideMenu on this page
      } else {
        hideMenu(newProject)
      }
    }

    watch(currProject, (newProject, oldProject)=> {
      console.log('watch currProject', newProject.type)
      switchProject(newProject, oldProject)
    }, {deep: true})

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

    const submitForm = async (project: any) => {
      console.log('submitForm', project)
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
      submitForm,
      cancel,
    }
  }
})
</script>