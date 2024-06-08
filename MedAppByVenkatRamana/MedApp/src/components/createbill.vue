<template>
  <div>
    <v-container class="mt-10 pink lighten-2 text-center">
      <h4>BILL ENTRY</h4>
      <v-row class="ma-5 text-center">
        <v-expansion-panels inset>
          <v-expansion-panel>
            <v-expansion-panel-header>Item</v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-row>
                <v-col>
                  <v-select
                    v-model="selectedMedicine"
                    :items="items"
                    item-text="medicine_name"
                    :rules="[(v) => !!v || 'Item is required']"
                    label="Medicine Name"
                    width="100%"
                    required
                  ></v-select>
                </v-col>
                <v-col>
                  <v-text-field
                    v-model="mqty"
                    type="number"
                    label="Quantity"
                    width="100%"
                  ></v-text-field>
                </v-col>
                <v-col max-width="300">
                  <v-btn @click="addtocart" width="100%" color="primary"
                    >Add</v-btn
                  >
                </v-col>
              </v-row>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-row>
    </v-container>
    <h2 id="msg" class="red--text text-center"></h2>
    <br />
    <br />
    <v-container class="pink lighten-2 mt-5">
      <v-data-table
        :items="filteredbill"
        :headers="headers"
        class="pa-10 text-center"
      >
        <template v-slot:top>
          <v-row justify="center">
            <v-col>
              <v-dialog v-model="dialog" persistent max-width="600px">
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                    color="primary"
                    width="100%"
                    dark
                    v-bind="attrs"
                    v-on="on"
                  >
                    preview
                  </v-btn>

                  <!-- Preview content Start-->
                </template>
                <v-card>
                  <v-card-title class="primary text-center">
                    <span class="text-center">PREVIEW BILL</span>
                    <v-spacer></v-spacer>
                    <v-btn icon  @click="dialog = false">
                      <v-icon >mdi-close</v-icon>
                    </v-btn>
                  </v-card-title>
                  <v-card-text>
                    <v-container>
                      <v-row>
                        <v-col> Medicine Name </v-col>
                        <v-col> Qty </v-col>
                        <v-col> Amount </v-col>
                      </v-row>
                      <v-row
                        m12
                        v-for="(item, index) in bill_details"
                        :key="index"
                      >
                        <v-col m3>{{ item.medicine_name }}</v-col>
                        <v-col m3>{{ item.quantity }}</v-col>
                        <v-col m3>{{ item.unit_price }}</v-col>
                      </v-row>
                      <v-row>
                        <v-col></v-col>
                        <v-col>Total </v-col>
                        <v-col
                          ><input v-model.number="t" label="Total" readonly />
                        </v-col>
                      </v-row>
                      <v-row>
                        <v-col></v-col>
                        <v-col> GST(18%) </v-col>
                        <v-col class="text-center"
                          ><input v-model.number="g" label="GST" readonly />
                        </v-col>
                      </v-row>
                      <v-row>
                        <v-col></v-col>
                        <v-col> NET PRICE </v-col>
                        <v-col
                          ><input v-model.number="np" label="Net Payable" />
                        </v-col>
                      </v-row>
                    </v-container>
                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" text @click="dialog = false">
                      <span @click="download()"> Print</span>
                    </v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>
            </v-col>
            <!-- Preview content End-->

            <v-col>
              <v-btn class="primary" width="100%" @click="save">Save</v-btn>
            </v-col>
            <v-col>
              <v-text-field
                v-model="billNo"
                label="Bill No"
                readonly
                disabled
              ></v-text-field>
            </v-col>
            <v-col>
              <v-text-field
                v-model="todayDate"
                label="Today Date"
                readonly
                disabled
              ></v-text-field>
            </v-col>
            <v-col>
              <v-text-field
                v-model.number="t"
                label="Total"
                readonly
                disabled
              ></v-text-field>
            </v-col>
            <v-col>
              <v-text-field
                v-model.number="g"
                label="GST(18%)"
                readonly
                disabled
              ></v-text-field>
            </v-col>
            <v-col>
              <v-text-field
                v-model.number="np"
                label="Net Payable"
                disabled
              ></v-text-field>
            </v-col>
          </v-row>
          <v-text-field
            v-model="search"
            label="Search"
            class="mx-4"
          ></v-text-field>
        </template>
        <template v-slot:item.actions="{ item }">
          <v-btn size="small" @click="deleteItem(item)">delete</v-btn>
        </template>
      </v-data-table>
    </v-container>
    <v-snackbar v-model="snackbar" timeout="3000">
      {{ text }}

      <template v-slot:action="{ attrs }">
        <v-btn color="pink" text v-bind="attrs" @click="snackbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script>
