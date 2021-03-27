// 用户相关url
// 用户信息
export const USERS_ME = "/users/v1/me";
// 用户详细信息
export const USERS_ME_DETAIL = "/users/v1/detail";
// 用户登录
export const USERS_LOGIN = "/users/v1/me/login";
export const USERS_INNER_LOGIN = "/users/inner/v1/me/login";
// 用户行为
export const USERS_ACTIONS = "/users/v1/actions";
// 用户登录记录
export const USERS_LOGINS = "/users/v1/login-records";
// 用户角色列表
export const USERS_ROLES = "/users/v1/roles";
// 用户列表
export const USERS = "/users/v1";
// 根据ID查询用户信息
export const USERS_ID = "/users/v1/:id";

// flux相关查询
// 用户行为日志列表
export const FLUXES_TRACKERS = "/fluxes/v1/trackers";
// http出错列表
export const FLUXES_HTTP_ERRORS = "/fluxes/v1/http-errors";
// 客户端上传的action日志列表
export const FLUXES_ACTIONS = "/fluxes/v1/actions";
// 后端HTTP调用列表
export const FLUXES_REQUESTS = "/fluxes/v1/requests";
// tag value列表
export const FLUXES_TAG_VALUES = "/fluxes/v1/tag-values/:measurement/:tag";

// 通用接口相关url
// 图形验证码
export const COMMONS_CAPTCHA = "/commons/captcha";
// schema状态列表
export const COMMONS_STATUSES = "/commons/schema-statuses";
// 路由列表
export const COMMONS_ROUTERS = "/commons/routers";
// 随机字符串
export const COMMONS_RANDOM_KEYS = "/commons/random-keys";
// HTTP性能指标统计
export const COMMONS_HTTP_STATS = "/commons/http-stats";

// 系统配置相关url
// 配置列表
export const CONFIGS = "/configurations/v1";
// 根据ID查询或更新配置
export const CONFIGS_ID = "/configurations/v1/:id";
// 当前有效配置
export const CONFIGS_CURRENT_VALID = "/configurations/v1/current-valid";

// 小说相关url
// 小说列表
export const NOVELS = "/novels/v1";
// 通过ID查找小说
export const NOVELS_ID = "/novels/v1/:id";
// 通过ID查找小说章节
export const NOVELS_CHAPTERS = "/novels/v1/:id/chapters";
// 通过ID与章节序号查找章节内容
export const NOVELS_CHAPTERS_ID = "/novels/v1/:id/chapters/:no";

// 管理员相关接口
export const ADMINS_SESSION_ID = "/@admin/v1/sessions/:id";
