package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

func generateCalendarLink(date time.Time, title string) (string, error) {
	endDate := date.Add(time.Hour * 1)
	formattedStartDate := date.Format("20060102T150405") + "UTC+2"
	formattedEndDate := endDate.Format("20060102T150405") + "UTC+2"

	encodedTitle := strings.ReplaceAll(title, " ", "+")

	linkTemplateString := "http://www.google.com/calendar/event?action=TEMPLATE&text=%s&dates=%s/%sUTC+02&details=&location=&trp=false"

	return fmt.Sprintf(linkTemplateString, encodedTitle, formattedStartDate, formattedEndDate), nil
}

type MeetingInfo struct {
	Date  time.Time `json:"date"`
	Title string    `json:"title"`
}

func parseMeetingInfo(result string) (MeetingInfo, error) {
	var meetingInfo MeetingInfo

	err := json.Unmarshal([]byte(result), &meetingInfo)
	return meetingInfo, err
}

func handleError(message string, err error, c tele.Context) error {
	str := fmt.Sprintf("%s: %v", message, err)
	c.Reply(str)
	return fmt.Errorf(str)
}
