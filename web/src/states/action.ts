import { userActions } from "../store";
import request from "../helpers/request";
import { USERS_ACTIONS } from "../constants/url";

// 定时flush的间隔
const flushInterval = 60 * 1000;
let timer: number;

// UserActionData 用户action信息
interface UserActionData {
  // 以上的为基本信息，所有记录都有
  category: string;
  route: string;
  path: string;
  result: number;
  time?: number;

  // 以下的为额外信息
  // 出错信息
  message?: string;
}

userActions.load();

// 默认的key，除此之外的都添加至extra
const defaultKeys: string[] = ["category", "route", "path", "result", "time"];
async function flush() {
  const data = await userActions.clear();
  if (!data || data.length === 0) {
    return;
  }
  // 需要将actions转换
  const actions = data.map((action) => {
    const extra: Record<string, unknown> = {};
    const newAction: Record<string, unknown> = {
      extra,
    };
    Object.keys(action).forEach((element) => {
      const value = action[element];
      if (defaultKeys.includes(element)) {
        newAction[element] = value;
      } else {
        extra[element] = value;
      }
    });
    return newAction;
  });

  request.post(USERS_ACTIONS, {
    actions,
  });
}

// 成功
export const SUCCESS = 0;
// 失败
export const FAIL = 1;

// 分类
// 点击（默认分类）
export const CLICK = "click";
// 登录
export const LOGIN = "login";
// 注册
export const REGISTER = "register";
// 路由跳转
export const ROUTE_CHANGE = "routeChange";
// 出错（仅用于异常类的出错，如undefined的使用异常，请求中断等）
export const ERROR = "error";

export function actionAdd(data: UserActionData): void {
  // 每次添加新的action时，清空定时器
  clearTimeout(timer);
  userActions.add(data);
  // 如果为0在添加时填充
  if (!data.time) {
    data.time = Math.floor(Date.now() / 1000);
  }
  if (!data.path) {
    data.path = "/";
  }
  if (!data.route) {
    data.route = "unknown";
  }
  if (userActions.size < 10) {
    // 重新启动定时器
    timer = setTimeout(flush, flushInterval);
    return;
  }
  flush();
}
