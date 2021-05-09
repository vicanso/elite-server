<template lang="pug">
el-col(
  :span="8"
): el-form-item(
  label="最新版本："
): el-input(
  placeholder="请输入最新版本号"
  v-model="form.latestVersion"
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
  },
  methods: {
    handleChange() {
      const value = JSON.stringify(this.form);
      this.$emit("change", value);
    },
  },
});
</script>
