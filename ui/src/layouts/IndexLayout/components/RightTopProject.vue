<template>
  <div>
    <div v-if="projects.length == 0" class="create-link" @click="selectProject('')">
      {{ t('create_project') }}
    </div>

    <a-dropdown v-if="projects.length > 0" class="dropdown">
      <a class="t-link-btn" @click.prevent>
        <span class="name">{{currProject.name}}</span>
        <span class="icon2"><icon-svg type="down"></icon-svg></span>
      </a>
      <template #overlay>
        <a-menu class="menu">
          <template v-for="item in projects" :key="item.path">
            <a-menu-item v-if="currProject.path !== item.path">
                <div class="line">
                  <div class="t-link name" @click="selectProject(item)">{{ item.name }}</div>
                  <div class="space"></div>
                  <div class="t-link icon" @click="setDeleteModel(item)">
                    <icon-svg type="delete" class="menu-icon"></icon-svg>
                  </div>
                </div>
            </a-menu-item>
          </template>

          <a-menu-divider v-if="projects.length > 1"/>

          <a-menu-item key="" class="create">
            <span class="t-link name" @click="selectProject('')">
              <icon-svg type="add" class="menu-icon"></icon-svg>
              {{ t('create_project') }}
            </span>
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>

    <project-create-form
      :visible="formVisible"
      :onCancel="cancel"
      :onSubmit="submitForm"
    />

    <a-modal
        v-model:visible="deleteVisible"
        title="Modal"
        :ok-text="t('confirm')"
        :cancel-text="t('cancel')"
        @ok="removeProject"
    >
      <p>{{t('confirm_delete', {name: deleteModel.path})}}</p>
    </a-modal>

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

  deleteModel: Ref
  deleteVisible: Ref<boolean>;
  setDeleteModel: (val, model) => void;
}

export default defineComponent({
  name: 'RightTopProject',
  components: {ProjectCreateForm, IconSvg},
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

    let deleteVisible = ref(false)
    let deleteModel = ref({} as any)
    const setDeleteModel = (model: any) => {
      if (model) {
        deleteModel.value = model;
        deleteVisible.value = true;
      } else {
        deleteVisible.value = false
      }
    }
    const removeProject = (): void => {
      console.log('removeProject', deleteModel)

      store.dispatch('project/removeProject', deleteModel.value.path).then(() => {
        deleteModel.value = {}
        deleteVisible.value = false
      })
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

      deleteVisible,
      deleteModel,
      setDeleteModel,
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
.dropdown {
  display: inline-block;
  padding-top: 13px;
  width: 150px;
  font-size: 15px !important;
  text-align: right;

  .name {
    display: inline-block;
    width: 130px;
    text-overflow: ellipsis;
    overflow: hidden;
    margin-right: 5px;
  }
  .icon2 {
    .svg-icon {
      vertical-align: 3px !important;
    }
  }
  .anticon-down {
    font-size: 16px !important;
    line-height: 20px;
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

    &.create {
      text-align: center;
    }
  }

}
</style>