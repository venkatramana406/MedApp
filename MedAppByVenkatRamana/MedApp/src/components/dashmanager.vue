<template>
  <div class="text-center">
    <v-container>
    <v-card justify="center"
    height="50" 
    class="px-2 mx-auto ">
      <h1 class="text-center text-h4 font-weight-regular black--text">Welcome Manager</h1>
    </v-card>
  </v-container>
    <v-container >
      <v-card  justify="center"
      height="100"
      width="400"
      color="pink lighten-2"
      class="pt-2 mx-auto white--text">
      <h2 class="text-center">TODAY SALES</h2>
      <h2 class="text-center">{{ "Rs. "+ this.todaysales }}</h2>
      </v-card>
    </v-container>
      <v-container >
        <v-card  justify="center"
        height="100"
        width="400"
        color="pink lighten-2"
        class="pt-2 mx-auto white--text">
        <h2 class="text-center">INVENTORY VALUE</h2>
          <h2 class="text-center">{{ "Rs. "+ this.CurrentInvenValue }}</h2>
        </v-card>
        <br>
        <h1>STOCK VIEW</h1>
        <v-sheet v-show="display">
          <apexchart width="90%" type="bar" :options="options" :series="series"></apexchart>
        </v-sheet>
        <h1>DAILY SALES TREND</h1>
        <v-sheet v-show="display">
          <apexchart width="100%" type="line" :options="options1" :series="series1"></apexchart>
        </v-sheet>
        <h1>MONTHLY SALES TREND</h1>
        <v-sheet v-show="display">
          <apexchart width="100%" type="line" :options="options2" :series="series2"></apexchart>
        </v-sheet>
        <h1>SALESMAN PERFORMANCE</h1>
        <div id="chart">
          <apexchart  type="donut"  width="80%"  :options="chartOptions" :series="series3"></apexchart>
        </div>
    </v-container>
  </div>
</template>
<script>
import crypto from 'crypto-js';
import eventservices from "../services/eventservices";
import apexchart from 'vue-apexcharts';
export default {
  name: "dashmanager",
  components:{
    apexchart
  },
  data() {
    return {
      tDate: "",
      todaysales: 0,
      CurrentInvenValue: 0,
      display:false,
      options: {
        chart: {
          id: 'vuechart-example'
        },
        xaxis: {
          categories: []
        }
      },
      series: [{
        name: 'series',
        data:[]
      }],

      options1: {
        chart: {
          id: 'vuechart-example'
        },
        xaxis: {
          weeklysales: []
        }
      },
      series1: [{
        name: 'series-1',
        data:[]
      }],

      options2: {
        chart: {
          id: 'vuechart-example'
        },
        xaxis: {
          monthlysales: []
        }
      },
      series2: [{
        name: 'series-2',
        data:[]
      }],

      user:'',

      series3: [],
          chartOptions: {
            chart: {
              width: 380,
              type: 'donut',
            },
            labels: [],
            responsive: [{
              breakpoint: 480,
              options: {
                chart: {
                  width: 200
                },
                legend: {
                  position: 'bottom'
                }
              }
            }]
          },
          
    };
  },
  created() {
    this.getdetail();
    this.getData();
    this.getWeekSale();
    this.getMonthlySale();
    this.getSalesManPerformance();
  },
  methods: {
    getData(){
      eventservices
        .stockview()
        .then((response) => {
          if (response.data.status == "S") {
            // console.log(response);
            let stock = response.data.stockviewarr || [];
            for (let i = 0; i < stock.length; i++) {
              this.options.xaxis.categories.push(stock[i].medicine_name);
              this.series[0].data.push(stock[i].quantity)
            }
            this.display=true
          }
        })
        .catch((error) => {
          console.log(error);
        });
    },
    getdetail() {
      const encryptKey="venkat434";
      let user=localStorage.getItem('user');
      if (user) {
      const decryptData=crypto.AES.decrypt(user,encryptKey).toString(crypto.enc.Utf8);
       this.user=JSON.parse(decryptData)
        if (this.user) {
          this.user=this.user.user_id
        }
      }
      let obj={user_id:user}
      eventservices.dashboardapi(obj)
        .then((response) => {
            this.todaysales = response.data.manager_detail[0].today_sales;
            this.todaysales=Number(this.todaysales).toFixed(2)
            this.CurrentInvenValue = response.data.manager_detail[0].inventory_value;
            this.CurrentInvenValue=Number(this.CurrentInvenValue).toFixed(2)
            // console.log(this.todaysales," ",this.CurrentInvenValue);
        })
        .catch((error) => {
          console.log(error);
        });
    },
    getWeekSale(){
      eventservices.weeklysales()
        .then((response) => {
          if (response.data.status == "S") {
            // console.log(response);
            let weeklysale = response.data.last_week || [];
            for (let i = 0; i < weeklysale.length; i++) {
              this.options1.xaxis.weeklysales.push(weeklysale[i].bill_date);
              this.series1[0].data.push(weeklysale[i].daily_sales)
            }
            // console.log(this.series1[0].data);
            // console.log( this.options1.xaxis.weeklysales);
            this.display=true
          }
        })
        .catch((error) => {
          console.log(error);
        });
    },
    getMonthlySale(){
      eventservices.monthlysales()
        .then((response) => {
          if (response.data.status == "S") {
            // console.log(response);
            let monthlysale = response.data.last_month || [];
            for (let i = 0; i < monthlysale.length; i++) {
              this.options2.xaxis.monthlysales.push(monthlysale[i].bill_date);
              this.series2[0].data.push(monthlysale[i].daily_sales)
            }
            // console.log(this.series1[0].data);
            // console.log( this.options1.xaxis.weeklysales);
            this.display=true
          }
        })
        .catch((error) => {
          console.log(error);
        });
    },
    getSalesManPerformance(){
      eventservices.salesmanperformance()
      .then((response)=>{
        if(response.data.status=="S"){
          let sales=response.data.result_data || [];
          this.series3 = [];
          // console.log("sales => ",sales);
          for( let i=0;i<sales.length;i++){
            this.series3.push(Number( sales[i].bill_amount))
            this.chartOptions.labels.push(sales[i].user_id)
          }
          // console.log("series 3 => ",this.series3);
        }
      })
      .catch((error)=>{
        console.log(error);
      });
    }
  },
};
</script>
