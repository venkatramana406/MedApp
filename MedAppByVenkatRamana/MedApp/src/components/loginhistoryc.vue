
 <template>
  <div>
    <h1 class="text-center">LOGIN HISTORY</h1>
    <v-container class="pink lighten-2 mt-5">
      
      <v-data-table :items="filteredHistory" :headers="headers" class="pa-10 text-center">
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
  name: "LoginHistoryComponent",
  data() {
    return {
      search: "",
      loginhistoryarr: [],
      headers: [
        {
          text: 'USER_ID',
          align: 'start',
          sortable: false,
          value: 'user_id',
        },
        { text: 'Login Date', value: 'login_date' },
        { text: 'Login Time', value: 'login_time' },
        { text: 'Logout Date', value: 'logout_date' },
        { text: 'Logout Time', value: 'logout_time' },
      ],
    };
  },
  computed: {
    filteredHistory() {
      if (this.search) {
        return this.loginhistoryarr.filter((item) =>
          Object.values(item).some((value) =>
            value.toString().toLowerCase().includes(this.search.toLowerCase())
          )
        );
      }
      return this.loginhistoryarr;
    },
  },
  mounted() {
    eventservices.loginhistoryapi()
    .then((response) => {
      if (response.data.status == "S") {
        this.loginhistoryarr = response.data.loginhistoryarr||[];
      }
    })
    .catch((error) => {
      console.log(error);
    });
  },
  updated() {
    eventservices.loginhistoryapi()
    .then((response) => {
      if (response.data.status == "S") {
        this.loginhistoryarr = response.data.loginhistoryarr||[];
      }
    })
    .catch((error) => {
      console.log(error);
    });
  }
};
</script>

