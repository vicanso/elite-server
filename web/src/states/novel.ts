import { reactive, readonly, DeepReadonly } from "vue";

import request from "../helpers/request";

import {
  NOVELS,
  NOVELS_ID,
  NOVELS_CHAPTERS,
  NOVELS_CHAPTERS_ID,
} from "../constants/url";

// 小说信息

interface Novel {
  id: number;
  name: string;
  author: string;
  createdAt?: string;
  updatedAt?: string;
  source?: number;
  status?: number;
  cover?: string;
  summary?: string;
  wordCount?: number;
}
interface Novels {
  processing: boolean;
  count: number;
  items: Novel[];
}
const novels: Novels = reactive({
  processing: false,
  count: -1,
  items: [],
});

interface Chapter {
  id: number;
  title: string;
  no: number;
  content?: string;
  wordCount?: number;
  createdAt?: string;
  updatedAt?: string;
}
interface Chapters {
  novel: number;
  processing: boolean;
  count: number;
  items: Chapter[];
}

const chapters: Chapters = reactive({
  novel: 0,
  processing: false,
  count: -1,
  items: [],
});

interface NovelStatuses {
  items: string[];
}
const statues: NovelStatuses = reactive({
  items: ["未知", "连载中", "已完结", "禁止状态"],
});

interface NovelDetail {
  processing: boolean;
  data: Novel;
}
const detail: NovelDetail = {
  processing: false,
  data: {
    id: 0,
    name: "",
    author: "",
  },
};

interface ReadonlyNovelState {
  novels: DeepReadonly<Novels>;
  statuses: DeepReadonly<NovelStatuses>;
  detail: DeepReadonly<NovelDetail>;
  chapters: DeepReadonly<Chapters>;
}

function fillInfo(item: Novel): Novel {
  item.id = item.id || 0;
  item.status = item.status || 0;
  return item;
}

function fillChapterInfo(item: Chapter): Chapter {
  item.no = item.no || 0;
  return item;
}

// novelList 获取小说列表
export async function novelList(params: {
  keyword?: string;
  limit: number;
  offset: number;
  fields?: string;
  mustCount?: string;
  ignoreCount?: string;
}): Promise<void> {
  if (novels.processing) {
    return;
  }
  try {
    novels.processing = true;
    const listParams = Object.assign({}, params);
    // 如果总数为-1（从其它返回或直接刷新 ），则强制获取总数
    if (novels.count === -1) {
      listParams.mustCount = "1";
    }
    const { data } = await request.get(NOVELS, {
      params: listParams,
    });
    const count = data.count || 0;
    if (count >= 0) {
      novels.count = count;
    }
    novels.items = data.novels || [];
    novels.items.forEach(fillInfo);
  } finally {
    novels.processing = false;
  }
}

// novelListClear 清空小说列表记录
export function novelListClear(): void {
  novels.count = -1;
  novels.items.length = 0;
}

// novelGetDetail 获取小说详情信息
export async function novelGetDetail(id: number): Promise<void> {
  if (detail.processing) {
    return;
  }
  // 如果在列表中能获取，则直接使用获取值
  const found = novels.items.find((item) => item.id === id);
  if (found) {
    Object.assign(detail.data, found);
    return;
  }
  // 调用接口查询
  try {
    detail.processing = true;
    const { data } = await request.get(NOVELS_ID.replace(":id", `${id}`));
    Object.assign(detail.data, fillInfo(data));
  } finally {
    detail.processing = false;
  }
}

// novelUpdateDetail 更新小说信息
export async function novelUpdateDetail(
  id: number,
  data: {
    status?: number;
    summary?: string;
  }
): Promise<void> {
  if (detail.processing) {
    return;
  }
  try {
    if (detail.data.id === id) {
      detail.processing = true;
    }
    await request.patch(NOVELS_ID.replace(":id", `${id}`), data);
    Object.assign(detail.data, data);
  } finally {
    if (detail.data.id === id) {
      detail.processing = false;
    }
  }
}

// novelListChapter 获取小说章节
export async function novelListChapter(
  novelID: number,
  params: {
    limit: number;
    offset: number;
    fields?: string;
    mustCount?: string;
    ignoreCount?: string;
  }
): Promise<void> {
  if (chapters.processing) {
    return;
  }
  try {
    chapters.novel = novelID;
    chapters.processing = true;
    const listParams = Object.assign({}, params);
    // 如果总数为-1（从其它返回或直接刷新 ），则强制获取总数
    if (chapters.count === -1) {
      listParams.mustCount = "1";
    }
    const { data } = await request.get(
      NOVELS_CHAPTERS.replace(":id", `${novelID}`),
      {
        params: listParams,
      }
    );
    const count = data.count || 0;
    if (count >= 0) {
      chapters.count = count;
    }
    chapters.items = data.chapters || [];
    chapters.items.forEach(fillChapterInfo);
  } finally {
    chapters.processing = false;
  }
}

// novelChaptersClear 清空小说章节
export function novelChaptersClear(): void {
  chapters.count = -1;
  chapters.items.length = 0;
}

const state: ReadonlyNovelState = {
  novels: readonly(novels),
  statuses: readonly(statues),
  detail: readonly(detail),
  chapters: readonly(chapters),
};

// useNovelState 用户小说相关state
export default function useNovelState(): ReadonlyNovelState {
  return state;
}
