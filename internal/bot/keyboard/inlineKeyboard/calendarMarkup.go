package bot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NicoNex/echotron/v3"
)

var weekDays = [7]string{
	"Пн",
	"Вт",
	"Ср",
	"Чт",
	"Пт",
	"Сб",
	"Вс",
}

var monthNamesByNumbers = map[int]string{
	1:  "Январь",
	2:  "Февраль",
	3:  "Март",
	4:  "Апрель",
	5:  "Май",
	6:  "Июнь",
	7:  "Июль",
	8:  "Август",
	9:  "Сентябрь",
	10: "Октябрь",
	11: "Ноябрь",
	12: "Декабрь",
}

type CalendarMarkup struct {
	Month    int
	Year     int
	keyboard [][]echotron.InlineKeyboardButton
}

func New(month, year int) *CalendarMarkup {
	return &CalendarMarkup{
		Month: month,
		Year:  year,
	}
}

func (c *CalendarMarkup) BuildMarkup() echotron.InlineKeyboardMarkup {
	c.title()
	c.daysHeader()
	c.monthDays()
	c.addNavigationButtons()

	return echotron.InlineKeyboardMarkup{
		InlineKeyboard: c.keyboard,
	}
}

func (c *CalendarMarkup) NextMonth() echotron.InlineKeyboardMarkup {
	nextTime := time.Date(c.Year, time.Month(c.Month+2), 0, 0, 0, 0, 0, time.UTC)
	c = New(
		int(nextTime.Month()),
		nextTime.Year(),
	)
	return c.BuildMarkup()
}

func (c *CalendarMarkup) PreviousMonth() echotron.InlineKeyboardMarkup {
	previousTime := time.Date(c.Year, time.Month(c.Month), 0, 0, 0, 0, 0, time.UTC)
	c = New(
		int(previousTime.Month()),
		previousTime.Year(),
	)
	return c.BuildMarkup()
}

func (c *CalendarMarkup) getRangeMonth() (duration int, firstDay int) {
	duration = time.Date(c.Year, time.Month(c.Month+1), 0, 0, 0, 0, 0, time.UTC).Day()
	firstDay = (int(time.Date(c.Year, time.Month(c.Month), 1, 0, 0, 0, 0, time.UTC).Weekday()) + 6) % 7
	return
}

func (c *CalendarMarkup) title() {
	c.keyboard = append(
		c.keyboard,
		addButton(fmt.Sprintf("%s %v", c.monthName(), c.Year), "nil"),
	)
}

func (c *CalendarMarkup) daysHeader() {
	c.keyboard = append(c.keyboard, rowButtons(weekDays[:]))
}

func (c *CalendarMarkup) monthDays() {
	duration, firstDay := c.getRangeMonth()
	row := c.addEmptyButtons(firstDay)
	for i := 1; i <= duration; i++ {
		date := fmt.Sprintf("%.2v.%.2v.%.4v", i, c.Month, c.Year)
		if len(row) == 7 {
			c.keyboard = append(c.keyboard, row)
			row = nil
		}
		row = append(row, addButton(strconv.Itoa(i), date)...)
	}
	emptyButtons := 7 - len(row)
	row = append(row, c.addEmptyButtons(emptyButtons)...)
	c.keyboard = append(c.keyboard, row)
}

func (c *CalendarMarkup) addNavigationButtons() {
	date := fmt.Sprintf("%.2v.%.4v", c.Month, c.Year)
	navigationButtons := []echotron.InlineKeyboardButton{
		addButton("<", fmt.Sprintf("back %s", date))[0],
		addButton("Меню", "menu")[0],
		addButton(">", fmt.Sprintf("next %s", date))[0],
	}
	c.keyboard = append(c.keyboard, navigationButtons)
}

func (c *CalendarMarkup) monthName() string {
	return monthNamesByNumbers[c.Month]
}

func (c *CalendarMarkup) addEmptyButtons(emptyButtons int) (emptyInlineButtons []echotron.InlineKeyboardButton) {
	for i := 0; i < emptyButtons; i++ {
		emptyInlineButtons = append(
			emptyInlineButtons,
			addButton(" ", "nil")...,
		)
	}
	return
}
