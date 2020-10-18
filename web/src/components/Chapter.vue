<template>
  <div class="chapter">
    <BaseEditor
      v-if="!processing && fields"
      title="更新小说章节"
      icon="el-icon-user"
      :id="id"
      :findByID="getNovelChapterByID"
      :updateByID="updateNovelChapterByID"
      :fields="fields"
    />
  </div>
</template>
<script>
import { mapActions } from "vuex";
import BaseEditor from "@/components/base/Editor.vue";
const fields = [
  {
    label: "名称： ",
    key: "name",
    disabled: true
  },
  {
    label: "作者： ",
    key: "author",
    disabled: true
  },
  {
    label: "章节序号：",
    key: "chapterNO",
    disabled: true,
    labelWidth: "100px"
  },
  {
    label: "章节名称：",
    key: "title",
    labelWidth: "100px",
    span: 24
  },
  {
    label: "章节内容：",
    key: "content",
    type: "textarea",
    labelWidth: "100px",
    autosize: {
      minRows: 10
    },
    span: 24
  }
];
export default {
  name: "Chapter",
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
    ...mapActions(["getNovelChapterByID", "updateNovelChapterByID"])
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
