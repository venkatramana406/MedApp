import Vue from "vue";
import VueRouter from "vue-router";
import login from "../views/login.vue";
import adduser from "../views/adduser.vue";

import stockview from "../views/stockview.vue";
import billentry from "../views/billentry.vue";
import stockentry from "../views/stockentry.vue"
import salesreport from "../views/salesreport.vue"
import loginhistory from "../views/loginhistory.vue"
import crypto from 'crypto-js';
import dash from "../views/dashboard.vue";


Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "login",
    component: login,
    meta:{auth:false }//authentication not required
  },
  {
    path: "/adduser",
    name: "adduser",
    component: adduser,
    meta:{auth:true ,role:['System Admin']},//authentication  required
  },
  {
    path: "/stockview",
    name:"stockview",
    component: stockview,
    meta:{auth:true ,role:['Manager','Biller','Inventory']},
  },
  {
    path: "/billentry",
    name:"billentry",
    component: billentry,
    meta:{auth:true ,role:['Biller']},
  },
  {
    path: "/stockentry",
    name:"stockentry",
    component: stockentry,
    meta:{auth:true ,role:['Manager','Inventory']},
  },
  {
    path: "/salesreport",
    name:"salesreport",
    component: salesreport,
    meta:{auth:true ,role:['Manager']},
  },
  {
    path: "/loginhistory",
    name:"loginhistory",
    component: loginhistory,
    meta:{auth:true ,role:['System Admin']},
  },
  {
    path:"/dashboard",
    name:"dash",
    component:dash,
    meta:{auth:true ,role:['Manager','Biller','System Admin','Inventory']},
  },
  
  
];

const router = new VueRouter({
  routes,
});

router.beforeEach((to,from,next)=>{
  if (to.matched.some(record=>record.meta.auth)) {
    let user;
    let userData=localStorage.getItem('user');
    if (userData) {
      const encryptKey="venkat434";
      const decryptData=crypto.AES.decrypt(userData,encryptKey).toString(crypto.enc.Utf8);
       user=JSON.parse(decryptData)
    }
    // console.log(user);
    if (!user) {
      next({name:"login"})
    }else{
      let role=user.role
      if (to.matched.some(record=>record.meta.role.includes(role))) {
        next()
      }else{
        next({name:"dash"})
      }
    }
  } 
  else {
    next();
  }
})

export default router;
