import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import vuetify from "./plugins/vuetify";
import dashmanager from "./components/dashmanager.vue";
import dashbiller from "./components/dashbiller.vue";
import dashsysadmin from "./components/dashsysadmin.vue";
import dashinventory from "./components/dashinventory.vue"
Vue.config.productionTip = false;
Vue.component("dashmanager",dashmanager);
Vue.component("dashbiller",dashbiller);
Vue.component("dashsysadmin",dashsysadmin);
Vue.component("dashinventory",dashinventory);

new Vue({
  router,
  store,
  vuetify,
  render: (h) => h(App),
}).$mount("#app");
