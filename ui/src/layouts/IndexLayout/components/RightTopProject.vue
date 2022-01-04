<template>
    <a-select
        ref="select"
        v-model:value="currProject.id"
        :bordered="true"
        style="width: 160px"
        @change="selectProject"
    >
      <a-select-option v-for="item in projects" :key="item.id" :value="item.id">{{ item.name }}</a-select-option>
    </a-select>
</template>
<script lang="ts">
import { computed, ComputedRef, defineComponent, onMounted } from "vue";
import { useStore } from "vuex";

import {ProjectData} from "@/store/project";

interface RightTopProject {
  projects: ComputedRef<any[]>;
  currProject: ComputedRef<any>;

  selectProject: () => void;
}

export default defineComponent({
    name: 'RightTopProjectSelection',
    components: {
    },
    setup(): RightTopProject {
      const store = useStore<{project: ProjectData}>();

      const projects = computed<any[]>(()=> store.state.project.projects);
      const currProject = computed<any>(()=> store.state.project.currProject);

      store.dispatch('project/fetchProject');

      onMounted(()=>{
        console.log('onMounted')
      })

      const selectProject = (): void => {
        console.log('selectProject')
      }

        return {
          selectProject,
          projects,
          currProject
        }
    }
})
</script>