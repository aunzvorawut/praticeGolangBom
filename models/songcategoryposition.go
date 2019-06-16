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

type SongCategoryPosition struct {
	Id            int64     `orm:"pk;auto"`
	Version       int64     `orm:"default(0)"`
	DateCreated   time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated   time.Time `orm:"auto_now;type(datetime)"`
	Song          *Song     `orm:"rel(fk)"`
	Category      *Category `orm:"rel(fk)"`
	OrderPosition int64     `orm:"default(0)"`
}

func init() {
	orm.RegisterModel(new(SongCategoryPosition))
}

// AddSongCategoryPosition insert a new SongCategoryPosition into database and returns
// last inserted Id on success.
func AddSongCategoryPosition(m *SongCategoryPosition) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetSongCategoryPositionById retrieves SongCategoryPosition by Id. Returns error if
// Id doesn't exist
func GetSongCategoryPositionById(id int64) (v *SongCategoryPosition, err error) {
	o := orm.NewOrm()
	v = &SongCategoryPosition{Id: id}
	if err = o.QueryTable(new(SongCategoryPosition)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSongCategoryPosition retrieves all SongCategoryPosition matches certain condition. Returns empty list if
// no records exist
func GetAllSongCategoryPosition(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SongCategoryPosition))
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

	var l []SongCategoryPosition
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

// UpdateSongCategoryPosition updates SongCategoryPosition by Id and returns error if
// the record to be updated doesn't exist
func UpdateSongCategoryPositionById(m *SongCategoryPosition) (err error) {
	o := orm.NewOrm()
	v := SongCategoryPosition{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSongCategoryPosition deletes SongCategoryPosition by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSongCategoryPosition(id int64) (err error) {
	o := orm.NewOrm()
	v := SongCategoryPosition{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SongCategoryPosition{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAllSongIdByCatId(max, offset int, catId int64, clientTypeFunc string) (result []int64) {

	today := time.Now()
	todayStr := today.Format("2006-01-02 15:04:05")
	extraSql := ""

	if clientTypeFunc == "box" {
		extraSql = " and s.is_box = true "
	}

	if clientTypeFunc == "mobile" {
		extraSql = " and s.is_mobile = true "
	}

	stringQuery := " select scp.song_id " +
		" from song_category_position scp " +
		" join song s " +
		" on scp.song_id = s.id " +
		" where scp.category_id = ? " +
		" and s.enabled = true " +
		" and s.start_date <= ? " +
		" and s.expired_date >= ? " +
		extraSql +
		" order by scp.order_position " +
		" limit ? , ? "

	o := orm.NewOrm()
	_, err := o.Raw(stringQuery, catId, todayStr, todayStr, offset, max).QueryRows(&result)
	if err != nil {
		beego.Error(err.Error())
	}
	return result

}

func GetCountSongIdByCatId(catId int64, clientTypeFunc string) (result []int64) {

	today := time.Now()
	todayStr := today.Format("2006-01-02 15:04:05")
	extraSql := ""

	if clientTypeFunc == "box" {
		extraSql = " and s.is_box = true "
	}

	if clientTypeFunc == "mobile" {
		extraSql = " and s.is_mobile = true "
	}

	stringQuery := " select count(scp.song_id) " +
		" from song_category_position scp " +
		" join song s " +
		" on scp.song_id = s.id " +
		" where scp.category_id = ? " +
		" and s.enabled = true " +
		" and s.start_date <= ? " +
		" and s.expired_date >= ? " +
		extraSql +
		" order by scp.order_position "

	o := orm.NewOrm()
	_, err := o.Raw(stringQuery, catId, todayStr, todayStr).QueryRows(&result)
	if err != nil {
		beego.Error(err.Error())
	}
	return result

}

func GetAllSongObjByCatIdAndClientType(catId int64, clientTypeFunc string) (ml []*SongCategoryPosition) {

	today := time.Now()

	o := orm.NewOrm()
	qs := o.QueryTable(new(SongCategoryPosition))

	qs = qs.Filter("Category", catId).Filter("song__enabled", true).Filter("song__enabled", true).Filter("song__startDate__lt", today).Filter("song__expiredDate__gt", today)

	if clientTypeFunc == "box" {
		qs = qs.Filter("song__isBox", true)
	}

	if clientTypeFunc == "mobile" {
		qs = qs.Filter("song__isMobile", true)
	}
	_, err := qs.OrderBy("orderPosition").RelatedSel().All(&ml)

	if err == nil {
		return ml
	}

	return nil

}

func GetSongCategoryPositionBySongIdAndCatId(songId, catId int64) *SongCategoryPosition {

	o := orm.NewOrm()
	d := new(SongCategoryPosition)
	err := o.QueryTable(new(SongCategoryPosition)).Filter("Song", songId).Filter("Category", catId).RelatedSel().Limit(1).One(d)
	if err != nil {
		return nil
	}

	return d
}
