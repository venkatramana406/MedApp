<template>
  <div id="app" >
    <div class="mt-15 text-center">
      <h1>LOGIN PAGE</h1>
      <br />
      <v-container>
        <v-card
          justify="center"
          height="400"
          width="400"
          color="whitesmoke lighten-2"
          class="pa-10 mx-auto"
        >
          <v-form>
            <v-text-field
              v-model="validation.user_id"
              label="Username"
              :counter="20"
              required
              class="mt-10"
            ></v-text-field>
            <v-text-field
              v-model="validation.password"
              label="Password"
              :counter="5"
              type="password"
              required
              class="mt-5"
            ></v-text-field>
            <div v-if="validation.user_id != '' && validation.password != ''">
              <v-btn v-show="true" class="primary" @click="logIn">LOG IN</v-btn>
            </div>
            <br />
            <div v-show="showerr">
              <v-alert shaped dense prominent type="error"
                >Please verify the username and password</v-alert
              >
            </div>
          </v-form>
        </v-card>
      </v-container>
    </div>
  </div>
</template>

<script>
import crypto from  'crypto-js';
import eventservices from "../services/eventservices";
export default {
  name: "loginpage",
  data() {
    return {
      button: false,
      validation: {
        user_id: "",
        password: "",
      },
      role: "",
      showerr: false,
    };
  },

  watch: {
    validation: {
      handler(value) {
        if (value.user_id != "") {
          this.showerr = false;
        }
      },
      deep: true,
    },
  },
  methods: {
    logIn() {
      if (this.validation.user_id != "" && this.validation.password != "") {
        eventservices.loginapi(this.validation)
          .then((response) => {
            // console.log("login=>",response);
            if (response.data.status == "S") {
           
            let role = response.data.userdetails.role;
            let login_id = response.data.userdetails.login_id;
            let user_id = this.validation.user_id;
            let login_history_id = response.data.userdetails.login_history_id;
            let user={user_id:user_id,role:role,login_id:login_id,login_history_id:login_history_id}
            let userData=JSON.stringify(user);
            const encryptKey="venkat434";
            const encryptData=crypto.AES.encrypt(userData,encryptKey).toString();
            localStorage.setItem('user',encryptData);
            this.$store.commit('UPDATE_USER',user)
              this.$router.push("/dashboard");
            } else {
              this.showerr = true;
              this.validation.user_id = "";
              this.validation.password = "";
              document.getElementById("error").innerHTML =
                "Please verify the username and password ";
            }
          })
          .catch((error) => {
            console.log(error);
          });
      }
    },
  },
};
</script>
