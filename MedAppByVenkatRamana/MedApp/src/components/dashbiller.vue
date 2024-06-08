<template>
  <div>
    <v-container>
      <v-card justify="center"
      height="50"
     
     
      class="px-2 mx-auto ">
        <h1 class="text-center text-h3 font-weight-regular black--text">Welcome Biller</h1>
      </v-card>
    </v-container>
    <br>
    <v-card  justify="center"
    height="50"
    width="400"
    color="pink lighten-2"
    class="pt-2 mx-auto white--text">
      <h2 class="text-center">YOUR TODAY SALES</h2>
    </v-card>
    <v-container  color="black">
      <v-card
        justify="center"
        height="150"
        width="400"
        color="pink lighten-2"
        class="pa-10 mx-auto white--text"
      >
    <span v-if="profit" class="text-h2 text-center"
      >{{"Rs. "+  totalSales }}
      <v-icon x-large color="green">mdi-arrow-up</v-icon></span
    >

    <span v-else class="text-h2 text-center"
      >{{ "Rs. "+ totalSales }}
      <v-icon x-large color="red">mdi-arrow-down</v-icon></span
    >
  </v-card>
</v-container>
  </div>
</template>

<script>
import eventservices from "../services/eventservices";
import crypto from 'crypto-js';
export default {
  name: "dashbiller",
  data() {
    return {
      todaySales: 0,
      yesterdaySales: 0,
      profit: false,
      totalSales: 0,
      user_id:""
    };
  },
  methods: {
    getdetail() {
      let user=localStorage.getItem('user');
      if (user) {
        const encryptKey="venkat434";
      const decryptData=crypto.AES.decrypt(user,encryptKey).toString(crypto.enc.Utf8);
        user=JSON.parse(decryptData)
        // console.log(user," <- ");
      }
        if (user) {
         this.user_id=user.user_id
         eventservices.dashboardapi({user_id:this.user_id})
        .then((response) => {
          this.todaySales = response.data.biller_detail[0].today_sales;
          this.yesterdaySales = response.data.biller_detail[0].yesterday_sales;
          if (this.todaySales > this.yesterdaySales) {
        this.totalSales = this.todaySales.toFixed(2);
        // console.log("totalsales - ",this.totalSales);
        this.profit = true;
      } else {
        this.totalSales = this.todaySales.toFixed(2);
        this.profit = false;
      }
        })
        .catch((error) => {
          console.log(error);
        });
        }
     
    },
  },
  mounted () {
    this.getdetail();
   
  },
};
</script>
