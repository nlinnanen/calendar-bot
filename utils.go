package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

func generateCalendarLink(date time.Time, title string) (string, error) {
	end_date := date.Add(time.Hour * 1)
	formatted_start_date := date.Format("20060102T150405") + "UTC+2"
	formatted_end_date := end_date.Format("20060102T150405") + "UTC+2"

	encoded_title := strings.ReplaceAll(title, " ", "+")

	link_template_string := "http://www.google.com/calendar/event?action=TEMPLATE&text=%s&dates=%s/%sUTC+02&details=&location=&trp=false"

	return fmt.Sprintf(link_template_string, encoded_title, formatted_start_date, formatted_end_date), nil
}

type MeetingInfo struct {
	Date  time.Time `json:"date"`
	Title string    `json:"title"`
}

func parseMeetingInfo(result string) (MeetingInfo, error) {
	var meeting_info MeetingInfo

	err := json.Unmarshal([]byte(result), &meeting_info)
	return meeting_info, err
}

func handleError(message string, err error, c tele.Context) error {
	str := fmt.Sprintf("%s: %v", message, err)
	c.Reply(str)
	return fmt.Errorf(str)
}
