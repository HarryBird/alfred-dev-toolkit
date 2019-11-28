package toolkit

import (
	"log"
	"strconv"
	"strings"
	"time"

	"regexp"

	alfred "github.com/HarryBird/alfred-toolkit-go"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func TimeAction(ctx *cli.Context, al *alfred.Alfred) {
	args := []string(ctx.Args())
	log.Println("Args:", args)

	t := time.Now()

	if len(args) > 0 {
		it := args[0]
		log.Println("Time:", it)

		log.Println(regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`).MatchString(it))

		var err error

		switch {
		// format: 2016-09-20
		case regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`).MatchString(it):
			t, err = time.ParseInLocation("2006-01-02", it, t.Location())
		// format: 2016/09/20
		case regexp.MustCompile(`^[0-9]{4}/[0-9]{2}/[0-9]{2}$`).MatchString(it):
			t, err = time.ParseInLocation("2006/01/02", it, t.Location())
		// format: 2016-09-20 14:30:30
		case regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}\s+[0-9]{2}:[0-9]{2}:[0-9]{2}$`).MatchString(it):
			t, err = time.ParseInLocation("2006-01-02 15:04:05", it, t.Location())
		// format: 2016/09/20 14:30:30
		case regexp.MustCompile(`^[0-9]{4}/[0-9]{2}/[0-9]{2}\s+[0-9]{2}:[0-9]{2}:[0-9]{2}$`).MatchString(it):
			t, err = time.ParseInLocation("2006/01/02 15:04:05", it, t.Location())
		// format: 1571549947.343234
		//case regexp.MustCompile(`^[0-9]{1,10}\.*[0-9]*`).MatchString(it):
		case regexp.MustCompile(`^[1-9][.\d]*`).MatchString(it):
			parts := strings.Split(it, ".")
			l := len(parts)
			var sec, nsec int64

			switch l {
			case 1:
				sec, err = strconv.ParseInt(parts[0], 10, 64)
				t = time.Unix(sec, 0)
			case 2:
				var err1, err2 error
				sec, err1 = strconv.ParseInt(parts[0], 10, 64)
				nsec, err2 = strconv.ParseInt(parts[1], 10, 64)

				if err1 != nil {
					err = err1
				} else {
					if err2 != nil {
						err = err2
						t = time.Unix(sec, 0)
					} else {
						t = time.Unix(sec, nsec)
					}
				}
			}
		default:
			err = errors.New("UnSupport Time Format")
		}

		if err != nil {
			log.Println("Parse Time Fail:" + err.Error())
			al.ResultAppend(alfred.NewErrorTitleItem("Parse Time Fail, Invalid Format", err.Error())).Output()
			return
		}
	}

	ts := strconv.FormatInt(t.Unix(), 10) + "." + strconv.Itoa(t.Nanosecond())
	local := t.Format("2006-01-02 15:04:05.999999999 -0700 MST")
	weekday := t.Weekday().String()
	utc := t.UTC().Format("2006-01-02 15:04:05.999999999 -0700 MST")

	al.ResultAppend(buildDateItem("Stamp: "+ts, "", ts))

	al.ResultAppend(buildDateItem("Day: "+weekday, "", weekday))

	al.ResultAppend(buildDateItem("Local: "+local, "", local))

	al.ResultAppend(buildDateItem("UTC: "+utc, "", utc))

	al.Output()
}

func buildDateItem(title, subTitle, arg string) alfred.Item {
	return alfred.NewItem(title, subTitle, arg, arg, "", "", true, alfred.NewIcon("", "./icons/date/clock.png"))
}
