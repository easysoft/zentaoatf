<template>
  <div>
    <div v-if="workspaces.length == 0" class="create-link" @click="selectWorkspace('')">
      {{ t('create_workspace') }}
    </div>

    <!-- zentao site selection -->
    <a-dropdown
        v-if="workspaces.length > 0"
        :dropdownMatchSelectWidth="false"
        class="dropdown-list">

      <a class="t-link-btn" @click.prevent>
        <span class="name">{{currWorkspace.name}}</span>
        <span class="icon2"><icon-svg type="down"></icon-svg></span>
      </a>
      <template #overlay>
        <a-menu class="menu">
          <template v-for="item in workspaces" :key="item.path">
            <a-menu-item v-if="currWorkspace.path !== item.path">
                <div class="line">
                  <div class="t-link name" @click="selectWorkspace(item)">{{ item.name }}</div>
                  <div class="space"></div>
                  <div class="t-link icon" @click="setDeleteModel(item)">
                    <icon-svg type="delete" class="menu-icon"></icon-svg>
                  </div>
                </div>
            </a-menu-item>
          </template>

          <a-menu-divider v-if="workspaces.length > 1"/>

          <a-menu-item key="" class="create">
            <span class="t-link name" @click="selectWorkspace('')">
              <icon-svg type="add" class="menu-icon"></icon-svg>
              {{ t('create_workspace') }}
            </span>
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>

    <!-- zentao product selection -->
    <a-dropdown
        v-if="workspaces.length > 0"
        :dropdownMatchSelectWidth="false"
        class="dropdown-list">

      <a class="t-link-btn" @click.prevent>
        <span class="name">{{currWorkspace.name}}</span>
        <span class="icon2"><icon-svg type="down"></icon-svg></span>
      </a>
      <template #overlay>
        <a-menu class="menu">
          <template v-for="item in workspaces" :key="item.path">
            <a-menu-item v-if="currWorkspace.path !== item.path">
              <div class="line">
                <div class="t-link name" @click="selectWorkspace(item)">{{ item.name }}</div>
                <div class="space"></div>
                <div class="t-link icon" @click="setDeleteModel(item)">
                  <icon-svg type="delete" class="menu-icon"></icon-svg>
                </div>
              </div>
            </a-menu-item>
          </template>

          <a-menu-divider v-if="workspaces.length > 1"/>

          <a-menu-item key="" class="create">
            <span class="t-link name" @click="selectWorkspace('')">
              <icon-svg type="add" class="menu-icon"></icon-svg>
              {{ t('create_workspace') }}
            </span>
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>

    <workspace-create-form
      :visible="formVisible"
      :onCancel="cancel"
      :onSubmit="submitForm"
    />

    <a-modal
        v-model:visible="deleteVisible"
        title="Modal"
        :ok-text="t('confirm')"
        :cancel-text="t('cancel')"
        @ok="removeWorkspace"
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
import {WorkspaceData} from "@/store/workspace";
import WorkspaceCreateForm from "@/views/component/workspace/create.vue";
import {createWorkspace} from "@/services/workspace";
import {hideMenu} from "@/utils/dom";
import {useI18n} from "vue-i18n";

interface RightTopWorkspace {
  t: (key: string | number) => string;
  workspaces: ComputedRef<any[]>;
  currWorkspace: ComputedRef;

  selectWorkspace: (item) => void;
  removeWorkspace: (item) => void;
  formVisible: Ref<boolean>;
  setFormVisible:  (val: boolean) => void;
  submitForm: (workspace: any) => Promise<void>;
  cancel: () => void;

  deleteModel: Ref
  deleteVisible: Ref<boolean>;
  setDeleteModel: (val, model) => void;
}

export default defineComponent({
  name: 'RightTopWorkspace',
  components: {WorkspaceCreateForm, IconSvg},
  setup(): RightTopWorkspace {
    const { t } = useI18n();

    const router = useRouter();
    const store = useStore<{ workspace: WorkspaceData }>();

    const workspaces = computed<any[]>(() => store.state.workspace.workspaces);
    const currWorkspace = computed<any>(() => store.state.workspace.currWorkspace);

    store.dispatch('workspace/fetchWorkspace', '')

    const switchWorkspace = (newWorkspace, oldWorkspace) => {
      const routerPath = router.currentRoute.value.path
      if ( (oldWorkspace.id && oldWorkspace.id !== newWorkspace.id && routerPath.indexOf('/exec/history/') > -1)
          || (newWorkspace.type === 'unit' && (routerPath === '/sync' || routerPath === '/script/list'))) {
        router.push(`/exec/history`) // will call hideMenu on this page
      } else {
        hideMenu(newWorkspace)
      }
    }

    watch(currWorkspace, (newWorkspace, oldWorkspace)=> {
      console.log('watch currWorkspace', newWorkspace.type)
      switchWorkspace(newWorkspace, oldWorkspace)
    }, {deep: true})

    onMounted(() => {
      console.log('onMounted')
    })

    const selectWorkspace = (item): void => {
      console.log('selectWorkspace', item)

      if (!item) {
        setFormVisible(true)
      } else {
        store.dispatch('workspace/fetchWorkspace', item.path)
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
    const removeWorkspace = (): void => {
      console.log('removeWorkspace', deleteModel)

      store.dispatch('workspace/removeWorkspace', deleteModel.value.path).then(() => {
        deleteModel.value = {}
        deleteVisible.value = false
      })
    }

    const formVisible = ref<boolean>(false);
    const setFormVisible = (val: boolean) => {
      formVisible.value = val;
    };

    const submitForm = async (workspace: any) => {
      console.log('submitForm', workspace)
      createWorkspace(workspace).then(() => {
        store.dispatch('workspace/fetchWorkspace', workspace.path);
        setFormVisible(false);
      }).catch(err => { console.log('') })
    }
    const cancel = () => {
      store.dispatch('workspace/fetchWorkspace', currWorkspace.value.path);
      setFormVisible(false);
    }

    return {
      t,
      selectWorkspace,
      removeWorkspace,
      workspaces,
      currWorkspace,
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
.dropdown-list {
  display: inline-block;
  margin-right: 26px;
  padding-top: 13px;
  font-size: 15px !important;

  .name {
    margin-right: 5px;
  }
  .icon2 {
    .svg-icon {
      vertical-align: -3px !important;
    }
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