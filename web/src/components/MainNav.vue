<template lang="pug">
//- 切换侧边栏
mixin ToggleNav
  a.toggleNav(
    href="#"
    @click.stop="toggleNav"
  )
    i(
      :class=`$props.shrinking ? "el-icon-s-unfold" : "el-icon-s-fold"`
    )

//- 应用图标
mixin HomeLogo
  h1
    router-link(
      v-if="!$props.shrinking"
      :to='{name: homeRoute}'
    )
      i.el-icon-cpu
      | Elite

//- 菜单栏
mixin Menu
  nav: el-menu.menu(
    :collapse="$props.shrinking"
    :default-active="active"
    background-color="#000c17"
    text-color="#fff"
    active-text-color="#fff"
  )
    el-submenu.submenu(
      v-for="(nav, i) in navs"
      :index="`${i}`"
      :key="`${i}`"
    )
      template(
        #title
      )
        i(
          :class="nav.icon"
        )
        span {{nav.name}}
      //- 子菜单栏
      el-menu-item.menuItem(
        v-for="(subItem, j) in nav.children"
        :index="`${i}-${j}`"
        :key="`${i}-${j}`"
        @click="goTo(subItem)"
      )
        span {{subItem.name}}
.mainNav
  //- 切换侧边栏
  +ToggleNav

  //- 应用图标
  +HomeLogo

  //- 菜单栏
  +Menu

</template>

<script lang="ts">
import { defineComponent } from "vue";
import {
  ROUTE_HOME,
  ROUTE_LOGINS,
  ROUTE_USERS,
  ROUTE_TRACKERS,
  ROUTE_APPLICATION_SETTING,
  ROUTE_MOCK_TIME,
  ROUTE_BLOCK_IP,
  ROUTE_SIGNED_KEY,
  ROUTE_ROUTER_MOCK,
  ROUTE_ROUTER_CONCURRENCY,
  REQUEST_CONCURRENCY,
  ROUTE_SESSION_INTERCEPTOR,
  ROUTE_CONFIGURATION,
  ROUTE_OTHERS,
  ROUTE_HTTP_ERRORS,
  ROUTE_ACTIONS,
  ROUTE_REQUESTS,
  NOVEL_LIST,
} from "../router";
import { USER_ADMIN, USER_SU } from "../constants/user";
import useUserState from "../states/user";
import { isAllowedUser } from "../helpers/util";

const navs = [
  {
    name: "小说",
    icon: "el-icon-notebook-1",
    roles: [USER_ADMIN, USER_SU],
    groups: [],
    children: [
      {
        name: "小说列表",
        route: NOVEL_LIST,
        roules: [],
        groups: [],
      },
    ],
  },
  {
    name: "用户",
    icon: "el-icon-user",
    roles: [USER_ADMIN, USER_SU],
    groups: [],
    children: [
      {
        name: "用户列表",
        route: ROUTE_USERS,
        roles: [],
        groups: [],
      },
      {
        name: "登录记录",
        route: ROUTE_LOGINS,
        roles: [],
        groups: [],
      },
      {
        name: "用户行为",
        route: ROUTE_TRACKERS,
        roles: [],
        groups: [],
      },
    ],
  },
  {
    name: "配置",
    icon: "el-icon-setting",
    roles: [USER_SU],
    groups: [],
    children: [
      {
        name: "所有配置",
        route: ROUTE_CONFIGURATION,
        roles: [],
        groups: [],
      },
      {
        name: "应用配置",
        route: ROUTE_APPLICATION_SETTING,
        roles: [],
        groups: [],
      },
      {
        name: "MockTime配置",
        route: ROUTE_MOCK_TIME,
        roles: [],
        groups: [],
      },
      {
        name: "黑名单IP",
        route: ROUTE_BLOCK_IP,
        roles: [],
        groups: [],
      },
      {
        name: "SignedKey配置",
        route: ROUTE_SIGNED_KEY,
        roles: [],
        groups: [],
      },
      {
        name: "路由Mock配置",
        route: ROUTE_ROUTER_MOCK,
        roles: [],
        groups: [],
      },
      {
        name: "路由并发配置",
        route: ROUTE_ROUTER_CONCURRENCY,
        roles: [],
        groups: [],
      },
      {
        name: "HTTP实例并发配置",
        route: REQUEST_CONCURRENCY,
        roles: [],
        groups: [],
      },
      {
        name: "Session拦截配置",
        route: ROUTE_SESSION_INTERCEPTOR,
        roles: [],
        groups: [],
      },
    ],
  },
  {
    name: "其它",
    icon: "el-icon-set-up",
    roles: [USER_SU],
    groups: [],
    children: [
      {
        name: "响应出错记录",
        route: ROUTE_HTTP_ERRORS,
        roles: [],
        groups: [],
      },
      {
        name: "后端HTTP调用",
        route: ROUTE_REQUESTS,
        roles: [],
        groups: [],
      },
      {
        name: "客户端行为记录",
        route: ROUTE_ACTIONS,
        roules: [],
        groups: [],
      },
      {
        name: "其它",
        route: ROUTE_OTHERS,
        roles: [],
        groups: [],
      },
    ],
  },
];

