import request from "@/helpers/request";

import { NOVEL_SOURCES, NOVELS, NOVELS_ID } from "@/constants/url";
import { addNoCacheQueryParam } from "@/helpers/util";

const prefix = "novel";

const mutationNovelSourceList = `${prefix}.source.list`;
const mutationNovelSourceListProcessing = `${mutationNovelSourceList}.processing`;

const mutationNovelPublishing = `${prefix}.publishing`;
const mutationNovelPublished = `${prefix}.published`;

const mutationNovelList = `${prefix}.list`;
const mutationNovelListProcessing = `${mutationNovelList}.processing`;

const mutationNovelUpdate = `${prefix}.update`;
const mutationNovelUpdateProcessing = `${mutationNovelUpdate}.processing`;

const state = {
  // 是否正在发布小说
  publishing: false,

  // 是否正在拉取小说源列表
  sourceListProcessing: false,
  // 小说源列表数据
  sourceList: {
    data: null,
    count: -1
  },

  // 是否正在拉取小说列表
  listProcessing: false,
  // 小说列表数据
  list: {
    data: null,
    count: -1
  },
  updateProcessing: false
};

export default {
  state,
  mutations: {
    [mutationNovelSourceListProcessing](state, processing) {
      state.sourceListProcessing = processing;
    },
    [mutationNovelSourceList](state, { novelSources = [], count = 0 }) {
      if (count >= 0) {
        state.sourceList.count = count;
      }
      state.sourceList.data = novelSources;
    },
    [mutationNovelPublishing](state, processing) {
      state.publishing = processing;
    },
    [mutationNovelPublished](state, { name, author }) {
      const arr = state.sourceList.data.slice(0);
      arr.forEach(item => {
        if (item.name === name && item.author === author) {
          item.status = 2;
        }
      });
      state.sourceList.data = arr;
    },
    [mutationNovelListProcessing](state, processing) {
      state.listProcessing = processing;
    },
    [mutationNovelList](state, { novels = [], count = 0 }) {
      if (count >= 0) {
        state.list.count = count;
      }
      state.list.data = novels;
    },
    [mutationNovelUpdateProcessing](state, processing) {
      state.updateProcessing = processing;
    },
    [mutationNovelUpdate](state, { id, data }) {
      if (!state.list.data) {
        return;
      }
      const arr = state.list.data.slice(0);
      arr.forEach(item => {
        if (item.id === id) {
          item = Object.assign(item, data);
        }
      });
      state.list.data = arr;
    }
  },
  actions: {
    // listNovelSource 获取小说源列表
    async listNovelSource({ commit }, params) {
      commit(mutationNovelSourceListProcessing, true);
      try {
        const { data } = await request.get(NOVEL_SOURCES, {
          params: addNoCacheQueryParam(params)
        });
        commit(mutationNovelSourceList, data);
      } finally {
        commit(mutationNovelSourceListProcessing, false);
      }
    },
    // publishNovel 发布小说
    async publishNovel({ commit }, params) {
      commit(mutationNovelPublishing, true);
      try {
        const { data } = await request.post(NOVELS, params);
        commit(mutationNovelPublished, data);
        return data;
      } finally {
        commit(mutationNovelPublishing, false);
      }
    },
    // listNovel 获取小说
    async listNovel({ commit }, params) {
      commit(mutationNovelListProcessing, true);
      try {
        const { data } = await request.get(NOVELS, {
          params: addNoCacheQueryParam(params)
        });
        commit(mutationNovelList, data);
        return data;
      } finally {
        commit(mutationNovelListProcessing, false);
      }
    },
    async getNovelByID(_, id) {
      const { data } = await request.get(NOVELS_ID.replace(":id", id));
      return data;
    },
    async updateNovelByID({ commit }, { id, data }) {
      commit(mutationNovelUpdateProcessing, true);
      try {
        await request.patch(NOVELS_ID.replace(":id", id), data);
        commit(mutationNovelUpdate, {
          id,
          data
        });
      } finally {
        commit(mutationNovelUpdateProcessing, false);
      }
    }
  }
};
