<template lang="pug">
//- 小说名称
mixin NameField
  el-col(
    :span="8"
  ): el-form-item(
    label="名称："
  ) {{detail.data.name}}

//- 小说作者
mixin AuthorField
  el-col(
    :span="8"
  ): el-form-item(
    label="作者："
  ) {{detail.data.author}}

//- 小说状态
mixin StatusField
  el-col(
    :span="8"
  ): el-form-item(
    label="状态："
  ): el-select.fullFill(
    v-model="form.status"
  ): el-option(
    v-for="(item, index) in statuses.items"
    :key="item"
    :label="item"
    :value="index"
  )

//- 小说总字数
mixin WordCountField
  el-col(
    :span="8"
  ): el-form-item(
    label="总字数："
  ) {{detail.data.wordCount && detail.data.wordCount.toLocaleString()}}

//- 小说更新时间
mixin UpdatedAtField
  el-col(
    :span="8"
  ): el-form-item(
    label="更新于："
  ): time-formater(
    :time="detail.data.updatedAt"
  )

//- 小说简介
mixin SummaryField
  el-col(
    :span="24"
  ): el-form-item(
    label="简介："
  ): el-input(
    type="textarea"
    :autosize="{ minRows: 4, maxRows: 8}"
    placeholder="请输入小说简介"
    v-model="form.summary"
  )

.novelDetail
  el-card(
    v-loading="detail.processing"
  )
    template(
      #header
    )
      a.mright10.bold(
        @click.prevent="$router.back()"
        href="#"
      ): i.el-icon-arrow-left
      | 小说详情
    el-form(
      ref="form"
      v-if="detail.data.id"
      :model="form"
      label-width="100px"
    ): el-row(
      :gutter="20"
    )
      +NameField

      +AuthorField

      +WordCountField

      +UpdatedAtField

      +StatusField

      +SummaryField

  el-button.fullFill.mtopMain(
    type="primary"
    @click="update"
  ) 更新
</template>

<script lang="ts">
import { defineComponent } from "vue";

import useNovelState, {
  novelGetDetail,
  novelUpdateDetail,
} from "../../states/novel";
import TimeFormater from "../../components/TimeFormater.vue";
import { diff } from "../../helpers/util";

export default defineComponent({
  name: "NovelDetail",
  components: {
    TimeFormater,
  },
  setup() {
    const novelState = useNovelState();
    return {
      detail: novelState.detail,
      statuses: novelState.statuses,
    };
  },
  data() {
    return {
      form: {
        status: 0,
        summary: "",
      },
    };
  },
  beforeMount() {
    this.fetch();
  },
  methods: {
    async fetch() {
      const { id } = this.$route.params;
      try {
        const { form, detail } = this;
        await novelGetDetail(Number(id));
        form.status = detail.data.status;
        form.summary = detail.data.summary;
        // 由于数据均在state中，不知为啥无法更新组件，
        // 因此强制更新
        this.$forceUpdate();
      } catch (err) {
        this.$error(err);
      }
    },
    // 更新小说信息
    async update() {
      const { form, detail } = this;
      const { modifiedCount, data } = diff(form, detail.data);
      if (modifiedCount === 0) {
        this.$message.warning("请先修改要更新的信息");
        return;
      }
      if (data.status === 0) {
        this.$message.warning("状态不允许设置为未知");
        return;
      }
      try {
        await novelUpdateDetail(detail.data.id, data);
        this.$message.info("已成功更新信息");
      } catch (err) {
        this.$error(err);
      }
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "../../common";
.novelDetail
  margin $mainMargin
</style>
