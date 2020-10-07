<template>
  <el-card class="novelSources">
    <div slot="header">
      <span>小说源查询</span>
    </div>
    <div v-loading="processing">
      <BaseFilter :fields="filterFields" @filter="filter" />
      <el-table :data="novelSources" row-key="id" stripe>
        <el-table-column prop="name" key="name" label="名称" />
        <el-table-column prop="author" key="author" label="作者" />
        <el-table-column prop="source" key="source" label="来源">
          <template slot-scope="scope">
            {{ sourceNameList[scope.row.source] || sourceNameList[0] }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template slot-scope="scope">
            <el-button
              :disabled="scope.row.status === 2"
              class="op"
              type="text"
              size="small"
              @click="publish(scope.row)"
              >{{ scope.row.status === 2 ? "已发布" : "发布" }}</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        class="pagination"
        layout="prev, pager, next, sizes"
        :current-page="currentPage"
        :page-size="query.limit"
        :page-sizes="pageSizes"
        :total="novelSourceCount"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </el-card>
</template>
<script>
import { mapActions, mapState } from "vuex";
import BaseTable from "@/components/base/Table.vue";
import BaseFilter from "@/components/base/Filter.vue";
import { NOVEL_SOURCES, PAGE_SIZES } from "@/constants/common";

const filterFields = [
  {
    label: "关键字：",
    key: "keyword",
    placeholder: "请输入要查询的关键字",
    clearable: true,
    span: 8
  },
  {
    label: "状态：",
    key: "status",
    placeholder: "请选择状态",
    type: "select",
    span: 8,
    options: [
      {
        name: "未发布",
        value: 1
      },
      {
        name: "已发布",
        value: 2
      }
    ]
  },
  {
    label: "",
    type: "filter",
    labelWidth: "0px",
    span: 8
  }
];

export default {
  name: "NovelSources",
  extends: BaseTable,
  components: {
    BaseFilter
  },
  data() {
    return {
      sourceNameList: NOVEL_SOURCES,
      filterFields,
      pageSizes: PAGE_SIZES,
      query: {
        offset: 0,
        limit: PAGE_SIZES[0],
        order: "-createdAt"
      }
    };
  },
  computed: {
    ...mapState({
      processing: state =>
        state.novel.sourceListProcessing || state.novel.publishing,
      novelSources: state => state.novel.sourceList.data || [],
      novelSourceCount: state => state.novel.sourceList.count
    })
  },
  methods: {
    ...mapActions(["listNovelSource", "publishNovel"]),
    async fetch() {
      const { query, processing } = this;
      if (processing) {
        return;
      }
      try {
        await this.listNovelSource(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    },
    async publish(item) {
      try {
        await this.publishNovel(item);
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.novelSources
  margin: $mainMargin
  i
    margin-right: 5px
.pagination
  text-align: right
  margin-top: 15px
</style>
