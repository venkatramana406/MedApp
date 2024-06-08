import axios from "axios";

const baseApiclient=axios.create({
    baseURL:'http://localhost:22991',
    withCredentials:false,
    headers:{
        Accept:'application/json',
        'Content-Type':'application/json'
    }
})
export default{
    loginapi(body){
        return baseApiclient.post('/login',body)
    },
    loginhistoryapi(){
        return baseApiclient.get('/loginhistory')
    },
    dashboardroleapi(){ 
        return baseApiclient.post('/dashboardrole')
    },
    dashboardapi(body){ 
        return baseApiclient.post('/dashboard',body)
    },
    logoutapi(body){
        // console.log("body = > ",body);
        return baseApiclient.put('/logout',body)
    },
    stockview(){
        return baseApiclient.get('stockview')
    },
    adduser(body){
        return baseApiclient.put('/adduser',body)
    },
    addmedicine(body){
        return baseApiclient.put('/addmedicine',body)
    },
    stockapi(){
        return baseApiclient.get('/stock')
    },
    updatestockapi(body){
        return baseApiclient.put('/updatestock',body)
    },
    salesreport(body){
        return baseApiclient.put('/salesreport',body)
    },
    savebill(body){
        return baseApiclient.post('/savebill',body)
    },
    savebillmaster(body){
        return baseApiclient.post('/savebillmaster',body)
    },
    updateStock(body){
        return baseApiclient.post('/upstock',body)
    },
    weeklysales(){
        return baseApiclient.post('/lastweeksales')
    },
    monthlysales(){
        return baseApiclient.post('/lastmonthsales')
    },
    salesmanperformance(){
        return baseApiclient.post('/salesmanperformance')
    }

}

