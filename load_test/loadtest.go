package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/joshsoftware/mqtt-sdk-go/gosdk"
)

var wg sync.WaitGroup

func getData() map[string]interface{} {
	return map[string]interface{}{
		"pulsel":       rand.Intn(60),
		"pulseh":       rand.Intn(30),
		"meter_status": 0,
		"valve_status": 2,
	}
}

func sendEvent(gw gosdk.Gateway, t *gosdk.Thing) {
	// Send Event data
	data := getData()

	eventData := gosdk.CreateThingEvent(t, data, 0)
	retVal := gw.ThingEvent(eventData)
	if retVal == nil {
		fmt.Println("Send thing event: ", eventData)
	} else {
		fmt.Println("Unable to send thing event: ", retVal)
	}

}

func sendAlert(gw gosdk.Gateway, t *gosdk.Thing, alertLevel int, alertLevelStr string) error {
	data := map[string]interface{}{"tempearture": 120}
	if err := gw.Alert(t, "This is an example "+alertLevelStr+" alert using go SDK.", alertLevel, data); err != nil {
		fmt.Println("Unable to send Alert: ", err)
	}
	return nil
}

func sendAlerts(gw gosdk.Gateway, t *gosdk.Thing) error {
	sendAlert(gw, t, 0, "INFO")
	sendAlert(gw, t, 1, "WARNING")
	sendAlert(gw, t, 2, "ERROR")
	sendAlert(gw, t, 3, "CRITICAL")
	return nil
}

func handleInstruciton(gw gosdk.Gateway, ts int64, t *gosdk.Thing, alertKey string, instruction map[string]interface{}) {
	fmt.Println("Recieved Instruction Data for thing with key-device_key ", t.Key, ": ", instruction)
	alertMsg := fmt.Sprintf("Instruction with alert key %s executed", alertKey)
	if err := gw.InstructionAck(t, alertKey, alertMsg, 0, map[string]interface{}{"execution_status": "success"}); err == nil {
		fmt.Println("Sent instruction execution ACK back to datonis")
	} else {
		fmt.Println("Could not send instruction execution ACK back to datonis")
	}
}

