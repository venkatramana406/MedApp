<template>
  <v-app >
    <v-app-bar class="navigation" color="pink accent-2 " max-height="65">
      <h1 class="text-h2 font-weight-regular deep-red--text">Med</h1>
      <h1 class="text-h4 font-weight-light black--text">App</h1>


          <v-row v-if="role === 'Manager'">
            <v-col
              v-for="n in navManager"
              :key="n.route"
              class="d-flex flex-row-reverse"
              
            >
              <v-toolbar-title>
                <router-link :to="n.route"
                  ><v-btn class="pink lighten-2 white--text">
                    <v-icon left>
                      {{ n.icon }}
                    </v-icon>
                    {{ n.names }}</v-btn
                  ></router-link
                >
              </v-toolbar-title>
            </v-col>
            <v-col class="ml-6">
              <v-btn @click="navigate()" class="pink lighten-2 white--text" max-height="40">
                <v-icon left>mdi-logout</v-icon>
                Logout</v-btn
              >
            </v-col>
          </v-row>

          <v-row v-else-if="role === 'System Admin'">
            <v-col
              v-for="n in navSysadmin"
              :key="n.route"
              class="d-flex justify-center"
            >
              <v-toolbar-title>
                <router-link :to="n.route"
                  ><v-btn class="pink lighten-2 white--text">
                    <v-icon left>
                      {{ n.icon }}
                    </v-icon>
                    {{ n.names }}</v-btn
                  ></router-link
                >
              </v-toolbar-title>
            </v-col>
            <v-col class="ml-6">
              <v-btn @click="navigate()" class="pink lighten-2 white--text" max-height="40">
                <v-icon left>mdi-logout</v-icon>
                Logout</v-btn
              >
            </v-col>
          </v-row>

          <v-row v-else-if="role === 'Biller'">
            <v-col
              v-for="n in navBiller"
              :key="n.route"
              class="d-flex justify-center"
            >
              <v-toolbar-title>
                <router-link :to="n.route"
                  ><v-btn class="pink lighten-2 white--text">
                    <v-icon left>
                      {{ n.icon }}
                    </v-icon>
                    {{ n.names }}</v-btn
                  ></router-link
                >
              </v-toolbar-title>
            </v-col>
            <v-col class="ml-6">
              <v-btn @click="navigate()" class="pink lighten-2 white--text" max-height="40">
                <v-icon left>mdi-logout</v-icon>
                Logout</v-btn
              >
            </v-col>
          </v-row>

          <v-row v-else-if="role === 'Inventory'">
            <v-col
              v-for="n in navInventory"
              :key="n.route"
              class="d-flex justify-center"
            >
              <v-toolbar-title>
                <router-link :to="n.route"
                  ><v-btn class="pink lighten-2 white--text">
                    <v-icon left>
                      {{ n.icon }}
                    </v-icon>
                    {{ n.names }}</v-btn
                  ></router-link
                >
              </v-toolbar-title>
            </v-col>
            <v-col class="ml-6">
              <v-btn @click="navigate()" class="pink lighten-2 white--text" max-height="40">
                <v-icon left>mdi-logout</v-icon>
                Logout</v-btn
              >
             
              
            </v-col>
          </v-row>

    </v-app-bar>
    <v-main>
      <router-view></router-view>
    </v-main>
    <v-dialog v-model="dialog" persistent max-width="600">
      <v-card color="">
        <v-card-title class="headline">Confirm Logout</v-card-title>
        <v-card-text>Are you sure you want to proceed with logout?</v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="red darken-1" text @click="dialog = false">Cancel</v-btn>
          <v-btn color="primary" @click="logout">OK</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<script>
import crypto from 'crypto-js';
import eventservices from "./services/eventservices";
export default {
  name: "navigation",
  data() {
    return {
      navBiller: [
        { names: "Stockview", route: "/stockview", icon: "mdi-view-dashboard" },
        { names: "Bill Entry", route: "/billentry", icon: "mdi-cash-register" },
        { names: "Dashboard", route: "/dashboard", icon: "mdi-home" },
      ],
      navManager: [
        { names: "Stockview", route: "/stockview", icon: "mdi-view-dashboard" },
        { names: "Stock Entry", route: "/stockentry", icon: "mdi-plus-box" },
        {
          names: "Sales report",
          route: "/salesreport",
          icon: "mdi-chart-line",
        },
        { names: "Dashboard", route: "/dashboard", icon: "mdi-home" },
      ],
      navSysadmin: [
        { names: "Add user", route: "/adduser", icon: "mdi-account-plus" },
        { names: "Login History", route: "/loginhistory", icon: "mdi-history" },
        { names: "Dashboard", route: "/dashboard", icon: "mdi-home" },
      ],

      navInventory: [
        { names: "Stock Entry", route: "/stockentry", icon: "mdi-plus-box" },
        { names: "Stockview", route: "/stockview", icon: "mdi-view-dashboard" },
        { names: "Dashboard", route: "/dashboard", icon: "mdi-home" },
      ],
      navi: true,
      checkuser: "",
      shownav: false,
      role: "",
      user:null,
      logIn: "",
      userdetail: {
        login_history_id: null,
        user_id: "",
      },
      dialog:false,
    };
  },
 
  methods: {
    navigate() {
     this.dialog=true;
      
    },
    
    logout() {
      if (this.user) {
        this.userdetail.login_history_id=parseInt(this.user.login_history_id)
        eventservices.logoutapi(this.userdetail)
    .then((response) => {
      // console.log(response);
      if (response.data.status === "S") {
        this.role=''
        localStorage.removeItem('user')
        if (this.$route.path !== '/') {
           this.$router.push("/");
}
        this.shownav = false;
       
        this.dialog = false; 
      }
    })
    .catch((error) => {
      console.log(error);
    });
      }
 
},
   
  },
  created(){
      let encryptKey="venkat434";
      let user=localStorage.getItem('user');
      // console.log(user);
     if (user) {
      let decryptData=crypto.AES.decrypt(user,encryptKey).toString(crypto.enc.Utf8);
       this.user=JSON.parse(decryptData)
     }
    if (this.user) {
      this.role=this.user.role;
    }
  },
  beforeRouteLeave(to, from, next) {
 if (to.path === '/') {
    localStorage.removeItem('user');
    this.role=''
 }
 next();
},
  updated() {
    if (this.$route.path === '/') {
          this.role=''
          localStorage.removeItem('user')
     }
    else{
      const encryptKey="venkat434";
      let user=localStorage.getItem('user');
      if (user) {
      const decryptData=crypto.AES.decrypt(user,encryptKey).toString(crypto.enc.Utf8);
       this.user=JSON.parse(decryptData)
        if (this.user) {
          this.role=this.user.role
        }
      }
    }   
  },
  beforeDestroy(){  
    localStorage.removeItem('user')
  }
};
</script>


