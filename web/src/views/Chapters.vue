<template>
  <div class="chapters">
    <el-card v-if="!editMode">
      <div slot="header">
        <span>小说章节查询</span>
      </div>
      <div v-loading="processing">
        <BaseFilter :fields="filterFields" @filter="filter" />
        <el-table
          :data="chapters"
          row-key="id"
          stripe
          @sort-change="handleSortChange"
        >
          <el-table-column
            prop="updatedAtDesc"
            key="updatedAtDesc"
            label="更新时间"
            width="180"
          />
          <el-table-column prop="title" key="title" label="章节名称" />
          <el-table-column
            prop="wordCountDesc"
            key="wordCountDesc"
            label="章节字数"
            width="100"
          />
          <el-table-column
            prop="contentDesc"
            key="contentDesc"
            label="章节内容"
          />
        </el-table>
        <el-pagination
          class="pagination"
          layout="prev, pager, next, sizes"
          :current-page="currentPage"
          :page-size="query.limit"
          :page-sizes="pageSizes"
          :total="chapterCount"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>
<script>
import { mapActions, mapState } from "vuex";
import BaseTable from "@/components/base/Table.vue";
import BaseFilter from "@/components/base/Filter.vue";
import { PAGE_SIZES } from "@/constants/common";

const filterFields = [
  {
    label: "关键字：",
    key: "id",
    type: "novelSelect",
    placeholder: "请输入小说关键字",
    span: 12
  },
  {
    label: "",
    type: "filter",
    labelWidth: "0px",
    span: 12
  }
];
export default {
  name: "Chapters",
  extends: BaseTable,
  components: {
    BaseFilter
  },
  data() {
    return {
      filterFields,
      pageSizes: PAGE_SIZES,
      query: {
        id: 0,
        offset: 0,
        limit: PAGE_SIZES[0],
        order: "-no"
      }
    };
  },
  computed: {
    ...mapState({
      processing: state => state.novel.listChapterProcessing,
      chapters: state => state.novel.chapterList.data || [],
      chapterCount: state => state.novel.chapterList.count
    })
  },
  methods: {
    ...mapActions(["listNovelChapter"]),
    async fetch() {
      const { query, processing } = this;
      const id = query.id;
      if (!id || processing) {
        return;
      }
      try {
        const params = Object.assign({}, query);
        delete params.id;
        await this.listNovelChapter({
          id: id,
          params
        });
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.chapters
  margin: $mainMargin
  i
    margin-right: 5px
.pagination
  text-align: right
  margin-top: 15px
</style>
