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
  createdAt: string;
  updatedAt: string;
  name: string;
  author: string;
  source: number;
  status: number;
  cover: string;
  summary: string;
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

interface NovelStatuses {
  items: string[];
}
const statues: NovelStatuses = reactive({
  items: ["未知", "连载中", "已完结", "禁止状态"],
});

interface ReadonlyNovelState {
  novels: DeepReadonly<Novels>;
  statuses: DeepReadonly<NovelStatuses>;
}

function fillInfo(item: Novel) {
  item.id = item.id || 0;
  item.status = item.status || 0;
}

// novelList 获取小说列表
export async function novelList(params: {
  keyword?: string;
  limit: number;
  offset: number;
  fields?: string;
}): Promise<void> {
  if (novels.processing) {
    return;
  }
  try {
    novels.processing = true;
    const { data } = await request.get(NOVELS, {
      params,
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

const state: ReadonlyNovelState = {
  novels: readonly(novels),
  statuses: readonly(statues),
};

// useNovelState 用户小说相关state
export default function useNovelState(): ReadonlyNovelState {
  return state;
}
