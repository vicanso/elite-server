<template>
  <div class="novels">
    <el-card v-if="!editMode">
      <div slot="header">
        <span>小说查询</span>
      </div>
      <div v-loading="processing">
        <BaseFilter :fields="filterFields" @filter="filter" />
        <el-table :data="novels" row-key="id" stripe>
          <el-table-column prop="name" key="name" label="名称" width="200" />
          <el-table-column
            prop="author"
            key="author"
            label="作者"
            width="150"
          />
          <el-table-column prop="source" key="source" label="来源" width="80">
            <template slot-scope="scope">
              {{ sourceNameList[scope.row.source] || sourceNameList[0] }}
            </template>
          </el-table-column>
          <el-table-column prop="status" key="status" label="状态" width="80">
            <template slot-scope="scope">
              {{ statusList[scope.row.status] || statusList[0] }}
            </template>
          </el-table-column>
          <el-table-column prop="summary" key="summary" label="简介" />
          <el-table-column label="操作" width="80">
            <template slot-scope="scope">
              <el-button
                class="op"
                type="text"
                size="small"
                @click="modify(scope.row)"
                >编辑</el-button
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
          :total="novelCount"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    <Novel v-else />
  </div>
</template>
<script>
import { mapActions, mapState } from "vuex";
import BaseTable from "@/components/base/Table.vue";
import BaseFilter from "@/components/base/Filter.vue";
import { NOVEL_SOURCES, NOVEL_STATUSES, PAGE_SIZES } from "@/constants/common";
import Novel from "@/components/Novel.vue";

const filterFields = [
  {
    label: "关键字：",
    key: "keyword",
    placeholder: "请输入要查询的关键字",
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
  name: "Novels",
  extends: BaseTable,
  components: {
    Novel,
    BaseFilter
  },
  data() {
    return {
      sourceNameList: NOVEL_SOURCES,
      statusList: NOVEL_STATUSES,
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
      processing: state => state.novel.listProcessing,
      novels: state => state.novel.list.data || [],
      novelCount: state => state.novel.list.count
    })
  },
  methods: {
    ...mapActions(["listNovel"]),
    async fetch() {
      const { query, processing } = this;
      if (processing) {
        return;
      }
      try {
        await this.listNovel(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  }
};
</script>

<style lang="sass" scoped>
@import "@/common.sass"
.novels
  margin: $mainMargin
  i
    margin-right: 5px
.pagination
  text-align: right
  margin-top: 15px
</style>
