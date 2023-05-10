package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tg-bot/pkg/types"
	"time"
)

func UserExistsInBx(id int, apiEndpoint string) (bool, bool, error) {
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return false, false, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "bot-app")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, false, err
	}
	defer resp.Body.Close()

	var data types.BxResult
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false, false, err
	}

	for _, userInfo := range data.Result {
		bxId, err := strconv.ParseInt(userInfo.Id, 10, 32)
		if err != nil {
			return false, false, err
		}
		if int(bxId) == id {
			return true, userInfo.Active, nil
		}
	}

	return false, false, nil
}

func CheckEmployee(usersIdsChan chan []int, errChan chan error, apiEndpoint string) {
	alreadyCheck := false

	start := 0
	limit := 50

	for {
		usersIDs := make([]int, 0)
		now := time.Now().In(time.FixedZone("MSK", 3*60*60))

		if now.Hour() == 16 && now.Minute() == 13 && !alreadyCheck {
			var data types.BxResult

			for {
				var currData types.BxResult
				endpoint := fmt.Sprintf("%s?start=%d", apiEndpoint, start)
				req, err := http.NewRequest("GET", endpoint, nil)

				if err != nil {
					errChan <- err
					return
				}

				req.Header.Set("Accept", "application/json")
				req.Header.Set("User-Agent", "bot-app")

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					errChan <- err
					continue
				}
				defer resp.Body.Close()

				err = json.NewDecoder(resp.Body).Decode(&currData)
				if err != nil {
					errChan <- err
					continue
				}

				data.Result = append(data.Result, currData.Result...)
				if len(currData.Result) < limit {
					break
				}
				start += 50
			}

			for _, userInfo := range data.Result {
				if !userInfo.Active {
					bxId, err := strconv.ParseInt(userInfo.Id, 10, 32)
					if err != nil {
						errChan <- err
						return
					}
					if bxId < 9 { // TODO: REMOVE CONDITION
						usersIDs = append(usersIDs, int(bxId))
					}
				}
			}
			alreadyCheck = true
			usersIdsChan <- usersIDs
		} else {
			alreadyCheck = false
		}

		time.Sleep(1 * time.Minute)
	}
}