export default defineComponent({
  name: "MainNav",
  props: {
    shrinking: {
      type: Boolean,
      default: false,
    },
    onToggle: {
      type: Function,
      default: null,
    },
  },
  emits: ["toggle"],

  setup() {
    const userState = useUserState();
    return {
      user: userState.info,
    };
  },
  data() {
    return {
      homeRoute: ROUTE_HOME,
      active: "",
    };
  },
  computed: {
    navs() {
      const { user } = this;
      if (!user || !user.account) {
        return [];
      }
      const { roles, groups } = user;
      const filterNavs = [];
      navs.forEach((item) => {
        // 如果该栏目有配置权限，而且用户无该权限
        if (item.roles && !isAllowedUser(item.roles, roles)) {
          return;
        }
        // 如果该栏目配置了允许分级，而该用户不属于该组
        if (item.groups && !isAllowedUser(item.groups, groups)) {
          return;
        }
        const clone = Object.assign({}, item);
        const children = item.children.map((subItem) =>
          Object.assign({}, subItem)
        );
        clone.children = children.filter((subItem) => {
          // 如果未配置色色与分组限制
          if (!subItem.roles && !subItem.groups) {
            return true;
          }
          if (subItem.roles && !isAllowedUser(subItem.roles, roles)) {
            return false;
          }
          if (subItem.groups && !isAllowedUser(subItem.groups, groups)) {
            return false;
          }
          return true;
        });
        filterNavs.push(clone);
      });
      return filterNavs;
    },
  },
  watch: {
    // 如果nav变化时，根据当前route定位
    navs() {
      this.updateActive(this.$route.name);
    },
    // 路由变化时设置对应的导航为活动状态
    $route(to) {
      this.updateActive(to.name);
    },
  },
  beforeMount() {
    this.updateActive(this.$route.name);
  },
  methods: {
    toggleNav() {
      this.$emit("toggle");
    },
    goTo({ route }) {
      if (!route || this.$route.name === route) {
        return;
      }
      this.$router.push({
        name: route,
      });
    },
    // 查询定位当前选中菜单
    updateActive(name) {
      const { navs } = this;
      let active = "";
      navs.forEach((nav, i) => {
        nav.children.forEach((item, j) => {
          if (item.route === name) {
            active = `${i}-${j}`;
          }
        });
      });
      this.active = active;
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "../common";
$mainNavColor = #000c17
.mainNav
  min-height 100vh
  overflow-y auto
  background-color $mainNavColor
.toggleNav
  height $mainHeaderHeight
  line-height $mainHeaderHeight
  display block
  float right
  width $mainNavShrinkingWidth
  text-align center
h1
  height $mainHeaderHeight
  line-height $mainHeaderHeight
  color $white
  padding-left 20px
  font-size 18px
  margin-right $mainNavShrinkingWidth
  i
    font-weight bold
    margin-right 5px
nav
  border-top 1px solid rgba($white, 0.3)
.menu
  border-right 1px solid $mainNavColor
.menuItem
  color rgba($white, 0.65)
  &.is-active
    background-color $darkBlue !important
</style>
