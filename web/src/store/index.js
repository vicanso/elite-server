import Vue from "vue";
import Vuex from "vuex";

import userStore from "@/store/modules/user";
import configStore from "@/store/modules/config";
import commonStore from "@/store/modules/common";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    user: userStore,
    config: configStore,
    common: commonStore
  }
});