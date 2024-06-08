package savebill

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"net/http"
)

type UpdateStockReqStruct struct {
	UpdateStockArr []UpdateStockStruct `json:"update_stock"`
}
type CommonRespStruct struct {
	ErrMsg string `json:"errmsg"`
	Status string `json:"status"`
	Msg    bool   `json:"msg"`
}

// store updated stock MedicineName and Quantity structure
type UpdateStockStruct struct {
	Quantity     int    `json:"qty"`
	MedicineName string `json:"medicine_name"`
}

// update Stock using this api
func UpdateStockBillAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-headers", "USER,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization")
	log.Println("update stock(+)")
	if r.Method == "POST" {
		fmt.Println("UpdateStockBillAPI(+)")
		var lUpStock UpdateStockReqStruct

		var resp CommonRespStruct
		body, err := io.ReadAll(r.Body)
		resp.Status = "S"
		if err != nil {
			log.Println(err)
			resp.Status = "E"
			resp.ErrMsg = "Error" + err.Error()
		} else {
			err := json.Unmarshal(body, &lUpStock)
			if err != nil {
				log.Println(err)
				resp.Status = "E"
				resp.ErrMsg = "Error" + err.Error()
			} else {
				resp = BillUpdateStock(lUpStock.UpdateStockArr)
			}

		}
		data, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(w, "Error taking data"+err.Error())
		} else {
			w.Write(data)
		}
		fmt.Println("UpdateStockBillAPI(-)")
	}
}

func BillUpdateStock(upStock []UpdateStockStruct) CommonRespStruct {
	var resp CommonRespStruct
	log.Println(upStock)
	db, err := LocalDBConnect()
	resp.Status = "S"
	resp.Msg = true
	if err != nil {
		log.Println("BUS01", err)
		resp.ErrMsg = err.Error()
		resp.Status = "E"
		resp.Msg = false
		return resp
	}
	defer db.Close()
	updateStock := `UPDATE MEDAPP_STOCK STOCK
	JOIN MEDAPP_MEDICINE_MASTER  ON STOCK.MEDICINE_MASTER_ID = MEDAPP_MEDICINE_MASTER.MEDICINE_MASTER_ID
	SET STOCK.QUANTITY  = NVL(STOCK.QUANTITY,0)  - ?
	WHERE MEDAPP_MEDICINE_MASTER.MEDICINE_NAME  = ?;`
	for _, update := range upStock {
		result, err := db.Exec(updateStock, update.Quantity, update.MedicineName)
		if err != nil {
			log.Println("BUS02", err)
			resp.ErrMsg = err.Error()
			resp.Status = "E"
			resp.Msg = false
			return resp
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Println("BUS03", err)
			resp.ErrMsg = err.Error()
			resp.Status = "E"
			resp.Msg = false
			return resp
		}
		if rowsAffected == 0 {
			log.Println("no row", update.MedicineName)
		}
	}
	return resp
}
