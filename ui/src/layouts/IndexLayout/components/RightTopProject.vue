<template>
  <div>
    <a-dropdown class="dropdown" :trigger="['click']">
      <a class=" t-link-btn" @click.prevent>
        <span class="name">{{currProject.name}}</span>
        <DownOutlined />
      </a>
      <template #overlay>
        <a-menu class="menu">
          <a-menu-item  v-for="item in projects" :key="item.path">
            <template v-if="currProject.path !== item.path">
              <div class="line">
                <div class="t-link name" @click="selectProject(item)">{{ item.name }}</div>
                <div class="space"></div>
                <div class="t-link icon" @click="removeProject(item)">
                  <icon-svg type="close"></icon-svg>
                </div>
              </div>
            </template>
          </a-menu-item>

          <a-menu-divider />
          <a-menu-item key="">
            <div class="t-link name" @click="selectProject('')">{{ t('create') }}</div>
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>

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
import IconSvg from "@/components/IconSvg/index";
import {ProjectData} from "@/store/project";
import ProjectCreateForm from "@/views/component/project/create.vue";
import {createProject} from "@/services/project";
import {hideMenu} from "@/utils/dom";
import {useI18n} from "vue-i18n";
import { DownOutlined } from '@ant-design/icons-vue';

interface RightTopProject {
  t: (key: string | number) => string;
  projects: ComputedRef<any[]>;
  currProject: ComputedRef;

  selectProject: (item) => void;
  removeProject: (item) => void;
  formVisible: Ref<boolean>;
  setFormVisible:  (val: boolean) => void;
  submitForm: (project: any) => Promise<void>;
  cancel: () => void;
}

export default defineComponent({
  name: 'RightTopProject',
  components: {DownOutlined, ProjectCreateForm, IconSvg},
  setup(): RightTopProject {
    const { t } = useI18n();

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

    const selectProject = (item): void => {
      console.log('selectProject', item)

      if (!item) {
        setFormVisible(true)
      } else {
        store.dispatch('project/fetchProject', item.path)
      }
    }
    const removeProject = (item): void => {
      console.log('removeProject', item)

      if (item.key === '') {
        setFormVisible(true)
      } else {
        store.dispatch('project/removeProject', item.path)
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
      t,
      selectProject,
      removeProject,
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

<style lang="less">
.dropdown {
  display: inline-block;
  padding: 13px 0;
  width: 150px;
  font-size: 15px !important;
  text-align: right;

  .name {
    display: inline-block;
    width: 130px;
    line-height: 15px;
    text-overflow: ellipsis;
    overflow: hidden;
    margin-right: 3px;
  }
  .anticon-down {
    font-size: 16px !important;
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
          font-size: 10px;
        }
      }

    }
  }

}
</style>