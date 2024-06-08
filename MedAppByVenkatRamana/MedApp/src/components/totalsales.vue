<template>
  <div class="mt-15">
    <v-container class="pink lighten-2">
      <v-row>
        <v-col>SALES REPORT</v-col>
        <v-col></v-col>
        <v-col></v-col>
      </v-row>
      <v-row class="text-center">
        <v-col>
          <v-text-field v-model="startDate" :max="maxDate" type="date">
          </v-text-field>
        </v-col>
        <v-col>
          <v-text-field v-model="endDate" :max="maxDate" type="date">
          </v-text-field>
        </v-col>
        <v-col>
          <v-btn class="mt-4 primary" @click="loadData">SEARCH</v-btn>
        </v-col>
      </v-row>
    </v-container>
    <h2 id="updateMsg" class="red--text text-center"></h2>
    <v-container class="pink lighten-2   mt-5">
      <v-data-table
        :items="filteredsales"
        :headers="headers"
        :search="search"
        class="pa-10 text-center"
      >
        <template v-slot:top>
          <v-row>
            <v-col></v-col>
            <v-col></v-col>
            <v-col></v-col>
            <v-col></v-col>
            <v-col>
              <v-btn class="mt-4 primary" @click="downloadBill" >DOWNLOAD</v-btn>
            </v-col>
          </v-row>
          <v-text-field
            v-model="search"
            label="Search"
            class="mx-4"
          ></v-text-field>
        </template>
      </v-data-table>
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
import eventservices from "../services/eventservices";
import Papa from "papaparse";
export default {
  name: "totalsales",
  data() {
    return {
      snackbar:false,
      text:null,  
      detailArr: [],
      search: "",
      startDate: "",
      endDate: "",
      filteredsales: [],
      headers: [
        {
          text: "Bill no",
          align: "start",
          sortable: false,
          value: "bill_no",
        },
        { text: "Bill Date", value: "bill_date" },
        { text: "Medicine Name", value: "medicine_name" },
        { text: "Qty", value: "quantity" },
        { text: "Amount", value: "amount" },
      ],
    };
  },
  computed: {
    maxDate() {
      const today = new Date().toISOString().split("T")[0];
      return today;
    },

  },
  methods: {
    loadData() {
      if (
        this.startDate != "" &&
        this.endDate != "" &&
        this.startDate <= this.endDate
      ) {
        let fdate = new Date(this.startDate).toLocaleDateString("en-ca");
        let tdate = new Date(this.endDate).toLocaleDateString("en-ca");
        let givendates = {
          from_date: fdate,
          to_date: tdate,
        };
        eventservices
          .salesreport(givendates)
          .then((response) => {
            if (response.data.status == "S") {
              this.filteredsales = response.data.salesarr ||[];
            }
          })
          .catch((error) => {
            console.log(error);
          });
      }else if( this.startDate == "" && this.endDate=="") {
        
        this.snackbar=true
      this.text=`Please select the From date and To date`
        return
      } else if( this.startDate == "" ) {
        
        this.snackbar=true
      this.text=`Please select the (From date)`
        return
      }else if( this.endDate == "" ) {
        
        this.snackbar=true
      this.text=`Please select the (To date)`
        return
      }else if(this.startDate>this.endDate  ) {
        
        this.snackbar=true
      this.text=`Invalid From date`
        return
      }
    },
    clearSearch() {
      this.search = "";
    },
    downloadBill(){
      const csv = Papa.unparse(this.filteredsales);
      // console.log("data", csv, typeof (csv), [csv]);
      const blob = new Blob([csv], { type: "text/csv;charset=utf-8;" });
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.setAttribute("href", url);
      link.setAttribute("download", "data.csv");
      link.style.visibility = "visible";
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  },
  watch: {
    startDate() {
      this.filteredsales=[]
    },
    endDate() {
      this.filteredsales=[]
    },

  },
};
</script>