func main() {

	file, err := os.Open("thing_key.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	rand.Seed(time.Now().UnixNano())
	fmt.Println("Here we go..")
	config := &gosdk.GatewayConfig{
		AccessKey: "78d9fa4teebt5add59ctb86e1a286477cb147392",
		Username:  "user",
		Password:  "user1",
	}

	config = gosdk.InitGatewayConfig(config)
	gw := gosdk.CreateGateway(config)

	for _, val := range records[0] {
		wg.Add(1)
		go launch_thing_test(gw, val)
	}

	wg.Wait()
}

//smly_xSnUil78K7,smly_ngCoKivaJM,smly_o1fSVMb4Pc,smly_ivkofURK58,smly_8DVFdRLCaQ,smly_Av4Ig6dQyA,smly_o6U0WsBMMk,smly_HxO7jH19Tb,smly_zTsjWtHa7f,smly_LStNq38fkx,smly_h3ePJM2cGj,smly_6SbheaDx0q,smly_jlVV4cdMhI,smly_envHsfNnnD,smly_soq7faYytj,smly_fg8MfWPImI,smly_rN9RNQaru8,smly_kkunGuxEn6,smly_MPWSTeYr0T,smly_YCBtnh1fjN,smly_EMwvivZoEg,smly_LQU92cVMrq,smly_WkMA9lLbPQ,smly_ZaEEeDvyAy,smly_wVhL4RQSqs,smly_9onJoJYtCd,smly_8jeez0l78d,smly_RCWctacVTJ,smly_8GIkKIM76K,smly_2a62ebEjXT,smly_6OLJ0ZGjmN,smly_p9gQvWAwyA,smly_PLilAowvT4,smly_4z4FEPaijR,smly_4DXYMzKymh,smly_6dgM9Qggdt,smly_IiBbisk2TA,smly_6d1eh8WlqF,smly_ywczQcQ0hU,smly_9fPdH8dTmi,smly_CwodCIogR1,smly_oBIsaPucSI,smly_TOflAw8nYZ,smly_AJPLeivPjb,smly_EvRAgunllY,smly_M0D8JXqaIo,smly_4mlWQPCj0l,smly_CN3VPjlGlr,smly_mZCEPCQJl9,smly_OAFu77cfV4,smly_Hm93mRfLW0,smly_thbjyba9RQ,smly_emqT1NU3JA,smly_P0iPTMexm2,smly_f1zkySCdzc,smly_fTrW7RmT9B,smly_Q82rxHNExm,smly_KIs6GRZTPS,smly_DdW04Jd2oR,smly_ptpZrCXQ73,smly_64AIl3CwlR,smly_OhDYHufe7Y,smly_rNSLzHpuLw,smly_YBTjfdX9sw,smly_OfOloeEF5R,smly_etWALDOKLC,smly_bY9SaHpyiY,smly_aV27SS0WCp,smly_q6xY4aUiLg,smly_5k6Z7LPQ3U,smly_TXltCz3Dip,smly_JETGX5JAsv,smly_tfbOh2OQdM,smly_2RgrN0wJIE,smly_O6RdFPLhtq,smly_nkx74hriKE,smly_awXtqISwG9,smly_D63F5vcoo0,smly_lVK041ZqAE,smly_d9PMGf97qr,smly_Lzvgxysp7J,smly_CCTQnMll1G,smly_Xa39K0H0cI,smly_319lGhrNm2,smly_JfvmG0aulh,smly_B2tpObhZZE,smly_mloyhGwaMd,smly_z1AH79eou5,smly_bVwVTs46Q7,smly_flt9M3TZNU,smly_wWmMSS6wkj,smly_Zcn0073jyn,smly_0JTjQARKFm,smly_JegxD4qnn5,smly_lWxxstLf2O,smly_Tbi7ep3R6B,smly_FUIh1huTYo,smly_NXvXzobJ5O,smly_6G3rqoqnF8,smly_iRKUmTV3uV,smly_85IRCePdLR,smly_ZITIye3Abo,smly_NtsKuFkGua,smly_3Vtyt17YLE,smly_6CXIH7ltrP,smly_WIH8JEVAqO,smly_L7odEj68uy,smly_G1ITHUfzbr,smly_uRzLuJ3L3t,smly_8E0AFKFYza,smly_efNwZqAHwj,smly_1N8CDiyG4O,smly_8351jupBRC,smly_rqMnsJ829q,smly_ccMsg84jMH,smly_qEqrSaMcTb,smly_Kb0tuNFDUj,smly_v6O3AbCDH6,smly_r5ukHKii8q,smly_E9ztFhHJzj,smly_SyxJvCSosf,smly_18Bq3ACmoe,smly_O37TWprtWx,smly_A8BsKTywSJ,smly_HfRpSKXsG8,smly_nA06DPurV4,smly_BG977WbJUI,smly_3PYzcI6v5M,smly_gdXZ8E6RJ8,smly_fubov6b09p,smly_Tx4SXsBg9B,smly_MhCZknycpk,smly_5AXcOGxUNh,smly_S7Lki1uBs0,smly_iJxIamocpo,smly_qYuX1e2qGV,smly_g3I4GJzORW,smly_8LNn6rcYYf,smly_Pml6Y41GuI,smly_mCDHoEZoDC,smly_69sujXykj2,smly_CHjcEYn36W,smly_dtVusHXhTX,smly_eZYVpPThWj,smly_UmOe0AMkk3,smly_dNWCu2FK6b,smly_1bzM3QD1Bc,smly_0YIUjLwHSw,smly_oNsC3QOh3o,smly_4LDMRGtH6u,smly_ktVuRZpQef,smly_raoltJkUqN,smly_qgAanTdf4A,smly_7BKkJCBKZE,smly_pAyEa1L49j,smly_Pp5l2ZjRPM,smly_gNwv8yhjWR,smly_2DwwM8tkmY,smly_7qJ2VGIpGK,smly_OVNtLELqez,smly_NGfebJ2oiY,smly_PzOJ2k1DYc,smly_tzWWuuCS42,smly_a7bijtKgRc,smly_6QDWHBXb1L,smly_Gql64rdtgS,smly_wHxE91z2n8,smly_kMBpVlavYD,smly_E20zhrygXo,smly_XpBzpGw9RX,smly_hlopJh8Wi2,smly_VLBeEF2ntH,smly_l0dPnG2X2Q,smly_OEu0DSI1DF,smly_T6AQAmCzbS,smly_J0NG6jp0FZ,smly_mU3gzYhkGG,smly_WwaS0PPzW9,smly_gXt6ndQ5mC,smly_3syTRX5zHG,smly_Re4wOCsvkh,smly_h4E8H76GCZ,smly_udiC8uRpkB,smly_Gq6SFhSM73,smly_hKxRpuU8lE,smly_wTgMkjD9HO,smly_QjnKWjPZSl,smly_lFolX16ul5,smly_70w67kPN1L,smly_DRZ28tzJsT,smly_3v3oymtQgv,smly_689M4O0wyz,smly_2g8ALzsWNd,smly_pHG2UAXhi8,smly_jHwcUpnqN1,smly_sm2cZ7gPMV,smly_CXOsM3iw8Y,smly_yUDQcweKsI,smly_SlV6FYGpCx,smly_hb6JbwSzSl,smly_zW3Ik1FEdX,smly_OGQyQM8GRI,smly_ubuEOIb3uR,smly_5jjJr3v1DN,smly_zawhC2qEfy,smly_U9cp6bJH0A

func launch_thing_test(gw gosdk.Gateway, test_key string) {

	defer wg.Done()

	t := gosdk.NewThing()
	t.Name = "water_meter"
	t.Key = test_key
	if err := gw.ThingRegister(t); err != nil {
		fmt.Println("Unable to Register: ", err)
	}

	// Handle instructions if using the mqtt or mqtts protocol.
	if err := gw.SetInstructionHandle(handleInstruciton); err == nil {
		fmt.Println("Started handling the instruction.")
	} else {
		fmt.Println("Could not started handling the instruction: ", err)
	}

	// Send Alerts
	sendAlerts(gw, t)
	count := 0
	for i := 0; i < 1000; i++ {
		if count == 0 {
			// Send a heartbeat.
			gw.ThingHeartbeat(t, 0)
		}
		sendEvent(gw, t)
		time.Sleep(time.Duration(1 * time.Second))
		count++
		count = count % 3
	}

	fmt.Println("Waiting...")
	time.Sleep(time.Duration(1 * time.Second))
	fmt.Println("Bye Bye")

}
