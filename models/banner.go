package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Banner struct {
	Id          int64     `orm:"pk;auto"`
	Version     int64     `orm:"default(0)"`
	DateCreated time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`

	BannerCoverImage   string    `orm:"null"`
	NextTo             string    `orm:"null"`
	ContentUrl         string    `orm:"null"`
	OrderPosition      int64     `orm:"default(0)"`
	DescriptionClient  string    `orm:"null"`
	Enabled            bool      `orm:"null;default(true)"`
	SongObj            *Song     `orm:"rel(fk);null"`
	CategoryObj        *Category `orm:"rel(fk);null"`
	MemoDetail         string    `orm:"null"`
	IsBox              bool
	IsMobile           bool
	IsVersionOne       bool
	IsVersionTwo       bool
	ShowOnPage         string `orm:"null;default(normal)"`
	CheckPurchasedUser bool   `orm:"default(false)"`
}

func init() {
	orm.RegisterModel(new(Banner))
}

// AddBanner insert a new Banner into database and returns
// last inserted Id on success.
func AddBanner(m *Banner) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetBannerById retrieves Banner by Id. Returns error if
// Id doesn't exist
func GetBannerById(id int64) (v *Banner, err error) {
	o := orm.NewOrm()
	v = &Banner{Id: id}
	if err = o.QueryTable(new(Banner)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBanner retrieves all Banner matches certain condition. Returns empty list if
// no records exist
func GetAllBanner(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Banner))
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

	var l []Banner
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

// UpdateBanner updates Banner by Id and returns error if
// the record to be updated doesn't exist
func UpdateBannerById(m *Banner) (err error) {
	o := orm.NewOrm()
	v := Banner{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBanner deletes Banner by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBanner(id int64) (err error) {
	o := orm.NewOrm()
	v := Banner{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Banner{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAllBannerHomePage(max, offset int64, clientTypeFunc string) (ml []*Banner) {

	o := orm.NewOrm()
	qs := o.QueryTable(new(Banner))

	qs = qs.Limit(max, offset).Filter("IsVersionTwo", true).Filter("Enabled", true)

	if clientTypeFunc == "box" {
		qs = qs.Filter("IsBox", true)
	}

	if clientTypeFunc == "mobile" {
		qs = qs.Filter("IsMobile", true)
	}

	qs = qs.OrderBy("OrderPosition")

	_, err := qs.RelatedSel().All(&ml)
	if err == nil {
		return ml
	}

	return nil

}
