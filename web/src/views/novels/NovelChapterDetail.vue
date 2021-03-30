<template lang="pug">
.novelChapterDetail
  el-card(
    v-loading="chapter.processing"
  )
    template(
      #header
    )
      a.mright10.bold(
        @click.prevent="$router.back()"
        href="#"
      ): i.el-icon-arrow-left
      | {{chapter.data.title}}
    el-input(
      type="textarea"
      :autosize="{ minRows: 20, maxRows: 10 }" 
      v-model="content"
    ) 
    el-button.fullFill.mtopMain(
      type="primary"
      @click="update"
    ) 更新
</template>

<script lang="ts">
import { defineComponent } from "vue";

import useNovelState, {
  novelGetChapterDetail,
  novelUpdateChapterDetail,
} from "../../states/novel";

export default defineComponent({
  name: "NovelChapterDetail",
  setup() {
    const novelState = useNovelState();
    return {
      chapter: novelState.chapterDetail,
    };
  },
  data() {
    return {
      content: "",
    };
  },
  beforeMount() {
    this.fetch();
  },
  methods: {
    async fetch() {
      const { id, no } = this.$route.params;
      try {
        await novelGetChapterDetail({
          id: Number(id),
          no: Number(no),
        });
        this.content = this.chapter.data.content;
      } catch (err) {
        this.$error(err);
      }
    },
    async update() {
      const { content } = this;
      if (content == this.chapter.data.content) {
        this.$message.warning("请先修改再更新");
        return;
      }
      const { id, no } = this.$route.params;
      try {
        await novelUpdateChapterDetail({
          id: Number(id),
          no: Number(no),
          content,
        });
        this.$message.info("章节内容更新成功");
      } catch (err) {
        this.$error(err);
      }
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "../../common";

.novelChapterDetail
  margin $mainMargin
</style>
