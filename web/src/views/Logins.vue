<template lang="pug">
mixin AccountColumn
  el-table-column(
    prop="account"
    key="account"
    label="账户"
    width="120"
  )

mixin IPColumn
  el-table-column(
    prop="ip"
    key="ip"
    label="IP"
    width="120"
  )

mixin LocationColumn
  el-table-column(
    label="定位"
    width="180"
  ): template(
    #default="scope"
  )
    template(
      v-if="scope.row.location"
    ) {{scope.row.location}}
    template(
      v-else
    ) --
  

mixin ISPColumn
  el-table-column(
    label="运营商"
    width="80"
  ): template(
    #default="scope"
  ) {{ scope.row.isp || "--" }}

mixin SessionColumn
  el-table-column(
    label="Session ID"
    width="100"
  ): template(
    #default="scope"
  ): base-tooltip(
    :content="scope.row.sessionID"
  )

mixin TrackIDColumn
  el-table-column(
    label="Track ID"
    width="100"
  ): template(
    #default="scope"
  ): base-tooltip(
    :content="scope.row.trackID"
  )

mixin ForwardedForColumn
  el-table-column(
    label="Forwarded For"
  ): template(
    #default="scope"
  ): base-tooltip(
    icon="el-icon-info"
    :content="scope.row.xForwardedFor"
  )

mixin UserAgentColumn
  el-table-column(
    label="User Agent"
  ): template(
    #default="scope"
  ): base-tooltip(
    icon="el-icon-mobile-phone"
    :content="scope.row.userAgent"
  )

mixin CreatedAtColumn
  el-table-column(
    prop="createdAt"
    key="createdAt"
    label="时间"
    width="160"
  ): template(
    #default="scope"
  ): time-formater(
    :time="scope.row.createdAt"
  )

mixin Pagination
  el-pagination.pagination(
    layout="prev, pager, next, sizes"
    :current-page="currentPage"
    :page-size="query.limit"
    :page-sizes="pageSizes"
    :total="logins.count"
    @size-change="handleSizeChange"
    @current-change="handleCurrentChange"
  )

el-card.logins
  template(
    #header
  )
    i.el-icon-user-solid
    span 用户登录查询
  div 
    base-filter(
      :fields="filterFields"
      :filter="filter"
    )
    div(
      v-loading="logins.processing"
    ): el-table(
      :data="logins.items"
      row-key="id"
      stripe
    )
      //- 账号
      +AccountColumn

      //- IP
      +IPColumn

      //- location
      +LocationColumn

      //- isp
      +ISPColumn
      
      //- session id
      +SessionColumn
      
      //- track id
      +TrackIDColumn

      //- forwarded for
      +ForwardedForColumn

      //- user agent
      +UserAgentColumn

      //- 创建时间
      +CreatedAtColumn

    //- 分页
    +Pagination
</template>

<script lang="ts">
import { defineComponent, onUnmounted } from "vue";

import useUserState, { userListLogin, userLoginClear } from "../states/user";
import {
  today,
  getDateDayShortcuts,
  formatBegin,
  formatEnd,
} from "../helpers/util";
import BaseFilter from "../components/base/Filter.vue";
import BaseTooltip from "../components/Tooltip.vue";
import TimeFormater from "../components/TimeFormater.vue";
import { PAGE_SIZES } from "../constants/common";
import FilterTable from "../mixins/FilterTable";

const defaultDateRange = [today(), today()];
const filterFields = [
  {
    label: "账号：",
    key: "account",
    placeholder: "请输入要查询的账号",
    clearable: true,
    span: 6,
  },
  {
    label: "时间：",
    key: "dateRange",
    type: "dateRange",
    placeholder: ["开始日期", "结束日期"],
    shortcuts: getDateDayShortcuts(["1d", "2d", "3d", "7d"]),
    defaultValue: defaultDateRange,
    span: 12,
  },
  {
    label: "",
    type: "filter",
    labelWidth: "0px",
    span: 6,
  },
];

export default defineComponent({
  name: "Logins",
  components: {
    BaseFilter,
    BaseTooltip,
    TimeFormater,
  },
  mixins: [FilterTable],
  setup() {
    // 必须以这种方式才会调用
    onUnmounted(() => {
      userLoginClear();
    });
    const userState = useUserState();
    return {
      logins: userState.logins,
    };
  },
  data() {
    return {
      filterFields,
      pageSizes: PAGE_SIZES,
      summaryFields: ["account", "ip", "sid", "tid", "location"],
      query: {
        dateRange: defaultDateRange,
        offset: 0,
        limit: PAGE_SIZES[0],
        account: "",
        order: "-createdAt",
      },
    };
  },
  methods: {
    async fetch() {
      const { query } = this;
      const params = Object.assign({}, query);
      const value = params.dateRange;
      if (value) {
        params.begin = formatBegin(value[0]);
        params.end = formatEnd(value[1]);
      } else {
        params.begin = "";
        params.end = "";
      }
      delete params.dateRange;
      try {
        await userListLogin(params);
      } catch (err) {
        this.$error(err);
      }
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "../common";

.logins
  margin $mainMargin
  i
    margin-right 5px
.pagination
  text-align right
  margin-top 15px
</style>
