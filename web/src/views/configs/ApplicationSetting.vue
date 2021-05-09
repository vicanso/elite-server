<template lang="pug">
//- 表格
mixin Table
  config-table(
    :category="category"
    name="应用配置"
  )
  .add
    el-button.addBtn(
      type="primary"
      @click="add"
    ) 添加
//- 配置编辑
mixin Editor
  config-editor(
    name="添加/更新应用配置"
    summary="用于配置应用程序相关设置"
    :category="category"
    :defaultValue="defaultValue"
  ): template(
    #data="configProps"
  ): application-setting-data(
    :data="configProps.form.data"
    @change.self="configProps.form.data = $event"
  )
.aplicationSetting
  div(
    v-if="!editMode"
  )
    +Table

  //- 编辑
  template(
    v-else
  )
    +Editor
</template>

<script lang="ts">
import { defineComponent } from "vue";

import ConfigEditor from "../../components/configs/Editor.vue";
import ApplicationSettingData from "../../components/configs/ApplicationSettingData.vue";
import ConfigTable from "../../components/configs/Table.vue";
import { APPLICATION_SETTING, CONFIG_EDIT_MODE } from "../../constants/common";

export default defineComponent({
  name: "ApplicationSetting",
  components: {
    ApplicationSettingData,
    ConfigTable,
    ConfigEditor,
  },
  data() {
    return {
      defaultValue: {
        category: APPLICATION_SETTING,
      },
      category: APPLICATION_SETTING,
    };
  },
  computed: {
    editMode() {
      return this.$route.query.mode === CONFIG_EDIT_MODE;
    },
  },
  methods: {
    add() {
      this.$router.push({
        query: {
          mode: CONFIG_EDIT_MODE,
        },
      });
    },
  },
});
</script>
<style lang="stylus" scoped>
@import "../../common";

.add
  margin $mainMargin
.addBtn
  width 100%
</style>
