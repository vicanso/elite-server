<template lang="pug">
//- 配置名称输入
mixin NameInput
  el-col(
    :span="8"
  ): el-form-item(
    label="名称："
  ): el-input(
    placeholder="请输入配置名称"
    v-model="form.name"
    clearable
    :disabled="!!$props.defaultValue.name"
  )

//- 配置分类输入
mixin CategoryInput
  el-col(
    :span="8"
  ): el-form-item(
    label="分类："
  ): el-input(
    placeholder="请输入配置分类（可选）"
    v-model="form.category"
    clearable
    :disabled="!!$props.defaultValue.category"
  )

//- 配置状态选择
mixin StatusSelect
  el-col(
    :span="8"
  ): el-form-item(
    label="是否启用："
  ): el-select.selector(
    v-model="form.status"
    placeholder="请选择配置状态"
  ): el-option(
    v-for="item in statuses.items"
    :key="item.value"
    :label="item.label"
    :value="item.value"
  )

//- 配置开始时间选择
mixin StartedSelect
  el-col(
    :span="8"
  ): el-form-item(
    label="开始时间："
  ): el-date-picker.datePicker.fullFill(
    v-model="form.startedAt"
    type="datetime"
    placeholder="选择开始时间"
  )

//- 配置结束时间选择
mixin EndedSelect
  el-col(
    :span="8"
  ): el-form-item(
    label="结束时间："
  ): el-date-picker.datePicker.fullFill(
    v-model="form.endedAt"
    type="datetime"
    placeholder="选择结束时间"
  )

.configurationEditor: el-card(
  v-loading="fetching"
)
  template(
    #header
  )
    i.el-icon-s-tools
    span {{ $props.name || "添加/更新配置" }}
  el-form.form(
    label-width="90px"
    v-loading="configs.processing"
    v-if="!fetching"
  )
    p
      i.el-icon-info
      | {{ $props.summary || "添加或更新系统配置信息" }}
    el-row(
      :gutter="15"
    )
      //- 配置名称
      +NameInput
 
      //- 配置分类
      +CategoryInput

      //- 是否启用
      +StatusSelect

      //- 开始时间
      +StartedSelect
 
      //- 结束时间
      +EndedSelect

      //- 配置内容
      slot(
        :form="form"
        name="data"
      )
      //- 确认提交按钮
      el-col(
        :span="primarySpan"
      ): el-form-item: ex-button(
        :onClick="submit"
      ) {{ submitText }}
      //- 返回
      el-col(
        :span="12"
        v-if="!$props.backDisabled"
      ): el-form-item: el-button.submit(
        @click="goBack"
      ) 返回
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { diff } from "../../helpers/util";
import useConfigState, {
  configFindByID,
  configAdd,
  configUpdateByID,
} from "../../states/config";

import ExButton from "../ExButton.vue";

export default defineComponent({
  name: "ConfigEditor",
  components: {
    ExButton,
  },
  props: {
    defaultValue: {
      type: Object,
      default: () => {
        return {};
      },
    },
    category: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      default: "",
    },
    summary: {
      type: String,
      default: "",
    },
    // 返回函数
    back: {
      type: Function,
      default: null,
    },
    backDisabled: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const configState = useConfigState();
    return {
      statuses: configState.statuses,
      configs: configState.configs,
      current: configState.configs.current,
    };
  },
  data() {
    const { $props, $route } = this;
    const { defaultValue, backDisabled } = $props;
    const submitText = $route.query.id ? "更新" : "提交";
    const primarySpan = backDisabled ? 24 : 12;
    return {
      primarySpan,
      originalValue: null,
      fetching: false,
      submitText,
      id: 0,
      form: {
        name: defaultValue.name || "",
        category: defaultValue.category || "",
        status: null,
        startedAt: "",
        endedAt: "",
        data: "",
      },
    };
  },
  watch: {
    $route() {
      this.fetchCurrent();
    },
  },
  beforeMount() {
    this.fetchCurrent();
  },
  methods: {
    async submit(): Promise<boolean> {
      let isSuccess = false;
      const { name, category, status, startedAt, endedAt, data } = this.form;
      if (!name || !status || !startedAt || !endedAt || !data) {
        this.$message.warning("名称、状态、开始结束日期以及配置数据均不能为空");
        return isSuccess;
      }
      const { id } = this;
      try {
        const config = {
          name,
          status,
          category,
          startedAt,
          endedAt,
          data,
        };
        if (startedAt.toISOString) {
          config.startedAt = startedAt.toISOString();
        }
        if (endedAt.toISOString) {
          config.endedAt = endedAt.toISOString();
        }
        // 更新
        if (id) {
          const info = diff(config, this.originalValue);
          if (!info.modifiedCount) {
            this.$message.warning("未修改配置无法更新");
            return isSuccess;
          }
          await configUpdateByID({
            id,
            data: info.data,
          });
          this.$message.info("修改配置成功");
          isSuccess = true;
        } else {
          await configAdd(config);
          this.$message.info("添加配置成功");
        }
        this.goBack();
      } catch (err) {
        this.$error(err);
      }
      return isSuccess;
    },
    goBack() {
      if (this.$props.back) {
        this.$props.back();
        return;
      }
      this.$router.back();
    },
    // 拉取当前配置
    async fetchCurrent() {
      const { query } = this.$route;
      let currentID = this.id;
      if (query.id) {
        currentID = Number(query.id);
      }
      if (currentID === this.id) {
        return;
      }
      this.fetching = true;
      try {
        const data = await configFindByID(currentID);
        const config = {};
        Object.keys(this.form).forEach((key) => {
          config[key] = data[key];
        });
        this.originalValue = config;
        Object.assign(this.form, config);
      } catch (err) {
        this.$error(err);
      } finally {
        this.fetching = false;
      }
      this.id = currentID;
    },
  },
});
</script>
<style lang="stylus" scoped>
@import "../../common";

.configurationEditor
  margin $mainMargin
  i
    margin-right 3px
  p
    color $darkGray
    font-size 13px
    margin 0 0 15px 0
.selector, .datePicker, .submit
  width 100%
</style>
