package main

import (
	"log"
	"medapp/addmedicine"
	"medapp/adduser"
	"medapp/apexchart"
	"medapp/dashboard"
	"medapp/login"
	"medapp/loginhistory"
	"medapp/logout"
	"medapp/salesreport"
	"medapp/savebill"
	"medapp/stock"
	"medapp/stockview"
	"medapp/updatestock"
	"net/http"
)

func main() {
	log.Println("Server started")
	http.HandleFunc("/login", login.LoginAPI)

	http.HandleFunc("/adduser", adduser.AddUserAPI)

	http.HandleFunc("/addmedicine", addmedicine.AddMedicineAPI)
	http.HandleFunc("/updatestock", updatestock.UpdateStockAPI)
	http.HandleFunc("/salesreport", salesreport.SalesReportAPI)
	http.HandleFunc("/loginhistory", loginhistory.LoginHistoryAPI)
	http.HandleFunc("/logout", logout.LogoutAPI)
	http.HandleFunc("/stockview", stockview.StockViewAPI)
	http.HandleFunc("/stock", stock.StockAPI)
	http.HandleFunc("/dashboard", dashboard.DashboardAPI)
	http.HandleFunc("/dashboardrole", dashboard.DashboardRoleAPI)
	http.HandleFunc("/savebill", savebill.SaveBillAPI)
	http.HandleFunc("/savebillmaster", savebill.SaveBillMasterAPI)
	http.HandleFunc("/upstock", savebill.UpdateStockBillAPI)
	http.HandleFunc("/lastweeksales", apexchart.LastWeekSalesAPI)
	http.HandleFunc("/lastmonthsales", apexchart.LastMonthSalesAPI)
	http.HandleFunc("/salesmanperformance", apexchart.SalesManAPI)
	http.ListenAndServe(":22991", nil)
	log.Println("server Ended")
}
