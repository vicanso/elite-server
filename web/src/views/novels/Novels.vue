<template lang="pug">
//- 小说ID
mixin IDColumn
  el-table-column(
    prop="id"
    key="id"
    label="ID"
    width="80"
    sortable
  )

//- 小说名称
mixin NameColumn
  el-table-column(
    prop="name"
    key="name"
    label="名称"
    width="150"
  )

//- 小说作者
mixin AuthorColumn
  el-table-column(
    prop="author"
    key="author"
    label="作者"
    width="150"
  )

//- 小说状态
mixin StatusColumn
  el-table-column(
    label="状态"
    width="80"
  ): template(
    #default="scope"
  ) {{statuses.items[scope.row.status]}}

//- 小说更新时间 
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

//- 简介
mixin SummaryColumn
  el-table-column(
    label="简介"
  ): template(
    #default="scope"
  ): base-tooltip(
    :viewSize="-50"
    :content="scope.row.summary"
  )

//- 操作
mixin OpColumn
  el-table-column(
    fixed="right"
    label="操作"
    width="120"
  ): template(
    #default="scope"
  ): .tac
    router-link.mright10(
      :to="{name: detailRoute, params: {id: scope.row.id}}"
    )
      i.el-icon-edit
      span 编辑
    router-link(
      :to="{name: chaptersRoute, params: {id: scope.row.id}}"
    )
      i.el-icon-s-operation
      span 章节
mixin Pagination
  el-pagination.pagination(
    v-if="novels.count >= 0"
    layout="prev, pager, next, sizes"
    :current-page="query.page"
    :page-size="query.limit"
    :page-sizes="pageSizes"
    :total="novels.count"
    @size-change="handleSizeChange"
    @current-change="handleCurrentChange"
  )
.novels
  el-card
    template(
      #header
    )
      i.el-icon-notebook-1
      span 小说列表
    div(
      v-loading="novels.processing"
    ): el-table(
      :data="novels.items"
      row-key="id"
      stripe
      @sort-change="handleSortChange"
    )
      //- 小说ID
      +IDColumn

      //- 小说名称
      +NameColumn

      //- 小说作者
      +AuthorColumn

      //- 小说状态
      +StatusColumn

      //- 小说更新时间
      +UpdatedAtColumn

      //- 小说简介
      +SummaryColumn

      +OpColumn
    //- 分页设置
    +Pagination
</template>

<script lang="ts">
import { defineComponent, onUnmounted } from "vue";

import { PAGE_SIZES } from "../../constants/common";
import TimeFormater from "../../components/TimeFormater.vue";
import BaseTooltip from "../../components/Tooltip.vue";
import useNovelState, { novelList, novelListClear } from "../../states/novel";
import { NOVEL_DETAIl, NOVEL_CHAPTERS } from "../../router";

export default defineComponent({
  name: "Novels",
  components: {
    TimeFormater,
    BaseTooltip,
  },
  setup() {
    onUnmounted(() => {
      novelListClear();
    });
    const novelState = useNovelState();
    return {
      pageSizes: PAGE_SIZES,
      statuses: novelState.statuses,
      novels: novelState.novels,
      detailRoute: NOVEL_DETAIl,
      chaptersRoute: NOVEL_CHAPTERS,
    };
  },
  data() {
    const { query } = this.$route;
    return {
      query: {
        page: Number(query.page || 1),
        limit: Number(query.limit || PAGE_SIZES[0]),
        order: query.order || "-updatedAt",
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
      const { query, novels } = this;
      if (novels.processing) {
        return;
      }
      try {
        const params = Object.assign({}, query);
        params.offset = (params.page - 1) * params.limit;
        delete params.page;
        await novelList(params);
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

.novels
  margin $mainMargin
i
  margin-right 3px
.pagination
  text-align right
  margin-top 15px
</style>
