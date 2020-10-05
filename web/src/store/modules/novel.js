import request from "@/helpers/request";

import { NOVEL_SOURCES, NOVELS } from "@/constants/url";

const prefix = "novel";

const mutationNovelSourceList = `${prefix}.source.list`;
const mutationNovelSourceListProcessing = `${mutationNovelSourceList}.processing`;

const mutationNovelPublishing = `${prefix}.publishing`;
const mutationNovelPublished = `${prefix}.published`;

const mutationNovelList = `${prefix}.list`;
const mutationNovelListProcessing = `${mutationNovelList}.processing`;

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
  }
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
    }
  },
  actions: {
    async listNovelSource({ commit }, params) {
      commit(mutationNovelSourceListProcessing, true);
      try {
        const { data } = await request.get(NOVEL_SOURCES, {
          params
        });
        commit(mutationNovelSourceList, data);
      } finally {
        commit(mutationNovelSourceListProcessing, false);
      }
    },
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
    async listNovel({ commit }, params) {
      commit(mutationNovelListProcessing, true);
      try {
        const { data } = await request.get(NOVELS, {
          params
        });
        commit(mutationNovelList, data);
        return data;
      } finally {
        commit(mutationNovelListProcessing, false);
      }
    }
  }
};
