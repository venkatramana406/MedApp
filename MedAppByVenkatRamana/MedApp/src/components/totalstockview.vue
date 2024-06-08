<template>
<div>
      <h1 class="text-center mt-5">STOCK VIEW</h1>
    <v-container class="pink lighten-2 mt-5">
      <v-data-table :items="filteredstocks" :headers="headers" class="pa-10 text-center">
        <template v-slot:top>
          <v-text-field v-model="search" label="Search" class="mx-4"></v-text-field>
        </template>
      </v-data-table>
    </v-container>
</div>
</template>
<script>
import eventservices from '../services/eventservices';
export default {
  name: "totalstockview",
  data() {
    return {
      search:'',
      stocks:[],
      brand:"", 
      headers: [
        {
          text: 'MedicineName',
          align: 'start',
          sortable: false,
          value: 'medicine_name',
        },
        { text: 'Brand', value: 'brand' },
        { text: 'Quantity', value: 'quantity' },
        { text: 'Unit Price', value: 'unit_price' },
      ],
    };
  },
  computed: {
    filteredstocks() {
      if (this.search) {
        return this.stocks.filter((item) =>
          Object.values(item).some((value) =>
            value.toString().toLowerCase().includes(this.search.toLowerCase())
          )
        );
      }
      return this.stocks;
    },
  },
  methods:{
    getMed(){
      eventservices.stockview()
    .then((response)=>{
     if(response.data.status=="S"){
      this.stocks=response.data.stockviewarr||[]
     }
    }).catch((error)=>{
      console.log(error)
    })
    }
  },
  mounted(){
   this.getMed()
  },

  
 
};
</script>