<template>
  <div class="novel">
    <BaseEditor
      v-if="!processing && fields"
      title="更新小说信息"
      icon="el-icon-user"
      :id="id"
      :findByID="getNovelByID"
      :updateByID="updateNovelByID"
      :fields="fields"
    />
  </div>
</template>
<script>
import { mapActions } from "vuex";
import BaseEditor from "@/components/base/Editor.vue";
const fields = [
  {
    label: "名称：",
    key: "name",
    disabled: true
  },
  {
    label: "作者：",
    key: "author",
    disabled: true
  },
  {
    label: "状态：",
    key: "status",
    type: "select",
    placeholder: "请选择状态",
    options: [
      {
        name: "连载中",
        value: 1
      },
      {
        name: "已完结",
        value: 2
      },
      {
        name: "下架",
        value: 3
      }
    ]
  },
  {
    label: "简介：",
    key: "summary",
    type: "textarea",
    autosize: {
      minRows: 5
    },
    span: 24
  }
];
export default {
  name: "Novel",
  components: {
    BaseEditor
  },
  data() {
    return {
      fields: null,
      processing: false,
      id: 0
    };
  },
  methods: {
    ...mapActions(["getNovelByID", "updateNovelByID"])
  },
  async beforeMount() {
    this.processing = true;
    const { id } = this.$route.query;
    if (id) {
      this.id = Number(id);
    }
    try {
      this.fields = fields;
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
