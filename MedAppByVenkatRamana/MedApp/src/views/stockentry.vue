<template>
  <div>
    <h1 class="text-center mt-15">StockEntry</h1>
    <updatestockentry/>
  </div>
</template>
<script>
import updatestockentry from "../components/updatestockentry.vue";

export default {
  name: "stockentry",
  components: {
    updatestockentry
  },
};
</script>
