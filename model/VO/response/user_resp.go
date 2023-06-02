package response

import "dou_yin/service/DO"

type QueryContactorList struct {
	ContactorList DO.ContactList `json:"contactor_list"`
}