import crypto from "crypto-js";
import eventservices from "../services/eventservices";
import Papa from "papaparse";

export default {
  name: "createbill",
  data() {
    return {
      tempArr: [],
      snackbar: false,
      text: null,
      selectedMedicine: "",
      mqty: 0,
      billarr: [],
      items: [],
      Eqty: [],
      todayDate: new Date().toLocaleDateString("en-ca"),
      total: 1,
      gst: 1,
      netPayable: 0,
      price: "",
      dialog: false,
      t: 0,
      g: 0,
      np: 0,
      billNo: "",
      tempbd: [],
      bill_details: [],
      mediAndQuan: [],
      checkuser: "",
      count: 0,
      mname: "",
      user_id: "",
      search: "",
      headers: [
        {
          text: "MedicineName",
          align: "start",
          sortable: false,
          value: "medicine_name",
        },
        { text: "Brand", value: "brand" },
        { text: "Quantity", value: "quantity" },
        { text: "Amount", value: "amount" },
        { text: "Delete", value: "actions" },
      ],
    };
  },
  created() {
    this.generateBillNo();
    this.getMed();
  },
  methods: {
    getMed() {
      eventservices
        .stockview()
        .then((response) => {
          if (response.data.status == "S") {
            this.items = response.data.stockviewarr || [];
            this.tempArr = response.data.stockviewarr || [];
          }
        })
        .catch((error) => {
          console.log(error);
        });
    },

    generateBillNo() {
      const prefix = "BILL";
      const timestamp = new Date().getTime();
      this.billNo = `${prefix}-${timestamp}`;
    },
    preview() {
      this.dialog = true;
    },
    Upstock(body) {
      eventservices
        .updateStock(body)
        .then((response) => {
          if (response.data.status == "S") {
            // console.log("yes");
          }
        })
        .catch((error) => {
          console.log(error);
        });
    },
    BillMaster(body) {
      eventservices
        .savebillmaster(body)
        .then((response) => {
          if (response.data.status == "S") {
            // console.log("yes");
          }
        })
        .catch((error) => {
          console.log(error);
        });
    },
    deleteItem(item) {
      const index = this.bill_details.indexOf(item);

      for (let i = 0; i < this.items.length; i++) {
        if (
          this.items[i].medicine_name == this.bill_details[index].medicine_name
        ) {
          this.tempArr[i].quantity = +Number(this.bill_details[index].quantity);
        }
      }
      this.bill_details.splice(index, 1);
      this.t -= Number(item.quantity) * Number(item.unit_price);
      this.g = this.t * (18 / 100);
      this.np = this.t + this.g;
    },
    save() {
      let update = [];
      this.bill_details.forEach((element) => {
        update.push({
          medicine_name: element.medicine_name,
          qty: element.quantity,
        });
      });

      let obj = { bill_details: this.bill_details };
      eventservices
        .savebill(obj)
        .then((response) => {
          if (response.data.status == "S") {
            // console.log("yes");
          }
        })
        .catch((error) => {
          console.log(error);
        });
      let billmaster = {
        bill_no: this.billNo,
        billdate: this.todayDate,
        bill_amount: this.t,
        gst: Number(this.g.toFixed(2)),
        net_price: Number(this.np.toFixed(2)),
        user_id: this.user_id,
      };
      this.BillMaster(billmaster);
      let obj1 = { update_stock: update };
      this.Upstock(obj1);
      this.getMed();

      this.count = 0;
      this.tempbd = []; // Clearing bill details array

      this.generateBillNo();
      (this.t = 0),
        (this.g = 0),
        (this.billarr = []),
        (this.np = 0),
        (this.bill_details = []);
    },
    download() {
      const csv = Papa.unparse(this.bill_details);
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
    },
    addtocart() {
      if (this.selectedMedicine === "") {
        this.snackbar = true;
        this.text = "Please select the Medicine name ";
        return;
      }
      if (this.mqty === "" || this.mqty < 0 || this.mqty === 0) {
        this.snackbar = true;
        this.text = "Please Enter the Quantity greater than Zero";
        return;
      }
      if ((this.selectedMedicine === " ") & (this.mqty > 0)) {
        this.snackbar = true;
        this.text = "Please select the Medicine name ";
        return;
      }
      if (!Number.isInteger(Number(this.mqty))) {
        this.snackbar = true;
        this.text = "Decimal number not allowed";
        return;
      }
      if (
        (this.selectedMedicine === "") &
        (this.mqty == "" || this.mqty == 0)
      ) {
        this.snackbar = true;
        this.text = "Please Select Medicine name and  Quantity ";
        return;
      }
      if (this.selectedMedicine != "" && this.mqty > 0) {
        for (let i = 0; i < this.items.length; i++) {
          let old = false;
          let index = -1;
          if (
            this.selectedMedicine == this.items[i].medicine_name &&
            this.mqty > this.tempArr[i].quantity
          ) {
            this.snackbar = true;
            this.text =
              "Out of quantity, available quantity is " +
              this.items[i].quantity;
            return;
          }
          for (let j = 0; j < this.bill_details.length; j++) {
            if (this.bill_details[j].medicine_name == this.selectedMedicine) {
              old = true;
              index = j;
              break;
            }
          }
          if (old) {
            if (this.selectedMedicine == this.items[i].medicine_name) {
              this.bill_details[index].quantity =
                Number(this.bill_details[index].quantity) + Number(this.mqty);
              this.tempArr[i].quantity =
                Number(this.tempArr[i].quantity) - this.mqty;
              this.mname = this.items[i].medicine_name;
              this.price = this.items[i].unit_price;
              this.bill_details[index].amount =
                Number(this.bill_details[index].quantity) * Number(this.price);

              this.t += this.mqty * Number(this.price);
              this.g = (this.t * 18) / 100;
              this.np = this.t + this.g;
            }
          } else {
            if (this.selectedMedicine == this.items[i].medicine_name) {
              this.tempArr[i].quantity =
                Number(this.tempArr[i].quantity) - this.mqty;
              this.mname = this.items[i].medicine_name;
              this.price = this.items[i].unit_price;
              this.bill_details.push({
                bill_no: this.billNo,
                amount: Number(this.mqty) * Number(this.items[i].unit_price),
                medicine_name: this.items[i].medicine_name,
                brand: this.items[i].brand,
                user_id: this.user_id,
                quantity: Number(this.mqty),
                unit_price: this.items[i].unit_price,
              });
              this.t += this.mqty * Number(this.price);
              this.g = (this.t * 18) / 100;
              this.np = this.t + this.g;
            }
          }
        }

        this.selectedMedicine = " ";
        this.mqty = "";
      }
    },
  },

  watch: {
    selectedMedicine() {
      this.snackbar = false;
    },
    mqty() {
      this.snackbar = false;
    },
  },
  mounted() {
    let user = localStorage.getItem("user");
    if (user) {
      const encryptKey = "venkat434";
      const decryptData = crypto.AES.decrypt(user, encryptKey).toString(
        crypto.enc.Utf8
      );
      let userData = JSON.parse(decryptData);
      this.user_id = userData.user_id;
    }
  },
  computed: {
    filteredbill() {
      if (this.search) {
        return this.bill_details.filter((item) =>
          Object.values(item).some((value) =>
            value.toString().toLowerCase().includes(this.search.toLowerCase())
          )
        );
      }
      return this.bill_details;
    },
  },
};
</script>
