<template lang="pug">
el-col(
  :span="8"
): el-form-item(
  label="最新版本："
): el-input(
  placeholder="请输入最新版本号"
  v-model="form.latestVersion"
  clearable
)
el-col(
  :span="8"
): el-form-item(
  label="适用版本："
): el-input(
  placeholder="请输入适用版本"
  v-model="form.applIcableVersion"
  clearable
)
el-col(
  :span="8"
): el-form-item(
  label-width="100px"
  label="预加载章节："
): el-input(
  placeholder="请输入预加载章节数"
  type="number"
  v-model="form.prefetchSize"
  clearable
)
el-col(
  :span="8"
)
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "ApplicationSettingData",
  props: {
    data: {
      type: String,
      default: "",
    },
  },
  emits: ["change"],
  data() {
    const form = {
      latestVersion: "",
      applIcableVersion: "",
      prefetchSize: null,
    };
    if (this.$props.data) {
      const data = JSON.parse(this.$props.data);
      Object.assign(form, data);
    }
    return {
      form,
    };
  },
  watch: {
    "form.latestVersion": function () {
      this.handleChange();
    },
    "form.applIcableVersion": function () {
      this.handleChange();
    },
    "form.prefetchSize": function () {
      this.handleChange();
    },
  },
  methods: {
    handleChange() {
      const data = Object.assign({}, this.form);
      if (data.prefetchSize) {
        data.prefetchSize = Number(data.prefetchSize);
      }
      const value = JSON.stringify(data);
      this.$emit("change", value);
    },
  },
});
</script>
