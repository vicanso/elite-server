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

const filterFields = [
  {
    label: "关键字：",
    key: "keyword",
    placeholder: "请输入要查询的账号",
    clearable: true,
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
  name: "NovelSources",
  extends: BaseTable,
  components: {
    BaseFilter
  },
  data() {
    const pageSizes = [10, 20, 30, 50];
    return {
      sourceNameList: ["未知", "笔趣阁"],
      filterFields,
      query: {
        offset: 0,
        limit: pageSizes[0],
        order: "-createdAt"
      }
    };
  },
  computed: {
    ...mapState({
      processing: state =>
        state.novel.novelSourceListProcessing || state.novel.novelPublishing,
      novelSources: state => state.novel.novelSourceList.data || [],
      novelSourceCount: state => state.novel.novelSourceList.count
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
