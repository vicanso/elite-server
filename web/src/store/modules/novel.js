import request from "@/helpers/request";

import { NOVEL_SOURCES, NOVELS } from "@/constants/url";

const prefix = "novel";

const mutationNovelSourceList = `${prefix}.source.list`;
const mutationNovelSourceListProcessing = `${mutationNovelSourceList}.processing`;

const mutationNovelPublishing = `${prefix}.publishing`;

const state = {
  // 是否正在发布小说
  novelPublishing: false,

  novelSourceListProcessing: false,
  novelSourceList: {
    data: null,
    count: -1
  }
};

export default {
  state,
  mutations: {
    [mutationNovelSourceListProcessing](state, processing) {
      state.novelSourceListProcessing = processing;
    },
    [mutationNovelSourceList](state, { novelSources = [], count = 0 }) {
      if (count >= 0) {
        state.novelSourceList.count = count;
      }
      state.novelSourceList.data = novelSources;
    },
    [mutationNovelPublishing](state, processing) {
      state.novelPublishing = processing;
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
        return data;
      } finally {
        commit(mutationNovelPublishing, false);
      }
    }
  }
};
