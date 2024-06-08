<template>
  <div>
    <v-container class="text-center  mt-15">
      <v-row class="pink lighten-2" justify="center">
        <v-col></v-col>
        <v-col></v-col>
        <v-col></v-col>
        <v-col>
          <v-dialog v-model="dialog" persistent max-width="600px">
            <template v-slot:activator="{ on, attrs }">
              <v-btn color="primary" dark v-bind="attrs" v-on="on"> ADD </v-btn>
            </template>
            <v-card>
              <v-card-title class="pink lighten-2">
                <v-btn icon  @click="dialog = false">
                  <v-icon>mdi-close</v-icon>
                </v-btn>
                <span class="  text-h5">ADD STOCK</span>
              </v-card-title>
              <v-card-title class="text-center">
                <span class="text-h5 text-center">MedApp</span>
                
              </v-card-title>
              <v-card-text>
                <v-container>
                  <v-row>
                    <v-col>
                      <v-text-field
                        v-model="mName"
                        label=" Medicine Name"
                      ></v-text-field>
                    </v-col>
                  </v-row>
                  <v-row>
                    <v-col>
                      <v-text-field
                        v-model="bName"
                        label="Brand Name "
                      ></v-text-field>
                    </v-col>
                  </v-row>
                </v-container>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <h2 id="msg" class="red--text"></h2>
                <v-btn color="blue darken-1  " text @click="addtostock">
                  ADD
                </v-btn>
              </v-card-actions>
             
            </v-card>
           
          </v-dialog>
        </v-col>
      </v-row>

      <v-row>
        <v-col>
          <v-select
            :items="brands"
            item-text="medicine_name"
            label="Medicine name"
            v-model="choosedMedicine"
            @change="selectMedicine"
          ></v-select>
        </v-col>
        <v-col>
          <v-text-field
            Read-only
            disabled
            :value="this.localbrand"
            label="Brand"
          >
          </v-text-field>
        </v-col>
        <v-col>
          <v-text-field
            v-model.number="qty"
            type="number"
            label="Qty"
          ></v-text-field>
        </v-col>
        <v-col>
          <v-text-field
            v-model.number="unitprice"
            label="Unit Price"
            type="number"
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row class="pink lighten-2">
        <v-col>
          <v-btn class="white--text blue darken-1" @click="updated">UPDATE</v-btn>
        </v-col>
       
      </v-row>
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
    </v-container>
    <h2 id="updateMsg" class="red--text text-center"></h2>
  </div>
</template>
<script>
import crypto from 'crypto-js';
import eventservices from '../services/eventservices';
export default {
  name: "updatestockentry",
  data() {
    return {
      brands:[],
      choosedMedicine: "",
      localbrand: "",
      snackbar:false,
      text:'',
      qty: 0,
      unitprice: "",
      mName: "",
      bName: "",
      user_id:'',
      dialog: false,
      checkadd:false,
    
    };
  },
  methods: {
    getMed(){
      eventservices.stockapi()
      .then((response)=>{
        if(response.data.status == 'S'){
           this.brands=response.data.stockarr ||[]
        }
      }).catch((error)=>{
        console.log(error);
      });
    },
    addtostock() {
     
     if(this.mName!="" && this.bName!="" && this.mName!=" " && this.bName!=" "){
      let obj={medicine_name:this.mName,brand:this.bName,user_id:''}
      eventservices.addmedicine(obj)
      .then((response)=>{
        if(response.data.status == 'S'){
            this.mName=""
            this.bName=""
            this.snackbar=true
            this.text="Added Successfully"
            this.dialog = false;
            this.getMed()
        }else{
          this.mName=""
            this.bName=""
            this.snackbar=true
            this.text="Medicine Already Exists"
        }
      }).catch((error)=>{
        console.log(error);
      });
     } else if (this.mName=="" && this.bName=="") {
      this.snackbar=true
      this.text=`Please fill the details`
        return
      }else if (this.mName=="") {
      this.snackbar=true
      this.text=`Please Enter the Medicine Name`
        return
      }else if (this.bName=="") {
      this.snackbar=true
      this.text=`Please Enter the Brand Name`
        return
      }else if (this.bName==" " || this.mName==" ") {
      this.snackbar=true
      this.text=`Don't Enter Space`
        return
      }
    },
    mounted(){
      let userData=localStorage.getItem('user')
      let user;
     if (userData) {
      const encryptKey="venkat434";
      const decryptData=crypto.AES.decrypt(userData,encryptKey).toString(crypto.enc.Utf8);
       user=JSON.parse(decryptData)
     }
       
        if (user) {
          this.user_id=user.user_id
        }
    },
    updated() {
      if(this.choosedMedicine=="" && this.unitprice=="" && this.qty==""){
        this.snackbar=true
      this.text=`Please fill the details`
        return
      }if(!Number.isInteger(this.qty) || this.qty<=0){
        this.snackbar=true
      this.text=`Enter Valid Quantity`
        return
      }if(!Number.isInteger(this.unitprice) || this.unitprice<=0){
        this.snackbar=true
      this.text=`Enter Valid unit price`
        return
      } if(this.choosedMedicine==""){
        this.snackbar=true
      this.text=`Please Select the Medicine Name`
        return
      }if( this.qty=="" || this.qty==0){
        this.snackbar=true
      this.text=`Please Enter the Quantity greater than zero `
        return
      } if( this.unitprice=="" || this.unitprice==0){
        this.snackbar=true
      this.text=`Please Enter the Unitprice greater than zero `
        return
      }
     
      if(this.choosedMedicine!="" & this.unitprice!="" & this.unitprice>0 & this.qty>0 & this.qty!=""){
        let updatemedicine = {
          medicine_name:this.choosedMedicine,
          brand:this.localbrand,
          unit_price:this.unitprice,
          quantity:this.qty,
          user_id:this.user_id
        }
        eventservices.updatestockapi(updatemedicine)
        .then((response)=>{
        if(response.data.status == 'S'){
          this.snackbar=true
          this.text="Updated Successfully"
          this.choosedMedicine=""
          this.unitprice=="" 
           this.unitprice=0 
            this.qty=0
        }
      }).catch((error)=>{
        console.log(error);
      });
      }
      this.getMed()
    },
    selectMedicine() {
      if (this.choosedMedicine != "") {
        for (let i = 0; i < this.brands.length; i++) {
          if (this.brands[i].medicine_name == this.choosedMedicine) {
            this.localbrand = this.brands[i].brand;
          }
        }
      }
    },
    openDialog() {
      this.dialog = true;
    },
  },
  watch:{
    mName(){
      if(this.mName!="" ){
        this.snackbar=false
      this.text=''
      }
    },
    bName(){
      if(this.bName!="" ){
        this.snackbar=false
      this.text=''
      }
    },
    choosedMedicine(){
      if(this.choosedMedicine!="" ){
        this.snackbar=false
      this.text=''
      }
    },
    qty(){
      if(this.qty!="" ){
        this.snackbar=false
      this.text=''
      }
    },
    unitprice(){
      if(this.unitprice!="" ){
        this.snackbar=false
      this.text=''
      }
    },
    dialog(){
      this.mName=""
      this.bName=""
    }
  },
  created(){
    this.getMed()
  }
};
</script>
