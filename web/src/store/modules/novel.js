import request from "@/helpers/request";

import {
  NOVEL_SOURCES,
  NOVELS,
  NOVELS_ID,
  NOVEL_CHAPTERS,
  NOVEL_CHAPTERS_UPDATE,
  NOVEL_CHAPTERS_ID,
  NOVEL_COVER
} from "@/constants/url";
import { addNoCacheQueryParam, formatDate } from "@/helpers/util";

const prefix = "novel";

const mutationNovelSourceList = `${prefix}.source.list`;
const mutationNovelSourceListProcessing = `${mutationNovelSourceList}.processing`;

const mutationNovelPublishing = `${prefix}.publishing`;
const mutationNovelPublished = `${prefix}.published`;

const mutationNovelList = `${prefix}.list`;
const mutationNovelListProcessing = `${mutationNovelList}.processing`;

const mutationNovelUpdate = `${prefix}.update`;
const mutationNovelUpdateProcessing = `${mutationNovelUpdate}.processing`;

const mutationNovelChapterList = `${prefix}.chapter.list`;
const mutationNovelChapterListProcessing = `${mutationNovelChapterList}.processing`;

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
  updateProcessing: false,

  // 是否正在拉取小说章节列表
  listChapterProcessing: false,
  chapterList: {
    data: null,
    count: -1
  }
};

function adjustNovelFields(item) {
  ["views", "downloads", "favorites", "wordCount"].forEach(key => {
    if (!item[key]) {
      item[key] = 0;
    }
  });
  item.coverURL =
    NOVEL_COVER.replace(":id", item.id) + "?quality=70&width=80&type=jpg";
  item.summaryCut = item.summary;
  const max = 40;
  if (item.summary || item.summary.length > max) {
    item.summaryCut = item.summary.substring(0, max) + "...";
  }
}

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
      novels.forEach(adjustNovelFields);
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
          item = adjustNovelFields(Object.assign(item, data));
        }
      });
      state.list.data = arr;
    },
    [mutationNovelChapterListProcessing](state, processing) {
      state.listChapterProcessing = processing;
    },
    [mutationNovelChapterList](state, { chapters = [], count = 0 }) {
      if (count >= 0) {
        state.chapterList.count = count;
      }
      chapters.forEach(item => {
        if (item.wordCount) {
          item.wordCountDesc = `${item.wordCount.toLocaleString()}字`;
        } else {
          item.wordCountDesc = "--";
        }
        item.contentDesc = item.content || "--";
        if (item.content && item.content.length > 20) {
          item.contentDesc = `${item.content.substring(0, 20)}...`;
        }
        item.updatedAtDesc = formatDate(item.updatedAt);
      });
      state.chapterList.data = chapters;
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
    // searchNovel
    async searchNovel(_, params) {
      const { data } = await request.get(NOVELS, {
        params: addNoCacheQueryParam(params)
      });
      return data;
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
    },
    // listNovelChapter
    async listNovelChapter({ commit }, { id, params }) {
      commit(mutationNovelChapterListProcessing, true);
      try {
        const { data } = await request.get(NOVEL_CHAPTERS.replace(":id", id), {
          params
        });
        commit(mutationNovelChapterList, data);
        return data;
      } finally {
        commit(mutationNovelChapterListProcessing, false);
      }
    },
    updateNovelChapters(_, id) {
      return request.post(NOVEL_CHAPTERS_UPDATE.replace(":id", id));
    },
    async getNovelChapterByID(_, id) {
      const { data } = await request.get(NOVEL_CHAPTERS_ID.replace(":id", id));
      data.no = data.no || 0;
      data.chapterNO = data.no + 1;
      const res = await request.get(NOVELS_ID.replace(":id", data.novel));
      if (res.data) {
        data.name = res.data.name;
        data.author = res.data.author;
      }
      return data;
    },
    async updateNovelChapterByID(_, { id, data }) {
      return request.patch(NOVEL_CHAPTERS_ID.replace(":id", id), data);
    },
    async updateNovelCoverByID(_, id) {
      return request.patch(NOVEL_COVER.replace(":id", id));
    }
  }
};
