<template>
  <div>
    <h4 class="text-center mt-10">USER DETAIL</h4>
    <v-container class="mt-10 pink lighten-2 text-center">
      <v-row class="ma-5 pa-4" justify="center">
        <v-expansion-panels inset>
          <v-expansion-panel>
            <v-expansion-panel-header>User Detail</v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-row>
                <v-col>
                  <v-text-field v-model="adduser.user_id" label="User Id"></v-text-field>
                </v-col>
                <v-col>
                  <v-text-field
                    v-model="adduser.password"
                    label="Password"
                  ></v-text-field>
                </v-col>
                <v-col>
                  <v-select
                    v-model="adduser.role"
                    :items="items"
                    label="Role"
                    required
                  ></v-select>
                </v-col>
                <v-col>
                  <v-btn @click="addtousers" color="blue lighten-1">Add</v-btn>
                </v-col>
              </v-row>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-row>
    </v-container>
    <v-snackbar
    v-model="snackbar"
    timeout="3000"
  >
    {{ text }}

    <template v-slot:action="{ attrs }">
      <v-btn
        color="pink"
        text
        v-bind="attrs"
        @click="snackbar = false"
      >
        Close
      </v-btn>
    </template>
  </v-snackbar>
  
  </div>
</template>

<script>
import eventservices from '../services/eventservices';

export default {
  name: "addnewuser",
  data() {
    return {
      snackbar:'',
      text:'',
      adduser:{
        user_id: "",
      password: "",
      role: "",
      },
      items: ["Biller", "Manager", "Systes Admin", "Inventory"],
      push:false
    };
  },
  methods: {
    addtousers() {
      if(this.adduser.user_id!="" && this.adduser.password!="" && this.adduser.role!=""){
        eventservices.adduser(this.adduser)
        .then((response)=>{
          if(response.data.status=="S"){

            this.snackbar=true
        this.text="User Added Successfully"
          this.adduser.user_id = "",
          this.adduser.password = "",
          this.adduser.role = "";
          }else{
            this.snackbar=true
        this.text="User Already Exists"
        this.adduser.user_id = "",
          this.adduser.password = "",
          this.adduser.role = "";
        return
          }
        }).catch((error)=>{
          console.log(error)
        }) 
        
        }else if( this.adduser.user_id=="" & (this.adduser.password=="" || this.adduser.password==0)){
        this.snackbar=true
        this.text="Set User Name and  Password "
        return
      }
      else if(this.adduser.user_id==""){
        this.snackbar=true
        this.text="Set Enter the User name "
        return
      }else if(this.adduser.password==""){
        this.snackbar=true
        this.text="Set the password"
        return
      } else if(this.adduser.role==""){
        this.snackbar=true
        this.text=" Please Select the role  "
        return
      }  
    },
  },
  watch:{
    adduser:{
      handler(value){
      if(value.user_id!="" || value.password!="" ||  value.role!=""){
        this.snackbar=false
        this.text=""
      }},
      deep:true
    },
  }
};
</script>
