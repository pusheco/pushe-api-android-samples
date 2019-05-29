package main

/* get user pushe_id with this function in your application:
  ** Android **
// java code //
String pid = Pushe.getPusheId(this);
*/

/* get user device_id with this function in your website:
   ** Website **
   See:  https://pushe.co/docs/webpush/#unique_id

   Pushe.getDeviceId()
       .then((deviceId) => {
           console.log(`deviceId is: ${deviceId}`);
       });
*/

// send simple notification to 'YOUR_APPLICATION_ID'

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	// obtain token -> https://pushe.co/docs/api/#api_get_token
	const token = "2746e8cf5ebe571670166ed84621a5c15b13bb2a"

	// Documentation :  (android) https://pushe.co/docs/api/#api_send_push_notification_to_single_users
	// 					(web) https://pushe.co/docs/webpush-api/#api_send_push_notification_according_to_device_id
	// ******************************************************
	// **************** filtered by imei ********************
	// ********************* For Android ********************
	// ******************************************************
	IMEIFilteredData := map[string]interface{}{
		"app_ids": []string{"co.pushe.apushe163test"}, // a list of app_id, like: [app_id_1 , ...] (compulsive)
		// send notification to all my android(web) applications
		// "app_ids":  []string{"__all__"}
		"platform": 1, // optional for android,
		"data": map[string]interface{}{
			"title":   "This is a simple push",         // (compulsive)
			"content": "All of your users will see me", // (compulsive)
		},
		"filters": map[string]interface{}{
			"imei": []string{"212**********23"}, // filter user by imei key (deprecated)
		},
		// extra parameters on Documentation -> https://pushe.co/docs/api/#api_send_advance_notification
	}

	// ******************************************************
	// **************** filtered by pushe_id ****************
	// ********************* For Android ********************
	// ******************************************************
	PusheIdFilteredData := map[string]interface{}{
		"app_ids": []string{"co.pushe.apushe163test"}, // a list of app_id, like: [app_id_1 , ...] (compulsive)
		// send notification to all my android(web) applications
		// "app_ids":  []string{"__all__"}
		"platform": 1, // optional for android,
		"data": map[string]interface{}{
			"title":   "This is a simple push",         // (compulsive)
			"content": "All of your users will see me", // (compulsive)
		},
		"filters": map[string]interface{}{
			"pushe_id": []string{"pid_************"}, // filter user by pushe_id key
		},
		// extra parameters on Documentation -> https://pushe.co/docs/api/#api_send_advance_notification
	}
	fmt.Println("PusheIdFilteredData", PusheIdFilteredData)

	// *******************************************************
	// **************** filtered by device_id ****************
	// **************** For Android and Web ******************
	// *******************************************************
	DeviceIdFilteredData := map[string]interface{}{
		"app_ids": []string{"co.pushe.apushe163test"}, // a list of app_id, like: [app_id_1 , ...] (compulsive)
		// send notification to all my android(web) applications
		// "app_ids":  []string{"__all__"}
		"platform": 1, // optional for android,
		// "platform": 2, for web (compulsive for web)
		"data": map[string]interface{}{
			"title":   "This is a simple push",         // (compulsive)
			"content": "All of your users will see me", // (compulsive)
		},
		"filters": map[string]interface{}{
			"device_id": []string{"DEVICE_ID_1"}, // filter user by device_id key
		},
		// extra parameters on Documentation -> https://pushe.co/docs/api/#api_send_advance_notification
	}
	fmt.Println("DeviceIdFilteredData", DeviceIdFilteredData)

	// Marshal returns the JSON encoding of reqData.
	// Choices: IMEIFilteredData, DeviceIdFilteredData, PusheIdFilteredData
	reqJSON, err := json.Marshal(IMEIFilteredData)

	// check encoded json
	if err != nil {
		fmt.Println("json:", err)
		return
	}

	// create request obj
	request, err := http.NewRequest(
		http.MethodPost,
		"https://api.pushe.co/v2/messaging/notifications/",
		bytes.NewBuffer(reqJSON),
	)

	// check request
	if err != nil {
		fmt.Println("Req error:", err)
		return
	}

	// set header
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Token "+token)

	// send request and get response
	client := http.Client{}
	response, err := client.Do(request)

	// check response
	if err != nil {
		fmt.Println("Resp error:", err)
		return
	}

	defer response.Body.Close()

	// check status_code and response
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(response.Body)
	respContent := buf.String()

	fmt.Println("status code =>", response.StatusCode)
	fmt.Println("response =>", respContent)
	fmt.Println("==========")

	if response.StatusCode == http.StatusCreated {
		fmt.Println("success!")

		var respData map[string]interface{}
		_ = json.Unmarshal(buf.Bytes(), &respData)

		var reportURL string

		// hashed_id just generated for Non-Free plan
		if respData["hashed_id"] != nil {
			reportURL = "https://pushe.co/report?id=" + respData["hashed_id"].(string)
		} else {
			reportURL = "no report url for your plan"
		}

		fmt.Println("report_url:", reportURL)
		fmt.Println("notification id:", int(respData["wrapper_id"].(float64)))
	} else {
		fmt.Println("failed!")
	}
}
