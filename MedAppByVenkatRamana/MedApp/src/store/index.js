import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    user:null,
  },
  mutations: {
    UPDATE_USER(state,val){
      state.user=val
    }
    
  },
  actions: {},
  modules: {},
  });
