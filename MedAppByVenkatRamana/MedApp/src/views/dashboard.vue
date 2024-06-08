<template>
<!-- <v-app> -->
 <div>
  
  <div v-if="this.role=='Manager'" >
        <dashmanager/>
      </div>
      <div v-else-if="this.role=='System Admin'" >
        <dashsysadmin/>
      </div>
      <div v-else-if="this.role=='Biller'">
        <dashbiller/>
      </div>
      <div v-else-if="this.role=='Inventory'" >
        <dashinventory/>
      </div>
</div>
</template>

<script>

import crypto from 'crypto-js';
export default {
  name: "Dash",
  data(){
    return{
    role:'',
    }
  },
  mounted() {
    let userData=localStorage.getItem('user')
     if (userData) {
      const encryptKey="venkat434";
      const decryptData=crypto.AES.decrypt(userData,encryptKey).toString(crypto.enc.Utf8);
       let user=JSON.parse(decryptData)
        if (user) {
          this.role=user.role
          // console.log(this.role);
        }
     }
  }

}
  
 
</script>