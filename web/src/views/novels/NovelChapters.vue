<template lang="pug">
//- 章节序列
mixin NoColumn
  el-table-column(
    prop="no"
    key="no"
    label="章节序号"
    width="100"
    sortable
  )

//- 章节名称
mixin TitleColumn
  el-table-column(
    prop="title"
    key="title"
    label="标题"
  )

//- 章节字数
mixin WordCountColumn
  el-table-column(
    prop="wordCount"
    key="wordCount"
    label="章节字数"
    width="100"
  )

//- 更新时间
mixin UpdatedAtColumn
  el-table-column(
    sortable
    prop="updatedAt"
    key="updatedAt"
    label="更新时间"
    width="160"
  ): template(
    #default="scope"
  ): time-formater(
    :time="scope.row.updatedAt"
  )
//- 操作
mixin OpColumn
  el-table-column(
    fixed="right"
    label="操作"
    width="80"
  ): template(
    #default="scope"
  ): .tac
    router-link.mright10(
      :to="{name: chapterDetailRoute, params: {id: $route.params.id, no: scope.row.no}}"
    )
      i.el-icon-edit
      span 编辑

mixin Pagination
  el-pagination.pagination(
    v-if="chapters.count >= 0"
    layout="prev, pager, next, sizes"
    :current-page="query.page"
    :page-size="query.limit"
    :page-sizes="pageSizes"
    :total="chapters.count"
    @size-change="handleSizeChange"
    @current-change="handleCurrentChange"
  )

.novelChapters
  el-card
    template(
      #header
    )
      a.mright10.bold(
        @click.prevent="$router.back()"
        href="#"
      ): i.el-icon-arrow-left
      i.el-icon-s-operation
      span 小说章节列表
    div(
      v-loading="chapters.processing"
    ): el-table(
      :data="chapters.items"
      row-key="id"
      stripe
      @sort-change="handleSortChange"
    )
      +NoColumn

      +TitleColumn

      +WordCountColumn

      +UpdatedAtColumn

      +OpColumn

    //- 分页设置
    +Pagination
</template>

<script lang="ts">
import { defineComponent, onUnmounted } from "vue";

import { PAGE_SIZES } from "../../constants/common";
import TimeFormater from "../../components/TimeFormater.vue";
import useNovelState, {
  novelListChapter,
  novelChaptersClear,
} from "../../states/novel";
import { NOVEL_CHAPTER_DETAIL } from "../../router";

export default defineComponent({
  name: "NovelChapters",
  components: {
    TimeFormater,
  },
  setup() {
    onUnmounted(() => {
      novelChaptersClear();
    });

    const novelState = useNovelState();
    return {
      chapterDetailRoute: NOVEL_CHAPTER_DETAIL,
      pageSizes: PAGE_SIZES,
      chapters: novelState.chapters,
    };
  },
  data() {
    const { query } = this.$route;
    return {
      query: {
        page: Number(query.page || 1),
        limit: Number(query.limit || PAGE_SIZES[0]),
        order: query.order || "-no",
        fields: "id,title,no,wordCount,updatedAt",
      },
    };
  },
  beforeMount() {
    this.fetch();
  },
  methods: {
    updateRouteQuery() {
      this.$router.replace({
        name: this.$route.name,
        query: this.query,
      });
      this.fetch();
    },
    async fetch() {
      const { query, chapters, $route } = this;
      if (chapters.processing) {
        return;
      }
      try {
        const params = Object.assign({}, query);
        params.offset = (params.page - 1) * params.limit;
        delete params.page;
        const novelID = Number($route.params.id);
        await novelListChapter(novelID, params);
      } catch (err) {
        this.$error(err);
      }
    },
    handleCurrentChange(page: number): void {
      this.query.page = page;
      this.updateRouteQuery();
    },
    handleSizeChange(pageSize: number): void {
      this.query.limit = pageSize;
      this.query.page = 1;
      this.updateRouteQuery();
    },
    handleSortChange(params: { prop: string; order: string }): void {
      let key = params.prop;
      if (!key) {
        return;
      }
      if (params.order === "descending") {
        key = `-${key}`;
      }
      this.query.order = key;
      this.query.page = 1;
      this.updateRouteQuery();
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "../../common";

.novelChapters
  margin $mainMargin

.btn
  margin-top $mainMargin

.pagination
  text-align right
  margin-top 15px
</style>
