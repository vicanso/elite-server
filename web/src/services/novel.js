import axios from "axios";

import { NOVELS, NOVELS_UPDATE, NOVELS_COVER_UPDATE } from "../urls";

// list 获取书籍信息列表
export async function list(params) {
  const { data } = await axios.get(NOVELS, {
    params: Object.assign(
      {
        nocache: true
      },
      params
    )
  });
  return data;
}

// updateByID 根据书籍ID更新书籍信息
export async function updateByID(id, params) {
  const url = NOVELS_UPDATE.replace(":id", id);
  const { data } = await axios.patch(url, params);
  return data;
}

// updateCoverByID 根据ID更新书籍封面
export async function updateCoverByID(id, imageURL) {
  const url = NOVELS_COVER_UPDATE.replace(":id", id);
  const { data } = await axios.patch(url, {
    cover: imageURL
  });
  return data;
}
