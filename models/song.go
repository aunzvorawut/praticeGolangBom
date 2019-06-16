package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Song struct {
	Id          int64     `orm:"pk;auto"`
	Version     int64     `orm:"default(0)"`
	DateCreated time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`

	Link              string
	NoVocalLink       string
	SdLink            string
	SdNoVocalLink     string
	AutoLink          string
	AutoNoVocalLink   string
	P240link          string
	P240noVocalLink   string
	P144sdLink        string
	P144sdNoVocalLink string
	P480link          string
	P480noVocalLink   string
	P1080link         string
	P1080noVocalLink  string

	IsRecordable            bool `orm:"default(false)"`
	IsRecordableTeamAisTest bool `orm:"default(false)"`
	IsPurchasedTest         bool
	SongName                string
	SongNameEng             string
	AlbumName               string
	AlbumNameEng            string
	ArtistName              string
	ArtistNameEng           string
	Duration                string
	Company                 string
	CoverPath               string
	SongCode                string
	ReleaseDate             string
	AlbumCode               string `orm:"null;size(255)"`
	ReleaseYear             string
	CategoryStr             string
	Enabled                 bool `orm:"default(true)"`
	IsFree                  bool `orm:"default(false)"`
	Checkm3u8               bool `orm:"default(false)"`
	Completem3u8            bool `orm:"default(true)"`
	IsBox                   bool
	IsMobile                bool
	IsForShared             bool `orm:"default(false)"`
	CountRecord             int64
	CountView               int64 `orm:"default(0)"`

	StartDate   time.Time
	ExpiredDate time.Time
}

func init() {
	orm.RegisterModel(new(Song))
}

// AddSong insert a new Song into database and returns
// last inserted Id on success.
func AddSong(m *Song) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetSongById retrieves Song by Id. Returns error if
// Id doesn't exist
func GetSongById(id int64) (v *Song, err error) {
	o := orm.NewOrm()
	v = &Song{Id: id}
	if err = o.QueryTable(new(Song)).Filter("Id", id).RelatedSel().One(v); err == nil {

		return v, nil
	}
	return nil, err
}

// GetAllSong retrieves all Song matches certain condition. Returns empty list if
// no records exist
func GetAllSong(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Song))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Song
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateSong updates Song by Id and returns error if
// the record to be updated doesn't exist
func UpdateSongById(m *Song) (err error) {
	o := orm.NewOrm()
	v := Song{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSong deletes Song by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSong(id int64) (err error) {
	o := orm.NewOrm()
	v := Song{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Song{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAllRecentlySongIdByUserObj(max, offset int64, userObj *SecUser, clientTypeFunc string) (result []int64) {

	var junkResult []string
	today := time.Now()
	todayStr := today.Format("2006-01-02 15:04:05")
	extraSql := ""

	if clientTypeFunc == "box" {
		extraSql = " and s.is_box = true "
	}

	if clientTypeFunc == "mobile" {
		extraSql = " and s.is_mobile = true "
	}

	stringQuery := " SELECT song_id as songId , MAX(sps.date_created) as date " +
		" FROM `stats_play_song`  sps join song s on sps.song_id = s.id " +
		" where sps.song_id is not null and sps.sec_user_id = ? " +
		extraSql +
		" and s.start_date <= ? and s.expired_date >= ?  " +
		" GROUP BY sps.song_id " +
		" order by date desc  " +
		" limit ? , ? "

	o := orm.NewOrm()
	_, err := o.Raw(stringQuery, userObj.Id, todayStr, todayStr, offset, max).QueryRows(&result , &junkResult)
	if err != nil {
		beego.Error("stringQuery = ",stringQuery)
		beego.Error(err.Error())
	}
	return result

}
