import { createRouter, createWebHashHistory } from "vue-router";

import { actionAdd, ROUTE_CHANGE, SUCCESS } from "../states/action";

import Home from "../views/Home.vue";
import Profile from "../views/Profile.vue";
import Login from "../views/Login.vue";
import Register from "../views/Register.vue";
import Logins from "../views/Logins.vue";
import Users from "../views/Users.vue";
import Trackers from "../views/Trackers.vue";
import Actions from "../views/Actions.vue";
import HTTPErrors from "../views/HTTPErrors.vue";
import Requests from "../views/Requests.vue";

// 系统配置
import ApplicationSetting from "../views/configs/ApplicationSetting.vue";
import MockTime from "../views/configs/MockTime.vue";
import BlockIP from "../views/configs/BlockIP.vue";
import SignedKey from "../views/configs/SignedKey.vue";
import RouterMock from "../views/configs/Router.vue";
import RouterConcurrency from "../views/configs/RouterConcurrency.vue";
import RequestConcurrency from "../views/configs/RequestConcurrency.vue";
import SessionInterceptor from "../views/configs/SessionInterceptor.vue";
import Configuration from "../views/configs/Configuration.vue";
import Others from "../views/Others.vue";

// 小说相关页面
import Novels from "../views/novels/Novels.vue";
import NovelDetail from "../views/novels/NovelDetail.vue";
import NovelChapters from "../views/novels/NovelChapters.vue";
import NovelChapterDetail from "../views/novels/NovelChapterDetail.vue";

export const ROUTE_HOME = "home";
export const ROUTE_PROFILE = "profile";
export const ROUTE_LOGIN = "login";
export const ROUTE_REGISTER = "register";
export const ROUTE_LOGINS = "logins";
export const ROUTE_USERS = "users";
export const ROUTE_TRACKERS = "trackers";
export const ROUTE_ACTIONS = "actions";
export const ROUTE_HTTP_ERRORS = "httpErrors";
export const ROUTE_REQUESTS = "requests";

// 系统配置
export const ROUTE_APPLICATION_SETTING = "applicationSetting";
export const ROUTE_MOCK_TIME = "mockTime";
export const ROUTE_BLOCK_IP = "blockIP";
export const ROUTE_SIGNED_KEY = "signedKey";
export const ROUTE_ROUTER_MOCK = "routerMock";
export const ROUTE_ROUTER_CONCURRENCY = "routerConcurrency";
export const REQUEST_CONCURRENCY = "requestConcurrency";
export const ROUTE_SESSION_INTERCEPTOR = "sessionInterceptor";
export const ROUTE_CONFIGURATION = "configuration";
export const ROUTE_OTHERS = "others";

// 小说相关
export const NOVEL_LIST = "novels";
export const NOVEL_DETAIl = "novelDetail";
export const NOVEL_CHAPTERS = "novelChapters";
export const NOVEL_CHAPTER_DETAIL = "novelChapterDetail";

interface Location {
  name: string;
  path: string;
}

const currentLocation: Location = {
  name: "",
  path: "",
};
const prevLocation: Location = {
  name: "",
  path: "",
};

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      name: ROUTE_HOME,
      component: Home,
    },
    {
      path: "/profile",
      name: ROUTE_PROFILE,
      component: Profile,
    },
    {
      path: "/login",
      name: ROUTE_LOGIN,
      component: Login,
    },
    {
      path: "/register",
      name: ROUTE_REGISTER,
      component: Register,
    },
    {
      path: "/users",
      name: ROUTE_USERS,
      component: Users,
    },
    {
      path: "/logins",
      name: ROUTE_LOGINS,
      component: Logins,
    },
    {
      path: "/trackers",
      name: ROUTE_TRACKERS,
      component: Trackers,
    },
    {
      path: "/actions",
      name: ROUTE_ACTIONS,
      component: Actions,
    },
    {
      path: "/http-errors",
      name: ROUTE_HTTP_ERRORS,
      component: HTTPErrors,
    },
    {
      path: "/requests",
      name: ROUTE_REQUESTS,
      component: Requests,
    },
    {
      path: "/application-settings",
      name: ROUTE_APPLICATION_SETTING,
      component: ApplicationSetting,
    },
    {
      path: "/mock-time",
      name: ROUTE_MOCK_TIME,
      component: MockTime,
    },
    {
      path: "/block-ip",
      name: ROUTE_BLOCK_IP,
      component: BlockIP,
    },
    {
      path: "/signed-key",
      name: ROUTE_SIGNED_KEY,
      component: SignedKey,
    },
    {
      path: "/router-mock",
      name: ROUTE_ROUTER_MOCK,
      component: RouterMock,
    },
    {
      path: "/router-concurrency",
      name: ROUTE_ROUTER_CONCURRENCY,
      component: RouterConcurrency,
    },
    {
      path: "/request-concurrency",
      name: REQUEST_CONCURRENCY,
      component: RequestConcurrency,
    },
    {
      path: "/session-interceptor",
      name: ROUTE_SESSION_INTERCEPTOR,
      component: SessionInterceptor,
    },
    {
      path: "/configuration",
      name: ROUTE_CONFIGURATION,
      component: Configuration,
    },
    {
      path: "/others",
      name: ROUTE_OTHERS,
      component: Others,
    },
    {
      path: "/novels",
      name: NOVEL_LIST,
      component: Novels,
    },
    {
      path: "/novels/:id",
      name: NOVEL_DETAIl,
      component: NovelDetail,
    },
    {
      path: "/novels/:id/chapters",
      name: NOVEL_CHAPTERS,
      component: NovelChapters,
    },
    {
      path: "/novels/:id/chapters/:no",
      name: NOVEL_CHAPTER_DETAIL,
      component: NovelChapterDetail,
    },
  ],
});

export function getCurrentLocation(): Location {
  return currentLocation;
}

router.beforeEach((to, from) => {
  if (from.name) {
    prevLocation.name = from.name.toString();
    prevLocation.path = from.fullPath;
  }
  if (to.name) {
    currentLocation.name = to.name.toString();
    currentLocation.path = to.fullPath;
  }
  actionAdd({
    category: ROUTE_CHANGE,
    route: currentLocation.name,
    path: currentLocation.path,
    result: SUCCESS,
    time: Math.floor(Date.now() / 1000),
  });
});

export default router;
